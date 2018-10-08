package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/spf13/viper"
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
    trustAddress common.Address
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
    res.rr = NewRoundRobinSigner(log, man, nodeID, kM, db, q, peer)
    res.trustAddress = common.HexToAddress(viper.GetString("TRUST_ETHEREUM_ADDR"))
    return res
}

/**
 * getNewCoinBlockIDs
 *
 * Returns a list of new blocks added to the coin block chain.
 */
func (c *CoinToQuanta) getNewCoinBlockIDs() []int64 {
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
        c.log.Debug(fmt.Sprintf("Coin2Quanta: No new block last=%d top=%d", lastProcessed, currentTop))
        return nil
    }
    blocks := make([]int64, 0)
    for i := lastProcessed+1; i <= currentTop; i++ {
        blocks = append(blocks, i)
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
    deposits, err := c.coinChannel.GetDepositsInBlock(blockID, c.man.ContractAddress)
    if err != nil {
        c.log.Error("Failed to get deposits from block")
        return nil
    }

    // let's make sure we can associate with quanta addresses
    linkedDeposits := []*coin.Deposit{}
    for _, d := range deposits {
        if q, e := c.db.GetValue(ETHADDR_QUANTAADDR, d.SenderAddr); e == nil && q != nil {
            d.QuantaAddr = *q
            linkedDeposits = append(linkedDeposits, d)
            c.log.Infof("Deposit is linked to %s", d.QuantaAddr)
        }
    }

    return linkedDeposits
}

/**
 * submitMessages
 *
 * Send withdrawal messagaes to quanta core
 */
func (c *CoinToQuanta) submitMessages(msgs []*peer_contact.PeerMessage) {
    for _, msg := range msgs {
        err := c.quantaChannel.ProcessDeposit(*msg)
        c.log.Infof("Process deposit to %s %d",
                        msg.Proposal.QuantaAdress, msg.Proposal.Amount)
        if err != nil {
            c.log.Error("Failed to submit deposit " + err.Error())
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
    c.log.Info(fmt.Sprintf("***** Start of Epoch %d *** ", c.rr.curEpoch))

    blockIDs := c.getNewCoinBlockIDs()
    if blockIDs != nil {
        for _, blockID := range blockIDs {
            deposits := c.getDepositsInBlock(blockID)
            c.log.Info(fmt.Sprintf("Got deposits %v", deposits))

            if deposits == nil {
                continue
            }
            c.rr.processNewDeposits(deposits)
        }

        for _, blockID := range blockIDs {
            addresses, err := c.coinChannel.GetForwardersInBlock(blockID)
            if err != nil {
                continue
            }

            for _, addr := range addresses {
                if addr.Trust.Hex() == c.trustAddress.Hex() {
                    c.db.SetValue(ETHADDR_QUANTAADDR, addr.ContractAddress.Hex(),"", addr.QuantaAddr)
                    c.log.Infof("Got new ETH->QUANTA address, %s -> %s", addr.ContractAddress.Hex(), addr.QuantaAddr)
                } else {
                    c.log.Error("Forward does not point to our trust address " + addr.Trust.Hex())
                }
            }
        }
    }

    allMsgs := c.rr.getExpiredMsgs()
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
        c.submitMessages(toSend)
    }

    if len(blockIDs) > 0 {
        lastBlockId := blockIDs[len(blockIDs)-1]
        c.log.Infof("set last block coin=%s height=%d", c.coinName, lastBlockId)
        setLastBlock(c.db, c.coinName, lastBlockId)
    }
}
