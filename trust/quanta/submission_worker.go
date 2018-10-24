package quanta

import (
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/stellar/go/clients/horizon"
	"net/http"
	"time"
)

//TODO: add kvstore, keep 2 buckets 1) pending_quanta_tx  2) completed_quanta_tx
type SubmitWorkerImpl struct {
	horizonClient *horizon.Client
	logger        logger.Logger
	queue         queue.Queue
	horizonUrl    string
	kv            kv_store.KVStore
}

func (s *SubmitWorkerImpl) Dispatch() {
	if s.logger == nil {
		panic("Missing logger")
	}
	s.logger.Infof("Submitworker started")

	for {
		//println("Wake up")
		time.Sleep(time.Second)
		data, err := s.kv.GetAllValues(kv_store.PENDING_QUANTA_TX)
		if err != nil {
			continue
		}

		var deposit peer_contact.PeerMessage
		for k, v := range data {

			err = json.Unmarshal([]byte(v), &deposit)
			if err != nil {
				s.logger.Error("could not unmarshall")
				continue
			}
			s.logger.Infof("Submit TX: %s signed=%v %s", deposit.Proposal.CoinName, deposit.SignedBy, deposit.MSG)

			res, err := s.horizonClient.SubmitTransaction(deposit.MSG)
			if err != nil {
				err2 := err.(*horizon.Error)

				s.logger.Error("could not submit transaction " + res.Hash + " " + err2.Error() + err2.Problem.Detail)
				fmt.Printf("%v", err2.Problem.Extras)

				//xdr.TransactionResult{}

			} else {
				s.logger.Infof("Successful tx submission %s", res.Hash)
				s.kv.RemoveKey(kv_store.PENDING_QUANTA_TX, k)
				s.kv.SetValue(kv_store.COMPLETED_QUANTA_TX, k, "", v)
			}
		}

	}
}

func (s *SubmitWorkerImpl) AttachQueue(kv kv_store.KVStore) error {
	s.kv = kv
	//s.queue.CreateQueue(queue.QUANTA_TX_QUEUE)

	s.horizonClient = &horizon.Client{
		URL:  s.horizonUrl,
		HTTP: http.DefaultClient,
	}

	return nil
}
