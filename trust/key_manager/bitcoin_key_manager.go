package key_manager

import (
	"crypto/ecdsa"
)

type BitcoinKeyManager struct {
}

func (b *BitcoinKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (b *BitcoinKeyManager) LoadNodeKeys(privKey string) error {
	panic("implement me")
}

func (b *BitcoinKeyManager) GetPublicKey() (string, error) {
	panic("implement me")
}

func (b *BitcoinKeyManager) GetPrivateKey() (*ecdsa.PrivateKey) {
	panic("implement me")
}

func (b *BitcoinKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("implement me")
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
