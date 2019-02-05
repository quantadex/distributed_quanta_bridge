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
	DoLoop(blockIDs []int64) []*coin.Deposit
	GetNewCoinBlockIDs() []int64
	PostProcessBlock(blockID int64) error
	Run()
	Stop()
}

func NewEthereumSync(coin coin.Coin,
	trustAddress string,
	issuingSymbol string,
	quantaChannel quanta.Quanta,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64) DepositSyncInterface {

	parent := NewDepositSync(coin, quantaChannel, issuingSymbol, db, rDb, logger, blockStartID)
	eth := &EthereumSync{
		*parent,
		common.HexToAddress(trustAddress),
		issuingSymbol,
	}
	eth.Setup()

	return eth
}

func NewBitcoinSync(coin coin.Coin,
	issuingSymbol string,
	quantaChannel quanta.Quanta,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64) DepositSyncInterface {

	parent := NewDepositSync(coin, quantaChannel, issuingSymbol, db, rDb, logger, blockStartID)
	btc := &BitcoinSync{
		*parent,
	}
	btc.Setup()

	return btc
}
