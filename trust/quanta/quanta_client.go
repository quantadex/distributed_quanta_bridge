package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"net/http"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"encoding/json"
)

type QuantaClient struct {
	network string
	issuer string // pub key
	horizonClient *horizon.Client
	queue queue.Queue
	worker SubmitWorker
}

// remember to test coins < 10^7
func (q *QuantaClient) CreateProposeTransaction(deposit *coin.Deposit) (string, error) {
	amount := fmt.Sprintf("%.7f",float64(deposit.Amount)/10000000)
	println(amount)

	tx, err := b.Transaction(
		b.Network{q.network},
		b.SourceAccount{q.issuer},
		b.AutoSequence{q.horizonClient},
		//b.Sequence{ 0 },
		b.Payment(
			b.Destination{deposit.QuantaAddr},
			b.NativeAmount{ amount},
		),
	)
	println("Seq", tx.TX.SeqNum)

	if err != nil {
		return "", err
	}

	txe, err := tx.Sign()

	if err != nil {
		return "", err
	}

	return txe.Base64()
}

func (q *QuantaClient) Attach() error {
	q.network = viper.GetString("NETWORK_PASSPHRASE")
	viper.GetString("NETWORK_PASSPHRASE")
	q.issuer = viper.GetString("ISSUER_ADDRESS")
	q.horizonClient = &horizon.Client{
		URL:  viper.GetString("HORIZON_URL"),
		HTTP: http.DefaultClient,
	}

	return nil
}

func (q *QuantaClient) AttachQueue(queueIn queue.Queue) error {
	q.queue = queueIn
	q.queue.CreateQueue(queue.QUANTA_TX_QUEUE)

	q.worker = NewSubmitWorker()
	q.worker.AttachQueue(q.queue)
	go q.worker.Dispatch()

	return nil
}

func (q *QuantaClient) GetTopBlockID() (int, error) {
	return 0, nil
}

func (q *QuantaClient) GetRefundsInBlock(blockID int, trustAddress string) ([]Refund, error) {
	return nil, nil
}

func (q *QuantaClient) ProcessDeposit(deposit peer_contact.PeerMessage) error {
	data,err := json.Marshal(deposit)
	if err != nil {
		return err
	}
	return q.queue.Put(queue.QUANTA_TX_QUEUE, data)
}




