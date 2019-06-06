package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"time"
	"crypto/sha256"
	"github.com/quantadex/distributed_quanta_bridge/common/consensus"
)

type AddressConsensus struct {
	logger    logger.Logger
	trustPeer *peer_contact.TrustPeerNode
	cosi      *consensus.Cosi
	pool      AddressRequestPool
	db *db.DB
	poolNotify map[string]chan error
	msgChan chan MsgAsync
	doneChan chan bool
	stateTracker map[string]map[string]int  // tracking blockhash -> state -> count
}

type AddressChange struct {
	Blockchain string
	QuantaAddr string
	Address    string
	Counter	   uint64
}

type MsgAsync struct {
	data AddressChange
	notify chan error
}

func GetAddressId(addr AddressChange) string {
	return fmt.Sprintf("%s:%s:s", addr.Blockchain, addr.QuantaAddr, addr.Address)
}

type AddressBlock struct {
	Transactions []AddressChange  // batch multiple address request
	State []db.CrosschainAddress  // state
}

// tx pool
type AddressRequestPool []AddressChange  // grab from top

func NewAddressConsensus(logger logger.Logger, trustNode *TrustNode, db *db.DB, kv kv_store.KVStore, minBlock int64) *AddressConsensus {
	var res AddressConsensus
	res.trustPeer = peer_contact.NewTrustPeerNode(trustNode.man, trustNode.peer, trustNode.nodeID, trustNode.queue, queue.ADDR_MSG_QUEUE, "/node/api/address", trustNode.quantakM)
	res.cosi = consensus.NewProtocol(res.trustPeer, trustNode.nodeID == 0, time.Second*6)
	res.logger = logger
	res.db = db
	res.msgChan = make(chan MsgAsync, 1)
	res.poolNotify = map[string]chan error {}
	res.doneChan = make(chan bool, 1)

	res.cosi.Verify = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressBlock{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		// if my state is the same -- confirm.
		my_state := db.GetCrosschainAll()
		my_state_hash, _ := json.Marshal(my_state)
		//println(string(my_state_hash))

		my_state_hashb := sha256.Sum256(my_state_hash)
		my_state_hash2 := common.Bytes2Hex(my_state_hashb[:])

		in_state_hash, _ := json.Marshal(msg.State)
		//println("IN", string(in_state_hash))

		in_state_hashb := sha256.Sum256(in_state_hash)
		in_state_hash2 := common.Bytes2Hex(in_state_hashb[:])

		if my_state_hash2 != in_state_hash2 {
			return errors.New(fmt.Sprintf("AddressConsensus state hash mismatch mine=%s in=%s", my_state_hash2, in_state_hash2))
		}

		for _, tx := range msg.Transactions  {
			if tx.Blockchain == coin.BLOCKCHAIN_ETH {
				headBlock, _ := control.GetLastBlock(kv, coin.BLOCKCHAIN_ETH)
				addrAvailable, err := db.GetAvailableShareAddress(headBlock, minBlock)
				if err != nil {
					return err
				}

				found := false
				for _, a := range addrAvailable {
					if a.Address == tx.Address {
						found = true
					}
				}

				if !found {
					return errors.New(fmt.Sprintf("Address available not match - msg=%s", tx.Address))
				}
			} else {
				addr, err := trustNode.CreateMultisig(tx.Blockchain, tx.QuantaAddr)
				if err != nil {
					return errors.New("Unable to generate address for " + tx.Blockchain + "," + err.Error())
				}
				if addr.ContractAddress != tx.Address || addr.QuantaAddr != tx.QuantaAddr {
					return errors.New("Unable to generate address for " + tx.Blockchain + ", somethign don't match")
				}
			}
		}

		return nil
	}

	res.cosi.Persist = func(encoded string, repair bool) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressBlock{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		if repair {
			logger.Infof("***** REPAIRING ADDRESS TABLE *****")
		}
		logger.Infof("Persisting number of txs=%d", len(msg.Transactions))

		for _, tx := range msg.Transactions {
			if tx.Blockchain == coin.BLOCKCHAIN_ETH {
				headBlock, _ := control.GetLastBlock(kv, coin.BLOCKCHAIN_ETH)
				err = db.UpdateShareAddressDestination(tx.Address, tx.QuantaAddr, uint64(headBlock))
			} else {
				addr, err := trustNode.CreateMultisig(tx.Blockchain, tx.QuantaAddr)
				if err != nil {
					return errors.New("Cannot persist due to error : " + err.Error())
				}
				err = db.AddCrosschainAddress(addr)
			}
		}

		return err
	}

	res.cosi.SignMsg = func(encoded string) (string, error) {
		return encoded, nil
	}

	res.cosi.Start()

	go func() {
		for {
			// feed messages back
			msg := res.trustPeer.GetMsg()

			if msg != nil {
				//log.Infof("Got peer message %v", msg.Signed)
				res.cosi.CosiMsgChan <- &consensus.CosiMessage{ msg.Msg, msg.Signed, consensus.Phase(msg.Phase), msg.Initial }
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	go res.StartConsensusIfNeeded()

	return &res
}

// logic to submit
// respond back to api
func (c *AddressConsensus) GetAddress(msg AddressChange) error {
	done := make(chan error, 1)

	c.msgChan <- MsgAsync{ msg, done}
	err := <- done
	return err
}

func (c *AddressConsensus) StartConsensusIfNeeded() error {
	c.logger.Infof("Started AddressBlock Block Producer")

	// gather enough txs
	pendingConsensus := false
	doneConsensus := make(chan error, 1)
	pendingTxs := []AddressChange{}
	doneFlag := false

	for {
		select {
			case msg := <-c.msgChan:
				c.logger.Infof("Enqueue new address %s %s", msg.data.Blockchain, msg.data.QuantaAddr)
				c.pool = append(c.pool, msg.data) // modify
				c.poolNotify[GetAddressId(msg.data)] = msg.notify

			// start new block every 3 sec for at least 1 address change
			case <-time.After(time.Second * 3):
				if !pendingConsensus && len(c.pool) > 0 {
					txsToProcess := c.pool        // read
					c.pool = AddressRequestPool{} // modify
					pendingConsensus = true
					pendingTxs = txsToProcess

					c.startNewBlock(txsToProcess, doneConsensus)
				}

			// notify all the callers
			case err := <-doneConsensus:
				for _, tx := range pendingTxs {
					notify := c.poolNotify[GetAddressId(tx)]
					notify <- err
					delete(c.poolNotify, GetAddressId(tx))
					pendingTxs = []AddressChange{}
					pendingConsensus = false

					c.logger.Infof("Notify addressBlock done %v", err)
				}
			case <-c.doneChan:
				doneFlag = true
				break
		}

		if doneFlag {
			c.logger.Infof("Exiting address consensus.")
			break
		}
	}
	return nil
}

func (c *AddressConsensus) Stop() {
	c.doneChan <- true
}

func (c *AddressConsensus) startNewBlock(txsToProcess []AddressChange, done chan error) {
	c.logger.Infof("Generate new address block. txs=%d", len(txsToProcess))

	state := c.db.GetCrosschainAll()

	// create the block
	block := AddressBlock{
		Transactions: txsToProcess,
		State: state,
	}

	data, err := json.Marshal(&block)

	if err != nil {
		c.logger.Error("Failed to encode check address:" + err.Error())
	} else {
		// wait for other node to see the tx
		c.cosi.StartNewRound(common.Bytes2Hex(data))

		result := <-c.cosi.FinalSigChan

		if result.Msg == nil {
			err = errors.New("Unable to agree for address change")
			c.logger.Error(err.Error())
		} else {
			c.logger.Infof("Cosi successfully agreed on address")
		}
	}

	done <- err
}
