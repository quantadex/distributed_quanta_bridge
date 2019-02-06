package sync

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"strings"
)

type BitcoinSync struct {
	DepositSync
}

func (c *BitcoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
}

func (c *BitcoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := db.GetCrosschainByBlockchain(c.rDb, coin.BLOCKCHAIN_BTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[strings.ToLower(w.Address)] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	if err != nil {
		c.logger.Info("getDepositsInBlock failed " + err.Error())
		return nil, err
	}

	if len(deposits) > 0 {
		msg,_ := json.Marshal(deposits)
		fmt.Printf("events = %v\n", string(msg))
	}

	return deposits, err
}

func (c *BitcoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}