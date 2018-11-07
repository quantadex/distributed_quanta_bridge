package control

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
	"github.com/magiconair/properties/assert"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"os"
	"strconv"
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
func TestSubmitWithdrawalRecovery(t *testing.T) {
	w := &coin.Withdrawal{
		TxId:               1,
		CoinName:           "ETH",
		DestinationAddress: "0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef",
		QuantaBlockID:      1,
		Amount:             12345,
		Signatures:         nil,
	}

	m := &MockEthereumCoin{}
	os.Remove("withdrawal_recovery")
	kv := &kv_store.BoltStore{}
	kv.Connect("withdrawal_recovery")
	quanta2Coin := &QuantaToCoin{
		coinChannel:        m,
		db:                 kv,
		quantaTrustAddress: "0xdda6327139485221633a1fcd65f4ac932e60a2e1",
	}
	key := strconv.Itoa(int(w.TxId))
	m.setWithdrawalOutput(true, "ERR:  failed to retrieve account nonce: Post http://localhost:7545: dial tcp 127.0.0.1:7545: connect: connection refused")
	quanta2Coin.SubmitWithdrawal(w)
	msgSigned, _ := kv.GetValue(kv_store.ETH_TX_LOG_SIGNED, key)
	msgRetry, _ := kv.GetValue(kv_store.ETH_TX_LOG_RETRY, key)

	msg, _ := json.Marshal(w)
	value := string(msg)
	assert.Equal(t, *msgSigned, value)
	assert.Equal(t, *msgRetry, value)

	m.setWithdrawalOutput(false, "")
	quanta2Coin.SubmitWithdrawal(w)

	msgSubmittedAfter, _ := kv.GetValue(kv_store.ETH_TX_LOG_SUBMITTED, key)
	assert.Equal(t, *msgSubmittedAfter, value)
}
=======

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
	w := &coin.Withdrawal{
		TxId:               1,
		CoinName:           "ETH",
		DestinationAddress: "0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef",
		QuantaBlockID:      1,
		Amount:             12345,
		Signatures:         nil,
	}

	m := &MockEthereumCoin{}
	os.Remove("withdrawal_recovery")
	kv := &kv_store.BoltStore{}
	kv.Connect("withdrawal_recovery")
	quanta2Coin := &QuantaToCoin{
		coinChannel:        m,
		db:                 kv,
		quantaTrustAddress: "0xdda6327139485221633a1fcd65f4ac932e60a2e1",
	}
	key := strconv.Itoa(int(w.TxId))
	m.setWithdrawalOutput(true, "ERR:  failed to retrieve account nonce: Post http://localhost:7545: dial tcp 127.0.0.1:7545: connect: connection refused")
	quanta2Coin.SubmitWithdrawal(w)
	msgSigned, _ := kv.GetValue(kv_store.ETH_TX_LOG_SIGNED, key)
	msgRetry, _ := kv.GetValue(kv_store.ETH_TX_LOG_RETRY, key)

	msg, _ := json.Marshal(w)
	value := string(msg)
	assert.Equal(t, *msgSigned, value)
	assert.Equal(t, *msgRetry, value)

	m.setWithdrawalOutput(false, "")
	quanta2Coin.SubmitWithdrawal(w)

	msgSubmittedAfter, _ := kv.GetValue(kv_store.ETH_TX_LOG_SUBMITTED, key)
	assert.Equal(t, *msgSubmittedAfter, value)

	
}


>>>>>>> origin/master
