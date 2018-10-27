package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/stellar/go/amount"
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
	AssetCode   string    `json:"asset_code"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	TxHash      string    `json:"transaction_hash"`
	Amount      string    `json:"amount"`
}

type Balances struct {
	Links struct {
		Self horizon.Link `json:"self"`
		Next horizon.Link `json:"next"`
		Prev horizon.Link `json:"prev"`
	} `json:"_links"`
	Records []Balance `json:"balances"`
}

type Balance struct {
	Balance     string `json:"balance"`
	Limit       string `json:"limit"`
	AssetType   string `json:"asset_type"`
	AssetCode   string `json:"asset_code"`
	AssetIssuer string `json:"asset_issuer"`
}

// remember to test coins < 10^7
func (q *QuantaClient) CreateProposeTransaction(deposit *coin.Deposit) (string, error) {
	amount := fmt.Sprintf("%.7f", float64(deposit.Amount)/10000000)
	println("Propose TX: ", deposit.CoinName, q.Issuer, amount, deposit.QuantaAddr)

	tx, err := b.Transaction(
		b.Network{q.Network},
		b.SourceAccount{q.Issuer},
		b.AutoSequence{q.horizonClient},
		//b.Sequence{ 0 },
		b.Payment(
			b.Destination{deposit.QuantaAddr},
			b.CreditAmount{deposit.CoinName, q.Issuer, amount},
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
	q.kv = kv
	q.kv.CreateTable(kv_store.PENDING_QUANTA_TX)
	q.kv.CreateTable(kv_store.COMPLETED_QUANTA_TX)
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

func (q *QuantaClient) GetBalance(assetName string, quantaAddress string) (float64, error) {
	m, err := q.GetAllBalances(quantaAddress)
	for k, v := range m {
		if k == assetName {
			return v, nil
		}
	}
	return 0, err
}

func (q *QuantaClient) GetAllBalances(quantaAddress string) (map[string]float64, error) {
	var m map[string]float64
	m = make(map[string]float64)
	url := fmt.Sprintf("%s/accounts/%s", q.horizonClient.URL, quantaAddress)
	resp, err := q.horizonClient.HTTP.Get(url)
	if err != nil {
		return m, err
	}

	var balances Balances
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return m, errors.New("failed to decode operations: " + err.Error())
	}
	for i := 0; i < len(balances.Records); i++ {
		if balances.Records[i].AssetCode == "" {
			m["native"], _ = strconv.ParseFloat(balances.Records[i].Balance, 64)
		} else {
			m[balances.Records[i].AssetCode], _ = strconv.ParseFloat(balances.Records[i].Balance, 64)
		}
	}
	return m, err
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
					am, err := amount.ParseInt64(op.Amount)
					if err != nil {
						return nil, cursor, nil
					}
					//TODO: Add ledger
					// it's a refund
					newRefund := Refund{
						CoinName:           op.AssetCode,
						DestinationAddress: op.To,
						OperationID:        num,
						Amount:             uint64(am),
						PageTokenID:        pt,
						TransactionId:      op.TxHash,
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

func (q *QuantaClient) ProcessDeposit(deposit peer_contact.PeerMessage) error {
	//x, _ := q.GetBalance("poojakishoreshah", "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN")
	//x, _ := q.GetTopBlockID("QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN")
	//fmt.Println(x)
	data, err := json.Marshal(deposit)
	if err != nil {
		return err
	}
	key := strconv.Itoa(int(deposit.Proposal.BlockID)) + deposit.Proposal.CoinName + deposit.Proposal.QuantaAdress
	return q.kv.SetValue(kv_store.PENDING_QUANTA_TX, key, "", string(data))
}
