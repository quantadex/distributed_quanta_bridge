package coin

import (
	"fmt"
	"sync"
)

type DummyCoin struct {
	index int
	deposits map[int][]*Deposit
}

var instance *DummyCoin
var once sync.Once

func GetDummyInstance() *DummyCoin {
	once.Do(func () {
		instance = &DummyCoin{ deposits: map[int][]*Deposit{}}
	})
	return instance
}

func (c *DummyCoin) AddDeposit(deposit *Deposit) (error) {
	fmt.Printf("insert deposit into block %d\n", c.index)
	c.deposits[c.index] = append(c.deposits[c.index], deposit)

	return nil
}

func (c *DummyCoin) CreateNewBlock() {
	c.index = c.index + 1
}

func (c *DummyCoin) GetTopBlockID() (int, error) {
	return c.index, nil
}

func (c *DummyCoin) GetDepositsInBlock(blockID int, trustAddress string) ([]*Deposit, error) {
	return c.deposits[blockID], nil
}

func (c *DummyCoin) SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error {
	panic("implement me")
}

