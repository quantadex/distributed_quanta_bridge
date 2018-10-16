package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "encoding/json"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/quanta_book/consensus/cosi"
    "time"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
    "github.com/pkg/errors"
    "encoding/base64"
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
    rr *RoundRobinSigner
    cosi *cosi.Cosi
    trustPeer *peer_contact.TrustPeerNode
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
                        man *manifest.Manifest,
                        quantaTrustAddress string,
                        coinContractAddress string,
                        kM key_manager.KeyManager,
                        coinName string,
                        peer peer_contact.PeerContact,
                        queue_ queue.Queue,
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
    res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_)
    res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

    res.cosi.Verify = func(msg string) error {
        withdrawal, err := res.coinChannel.DecodeRefund(msg)
        if err != nil {
            log.Error("Unable to decode refund at persist")
            return err
        }

        refKey := getKeyName(withdrawal.CoinName, withdrawal.DestinationAddress, withdrawal.QuantaBlockID)
        state := getState(db, QUANTA_CONFIRMED, refKey)
        if state == CONFIRMED {
            return nil
        }
        return errors.New("Unable to verify: " + refKey)
    }

    res.cosi.Persist = func(msg string) error {
        withdrawal, err := res.coinChannel.DecodeRefund(msg)
        if err != nil {
            log.Error("Unable to decode refund at persist")
            return err
        }

        refKey := getKeyName(withdrawal.CoinName, withdrawal.DestinationAddress, withdrawal.QuantaBlockID)

        success := signTx(db, QUANTA_CONFIRMED, refKey)
        if !success {
            res.logger.Error("Failed to confirm transaction")
            return errors.New("Failed to confirm transaction")
        }
        return nil
    }

    res.cosi.SignMsg = func(msg string) (string, error) {
        decoded, err := base64.StdEncoding.DecodeString(msg)
        if err != nil {
            log.Error("Unable to Sign refund msg")
            return "", err
        }
        encodedSig, err := res.kM.SignMessage(decoded)
        if err != nil {
            log.Error("Unable to Sign/encode refund msg")
            return "", err
        }

        return base64.StdEncoding.EncodeToString(encodedSig), nil
    }

    res.cosi.Start()
    return res
}

/**
 * getNewBlockIDs
 *
 * Gets a list of new quanta blocks since last processed.
 */
func (c *QuantaToCoin) GetNewBlockIDs() []int64 {
    lastProcessed, valid := getLastBlock(c.db, QUANTA)
    if !valid {
        c.logger.Error("Failed to get last processed block ID")
        return nil
    }

    currentTop, err := c.quantaChannel.GetTopBlockID(c.quantaTrustAddress)
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
func (c *QuantaToCoin) getRefundsInBlock(blockID int64) ([]quanta.Refund, int64) {
    refunds, nextBlockId, err := c.quantaChannel.GetRefundsInBlock(blockID, c.quantaTrustAddress)
    if err != nil {
        c.logger.Error("Failed to get refunds in quanta block")
        return nil, 0
    }
    return refunds, nextBlockId
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
func (c *QuantaToCoin) DoLoop(blockIDs []int64) {
    if blockIDs == nil {
        return
    }

    for _, blockID := range blockIDs {
        success := setLastBlock(c.db, QUANTA, blockID)
        if !success {
            c.logger.Error("Failed to mark block as signed")
            continue // should skip
        }
        refunds, nextBlockID := c.getRefundsInBlock(blockID)
        if refunds == nil {
            continue
        }

        // feed messages back
        msg := c.trustPeer.GetMsg()
        c.cosi.CosiMsgChan <- msg

        for _, refund := range refunds {
            refKey := getKeyName(refund.CoinName, refund.DestinationAddress, refund.BlockID)
            confirmTx(c.db, QUANTA_CONFIRMED, refKey)

            if c.nodeID == 0 {
                encoded, err := c.coinChannel.EncodeRefund(coin.Withdrawal{
                    CoinName: refund.CoinName,
                    DestinationAddress: refund.DestinationAddress,
                    QuantaBlockID: refund.BlockID,
                })

                if err != nil {
                    c.logger.Error("Failed to encode refund " + err.Error())
                    continue
                }

                c.cosi.StartNewRound(encoded)

                result := <- c.cosi.FinalSigChan

                if result.Msg == nil {
                    c.logger.Error("Unable to sign for refund")
                } else {
                    c.logger.Infof("Great! Cosi successfully signed refund")
                    c.submitRefund(nil)
                }
            }
        }
    }

}
