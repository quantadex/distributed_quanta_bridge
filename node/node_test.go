package main

import (
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"fmt"
	"testing"
	"sync"
	"os"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/stretchr/testify/assert"
)

func generateConfig(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherNet test.EthereumEnv, index int) *common.Config {
	return &common.Config {
		ExternalListenPort: 5200+index,
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
		EthereumNetworkId: etherNet.NetworkId,
		EthereumBlockStart: 0,
		EthereumRpc: etherNet.Rpc,
		EthereumKeyStore: ethereum.NodeSecrets[index],
		EthereumTrustAddr: ethereum.TrustContract,
	}
}

func assertMsgCountEqualDoLoop(t *testing.T, label string, expected int, actual int, blockNum int64, nodeNum int, totalNodes int, node *TrustNode) {
	assert.Equal(t, expected, actual, "%s message count was incorrect for block #%d [node #%d/%d id=%d]", label, blockNum, nodeNum, totalNodes, node.nodeID)
}

func StartNodes(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherEnv test.EthereumEnv) []*TrustNode {
	println("Starting nodes")

	mutex := sync.Mutex{}

	nodes := []*TrustNode{}
	var wg sync.WaitGroup

	for i := 0; i < len(quanta.NodeSecrets); i++ {
		wg.Add(1)

		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))

		config := generateConfig(quanta, ethereum, etherEnv, i)

		go func(config common.Config) {
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