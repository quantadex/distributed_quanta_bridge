package main

import (
	"github.com/spf13/viper"
	"bytes"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"sync"
	"os"
	"testing"
	"time"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/ethereum/go-ethereum/common"
)

var NODE_KEYS = []string {
	"ZBYEUJIWP2AXG2V6ZW4F5OTM5APW3SOTTM6YGMKO6MQSY7U3IHFJZHWQ",
	"ZAFYSHEOQIK67O6S6SD5X7PVTLULQH3WQ3AMAGOO4NHSRM5SIKWCWFZB",
	"ZC4U5P5DWNXGRUENOCOKZFHAWFKBE7JFOB2BCEKCM7BKXXKQE3DARXIJ",
}

var ETHKEYS = []string {
	"c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3",
	"ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f",
	"0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1",
}

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
EthereumNetworkId: 1540234608622
EthereumBlockStart: 0
EthereumRpc: http://localhost:7545
EthereumKeyStore: %s
`, port, port, key, ethPrivKey))

	viper.ReadConfig(bytes.NewBuffer(config))
}

func StartNodes(n int, trustAddress common.Address)[]*TrustNode {
	println("Starting nodes with trust ", trustAddress.Hex())

	nodes := []*TrustNode{}
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))
		SetConfig(NODE_KEYS[i], 5100 + i, ETHKEYS[i])
		config := Config {}
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
	StartRegistry()
	nodes := StartNodes(3, common.HexToAddress("0xe0006458963c3773B051E767C5C63FEe24Cd7Ff9"))
	time.Sleep(time.Millisecond*250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	//DoLoopDeposit(nodes, []int64{4196673})  // we make deposit

	// DEPOSIT to TEST2
	DoLoopDeposit(nodes, []int64{4248970})
	DoLoopDeposit(nodes, []int64{4248971})  // we make deposit
	DoLoopDeposit(nodes, []int64{4249018})
	DoLoopDeposit(nodes, []int64{4249019})
	time.Sleep(time.Second*4)
}

func TestRopstenERC20Token(t *testing.T) {
	//StartRegistry()
	//nodes := StartNodes(3)
	//time.Sleep(time.Millisecond*250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	//DoLoopDeposit(nodes, []int64{4196673})  // we make deposit
	//DoLoopDeposit(nodes, []int64{4196674})
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
	StartRegistry()

	time.Sleep(time.Millisecond*250)

	ethereumClient, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		t.Error(err)
		return
	}

	trustAddress := common.HexToAddress("0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f")
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	if err != nil {
		t.Error(err)
		return
	}

	num, err := contract.TotalSigners(nil)
	if err != nil {
		println(err.Error())
	}

	println("Num of signers=", num.Uint64())

	txId, err := contract.TxIdLast(nil)
	if err != nil {
		println(err.Error())
	}

	println("txID of signers=", txId)

	nodes := StartNodes(3, trustAddress)
	DoLoopWithdrawal(nodes, 0)

	time.Sleep(8 * time.Second)
}