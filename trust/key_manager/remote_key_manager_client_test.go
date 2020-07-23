package key_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyManager(t *testing.T) {
	eth, _ := NewEthKeyManager()
	eth.LoadNodeKeys("0afd879321bf647a8ae9484e780916e681710255d733b3d5033710aa09ddfcd1")

	km, _ := NewBitCoinKeyManager("localhost:18332", "regnet", "user", "123", []string{"1", "2", "3"})
	err := km.LoadNodeKeys("92REaZhgcw6FF2rz8EnY1HMtBvgh3qh4gs9PxnccPrju6ZCFetk")
	assert.NoError(t, err)

	kms := map[string]KeyManager{
		"ETH": eth,
		"BTC": km,
	}

	service := NewRemoteKeyManagerService(kms, common2.Secrets{DatabaseUrl: "db://secret"})
	service.Serve(":4444")

	//time.Sleep(time.Second)

	println("Connecting..")
	client, err := NewRemoteKeyManager("ETH", "localhost:4444")
	if err != nil {
		println("connected? ", err.Error())
	}

	var reply SignResponse
	req := &SignMessage{"ETH", common.Bytes2Hex([]byte("This is proof that I, user A, have access to this address"))}
	err = client.Client.Call("Signer.SignTx", req, &reply)
	println("signed", reply.Signed, "ERR=", err)
	if err != nil {
		println(err.Error())
	}

	pubkey, _ := client.GetPublicKey()
	println("pub key=", pubkey)

	var res common2.Secrets
	res, err = client.GetSecretsWithoutKeys()

	if err != nil {
		print(err.Error())
	}
	fmt.Printf("URL=%v\n", res.DatabaseUrl)

	client, err = NewRemoteKeyManager("BTC", "localhost:4444")
	signers := client.GetSigners()
	fmt.Printf("signers=%v\n", signers)

}
