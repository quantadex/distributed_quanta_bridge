package coin

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	common2 "github.com/ethereum/go-ethereum/common"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	chaincfg3 "github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	chaincfg2 "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcutil"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/scorum/bitshares-go/types"
	"time"
)

const BLOCKCHAIN_ETH = "ETH"
const CONST_PRECISION = 1e5

/**
 * Deposit
 *
 * The deposit struct captures the data of making a deposit into the trust
 */
type Deposit struct {
	Tx         string
	Type       types.OpType
	CoinName   string // Type of coin (e.g. ETH)
	SenderAddr string
	QuantaAddr string   // Destination quanta acount
	Amount     int64    // Deposit size
	BlockID    int64    // The blockID in which this deposit was found
	Signatures []string // hex signatures via quanta
	BlockHash  string
}

/**
 * Withdrawal
 *
 * The data structure that needs to be filled out to do a succesful withdrawal.
 */
type Withdrawal struct {
	Tx                 string
	TxId               uint64 // The Node authorizing this
	CoinName           string // The issued coin
	Blockchain         string
	SourceAddress      string
	DestinationAddress string   // Where this money is going
	QuantaBlockID      int64    // Which block this transaction was processed in quanta
	Amount             uint64   // The withdrawal size
	Signatures         []string // hex signatures via ethereum
	BlockHash          string
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

	Blockchain() string

	/**
	 * GetTopBlockID
	 *
	 * Returns the ID of the newest block in the chain.
	 */
	GetTopBlockID() (int64, error)

	GetTxID(trustAddress common2.Address) (uint64, error)

	/**
	 * GetDepositsInBlock
	 *
	 * For the specified blockID return all transactions that sent money to the specified address.
	 * The quanta address is specified as data in the transaction.
	 * Returns as a list of deposits.
	 * Returns nil if no matching deposits.
	 * Returns error of error encountered
	 */
	GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error)

	GetPendingTx(map[string]string) ([]*Deposit, error)

	/**
	 * GetForwardersInBlock
	 *
	 * Forwarders are smart contracts that are pointing into our trust address
	 * with information about QUANTA Address
	 * We will record this in our KV later, to know where deposits came from.
	 */
	GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error)

	/**
	 * GenerateMultisig - this is for creating multisig wallet
	 */
	GenerateMultisig(accountId string) (string, error)

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
	SendWithdrawal(trustAddress common2.Address,
		ownerKey *ecdsa.PrivateKey,
		w *Withdrawal) (string, error)

	EncodeRefund(w Withdrawal) (string, error)
	DecodeRefund(encoded string) (*Withdrawal, error)

	FillCrosschainAddress(crosschainAddr map[string]string)
	FlushCoin(forwarder string, address string) error

	CheckValidAddress(address string) bool

	GetBlockInfo(hash string) (string, int64, error)

	GetBlockTime(blockId int64) (time.Time, error)

	CheckValidAmount(amount uint64) bool

	SetIssuerAddress(address string)
}

func NewDummyCoin() (Coin, error) {
	return &DummyCoin{}, nil
}

func NewEthereumCoin(networkId string, ethereumRpc string, secret string, erc20map map[string]string, withdrawMin, withdrawFee float64, gasFee int64, blackList map[string]bool) (Coin, error) {
	key, err := crypto2.HexToECDSA(secret)
	if err != nil {
		return nil, err
	}
	return &EthereumCoin{maxRange: common.MaxNumberInt64, networkId: networkId, ethereumRpc: ethereumRpc, ethereumSecret: key, erc20map: erc20map, EthWithdrawMin: withdrawMin, EthWithdrawFee: withdrawFee, EthWithdrawGasFee: gasFee, BlackList: blackList}, nil
}

func NewBitcoinCoin(rpcHost string, params *chaincfg.Params, signers []string, rpcUser, rpcPassword, grapheneSeedPrefix string, withdrawMin, withdrawFee float64, blackList map[string]bool) (Coin, error) {
	signersA := []btcutil.Address{}
	for _, s := range signers {
		addr, err := btcutil.DecodeAddress(s, params)
		if err != nil {
			panic("corrupted btc address")
		}
		signersA = append(signersA, addr)
	}

	return &BitcoinCoin{rpcHost: rpcHost, chaincfg: params, signers: signersA, rpcUser: rpcUser, rpcPassword: rpcPassword, grapheneSeedPrefix: grapheneSeedPrefix, BtcWithdrawMin: withdrawMin, BtcWithdrawFee: withdrawFee, BlackList: blackList}, nil
}

func NewLitecoinCoin(rpcHost string, params *chaincfg2.Params, signers []string, rpcUser, rpcPassword, grapheneSeedPrefix string, withdrawMin, withdrawFee float64, blackList map[string]bool) (Coin, error) {
	signersA := []ltcutil.Address{}
	for _, s := range signers {
		addr, err := ltcutil.DecodeAddress(s, params)
		if err != nil {
			panic("corrupted ltc address")
		}
		signersA = append(signersA, addr)
	}

	return &LiteCoin{rpcHost: rpcHost, chaincfg: params, signers: signersA, rpcUser: rpcUser, rpcPassword: rpcPassword, grapheneSeedPrefix: grapheneSeedPrefix, LtcWithdrawMin: withdrawMin, LtcWithdrawFee: withdrawFee, BlackList: blackList}, nil
}

func NewBCHCoin(rpcHost string, params *chaincfg3.Params, signers []string, rpcUser, rpcPassword, grapheneSeedPrefix string, withdrawMin, withdrawFee float64, blackList map[string]bool) (Coin, error) {
	signersA := []bchutil.Address{}
	for _, s := range signers {
		addr, err := bchutil.DecodeAddress(s, params)
		if err != nil {
			panic("corrupted bch address")
		}
		signersA = append(signersA, addr)
	}

	return &BCH{rpcHost: rpcHost, chaincfg: params, signers: signersA, rpcUser: rpcUser, rpcPassword: rpcPassword, grapheneSeedPrefix: grapheneSeedPrefix, BchWithdrawMin: withdrawMin, BchWithdrawFee: withdrawFee, BlackList: blackList}, nil
}

/**
 * Used for testing
 */
func NewEthereumCoinWithMax(max int64) (Coin, error) {
	return &EthereumCoin{maxRange: max}, nil
}
