package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	metric2 "github.com/quantadex/distributed_quanta_bridge/common/metric"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/stretchr/testify/assert"
	"github.com/zserge/metric"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()
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
}

func TestAddress(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()
	time.Sleep(time.Millisecond * 250)

	address := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	address2 := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec888",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	address3 := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec889",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_ETH,
	}

	// test crosschain
	for _, n := range nodes {
		n.rDb.AddCrosschainAddress(address)
		n.rDb.AddCrosschainAddress(address2)
		control.SetLastBlock(n.db, coin.BLOCKCHAIN_ETH, 700000)
	}

	// wait for node to bootup
	time.Sleep(time.Millisecond * 2000)

	// test crosschain
	res, err := http.Get("http://localhost:5200/api/address/eth/pooja")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	res, err = http.Post("http://localhost:5200/api/address/eth/pooja", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	res, err = http.Post("http://localhost:5200/api/address/eth/alpha", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	fmt.Println("We're out of addresses, expect to fail")
	res, err = http.Post("http://localhost:5200/api/address/eth/charlie", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data 3", res.StatusCode, string(bodyBytes))
	assert.True(t, strings.Contains(string(bodyBytes), "Unable to agree for address change"))

	fmt.Println("*** TESTING FOR REPAIR ***")
	// show the address only for first 2 nodes, 3rd node will attempt to repair.
	for _, n := range nodes[0:2] {
		n.rDb.AddCrosschainAddress(address3)
	}

	res, err = http.Post("http://localhost:5200/api/address/eth/charlie", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	// expect repair

	// expect no repair
	for _, n := range nodes {
		address3.ContractAddress = "0xba420ef5d725361d8fdc58cb1e4fa62eda9ec111"
		n.rDb.AddCrosschainAddress(address3)
	}
	res, err = http.Post("http://localhost:5200/api/address/eth/beta", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
}

func TestStatus(t *testing.T) {
	fmt.Println(time.Now())
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()
	time.Sleep(time.Millisecond * 250)

	address := &crypto.ForwardInput{
		"hdfvh",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_LTC,
	}

	// test crosschain
	nodes[0].rDb.AddCrosschainAddress(address)
	//nodes[1].rDb.AddCrosschainAddress(address)

	address = &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"0x01",
		coin.BLOCKCHAIN_LTC,
	}

	nodes[0].rDb.AddCrosschainAddress(address)
	//nodes[1].rDb.AddCrosschainAddress(address)

	control.SetLastBlock(nodes[0].db, control.QUANTA, 9593474)
	//control.SetLastBlock(nodes[1].db, control.QUANTA, 9593474)

	time.Sleep(time.Millisecond * 2000)

	// test crosschain
	res, err := http.Get("http://localhost:5200/api/status")
	//res, err = http.Get("http://localhost:5201/api/status")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	//res, err = http.Get("http://localhost:5201/api/status")
	//assert.NoError(t, err)
	//bodyBytes, _ = ioutil.ReadAll(res.Body)
	//println("data", res.StatusCode, string(bodyBytes))

}

type Metric struct {
	M metric.Metric
}

func (s Metric) Print(i int) {
	s.M.Add(1)
	fmt.Println(s.M.String())
	var t *metric2.Counter
	json.Unmarshal([]byte(s.M.String()), &t)
	fmt.Println(t.Samples[0].Count)

}

func TestMetric(t *testing.T) {
	s := Metric{
		M: metric.NewCounter("1m"),
	}

	for i := 0; i < 10; i++ {
		if i == 10 {
			s.Print(i)
		} else {
			s.Print(0)
		}
	}
}

func TestAddressAllNodes(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()
	time.Sleep(time.Millisecond * 250)

	//btc
	res, err := http.Post("http://localhost:5200/api/address/btc/pooja", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	res, err = http.Get("http://localhost:5200/api/address/btc/pooja")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	//res, err = http.Get("http://localhost:5201/api/address/btc/pooja")
	//assert.NoError(t, err)
	//bodyBytes, _ = ioutil.ReadAll(res.Body)
	//println("data", res.StatusCode, string(bodyBytes))
	//assert.Equal(t, 200, res.StatusCode)

	//bch
	res, err = http.Post("http://localhost:5200/api/address/bch/pooja", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	res, err = http.Get("http://localhost:5200/api/address/bch/pooja")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	//res, err = http.Get("http://localhost:5201/api/address/bch/pooja")
	//assert.NoError(t, err)
	//bodyBytes, _ = ioutil.ReadAll(res.Body)
	//println("data", res.StatusCode, string(bodyBytes))
	//assert.Equal(t, 200, res.StatusCode)

	//ltc
	res, err = http.Post("http://localhost:5200/api/address/ltc/pooja", "", nil)
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	res, err = http.Get("http://localhost:5200/api/address/ltc/pooja")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))
	assert.Equal(t, 200, res.StatusCode)

	//res, err = http.Get("http://localhost:5201/api/address/ltc/pooja")
	//assert.NoError(t, err)
	//bodyBytes, _ = ioutil.ReadAll(res.Body)
	//println("data", res.StatusCode, string(bodyBytes))
	//assert.Equal(t, 200, res.StatusCode)

}

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func TestStressTest(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()

	// wait for node to bootup
	time.Sleep(time.Millisecond * 1000)

	for i := 0; i < 10; i++ {
		contractAddress := "0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999" + strconv.Itoa(i)
		address := &crypto.ForwardInput{
			contractAddress,
			common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
			"address-pool",
			"0x01",
			coin.BLOCKCHAIN_ETH,
		}
		for _, node := range nodes {
			node.rDb.AddCrosschainAddress(address)
		}
	}

	// test crosschain
	for _, n := range nodes {
		control.SetLastBlock(n.db, coin.BLOCKCHAIN_ETH, 700000)
	}

	//reducing to 10 to test on circleci
	resultChan := make(chan interface{}, 10)
	numTests := 10
	baseStr := RandomString(20)

	// test crosschain
	for i := 0; i < numTests; i++ {
		go func(i int) {
			str := baseStr + strconv.Itoa(i)
			res, err := http.Post("http://localhost:5200/api/address/eth/"+str, "", nil)

			if err != nil {
				resultChan <- err
			} else if res.StatusCode != 200 {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- errors.New("Status code:" + res.Status + " " + string(bodyBytes))
			} else {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- string(bodyBytes)
			}
		}(i)
	}

	for i := 0; i < numTests; i++ {
		res := <-resultChan
		switch v := res.(type) {
		case string:
			println("result ", v)
		case error:
			println("result ", v.Error())
		}
	}
}

func TestDuplicate(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()

	// wait for node to bootup
	time.Sleep(time.Millisecond * 1000)

	resultChan := make(chan interface{}, 3)
	numTests := 3
	baseStr := RandomString(20)

	// test crosschain
	for i := 0; i < numTests; i++ {
		go func(i int) {
			res, err := http.Post("http://localhost:5200/api/address/btc/"+baseStr, "", nil)

			if err != nil {
				resultChan <- err
			} else if res.StatusCode != 200 {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- errors.New("Status code:" + res.Status + " " + string(bodyBytes))
				assert.Contains(t, string(bodyBytes), "Duplicate address request")
			} else {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- string(bodyBytes)
			}
		}(i)
	}

	for i := 0; i < numTests; i++ {
		res := <-resultChan
		switch v := res.(type) {
		case string:
			println("result ", v)
		case error:
			println("result ", v.Error())
		}
	}
}

func TestRepair(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()

	// wait for node to bootup
	time.Sleep(time.Millisecond * 1000)
	addr, err := nodes[0].CreateMultisig("BTC", "pooja5")
	assert.NoError(t, err)
	addr2, err := nodes[0].CreateMultisig("BTC", "pooja7")
	assert.NoError(t, err)

	for _, n := range nodes {
		n.rDb.AddCrosschainAddress(addr)
		if n.nodeID != 2 {
			n.rDb.AddCrosschainAddress(addr2)
		}
	}

	resultChan := make(chan interface{}, 1)
	numTests := 1
	baseStr := RandomString(20)

	// test crosschain
	for i := 0; i < numTests; i++ {
		go func(i int) {
			str := baseStr + strconv.Itoa(i)
			res, err := http.Post("http://localhost:5200/api/address/bch/"+str, "", nil)

			if err != nil {
				fmt.Println("error = ", err)
				resultChan <- err
			} else if res.StatusCode != 200 {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- errors.New("Status code:" + res.Status + " " + string(bodyBytes))
			} else {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- string(bodyBytes)
			}
		}(i)
	}

	for i := 0; i < numTests; i++ {
		res := <-resultChan
		switch v := res.(type) {
		case string:
			println("result ", v)
		case error:
			println("result ", v.Error())
		}
	}
}

func TestVariationTimming(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()

	// wait for node to bootup
	time.Sleep(time.Millisecond * 1000)

	resultChan := make(chan interface{}, 1)
	numTests := 10
	baseStr := RandomString(20)

	// test crosschain
	for i := 0; i < numTests; i++ {
		go func(i int) {
			if i%5 == 0 {
				time.Sleep(time.Second * 15)
			} else if i%3 == 0 {
				time.Sleep(time.Second * 10)
			} else if i%2 == 0 {
				time.Sleep(time.Second * 5)
			}
			str := baseStr + strconv.Itoa(i)
			res, err := http.Post("http://localhost:5200/api/address/ltc/"+str, "", nil)

			if err != nil {
				resultChan <- err
			} else if res.StatusCode != 200 {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- errors.New("Status code:" + res.Status + " " + string(bodyBytes))
			} else {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				resultChan <- string(bodyBytes)
			}
		}(i)
	}

	for i := 0; i < numTests; i++ {
		res := <-resultChan
		switch v := res.(type) {
		case string:
			println("result ", v)
		case error:
			println("result ", v.Error())
		}
	}
}

func TestAddressToQuanta(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		time.Sleep(time.Second * 1)
	}()
	time.Sleep(time.Millisecond * 250)

	address := &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec9084",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"address-pool",
		"0x01",
		coin.BLOCKCHAIN_LTC,
	}

	nodes[0].rDb.AddCrosschainAddress(address)
	nodes[1].rDb.AddCrosschainAddress(address)

	// test crosschain
	nodes[0].rDb.AddCrosschainAddress(address)
	//nodes[1].rDb.AddCrosschainAddress(address)

	address = &crypto.ForwardInput{
		"0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999",
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"0x01",
		coin.BLOCKCHAIN_LTC,
	}

	nodes[0].rDb.AddCrosschainAddress(address)
	nodes[1].rDb.AddCrosschainAddress(address)

	control.SetLastBlock(nodes[0].db, control.QUANTA, 9593474)
	control.SetLastBlock(nodes[1].db, control.QUANTA, 9593474)

	time.Sleep(time.Second * 2)

	// test crosschain
	res, err := http.Get("http://localhost:5200/api/address/ltc")
	//res, err = http.Get("http://localhost:5201/api/status")
	assert.NoError(t, err)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

	res, err = http.Get("http://localhost:5201/api/address_to_quanta/ltc/0xba420ef5d725361d8fdc58cb1e4fa62eda9ec999")
	assert.NoError(t, err)
	bodyBytes, _ = ioutil.ReadAll(res.Body)
	println("data", res.StatusCode, string(bodyBytes))

}
