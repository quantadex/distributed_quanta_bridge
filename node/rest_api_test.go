package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
)

func TestAPI(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	address := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"alpha",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	// test crosschain
	nodes[0].rDb.AddCrosschainAddress(address)
	res, err := http.Get("http://localhost:5200/api/address/eth/alpha")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)

	// test history
	for i := 0; i < 10; i++ {
		deposit := &coin.Deposit{
			Tx:         strconv.Itoa(i),
			CoinName:   "TEST",
			SenderAddr: "Sender",
			QuantaAddr: "To",
			Amount:     int64(i),
			BlockID:    int64(i),
		}
		db.ConfirmDeposit(nodes[0].rDb, deposit, false)
	}
	for i := 10; i < 15; i++ {
		deposit := &coin.Deposit{
			Tx:         strconv.Itoa(i),
			CoinName:   "TEST",
			SenderAddr: "Sender",
			QuantaAddr: "Joe",
			Amount:     int64(i),
			BlockID:    int64(i),
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
	assert.Equal(t, 5, len(addresses))
	assert.Equal(t, int64(14), addresses[0].BlockId)

	// get offset 5
	res, err = http.Get("http://localhost:5200/api/history?offset=5&limit=5")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &addresses)

	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, 5, len(addresses))
	assert.Equal(t, int64(9), addresses[0].BlockId)

	// get filter
	res, err = http.Get("http://localhost:5200/api/history?user=Joe&limit=15")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &addresses)

	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, 5, len(addresses))
	assert.Equal(t, int64(14), addresses[0].BlockId)

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func TestAddress(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	address := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	control.SetLastBlock(nodes[0].db, coin.BLOCKCHAIN_ETH, 700000)
	control.SetLastBlock(nodes[1].db, coin.BLOCKCHAIN_ETH, 700000)

	// test crosschain
	nodes[0].rDb.AddCrosschainAddress(address)
	nodes[1].rDb.AddCrosschainAddress(address)
	//nodes[0].rDb.UpdateLastBlockNumber("0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990", 1)
	//nodes[1].rDb.UpdateLastBlockNumber("0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990", 1)

	// wait for node to bootup
	time.Sleep(time.Millisecond * 2000)

	// test crosschain
	res, err := http.Get("http://localhost:5200/api/address/eth/pooja")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	res, err = http.Get("http://localhost:5201/api/address/eth/pooja")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)

}
