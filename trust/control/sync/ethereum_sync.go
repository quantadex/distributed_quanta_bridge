package sync

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"math/big"
	"strconv"
	"strings"
)

type EthereumSync struct {
	DepositSync
	trustAddress  common.Address
	issuingSymbol map[string]string
	ethFlush      bool
	ethMinConfirm int64
}

func (c *EthereumSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
	c.fnPostProcessBlock = c.PostProcessBlock
	c.fnGetWatchAddress = c.GetWatchAddress
	c.fnFindAllAndConfirm = c.FindAllAndConfirm
	c.fnGetMinConfirmation = c.GetMinConfirmation
}

/**
 * getDepositsInBlock
 *
 * Returns deposits made to the coin trust account in this block
 */

func (c *EthereumSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_ETH)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[strings.ToLower(w.Address)] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	if err != nil {
		c.logger.Info("getDepositsInBlock failed " + err.Error())
		return nil, err
	}

	for _, dep := range deposits {

		err := c.rDb.UpdateCrosschainAddrBlockNumber(dep.SenderAddr, uint64(blockID))
		if err != nil {
			c.logger.Errorf("Could not update the last block number for %s", dep.SenderAddr)
		}

		// define custom token issuance
		if dep.CoinName == "ETH" {
			dep.CoinName = c.issuingSymbol["eth"]

			// ethereum converts to precision 5, now we need to convert to precision of the asset
			dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 5, int(c.coinInfo[c.issuingSymbol["eth"]].Precision))
		} else {
			dep.CoinName = strings.ToUpper(dep.CoinName)
			// we assume precision is always 5
			// we need to flush erc-20 coins
			if strings.Contains(dep.CoinName, "0X") {

				if c.ethFlush {
					parts := strings.Split(dep.CoinName, "0X")
					ercAddr := "0x" + parts[1]

					err = c.coinChannel.FlushCoin(dep.SenderAddr, ercAddr)
					if err != nil {
						c.logger.Error(err.Error())
					}
				}

				asset, err := c.quantaChannel.GetAsset(dep.CoinName)
				var precision int
				if err != nil {
					precision = 5
				} else {
					precision = int(asset.Precision)
				}

				dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 5, precision)
			}
		}

		// Need to convert to uppercase, which graphene requires
		dep.CoinName = strings.ToUpper(dep.CoinName)
	}

	return deposits, nil
}

func (c *EthereumSync) PostProcessBlock(blockID int64) error {
	addresses, err := c.coinChannel.GetForwardersInBlock(blockID)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	for _, addr := range addresses {
		if addr.Trust.Hex() == c.trustAddress.Hex() {
			c.logger.Infof("New Forwarder Address ETH->QUANTA address, %s -> %s", addr.ContractAddress, addr.QuantaAddr)
			c.rDb.AddCrosschainAddress(addr)
			//db.AddCrosschainAddress(c.rDb, addr)
		} else {
			c.logger.Error(fmt.Sprintf("MISMATCH: Forwarder address[%s] in blockID=%d does not match our trustAddress[%s]",
				addr.Trust.Hex(), blockID, c.trustAddress.Hex()))
		}
	}
	return nil
}

func (c *EthereumSync) GetWatchAddress() map[string]string {
	return nil
}

func (c *EthereumSync) FindAllAndConfirm() error {
	txs, err := db.QueryAllWaitForConfirmTxETH(c.rDb, c.issuingSymbol["eth"])
	if err != nil {
		return err
	}

	for _, t := range txs {
		blockHash, confirm, err := c.coinChannel.GetBlockInfo(t.Tx)
		if err != nil {
			return errors.Wrap(err, "Could not get the transaction")
		}
		tx, err := db.GetAllWaitForConfirmTransaction(c.rDb, t.Tx, db.DEPOSIT)
		for _, d := range tx {
			if d.BlockHash == blockHash {
				if confirm > c.ethMinConfirm {
					err := db.ChangeSubmitState(c.rDb, d.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT, d.BlockHash)
					if err != nil {
						return errors.Wrap(err, "Could not change the submit state to consensus")
					}
				} else {
					submitState := db.WAIT_FOR_CONFIRMATION + " " + strconv.Itoa(int(confirm)) + "/" + strconv.Itoa(int(c.ethMinConfirm))
					err := db.ChangeSubmitState(c.rDb, d.Tx, submitState, db.DEPOSIT, blockHash)
					if err != nil {
						return errors.Wrap(err, "Could not change state to wait for confirmation")
					}
					c.logger.Infof("Transaction %s has %d confirmations", d.Tx, confirm)
				}
			} else {
				c.logger.Infof("BlockHash is different for %s, it is now an orphan", d.Tx)
				err := db.ChangeSubmitState(c.rDb, d.Tx, db.ORPHAN, db.DEPOSIT, d.BlockHash)
				if err != nil {
					return errors.Wrap(err, "Could not change the submit state to orphan")
				}

			}
		}
	}
	return nil
}

func (c *EthereumSync) GetMinConfirmation() int64 {
	return c.ethMinConfirm
}
