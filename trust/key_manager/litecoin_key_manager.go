package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common"
)

type LitecoinKeyManager struct {
	privateKey  *ltcutil.WIF
	client      *rpcclient.Client
	chaincfg    *chaincfg.Params
	bitcoinRPC  string
	rpcUser     string
	rpcPassword string
}

func (b *LitecoinKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (b *LitecoinKeyManager) LoadNodeKeys(privKey string) error {
	var err error
	// TODO: make this configurable
	b.client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.bitcoinRPC,
		User:         b.rpcUser,
		Pass:         b.rpcPassword,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)

	if err != nil {
		return errors.Wrap(err, "Cannot attach client for LTC")
	}

	//err = crypto.ValidateNetwork(b.client, "Litecoin")
	//if err != nil {
	//	return err
	//}

	b.privateKey, err = ltcutil.DecodeWIF(privKey)

	return err
}

func (b *LitecoinKeyManager) GetPublicKey() (string, error) {
	pub, err := ltcutil.NewAddressPubKey(b.privateKey.SerializePubKey(), b.chaincfg)
	if err != nil {
		return "", err
	}
	return pub.EncodeAddress(), nil
}

func (b *LitecoinKeyManager) GetPrivateKey() *ecdsa.PrivateKey {
	return b.privateKey.PrivKey.ToECDSA()
}

func (b *LitecoinKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("not required")
}

func (b *LitecoinKeyManager) SignMessageObj(original interface{}) *string {
	panic("not required")
}

func (b *LitecoinKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (b *LitecoinKeyManager) SignTransaction(encoded string) (string, error) {
	var res common.TransactionLitecoin

	err := json.Unmarshal([]byte(encoded), &res)
	if err != nil {
		return "", err
	}

	dataBytes, err := hex.DecodeString(res.Tx)
	if err != nil {
		return "", err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	err = tx.Deserialize(bytes.NewBuffer(dataBytes))

	if err != nil {
		return "", err
	}

	txSigned, _, err := b.client.SignRawTransaction3(tx, res.RawInput, []string{b.privateKey.String()})
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = txSigned.Serialize(&buf)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), err
}

func (b *LitecoinKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}
