package coin

import (
	"testing"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/btcsuite/btcutil"
	"log"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/ethereum/go-ethereum/common"
)

/*
cSJ2vqDoT9p6PXqNdNzaLNtMKePVjEzvoAryUUN7qAAB4njLKMXa 2018-10-29T23:34:06Z
reserve=1 # addr=2N39mAkxmLNnnL9WYecjkkTTtHTVQ3RtfZx hdkeypath=m/0'/0'/5'


cUgyLdmWgiMZcnCgTnmV1ag2evz5Eid6HQVacqfXPpWUbxgJcGt6 2018-10-29T23:34:06Z
reserve=1 # addr=2NF63kkxcegxtMuTKartK4tsyXsoHxhRvpN hdkeypath=m/0'/0'/13'

 */
func TestBitcoinEncodeRefund(t *testing.T) {
	client, err := NewBitcoinCoin(&chaincfg.RegressionNetParams)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	addr1, err := btcutil.DecodeAddress("2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", &chaincfg.RegressionNetParams)
	addr2, err := btcutil.DecodeAddress("2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR", &chaincfg.RegressionNetParams)
	if err != nil {
		println(err)
		assert.NoError(t, err)
	}

	println(addr1.String(), addr1.EncodeAddress())

	bitcoin := client.(*BitcoinCoin)
	msig, err := bitcoin.GenerateMultisig([]btcutil.Address {
		addr1, addr2,
	})

	log.Println("multisig: ", msig, err)

	w := Withdrawal{
		SourceAddress: msig,
		DestinationAddress: "2NGYCnkuo62kL1QpAzV3bRaf747bSM8suQm",
		Amount: 1000,
	}
	tx, err := client.EncodeRefund(w)
	assert.NoError(t, err)

	km,_ := key_manager.NewBitCoinKeyManager()

	err = km.LoadNodeKeys("cNxQax7BfpbikeuCebPGCgTefTah5h1XhVDfaotVdFmXtaLCWLd9")
	assert.NoError(t, err)

	tx_signed1, err := km.SignTransaction(tx)
	assert.NoError(t, err)

	err = km.LoadNodeKeys("cUixT9PYjTtNzcVjF8sB7iM9JeEf8tLHm9Wjgo972x8opCRNTasS")
	tx_signed2, err := km.SignTransaction(tx)
	assert.NoError(t, err)

	fmt.Println(tx)
	fmt.Println(tx_signed1)
	fmt.Println(tx_signed2)

	w.Signatures = []string{ tx_signed1, tx_signed2}
	hash,err := bitcoin.SendWithdrawal(common.HexToAddress("0x0"), nil, &w)

	assert.NoError(t, err)
	fmt.Println("hash", hash, err)
}