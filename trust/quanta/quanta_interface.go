package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/scorum/bitshares-go/apis/database"
)

const QUANTA_PRECISION = 10000000

/**
 * Refund
 *
 * User's quanta return to trust in order to get a refund
 */
type Refund struct {
	TransactionId      string
	LedgerID           int32
	OperationID        int64
	PageTokenID        int64 // use this as your blockID
	CoinName           string
	SourceAddress      string
	DestinationAddress string // extract from memo
	Amount             uint64
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
	AttachQueue(kv kv_store.KVStore) error

	/**
	 * GetTopBlockID
	 *
	 * Returns the id of the latest quanta block.
	 */
	GetTopBlockID() (int64, error)

	/**
	 * GetRefundsInBlock
	 *
	 * Returns a list of refunds that were made to the specified address in the given block.
	 * Return nil if no matching deposits.
	 */
	GetRefundsInBlock(blockID int64, trustAddress string) ([]Refund, int64, error)

	/**
	 * ProcessDeposit
	 *
	 * Once enough nodes have signed the deposit the last node sends it to quanta to
	 * transfer the funds into the user's quanta account
	 */
	ProcessDeposit(deposit *coin.Deposit, proposed string) error

	GetBalance(assetName string, quantaAddress string) (float64, error)
	GetAllBalances(quantaAddress string) (map[string]float64, error)
	DecodeTransaction(base64 string) (*coin.Deposit, error)
	Broadcast(stx string) error

	CreateTransferProposal(dep *coin.Deposit) (string, error)
	CreateNewAssetProposal(issuer string, symbol string, precision uint8) (string, error)
	CreateIssueAssetProposal(dep *coin.Deposit) (string, error)
	AssetExist(issuer string, symbol string) (bool, error)
	AccountExist(quantaAddr string) bool
	GetAsset(assetName string) (*database.Asset, error)
}

func NewQuanta(options QuantaClientOptions) (Quanta, error) {
	return &QuantaClient{QuantaClientOptions: options}, nil
}

func NewQuantaGraphene(options QuantaClientOptions) (Quanta, error) {
	return &QuantaGraphene{QuantaClientOptions: options}, nil
}

/**
 * Submitworker's job is to submit the deposits into
 * the horizon service, and retry as neccessary
 * Decouples whether horizon is online or not
 */
type SubmitWorker interface {
	Dispatch()
	//AttachQueue(kv queue.Queue) error
	AttachQueue(kv kv_store.KVStore) error
}

func NewSubmitWorker(options QuantaClientOptions) SubmitWorker {
	return &SubmitWorkerImpl{QuantaClientOptions: options}
}
