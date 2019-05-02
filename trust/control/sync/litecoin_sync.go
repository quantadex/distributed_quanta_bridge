package sync

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"math/big"
)

type LitecoinSync struct {
	DepositSync
	issuingSymbol map[string]string
	ltcMinConfirm int64
}

func (c *LitecoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
	c.fnGetWatchAddress = c.GetWatchAddress
	c.fnTransformCoin = c.TransformCoin
	c.fnFindAllAndConfirm = c.FindAllAndConfirm
}

func (c *LitecoinSync) TransformCoin(dep *coin.Deposit) *coin.Deposit {
	if dep.CoinName == "LTC" {
		dep.CoinName = c.issuingSymbol["ltc"]
	}
	dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 8, int(c.coinInfo[c.issuingSymbol["ltc"]].Precision))
	return dep
}

func (c *LitecoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_LTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	for _, dep := range deposits {
		dep = c.TransformCoin(dep)
	}

	if err != nil {
		c.logger.Info("getDepositsInBlock failed " + err.Error())
		return nil, err
	}

	if len(deposits) > 0 {
		msg, _ := json.Marshal(deposits)
		fmt.Printf("events = %v\n", string(msg))
	}

	return deposits, err
}

func (c *LitecoinSync) GetWatchAddress() map[string]string {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_LTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	return watchMap
}

func (c *LitecoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}

func (c *LitecoinSync) FindAndConfirm(tx db.Transaction, blockHash string, confirmations int64) error {
	if confirmations == -1 {
		err := db.ChangeSubmitState(c.rDb, tx.Tx, db.ORPHAN, db.DEPOSIT, blockHash)
		if err != nil {
			return errors.Wrap(err, "Could not change state to orphan")
		}
	} else if confirmations > c.ltcMinConfirm {
		err := db.ChangeSubmitState(c.rDb, tx.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT, blockHash)
		if err != nil {
			return errors.Wrap(err, "Could not change state to consensus")
		}
	} else {
		c.logger.Infof("Transaction %s has %d confirmations", tx.Tx, confirmations)
	}
	return nil
}

func (c *LitecoinSync) FindAllAndConfirm() error {
	txs, err := db.QueryAllWaitForConfirmTx(c.rDb, c.issuingSymbol["ltc"])
	if err != nil {
		return err
	}
	for _, tx := range txs {
		blockHash, confirmations, err := c.coinChannel.GetBlockInfo(tx.BlockHash)
		if err != nil {
			return errors.Wrap(err, "Could not get block info")
		}
		err = c.FindAndConfirm(tx, blockHash, confirmations)
		if err != nil {
			return err
		}
	}
	return nil
}
