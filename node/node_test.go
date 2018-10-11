package main

import (
	"testing"
	"github.com/spf13/viper"
	"bytes"
	"fmt"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"sync"
	"os"
)

var NODE_KEYS = []string {
	"ZBYEUJIWP2AXG2V6ZW4F5OTM5APW3SOTTM6YGMKO6MQSY7U3IHFJZHWQ",
	"ZAFYSHEOQIK67O6S6SD5X7PVTLULQH3WQ3AMAGOO4NHSRM5SIKWCWFZB",
	"ZC4U5P5DWNXGRUENOCOKZFHAWFKBE7JFOB2BCEKCM7BKXXKQE3DARXIJ",
}

//address:QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN weight:1
//address:QCNKL7QKKQZD63UW27JLY7LDLR6MME3WNLUJ47VP25EZH5THRPEZRSAK weight:1
//address:QCN2DWLVXNAZW6ALR6KXJWGQB4J2J5TBJVPYLQMIU2TDCXIOBID5WRU5 weight:1
//address:QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB weight:1 *** Issuer

func SetConfig(key string, port int) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var config = []byte(fmt.Sprintf(`
ListenIp: 0.0.0.0
ListenPort: %d
UsePrevKeys: true
KvDbName: kv_db_%d
CoinName: ETH
IssuerAddress: QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB
NodeKey: %s
HorizonUrl: http://testnet-02.quantachain.io:8000/
NetworkPassphrase: QUANTA Test Network ; September 2018
RegistrarIp: localhost
RegistrarPort: 5001
EthereumNetworkId: 3
EthereumBlockStart: 4186070
EthereumRpc: https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6
EthereumTrustAddr: 0xe0006458963c3773B051E767C5C63FEe24Cd7Ff9
`, port, port, key))

	viper.ReadConfig(bytes.NewBuffer(config))
}

func StartNodes(n int)[]*TrustNode {
	nodes := []*TrustNode{}
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))
		SetConfig(NODE_KEYS[i], 5100 + i)
		config := Config {}
		err := viper.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		go func(config Config) {
			defer wg.Done()

			coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
			if err != nil {
				panic("Cannot create ethereum listener")
			}

			nodes = append(nodes, bootstrapNode(config, coin))
		}(config)

	}

	wg.Wait()

	return nodes
}

func StartRegistry() {
	logger, _ := logger.NewLogger("registrar")
	s := service.NewServer(service.NewRegistry(), "localhost:5001", logger)
	s.DoHealthCheck(5)
	go s.Start()
}

func DoLoop(nodes []*TrustNode, blockIds []int64) {
	for _, n := range nodes {
		n.cTQ.DoLoop(blockIds)
	}
}

func TestNode(t *testing.T) {
	StartRegistry()

	nodes := StartNodes(3)

	time.Sleep(time.Millisecond*250)

	//dummy := coin.GetDummyInstance()
	//dummy.CreateNewBlock()
	//
	// user generated ETH address from Quanta address
	// assume it is deposited to ETH address
	// 0x0f8d1c23a90795a7a738d90380ec8bb5e984ce9259b78ee7c5d1592253e4798c
	// with our system associating to QDCFARPB4ZR7VGTEL2XII5OPPUPPX2PQAYZURXRVR6Z34GNWTUHGVSXT
	//dummy.AddDeposit(&coin.Deposit{"ETH", "",
	//				"QDCFARPB4ZR7VGTEL2XII5OPPUPPX2PQAYZURXRVR6Z34GNWTUHGVSXT",
	//				15*10000000, 1})

	DoLoop(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	DoLoop(nodes, []int64{4196673})  // we make deposit
	DoLoop(nodes, []int64{4196674})

}