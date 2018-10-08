package coin

import (
	"fmt"
	"sync"
)

type DummyCoin struct {
	index int64
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
	c.deposits[int(c.index)] = append(c.deposits[int(c.index)], deposit)

	return nil
}

func (c *DummyCoin) CreateNewBlock() {
	c.index = c.index + 1
}

func (c *DummyCoin) GetTopBlockID() (int64, error) {
	return c.index, nil
}

func (c *DummyCoin) GetDepositsInBlock(blockID int64, trustAddress string) ([]*Deposit, error) {
	return c.deposits[int(blockID)], nil
}

func (c *DummyCoin) SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error {
	panic("implement me")
}

func (c *DummyCoin) GetForwardersInBlock(blockID int64, trustAddress string) ([]*ForwardInput, error) {
	return []*ForwardInput{}, nil
}

func (c *DummyCoin) Attach() error {
	return nil
}

