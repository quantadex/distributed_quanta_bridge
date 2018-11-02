package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/stretchr/testify/assert"
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
const QUANTA_ASSET = "0xAc2AFb5463F5Ba00a1161025C2ca0311748BfD2c"
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

/**
 * This one test native token from block 4186072
 * DATA DEPENDENT on ROPSTEN
 */
func TestRopstenNativeETH(t *testing.T) {
	r := StartRegistry()
	nodes := StartNodes(QUANTA_ISSUER, ROPSTEN_TRUST, ETHER_NETWORKS[ROPSTEN])
	time.Sleep(time.Millisecond*250)

	// DEPOSIT to TEST2
	block := int64(4248970)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// Check for the deposit
	block = 4249018
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 1, len(allDeposits), block, i+1, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4249019
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		// TODO: inspect the messages for the right content

		if i == 0 {
			// the first node doesn't receive any peer messages (yet), we will check for it in block #4249020
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		} else {
			// the rest of the nodes each receive a peer message
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		}
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// check for the signed peer messages
	block = 4249020
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		// TODO: inspect the messages for the right content
		if i == 0 {
			// the first node is the last to receive the peer message, and will go ahead and goes ahead and sends the submission message
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "sent", 1, len(allSentMsgs), block, i+1, len(nodes), node)

			// TODO: verify the content of that sent message
		} else {
			// the rest of the nodes have nothing to do
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		}
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// TODO: how to detect this successful tx submission handling?
	// S2018/11/01 16:04:04 I [5100] Successful tx submission 1ef7303e5f49d80feb0c4955a97d63336d79ae7734bd329d4d899f15db43a60d,remove 4249018ETHQDIX3EOMEWN7OLZ3BEIN5DE7MCVSAP6547FFM3FFITQSTFXWUK4XA2NB

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
	assert.NoError(t, err)

	fmt.Printf("[ASSET %s] [ACCOUNT %s] initial_balance = %.9f\n", QUANTA_ASSET, QUANTA_ACCOUNT, initialBalance)

	time.Sleep(time.Millisecond * 250)
	//DoLoopDeposit(nodes, []int64{4186072, 4186072, 4186074}) // we create the original smart contract on 74
	DoLoopDeposit(nodes, []int64{6,7}) // we make deposit
	time.Sleep(time.Second * 6)
	newBalance, err := nodes[0].q.GetBalance(QUANTA_ASSET, QUANTA_ACCOUNT)
	assert.NoError(t, err)
	assert.Equal(t, newBalance, initialBalance+0.0000001)
	//DoLoopDeposit(nodes, []int64{4196674})
	time.Sleep(time.Second * 15)
	StopNodes(nodes)
	StopRegistry(r)
}

func TestWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(ETHER_NETWORKS[ROPSTEN].rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(ROPSTEN_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry()

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)

	println("latest TXID=", txId)

  nodes := StartNodes(QUANTA_ISSUER, ROPSTEN_TRUST, ETHER_NETWORKS[ROPSTEN])

	cursor := int64(0)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, txId, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 28, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, "0x0", txId, "txId: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}
	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	// check for the Cosi!
	cursor = int64(1912892534296577)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, txId, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Zero(t, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)

		if (i == 0) {
			// first one should of submitted a withdrawal transaction
			assert.NotEqual(t, "0x0", txId, "txId: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		} else {
			assert.Equal(t, "0x0", txId, "txId: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		}
	}
	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	// check for the refund confirmations
	cursor = int64(1912892534296100)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, txId, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		// all three nodes should confirm a refund (from loop2)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, "0x0", txId, "txId: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}
	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	StopNodes(nodes)
	StopRegistry(r)
}
