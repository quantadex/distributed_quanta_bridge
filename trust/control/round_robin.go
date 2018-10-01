package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    dll "github.com/emirpasic/gods/lists/doublylinkedlist"
    "fmt"
    "github.com/quantadex/distributed_quanta_bridge/common"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
)

const DELAY_PENALTY = 5

/**
 * RoundRobinSigner
 *
 * Implements the peer node distributed signing algorithm
 */
type RoundRobinSigner struct {
    log logger.Logger
    man *manifest.Manifest
    myNodeID int
    kM key_manager.KeyManager
    db kv_store.KVStore
    peer peer_contact.PeerContact
    deferQ map[int]*dll.List
    quanta quanta.Quanta
    curEpoch int
}

/**
 * NewRoundRobinSigner
 *
 * Creates a new round-robin signer.
 * Does not initialize any module.
 * All modules must already by initialized and passed in.
 *
 */
func NewRoundRobinSigner(   log logger.Logger,
                            man *manifest.Manifest,
                            myNodeID int,
                            kM key_manager.KeyManager,
                            db kv_store.KVStore,
                            quanta quanta.Quanta,
                            peer peer_contact.PeerContact ) *RoundRobinSigner {

    res := &RoundRobinSigner{}
    res.log = log
    res.man = man
    res.myNodeID = myNodeID
    res.kM = kM
    res.db = db
    res.peer = peer
    res.quanta = quanta
    res.deferQ = make(map[int]*dll.List, 0)
    res.curEpoch = 0
    return res
}

/**
 * addToDeferQ
 *
 * Adds a message to a defer queue with a given time penalty
 * When the time penalty has expired this message will be made available
 */
func (r *RoundRobinSigner) addToDeferQ(msg *peer_contact.PeerMessage) {
    expires := r.curEpoch + (msg.NodesMissed*DELAY_PENALTY)
    r.log.Info(fmt.Sprintf("Added msg to DeferQ at expiration=%d epoch=%d", expires, r.curEpoch))

    var deferList *dll.List
    var found bool

    deferList, found = r.deferQ[expires]
    if !found {
        deferList = dll.New()
        r.deferQ[expires] = deferList
    }
    deferList.Append(msg)
}

/**
 * addTick
 *
 * Increment the internal clock.
 * This is used for knowing when to dequeue from deferred queue.
 *
 */
func (r *RoundRobinSigner) addTick() {
    r.curEpoch +=1
}

/**
 * getExpiredMsgs
 *
 * Returns all deferred messages that have not been previously consumed by which have
 * had their time penalty expire
 */
func (r *RoundRobinSigner) getExpiredMsgs() []*peer_contact.PeerMessage {
    results := make([]*peer_contact.PeerMessage, 0)
    hits := make([]int, 0)
    for k, v := range r.deferQ {
        if k <= r.curEpoch {
            hits = append(hits, k)
            it := v.Iterator()
            for it.Next() {
                results = append(results, it.Value().(*peer_contact.PeerMessage))
            }
        }
    }
    if len(hits) == 0 {
        return nil
    }
    for _, k := range hits {
        delete(r.deferQ, k)
    }
    r.log.Info(fmt.Sprintf("Expired? %v %v\n", hits, results))

    return results
}

/**
 * validateTransaction
 *
 * Returns true if the given message has indeed been seen by this node and is in the
 * CONFIRMED state meaning it has not been processed previosly.
 *
 */
func (r *RoundRobinSigner) validateTransaction(msg *peer_contact.PeerMessage) bool {
    txKey := getKeyName(msg.Proposal.CoinName, msg.Proposal.QuantaAdress, msg.Proposal.BlockID)
    state := getState(r.db, COIN_CONFIRMED, txKey)
    r.log.Infof("Validate key=%s state=%s", txKey, state)
    if state == CONFIRMED {
        return true
    }
    return false
}

/**
 * validateIntegrity
 *
 * Returns true if the signed content can be decreptyed by the series of node keys that
 * are claimed in the signing history and if the decrypted contents match the raw contents
 */
func (r *RoundRobinSigner) validateIntegrity(msg *peer_contact.PeerMessage) bool {
    if len(msg.SignedBy) == 0 {
        return true
    }
    valid, err := r.kM.VerifyTransaction(msg.MSG)
    if !valid {
        r.log.Error("validateIntegrity: failed to verify message " + err.Error())
        return false
    }

    decoded, err := r.kM.DecodeTransaction(msg.MSG)
    if err != nil {
        r.log.Error("validateIntegrity: failed to decode message")
        return false
    }
    //if decoded.BlockID != msg.Proposal.BlockID {
    //    return false
    //}
    //if decoded.CoinName != msg.Proposal.CoinName {
    //    return false
    //}
    if decoded.QuantaAddr != msg.Proposal.QuantaAdress {
        r.log.Error("validateIntegrity: address mismatch")
        return false
    }
    if decoded.Amount != msg.Proposal.Amount {
        r.log.Error("validateIntegrity: amount mismatch")
        return false
    }
    return true
}

/**
 * createNewPeerMsg
 *
 * This is the first node to process this deposit. Start the peer message.
 *
 */
func (r *RoundRobinSigner) createNewPeerMsg(deposit *coin.Deposit, missedNodes int) *peer_contact.PeerMessage {
    payment := &peer_contact.PaymentReq{}
    payment.BlockID = deposit.BlockID
    payment.CoinName = deposit.CoinName
    payment.QuantaAdress = deposit.QuantaAddr
    payment.Amount = deposit.Amount

    msg := &peer_contact.PeerMessage{}
    msg.Proposal = *payment
    msg.SignedBy = make([]int, 0)
    msg.NodesMissed = missedNodes
    var err error
    msg.MSG, err = r.quanta.CreateProposeTransaction(deposit)
    if err != nil {
        r.log.Error("error creating tx " + err.Error())
    }
    return msg
}

/**
 * signPeerMsg
 *
 * Mark the message as signed in DB ensuring a node only ever signs 1 msg
 * Encrypt the chained contents with private key.
 *
 */
func (r *RoundRobinSigner) signPeerMsg(msg *peer_contact.PeerMessage) bool {
    txKey := getKeyName(msg.Proposal.CoinName, msg.Proposal.QuantaAdress, msg.Proposal.BlockID)
    success := signTx(r.db, COIN_CONFIRMED, txKey)
    if !success {
        r.log.Error("Failed to mark as signed")
        return false
    }
    data := msg.MSG
    var err error

    // need to check if we signed it

    //if len(msg.SignedBy) == 0 {
    //    data, err = json.Marshal(msg.Proposal)
    //    if err != nil {
    //        r.log.Error("Failed to marshal payment req")
    //        return false
    //    }
    //}
    r.log.Infof("sign peer msg %v", msg.SignedBy)

    data, err = r.kM.SignTransaction(data)
    if err != nil {
        r.log.Error("Failed to sign the message: " + err.Error())
        return false
    }
    msg.MSG = data
    msg.SignedBy = append(msg.SignedBy, r.myNodeID)
    return true
}

/**
 * sendMessage
 *
 * Sends the message to the next peer inline.
 * On failure tries subsequent peer so long as the number of missed
 * nodes is less than the quorum tolerance
 */
func (r *RoundRobinSigner) sendMessage(msg *peer_contact.PeerMessage) bool {
        destination := (r.myNodeID + 1) % r.man.N
        tolerance := common.MaxInt(1, r.man.N - r.man.Q)
        r.log.Infof("sendMessage to peer %d missed=%d tolerance=%d", destination, msg.NodesMissed, tolerance)

        for msg.NodesMissed < tolerance {
            err := r.peer.SendMsg(r.man, destination, msg)
            if err == nil {
                return true
            }
            destination = (destination + 1) % r.man.N
            msg.NodesMissed++
        }
        return false
}

/**
 * processNewDeposits
 *
 * Called from higher up with a list of new deposits sent to the coin trust.
 * For any deposit where this can be the first node. Create a message and insert
 * into defered queue.
 */
func (r *RoundRobinSigner) processNewDeposits(deposits []*coin.Deposit) {

    for _, deposit := range deposits {

        // let's confirm them
        r.log.Infof("Confirm transaction %d", deposit.Amount)
        confirmTx(r.db, COIN_CONFIRMED, getKeyName(deposit.CoinName, deposit.QuantaAddr, deposit.BlockID))

        startNode := deposit.BlockID % r.man.N
        missedNodes := 0
        for i := 0; i < r.man.N; i++ {
           nodeID := (r.myNodeID + i) % r.man.N
           if nodeID == startNode {
               break
           }
           missedNodes++
        }
        // if too many nodes missed, just skip for this deposit
        tolerance := r.man.N - r.man.Q
        if missedNodes > tolerance {
           continue
        }
        msg := r.createNewPeerMsg(deposit, 0)
        msg.Proposer = r.myNodeID
        r.addToDeferQ(msg)
    }
}

/**
 * processNewPeerMsgs
 *
 * Called from above with new peer messages that come either from the deferred queue or
 * have arrived from peers. Signs any valid message and sends to next peer.
 *
 * Returns those messages where that have reached signature quorum and can go to quanta
 */
func (r *RoundRobinSigner) processNewPeerMsgs(msgs []*peer_contact.PeerMessage) []*peer_contact.PeerMessage {

    toSend := make([]*peer_contact.PeerMessage, 0)
    for _, msg := range msgs {
        success := r.validateTransaction(msg)
        if !success {
            r.log.Infof("Msg failed validate transaction proposer=%d", msg.Proposer)
            continue
        }
        success = r.validateIntegrity(msg)
        if !success {
            r.log.Infof("Msg failed validate integrity")
            continue
        }
        success = r.signPeerMsg(msg)
        if !success {
            continue
        }
        if len(msg.SignedBy) > r.man.Q {
            r.log.Error("Too many signatures")
            continue
        }

        r.log.Infof("signed by n=%d q=%d", len(msg.SignedBy), r.man.Q)
        if len(msg.SignedBy) == r.man.Q {
            toSend = append(toSend, msg)
            continue
        }
        success = r.sendMessage(msg)
        if !success {
            r.log.Error("Failed to send message to peers")
        }
    }
    return toSend
}
