package sync

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/node/webhook"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"math/big"
	"strconv"
)

type BitcoinSync struct {
	DepositSync
	issuingSymbol map[string]string //TODO: pass in only neccessary data (eg. issuingSymbol)
	btcMinConfirm int64
}

func (c *BitcoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
	c.fnGetWatchAddress = c.GetWatchAddress
	c.fnTransformCoin = c.TransformCoin
	c.fnFindAllAndConfirm = c.FindAllAndConfirm
	c.fnGetMinConfirmation = c.GetMinConfirmation
}

func (c *BitcoinSync) TransformCoin(dep *coin.Deposit) *coin.Deposit {
	if dep.CoinName == "BTC" {
		dep.CoinName = c.issuingSymbol["btc"]
	}
	dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 8, int(c.coinInfo[c.issuingSymbol["btc"]].Precision))
	return dep
}

func (c *BitcoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_BTC)
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

func (c *BitcoinSync) GetWatchAddress() map[string]string {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_BTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	return watchMap
}

func (c *BitcoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}

func (c *BitcoinSync) FindAndConfirm(tx db.Transaction, blockHash string, confirmations int64) error {
	if confirmations == -1 {
		err := db.ChangeSubmitState(c.rDb, tx.Tx, db.ORPHAN, db.DEPOSIT, blockHash)
		if err != nil {
			return errors.Wrap(err, "Could not change state to orphan")
		}
	} else if confirmations > c.btcMinConfirm {
		c.eventsChan <- webhook.Event{control.Deposit_In_Consensus, tx.To, tx.Tx}

		err := db.ChangeSubmitState(c.rDb, tx.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT, blockHash)
		if err != nil {
			return errors.Wrap(err, "Could not change state to consensus")
		}
	} else {
		c.eventsChan <- webhook.Event{control.Deposit_Wait_For_Confirmation, tx.To, tx.Tx}

		submitState := db.WAIT_FOR_CONFIRMATION + " " + strconv.Itoa(int(confirmations)) + "/" + strconv.Itoa(int(c.btcMinConfirm))
		err := db.ChangeSubmitState(c.rDb, tx.Tx, submitState, db.DEPOSIT, blockHash)
		if err != nil {
			return errors.Wrap(err, "Could not change state to wait for confirmation")
		}
		c.logger.Infof("Transaction %s has %d confirmations", tx.Tx, confirmations)
	}
	return nil
}

func (c *BitcoinSync) FindAllAndConfirm() error {
	txs, err := db.QueryAllWaitForConfirmTx(c.rDb, c.issuingSymbol["btc"])
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

func (c *BitcoinSync) GetMinConfirmation() int64 {
	return c.btcMinConfirm
}
