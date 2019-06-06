package main

import (
	"crypto/sha256"
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
	"github.com/quantadex/quanta_book/consensus/cosi"
	"time"
)

type AddressConsensus struct {
	logger       logger.Logger
	trustPeer    *peer_contact.TrustPeerNode
	cosi         *cosi.Cosi
	pool         AddressRequestPool
	db           *db.DB
	poolNotify   map[string]chan error
	msgChan      chan MsgAsync
	stateTracker map[string]map[string]int // tracking blockhash -> state -> count
}

type AddressChange struct {
	Blockchain string
	QuantaAddr string
	Address    string
	Counter    uint64
}

type MsgAsync struct {
	data   AddressChange
	notify chan error
}

func GetAddressId(addr AddressChange) string {
	return fmt.Sprintf("%s:%s:%s", addr.Blockchain, addr.QuantaAddr, addr.Address)
}

type AddressBlock struct {
	transactions []AddressChange        // batch multiple address request
	state        []db.CrosschainAddress // state
}

// tx pool
type AddressRequestPool []AddressChange // grab from top

func NewAddressConsensus(logger logger.Logger, trustNode *TrustNode, db *db.DB, kv kv_store.KVStore, minBlock int64) *AddressConsensus {
	var res AddressConsensus
	res.trustPeer = peer_contact.NewTrustPeerNode(trustNode.man, trustNode.peer, trustNode.nodeID, trustNode.queue, queue.ADDR_MSG_QUEUE, "/node/api/address", trustNode.quantakM)
	res.cosi = cosi.NewProtocol(res.trustPeer, trustNode.nodeID == 0, time.Second*6)
	res.logger = logger
	res.db = db
	res.msgChan = make(chan MsgAsync, 1)
	res.poolNotify = make(map[string]chan error)
	res.stateTracker = make(map[string]map[string]int)
	res.cosi.Verify = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressBlock{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}
		decodedSha := sha256.Sum256(decoded)

		blockHash := common.Bytes2Hex(decodedSha[:])
		counter, ok := res.stateTracker[blockHash]
		if !ok {
			counter = map[string]int{}
			res.stateTracker[blockHash] = counter
		}
		data, err := json.Marshal(msg.state)

		// not sure what happens when item does not exist in counter
		counter[string(data)] += 1

		//
		//if msg.Blockchain == coin.BLOCKCHAIN_ETH {
		//	headBlock, _ := control.GetLastBlock(kv, coin.BLOCKCHAIN_ETH)
		//	addrAvailable, err := db.GetAvailableShareAddress(headBlock, minBlock)
		//	if err != nil {
		//		return err
		//	}
		//	for _, a := range addrAvailable {
		//		if a.Address == msg.Address {
		//			return nil
		//		}
		//	}
		//	return errors.New(fmt.Sprintf("Address available not match - msg=%s", msg.Address))
		//
		//} else {
		//	addr, err := trustNode.CreateMultisig(msg.Blockchain, msg.QuantaAddr)
		//	if err != nil {
		//		return errors.New("Unable to generate address for " + msg.Blockchain + "," + err.Error())
		//	}
		//	if addr.ContractAddress != msg.Address || addr.QuantaAddr != msg.QuantaAddr {
		//		return errors.New("Unable to generate address for " + msg.Blockchain + ", somethign don't match")
		//	}
		//	return nil
		//}
		return nil
	}

	res.cosi.Persist = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressChange{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		if msg.Blockchain == "" {
			headBlock, _ := control.GetLastBlock(kv, coin.BLOCKCHAIN_ETH)
			err = db.UpdateShareAddressDestination(msg.Address, msg.QuantaAddr, uint64(headBlock))
		} else {
			addr, err := trustNode.CreateMultisig(msg.Blockchain, msg.QuantaAddr)
			if err != nil {
				return errors.New("Cannot persist due to error : " + err.Error())
			}
			err = db.AddCrosschainAddress(addr)
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
				res.cosi.CosiMsgChan <- msg
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return &res
}

// logic to submit
// respond back to api
func (c *AddressConsensus) GetAddress(msg AddressChange) error {
	done := make(chan error, 1)

	done <- nil
	c.msgChan <- MsgAsync{msg, done}
	go c.StartConsensusIfNeeded()
	err := <-done
	return err
}

func (c *AddressConsensus) StartConsensusIfNeeded() error {
	// gather enough txs
	pendingConsensus := false
	doneConsensus := make(chan error, 1)
	pendingTxs := []AddressChange{}

	for {
		select {
		case msg := <-c.msgChan:
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
				pendingConsensus = true
			}
		}
	}
}

func (c *AddressConsensus) startNewBlock(txsToProcess []AddressChange, done chan<- error) {

	state := c.db.GetCrosschainAll()

	// create the block
	block := AddressBlock{
		transactions: txsToProcess,
		state:        state,
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
