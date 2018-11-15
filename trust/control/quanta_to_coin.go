package control

import "C"
import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/quantadex/quanta_book/consensus/cosi"
	"strconv"
	"strings"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
)

const QUANTA = "QUANTA"
const DQ_QUANTA2COIN = "DQ_QUANTA2COIN"

/**
 * QuantaToCoin
 *
 * This module polls for new quanta blocks and find refunds submitted to the quanta trust.
 * It validates and signs the the withdrawls and issues them to the coin's smart contract.
 */
type QuantaToCoin struct {
	logger              logger.Logger
	db					kv_store.KVStore
	rDb					*db.DB
	coinChannel         coin.Coin
	quantaChannel       quanta.Quanta
	quantaTrustAddress  string
	coinContractAddress string
	coinkM              key_manager.KeyManager
	coinName            string
	nodeID              int
	rr                  *RoundRobinSigner
	cosi                *cosi.Cosi
	trustPeer           *peer_contact.TrustPeerNode
	deferQ              *queue.DeferQ
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
	c coin.Coin,
	q quanta.Quanta,
	man *manifest.Manifest,
	quantaTrustAddress string,
	coinContractAddress string,
	kM key_manager.KeyManager,
	coinName string,
	peer peer_contact.PeerContact,
	queue_ queue.Queue,
	nodeID int) *QuantaToCoin {
	res := &QuantaToCoin{}
	res.logger = log
	res.db = db_
	res.rDb = rDb
	res.coinChannel = c
	res.quantaChannel = q
	res.quantaTrustAddress = quantaTrustAddress
	res.coinContractAddress = coinContractAddress
	res.coinkM = kM
	res.coinName = coinName
	res.nodeID = nodeID
	res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_)
	res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

	res.deferQ = queue.NewDeferQ(DELAY_PENALTY)
	res.deferQ.CreateQueue(DQ_QUANTA2COIN)

	res.cosi.Verify = func(msg string) error {
		withdrawal, err := res.coinChannel.DecodeRefund(msg)
		if err != nil {
			log.Error("Unable to decode refund")
			return err
		}
		tx := db.GetTransaction(rDb, withdrawal.Tx)
		if tx != nil {
			// we're not going to sign again
			if tx.Signed == false {
				return nil
			}
			log.Error("Unable to verify refund " + tx.Tx)
		}

		return errors.New("Unable to verify: " + withdrawal.Tx)
	}

	res.cosi.Persist = func(msg string) error {
		withdrawal, err := res.coinChannel.DecodeRefund(msg)
		if err != nil {
			log.Error("Unable to decode refund at persist")
			return err
		}

		err = db.SignWithdrawal(rDb, withdrawal)
		if err != nil {
			res.logger.Error("Failed to confirm transaction: " + err.Error())
			return errors.New("Failed to confirm transaction: "+ err.Error())
		}
		return nil
	}

	res.cosi.SignMsg = func(encoded string) (string, error) {
		decoded := common.Hex2Bytes(encoded)
		msg := &coin.EncodedMsg{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return "", err
		}

		encodedSig, err := res.coinkM.SignTransaction(msg.Message)
		log.Infof("Sign msg %s", encodedSig)

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
			msg := res.trustPeer.GetRefundMsg()

			if msg != nil {
				//log.Infof("Got peer message %v", msg.Signed)
				res.cosi.CosiMsgChan <- msg
				continue
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	return res
}

func (c *QuantaToCoin) WithdrawalLoop(w *coin.Withdrawal) {
	signed, _ := c.db.GetAllValues(kv_store.ETH_TX_LOG_SIGNED)
	for key, value := range signed {
		_, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.quantaTrustAddress), nil, w)
		if err == nil {
			c.db.SetValue(kv_store.ETH_TX_LOG_SUBMITTED, key, "", value)
			c.db.RemoveKey(kv_store.ETH_TX_LOG_SIGNED, key)
		} else if strings.Contains(err.Error(), "connect: connection refused") {
			c.db.SetValue(kv_store.ETH_TX_LOG_RETRY, key, "", value)
		} else if strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
			c.db.SetValue(kv_store.ETH_TX_LOG_FAILED_RECOVERABLE, key, "", value)
		} else if err.Error() == "Failed" {
			c.db.SetValue(kv_store.ETH_TX_LOG_FAILED_UNRECOVERABLE, key, "", value)
		}
	}
	retry, _ := c.db.GetAllValues(kv_store.ETH_TX_LOG_RETRY)
	for key, msg := range retry {
		_, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.quantaTrustAddress), nil, w)
		if err == nil {
			c.db.SetValue(kv_store.ETH_TX_LOG_SUBMITTED, key, "", msg)
			c.db.RemoveKey(kv_store.ETH_TX_LOG_RETRY, key)
			c.db.RemoveKey(kv_store.ETH_TX_LOG_SIGNED, key)
		}
	}
}

// Receives withdrawal with all of the signatures - queue it properly
// for submission, and retry as neccessary
func (c *QuantaToCoin) SubmitWithdrawal(w *coin.Withdrawal) {
	//InitLedger(c.db)

	key := strconv.Itoa(int(w.TxId))
	msg, _ := json.Marshal(w)
	value := string(msg)

	c.db.SetValue(kv_store.ETH_TX_LOG, key, "", value)
	c.db.SetValue(kv_store.ETH_TX_LOG_SIGNED, key, "", value)
	c.WithdrawalLoop(w)
}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 * returns []Refund, txId string, error
 *     txId is the transaction id hex string if the withdrawal was a success, otherwise 0x0
 */
func (c *QuantaToCoin) DoLoop(cursor int64) ([]quanta.Refund, string, error) {
	refunds, _, err := c.quantaChannel.GetRefundsInBlock(cursor, c.quantaTrustAddress)

	if err != nil {
		c.logger.Error(err.Error())
		return refunds, "0x0", err
	}

	txResult := "0x0"
	errResult := error(nil)

	c.logger.Infof("QuantaToCoin Epoch=%d refunds %v", c.deferQ.Epoch(), refunds)

	// separate confirm, and sign as two different stages
	for _, refund := range refunds {
		c.logger.Infof("Confirm Refund tx=%s pt=%d", refund.TransactionId, refund.PageTokenID)

		w := &coin.Withdrawal{
			Tx: 				refund.TransactionId,
			CoinName:           refund.CoinName,
			DestinationAddress: refund.DestinationAddress,
			QuantaBlockID:      refund.OperationID,
			Amount:             coin.StellarToWei(refund.Amount),
		}
		db.ConfirmWithdrawal(c.rDb, w)

		c.deferQ.Put(DQ_QUANTA2COIN, w)
		cursor = refund.PageTokenID

		success := setLastBlock(c.db, QUANTA, refund.PageTokenID)
		if !success {
			c.logger.Error("Failed to mark block as signed")
			return refunds, "0x0", errors.New("Failed to mark block as signed")
		}
	}

	// TODO: make this process multiple refunds in one pass
	for {
		refundI, _ := c.deferQ.Get(DQ_QUANTA2COIN)

		if refundI == nil {
			break
		}
		w := refundI.(*coin.Withdrawal)

		//TODO: do checksum check, should bounce back the payment
		if w.DestinationAddress == "" {
			c.logger.Error("Refund is missing destination address, skipping.")
		} else if w.Amount == uint64(0) {
			c.logger.Error("Amount is too small")
		} else if c.nodeID == 0 {
			txId, err := c.coinChannel.GetTxID(common.HexToAddress(c.coinContractAddress))
			if err != nil {
				c.logger.Error("Could not get txID: " + err.Error() + " " + c.coinContractAddress)
				//TODO: How to handle this?
			}
			w.TxId = txId + 1

			c.logger.Infof("Start new round %s to=%s amount=%d", w.CoinName, w.DestinationAddress, w.Amount)
			encoded, err := c.coinChannel.EncodeRefund(*w)

			if err != nil {
				c.logger.Error("Failed to encode refund " + err.Error())
				return refunds, "0x0", err
			}

			// wait for other node to see the tx
			time.Sleep(time.Second * 3)
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
				tx, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.coinContractAddress), c.coinkM.GetPrivateKey(), w)
				if err != nil {
					c.logger.Error(err.Error())
				}
				c.logger.Infof("Submitted withdrawal in tx=%s", tx)
				txResult = tx
				errResult = err
			}
		}
	}

	c.deferQ.AddTick()
	c.logger.Infof("Next cursor is = %d, numRefunds=%d, txId=%s", cursor, len(refunds), txResult)

	return refunds, txResult, errResult
}
