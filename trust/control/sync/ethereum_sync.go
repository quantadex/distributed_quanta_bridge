package sync

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"strings"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type EthereumSync struct {
	DepositSync
	trustAddress common.Address
	issuingSymbol string
}

func (c *EthereumSync) Setup() {
	c.fnDepositInBlock =  c.GetDepositsInBlock
	c.fnPostProcessBlock = c.PostProcessBlock
}
/**
 * getDepositsInBlock
 *
 * Returns deposits made to the coin trust account in this block
 */

func (c *EthereumSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := db.GetCrosschainByBlockchain(c.rDb, coin.BLOCKCHAIN_ETH)
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
		// define custom token issuance
		if dep.CoinName == "ETH" {
			dep.CoinName = c.issuingSymbol

			// ethereum converts to precision 5, now we need to convert to precision of the asset
			dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 5, int(c.coinInfo.Precision))
		} else {
			// we assume precision is always 5
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
			c.logger.Infof("New Forwarder Address ETH->QUANTA address, %s -> %s", addr.ContractAddress.Hex(), addr.QuantaAddr)
			db.AddCrosschainAddress(c.rDb, addr)
		} else {
			c.logger.Error(fmt.Sprintf("MISMATCH: Forwarder address[%s] in blockID=%d does not match our trustAddress[%s]",
				addr.Trust.Hex(), blockID, c.trustAddress.Hex()))
		}
	}
	return nil
}