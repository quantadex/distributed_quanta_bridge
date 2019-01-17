package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
	"strings"
	"time"
	"github.com/scorum/bitshares-go/types"
	"github.com/scorum/bitshares-go/apis/database"
	"math/big"
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
	coinChannel   coin.Coin     // ethereum
	quantaChannel quanta.Quanta // stellar -> graphene
	db            kv_store.KVStore
	rDb           *db.DB
	man           *manifest.Manifest
	peer          peer_contact.PeerContact
	coinName      string
	trustAddress  common.Address
	trustPeer     *peer_contact.TrustPeerNode
	cosi          *cosi.Cosi
	coinInfo	  *database.Asset

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
	coinName string,
	nodeID int,
	peer peer_contact.PeerContact,
	queue_ queue.Queue,
	options C2QOptions,

	quantaOptions quanta.QuantaClientOptions) *CoinToQuanta {
	res := &CoinToQuanta{C2QOptions: options}
	res.logger = log
	res.coinChannel = c
	res.quantaChannel = q
	res.db = db_
	res.rDb = rDb
	res.man = man
	res.nodeID = nodeID
	res.coinName = coinName
	res.peer = peer
	res.trustAddress = common.HexToAddress(options.EthTrustAddress)
	res.doneChan = make(chan bool, 1)
	res.readyChan = make(chan bool, 1)
	res.quantaOptions = quantaOptions
	res.coinInfo,_ = q.GetAsset(coinName)

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
		if tx != nil {
			// we're not going to sign again
			if tx.Signed == false {
				return nil
			}
			log.Error("Unable to verify refund " + tx.Tx)
		}

		return errors.New("Unable to verify: " + deposit.Tx + " " + err.Error())
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

/**
 * getNewCoinBlockIDs
 *
 * Returns a list of new blocks added to the coin block chain.
 */
func (c *CoinToQuanta) GetNewCoinBlockIDs() []int64 {
	lastProcessed, valid := GetLastBlock(c.db, c.coinName)
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
	for i := common2.MaxInt64(c.BlockStartID, lastProcessed+1); i <= currentTop; i++ {
		blocks = append(blocks, i)
		if len(blocks) == MAX_PROCESS_BLOCKS {
			break
		}
	}
	c.logger.Info(fmt.Sprintf("Coin2Quanta: Got blocks %v", blocks))

	return blocks
}

/**
 * getDepositsInBlock
 *
 * Returns deposits made to the coin trust account in this block
 */

func (c *CoinToQuanta) getDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := db.GetCrosschainByBlockchain(c.rDb, c.coinName)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[strings.ToLower(w.Address)] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	if err != nil {
		c.logger.Info("getDepositsInBlock failed " + err.Error())
		return nil, err
	}

	for _, dep := range deposits {
		if dep.CoinName == "ETH" {
			dep.CoinName = c.coinName

			// ethereum converts to precision 5, now we need to convert to precision of the asset
			dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 5, int(c.coinInfo.Precision))
		} else {
			// we assume precision is always 5
		}

		// Need to convert to uppercase, which graphene requires
		dep.CoinName = strings.ToUpper(dep.CoinName)
	}

	return deposits, nil
}

func (c *CoinToQuanta) processDeposits() {
	txs := db.QueryDepositByAge(c.rDb, time.Now().Add(-time.Second*5), []string{db.SUBMIT_CONSENSUS})
	if len(txs) > 0 {
		tx := txs[0]
		w := &coin.Deposit{
			Tx:         tx.Tx,
			CoinName:   tx.Coin,
			QuantaAddr: tx.To,
			BlockID:    tx.BlockId,
			SenderAddr: tx.From,
			Amount:     tx.Amount,
		}
		// if not a native token, we need to flush it
		if tx.Coin != c.coinName {
			parts := strings.Split(c.coinName, "0X")
			if len(parts) > 1 {
				// flush
				// contract := parts[1]
			}
		}

		// check if asset exists
		//if not, then propose new asset
		exist, err := c.quantaChannel.AssetExist(c.quantaOptions.Issuer, tx.Coin)
		if err != nil {
			c.logger.Error(err.Error())
		}

		if !exist {
			_, err = c.StartConsensus(w, NEWASSET_CONSENSUS)
			if err != nil {
				c.logger.Error("failed to create asset, error = " + err.Error())
			}
		} else{
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

		err := c.quantaChannel.Broadcast(v.SubmitTx)
		if err != nil {
			db.ChangeSubmitState(c.rDb, v.Tx, db.SUBMIT_FATAL, db.DEPOSIT)
			msg := quanta.ErrorString(err, false)
			c.logger.Error("could not submit transaction " + msg)
			if strings.Contains(msg, "tx_bad_seq") || strings.Contains(msg, "op_malformed") {
				db.ChangeSubmitState(c.rDb, v.Tx, db.SUBMIT_FATAL, db.DEPOSIT)
			}
		} else {
			c.logger.Infof("Successful tx submission %s,remove %s", "", k)
			err = db.ChangeSubmitState(c.rDb, v.Tx, db.SUBMIT_SUCCESS, db.DEPOSIT)
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

	c.logger.Infof("Start new round %s %s to=%s amount=%d", tx.Tx, tx.CoinName, tx.QuantaAddr, tx.Amount)
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
		db.ChangeSubmitState(c.rDb, tx.Tx, db.ENCODE_FAILURE, db.DEPOSIT)
		return HEX_NULL, err
	}

	data, err := json.Marshal(&coin.EncodedMsg{encoded, tx.Tx, tx.BlockID})

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
			err = c.quantaChannel.Broadcast(txe)
			if err != nil {
				return HEX_NULL, err
			}
			fmt.Println("Asset Created")
			return HEX_NULL, nil
		} else {
			db.ChangeSubmitQueue(c.rDb, tx.Tx, txe, db.DEPOSIT)
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

/**
 * DoLoop
 *
 * Do one iteration of the loop. Get all new coin blocks and theie deposits.
 * Shoot those into round robin
 * Get all ready messages from RR and send these to quanta.
 *
 * returns allDeposits []*coin.Deposit
 */
func (c *CoinToQuanta) DoLoop(blockIDs []int64) []*coin.Deposit {
	c.logger.Debugf("***** Start # of blocks=%d man.N=%d,man.Q=%d *** ", len(blockIDs), c.man.N, c.man.Q)

	allDeposits := make([]*coin.Deposit, 0)

	if blockIDs != nil {
		for _, blockID := range blockIDs {
			deposits, err := c.getDepositsInBlock(blockID)
			if err != nil {
				c.logger.Error("Failed to get deposits from block: " + err.Error())
				return allDeposits
			}

			if deposits != nil {
				if len(deposits) > 0 {
					c.logger.Info(fmt.Sprintf("Block %d Got deposits %d %v", blockID, len(deposits), deposits))
				}

				for _, dep := range deposits {
					err = db.ConfirmDeposit(c.rDb, dep, false)
					if err != nil {
						c.logger.Error("Cannot insert into db:" + err.Error())
					}
					allDeposits = append(allDeposits, dep)

					if !c.quantaChannel.AccountExist(dep.QuantaAddr) {
						// if not exist, let's bounce money back
					} else if dep.Amount == 0 {
						c.logger.Error("Amount is too small")
					} else if c.nodeID == 0 {
						db.ChangeSubmitState(c.rDb, dep.Tx, db.SUBMIT_CONSENSUS, db.DEPOSIT)
					}
				}
			}

			addresses, err := c.coinChannel.GetForwardersInBlock(blockID)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}

			for _, addr := range addresses {
				if addr.Trust.Hex() == c.trustAddress.Hex() {
					c.logger.Infof("New Forwarder Address ETH->QUANTA address, %s -> %s", addr.ContractAddress.Hex(), addr.QuantaAddr)
					db.AddCrosschainAddress(c.rDb, addr)
				} else {
					c.logger.Error(fmt.Sprintf("MISMATCH: Forwarder address[%s] in blockID=%d does not match our trustAddress[%s]",
						addr.Trust.Hex(), blockID, c.trustAddress.Hex()))
				}
			}
		}
	}

	if len(blockIDs) > 0 {
		lastBlockId := blockIDs[len(blockIDs)-1]
		c.logger.Debugf("set last block coin=%s height=%d", c.coinName, lastBlockId)
		setLastBlock(c.db, c.coinName, lastBlockId)
	}

	return allDeposits
}
