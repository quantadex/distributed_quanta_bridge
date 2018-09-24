package key_manager

import (
	"github.com/stellar/go/keypair"
	"github.com/spf13/viper"
)

type QuantaKeyManager struct {
	key keypair.KP
}

func (k *QuantaKeyManager) CreateNodeKeys() error {
}

func (k *QuantaKeyManager) LoadNodeKeys() error {
	var err error
	publicKey := viper.GetString("NODE_KEY")
	k.key, err = keypair.Parse(publicKey)
	return err
}

func (k *QuantaKeyManager) GetPublicKey() (string, error) {
	return k.GetPublicKey()
}

func (k *QuantaKeyManager) SignMessage(original []byte) ([]byte, error) {
	return k.key.Sign(original)
}

func (k *QuantaKeyManager) DecodeMessage(original []byte, signature string) ([]byte, error) {
	return nil, nil
}


