package quanta

import (
    "trust/peer_contact"
)

/**
 * Refund
 *
 * User's quanta return to trust in order to get a refund
 */
type Refund struct {
    CoinName string
    DestinationAddress string
    Amount int
    BlockID int
}

/**
 * Quanta
 *
 * This is the module through which the trust node communicates with a quanta-core node.
 * It is used to get the funds sent to the trust account as well as dispense funds from the
 * quanta trust.
 */
type Quanta interface {
    /**
     * Attach
     *
     * Connects to the quanta-core node. Returns error if this fails.
     */
    Attach() error

    /**
     * GetTopBlockID
     *
     * Returns the id of the latest quanta block.
     */
    GetTopBlockID() (int, error)

    /**
     * GetRefundsInBlock
     *
     * Returns a list of refunds that were made to the specified address in the given block.
     * Return nil if no matching deposits.
     */
    GetRefundsInBlock(blockID int, trustAddress string) ([]Refund, error)

    /**
     * ProcessDeposit
     *
     * Once enough nodes have signed the deposit the last node sends it to quanta to
     * transfer the funds into the user's quanta account
     */
    ProcessDeposit(deposit *peer_contact.PeerMessage) error
}

func NewQuanta() (*Quanta, error) {
    return nil, nil
}
