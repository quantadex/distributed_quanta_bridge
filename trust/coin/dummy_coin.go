package coin

import (
	"crypto/ecdsa"
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"sync"
)

type DummyCoin struct {
	index    int64
	deposits map[int][]*Deposit
}

var instance *DummyCoin
var once sync.Once

func GetDummyInstance() *DummyCoin {
	once.Do(func() {
		instance = &DummyCoin{deposits: map[int][]*Deposit{}}
	})
	return instance
}

func (c *DummyCoin) Blockchain() string {
	return "DUMMY"
}

func (b *DummyCoin) GenerateMultisig(accountId string) (string, error) {
	panic("not implemented")
}

func (c *DummyCoin) AddDeposit(deposit *Deposit) error {
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

func (c *DummyCoin) GetTxID(trustAddress common2.Address) (uint64, error) {
	return 0, nil
}

func (c *DummyCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	return c.deposits[int(blockID)], nil
}

func (c *DummyCoin) SendWithdrawal(trustAddress common2.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {
	panic("implement me")
}

func (c *DummyCoin) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	return []*crypto.ForwardInput{}, nil
}

func (c *DummyCoin) Attach() error {
	return nil
}

func (c *DummyCoin) FillCrosschainAddress(crosschainAddr map[string]string) {

}

func (c *DummyCoin) EncodeRefund(w Withdrawal) (string, error) {
	return "", nil
}

func (c *DummyCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	return nil, nil
}

func (c *DummyCoin) CheckValidAddress(address string) bool {
	panic("Not implemented")
}
