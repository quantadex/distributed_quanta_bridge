/**
 * Most test cases assumes you have a local
 * btcd websocket enabled RPC server running on localhost:4444 with
 * rpcuser=test and rpcpass=test
 */
package coin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const BTCD_RPC_HOST = "127.0.0.1:4444"
const BTCD_RPC_USER = "test"
const BTCD_RPC_PASS = "test"

func TestBitcoinDefaultConstructor(t *testing.T) {
	coin := &BitcoinCoin{}
	assert.NotNil(t, coin)
}

func TestBitcoinGetClientBadHost(t *testing.T) {
	coin := &BitcoinCoin{}

	// we assume nothing is listening on 4440
	client, err := coin.GetClient("localhost:4440", "foo", "badpass")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestBitcoinGetClientCanConnect(t *testing.T) {
	coin := &BitcoinCoin{}

	client, err := coin.GetClient(BTCD_RPC_HOST, BTCD_RPC_USER, BTCD_RPC_PASS)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	if (client != nil) {
		// low level shutdown since we shouldn't trust BitcoinCoin.Detach() yet
		client.Shutdown()
		client.WaitForShutdown()
	}
}

func TestBitcoinAttachDetach(t *testing.T) {
	coin := &BitcoinCoin{}

	err := coin.Attach()
	assert.NoError(t, err)
	err = coin.Detach()
	assert.NoError(t, err)
}

func TestBitcoinGetBlockId(t *testing.T) {
	coin := &BitcoinCoin{}
	coin.Attach()
	defer coin.Detach()

	count, err := coin.GetTopBlockID()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, count)
}

func TestBitcoinEncodeDecode(t *testing.T) {
	w := &Withdrawal{
		TxId:               42,
		CoinName:           "BTC",
		DestinationAddress: "123456789",
		QuantaBlockID:      1234,
		Amount:             12345,
		Signatures:         nil,
	}

	coin := &BitcoinCoin{}
	encoded, err := coin.EncodeRefund(*w)
	assert.NoError(t, err)
	w2, err := coin.DecodeRefund(encoded)
	assert.NoError(t, err)

	assert.Equal(t, w.TxId, w2.TxId)
	assert.Equal(t, w.CoinName, w2.CoinName)
	assert.Equal(t, w.DestinationAddress, w2.DestinationAddress)
	assert.Equal(t, w.QuantaBlockID, w2.QuantaBlockID)
	assert.Equal(t, w.Amount, w2.Amount)
	assert.Equal(t, w.Signatures, w2.Signatures)
}
