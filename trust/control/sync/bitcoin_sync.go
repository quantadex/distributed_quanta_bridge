package sync

import "github.com/quantadex/distributed_quanta_bridge/trust/coin"

type BitcoinSync struct {
	DepositSync
}

func (c *BitcoinSync) Setup() {

}

func (c *BitcoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	panic("not imp")
}

func (c *BitcoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}