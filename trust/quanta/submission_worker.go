package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/stellar/go/clients/horizon"
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"time"
)

type SubmitWorkerImpl struct {
	horizonClient *horizon.Client
	logger logger.Logger
	queue queue.Queue
	horizonUrl string
}

func (s *SubmitWorkerImpl) Dispatch() {
	if s.logger == nil {
		panic("Missing logger")
	}
	s.logger.Infof("Submitworker started")

	for {
		//println("Wake up")
		time.Sleep(time.Second)

		data, err := s.queue.Get(queue.QUANTA_TX_QUEUE)
		if err != nil {
			continue
		}

		var deposit peer_contact.PeerMessage
		err = json.Unmarshal(data, &deposit)
		if err != nil {
			s.logger.Error("could not unmarshall")
			continue
		}

		s.logger.Infof("Submit TX: %s", deposit.MSG)

		_, err = s.horizonClient.SubmitTransaction(deposit.MSG)
		if err != nil {
			err2 := err.(*horizon.Error)
			s.logger.Error("could not submit transaction " + err2.Error() + err2.Problem.Detail)
		}
	}
}

func (s *SubmitWorkerImpl) AttachQueue(q queue.Queue) error {
	s.queue = q
	s.queue.CreateQueue(queue.QUANTA_TX_QUEUE)

	s.horizonClient = &horizon.Client{
		URL:  s.horizonUrl,
		HTTP: http.DefaultClient,
	}

	return nil
}


