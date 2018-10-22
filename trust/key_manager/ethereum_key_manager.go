package key_manager

import (
	"io/ioutil"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
	"strings"
	"fmt"
)

type EthereumKeyManager struct{
	key *ecdsa.PrivateKey
}

func (e *EthereumKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (e *EthereumKeyManager) LoadNodeKeys(filename string) error {

	if strings.HasPrefix(filename, "file://")  {
		keyjson, err := ioutil.ReadFile(strings.TrimPrefix(filename, "file://"))
		if err != nil {
			return err
		}

		key, err := keystore.DecryptKey(keyjson, "test123")
		if err != nil {
			return err
		}
		e.key = key.PrivateKey
		return nil
	} else {
		key, err := crypto.HexToECDSA(filename)
		if err != nil {
			return err
		}
		e.key = key
		return nil
	}
}

func (e *EthereumKeyManager) GetPublicKey() (string, error) {
	return crypto.PubkeyToAddress(e.key.PublicKey).Hex(), nil
}

func (e *EthereumKeyManager) GetPrivateKey() (*ecdsa.PrivateKey) {
	return e.key
}

func (e *EthereumKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("implement me")
}

func (e *EthereumKeyManager) SignMessageObj(original interface{}) (*string) {
	panic("implement me")
}

func (e *EthereumKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

// return hex signature
func (e *EthereumKeyManager) SignTransaction(hex string) (string, error) {
	dataBytes := common.Hex2Bytes(hex)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(dataBytes), dataBytes)
	h := crypto.Keccak256([]byte(msg))

	sig, err := crypto.Sign(h, e.key)
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(sig), nil
}

func (e *EthereumKeyManager) VerifyTransaction(hexSig string) (bool, error) {
	//dataBytes := common.Hex2Bytes(hexSig)
	//
	//crypto.VerifySignature(e.key.PrivateKey.PublicKey.)
	return false, nil
}
