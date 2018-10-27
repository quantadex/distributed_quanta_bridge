package main

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/magiconair/properties/assert"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/spf13/viper"
	"os"
	"sync"
	"testing"
	"time"
)

const ropsten_infura = "https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6"

var NODE_KEYS = []string {
	"ZBHK5VE5ZM5MJI3FM7JOW7MMUF3FIRUMV3BTLUTJWQHDFEN7MG3J4VAV",
	"ZDX6DGXBYAR3Z2BS4T4ITRTWPNJOSR5TPTVYN65UKEGP4ILOZ5GXU2KE",
	"ZCKSY3ZJ4KW4M5QZ4B5BUG2JCCA7JUUBCCFOIODEOYVPDHRPMP6AQ73V",
}

var ETHKEYS = []string {
	"c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3",
	"ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f",
	"0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1",
}

var ETH_TRUST_ADDRESS = "0xe0006458963c3773B051E767C5C63FEe24Cd7Ff9"

// var ETH_TRUST_ADDRESS = "0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d"

// var ETHKEYS = []string {
// 	"A7D7C6A92361590650AD0965970E186179F24F36B2B51CFE83F3AE8886BB6773",
// 	"4C7F96D0CB8F2C48FD22CCB974513E6E9B0DC89475286BB24D2010E8D82AA461",
// 	"2E563A40747FA56419FB168ADF507C596E1A604D073D0F9E646B803DFA5BE94C",
// }


//address:QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN weight:1
//address:QCNKL7QKKQZD63UW27JLY7LDLR6MME3WNLUJ47VP25EZH5THRPEZRSAK weight:1
//address:QCN2DWLVXNAZW6ALR6KXJWGQB4J2J5TBJVPYLQMIU2TDCXIOBID5WRU5 weight:1
//address:QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB weight:1 *** Issuer

func SetConfig(key string, port int, ethPrivKey string) {
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
EthereumBlockStart: 0
EthereumRpc: %s
EthereumKeyStore: %s
HEALTH_INTERVAL: 5
`, port, port, key, ropsten_infura, ethPrivKey))

	viper.ReadConfig(bytes.NewBuffer(config))
}

func StartNodes(n int, trustAddress common.Address) []*TrustNode {
	println("Starting nodes with trust ", trustAddress.Hex())

	mutex := sync.Mutex{}

	nodes := []*TrustNode{}
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))
		SetConfig(NODE_KEYS[i], 5100+i, ETHKEYS[i])
		config := Config{}
		err := viper.Unmarshal(&config)
		config.EthereumTrustAddr = trustAddress.Hex()

		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		go func(config Config) {
			defer wg.Done()

			coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
			if err != nil {
				panic("Cannot create ethereum listener")
			}

			mutex.Lock()
			node := bootstrapNode(config, coin)
			nodes = append(nodes, node)
			mutex.Unlock()

			registerNode(config, node)
		}(config)

	}

	wg.Wait()

	return nodes
}

func StopNodes(nodes []*TrustNode) {

	for _, n := range nodes {
		n.Stop()
	}
}

func StartRegistry() *service.Server {
	logger, _ := logger.NewLogger("registrar")
	s := service.NewServer(service.NewRegistry(), "localhost:5001", logger)
	s.DoHealthCheck(5)
	go s.Start()
	return s
}

func StopRegistry(s *service.Server) {
	s.Stop()
}

func DoLoopDeposit(nodes []*TrustNode, blockIds []int64) {
	for _, n := range nodes {
		n.cTQ.DoLoop(blockIds)
	}
}

func DoLoopWithdrawal(nodes []*TrustNode, cursor int64) {
	for _, n := range nodes {
		go n.qTC.DoLoop(cursor)
	}
}

/**
 * This one test native token from block 4186072
 */
func TestRopstenNativeETH(t *testing.T) {
	r := StartRegistry()
	nodes := StartNodes(3, common.HexToAddress(ETH_TRUST_ADDRESS))
	time.Sleep(time.Millisecond*250)

	// DEPOSIT to TEST2
	DoLoopDeposit(nodes, []int64{4248970})
	DoLoopDeposit(nodes, []int64{4249018}) // we make deposit
	DoLoopDeposit(nodes, []int64{4249019})
	time.Sleep(time.Second * 4)
	StopNodes(nodes)
	StopRegistry(r)
}

func TestRopstenERC20Token(t *testing.T) {
	//StartRegistry()
	//nodes := StartNodes(3)
	//time.Sleep(time.Millisecond*250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	//DoLoopDeposit(nodes, []int64{4196673})  // we make deposit
	//DoLoopDeposit(nodes, []int64{4196674})
	StartRegistry()
	nodes := StartNodes(3, common.HexToAddress("0xb1E02e31c9A2403FeAFA7E483Ebb3e1b5ffa3164"))
	initialBalance, _ := nodes[0].q.GetBalance("0xAc2AFb5463F5Ba00a1161025C2ca0311748BfD2c", "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN")
	fmt.Println("initial balance is: ", initialBalance)
	time.Sleep(time.Millisecond * 250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	DoLoopDeposit(nodes, []int64{4250980}) // we make deposit
	time.Sleep(time.Second * 6)
	newBalance, _ := nodes[0].q.GetBalance("0xAc2AFb5463F5Ba00a1161025C2ca0311748BfD2c", "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN")
	assert.Equal(t, newBalance, initialBalance+0.0000001)
	//DoLoopDeposit(nodes, []int64{4196674})
	time.Sleep(time.Second * 15)
}

func TestDummyCoin(t *testing.T) {

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
}

func TestWithdrawal(t *testing.T) {
	r := StartRegistry()

	ethereumClient, err := ethclient.Dial(ropsten_infura)
	if err != nil {
		t.Error(err)
		return
	}

	trustAddress := common.HexToAddress(ETH_TRUST_ADDRESS)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	if err != nil {
		t.Error(err)
		return
	}

	//num, err := contract.TotalSigners(nil)
	//if err != nil {
	//	println(err.Error())
	//}
	//
	//println("Num of signers=", num.Uint64())

	txId, err := contract.TxIdLast(nil)
	if err != nil {
		println(err.Error())
	}

	println("txID of signers=", txId)

	nodes := StartNodes(3, trustAddress)
	DoLoopWithdrawal(nodes, 0)

	time.Sleep(8 * time.Second)
	StopNodes(nodes)
	StopRegistry(r)
}
