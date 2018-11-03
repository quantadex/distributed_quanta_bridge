package main

import (
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"fmt"
	"testing"
	"sync"
	"os"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/stretchr/testify/assert"
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
		// 0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990
		"A7D7C6A92361590650AD0965970E186179F24F36B2B51CFE83F3AE8886BB6773",
		// 0xe0006458963c3773b051e767c5c63fee24cd7ff9
		"4C7F96D0CB8F2C48FD22CCB974513E6E9B0DC89475286BB24D2010E8D82AA461",
		// 0xba7573c0e805ef71acb7f1c4a55e7b0af416e96a
		"2E563A40747FA56419FB168ADF507C596E1A604D073D0F9E646B803DFA5BE94C",
	},
	TrustContract: "0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d",
}

const ROPSTEN = "ROPSTEN"
const LOCAL = "LOCAL"

// must match up with the HorizonUrl
const QUANTA_ACCOUNT = "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN"

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
		EthereumTrustAddr: ethereum.TrustContract,
	}
}

func assertMsgCountEqualDoLoop(t *testing.T, label string, expected int, actual int, blockNum int64, nodeNum int, totalNodes int, node *TrustNode) {
	assert.Equal(t, expected, actual, "%s message count was incorrect for block #%d [node #%d/%d id=%d]", label, blockNum, nodeNum, totalNodes, node.nodeID)
}

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
	for _, node := range nodes {
		node.cTQ.DoLoop(blockIds)
	}
}

func DoLoopWithdrawal(nodes []*TrustNode, cursor int64) {
	for _, node := range nodes {
		go node.qTC.DoLoop(cursor)
	}
}