package key_manager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"log"
	"net"
	"net/rpc"
)

type RemoteKeyManager struct {
	Client     *rpc.Client
	blockchain string
	address    string
}

func (r *RemoteKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (r *RemoteKeyManager) LoadNodeKeys(privKey string) error {
	panic("implement me")
}

func (r *RemoteKeyManager) GetPublicKey() (pub string, err error) {
	err = r.Connect()
	err = r.Client.Call("Signer.GetPublicKey", r.blockchain, &pub)
	return
}

func (r *RemoteKeyManager) GetPrivateKey() *ecdsa.PrivateKey {
	err := r.Connect()
	var priv string
	err = r.Client.Call("Signer.GetPrivateKey", r.blockchain, &priv)
	if err != nil {
		return nil
	}

	key, err := crypto.HexToECDSA(priv)
	if err != nil {
		return nil
	}
	return key
}

func (r *RemoteKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("implement me")
}

func (r *RemoteKeyManager) GetSigners() (ret []string) {
	r.Connect()
	r.Client.Call("Signer.GetSigners", r.blockchain, &ret)
	return
}

/**
 * pseudosignature for now.
 */
func (r *RemoteKeyManager) SignMessageObj(msg interface{}) *string {
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	//signed, _ := k.key.Sign(bData.Bytes())
	signedbase64 := base64.StdEncoding.EncodeToString(bData.Bytes())
	return &signedbase64
}

func (r *RemoteKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	return true
}

func (r *RemoteKeyManager) SignTransaction(encoded string) (string, error) {
	err := r.Connect()

	var reply SignResponse
	err = r.Client.Call("Signer.SignTx", &SignMessage{r.blockchain, encoded}, &reply)
	return reply.Signed, err
}

func (r *RemoteKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}

func (r *RemoteKeyManager) GetSecretsWithoutKeys() (secrets common.Secrets, err error) {
	err = r.Connect()
	err = r.Client.Call("Secrets.GetSecretsWithoutKeys", "", &secrets)
	return
}

func NewRemoteKeyManager(blockchain string, address string) (*RemoteKeyManager, error) {
	km := &RemoteKeyManager{nil, blockchain, address}
	err := km.Connect()
	return km, err
}

func (r *RemoteKeyManager) Connect() error {
	var e error
	var tcp net.Conn
	tcp, e = net.Dial("tcp", r.address)
	client := rpc.NewClient(tcp)

	if e != nil {
		log.Fatal("Unable to connect remote km", e)
		return e
	}
	r.Client = client
	return nil
}
