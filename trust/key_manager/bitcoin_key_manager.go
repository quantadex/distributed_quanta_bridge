package key_manager

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/chaincfg"
)

type BitcoinKeyManager struct {
	privateKey *btcutil.WIF
	client *rpcclient.Client
	chaincfg *chaincfg.Params
}

func (b *BitcoinKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (b *BitcoinKeyManager) LoadNodeKeys(privKey string) error {
	panic("implement me")
}

func (b *BitcoinKeyManager) GetPublicKey() (string, error) {
	pub, err := btcutil.NewAddressPubKey(b.privateKey.SerializePubKey(), b.chaincfg)
	if err != nil {
		return "", err
	}
	return pub.EncodeAddress(), nil
}

func (b *BitcoinKeyManager) GetPrivateKey() (*ecdsa.PrivateKey) {
	return b.privateKey.PrivKey.ToECDSA()
}

func (b *BitcoinKeyManager) SignMessage(original []byte) ([]byte, error) {

}

func (b *BitcoinKeyManager) SignMessageObj(original interface{}) (*string) {
	panic("implement me")
}

func (b *BitcoinKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (b *BitcoinKeyManager) SignTransaction(encoded string) (string, error) {
	panic("implement me")
}

func (b *BitcoinKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}
