package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func generateConfig(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherNet test.EthereumEnv, index int) (*common.Config, *common.Secrets) {
	return &common.Config{
			ExternalListenPort: 5200 + index,
			ListenIp:           "0.0.0.0",
			ListenPort:         5100 + index,
			UsePrevKeys:        true,
			MinNodes:           2,

			KvDbName:           fmt.Sprintf("kv_db_%d", 5100+index),
			CoinMapping:        map[string]string{"BCH": "TESTISSUE8", "LTC": "TESTISSUE2", "BTC": "TESTISSUE3", "ETH": "TESTETH"},
			IssuerAddress:      quanta.SourceAccount,
			NetworkUrl:         "ws://testnet-01.quantachain.io:8090",
			ChainId:            "bb2aeb9eebaaa29d79ed81699ee49a912c19c59b9350f8f8d3d81b12fa178495",
			RegistrarIp:        "localhost",
			RegistrarPort:      6000,
			EthereumNetworkId:  etherNet.NetworkId,
			EthereumBlockStart: 0,
			EthereumRpc:        etherNet.Rpc,

			EthereumTrustAddr:    ethereum.TrustContract,
			Erc20Mapping:         map[string]string{strings.ToLower("0xDfE1002c2e1AE5E8F4f34bf481900dAae5351992"): "DAI"},
			EthMinConfirmation:   1,
			EthDegradedThreshold: 2000,
			EthFailureThreshold:  4000,
			MinBlockReuse:        43200,

			BtcRpc:               "localhost:18332",
			BtcNetwork:           "regnet",
			BtcMinConfirmation:   1,
			BtcDegradedThreshold: 2000,
			BtcFailureThreshold:  4000,

			LtcRpc:               "localhost:19332",
			LtcNetwork:           "regnet",
			LtcMinConfirmation:   1,
			LtcDegradedThreshold: 2000,
			LtcFailureThreshold:  4000,

			BchRpc:               "localhost:18333",
			BchNetwork:           "regnet",
			BchMinConfirmation:   1,
			BchDegradedThreshold: 2000,
			BchFailureThreshold:  4000,

			QuantaDegradedThreshold:   2000,
			QuantaFailureThreshold:    4000,
			DepDegradedThreshold:      10,
			DepFailureThreshold:       20,
			WithdrawDegradedThreshold: 10,
			WithdrawFailureThreshold:  20,
		}, &common.Secrets{
			NodeKey:          quanta.NodeSecrets[index],
			EthereumKeyStore: ethereum.NodeSecrets[index],
			DatabaseUrl:      fmt.Sprintf("postgres://postgres:@localhost/crosschain_%d", index),

			BtcPrivateKey:  test.BTCSECRETS.NodeSecrets[index],
			BtcRpcUser:     "user",
			BtcRpcPassword: "123",
			BtcSigners:     []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"},

			LtcPrivateKey:  test.LTCSECRETS.NodeSecrets[index],
			LtcRpcUser:     "user",
			LtcRpcPassword: "123",
			LtcSigners:     []string{"047AABB69BBE1B5D9E2EFD10D0215A37AE835EAE08DFDF795E5A8411271F690CC8797CF4DEB3508844920E28A42A67D8A3F56D5B6B65401DEDB1E130F9F9908463", "04851D591308AFBE768566060C01A60A5F6AC6C78C3766559C835BEF0485628013ADC7D7E7676B0281FB83E788F4BC11E4CA597D1A53AF5F0BB90D555A28B55504"},

			BchPrivateKey:  test.BCHSECRETS.NodeSecrets[index],
			BchRpcUser:     "user",
			BchRpcPassword: "123",
			BchSigners:     []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"},

			GrapheneSeedPrefix: "",
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
	println("Starting nodes with ", ethereum.TrustContract)

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

		config, secrets := generateConfig(quanta, ethereum, etherEnv, currentIndex)

		go func(config common.Config, currentIndex int) {
			defer wg.Done()

			mutex.Lock()
			node := bootstrapNode(config, *secrets, true)
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
	etherEnv test.EthereumEnv, numOfNodes int) []*TrustNode {
	nodes := make([]*TrustNode, numOfNodes)
	indexes := make([]int, numOfNodes)
	for i := range indexes {
		indexes[i] = i
	}
	return StartNodesWithIndexes(quanta, ethereum, etherEnv, true, indexes, nodes)
}

func StartNodeListener(quanta *test.QuantaNodeSecrets, ethereum *test.EthereumTrustSecrets,
	etherEnv test.EthereumEnv, nodes []*TrustNode) []*TrustNode {
	fmt.Println("\nStarting Node again")
	config, _ := generateConfig(quanta, ethereum, etherEnv, 0)
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

func StartRegistry(minNodes int, url string) *service.Server {
	logger, _ := logger.NewLogger("registrar")
	path, _ := filepath.Abs(filepath.Dir("config.yml"))
	s := service.NewServer(service.NewRegistry(minNodes, path), url, logger)
	s.DoHealthCheck(5)
	go s.Start()
	return s
}

func StopRegistry(s *service.Server) {
	path, _ := filepath.Abs(filepath.Dir("manifest.yml"))
	file := path + "/manifest.yml"
	os.Remove(file)
	s.Stop()
}

func DoLoopWithdrawal(nodes []*TrustNode, cursor int64) {
	for _, node := range nodes {
		go node.qTC.DoLoop(cursor)
	}
}
