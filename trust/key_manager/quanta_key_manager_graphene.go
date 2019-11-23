package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/scorum/bitshares-go/sign"
)

type QuantaKeyGraphene struct {
	chain      string
	privateKey *btcec.PrivateKey
}

func (k *QuantaKeyGraphene) LoadNodeKeys(wif string) error {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return err
	}
	k.privateKey = w.PrivKey
	return err
}

func (k *QuantaKeyGraphene) SignTransaction(encoded string) (string, error) {
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(encoded), &tx)

	digest, err := tx.Digest(k.chain)
	if err != nil {
		return "", err
	}
	sig := sign.SignBufferSha256(digest, k.privateKey.ToECDSA())
	return hex.EncodeToString(sig), nil
}

func (k *QuantaKeyGraphene) VerifyTransaction(encoded string) (bool, error) {
	//var tx sign.SignedTransaction
	//err := json.Unmarshal([]byte(encoded), &tx)

	//if err != nil {
	//	println("ERROR!", err)
	//	return false, err
	//}
	return true, nil
}

func (k *QuantaKeyGraphene) SignMessageObj(msg interface{}) *string {
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	digest := sha256.Sum256(bData.Bytes())

	sig := sign.SignBufferSha256(digest[:], k.privateKey.ToECDSA())
	sigHex := hex.EncodeToString(sig)
	return &sigHex
}

func (k *QuantaKeyGraphene) CreateNodeKeys() error {
	return nil
}

func (k *QuantaKeyGraphene) GetPublicKey() (string, error) {
	return crypto.GetGraphenePublicKey(k.privateKey.PubKey())
}

func (k *QuantaKeyGraphene) GetPrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (k *QuantaKeyGraphene) VerifySignatureObj(msg interface{}, signature string) bool {
	//return crypto.VerifyMessage(msg, k.key.Address(), signature)
	publicKey, err := k.GetPublicKey()
	if err != nil {
		return false
	}
	return crypto.VerifyMessage(msg, publicKey, signature)
}

func (k *QuantaKeyGraphene) SignMessage(original []byte) ([]byte, error) {
	digest := sha256.Sum256(original)
	sig := sign.SignBufferSha256(digest[:], k.privateKey.ToECDSA())
	return sig, nil
}
