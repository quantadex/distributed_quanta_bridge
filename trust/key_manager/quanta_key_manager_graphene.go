package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/scorum/bitshares-go/sign"
)

type QuantaKeyGraphene struct {
	chain      string
	privateKey []btcec.PrivateKey
}

func (k *QuantaKeyGraphene) LoadNodeKeys(wif string) error {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return err
	}
	k.privateKey = append(k.privateKey, *w.PrivKey)

	fmt.Println(k.privateKey)
	return err
}

func (k *QuantaKeyGraphene) SignTransaction(encoded string) (string, error) {
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(encoded), &tx)

	digest, err := tx.Digest(k.chain)
	if err != nil {
		return "", err
	}
	var i int
	var sigHex string
	for i = 0; i < len(k.privateKey); i++ {
		sig := sign.SignBufferSha256(digest, k.privateKey[i].ToECDSA())
		sigHex = hex.EncodeToString(sig) + sigHex
	}
	fmt.Println(sigHex)
	return sigHex, nil
}

func (k *QuantaKeyGraphene) VerifyTransaction(encoded string) (bool, error) {
	var tx sign.SignedTransaction
	err := json.Unmarshal([]byte(encoded), &tx)

	if err != nil {
		println("ERROR!", err)
		return false, err
	}
	return true, nil
}

func (k *QuantaKeyGraphene) SignMessageObj(msg interface{}) *string {
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	digest := sha256.Sum256(bData.Bytes())

	var i int
	var sigHex string
	for i = 0; i < len(k.privateKey); i++ {
		sig := sign.SignBufferSha256(digest[:], k.privateKey[i].ToECDSA())
		sigHex = hex.EncodeToString(sig) + sigHex
	}

	return &sigHex
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
	digest := sha256.Sum256(original)

	var i int
	var sig []byte
	for i = 0; i < len(k.privateKey); i++ {
		s := sign.SignBufferSha256(digest[:], k.privateKey[i].ToECDSA())
		sig = append(sig, s...)
	}
	return sig, nil
}
