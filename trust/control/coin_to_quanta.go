package control

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
)

/**
 * CoinToQuanta
 *
 * This modules receives new deposits made to the coin trust
 * and using the round robin module creates transactions in quanta
 */
type CoinToQuanta struct {
	log           logger.Logger
	coinChannel   coin.Coin
	quantaChannel quanta.Quanta
	db            kv_store.KVStore
	man           *manifest.Manifest
	peer          peer_contact.PeerContact
	coinName      string
	trustAddress  common.Address
	rr            *RoundRobinSigner
	C2QOptions
}

type C2QOptions struct {
	EthTrustAddress string
	BlockStartID    int64
}

const MAX_PROCESS_BLOCKS = 100

/**
 * NewCoinToQuanta
 *
 * Returns a new instance of the module
 * Initializes nothing so it should all be already initialized.
 */
func NewCoinToQuanta(log logger.Logger,
	db kv_store.KVStore,
	c coin.Coin,
	q quanta.Quanta,
	man *manifest.Manifest,
	kM key_manager.KeyManager,
	coinName string,
	nodeID int,
	peer peer_contact.PeerContact,
	options C2QOptions) *CoinToQuanta {
	res := &CoinToQuanta{C2QOptions: options}
	res.log = log
	res.coinChannel = c
	res.quantaChannel = q
	res.db = db
	res.man = man
	res.coinName = coinName
	res.peer = peer
	res.rr = NewRoundRobinSigner(log, man, nodeID, kM, db, q, peer)
	res.trustAddress = common.HexToAddress(options.EthTrustAddress)
	return res
}

/**
 * getNewCoinBlockIDs
 *
 * Returns a list of new blocks added to the coin block chain.
 */
func (c *CoinToQuanta) GetNewCoinBlockIDs() []int64 {
	lastProcessed, valid := GetLastBlock(c.db, c.coinName)
	if !valid {
		c.log.Error("Failed to get last processed ID")
		return nil
	}

	currentTop, err := c.coinChannel.GetTopBlockID()
	if err != nil {
		c.log.Error("Failed to get current top block")
		return nil
	}

	if lastProcessed > currentTop {
		c.log.Error("Coin top block smaller than last processed")
		return nil
	}

	if lastProcessed == currentTop {
		c.log.Debug(fmt.Sprintf("Coin2Quanta: No new block last=%d top=%d", lastProcessed, currentTop))
		return nil
	}
	blocks := make([]int64, 0)
	for i := common2.MaxInt64(c.BlockStartID, lastProcessed+1); i <= currentTop; i++ {
		blocks = append(blocks, i)
		if len(blocks) == MAX_PROCESS_BLOCKS {
			break
		}
	}
	c.log.Info(fmt.Sprintf("Got blocks %v", blocks))

	return blocks
}

/**
 * getDepositsInBlock
 *
 * Returns deposits made to the coin trust account in this block
 */
func (c *CoinToQuanta) getDepositsInBlock(blockID int64) []*coin.Deposit {
	watchAddresses, err := c.db.GetAllValues(ETHADDR_QUANTAADDR)
	if err != nil {
		c.log.Error("Failed to get watch addresses")
		return nil
	}

	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchAddresses)

	//var m map[string]string
	//m = make(map[string]string)
	//m[strings.ToLower("0xb1E02e31c9A2403FeAFA7E483Ebb3e1b5ffa3164")] = "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN"

	//deposits, err := c.coinChannel.GetDepositsInBlock(blockID, m)
	if err != nil {
		c.log.Error("Failed to get deposits from block")
		return nil
	}

	return deposits
}

/**
 * submitMessages
 *
 * Send withdrawal messagaes to quanta core
 */
func (c *CoinToQuanta) submitMessages(msgs []*peer_contact.PeerMessage) {
	for _, msg := range msgs {
		err := c.quantaChannel.ProcessDeposit(*msg)
		if err != nil {
			c.log.Error("Failed to submit deposit " + err.Error())
		}
		c.log.Infof("Submit deposit to %s %d",
			msg.Proposal.QuantaAdress, msg.Proposal.Amount)
	}
}

/**
 * DoLoop
 *
 * Do one iteration of the loop. Get all new coin blocks and theie deposits.
 * Shoot those into round robin
 * Get all ready messages from RR and send these to quanta.
 *
 * returns allDeposits []*coin.Deposit,
 *         allPeerMsgs []*peer_contact.PeerMessage,
 *         allSentMsgs []*peer_contact.PeerMessage
 */
func (c *CoinToQuanta) DoLoop(blockIDs []int64) ([]*coin.Deposit, []*peer_contact.PeerMessage, []*peer_contact.PeerMessage) {
	c.rr.deferQ.AddTick()
	c.log.Info(fmt.Sprintf("***** Start of Epoch %d # of blocks=%d man.N=%d,man.Q=%d *** ",
		c.rr.deferQ.Epoch(), len(blockIDs), c.man.N, c.man.Q))

	allDeposits := make([]*coin.Deposit, 0)

	if blockIDs != nil {
		for _, blockID := range blockIDs {
			deposits := c.getDepositsInBlock(blockID)
			if deposits != nil {
				c.log.Info(fmt.Sprintf("Block %d Got deposits %d %v", blockID, len(deposits), deposits))
				c.rr.processNewDeposits(deposits)

				for i := 0; i < len(deposits); i++ {
					allDeposits = append(allDeposits, deposits[i])
				}
			}

			addresses, err := c.coinChannel.GetForwardersInBlock(blockID)
			if err != nil {
				c.log.Error(err.Error())
				continue
			}

			for _, addr := range addresses {
				if addr.Trust.Hex() == c.trustAddress.Hex() {
					c.log.Infof("New Forwarder Address ETH->QUANTA address, %s -> %s", addr.ContractAddress.Hex(), addr.QuantaAddr)
					c.db.SetValue(ETHADDR_QUANTAADDR, strings.ToLower(addr.ContractAddress.Hex()), "", addr.QuantaAddr)
				} else {
					c.log.Error(fmt.Sprintf("MISMATCH: Forwarder address[%s] in blockID=%d does not match our trustAddress[%s]",
						addr.Trust.Hex(), blockID, c.trustAddress.Hex()))
				}
			}
		}
	}

	allMsgs := make([]*peer_contact.PeerMessage, 0)
	expiredMsgs, _ := c.rr.deferQ.Get(DQ_NAME)

	if expiredMsgs != nil {
		allMsgs = append(allMsgs, expiredMsgs.(*peer_contact.PeerMessage))
	}

	for true {
		msg := c.peer.GetMsg()
		if msg == nil {
			break
		}
		c.log.Infof("Got peer message %v", msg)
		allMsgs = append(allMsgs, msg)
	}

	toSend := c.rr.processNewPeerMsgs(allMsgs)
	if len(toSend) > 0 {
		c.log.Infof("Submitting %d messages", len(toSend))
		c.submitMessages(toSend)
	}

	if len(blockIDs) > 0 {
		lastBlockId := blockIDs[len(blockIDs)-1]
		c.log.Infof("set last block coin=%s height=%d", c.coinName, lastBlockId)
		setLastBlock(c.db, c.coinName, lastBlockId)
	}

	return allDeposits, allMsgs, toSend
}
