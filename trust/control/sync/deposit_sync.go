package sync

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/scorum/bitshares-go/apis/database"
	"time"
)

/**
 * DepositSync is an interface which describes the lifecycle of syncing blockchain deposits into the database.
 */
type DepositSync struct {
	coinChannel        coin.Coin
	quantaChannel      quanta.Quanta // stellar -> graphene
	coinInfo           map[string]*database.Asset
	db                 kv_store.KVStore
	rDb                *db.DB
	logger             logger.Logger
	blockStartID       int64
	fnDepositInBlock   func(blockID int64) ([]*coin.Deposit, error)
	fnPostProcessBlock func(blockID int64) error

	doneChan chan bool
}

func NewDepositSync(coin coin.Coin,
	quantaChannel quanta.Quanta,
	issuingSymbol map[string]string,
	db kv_store.KVStore,
	rDb *db.DB,
	logger logger.Logger,
	blockStartID int64) *DepositSync {

	//coinInfo, _ := quantaChannel.GetAsset(issuingSymbol[""])
	coinInfo := make(map[string]*database.Asset)
	for _, v := range issuingSymbol {
		asset, _ := quantaChannel.GetAsset(v)
		coinInfo[v] = asset
	}

	return &DepositSync{
		coinChannel:   coin,
		quantaChannel: quantaChannel,
		coinInfo:      coinInfo,
		db:            db,
		rDb:           rDb,
		doneChan:      make(chan bool, 1),
		logger:        logger,
		blockStartID:  blockStartID,
	}
}

// run 1 iteration of syncing
func (c *DepositSync) DoLoop(blockIDs []int64) []*coin.Deposit {
	c.logger.Debugf("***** Start # of blocks=%d *** ", len(blockIDs))

	allDeposits := make([]*coin.Deposit, 0)

	if blockIDs != nil {
		for _, blockID := range blockIDs {
			deposits, err := c.fnDepositInBlock(blockID)
			if err != nil {
				//Retry the block that was not found
				if err.Error() == "Block not found" {
					lastBlockId := blockIDs[len(blockIDs)-2]
					control.SetLastBlock(c.db, c.coinChannel.Blockchain(), lastBlockId)
				}
				c.logger.Error("Failed to get deposits from block: " + err.Error())
				return allDeposits
			}

			if deposits != nil {
				if len(deposits) > 0 {
					c.logger.Info(fmt.Sprintf("Block %d Got deposits %d %v", blockID, len(deposits), deposits))
				}

				for _, dep := range deposits {
					// every node must mark the deposit
					err = db.ConfirmDeposit(c.rDb, dep, false)

					if err != nil {
						c.logger.Error("Cannot insert into db:" + err.Error())
					}
					allDeposits = append(allDeposits, dep)

					if !c.quantaChannel.AccountExist(dep.QuantaAddr) {
						// if not exist, let's bounce money back
					} else if dep.Amount == 0 {
						c.logger.Error("Amount is too small")
					}
					// have leader pick up automatically
					//else if c.nodeID == 0 {
					//	db.ChangeSubmitState(c.rDb, dep.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT)
					//}
				}
			}

			if c.fnPostProcessBlock != nil {
				err = c.fnPostProcessBlock(blockID)
				if err != nil {
					c.logger.Error(err.Error())
				}
			}
		}
	}

	if len(blockIDs) > 0 {
		lastBlockId := blockIDs[len(blockIDs)-1]
		c.logger.Debugf("set last block coin=%s height=%d", c.coinChannel.Blockchain(), lastBlockId)
		control.SetLastBlock(c.db, c.coinChannel.Blockchain(), lastBlockId)
	}

	return allDeposits
}

/**
 * getNewCoinBlockIDs
 *
 * Returns a list of new blocks added to the coin block chain.
 */
func (c *DepositSync) GetNewCoinBlockIDs() []int64 {
	lastProcessed, valid := control.GetLastBlock(c.db, c.coinChannel.Blockchain())
	if !valid {
		c.logger.Error("Failed to get last processed ID")
		return nil
	}

	currentTop, err := c.coinChannel.GetTopBlockID()
	if err != nil {
		c.logger.Error("Failed to get current top block")
		return nil
	}

	if lastProcessed > currentTop {
		c.logger.Error("Coin top block smaller than last processed")
		return nil
	}

	if lastProcessed == currentTop {
		c.logger.Debug(fmt.Sprintf("Coin2Quanta: No new block last=%d top=%d", lastProcessed, currentTop))
		return nil
	}
	blocks := make([]int64, 0)
	for i := common.MaxInt64(c.blockStartID, lastProcessed+1); i <= currentTop; i++ {
		blocks = append(blocks, i)
		if len(blocks) == control.MAX_PROCESS_BLOCKS {
			break
		}
	}
	c.logger.Info(fmt.Sprintf("Coin2Quanta: Got blocks %v", blocks))

	return blocks
}

// run in infinite loop
func (c *DepositSync) Run() {
	delayTime := time.Second
	//init := false
	for true {
		select {
		case <-time.After(delayTime):
			blockIDs := c.GetNewCoinBlockIDs()
			c.DoLoop(blockIDs)

			// scale up time
			if len(blockIDs) == control.MAX_PROCESS_BLOCKS {
				delayTime = time.Second
			} else {
				delayTime = time.Second * 3
			}

		case <-c.doneChan:
			c.logger.Infof("Exiting.")
			break
		}
	}
}

// stop the infinite loop
func (c *DepositSync) Stop() {
	c.doneChan <- true
}