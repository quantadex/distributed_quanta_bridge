package control

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
	"github.com/scorum/bitshares-go/types"
	"github.com/zserge/metric"
	"strconv"
	"strings"
	"time"
)

type DepositResult struct {
	D   *coin.Deposit
	Err error
	Tx  string
}

type ConsensusType int

const (
	NEWASSET_CONSENSUS = iota
	TRANSFER_CONSENSUS
	ISSUE_CONSENSUS
)

/**
 * CoinToQuanta
 *
 * This modules receives new deposits made to the coin trust
 * and using the round robin module creates transactions in quanta
 */
type CoinToQuanta struct {
	logger        logger.Logger
	quantaChannel quanta.Quanta // stellar -> graphene
	db            kv_store.KVStore
	rDb           *db.DB
	man           *manifest.Manifest
	peer          peer_contact.PeerContact
	trustPeer     *peer_contact.TrustPeerNode
	cosi          *cosi.Cosi
	counter       metric.Metric

	readyChan chan bool
	doneChan  chan bool
	SuccessCb func(DepositResult)
	nodeID    int
	C2QOptions
	quantaOptions quanta.QuantaClientOptions
}

type C2QOptions struct {
	EthTrustAddress string
	BlockStartID    int64
}

const MAX_PROCESS_BLOCKS = 100

/**
 * NewCoinToQuanta
 *
 * Returns a new instance of the module
 * Initializes nothing so it should all be already initialized.
 */
func NewCoinToQuanta(log logger.Logger,
	db_ kv_store.KVStore,
	rDb *db.DB,
	c coin.Coin,
	q quanta.Quanta,
	man *manifest.Manifest,
	kM key_manager.KeyManager,
	nodeID int,
	peer peer_contact.PeerContact,
	queue_ queue.Queue,
	options C2QOptions,

	quantaOptions quanta.QuantaClientOptions) *CoinToQuanta {
	res := &CoinToQuanta{C2QOptions: options}
	res.logger = log
	res.quantaChannel = q
	res.db = db_
	res.rDb = rDb
	res.man = man
	res.nodeID = nodeID
	res.peer = peer
	res.doneChan = make(chan bool, 1)
	res.readyChan = make(chan bool, 1)
	res.quantaOptions = quantaOptions
	res.counter = metric.NewCounter("24h1m")
	if res.nodeID == 0 {
		expvar.Publish("deposit_status", res.counter)
	}

	res.trustPeer = peer_contact.NewTrustPeerNode(man, peer, nodeID, queue_, queue.PEERMSG_QUEUE, "/node/api/peer")
	res.cosi = cosi.NewProtocol(res.trustPeer, nodeID == 0, time.Second*3)

	res.cosi.Verify = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &coin.EncodedMsg{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		deposit, err := res.quantaChannel.DecodeTransaction(msg.Message)
		if err != nil {
			log.Error("Unable to decode quanta tx")
			return err
		}

		deposit.Tx = msg.Tx

		if err != nil {
			log.Error("Unable to decode quanta tx")
			return err
		}

		tx, err := db.GetTransaction(rDb, deposit.Tx)
		if err != nil {
			return err
		}

		if tx != nil {
			// we're not going to sign again
			if tx.Signed == false {
				return nil
			}
			log.Error("Unable to verify refund " + tx.Tx)
		}

		return errors.New("Unable to verify: " + deposit.Tx)
	}

	res.cosi.Persist = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &coin.EncodedMsg{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		deposit, err := res.quantaChannel.DecodeTransaction(msg.Message)
		if err != nil {
			log.Error("Unable to decode quanta tx")
			return err
		}

		deposit.Tx = msg.Tx

		if err != nil {
			log.Error("Unable to decode refund at persist")
			return err
		}

		// if it is create, don't bother marking it, because it's okay to sign multiple time since
		// asset can only be created once.
		if deposit.Type != types.CreateAssetOpType {
			err = db.SignDeposit(rDb, deposit)
			if err != nil {
				return errors.New("Failed to confirm transaction: " + err.Error())
			}
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

		encodedSig, err := kM.SignTransaction(msg.Message)
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
			msg := res.trustPeer.GetMsg()

			if msg != nil {
				//log.Infof("Got peer message %v", msg.Signed)
				res.cosi.CosiMsgChan <- msg
				continue
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	if nodeID == 0 {
		go res.dispatchIssuance()
	}

	return res
}

func (c *CoinToQuanta) processDeposits() {
	// only leader - pick up  deposits with empty or null submit_state
	txs := db.QueryDepositByAge(c.rDb, time.Now().Add(-time.Second*5), []string{db.SUBMIT_CONSENSUS})
	if len(txs) > 0 {
		c.counter.Add(1)
		tx := txs[0]
		w := &coin.Deposit{
			Tx:         tx.Tx,
			CoinName:   tx.Coin,
			QuantaAddr: tx.To,
			BlockID:    tx.BlockId,
			SenderAddr: tx.From,
			Amount:     tx.Amount,
			BlockHash:  tx.BlockHash,
		}
		// if not a native token, we need to flush it
		//if tx.Coin != c.coinName {
		//	parts := strings.Split(c.coinName, "0X")
		//	if len(parts) > 1 {
		//		// flush
		//		// contract := parts[1]
		//	}
		//}

		// check if asset exists
		//if not, then propose new asset
		exist, err := c.quantaChannel.AssetExist(c.quantaOptions.Issuer, tx.Coin)
		if err != nil {
			if err.Error() == "issuer do not match" {
				db.ChangeSubmitState(c.rDb, tx.Tx, db.DUPLICATE_ASSET, db.DEPOSIT, tx.BlockHash)
				return
			}
			c.logger.Error(err.Error())
		}

		if !exist {
			_, err = c.StartConsensus(w, NEWASSET_CONSENSUS)
			if err != nil {
				c.logger.Error("failed to create asset, error = " + err.Error())
			}
		} else {
			fmt.Println("asset exists")
		}

		// if newasset was created successfully
		if err == nil {
			time.Sleep(3 * time.Second)

			if tx.IsBounced {
				c.StartConsensus(w, TRANSFER_CONSENSUS)
			} else {
				c.StartConsensus(w, ISSUE_CONSENSUS)
			}
		}
	}
}

func (c *CoinToQuanta) processSubmissions() {
	data := db.QueryDepositByAge(c.rDb, time.Now(), []string{db.SUBMIT_QUEUE})
	if len(data) > 0 {
		c.logger.Errorf("processSubmissions pending=%d", len(data))
	}

	for k, v := range data {

		c.logger.Infof("Submit TX: %s signed=%v %v", v.Tx, v.Signed, v.SubmitTx)

		resp, err := c.quantaChannel.Broadcast(v.SubmitTx)
		if err != nil {
			msg := quanta.ErrorString(err, false)
			c.logger.Error("could not submit transaction " + msg)
			if strings.Contains(msg, "tx_bad_seq") || strings.Contains(msg, "op_malformed") {
				db.ChangeSubmitState(c.rDb, v.Tx, db.SUBMIT_FATAL, db.DEPOSIT, v.BlockHash)
			}
		} else {
			c.logger.Infof("Successful tx submission %s,remove %s", "", k)

			txHash := strconv.Itoa(int(resp.BlockNum)) + "_" + strconv.Itoa(int(resp.TrxNum))
			err = db.ChangeDepositSubmitState(c.rDb, v.Tx, db.SUBMIT_SUCCESS, int(resp.BlockNum), txHash, v.BlockHash)
			//err = db.ChangeSubmitState(c.rDb, v.Tx, db.SUBMIT_SUCCESS, db.DEPOSIT)
			if err != nil {
				c.logger.Error("Error removing key=" + v.Tx)
			}
		}

	}
}

func (c *CoinToQuanta) dispatchIssuance() {
	ready := true
	for {
		select {
		case <-c.readyChan:
			ready = true
		case <-time.After(time.Second):
			if ready {
				_, err := c.quantaChannel.GetTopBlockID()
				if err != nil {
					if err.Error() == "connection is shut down" {
						c.quantaChannel.Reconnect()
					} else {
						c.logger.Error("Unhandled error. " + err.Error())
					}
				}
				c.processDeposits()
				c.processSubmissions()
			}

		case <-c.doneChan:
			c.logger.Infof("Exiting.")
			break
		}
	}
}

func (c *CoinToQuanta) StartConsensus(tx *coin.Deposit, consensus ConsensusType) (string, error) {
	txResult := HEX_NULL
	errResult := error(nil)

	c.logger.Infof("Start new round %s %s to=%s amount=%d type =%d", tx.Tx, tx.CoinName, tx.QuantaAddr, tx.Amount, consensus)
	var encoded string
	var err error

	switch consensus {
	case NEWASSET_CONSENSUS:
		encoded, err = c.quantaChannel.CreateNewAssetProposal(c.quantaOptions.Issuer, tx.CoinName, 5)
	case ISSUE_CONSENSUS:
		encoded, err = c.quantaChannel.CreateIssueAssetProposal(tx)
	case TRANSFER_CONSENSUS:
		encoded, err = c.quantaChannel.CreateTransferProposal(tx)

	}

	if err != nil {
		c.logger.Error("Failed to encode refund 1" + err.Error())
		db.ChangeSubmitState(c.rDb, tx.Tx, db.ENCODE_FAILURE, db.DEPOSIT, tx.BlockHash)
		return HEX_NULL, err
	}

	data, err := json.Marshal(&coin.EncodedMsg{encoded, tx.Tx, tx.BlockID, tx.CoinName})

	if err != nil {
		c.logger.Error("Failed to encode refund 2" + err.Error())
		return HEX_NULL, err
	}

	// wait for other node to see the tx
	c.cosi.StartNewRound(common.Bytes2Hex(data))

	result := <-c.cosi.FinalSigChan

	if result.Msg == nil {
		errResult = errors.New("Unable to sign for refund")
		c.logger.Error("Unable to sign for refund")
	} else {
		// save this to queue for later in case ETH RPC is down.
		tx.Signatures = result.Msg
		// save in eth_tx_log_signed (kvstore) [S=signed,X=submitted,F=failed(uncoverable), R=retry(connection failed)] ; recoverable=RPC not available
		c.logger.Infof("Great! Cosi successfully signed deposit")

		txe, err := quanta.ProcessGrapheneTransaction(encoded, tx.Signatures)

		if consensus == NEWASSET_CONSENSUS {
			_, err = c.quantaChannel.Broadcast(txe)
			if err != nil {
				return HEX_NULL, err
			}
			fmt.Println("Asset Created")
			return HEX_NULL, nil
		} else {
			db.ChangeSubmitQueue(c.rDb, tx.Tx, txe, db.DEPOSIT, tx.BlockHash)
		}

		txResult = ""
		errResult = err

		if c.SuccessCb != nil {
			c.SuccessCb(DepositResult{tx, errResult, txResult})
		}
	}
	return txResult, errResult
}

func (c *CoinToQuanta) Stop() {
	c.doneChan <- true
}
