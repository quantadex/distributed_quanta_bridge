package main

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/quanta_book/consensus/cosi"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/go-errors/errors"
	"fmt"
)

type AddressConsensus struct {
	logger        logger.Logger
	trustPeer	  *peer_contact.TrustPeerNode
	cosi		  *cosi.Cosi
}

type AddressChange struct {
	QuantaAddr string
	Address string
}

func NewAddressConsensus(logger logger.Logger, trustNode *TrustNode, db *db.DB,kv kv_store.KVStore, minBlock int64) *AddressConsensus {
	var res AddressConsensus
	res.trustPeer = peer_contact.NewTrustPeerNode(trustNode.man, trustNode.peer, trustNode.nodeID, trustNode.queue, queue.ADDR_MSG_QUEUE, "/node/api/address")
	res.cosi = cosi.NewProtocol(res.trustPeer, trustNode.nodeID == 0, time.Second*3)
	res.logger = logger
	res.cosi.Verify = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressChange{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		headBlock, _ := control.GetLastBlock(kv, coin.BLOCKCHAIN_ETH)
		addr, err := db.GetAvailableShareAddress(headBlock, minBlock)
		if err != nil {
			return err
		}
		if addr.Address != msg.Address {
			return errors.New(fmt.Sprintf("Address available not match - %s msg=%s",addr.Address,msg.Address))
		}
		return nil
	}

	res.cosi.Persist = func(encoded string) error {
		decoded := common.Hex2Bytes(encoded)
		msg := &AddressChange{}
		err := json.Unmarshal(decoded, msg)
		if err != nil {
			return err
		}

		err = db.UpdateShareAddressDestination(msg.Address, msg.QuantaAddr)
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

func (c *AddressConsensus) GetConsensus(msg AddressChange) error {
	data, err := json.Marshal(&msg)

	if err != nil {
		c.logger.Error("Failed to encode check address:" + err.Error())
		return err
	}

	// wait for other node to see the tx
	c.cosi.StartNewRound(common.Bytes2Hex(data))

	result := <-c.cosi.FinalSigChan

	if result.Msg == nil {
		c.logger.Error("Unable to agree for address change")
		return errors.New("Unable to agree for address change")
	} else {
		c.logger.Infof("Cosi successfully agreed on address %s for %s", msg.Address, msg.QuantaAddr)
	}
	return nil
}