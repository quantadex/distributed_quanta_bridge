package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/rpcclient"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common"
)

type BCHKeyManager struct {
	privateKey  *bchutil.WIF
	client      *rpcclient.Client
	chaincfg    *chaincfg.Params
	bitcoinRPC  string
	rpcUser     string
	rpcPassword string
	signers     []string
}

func (b *BCHKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (b *BCHKeyManager) LoadNodeKeys(privKey string) error {
	var err error
	// TODO: make this configurable
	b.client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.bitcoinRPC,
		User:         b.rpcUser,
		Pass:         b.rpcPassword,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)

	if err != nil {
		return errors.Wrap(err, "Cannot attach client for BCH")
	}

	//err = crypto.ValidateNetwork(b.client, "Bitcoin ABC")
	//if err != nil {
	//	return err
	//}

	b.privateKey, err = bchutil.DecodeWIF(privKey)

	return err
}

func (b *BCHKeyManager) GetPublicKey() (string, error) {
	pub, err := bchutil.NewAddressPubKey(b.privateKey.SerializePubKey(), b.chaincfg)
	if err != nil {
		return "", err
	}
	return pub.EncodeAddress(), nil
}

func (r *BCHKeyManager) GetSigners() []string {
	return r.signers
}

func (b *BCHKeyManager) GetPrivateKey() *ecdsa.PrivateKey {
	return b.privateKey.PrivKey.ToECDSA()
}

func (b *BCHKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("not required")
}

func (b *BCHKeyManager) SignMessageObj(original interface{}) *string {
	panic("not required")
}

func (b *BCHKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (b *BCHKeyManager) SignTransaction(encoded string) (string, error) {
	var res common.TransactionBCH

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
	fmt.Println("private key = ", b.privateKey.String())
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

func (b *BCHKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}
