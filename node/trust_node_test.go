package main

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gcash/bchutil"
	"github.com/ltcsuite/ltcutil"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os/exec"
	"testing"
	"time"

	chaincfg3 "github.com/gcash/bchd/chaincfg"
	chaincfg2 "github.com/ltcsuite/ltcd/chaincfg"
)

/*
 * INTEGRATED TESTING between ROPSTEN & 3-node setup
 */

/**
 * This one test native token from block 4186072
 * DATA DEPENDENT on ROPSTEN
 */

func GetEthSync(node *TrustNode, ethFlush bool, minConfirm int64) sync.DepositSyncInterface {
	return sync.NewEthereumSync(node.eth,
		test.GRAPHENE_TRUST.TrustContract,
		map[string]string{"eth": "TESTETH"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0,
		ethFlush,
		minConfirm)
}

func GetBtcSync(node *TrustNode, minConfirm int64) sync.DepositSyncInterface {
	return sync.NewBitcoinSync(node.btc,
		map[string]string{"btc": "TESTISSUE3"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0,
		minConfirm)
}

func GetLtcSync(node *TrustNode, minConfirm int64) sync.DepositSyncInterface {
	return sync.NewLitecoinSync(node.ltc,
		map[string]string{"ltc": "TESTISSUE2"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0,
		minConfirm)
}

func GetBchSync(node *TrustNode, minConfirm int64) sync.DepositSyncInterface {
	return sync.NewBCHSync(node.bch,
		map[string]string{"bch": "TESTISSUE8"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0,
		minConfirm)
}

func TestRopstenNativeETH(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	time.Sleep(time.Millisecond * 250)

	depositResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		depositResult <- c
	}
	config := generateConfig(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 0)

	// DEPOSIT to TEST2
	// 0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A transfers 0.01 ETH to forward address: 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 on block 4327101
	// Foward contract 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 created on 4327057
	block := int64(5066807)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// Check for the deposit
	block = int64(5413358)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 1, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = int64(5413359)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(time.Second * 5)

	var w *control.DepositResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-depositResult:
		w = &w_
	}

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	time.Sleep(time.Second * 5)

	// TODO: how to detect this successful tx submission handling?
	// S2018/11/01 16:04:04 I [5100] Successful tx submission 1ef7303e5f49d80feb0c4955a97d63336d79ae7734bd329d4d899f15db43a60d,remove 4249018ETHQDIX3EOMEWN7OLZ3BEIN5DE7MCVSAP6547FFM3FFITQSTFXWUK4XA2NB

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

// Block 4354954: ERC-20  0x719791a052ed86360015e659542d3b0b7f44182e created by 0x9b93a4be348ab36c16ec6861602f84321f24e544  tx=https://ropsten.etherscan.io/tx/0xfc50ed99430174aa791ea9fa0883ef12d1ba7aa3a2daf26f3bb255b8f5f9af1b
// Block 4354971: Forward contract created  0x527370DEF157BD1113DB9448BD05E6402fFb5A0d forwarding to trust contract from 0x9b93a4be348ab36c16ec6861602f84321f24e544
// Block 4355067: Nothing - get block to agree

// Block 4356004 ERC-20 0x541d973a7168dbbf413eab6993a5e504ec5accb0
// Block 4356013 sent .0001234  precision 9  tx=https://ropsten.etherscan.io/tx/0x51a1018c6b2afd7bffb52d05178c3b66c9336e8c2dfeca0110f405cc41613492
func TestRopstenERC20Token(t *testing.T) {
	r := StartRegistry(2, ":6000")
	//ercContract := "0x541d973a7168dbbf413eab6993a5e504ec5accb0"
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	initialBalance, err := nodes[0].q.GetBalance("SIMPLETOKEN0XDFE1002C2E1AE5E8F4F34BF481900DAAE5351992", "pooja")
	assert.NoError(t, err)

	fmt.Printf("[ASSET %s] [ACCOUNT %s] initial_balance = %.9f\n", "SIMPLETOKEN0XDFE1002C2E1AE5E8F4F34BF481900DAAE5351992", "pooja", initialBalance)

	time.Sleep(time.Millisecond * 250)
	config := generateConfig(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 0)

	block := int64(5066807)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// Check for the deposit
	block = int64(5066833)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 1, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = int64(5066834)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		var ethSync sync.DepositSyncInterface
		if i == 0 {
			ethSync = GetEthSync(node, true, config.EthMinConfirmation)
		} else {
			ethSync = GetEthSync(node, false, config.EthMinConfirmation)
		}

		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(time.Second * 6)
	newBalance, err := nodes[0].q.GetBalance("SIMPLETOKEN0XDFE1002C2E1AE5E8F4F34BF481900DAAE5351992", "pooja")
	assert.NoError(t, err)
	//assert.Equal(t, initialBalance+float64(0.1), newBalance)
	fmt.Printf("Initial balance=%f , expecting final balance = %f\n", initialBalance, newBalance)
	time.Sleep(time.Second * 5)
	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

//func TestTrustNode_Stop(t *testing.T) {
//	r := StartRegistry()
//	//indexes := []int{0,1,2}
//	//nodes := []*TrustNode{}
//	//nodes = StartNodesWithIndexes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN],true,indexes, nodes)
//
//	nodes := StartNodes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
//	time.Sleep(time.Millisecond * 250)
//
//	block := int64(4327057)
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
//			fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
//		}
//
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	block = 4327101
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
//			fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
//		}
//
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	block = 4327102
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
//			fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//
//			if int(4327101-1)%3 == i {
//				assertMsgCountEqualDoLoop(t, "sent", 1, len(allSentMsgs), block, i+1, len(nodes), node)
//			} else {
//				assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
//			}
//		}
//	}
//
//	//StopNodeListener(nodes[0])
//	StopNodes(nodes, []int{2})
//	time.Sleep(time.Second * 6)
//
//	block = int64(4327057)
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
//			fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
//			assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
//		}
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	block = 4327127
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
//			fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
//			assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
//		}
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	block = 4327128
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		if node != nil {
//			fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
//			allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
//			fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
//			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
//			if (4327128-1)%3 == i {
//				assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
//			} else {
//				assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
//			}
//		}
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	time.Sleep(time.Second * 6)
//
//	indexToStart := []int{2}
//	//nodes = StartNodeListener(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], nodes)
//	nodes = StartNodesWithIndexes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], false, indexToStart, nodes)
//	time.Sleep(time.Second * 6)
//	block = int64(4327057)
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
//		allDeposits := node.cTQ.DoLoop([]int64{block})
//		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits))
//		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//	block = 4327101
//	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
//	for i, node := range nodes {
//		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
//		allDeposits := node.cTQ.DoLoop([]int64{block})
//	}
//	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)
//
//
//	time.Sleep(time.Second * 6)
//	StopNodes(nodes, []int{0, 1, 2})
//	StopRegistry(r)
//}

func TestWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(test.ETHER_NETWORKS[test.ROPSTEN].Rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(test.GRAPHENE_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry(2, ":6000")

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)
	println("latest TXID=", txId)

	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	withdrawResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		withdrawResult <- c
	}

	cursor := int64(5373523)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}

	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	time.Sleep(time.Second * 6)

	var w *control.DepositResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-withdrawResult:
		w = &w_
	}

	fmt.Println("w = ", w)

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func TestBCHDeposit(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	depositResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		depositResult <- c
	}

	config := generateConfig(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 0)

	client, err := coin.NewBCHCoin("localhost:18332", &chaincfg3.RegressionNetParams, []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"})
	err = client.Attach()
	assert.NoError(t, err)

	msig, err := client.GenerateMultisig("pooja")
	assert.NoError(t, err)

	forwardAddress := &crypto.ForwardInput{
		msig,
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"",
		coin.BLOCKCHAIN_BCH,
	}
	nodes[0].rDb.AddCrosschainAddress(forwardAddress)
	nodes[1].rDb.AddCrosschainAddress(forwardAddress)

	assert.NoError(t, err)

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/BCH/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	address := string(bodyBytes)[13:62]
	err = ImportAddress(address, "bitcoin-cli")
	assert.NoError(t, err)

	amount, _ := bchutil.NewAmount(0.01)
	SendBCH(address, amount)
	GenerateBlock("bitcoin-cli")

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block := int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ltcSync := GetBchSync(node, config.LtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ltcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	GenerateBlock("bitcoin-cli")

	blockId, err = client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block = int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ltcSync := GetBchSync(node, config.LtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ltcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(10 * time.Second)

	var w *control.DepositResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-depositResult:
		w = &w_
	}

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	time.Sleep(time.Second * 5)

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func TestBCHWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(test.ETHER_NETWORKS[test.ROPSTEN].Rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(test.GRAPHENE_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry(2, ":6000")

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)
	println("latest TXID=", txId)

	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	withdrawResult := make(chan control.WithdrawalResult)

	nodes[0].qTC.SuccessCb = func(c control.WithdrawalResult) {
		withdrawResult <- c
	}

	client, err := coin.NewBCHCoin("localhost:18332", &chaincfg3.RegressionNetParams, []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"})
	err = client.Attach()
	assert.NoError(t, err)

	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)
	msig, err := client.GenerateMultisig(btec)

	forwardAddress := &crypto.ForwardInput{
		msig,
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"",
		coin.BLOCKCHAIN_BCH,
	}
	nodes[0].rDb.AddCrosschainAddress(forwardAddress)
	nodes[1].rDb.AddCrosschainAddress(forwardAddress)

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/BCH/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	amount, err := bchutil.NewAmount(0.9)
	address := string(bodyBytes)[13:62]
	err = ImportAddress(address, "bitcoin-cli")
	assert.NoError(t, err)

	SendBCH(address, amount)
	GenerateBlock("bitcoin-cli")

	cursor := int64(8529695)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}

	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	time.Sleep(time.Second * 8)

	var w *control.WithdrawalResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-withdrawResult:
		w = &w_
	}

	fmt.Println("w = ", w)

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func SendBCH(address string, amount bchutil.Amount) (string, error) {
	amountStr := fmt.Sprintf("%f", amount.ToBCH())
	fmt.Printf("Sending to %s amount of %s\n", address, amountStr)
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"sendtoaddress",
		address,
		amountStr,
	}

	cmd := exec.Command("bitcoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func SendLTC(address string, amount ltcutil.Amount) (string, error) {
	amountStr := fmt.Sprintf("%f", amount.ToBTC())
	fmt.Printf("Sending to %s amount of %s\n", address, amountStr)
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"sendtoaddress",
		address,
		amountStr,
	}

	cmd := exec.Command("litecoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func TestLTCDeposit(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	depositResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		depositResult <- c
	}

	config := generateConfig(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 0)

	client, err := coin.NewLitecoinCoin("localhost:19332", &chaincfg2.RegressionNetParams, []string{"047AABB69BBE1B5D9E2EFD10D0215A37AE835EAE08DFDF795E5A8411271F690CC8797CF4DEB3508844920E28A42A67D8A3F56D5B6B65401DEDB1E130F9F9908463", "04851D591308AFBE768566060C01A60A5F6AC6C78C3766559C835BEF0485628013ADC7D7E7676B0281FB83E788F4BC11E4CA597D1A53AF5F0BB90D555A28B55504"})
	err = client.Attach()
	assert.NoError(t, err)

	msig, err := client.GenerateMultisig("pooja")
	assert.NoError(t, err)

	forwardAddress := &crypto.ForwardInput{
		msig,
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"",
		coin.BLOCKCHAIN_LTC,
	}
	nodes[0].rDb.AddCrosschainAddress(forwardAddress)
	nodes[1].rDb.AddCrosschainAddress(forwardAddress)

	assert.NoError(t, err)

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/LTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	address := string(bodyBytes)[13:47]
	err = ImportAddress(address, "litecoin-cli")
	assert.NoError(t, err)

	amount, _ := ltcutil.NewAmount(0.01)
	SendLTC(address, amount)
	GenerateBlock("litecoin-cli")

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block := int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ltcSync := GetLtcSync(node, config.LtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ltcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	GenerateBlock("litecoin-cli")

	blockId, err = client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block = int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ltcSync := GetLtcSync(node, config.LtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ltcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(10 * time.Second)

	var w *control.DepositResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-depositResult:
		w = &w_
	}

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	time.Sleep(time.Second * 5)

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func TestLTCWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(test.ETHER_NETWORKS[test.ROPSTEN].Rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(test.GRAPHENE_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry(2, ":6000")

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)
	println("latest TXID=", txId)

	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	withdrawResult := make(chan control.WithdrawalResult)

	nodes[0].qTC.SuccessCb = func(c control.WithdrawalResult) {
		withdrawResult <- c
	}

	client, err := coin.NewLitecoinCoin("localhost:19332", &chaincfg2.RegressionNetParams, []string{"047AABB69BBE1B5D9E2EFD10D0215A37AE835EAE08DFDF795E5A8411271F690CC8797CF4DEB3508844920E28A42A67D8A3F56D5B6B65401DEDB1E130F9F9908463", "04851D591308AFBE768566060C01A60A5F6AC6C78C3766559C835BEF0485628013ADC7D7E7676B0281FB83E788F4BC11E4CA597D1A53AF5F0BB90D555A28B55504"})
	err = client.Attach()
	assert.NoError(t, err)

	msig, err := client.GenerateMultisig("pooja")
	assert.NoError(t, err)

	_, err = client.GenerateMultisig("pooja2")
	assert.NoError(t, err)

	forwardAddress := &crypto.ForwardInput{
		msig,
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"",
		coin.BLOCKCHAIN_LTC,
	}
	nodes[0].rDb.AddCrosschainAddress(forwardAddress)
	nodes[1].rDb.AddCrosschainAddress(forwardAddress)

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/LTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	amount, err := ltcutil.NewAmount(0.9)
	address := string(bodyBytes)[13:47]
	err = ImportAddress(address, "litecoin-cli")
	assert.NoError(t, err)

	SendLTC(address, amount)
	GenerateBlock("litecoin-cli")

	cursor := int64(8348073)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}

	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	time.Sleep(time.Second * 8)

	var w *control.WithdrawalResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-withdrawResult:
		w = &w_
	}

	fmt.Println("w = ", w)

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func TestBTCDeposit(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	depositResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		depositResult <- c
	}
	config := generateConfig(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 0)

	client, err := coin.NewBitcoinCoin("localhost:18332", &chaincfg.RegressionNetParams, []string{"2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", "2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	msig, err := client.GenerateMultisig("pooja")
	assert.NoError(t, err)

	_, err = client.GenerateMultisig("pooja2")
	assert.NoError(t, err)

	forwardAddress := &crypto.ForwardInput{
		msig,
		common.HexToAddress(test.GRAPHENE_TRUST.TrustContract),
		"pooja",
		"",
		coin.BLOCKCHAIN_BTC,
	}
	nodes[0].rDb.AddCrosschainAddress(forwardAddress)
	nodes[1].rDb.AddCrosschainAddress(forwardAddress)

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/BTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	address := string(bodyBytes)[13:48]
	err = ImportAddress(address, "bitcoin-cli")
	assert.NoError(t, err)

	amount, _ := btcutil.NewAmount(0.01)
	SendBTC(address, amount)
	GenerateBlock("bitcoin-cli")

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block := int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		btcSync := GetBtcSync(node, config.BtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := btcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	GenerateBlock("bitcoin-cli")

	blockId, err = client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)

	block = int64(blockId)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		btcSync := GetBtcSync(node, config.BtcMinConfirmation)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := btcSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(10 * time.Second)

	var w *control.DepositResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-depositResult:
		w = &w_
	}

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	time.Sleep(time.Second * 5)

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}

func ImportAddress(address string, command string) error {
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"importaddress",
		address,
	}

	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
		return err
	}
	return nil
}

func SendBTC(address string, amount btcutil.Amount) (string, error) {
	amountStr := fmt.Sprintf("%f", amount.ToBTC())
	fmt.Printf("Sending to %s amount of %s\n", address, amountStr)
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"sendtoaddress",
		address,
		amountStr,
	}

	cmd := exec.Command("bitcoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func GenerateBlock(command string) (string, error) {
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"generate",
		"1",
	}

	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func TestBTCWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(test.ETHER_NETWORKS[test.ROPSTEN].Rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(test.GRAPHENE_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry(2, ":6000")

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)
	println("latest TXID=", txId)

	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	withdrawResult := make(chan control.WithdrawalResult)

	nodes[0].qTC.SuccessCb = func(c control.WithdrawalResult) {
		withdrawResult <- c
	}

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/BTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	amount, err := btcutil.NewAmount(0.9)
	address := string(bodyBytes)[13:48]
	err = ImportAddress(address, "bitcoin-cli")
	assert.NoError(t, err)

	SendBTC(address, amount)
	GenerateBlock("bitcoin-cli")

	cursor := int64(5116140)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}

	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	time.Sleep(time.Second * 8)

	var w *control.WithdrawalResult
	select {
	case <-time.After(time.Second * 8):
		w = nil
	case w_ := <-withdrawResult:
		w = &w_
	}

	fmt.Println("w = ", w)

	assert.NotNil(t, w, "We expect withdrawal completed")
	assert.NoError(t, w.Err, "should not get an error")

	StopNodes(nodes, []int{0, 1})
	StopRegistry(r)
}
