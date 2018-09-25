package key_manager

import (
	"github.com/stellar/go/keypair"
	"github.com/spf13/viper"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
)

type QuantaKeyManager struct {
	key keypair.KP
}

func (k *QuantaKeyManager) CreateNodeKeys() error {
	return nil
}

func (k *QuantaKeyManager) LoadNodeKeys() error {
	var err error
	privateKey := viper.GetString("NODE_KEY")
	k.key, err = keypair.Parse(privateKey)
	return err
}

func (k *QuantaKeyManager) GetPublicKey() (string, error) {
	return k.key.Address(), nil
}

func (k *QuantaKeyManager) SignMessage(original []byte) ([]byte, error) {
	return k.key.Sign(original)
}

func (k *QuantaKeyManager) VerifySignatureObj(msg interface{}, signature string) bool {
	return crypto.VerifyMessage(msg, k.key.Address(), signature)
}

func (k *QuantaKeyManager) SignMessageObj(msg interface{}) *string {

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	signed, _ := k.key.Sign(bData.Bytes())
	signedbase64 := base64.StdEncoding.EncodeToString(signed)
	return &signedbase64
}
