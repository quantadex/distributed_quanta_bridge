package peer_contact

import (
        "common/manifest"
)

/**
 * PaymentReq
 *
 * This is the struct that actually gets signed.
 */
type PaymentReq struct {
    BlockID  int // The coin block where payment was made
    CoinName string // The coin type represented (e.g ETH)
    QuantaAddr string // The address where to pay
    Amount  int // The amount to pay
}

/**
 * PeerMessage
 * 
 * This is the message that is being passed between nodes to gain signatures
 */
type PeerMessage struct {
    Proposal PaymentReq // The unsigned (raw) version of the msg
    SignedBy []int // List of nodeIDs (in-order) of nodes that signed
    NodesMissed int // Number of skipped nodes (not-signed)
    MSG []byte // The actual signed message
}

/**
 * PeerContact
 *
 * This module is used to communicate with other trust nodes in the manifest.
 * All communication is async. Messages are sent to peers in the group and only an OK is expected.
 * Messages are received by the node's listener agent and placed in the listener queue.
 * This module attaches to the listener queue-service and pulls received peer messages when requested. 
 */
type PeerContact interface {
    /**
     * AttachToListener
     *
     * Connect to the queue-service for the node's listener.
     * The Queue's name is in env variable (NODE_LISTENER_QUEUE)
     * Stash the Queue object in the local object
     * Return error if no variable or propogate error from Connect()
     */
    AttachToListener() error

    /**
     *  SendMsg
     *
     * Sends the given message to the node idnetified by nodeID in the manifest.
     * Return error if did not receive OK.
     */
    SendMsg(m *manifest.Manifest, destinationNodeID int, msg *PeerMessage) error

    /**
     * GetMsg
     *
     * Returns any message received by listener from peers.
     * If queue "peer_msg" is empty returns nil
     * Otherwise returns next msg
     */
    GetMsg() *PeerMessage
}

func NewPeerContact() (*PeerContact, error) {
    return nil, nil
}
