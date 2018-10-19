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
	"efff3cf1e98e8b348041c04fdc1f3019d0dae19d6f6e489bf8cfd38cb5270ddd",
	"b28ee83828f3e96f5d9048f866fe9b59e1c9b8a201fb3a71d8b19a3db9959249",
	"6aa210915b26e4f48f2f525ad14759e298bb98d2071fd8032149f33d5baff094",
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
EthereumNetworkId: 3
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

	//k1 , _ := crypto.HexToECDSA(ETHKEYS[0])
	//k2 , _ := crypto.HexToECDSA(ETHKEYS[1])
	//k3 , _ := crypto.HexToECDSA(ETHKEYS[2])

	//u1 := bind.NewKeyedTransactor(k1)
	//u2 := bind.NewKeyedTransactor(k2)
	//u3 := bind.NewKeyedTransactor(k3)

	ethereumClient, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		t.Error(err)
		return
	}

	//trustAddress, tx, contract, _ := contracts.DeployTrustContract(u1, ethereumClient)
	//println("trustaddr=", trustAddress.Hex(), tx.ChainId())
	//
	//time.Sleep(time.Second*3)
	//
	//tx, err  = contract.AssignInitialSigners(u1, []common.Address{
	//	u1.From, u2.From, u3.From,
	//})

	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//time.Sleep(time.Second*3)

	trustAddress := common.HexToAddress("0x384c0e4b22abfc546bdb84a5b259a82b351619b6")
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	if err != nil {
		t.Error(err)
		return
	}

	num, err := contract.GetTotalSigners(nil)
	if err != nil {
		println(err.Error())
	}

	println("Num of signers=", num)

	txId, err := contract.TxIdLast(nil)
	if err != nil {
		println(err.Error())
	}

	println("txID of signers=", txId)

	nodes := StartNodes(3, trustAddress)
	DoLoopWithdrawal(nodes, 0)

	time.Sleep(8 * time.Second)
}