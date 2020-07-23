package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common"
)

type BitcoinKeyManager struct {
	privateKey  *btcutil.WIF
	client      *rpcclient.Client
	chaincfg    *chaincfg.Params
	bitcoinRPC  string
	rpcUser     string
	rpcPassword string
	signers     []string
}

func (b *BitcoinKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (b *BitcoinKeyManager) LoadNodeKeys(privKey string) error {
	var err error
	// TODO: make this configurable
	b.client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.bitcoinRPC,
		User:         b.rpcUser,
		Pass:         b.rpcPassword,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)

	if err != nil {
		return errors.Wrap(err, "Cannot attach client for BTC")
	}

	//err = crypto.ValidateNetwork(b.client, "Satoshi")
	//if err != nil {
	//	return err
	//}

	b.privateKey, err = btcutil.DecodeWIF(privKey)

	return err
}

func (b *BitcoinKeyManager) GetPublicKey() (string, error) {
	pub, err := btcutil.NewAddressPubKey(b.privateKey.SerializePubKey(), b.chaincfg)
	if err != nil {
		return "", err
	}
	return pub.EncodeAddress(), nil
}

func (r *BitcoinKeyManager) GetSigners() []string {
	return r.signers
}

func (b *BitcoinKeyManager) GetPrivateKey() *ecdsa.PrivateKey {
	return b.privateKey.PrivKey.ToECDSA()
}

func (b *BitcoinKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("not required")
}

func (b *BitcoinKeyManager) SignMessageObj(original interface{}) *string {
	panic("not required")
}

func (b *BitcoinKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (b *BitcoinKeyManager) SignTransaction(encoded string) (string, error) {
	var res common.TransactionBitcoin

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

func (b *BitcoinKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}
