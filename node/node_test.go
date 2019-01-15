package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func generateConfig(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherNet test.EthereumEnv, index int) *common.Config {
	return &common.Config{
		ExternalListenPort: 5200 + index,
		ListenIp:           "0.0.0.0",
		ListenPort:         5100 + index,
		UsePrevKeys:        true,
		KvDbName:           fmt.Sprintf("kv_db_%d", 5100+index),
		CoinName:           "ETH",
		IssuerAddress:      quanta.SourceAccount,
		NodeKey:            quanta.NodeSecrets[index],
		NetworkUrl:         "ws://testnet-01.quantachain.io:8090",
		ChainId:            "bb2aeb9eebaaa29d79ed81699ee49a912c19c59b9350f8f8d3d81b12fa178495",
		RegistrarIp:        "localhost",
		RegistrarPort:      5001,
		EthereumNetworkId:  etherNet.NetworkId,
		EthereumBlockStart: 0,
		EthereumRpc:        etherNet.Rpc,
		EthereumKeyStore:   ethereum.NodeSecrets[index],
		EthereumTrustAddr:  ethereum.TrustContract,
		DatabaseUrl:        fmt.Sprintf("postgres://postgres:@localhost/crosschain_%d", index),
	}
}

func assertMsgCountEqualDoLoop(t *testing.T, label string, expected int, actual int, blockNum int64, nodeNum int, totalNodes int, node *TrustNode) {
	assert.Equal(t, expected, actual, "%s message count was incorrect for block #%d [node #%d/%d id=%d]", label, blockNum, nodeNum, totalNodes, node.nodeID)
}

/**
 * StartNodesWithIndexes starts the nodes with indexesToStart, which indexes into the quanta/ethereum child keys.
 * Indexes of nodes []*TrustNode should always be the same as the index of the quanta/ethereum index.
 * We can assume that nodes[]*TrustNode is pre-allocated with len(quanta child keys)
 * nodes[]*TrustNode are not modified
 */
func StartNodesWithIndexes(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherEnv test.EthereumEnv, removePrevDB bool, indexesToStart []int, nodesIn []*TrustNode) []*TrustNode {
	println("Starting nodes")

	nodes := make([]*TrustNode, 2)
	copy(nodes, nodesIn)

	mutex := sync.Mutex{}
	var wg sync.WaitGroup

	for i := 0; i < len(indexesToStart); i++ {
		currentIndex := indexesToStart[i]

		if nodes[currentIndex] != nil {
			println("Error: Node is already started!")
			continue
		}

		wg.Add(1)

		if removePrevDB {
			os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+currentIndex))
		}

		config := generateConfig(quanta, ethereum, etherEnv, currentIndex)

		go func(config common.Config, currentIndex int) {
			defer wg.Done()

			coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
			if err != nil {
				panic("Cannot create ethereum listener")
			}

			mutex.Lock()
			node := bootstrapNode(config, coin)
			nodes[currentIndex] = node
			mutex.Unlock()

			db.EmptyTable(node.rDb)
			registerNode(config, node)

			// ensure they start on time
			time.Sleep(time.Millisecond * 250)
		}(*config, currentIndex)

		time.Sleep(time.Second)
	}
	wg.Wait()
	return nodes
}

func StartNodes(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherEnv test.EthereumEnv) []*TrustNode {
	nodes := make([]*TrustNode, 2)
	return StartNodesWithIndexes(quanta, ethereum, etherEnv, true, []int{0, 1}, nodes)
}

func StartNodeListener(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherEnv test.EthereumEnv, nodes []*TrustNode) []*TrustNode {
	fmt.Println("\nStarting Node again")
	config := generateConfig(quanta, ethereum, etherEnv, 0)
	nodes[0].StartListener(*config)

	return nodes
}

func StopNodeListener(node *TrustNode) {
	fmt.Println("\nStopping node")
	node.StopListener()
}

func StopNodes(nodes []*TrustNode, indexesToStart []int) {
	fmt.Println("Stopping Nodes")
	for _, n := range indexesToStart {
		nodes[n].Stop()
		nodes[n] = nil
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
