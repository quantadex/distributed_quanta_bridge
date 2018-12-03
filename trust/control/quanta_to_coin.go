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
	"strings"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
)

const QUANTA = "QUANTA"
const DQ_QUANTA2COIN = "DQ_QUANTA2COIN"

type WithdrawalResult struct {
	W *coin.Withdrawal
	Err error
	Tx string
}
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
	doneChan			chan bool
	SuccessCb			func(WithdrawalResult)
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
	res.doneChan = make(chan bool, 1)
	res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_, queue.REFUNDMSG_QUEUE, "/node/api/refund")
	res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

	res.cosi.Verify = func(msg string) error {
		withdrawal, err := res.coinChannel.DecodeRefund(msg)
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


func (c *QuantaToCoin) DispatchWithdrawal() {
	for {
		select {
			case <-time.After(time.Second):
				txs := db.QueryWithdrawalByAge(c.rDb, time.Now().Add(-time.Second*5), []string{db.SUBMIT_CONSENSUS})
				for _, tx := range txs {
					w := &coin.Withdrawal{
						Tx: tx.Tx,
						TxId: tx.TxId,
						CoinName: tx.Coin,
						SourceAddress: tx.From,
						DestinationAddress: tx.To,
						QuantaBlockID: tx.BlockId,
						Amount: uint64(tx.Amount),
					}
					c.StartConsensus(w)
				}
			case <- c.doneChan:
				c.logger.Infof("Exiting.")
				break
		}
	}
}

func (c *QuantaToCoin) Stop() {
	c.doneChan <- true
}

func (c *QuantaToCoin) StartConsensus(w *coin.Withdrawal) (string, error) {
	txResult := HEX_NULL
	errResult := error(nil)

	txId, err := c.coinChannel.GetTxID(common.HexToAddress(c.coinContractAddress))
	if err != nil {
		c.logger.Error("Could not get txID: " + err.Error() + " " + c.coinContractAddress)
		//TODO: How to handle this?
	}
	w.TxId = txId + 1

	c.logger.Infof("Start new round %s %s to=%s amount=%d", w.Tx, w.CoinName, w.DestinationAddress, w.Amount)
	encoded, err := c.coinChannel.EncodeRefund(*w)

	if err != nil {
		c.logger.Error("Failed to encode refund " + err.Error())
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
		tx, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.coinContractAddress), c.coinkM.GetPrivateKey(), w)

		if err != nil {
			c.logger.Errorf("Error submission: %s", err.Error())
			if strings.Contains(err.Error(), "known transaction") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_FATAL)
			} else if strings.Contains(err.Error(), "replacement transaction underpriced") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_FATAL)
			} else if strings.Contains(err.Error(), "connect: connection refused") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_RECOVERABLE)
			} else if strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
				db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_RECOVERABLE)
			}
		} else {
			db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_SUCCESS)
			c.logger.Infof("Submitted withdrawal in tx=%s SUCCESS", tx)
		}

		txResult = tx
		errResult = err

		if c.SuccessCb != nil {
			c.SuccessCb(WithdrawalResult{ w, errResult, tx})
		}
	}
	return txResult, errResult
}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 * returns []Refund, txId string, error
 *     txId is the transaction id hex string if the withdrawal was a success, otherwise 0x0
 */
func (c *QuantaToCoin) DoLoop(cursor int64) ([]quanta.Refund, error) {
	refunds, _, err := c.quantaChannel.GetRefundsInBlock(cursor, c.quantaTrustAddress)

	if err != nil {
		c.logger.Error(err.Error())
		return refunds, err
	}

	errResult := error(nil)

	c.logger.Debugf("QuantaToCoin refunds %v", refunds)

	// separate confirm, and sign as two different stages
	for _, refund := range refunds {
		c.logger.Infof("Confirm Refund tx=%s pt=%d", refund.TransactionId, refund.PageTokenID)

		w := &coin.Withdrawal{
			Tx: 				refund.TransactionId,
			CoinName:           refund.CoinName,
			SourceAddress:      refund.SourceAddress,
			DestinationAddress: refund.DestinationAddress,
			QuantaBlockID:      refund.OperationID,
			Amount:             coin.StellarToWei(refund.Amount),
		}

		db.ConfirmWithdrawal(c.rDb, w)
		cursor = refund.PageTokenID

		if w.DestinationAddress == "" {
			c.logger.Error("Refund is missing destination address, skipping.")

		} else if w.Amount == uint64(0) {
			c.logger.Error("Amount is too small")
		} else if c.nodeID == 0 {
			db.ChangeSubmitState(c.rDb, w.Tx, db.SUBMIT_CONSENSUS)
		}

		success := setLastBlock(c.db, QUANTA, refund.PageTokenID)
		if !success {
			c.logger.Error("Failed to mark block as signed")
			return refunds, errors.New("Failed to mark block as signed")
		}
	}

	c.logger.Debugf("Next cursor is = %d, numRefunds=%d", cursor, len(refunds))

	return refunds, errResult
}
