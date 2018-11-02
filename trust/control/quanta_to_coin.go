package control

import (
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
	"encoding/json"
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
	db                  kv_store.KVStore
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
	deferQ 				*queue.DeferQ
}

/**
 * NewQuantaToCoin
 *
 * Create a new instance of the class. This does not initialize modules.
 * All modules must already have been initialized and passed in here.
 */
func NewQuantaToCoin(log logger.Logger,
	db kv_store.KVStore,
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
	res.db = db
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

		refKey := getKeyName(withdrawal.CoinName, withdrawal.DestinationAddress, withdrawal.QuantaBlockID)
		state := getState(db, QUANTA_CONFIRMED, refKey)
		if state == CONFIRMED {
			return nil
		}
		log.Error("Unable to verify refund " + refKey + " " + state)
		return errors.New("Unable to verify: " + refKey)
	}

	res.cosi.Persist = func(msg string) error {
		withdrawal, err := res.coinChannel.DecodeRefund(msg)
		if err != nil {
			log.Error("Unable to decode refund at persist")
			return err
		}

		refKey := getKeyName(withdrawal.CoinName, withdrawal.DestinationAddress, withdrawal.QuantaBlockID)

		success := signTx(db, QUANTA_CONFIRMED, refKey)
		if !success {
			res.logger.Error("Failed to confirm transaction")
			return errors.New("Failed to confirm transaction")
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

func (c *QuantaToCoin) WithdrawalSubmitter() {
	for {
		time.Sleep(time.Second)
	}
}

// Receives withdrawal with all of the signatures - queue it properly
// for submission, and retry as neccessary
func (c *QuantaToCoin) SubmitWithdrawal(w *coin.Withdrawal) {

}

/**
 * DoLoop
 *
 * Does one complete loop. Pulling all new quanta blocks and processing any refunds in them
 */
func (c *QuantaToCoin) DoLoop(cursor int64) {
	refunds, _, err := c.quantaChannel.GetRefundsInBlock(cursor, c.quantaTrustAddress)

	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	c.logger.Infof("QuantaToCoin Epoch=%d refunds %v", c.deferQ.Epoch(), refunds)

	// separate confirm, and sign as two different stages
	for _, refund := range refunds {
		refKey := getKeyName(refund.CoinName, strings.ToLower(common.HexToAddress(refund.DestinationAddress).Hex()), int64(refund.OperationID))
		c.logger.Infof("Confirm Refund = %s tx=%s pt=%d", refKey, refund.TransactionId, refund.PageTokenID)
		confirmTx(c.db, QUANTA_CONFIRMED, refKey)

		c.deferQ.Put(DQ_QUANTA2COIN, &refund)
		cursor = refund.PageTokenID
		success := setLastBlock(c.db, QUANTA, refund.PageTokenID)
		if !success {
			c.logger.Error("Failed to mark block as signed")
			return
		}
	}

	// TODO: make this process multiple refunds in one pass
	refundI, _ := c.deferQ.Get(DQ_QUANTA2COIN)

	if refundI != nil {
		refund := refundI.(*quanta.Refund)

		//TODO: do checksum check, should bounce back the payment
		if refund.DestinationAddress == "" {
			c.logger.Error("Refund is missing destination address, skipping.")
		}

		// i'm the leader
		if c.nodeID == 0 && refund.DestinationAddress != "" {
			txId, err := c.coinChannel.GetTxID(common.HexToAddress(c.coinContractAddress))
			if err != nil {
				c.logger.Error("Could not get txID: " + err.Error() + " " + c.coinContractAddress)
				//TODO: How to handle this?
			}
			w := coin.Withdrawal{
				TxId:               txId + 1,
				CoinName:           refund.CoinName,
				DestinationAddress: refund.DestinationAddress,
				QuantaBlockID:      refund.OperationID,
				Amount:             coin.StellarToWei(refund.Amount),
			}
			c.logger.Infof("Start new round %s to=%s amount=%d", w.CoinName, w.DestinationAddress, w.Amount)
			encoded, err := c.coinChannel.EncodeRefund(w)

			if err != nil {
				c.logger.Error("Failed to encode refund " + err.Error())
				return
			}

			// wait for other node to see the tx
			time.Sleep(time.Second * 3)
			c.cosi.StartNewRound(encoded)

			result := <-c.cosi.FinalSigChan

			if result.Msg == nil {
				c.logger.Error("Unable to sign for refund")
			} else {
				// save this to queue for later in case ETH RPC is down.
				w.Signatures = result.Msg
				// save in eth_tx_log_signed (kvstore) [S=signed,X=submitted,F=failed(uncoverable), R=retry(connection failed)] ; recoverable=RPC not available
				c.logger.Infof("Great! Cosi successfully signed refund")
				tx, err := c.coinChannel.SendWithdrawal(common.HexToAddress(c.coinContractAddress), c.coinkM.GetPrivateKey(), &w)
				if err != nil {
					c.logger.Error(err.Error())
				}
				c.logger.Infof("Submitted withdrawal in tx=%s", tx)
			}
		}
	}
	c.deferQ.AddTick()
	c.logger.Infof("Next cursor is = %d", cursor)
}
