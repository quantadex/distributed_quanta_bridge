package control

import (
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "encoding/json"
)

const DELAY_PENALTY = 10

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
    deferQ map[int][]*peer_contact.PeerMessage
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
                            peer peer_contact.PeerContact ) *RoundRobinSigner {

    res := &RoundRobinSigner{}
    res.log = log
    res.man = man
    res.myNodeID = myNodeID
    res.kM = kM
    res.db = db
    res.peer = peer
    res.deferQ = make(map[int][]*peer_contact.PeerMessage, 0)
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
    var deferList []*peer_contact.PeerMessage
    var found bool

    deferList, found = r.deferQ[expires]
    if !found {
        deferList = make([]*peer_contact.PeerMessage, 0)
    }
    deferList = append(deferList, msg)
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
            for _, msg := range v {
                results = append(results, msg)
            }
        }
    }
    if len(hits) == 0 {
        return nil
    }
    for _, k := range hits {
        delete(r.deferQ, k)
    }
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
    target := make([]byte, 0)
    copy(target, msg.MSG)
    var err error

    for _, nodeID := range msg.SignedBy {
        pubKey := r.man.Nodes[nodeID].PubKey
        target, err = r.kM.DecodeMessage(target, pubKey)
        if err != nil {
            return false
        }
    }
    decoded := &peer_contact.PaymentReq{}
    err = json.Unmarshal(target, decoded)
    if err != nil {
        return false
    }
    if decoded.BlockID != msg.Proposal.BlockID {
        return false
    }
    if decoded.CoinName != msg.Proposal.CoinName {
        return false
    }
    if decoded.QuantaAdress != msg.Proposal.QuantaAdress {
        return false
    }
    if decoded.Amount != msg.Proposal.Amount {
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
    msg.MSG = make([]byte, 0)
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
    if len(msg.SignedBy) == 0 {
        data, err = json.Marshal(msg.Proposal)
        if err != nil {
            r.log.Error("Failed to marshal payment req")
            return false
        }
    }
    data, err = r.kM.SignMessage(data)
    if err != nil {
        r.log.Error("Failed to encrypt the message")
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
        tolerance := r.man.N - r.man.Q
        for msg.NodesMissed < tolerance {
            err := r.peer.SendMsg(r.man, destination, msg)
            if err == nil {
                return true
            }
            destination = (destination + 1) % r.man.Na
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
        startNode := deposit.BlockID % r.man.N
        missedNodes := 0
        for i := 0; i < r.man.N; i++ {
            nodeID := (r.myNodeID + i) % r.man.N
            if nodeID == startNode {
                break
            }
            missedNodes++
        }
        tolerance := r.man.N - r.man.Q
        if missedNodes > tolerance {
            continue
        }
        msg := r.createNewPeerMsg(deposit, missedNodes)
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
            continue
        }
        success = r.validateIntegrity(msg)
        if !success {
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
