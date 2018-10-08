package coin

/**
 * Deposit
 *
 * The deposit struct captures the data of making a deposit into the trust
 */
type Deposit struct {
    CoinName string // Type of coin (e.g. ETH)
    SenderAddr string
    QuantaAddr string // Destination quanta acount
    Amount int64 // Deposit size
    BlockID int // The blockID in which this deposit was found
}

/**
 * Withdrawal
 *
 * The data structure that needs to be filled out to do a succesful withdrawal.
 */
type Withdrawal struct {
    NodeID int // The Node authorizing this
    CoinName string // The type of coin (e.g. ETH)
    DestinationAddress string // Where this money is going
    QuantaBlockID int // Which block this transaction was processed in quanta
    Amount int // The withdrawal size
}

/**
 * Coin
 * 
 * This module attaches to a coin-core (e.g. eth-core) node running on the same machine.
 * It is used to get the block updates from the coin as well as submit payment to the contract.
 *
 */
type Coin interface {
    /**
     * Attach
     *
     * Connect to the specified coin core node. Return error if failed.
     */
    Attach() error

    /**
     * GetTopBlockID
     *
     * Returns the ID of the newest block in the chain.
     */
    GetTopBlockID() (int64, error)

    /**
     * GetDepositsInBlock
     *
     * For the specified blockID return all transactions that sent money to the specified address.
     * The quanta address is specified as data in the transaction.
     * Returns as a list of deposits.
     * Returns nil if no matching deposits.
     * Returns error of error encountered
     */
    GetDepositsInBlock(blockID int64, trustAddress string) ([]*Deposit, error)

    /**
     * GetForwardersInBlock
     *
     * Forwarders are smart contracts that are pointing into our trust address
     * with information about QUANTA Address
     * We will record this in our KV later, to know where deposits came from.
     */
    GetForwardersInBlock(blockID int64) ([]*ForwardInput, error)

    /**
     * SendWithdrawl
     *
     * Send the withdrawal to the contract.
     *  - w is the native withdrawal
     *  - s is w signed by the node's private key
     *  - if w and s do not match this will fail.
     *
     * Return error if one was encountered
     */
    SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error
}

func NewDummyCoin() (Coin, error) {
    return &DummyCoin{}, nil
}


func NewEthereumCoin() (Coin, error) {
    return &EthereumCoin{}, nil
}
