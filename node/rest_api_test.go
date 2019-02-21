package main

import (
	"testing"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"time"
	"net/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"encoding/json"
)

func TestAPI(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	address := &coin.ForwardInput{
		common.HexToAddress( "0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990"),
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"alpha",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	// test crosschain
	db.AddCrosschainAddress(nodes[0].rDb, address)
	res, err := http.Get("http://localhost:5200/api/address/eth/alpha")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)

	// test history
	for i := 0; i < 10; i++ {
		deposit := &coin.Deposit{
			Tx: strconv.Itoa(i),
			CoinName: "TEST",
			SenderAddr: "Sender",
			QuantaAddr: "To",
			Amount: int64(i),
			BlockID: int64(i),
		}
		db.ConfirmDeposit(nodes[0].rDb, deposit, false)
	}
	for i := 10; i < 15; i++ {
		deposit := &coin.Deposit{
			Tx: strconv.Itoa(i),
			CoinName: "TEST",
			SenderAddr: "Sender",
			QuantaAddr: "Joe",
			Amount: int64(i),
			BlockID: int64(i),
		}
		db.ConfirmDeposit(nodes[0].rDb, deposit, false)
	}

	// get offset 0
	res, err = http.Get("http://localhost:5200/api/history?limit=5")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	var addresses []db.Transaction
	json.Unmarshal(bodyBytes, &addresses)

	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t,5,len(addresses))
	assert.Equal(t, int64(14), addresses[0].BlockId)

	// get offset 5
	res, err = http.Get("http://localhost:5200/api/history?offset=5&limit=5")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &addresses)

	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t,5,len(addresses))
	assert.Equal(t, int64(9), addresses[0].BlockId)

	// get filter
	res, err = http.Get("http://localhost:5200/api/history?user=Joe&limit=15")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &addresses)

	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t,5,len(addresses))
	assert.Equal(t, int64(14), addresses[0].BlockId)


	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}
