package quanta

import (
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
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
     * AttachQueue
     *
     * Connects to the quanta-core node. Returns error if this fails.
     */
    AttachQueue(queue queue.Queue) error

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
    ProcessDeposit(deposit peer_contact.PeerMessage) error

    CreateProposeTransaction(*coin.Deposit) (string, error) // base64 tx envelope
}

func NewQuanta() (Quanta, error) {
    return &QuantaClient{}, nil
}

/**
 * Submitworker's job is to submit the deposits into
 * the horizon service, and retry as neccessary
 * Decouples whether horizon is online or not
 */
type SubmitWorker interface {
    Dispatch()
    AttachQueue(queue queue.Queue) error
}

func NewSubmitWorker() (SubmitWorker) {
    return &SubmitWorkerImpl{}
}
