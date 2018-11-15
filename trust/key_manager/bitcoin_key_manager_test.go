package key_manager

import (
	"fmt"
  "github.com/stretchr/testify/assert"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"testing"
)


func TestLoadNodeKeysFromFile(t *testing.T) {
	km, err := NewBtcKeyManager()
  assert.NoError(t, err)
  assert.NotNil(t, km)

	err = km.LoadNodeKeys("../../keystore/test_btc_wallet.db")
  assert.NoError(t, err)
}

func TestLoadNodeKeysFromBadFile(t *testing.T) {
  km, _ := NewBtcKeyManager()
	err := km.LoadNodeKeys("/dev/null")
  assert.Error(t, err)
}

func TestGetPublicKey(t *testing.T) {
  km, _ := NewBtcKeyManager()
  km.LoadNodeKeys("../../keystore/test_btc_wallet.db")
  publicKey, err := km.GetPublicKey()
  assert.NoError(t, err)
  assert.Equal(t, "0x772224f88274E6573e1Cf37351ECB828CC644856", publicKey)
}

func TestGetPrivateKey(t *testing.T) {
  km, _ := NewBtcKeyManager()
  km.LoadNodeKeys("../../keystore/test_btc_wallet.db")
  privateKey := km.GetPrivateKey()

  privateKeyHex := fmt.Sprintf("%x", privateKey.D.Bytes())
  assert.Equal(t, "2d70a8013248e8d8406917269bc05ef002aad62531d193ebae5ab7f25f2138e2", privateKeyHex)
}

func TestBitcoinWithdrawalTX(t *testing.T) {
	km, _ := NewBtcKeyManager()
	km.LoadNodeKeys("../../keystore/test_btc_wallet.db")

  wif, _ := BTC_NETWORK.CreatePrivateKey()
  dest, _ := BTC_NETWORK.GetAddress(wif)

	w := &coin.Withdrawal{
		TxId:               1,
		CoinName:           "BTC",
		DestinationAddress: dest.String(),
		QuantaBlockID:      1,
		Amount:             10000,
		Signatures:         nil,
	}

	coin := &coin.BitcoinCoin{}
	encoded, err := coin.EncodeRefund(*w)
	assert.NoError(t, err)
	fmt.Printf("encoded: %s\n", encoded)
	signed, err := km.SignTransaction(encoded)
	assert.NoError(t, err)
	fmt.Printf("signed: %s\n", signed)

  w.Signatures = []string{signed, signed, signed}

	// client := &Listener{NetworkID: ROPSTEN_NETWORK_ID}
	// tx, err := client.SendWithdrawal(sim, userAuth.From, userKey, w)

	// if err != nil {
	// 	println("ERR: ", err.Error())
	// }
	// println(tx)
}
