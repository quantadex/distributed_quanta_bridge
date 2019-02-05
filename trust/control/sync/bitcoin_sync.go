package sync

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"encoding/json"
	"fmt"
)

type BitcoinSync struct {
	DepositSync
}

func (c *BitcoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
}

func (c *BitcoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, nil)

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