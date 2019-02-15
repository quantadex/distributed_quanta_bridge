package coin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

/*
cSJ2vqDoT9p6PXqNdNzaLNtMKePVjEzvoAryUUN7qAAB4njLKMXa 2018-10-29T23:34:06Z
reserve=1 # addr=2N39mAkxmLNnnL9WYecjkkTTtHTVQ3RtfZx hdkeypath=m/0'/0'/5'


cUgyLdmWgiMZcnCgTnmV1ag2evz5Eid6HQVacqfXPpWUbxgJcGt6 2018-10-29T23:34:06Z
reserve=1 # addr=2NF63kkxcegxtMuTKartK4tsyXsoHxhRvpN hdkeypath=m/0'/0'/13'

*/

func TestCheckHash(t *testing.T) {
	scriptBytes, _ := hex.DecodeString("004730440220304a3f60b7f5510e80b086cee9e88e38672c6031c3c2905b39bf5b180ba463b602205a14fcffee23eda01193e4dd281be474b058cc5d58b523f6bbc882469e05647c01473044022036224d12535cb02d597e1e1d0a6baccca9c27531464935110bab23c7a40f0cdd02206de86e931083484f6c23da41455863e6ff0c6b1359f8edf5f54364b83deea39c0147522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e52ae")
	addr, _ := btcutil.NewAddressScriptHash(scriptBytes, &chaincfg.RegressionNetParams)
	println(addr.EncodeAddress(), addr.String())
}

func TestBitcoinEncodeRefund(t *testing.T) {
	client, err := NewBitcoinCoin(&chaincfg.RegressionNetParams, []string{"2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", "2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	bitcoin := client.(*BitcoinCoin)
	crosschainAddr := make(map[string]string)
	crosschainAddr["2NA4mXEQvB594k2cZX4maVayQDTw8o7PG5m"] = "pooja"
	bitcoin.crosschainAddr = crosschainAddr
	//btec, err := crypto.NewGraphenePublicKeyFromString("QA5nvEN2S7Dej2C9hrLJTHNeMGeHq6uyjMdoceR74CksyApeZHWS")
	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)
	msig, err := bitcoin.GenerateMultisig(btec)

	log.Println("multisig: ", msig, err)

	w := Withdrawal{
		SourceAddress:      "2NA4mXEQvB594k2cZX4maVayQDTw8o7PG5m",
		DestinationAddress: "2NGYCnkuo62kL1QpAzV3bRaf747bSM8suQm",
		Amount:             1000,
		QuantaBlockID:      0,
	}
	tx, err := client.EncodeRefund(w)
	assert.NoError(t, err)
	var encoded EncodedMsg
	json.Unmarshal([]byte(tx), &encoded)

	km, _ := key_manager.NewBitCoinKeyManager()

	err = km.LoadNodeKeys("cNxQax7BfpbikeuCebPGCgTefTah5h1XhVDfaotVdFmXtaLCWLd9")
	assert.NoError(t, err)

	tx_signed1, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	err = km.LoadNodeKeys("cUixT9PYjTtNzcVjF8sB7iM9JeEf8tLHm9Wjgo972x8opCRNTasS")
	tx_signed2, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	fmt.Println(tx)
	fmt.Println(tx_signed1)
	fmt.Println(tx_signed2)

	w.Signatures = []string{tx_signed1, tx_signed2}
	hash, err := bitcoin.SendWithdrawal(common.HexToAddress("0x0"), nil, &w)

	assert.NoError(t, err)
	fmt.Println("hash", hash, err)
}

func TestTopBlockId(t *testing.T) {
	client, err := NewBitcoinCoin(&chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)
}

func TestDeposits(t *testing.T) {
	client, err := NewBitcoinCoin(&chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	_, err = client.GetDepositsInBlock(130, nil)
	assert.NoError(t, err)
}

func TestDecode(t *testing.T) {
	client, err := NewBitcoinCoin(&chaincfg.RegressionNetParams, []string{"2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", "2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	bitcoin := client.(*BitcoinCoin)
	crosschainAddr := make(map[string]string)
	crosschainAddr["n2PNkvCSkkSKvgqLsQXAQACFETQwKvc16X"] = "pooja"
	bitcoin.crosschainAddr = crosschainAddr
	//btec, err := crypto.NewGraphenePublicKeyFromString("QA5nvEN2S7Dej2C9hrLJTHNeMGeHq6uyjMdoceR74CksyApeZHWS")
	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)

	msig, err := bitcoin.GenerateMultisig(btec)

	log.Println("multisig: ", msig, err)

	w := Withdrawal{
		SourceAddress:      "n2PNkvCSkkSKvgqLsQXAQACFETQwKvc16X",
		DestinationAddress: "2NGYCnkuo62kL1QpAzV3bRaf747bSM8suQm",
		Amount:             1000,
		Tx:                 "4418603_0",
		QuantaBlockID:      0,
	}
	tx, err := client.EncodeRefund(w)
	fmt.Println("Encoded = ", tx)
	assert.NoError(t, err)

	_, err = client.DecodeRefund(tx)
	assert.NoError(t, err)
}
