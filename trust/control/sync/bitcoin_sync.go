package sync

import "github.com/quantadex/distributed_quanta_bridge/trust/coin"

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

	return deposits, err
}

func (c *BitcoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}