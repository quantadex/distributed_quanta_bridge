package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

/*
 * INTEGRATED TESTING between ROPSTEN & 3-node setup
 */

/**
 * This one test native token from block 4186072
 * DATA DEPENDENT on ROPSTEN
 */

func GetEthSync(node *TrustNode) sync.DepositSyncInterface {
	return sync.NewEthereumSync(node.eth,
		test.GRAPHENE_TRUST.TrustContract,
		map[string]string{"eth": "TESTETH"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0)
}

func GetBtcSync(node *TrustNode) sync.DepositSyncInterface {
	return sync.NewBitcoinSync(node.btc,
		map[string]string{"btc": "TESTISSUE3"},
		node.q,
		node.db,
		node.rDb,
		node.log,
		0)
}
func TestRopstenNativeETH(t *testing.T) {
	r := StartRegistry(2, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

	time.Sleep(time.Millisecond * 250)

	depositResult := make(chan control.DepositResult)

	nodes[0].cTQ.SuccessCb = func(c control.DepositResult) {
		depositResult <- c
	}

	// DEPOSIT to TEST2
	// 0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A transfers 0.01 ETH to forward address: 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 on block 4327101
	// Foward contract 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 created on 4327057
	block := int64(5061200)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ethSync := GetEthSync(node)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// Check for the deposit
	block = int64(5061248)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ethSync := GetEthSync(node)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 1, len(allDeposits), block, i+1, len(nodes), node)
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

	block := int64(5061200)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ethSync := GetEthSync(node)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		//allDeposits := node.cTQ.DoLoop([]int64{block})
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	// Check for the deposit
	block = int64(5061261)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ethSync := GetEthSync(node)
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits := ethSync.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d [deposit]\n\n", block, i+1, len(nodes), len(allDeposits))

		// TODO: inspect the messages for the right content
		assertMsgCountEqualDoLoop(t, "deposit", 1, len(allDeposits), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = int64(5061262)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		ethSync := GetEthSync(node)
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

	trustAddress := common.HexToAddress(test.GRAPHENE_TRUST2.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry(2, ":6000")

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)
	println("latest TXID=", txId)

	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST2, test.ETHER_NETWORKS[test.ROPSTEN])

	withdrawResult := make(chan control.WithdrawalResult)

	nodes[0].qTC.SuccessCb = func(c control.WithdrawalResult) {
		withdrawResult <- c
	}

	cursor := int64(5290510)
	fmt.Printf("=======================\n[CURSOR %d] BEGIN\n\n", cursor)
	for i, node := range nodes {
		refunds, err := node.qTC.DoLoop(cursor)
		assert.NoError(t, err, "error: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
		assert.Equal(t, 1, len(refunds), "refunds: cursor #%d [node #%d/%d id=%d]", cursor, i+1, len(nodes), node.nodeID)
	}

	fmt.Printf("[CURSOR %d] END\n=======================\n\n", cursor)

	time.Sleep(time.Second * 48)

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

	pubKey := "pooja"
	res, err := http.Get("http://localhost:5200/api/address/new/BTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

	block := int64(146)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		btcSync := GetBtcSync(node)
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
	res, err := http.Get("http://localhost:5200/api/address/new/BTC/" + pubKey)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	println("Address created ", string(bodyBytes))

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
