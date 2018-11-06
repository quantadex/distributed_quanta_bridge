package control

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"crypto/ecdsa"
	"github.com/go-errors/errors"
	"testing"
)

type MockEthereumCoin struct {
	error bool
	message string
}

func (c *MockEthereumCoin) setWithdrawalOutput(error bool, message string) {
	c.error = error
	c.message = message
}

func (c *MockEthereumCoin) Attach() error {
	panic("implement me")
}

func (c *MockEthereumCoin) GetTopBlockID() (int64, error) {
	panic("implement me")
}

func (c *MockEthereumCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	panic("implement me")
}

func (c *MockEthereumCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*coin.Deposit, error) {
	panic("implement me")
}

func (c *MockEthereumCoin) GetForwardersInBlock(blockID int64) ([]*coin.ForwardInput, error) {
	panic("implement me")
}

func (c *MockEthereumCoin) SendWithdrawal(trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *coin.Withdrawal) (string, error) {

		if c.error {
			return "",errors.New(c.message)
		}
		return "tx_id", nil
}

func (c *MockEthereumCoin) EncodeRefund(w coin.Withdrawal) (string, error) {
	panic("implement me")
}

func (c *MockEthereumCoin) DecodeRefund(encoded string) (*coin.Withdrawal, error) {
	panic("implement me")
}

func TestFunc (t *testing.T) {
	//m := &MockEthereumCoin{}
	//quanta2Coin := &QuantaToCoin{
	//	coinChannel: m,
	//}
	//
	//m.setWithdrawalOutput(true, "")

	// call code

	// expect your code to handle the withdrawal error

	
}


