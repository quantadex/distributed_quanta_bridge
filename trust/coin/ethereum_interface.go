package coin

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/stellar/go/support/log"
	"math/big"
	"regexp"
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
	weiInEth               = new(big.Rat).SetInt(new(big.Int).Exp(ten, eighteen, nil))
	StellarAmountPrecision = new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(7), nil))
	ROPSTEN_NETWORK_ID     = "3"
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
// persistent db.
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

func WeiToStellar(valueInWei big.Int) int64 {
	valueEth := new(big.Rat).SetInt(&valueInWei)
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(11), nil))
	result := new(big.Rat)
	result = result.Quo(valueEth, powerDelta)
	num, _ := new(big.Int).SetString(result.FloatString(0), 10)
	return num.Int64()
}

func WeiToGraphene(valueInWei big.Int) int64 {
	valueEth := new(big.Rat).SetInt(&valueInWei)
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(13), nil))
	result := new(big.Rat)
	result = result.Quo(valueEth, powerDelta)
	num, _ := new(big.Int).SetString(result.FloatString(0), 10)
	return num.Int64()
}

func PowerDelta(value big.Int, curPrecision int, targetPrecision int) int64 {
	valueEth := new(big.Rat).SetInt(&value)
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(int64(common2.AbsInt(curPrecision-targetPrecision))), nil))
	result := new(big.Rat)
	if targetPrecision < curPrecision {
		result = result.Quo(valueEth, powerDelta)
	} else {
		result = result.Mul(valueEth, powerDelta)
	}
	num, _ := new(big.Int).SetString(result.FloatString(0), 10)
	return num.Int64()
}

func Erc20AmountToGraphene(valueInWei big.Int, dec uint8) int64 {
	valueEth := new(big.Rat).SetInt(&valueInWei)
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(18-int64(dec)), nil))
	result := new(big.Rat)
	result = result.Mul(valueEth, powerDelta)

	return WeiToGraphene(*result.Num())
}

func StellarToWei(valueInStellar uint64) *big.Int {
	valueWei := new(big.Rat)
	stellar := new(big.Rat).SetInt(big.NewInt(int64(valueInStellar)))
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(11), nil))

	return valueWei.Mul(stellar, powerDelta).Num()
}

func GrapheneToWei(valueInGraphene uint64) *big.Int {
	valueWei := new(big.Rat)
	stellar := new(big.Rat).SetInt(big.NewInt(int64(valueInGraphene)))
	powerDelta := new(big.Rat).SetInt(new(big.Int).Exp(ten, big.NewInt(13), nil))

	return valueWei.Mul(stellar, powerDelta).Num()
}

func CheckValidEthereumAddress(address string) bool {
	var add [20]byte
	copy(add[:], address)
	Ma := common.NewMixedcaseAddress(add)
	var validAddress = regexp.MustCompile(`^[0x]+[0-9a-fA-F]{40}$`)
	return Ma.ValidChecksum() && validAddress.MatchString(address)
}

type ForwardInput struct {
	ContractAddress common.Address
	Trust           common.Address
	QuantaAddr      string
	TxHash          string
	Blockchain      string
}
