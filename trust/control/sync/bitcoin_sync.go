package sync

import (
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
)

type BitcoinSync struct {
	DepositSync
	issuingSymbol map[string]string //TODO: pass in only neccessary data (eg. issuingSymbol)
}

func (c *BitcoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
}

func (c *BitcoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_BTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	for _, dep := range deposits {
		if dep.CoinName == "BTC" {
			dep.CoinName = c.issuingSymbol["btc"]
		}
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

func (c *BitcoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}