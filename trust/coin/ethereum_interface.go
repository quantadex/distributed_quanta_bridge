package coin

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stellar/go/support/log"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

var (
	ten      = big.NewInt(10)
	eighteen = big.NewInt(18)
	// weiInEth = 10^18
	weiInEth = new(big.Rat).SetInt(new(big.Int).Exp(ten, eighteen, nil))
	StellarAmountPrecision =  new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(7), nil))
	ROPSTEN_NETWORK_ID = "3"
)

type TransactionHandler func(transaction Transaction) error

type Transaction struct {
	Hash string
	// Value in Wei
	ValueWei *big.Int
	To       string
}

type Client interface {
	NetworkID(ctx context.Context) (*big.Int, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

// Storage is an interface that must be implemented by an object using
// persistent storage.
type Storage interface {
	// GetEthereumBlockToProcess gets the number of Ethereum block to process. `0` means the
	// processing should start from the current block.
	GetEthereumBlockToProcess() (uint64, error)
	// SaveLastProcessedEthereumBlock should update the number of the last processed Ethereum
	// block. It should only update the block if block > current block in atomic transaction.
	SaveLastProcessedEthereumBlock(block int64) error
}

type Listener struct {
	Enabled            bool
	Client             Client  `inject:""`
	Storage            Storage `inject:""`
	NetworkID          string
	TransactionHandler TransactionHandler

	log *log.Entry
}

func WeiToStellar(valueInWei int64) int64 {
	valueEth := new(big.Rat)
	valueEth.Quo(new(big.Rat).SetInt(big.NewInt(valueInWei)), weiInEth)

	result := new(big.Rat)
	return result.Mul(valueEth, StellarAmountPrecision).Num().Int64()
}

type ForwardInput struct {
	address common.Address
	quantaAddress string
}

