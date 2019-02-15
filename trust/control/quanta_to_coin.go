package control

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/quantadex/quanta_book/consensus/cosi"
	"github.com/scorum/bitshares-go/apis/database"
	"math/big"
	"strings"
	"time"
)

const QUANTA = "QUANTA"
const DQ_QUANTA2COIN = "DQ_QUANTA2COIN"

type WithdrawalResult struct {
	W   *coin.Withdrawal
	Err error
	Tx  string
}

/**
 * QuantaToCoin
 *
 * This module polls for new quanta blocks and find refunds submitted to the quanta trust.
 * It validates and signs the the withdrawls and issues them to the coin's smart contract.
 */
type QuantaToCoin struct {
	logger              logger.Logger
	db                  kv_store.KVStore
	rDb                 *db.DB
	coinChannel         map[string]coin.Coin
	quantaChannel       quanta.Quanta
	quantaTrustAddress  string
	coinContractAddress string
	coinkM              map[string]key_manager.KeyManager
	nodeID              int
	coinMapping         map[string]string
	coinInfo            map[string]*database.Asset

	rr        *RoundRobinSigner
	cosi      *cosi.Cosi
	trustPeer *peer_contact.TrustPeerNode
	doneChan  chan bool
	SuccessCb func(WithdrawalResult)
}

/**
 * NewQuantaToCoin
 *
 * Create a new instance of the class. This does not initialize modules.
 * All modules must already have been initialized and passed in here.
 */
func NewQuantaToCoin(log logger.Logger,
	db_ kv_store.KVStore,
	rDb *db.DB,
	c map[string]coin.Coin,
	q quanta.Quanta,
	man *manifest.Manifest,
	quantaTrustAddress string,
	coinContractAddress string,
	kM map[string]key_manager.KeyManager,
	coinMapping map[string]string,
	peer peer_contact.PeerContact,
	queue_ queue.Queue,
	nodeID int,
	coinInfo map[string]*database.Asset) *QuantaToCoin {
	res := &QuantaToCoin{}
	res.logger = log
	res.db = db_
	res.rDb = rDb
	res.coinChannel = c
	res.quantaChannel = q
	res.quantaTrustAddress = quantaTrustAddress
	res.coinContractAddress = coinContractAddress
	res.coinkM = kM
	res.coinMapping = coinMapping
	res.nodeID = nodeID
	res.doneChan = make(chan bool, 1)
	res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_, queue.REFUNDMSG_QUEUE, "/node/api/refund")
	res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

	res.coinInfo = coinInfo

	res.cosi.Verify = func(msg string) error {
		var encoded coin.EncodedMsg
		json.Unmarshal([]byte(msg), &encoded)
		blockchain, b := res.getBlockchainForCoin(encoded.CoinName)
		if !b {
			log.Error("Unable to verify refund coin name" + encoded.CoinName)
		} else {
			withdrawal, err := res.coinChannel[blockchain].DecodeRefund(msg)
			if err != nil {
				log.Error("Unable to decode refund")
				return err
			}
			tx, err := db.GetTransaction(rDb, withdrawal.Tx)
			if tx != nil {
				// we're not going to sign again
				if tx.Signed == false {
					return nil
				}
				log.Error("Unable to verify refund " + tx.Tx)
			}

			return errors.New("Unable to verify: " + withdrawal.Tx + " " + err.Error())
		}
		return nil
	}

	res.cosi.Persist = func(msg string) error {
		var encoded coin.EncodedMsg
		json.Unmarshal([]byte(msg), &encoded)
		blockchain, b := res.getBlockchainForCoin(encoded.CoinName)
		if !b {
			return nil
		} else {
			withdrawal, err := res.coinChannel[blockchain].DecodeRefund(msg)
			if err != nil {
				log.Error("Unable to decode refund at persist")
				return err
			}

			err = db.SignWithdrawal(rDb, withdrawal)
			if err != nil {
				res.logger.Error("Failed to confirm transaction: " + err.Error())
				return errors.New("Failed to confirm transaction: " + err.Error())
			}
			return nil
		}
		return nil
	}

	res.cosi.SignMsg = func(encoded string) (string, error) {
		//decoded := common.Hex2Bytes(encoded)
		msg := &coin.EncodedMsg{}
		err := json.Unmarshal([]byte(encoded), msg)
		if err != nil {
			return "", err
		}

		blockchain, b := res.getBlockchainForCoin(msg.CoinName)
		if !b {
			return "", nil
		} else {

		}

		encodedSig, err := res.coinkM[blockchain].SignTransaction(msg.Message)
		log.Infof("Sign msg %s -> %s", msg.Message, encodedSig)

		if err != nil {
			log.Error("Unable to Sign/encode refund msg")
			return "", err
		}

		return encodedSig, nil
	}

	res.cosi.Start()

	go func() {
		for {
			// feed messages back
			msg := res.trustPeer.GetMsg()

			if msg != nil {
				//log.Infof("Got peer message %v", msg.Signed)
				res.cosi.CosiMsgChan <- msg
				continue
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	go res.DispatchWithdrawal()

	return res
}

func (c *QuantaToCoin) GetNewCoinBlockIDs() []int64 {
	lastProcessed, valid := GetLastBlock(c.db, QUANTA)
	if !valid {
		c.logger.Error("Failed to get last processed ID")
		return nil
	}

	currentTop, err := c.quantaChannel.GetTopBlockID()
	if err != nil {
		c.logger.Error("Failed to get current top block")
		return nil
	}

	if lastProcessed > currentTop {
		c.logger.Error("Coin top block smaller than last processed")
		return nil
	}

	if lastProcessed == currentTop {
		c.logger.Debug(fmt.Sprintf("Quanta2Coin: No new block last=%d top=%d", lastProcessed, currentTop))
		return nil
	}

	blocks := make([]int64, 0)
	for i := common2.MaxInt64(0, lastProcessed+1); i <= currentTop; i++ {
		blocks = append(blocks, i)
		if len(blocks) == 3*MAX_PROCESS_BLOCKS {
			break
		}
	}
	c.logger.Info(fmt.Sprintf("Quanta2Coin: Got blocks %v", blocks))

	return blocks
}

func (c *QuantaToCoin) getBlockchainForCoin(coinName string) (string, bool) {
	if strings.Contains(coinName, "0X") {
		return coin.BLOCKCHAIN_ETH, true
	}
	for k, v := range c.coinMapping {
		if v == coinName {
			return k, true
		}
	}
	return "", false
}
func (c *QuantaToCoin) DispatchWithdrawal() {
	var lastBlock int64 = 0
	for {
		select {
		case <-time.After(time.Second * 10):
			txs := db.QueryWithdrawalByAge(c.rDb, time.Now().Add(-time.Second*5), []string{db.SUBMIT_CONSENSUS})

			if len(txs) > 0 {
				blockchain, ok := c.getBlockchainForCoin(txs[0].Coin)
				if !ok {
					c.logger.Errorf("Blockchain not found for %s", txs[0].Coin)
				} else {
					currentBlock, err := c.coinChannel[blockchain].GetTopBlockID()
					if err != nil {
						c.logger.Error(err.Error())

					} else {
						//to avoid multiple transactions in one block
						if currentBlock > lastBlock+1 {
							lastBlock = currentBlock
							w := &coin.Withdrawal{
								Tx:                 txs[0].Tx,
								TxId:               txs[0].TxId,
								CoinName:           txs[0].Coin,
								SourceAddress:      txs[0].From,
								DestinationAddress: txs[0].To,
								QuantaBlockID:      txs[0].BlockId,
								Amount:             uint64(txs[0].Amount),
							}
							c.StartConsensus(w)
						}
					}
				}

			}
		case <-c.doneChan:
			c.logger.Infof("Exiting.")
			break
		}
	}
}

func (c *QuantaToCoin) Stop() {
	c.doneChan <- true
}

func (c *QuantaToCoin) GetCrosschainAddress(w *coin.Withdrawal, blockchain string) map[string]string {
	crossAddr := c.rDb.GetCrosschainByBlockchain(blockchain)

	watchMap := make(map[string]string)

	for _, w := range crossAddr {
		watchMap[w.Address] = w.QuantaAddr
	}
	return watchMap
}

func (c *QuantaToCoin) StartConsensus(w *coin.Withdrawal) (string, error) {
	blockchain, ok := c.getBlockchainForCoin(w.CoinName)
	if !ok {
		c.logger.Errorf("Blockchain not found for %s", w.CoinName)
		return "", nil
	}

	c.coinChannel[blockchain].FillCrosschainAddress(c.GetCrosschainAddress(w, blockchain))

	txResult := HEX_NULL
	errResult := error(nil)

	txId, err := c.coinChannel[blockchain].GetTxID(common.HexToAddress(c.coinContractAddress))
	if err != nil {
		c.logger.Error("Could not get txID: " + err.Error() + " " + c.coinContractAddress)
		//TODO: How to handle this?
	}
	w.TxId = txId + 1
	c.logger.Infof("Start new round %s %s to=%s amount=%d", w.Tx, w.CoinName, w.DestinationAddress, w.Amount)
	encoded, err := c.coinChannel[blockchain].EncodeRefund(*w)

	if err != nil {
		c.logger.Error("Failed to encode refund " + err.Error())
		db.ChangeSubmitState(c.rDb, w.Tx, db.ENCODE_FAILURE, db.WITHDRAWAL)
		return HEX_NULL, err
	}

	// wait for other node to see the tx
	c.cosi.StartNewRound(encoded)

	result := <-c.cosi.FinalSigChan

	if result.Msg == nil {
		errResult = errors.New("Unable to sign for refund")
		c.logger.Error("Unable to sign for refund")
	} else {
		// save this to queue for later in case ETH RPC is down.
		w.Signatures = result.Msg
		// save in eth_tx_log_signed (kvstore) [S=signed,X=submitted,F=failed(uncoverable), R=retry(connection failed)] ; recoverable=RPC not available
		c.logger.Infof("Great! Cosi successfully signed refund")
		//c.SubmitWithdrawal(&w)
		tx, err := c.coinChannel[blockchain].SendWithdrawal(common.HexToAddress(c.coinContractAddress), c.coinkM[blockchain].GetPrivateKey(), w)

		if err != nil {
			c.logger.Errorf("Error submission: %s", err.Error())
			if strings.Contains(err.Error(), "known transaction") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_FATAL, db.WITHDRAWAL)
			} else if strings.Contains(err.Error(), "replacement transaction underpriced") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_FATAL, db.WITHDRAWAL)
			} else if strings.Contains(err.Error(), "connect: connection refused") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_RECOVERABLE, db.WITHDRAWAL)
			} else if strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_RECOVERABLE, db.WITHDRAWAL)
			}
		} else {
			db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_SUCCESS, db.WITHDRAWAL)
			c.logger.Infof("Submitted withdrawal in tx=%s SUCCESS", tx)
		}

		txResult = tx
		errResult = err

		if c.SuccessCb != nil {
			c.SuccessCb(WithdrawalResult{w, errResult, tx})
		}
	}
	return txResult, errResult
}

// VERY IMPORTANT CODE
func (c *QuantaToCoin) ComputeAmountToGraphene(coinName string, amount uint64) uint64 {
	// this is ETH, so we have to convert from system precision standard precision (5)
	for _, v := range c.coinInfo {
		if coinName == v.Symbol {
			return uint64(coin.PowerDelta(*big.NewInt(int64(amount)), int(v.Precision), 5))
		}
	}
	//if coinName == c.coinName {
	//	return uint64(coin.PowerDelta(*big.NewInt(int64(amount)), int(c.coinInfo.Precision), 5))
	//}

	return amount
}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 * returns []RefunTestTransactionQueryTestTransactionQueryd, txId string, error
 *     txId is the transaction id hex string if the withdrawal was a success, otherwise 0x0
 */
func (c *QuantaToCoin) DoLoop(cursor int64) ([]quanta.Refund, error) {
	refunds, _, err := c.quantaChannel.GetRefundsInBlock(cursor, c.quantaTrustAddress)

	if err != nil {
		c.logger.Error(err.Error())
		return refunds, err
	}
	//coin.StellarToWei()

	errResult := error(nil)

	c.logger.Debugf("QuantaToCoin refunds %v", refunds)

	// separate confirm, and sign as two different stages
	// refund gives issued token, withdrawal can convert into blockchain
	for _, refund := range refunds {
		// REPLACE WITH COMMON CODE
		blockchain, ok := c.getBlockchainForCoin(refund.CoinName)
		if !ok {
			c.logger.Errorf("Blockchain not found for %s", refund.CoinName)
			continue
		}

		c.logger.Infof("Confirm Refund tx=%s pt=%d", refund.TransactionId, refund.PageTokenID)
		w := &coin.Withdrawal{
			Tx:                 refund.TransactionId,
			CoinName:           refund.CoinName,
			Blockchain:         blockchain,
			SourceAddress:      refund.SourceAddress,
			DestinationAddress: refund.DestinationAddress,
			QuantaBlockID:      refund.OperationID,
			// TODO: Potentially losing precision when converting to wei
			Amount: c.ComputeAmountToGraphene(refund.CoinName, refund.Amount),
		}

		db.ConfirmWithdrawal(c.rDb, w)
		//cursor = refund.PageTokenID

		if w.DestinationAddress == "0x0000000000000000000000000000000000000000" || !c.coinChannel[blockchain].CheckValidAddress(w.DestinationAddress) {
			// create a deposit
			dep := &coin.Deposit{
				Tx:         refund.TransactionId,
				CoinName:   refund.CoinName, // coin,issuer
				SenderAddr: "",
				QuantaAddr: refund.SourceAddress,
				Amount:     int64(refund.Amount),
				BlockID:    int64(refund.LedgerID),
			}
			// mark as a bounced transaction
			err := db.ConfirmDeposit(c.rDb, dep, true)
			if err != nil {
				c.logger.Error("Cannot insert into db:" + err.Error())
			} else if c.nodeID == 0 {
				err = db.ChangeSubmitState(c.rDb, dep.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT)
				if err != nil {
					c.logger.Error("Cannot change submit state:" + err.Error())
				}
			}

			c.logger.Error("Refund is missing destination address, skipping.")

		} else if w.Amount == 0 {
			c.logger.Error("Amount is too small")
		} else if c.nodeID == 0 {
			db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_CONSENSUS, db.WITHDRAWAL)
		}

	}

	// mark the block for the next loop through
	success := SetLastBlock(c.db, QUANTA, cursor)
	if !success {
		c.logger.Error("Failed to mark block as signed")
		return refunds, errors.New("Failed to mark block as signed")
	}

	c.logger.Debugf("Next cursor is = %d, numRefunds=%d", cursor, len(refunds))

	return refunds, errResult
}
