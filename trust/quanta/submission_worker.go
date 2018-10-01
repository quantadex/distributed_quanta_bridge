package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/stellar/go/clients/horizon"
	"github.com/spf13/viper"
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"encoding/json"
	"time"
	"fmt"
)

type SubmitWorkerImpl struct {
	horizonClient *horizon.Client
	queue queue.Queue
}

func (s *SubmitWorkerImpl) Dispatch() {
	println("Submitworker started")
	for {
		data, err := s.queue.Get(queue.QUANTA_TX_QUEUE)
		if err != nil{
			continue
		}

		var deposit peer_contact.PeerMessage
		err = json.Unmarshal(data, &deposit)
		if err != nil {
			println("could not unmarshall")
			continue
		}

		println(deposit.MSG)
		_, err = s.horizonClient.SubmitTransaction(deposit.MSG)
		if err != nil {
			err2 := err.(*horizon.Error)
			fmt.Println("could not submit transaction ", err2.Error(), err2.Problem)
		}
		time.Sleep(time.Second)
	}
}

func (s *SubmitWorkerImpl) AttachQueue(q queue.Queue) error {
	s.queue = q
	s.queue.CreateQueue(queue.QUANTA_TX_QUEUE)

	s.horizonClient = &horizon.Client{
		URL:  viper.GetString("HORIZON_URL"),
		HTTP: http.DefaultClient,
	}

	return nil
}


