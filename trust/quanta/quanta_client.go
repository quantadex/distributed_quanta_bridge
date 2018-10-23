package quanta

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
	"net/http"
	"strconv"
	"time"
)

type QuantaClientOptions struct {
	Logger     logger.Logger
	Network    string
	Issuer     string // pub key
	HorizonUrl string
}

type QuantaClient struct {
	QuantaClientOptions
	horizonClient *horizon.Client
	queue         queue.Queue
	worker        SubmitWorker
	kv            kv_store.KVStore
}

type Operations struct {
	Links struct {
		Self horizon.Link `json:"self"`
		Next horizon.Link `json:"next"`
		Prev horizon.Link `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Records []Operation `json:"records"`
	} `json:"_embedded"`
}

type Operation struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	PagingToken string    `json:"paging_token"`
	CreatedAt   time.Time `json:"created_at"`
	AssetType   string    `json:"asset_type"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	TxHash      string    `json:"transaction_hash"`
}

// remember to test coins < 10^7
func (q *QuantaClient) CreateProposeTransaction(deposit *coin.Deposit) (string, error) {
	amount := fmt.Sprintf("%.7f", float64(deposit.Amount)/10000000)
	println("Propose TX: ", deposit.CoinName, q.Issuer)

	tx, err := b.Transaction(
		b.Network{q.Network},
		b.SourceAccount{q.Issuer},
		b.AutoSequence{q.horizonClient},
		//b.Sequence{ 0 },
		b.Payment(
			b.Destination{deposit.QuantaAddr},
			b.CreditAmount{"mnbvcxzlkjhgfdsapoiuytrewqmnbvcxzlkjhgfdsapoiuytrewq123456789", q.Issuer, amount},
		),
	)

	if err != nil {
		return "", err
	}

	txe, err := tx.Sign()

	if err != nil {
		return "", err
	}

	return txe.Base64()
}

func (k *QuantaClient) DecodeTransaction(base64 string) (*coin.Deposit, error) {
	txe := &xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(base64, txe)
	if err != nil {
		return nil, err
	}

	ops := txe.Tx.Operations
	if len(ops) != 1 {
		return nil, errors.New("no operations found")
	}

	paymentOp, success := ops[0].Body.GetPaymentOp()
	if !success {
		return nil, errors.New("no payment op found")
	}

	return &coin.Deposit{CoinName: paymentOp.Asset.String(),
		QuantaAddr: paymentOp.Destination.Address(),
		Amount:     int64(paymentOp.Amount),
		BlockID:    0,
	}, nil
}

func (q *QuantaClient) Attach() error {
	q.horizonClient = &horizon.Client{
		URL:  q.HorizonUrl,
		HTTP: http.DefaultClient,
	}

	return nil
}

func (q *QuantaClient) AttachQueue(kv kv_store.KVStore) error {

	//q.queue = queueIn
	//q.queue.CreateQueue(queue.QUANTA_TX_QUEUE)
	q.kv = kv
	q.kv.CreateTable("Pending_Quanta_Tx")
	q.kv.CreateTable("Completed_Quanta_Tx")
	q.worker = NewSubmitWorker(q.HorizonUrl, q.Logger)
	q.worker.AttachQueue(q.kv)
	go q.worker.Dispatch()

	return nil
}

func (q *QuantaClient) GetTopBlockID(accountId string) (int64, error) {
	url := fmt.Sprintf("%s/accounts/%s/operations?order=desc&limit=1", q.horizonClient.URL, accountId)

	resp, err := q.horizonClient.HTTP.Get(url)
	if err != nil {
		return 0, err
	}

	var operations Operations
	if err := json.NewDecoder(resp.Body).Decode(&operations); err != nil {
		return 0, errors.New("failed to decode operations: " + err.Error())
	}

	if len(operations.Embedded.Records) > 0 {
		num, _ := strconv.ParseInt(operations.Embedded.Records[0].ID, 10, 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	}

	return 0, nil
}

func (q *QuantaClient) GetTransactionWithHash(hash string) (*horizon.Transaction, error) {
	url := fmt.Sprintf("%s/transactions/%s", q.horizonClient.URL, hash)
	//println(url)

	resp, err := q.horizonClient.HTTP.Get(url)
	if err != nil {
		return nil, err
	}

	var tx horizon.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, errors.New("failed to decode operations: " + err.Error())
	}

	return &tx, nil
}

// returns nextPageToken
func (q *QuantaClient) GetRefundsInBlock(cursor int64, trustAddress string) ([]Refund, int64, error) {
	url := fmt.Sprintf("%s/accounts/%s/payments?order=asc&limit=100&cursor=%d", q.horizonClient.URL, trustAddress, cursor)
	println(url)

	resp, err := q.horizonClient.HTTP.Get(url)
	if err != nil {
		return nil, 0, err
	}

	var operations Operations
	if err := json.NewDecoder(resp.Body).Decode(&operations); err != nil {
		return nil, 0, errors.New("failed to decode operations: " + err.Error())
	}

	refunds := []Refund{}

	if len(operations.Embedded.Records) > 0 {
		var num int64
		var pt int64
		for index, op := range operations.Embedded.Records {
			num, _ = strconv.ParseInt(operations.Embedded.Records[index].ID, 10, 64)
			pt, _ = strconv.ParseInt(operations.Embedded.Records[index].PagingToken, 10, 64)
			if err != nil {
				return nil, 0, err
			}
			if op.Type == "payment" {
				if op.To == trustAddress {
					// it's a refund
					newRefund := Refund{
						CoinName:           op.AssetType,
						DestinationAddress: op.To,
						OperationID:        num,
					}

					tx, err := q.GetTransactionWithHash(op.TxHash)
					if err != nil {
						return nil, cursor, nil
					}
					newRefund.TransactionId = tx.ID
					memo, _ := base64.StdEncoding.DecodeString(tx.Memo)
					newRefund.DestinationAddress = string(memo)
					newRefund.LedgerID = tx.Ledger
					refunds = append(refunds, newRefund)
				}
			}
		}
		return refunds, pt, nil
	}

	return refunds, cursor, nil
}

//TODO: write into pending_quanta_tx  with unique key -> jsonBytes of PeerMessage
// submission worker will scan pending_quanta_tx and submit to the QUANTA blockchain
func (q *QuantaClient) ProcessDeposit(deposit peer_contact.PeerMessage) error {
	data, err := json.Marshal(deposit)
	if err != nil {
		return err
	}
	key := peer_contact.CreateUniqueKey(data, deposit)

	return q.kv.Put("Pending_Quanta_Tx", []byte(key), data)
}
