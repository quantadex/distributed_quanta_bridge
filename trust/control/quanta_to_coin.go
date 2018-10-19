package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/quanta_book/consensus/cosi"
    "time"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
    "github.com/pkg/errors"
    "github.com/ethereum/go-ethereum/common"
    "strings"
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
    coinkM key_manager.KeyManager
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
    res.coinkM = kM
    res.coinName = coinName
    res.nodeID = nodeID
    res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_)
    res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

    res.cosi.Verify = func(msg string) error {
        withdrawal, err := res.coinChannel.DecodeRefund(msg)
        if err != nil {
            log.Error("Unable to decode refund")
            return err
        }

        refKey := getKeyName(withdrawal.CoinName, withdrawal.DestinationAddress, 0)
        state := getState(db, QUANTA_CONFIRMED, refKey)
        if state == CONFIRMED {
            return nil
        }
        log.Error("Unable to verify refund " + refKey + " " + state)
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
        encodedSig, err := res.coinkM.SignTransaction(msg)
        log.Infof("Sign msg %s", encodedSig)

        if err != nil {
            log.Error("Unable to Sign/encode refund msg")
            return "", err
        }

        return encodedSig, nil
    }

    res.cosi.Start()

    go func() {
       for {
           // feed messages back
           msg := res.trustPeer.GetRefundMsg()

           if msg != nil {
               //log.Infof("Got peer message %v", msg.Signed)
               res.cosi.CosiMsgChan <- msg
               continue
           }

           time.Sleep(100 * time.Millisecond)
       }
    }()

    return res
}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 */
func (c *QuantaToCoin) DoLoop(cursor int64) {
    refunds, _ , err := c.quantaChannel.GetRefundsInBlock(cursor, c.quantaTrustAddress)

    if err != nil {
        c.logger.Error(err.Error())
        return
    }

    c.logger.Infof("Got refunds %v", refunds)

    // separate confirm, and sign as two different stages
    for _, refund := range refunds {
        refKey := getKeyName(refund.CoinName, strings.ToLower(common.HexToAddress(refund.DestinationAddress).Hex()), 0)
        c.logger.Infof("Confirm Refund = %s tx=%s", refKey, refund.TransactionId)
        confirmTx(c.db, QUANTA_CONFIRMED, refKey)

        //TODO: do checksum check, should bounce back the payment
        if refund.DestinationAddress == "" {
            c.logger.Error("Refund is missing destination address, skipping.")
            continue
        }

        // i'm the leader
        if c.nodeID == 0 {
            txId, err := c.coinChannel.GetTxID(common.HexToAddress(c.coinContractAddress))
            if err != nil {
                c.logger.Error("Could not get txID: " + err.Error() + " " + c.coinContractAddress)
                //TODO: How to handle this?
            }
            w := coin.Withdrawal{
                TxId: txId,
                CoinName: refund.CoinName,
                DestinationAddress: refund.DestinationAddress,
                QuantaBlockID: refund.OperationID,
                Amount: refund.Amount,
            }
            c.logger.Infof("Start new round %v", w)
            encoded, err := c.coinChannel.EncodeRefund(w)

            if err != nil {
                c.logger.Error("Failed to encode refund " + err.Error())
                continue
            }

            // wait for other node to see the tx
            time.Sleep(time.Second * 1)
            c.cosi.StartNewRound(encoded)

            result := <- c.cosi.FinalSigChan

            if result.Msg == nil {
                c.logger.Error("Unable to sign for refund")
            } else {
                // save this to queue for later in case ETH RPC is down.
                w.Signatures = result.Msg
                c.logger.Infof("Great! Cosi successfully signed refund")
                tx, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.coinContractAddress), c.coinkM.GetPrivateKey() ,&w)
                if err != nil {
                    c.logger.Error(err.Error())
                }
                c.logger.Infof("Submitted withdrawal in tx=%s", tx)
            }

            success := setLastBlock(c.db, QUANTA, refund.PageTokenID)
            if !success {
                c.logger.Error("Failed to mark block as signed")
                continue // should skip
            }

        }
    }
    c.logger.Infof("Next cursor is = %d", cursor)
}
