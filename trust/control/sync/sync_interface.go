package sync

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
)

type DepositSyncInterface interface {
	Setup()
	GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error)
	GetWatchAddress() map[string]string
	DoLoop(blockIDs []int64) []*coin.Deposit
	GetNewCoinBlockIDs() []int64
	PostProcessBlock(blockID int64) error
	Run()
	Stop()
}

func NewEthereumSync(coin coin.Coin,
	trustAddress string,
	issuingSymbol map[string]string,
	quantaChannel quanta.Quanta,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64,
	ethFlush bool,
	ethMinConfirmation int64) DepositSyncInterface {

	parent := NewDepositSync(coin, quantaChannel, issuingSymbol, db, rDb, logger, blockStartID)
	eth := &EthereumSync{
		*parent,
		common.HexToAddress(trustAddress),
		issuingSymbol,
		ethFlush,
		ethMinConfirmation,
	}
	eth.Setup()

	return eth
}

func NewBitcoinSync(coin coin.Coin,
	issuingSymbol map[string]string,
	quantaChannel quanta.Quanta,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64,
	btcMinConfirmation int64) DepositSyncInterface {

	parent := NewDepositSync(coin, quantaChannel, issuingSymbol, db, rDb, logger, blockStartID)
	btc := &BitcoinSync{
		*parent,
		issuingSymbol,
		btcMinConfirmation,
	}
	btc.Setup()

	return btc
}

func NewLitecoinSync(coin coin.Coin,
	issuingSymbol map[string]string,
	quantaChannel quanta.Quanta,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64,
	ltcMinConfirmation int64) DepositSyncInterface {

	parent := NewDepositSync(coin, quantaChannel, issuingSymbol, db, rDb, logger, blockStartID)
	ltc := &LitecoinSync{
		*parent,
		issuingSymbol,
		ltcMinConfirmation,
	}
	ltc.Setup()

	return ltc
}
