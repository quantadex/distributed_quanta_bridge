package key_manager

import (
	"github.com/spf13/viper"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/stellar_go/xdr"
	"github.com/quantadex/stellar_go/build"
	"github.com/quantadex/stellar_go/keypair"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/pkg/errors"
	"log"
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
	txe := &xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(base64, txe)
	if err != nil {
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

func (k *QuantaKeyManager) LoadNodeKeys() error {
	var err error
	privateKey := viper.GetString("NODE_KEY")
	k.key, err = keypair.Parse(privateKey)
	k.seed = privateKey
	k.network = viper.GetString("NETWORK_PASSPHRASE")

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

func (k *QuantaKeyManager) DecodeTransaction(base64 string) (*coin.Deposit, error) {
	txe := &xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(base64, txe)
	if err != nil {
		return nil, err
	}

	ops := txe.Tx.Operations
	if len(ops) != 1 {
		return nil, errors.New("no operations found")
	}

	paymentOp, success := ops[0].Body.GetPaymentOp()
	if !success {
		return nil, errors.New("no payment op found")
	}

	return &coin.Deposit{ CoinName: paymentOp.Asset.String(),
					QuantaAddr: paymentOp.Destination.Address(),
					Amount: int64(paymentOp.Amount),
					BlockID: 0,
	}, nil
}
