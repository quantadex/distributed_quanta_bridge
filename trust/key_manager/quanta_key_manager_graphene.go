package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"github.com/agl/ed25519"
	"github.com/scorum/bitshares-go/types"
)

type QuantaKeyGraphene struct {
	seed    string
	network string
}

func (k *QuantaKeyGraphene) LoadNodeKeys(privkey string) error {
	var err error
	k.seed = privkey
	return err
}

func (k *QuantaKeyGraphene) SignTransaction(encoded string) (string, error) {
	txe := &types.Transaction{}
	encodedBytes := make([]byte, len(encoded))
	copy(encodedBytes, encoded)

	err := json.Unmarshal(encodedBytes, txe)
	if err != nil {
		return "", err
	}
	result, err := json.Marshal(txe.Signatures[0])
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (k *QuantaKeyGraphene) VerifyTransaction(encoded string) (bool, error) {
	txe := &types.Transaction{}
	encodedBytes := make([]byte, len(encoded))
	copy(encodedBytes, encoded)

	err := json.Unmarshal(encodedBytes, txe)
	if err != nil {
		println("ERROR!", err)
		return false, err
	}
	return true, nil
}

func (k *QuantaKeyGraphene) SignMessageObj(msg interface{}) *string {
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	var privKeys *[64]byte
	signed := ed25519.Sign(privKeys, bData.Bytes())
	signedbase64 := base64.StdEncoding.EncodeToString(signed[:])
	return &signedbase64
}

func (k *QuantaKeyGraphene) CreateNodeKeys() error {
	return nil
}

func (k *QuantaKeyGraphene) GetPublicKey() (string, error) {
	panic("Implement me")
}

func (k *QuantaKeyGraphene) GetPrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (k *QuantaKeyGraphene) VerifySignatureObj(msg interface{}, signature string) bool {
	//return crypto.VerifyMessage(msg, k.key.Address(), signature)
	panic("Implement me")
}

func (k *QuantaKeyGraphene) SignMessage(original []byte) ([]byte, error) {
	panic("Implement me")
}
