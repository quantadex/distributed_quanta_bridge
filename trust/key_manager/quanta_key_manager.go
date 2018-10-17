package key_manager

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/stellar/go/xdr"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"log"
	"crypto/ecdsa"
)

type QuantaKeyManager struct {
	key keypair.KP
	seed string
	network string
}

func (k *QuantaKeyManager) SignTransaction(base64 string) (string, error) {
	txe := &xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(base64, txe)
	if err != nil {
		return "", err
	}

	b := &build.TransactionEnvelopeBuilder{E: txe}
	b.Init()

	err = b.MutateTX(build.Network{k.network})
	if err != nil {
		log.Fatal(err)
	}

	err = b.Mutate(build.Sign{k.seed})
	if err != nil {
		println("Sign error")
		return "", err
	}
	return xdr.MarshalBase64(b.E)
}


func (k *QuantaKeyManager) VerifyTransaction(base64 string) (bool, error) {
	println("verify ", base64)
	txe := &xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(base64, txe)
	if err != nil {
		println("ERROR!", err)
		return false, err
	}

	//b := &build.TransactionEnvelopeBuilder{E: txe}
	//b.Init()
	//
	//var txBytes bytes.Buffer
	//_, err = xdr.Marshal(&txBytes, txe.Tx)
	//if err != nil {
	//	return false, errors.New("unable to marshall bytes")
	//}

	//txHash := hash.Hash(txBytes.Bytes())
	//signature, err := skp.Sign(txHash[:])
	//if err != nil {
	//	panic(err)
	//}

	// let's delay transaction validation to stellar
	return true, nil
}

func (k *QuantaKeyManager) CreateNodeKeys() error {
	return nil
}

func (k *QuantaKeyManager) LoadNodeKeys(privkey string) error {
	var err error
	privateKey := privkey
	k.key, err = keypair.Parse(privateKey)
	k.seed = privateKey

	return err
}

func (k *QuantaKeyManager) GetPublicKey() (string, error) {
	return k.key.Address(), nil
}

func (k *QuantaKeyManager) GetPrivateKey() (*ecdsa.PrivateKey) {
	return nil
}

func (k *QuantaKeyManager) VerifySignatureObj(msg interface{}, signature string) bool {
	return crypto.VerifyMessage(msg, k.key.Address(), signature)
}


func (k *QuantaKeyManager) SignMessage(original []byte) ([]byte, error) {
	return k.key.Sign(original)
}

func (k *QuantaKeyManager) SignMessageObj(msg interface{}) *string {

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	signed, _ := k.key.Sign(bData.Bytes())
	signedbase64 := base64.StdEncoding.EncodeToString(signed)
	return &signedbase64
}
