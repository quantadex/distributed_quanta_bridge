package coin

import (
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
)

type BitcoinCoin struct {
	
}

func (b *BitcoinCoin) Attach() error {
	panic("implement me")
}

func (b *BitcoinCoin) GetTopBlockID() (int64, error) {
	panic("implement me")
}

func (b *BitcoinCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	panic("implement me")
}

func (b *BitcoinCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	panic("implement me")
}

func (b *BitcoinCoin) GetForwardersInBlock(blockID int64) ([]*ForwardInput, error) {
	panic("implement me")
}

func (b *BitcoinCoin) SendWithdrawal(trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {
	panic("implement me")
}

func (b *BitcoinCoin) EncodeRefund(w Withdrawal) (string, error) {
	panic("implement me")
}

func (b *BitcoinCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	panic("implement me")
}


