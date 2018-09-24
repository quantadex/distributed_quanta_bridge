package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
)

/**
 * CoinToQuanta
 *
 * This modules receives new deposits made to the coin trust
 * and using the round robin module creates transactions in quanta
 */
type CoinToQuanta struct {
    log logger.Logger
    coinChannel coin.Coin
    quantaChannel quanta.Quanta
    db kv_store.KVStore
    man *manifest.Manifest
    peer peer_contact.PeerContact
    coinName string
    rr *RoundRobinSigner
}

/**
 * NewCoinToQuanta
 *
 * Returns a new instance of the module
 * Initializes nothing so it should all be already initialized.
 */
func NewCoinToQuanta(   log logger.Logger,
                        db kv_store.KVStore,
                        c coin.Coin,
                        q quanta.Quanta,
                        man *manifest.Manifest,
                        kM key_manager.KeyManager,
                        coinName string,
                        nodeID int,
                        peer peer_contact.PeerContact ) *CoinToQuanta {
    res := &CoinToQuanta{}
    res.log = log
    res.coinChannel = c
    res.quantaChannel = q
    res.db = db
    res.man = man
    res.coinName = coinName
    res.peer = peer
    res.rr = NewRoundRobinSigner(log, man, nodeID, kM, db, peer)
    return res
}

/**
 * getNewCoinBlockIDs
 *
 * Returns a list of new blocks added to the coin block chain.
 */
func (c *CoinToQuanta) getNewCoinBlockIDs() []int {
    lastProcessed, valid := getLastBlock(c.db, c.coinName)
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
        c.log.Debug("No new block")
        return nil
    }
    blocks := make([]int, 0)
    for i := lastProcessed+1; i <= currentTop; i++ {
        blocks = append(blocks, i)
    }
    return blocks
}

/**
 * getDepositsInBlock
 *
 * Returns deposits made to the coin trust account in this block
 */
func (c *CoinToQuanta) getDepositsInBlock(blockID int) []*coin.Deposit {
    deposits, err := c.coinChannel.GetDepositsInBlock(blockID, c.man.ContractAddress)
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
            c.log.Error("Failed to submit deposut")
        }
    }
}

/**
 * DoLoop
 *
 * Do one iteration of the loop. Get all new coin blocks and theie deposits.
 * Shoot those into round robin
 * Get all ready messages from RR and send these to quanta.
 */
func (c *CoinToQuanta) DoLoop() {
    c.rr.addTick()
    blockIDs := c.getNewCoinBlockIDs()
    if blockIDs == nil {
        return
    }
    for _, blockID := range blockIDs {
        deposits := c.getDepositsInBlock(blockID)
        if deposits == nil {
            continue
        }
        c.rr.processNewDeposits(deposits)
    }
    allMsgs := c.rr.getExpiredMsgs()
    for true {
        msg := c.peer.GetMsg()
        if msg == nil {
            break
        }
        allMsgs = append(allMsgs, msg)
    }
    toSend := c.rr.processNewPeerMsgs(allMsgs)
    if len(toSend) > 0 {
        c.submitMessages(toSend)
    }
}
