package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/stretchr/testify/assert"
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
func TestRopstenNativeETH(t *testing.T) {
	r := StartRegistry()
	nodes := StartNodes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	// DEPOSIT to TEST2
	// 0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A transfers 0.01 ETH to forward address: 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 on block 4327101
	// Foward contract 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98 created on 4327057
	block := int64(4327057)
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
	block = 4327101
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

	block = 4327102
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		// TODO: inspect the messages for the right content
		// index 0 does not always send the message.
		// last node relative from node it was sent from - see round_robin
		if int(4327101-1)%3 == i {
			// the first node doesn't receive any peer messages (yet), we will check for it in block #4249020
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "sent", 1, len(allSentMsgs), block, i+1, len(nodes), node)
		} else {
			// the rest of the nodes each receive a peer message
			assertMsgCountEqualDoLoop(t, "deposit", 0, len(allDeposits), block, i+1, len(nodes), node)
			assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
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

// Block 4354954: ERC-20  0x719791a052ed86360015e659542d3b0b7f44182e created by 0x9b93a4be348ab36c16ec6861602f84321f24e544  tx=https://ropsten.etherscan.io/tx/0xfc50ed99430174aa791ea9fa0883ef12d1ba7aa3a2daf26f3bb255b8f5f9af1b
// Block 4354971: Forward contract created  0x527370DEF157BD1113DB9448BD05E6402fFb5A0d forwarding to trust contract from 0x9b93a4be348ab36c16ec6861602f84321f24e544
// Block 4355067: Nothing - get block to agree

// Block 4356004 ERC-20 0x541d973a7168dbbf413eab6993a5e504ec5accb0
// Block 4356013 sent .0001234  precision 9  tx=https://ropsten.etherscan.io/tx/0x51a1018c6b2afd7bffb52d05178c3b66c9336e8c2dfeca0110f405cc41613492
func TestRopstenERC20Token(t *testing.T) {
	r := StartRegistry()
	ercContract := "0x541d973a7168dbbf413eab6993a5e504ec5accb0"
	nodes := StartNodes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	initialBalance, err := nodes[0].q.GetBalance(ercContract, test.QUANTA_ACCOUNT)
	assert.NoError(t, err)

	fmt.Printf("[ASSET %s] [ACCOUNT %s] initial_balance = %.9f\n", ercContract, test.QUANTA_ACCOUNT, initialBalance)

	time.Sleep(time.Millisecond * 250)
	DoLoopDeposit(nodes, []int64{4354971}) // forward address
	DoLoopDeposit(nodes, []int64{4356013}) // we make deposit
	DoLoopDeposit(nodes, []int64{4356014}) // no-op

	time.Sleep(time.Second * 6)
	newBalance, err := nodes[0].q.GetBalance(ercContract, test.QUANTA_ACCOUNT)
	assert.NoError(t, err)
	assert.Equal(t, initialBalance+float64(0.001234), newBalance)
	fmt.Printf("Initial balance=%f , expecting final balance = %f\n", initialBalance, initialBalance+float64(0.001234))
	time.Sleep(time.Second * 5)
	StopNodes(nodes)
	StopRegistry(r)
}

func TestTrustNode_Stop(t *testing.T) {
	r := StartRegistry()
	//indexes := []int{0,1,2}
	//nodes := []*TrustNode{}
	//nodes = StartNodesWithIndexes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN],true,indexes, nodes)
	nodes := StartNodes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])
	time.Sleep(time.Millisecond * 250)

	block := int64(4327057)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327101
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327102
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		if int(4327101-1)%3 == i {
			assertMsgCountEqualDoLoop(t, "sent", 1, len(allSentMsgs), block, i+1, len(nodes), node)
		} else {
			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		}

	}

	//StopNodeListener(nodes[0])
	StopSingleNode(nodes[0])
	time.Sleep(time.Second * 6)

	block = int64(4327057)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
		fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327127
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
		fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
		assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327128
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := nodes[i].cTQ.DoLoop([]int64{block})
		fmt.Printf("[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i, len(nodes), node)
		if (4327128-1)%3 == i {
			assertMsgCountEqualDoLoop(t, "peer", 1, len(allPeerMsgs), block, i+1, len(nodes), node)
		} else {
			assertMsgCountEqualDoLoop(t, "peer", 0, len(allPeerMsgs), block, i+1, len(nodes), node)
		}

	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	time.Sleep(time.Second * 6)

	indexToStart := []int{0}
	//nodes = StartNodeListener(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], nodes)
	nodes = StartNodesWithIndexes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], false, indexToStart, nodes)
	time.Sleep(time.Second * 6)
	block = int64(4327057)
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327101
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))
		assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
	}
	fmt.Printf("[BLOCK %d] END\n=======================\n\n", block)

	block = 4327102
	fmt.Printf("=======================\n[BLOCK %d] BEGIN\n\n", block)
	for i, node := range nodes {
		fmt.Printf("[BLOCK %d] Node[#%d/%d id=%d] calling doLoop...\n", block, i+1, len(nodes), node.nodeID)
		allDeposits, allPeerMsgs, allSentMsgs := node.cTQ.DoLoop([]int64{block})
		fmt.Printf("...[BLOCK %d] Node[#%d/%d] counts %d/%d/%d [deposit/peer/sent]\n\n", block, i+1, len(nodes), len(allDeposits), len(allPeerMsgs), len(allSentMsgs))

		if int(4327101-1)%3 == i {
			assertMsgCountEqualDoLoop(t, "sent", 1, len(allSentMsgs), block, i+1, len(nodes), node)
		} else {
			assertMsgCountEqualDoLoop(t, "sent", 0, len(allSentMsgs), block, i+1, len(nodes), node)
		}
	}

	time.Sleep(time.Second * 6)
	StopNodes(nodes)
	StopRegistry(r)
}

func TestWithdrawal(t *testing.T) {
	ethereumClient, err := ethclient.Dial(test.ETHER_NETWORKS[test.ROPSTEN].Rpc)
	assert.Nil(t, err)

	trustAddress := common.HexToAddress(test.ROPSTEN_TRUST.TrustContract)
	contract, err := contracts.NewTrustContract(trustAddress, ethereumClient)
	assert.NoError(t, err)

	r := StartRegistry()

	txId, err := contract.TxIdLast(nil)
	assert.NoError(t, err)

	println("latest TXID=", txId)

	nodes := StartNodes(test.QUANTA_ISSUER, test.ROPSTEN_TRUST, test.ETHER_NETWORKS[test.ROPSTEN])

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

		if i == 0 {
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
