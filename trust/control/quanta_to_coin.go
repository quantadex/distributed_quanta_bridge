package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "encoding/json"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
)

const QUANTA = "QUANTA"

/**
 * QuantaToCoin
 *
 * This module polls for new quanta blocks and find refunds submitted to the quanta trust.
 * It validates and signs the the withdrawls and issues them to the coin's smart contract.
 */
type QuantaToCoin struct {
    logger logger.Logger
    db kv_store.KVStore
    coinChannel coin.Coin
    quantaChannel quanta.Quanta
    quantaTrustAddress string
    coinContractAddress string
    kM key_manager.KeyManager
    coinName string
    nodeID int
}

/**
 * NewQuantaToCoin
 *
 * Create a new instance of the class. This does not initialize modules.
 * All modules must already have been initialized and passed in here.
 */
func NewQuantaToCoin(   log logger.Logger,
                        db kv_store.KVStore,
                        c coin.Coin,
                        q quanta.Quanta,
                        quantaTrustAddress string,
                        coinContractAddress string,
                        kM key_manager.KeyManager,
                        coinName string,
                        nodeID int ) *QuantaToCoin {
    res := &QuantaToCoin{}
    res.logger = log
    res.db = db
    res.coinChannel = c
    res.quantaChannel = q
    res.quantaTrustAddress = quantaTrustAddress
    res.coinContractAddress = coinContractAddress
    res.kM = kM
    res.coinName = coinName
    res.nodeID = nodeID
    return res
}

/**
 * getNewBlockIDs
 *
 * Gets a list of new quanta blocks since last processed.
 */
func (c *QuantaToCoin) getNewBlockIDs() []int64 {
    lastProcessed, valid := getLastBlock(c.db, QUANTA)
    if !valid {
        c.logger.Error("Failed to get last processed block ID")
        return nil
    }

    currentTop, err := c.quantaChannel.GetTopBlockID()
    if err != nil {
        c.logger.Error("Failed to get top quanta block ID")
        return nil
    }

    if lastProcessed > currentTop {
        c.logger.Error("Quanta top block smaller than last processed")
        return nil
    }

    if lastProcessed == currentTop {
        c.logger.Debug("Quanta2Coin: No new block")
        return nil
    }
    blocks := make([]int64, 0)
    for i := lastProcessed+1; i <= currentTop; i++ {
        blocks = append(blocks, i)
    }
    return blocks
}

/**
 * getRedundsInBlock
 *
 * Gets all the quanta refunds in a given block
 */
func (c *QuantaToCoin) getRefundsInBlock(blockID int64) []quanta.Refund {
    refunds, err := c.quantaChannel.GetRefundsInBlock(blockID, c.quantaTrustAddress)
    if err != nil {
        c.logger.Error("Failed to get refunds in quanta block")
        return nil
    }
    return refunds
}

/**
 * validateAndSignRefund
 *
 * Checks that the refund has not been previously issued and marks it signed in DB.
 */
func (c *QuantaToCoin) validateAndSignRefund(refund *quanta.Refund) bool {
    refKey := getKeyName(refund.CoinName, refund.DestinationAddress, refund.BlockID)
    success := confirmTx(c.db, QUANTA_CONFIRMED, refKey)
    if !success {
        c.logger.Error("Failed to confirm transaction")
        return false
    }
    success = signTx(c.db, QUANTA_CONFIRMED, refKey)
    if !success {
        return false
    }
    return true
}

/**
 * submitRefund
 *
 * Signs the withdrawal and sends it to the coin's smart contract.
 */
func (c *QuantaToCoin) submitRefund(refund *quanta.Refund) bool {
    w := &coin.Withdrawal{}
    w.NodeID = c.nodeID
    w.CoinName = refund.CoinName
    w.DestinationAddress = refund.DestinationAddress
    w.QuantaBlockID = refund.BlockID
    w.Amount = refund.Amount

    raw, err := json.Marshal(w)
    if err != nil {
        c.logger.Error("Failed to marshal withdrawal")
        return false
    }

    s, err := c.kM.SignMessage(raw)
    if err != nil {
        c.logger.Error("Failed to sign withdrawal")
        return false
    }

    err = c.coinChannel.SendWithdrawal(c.coinContractAddress, *w, s)
    if err != nil {
        c.logger.Error("Failed to send withdrawal")
        return false
    }
    return true
}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 */
func (c *QuantaToCoin) DoLoop() {
    newBlocks := c.getNewBlockIDs()
    if newBlocks == nil {
        return
    }
    for _, blockID := range newBlocks {
        success := setLastBlock(c.db, QUANTA, blockID)
        if !success {
            c.logger.Error("Failed to mark block as signed")
            continue // should skip
        }
        refunds := c.getRefundsInBlock(blockID)
        if refunds == nil {
            continue
        }
        for _, refund := range refunds {
            valid := c.validateAndSignRefund(&refund)
            if !valid {
                continue //skip invalid. Probably log.
            }
            c.submitRefund(&refund)
        }
    }

}
