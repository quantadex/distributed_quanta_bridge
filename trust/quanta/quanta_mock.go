package quanta

import "fmt"

type DummyCoin struct {
	index int

	deposits map[int][]*Refund
}

var instance *DummyCoin

func GetDummyInstance() *DummyCoin {
	if instance == nil {
		instance = &DummyCoin{ deposits: map[int][]*Refund{}}
	}
	return instance
}

func (c *DummyCoin) AddRefund(deposit *Refund) (error) {
	fmt.Printf("insert deposit into block %d\n", c.index)
	c.deposits[c.index] = append(c.deposits[c.index], deposit)
	return nil
}

func (c *DummyCoin) CreateNewBlock() {
	c.index++
}

func (c *DummyCoin) GetTopBlockID() (int, error) {
	return c.index, nil
}

func (c *DummyCoin) GetDepositsInBlock(blockID int, trustAddress string) ([]*Refund, error) {
	fmt.Printf("Get deposit for block %d", blockID)
	return c.deposits[blockID], nil
}


