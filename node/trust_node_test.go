package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"os"
	"sync"
	"testing"
	"time"
)

type QuantaNodeSecrets struct {
	NodeSecrets []string
	SourceAccount string
}

type EthereumTrustSecrets struct {
	NodeSecrets []string
	TrustContract string
}

type EthereumEnv struct {
	rpc string
	networkId string
}

var QUANTA_ISSUER = &QuantaNodeSecrets{
	NodeSecrets:[]string{
		"ZBHK5VE5ZM5MJI3FM7JOW7MMUF3FIRUMV3BTLUTJWQHDFEN7MG3J4VAV",
		"ZDX6DGXBYAR3Z2BS4T4ITRTWPNJOSR5TPTVYN65UKEGP4ILOZ5GXU2KE",
		"ZC4U5P5DWNXGRUENOCOKZFHAWFKBE7JFOB2BCEKCM7BKXXKQE3DARXIJ",
	},
	SourceAccount: "QCISRUJ73RQBHB3C4LA6X537LPGSFZF3YUZ6MOPUOUJR5A63I5TLJML4",
}

var ROPSTEN_TRUST = &EthereumTrustSecrets {
	NodeSecrets: []string {
		"A7D7C6A92361590650AD0965970E186179F24F36B2B51CFE83F3AE8886BB6773",
		"4C7F96D0CB8F2C48FD22CCB974513E6E9B0DC89475286BB24D2010E8D82AA461",
		"2E563A40747FA56419FB168ADF507C596E1A604D073D0F9E646B803DFA5BE94C",
	},
	TrustContract: "0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d",
}

const ROPSTEN = "ROPSTEN"
const LOCAL = "LOCAL"

var ETHER_NETWORKS = map[string]EthereumEnv {
	ROPSTEN : EthereumEnv{ "https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6", "3" } ,
	LOCAL: EthereumEnv{ "http://localhost:7545", "10" },
}

func generateConfig(quanta *QuantaNodeSecrets, ethereum *EthereumTrustSecrets,
						etherNet EthereumEnv, index int) *Config {
	return &Config {
		ListenIp: "0.0.0.0",
		ListenPort: 5100+index,
		UsePrevKeys: true,
		KvDbName: fmt.Sprintf("kv_db_%d", 5100+index),
		CoinName: "ETH",
		IssuerAddress: quanta.SourceAccount,
		NodeKey: quanta.NodeSecrets[index],
		HorizonUrl: "http://testnet-02.quantachain.io:8000/",
		NetworkPassphrase: "QUANTA Test Network ; September 2018",
		RegistrarIp: "localhost",
		RegistrarPort: 5001,
		EthereumNetworkId: etherNet.networkId,
		EthereumBlockStart: 0,
		EthereumRpc: etherNet.rpc,
		EthereumKeyStore: ethereum.NodeSecrets[index],
	}
}

// must match up with the HorizonUrl
var QUANTA_ASSET = "0xAc2AFb5463F5Ba00a1161025C2ca0311748BfD2c"
var QUANTA_ACCOUNT = "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN"

func StartNodes(quanta *QuantaNodeSecrets, ethereum *EthereumTrustSecrets,
	etherEnv EthereumEnv) []*TrustNode {
	println("Starting nodes")

	mutex := sync.Mutex{}

	nodes := []*TrustNode{}
	var wg sync.WaitGroup

	for i := 0; i < len(quanta.NodeSecrets); i++ {
		wg.Add(1)

		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))

		config := generateConfig(quanta, ethereum, etherEnv, i)

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
		}(*config)

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
	nodes := StartNodes(QUANTA_ISSUER, ROPSTEN_TRUST, ETHER_NETWORKS[ROPSTEN])
	time.Sleep(time.Millisecond*250)

	// DEPOSIT to TEST2
	DoLoopDeposit(nodes, []int64{4248970})
	DoLoopDeposit(nodes, []int64{4249018}) // we make deposit
	DoLoopDeposit(nodes, []int64{4249019})
	DoLoopDeposit(nodes, []int64{4249020})
	time.Sleep(time.Second * 6)
	StopNodes(nodes)
	StopRegistry(r)
}

func TestRopstenERC20Token(t *testing.T) {
	//StartRegistry()
	//nodes := StartNodes(3)
	//time.Sleep(time.Millisecond*250)
	//DoLoopDeposit(nodes, []int64{4196673})  // we make deposit
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	//DoLoopDeposit(nodes, []int64{4196674})
	r := StartRegistry()
	nodes := StartNodes(QUANTA_ISSUER, ROPSTEN_TRUST, ETHER_NETWORKS[LOCAL])
	initialBalance, err := nodes[0].q.GetBalance(QUANTA_ASSET, QUANTA_ACCOUNT)
	assert.Nil(t, err)

	fmt.Printf("[ASSET %s] [ACCOUNT %s] initial_balance = %.9f\n", QUANTA_ASSET, QUANTA_ACCOUNT, initialBalance)

	time.Sleep(time.Millisecond * 250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	DoLoopDeposit(nodes, []int64{6,7}) // we make deposit
	time.Sleep(time.Second * 6)
	newBalance, err := nodes[0].q.GetBalance(QUANTA_ASSET, QUANTA_ACCOUNT)
	assert.Nil(t, err)
	assert.Equal(t, newBalance, initialBalance+0.0000001)
	//DoLoopDeposit(nodes, []int64{4196674})
	time.Sleep(time.Second * 15)
	StopNodes(nodes)
	StopRegistry(r)
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
	ethereumClient, err := ethclient.Dial(ETHER_NETWORKS[ROPSTEN].rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(ROPSTEN_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.Nil(t, err)

	r := StartRegistry()

	txId, err := contract.TxIdLast(nil)
	assert.Nil(t, err)

	println("latest TXID=", txId)

	nodes := StartNodes(QUANTA_ISSUER, ROPSTEN_TRUST, ETHER_NETWORKS[ROPSTEN])
	DoLoopWithdrawal(nodes, 0)
	DoLoopWithdrawal(nodes, 1912892534296577)
	DoLoopWithdrawal(nodes, 1912892534296100)

	time.Sleep(15 * time.Second)
	StopNodes(nodes)
	StopRegistry(r)
}
