package quanta

/**
 Links:
http://docs.bitshares.org/integration/traders/index.html#public-api
https://github.com/scorum/bitshares-go/blob/master/apis/database/api_test.go
*/
import (
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/scorum/bitshares-go/apis/database"
	"github.com/scorum/bitshares-go/apis/login"
	"github.com/scorum/bitshares-go/apis/networkbroadcast"
	"github.com/scorum/bitshares-go/sign"
	"github.com/scorum/bitshares-go/transport/websocket"
	"github.com/scorum/bitshares-go/types"
	"log"
	"math"
	"strconv"
	"time"
)

type QuantaGraphene struct {
	Database         *database.API
	kv               kv_store.KVStore
	Logger           logger.Logger
	Db               *db.DB
	NetworkBroadcast *networkbroadcast.API
}

type Asset struct {
	ID                 types.ObjectID
	Symbol             string
	Precision          uint8
	Issuer             string
	DynamicAssetDataID string
}

type Object struct {
	ID   types.ObjectID
	Name string
}

const url = "ws://testnet-01.quantachain.io:8090"

func (q *QuantaGraphene) Attach() error {
	transport, err := websocket.NewTransport(url)
	if err != nil {
		return err
	}
	databaseAPIID, err := login.NewAPI(transport).Database()
	api := database.NewAPI(databaseAPIID, transport)
	q.Database = api

	networkAPIID, err := login.NewAPI(transport).NetworkBroadcast()
	apiNetwork := networkbroadcast.NewAPI(networkAPIID, transport)
	q.NetworkBroadcast = apiNetwork
	return nil
}

func (q *QuantaGraphene) Broadcast(stx string) error {
	// broadcast here
	var err error
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(stx), &tx)
	//q.NetworkBroadcast.BroadcastTransaction(tx.Transaction)

	resp, err := q.NetworkBroadcast.BroadcastTransactionSynchronous(tx.Transaction)
	fmt.Println("response = ", resp)
	return err
}

func (q *QuantaGraphene) AttachQueue(kv kv_store.KVStore) error {
	panic("implement me")
}

// get_dynamics
func (q *QuantaGraphene) GetTopBlockID(accountId string) (int64, error) {
	res, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	blockId := res.HeadBlockNumber

	return int64(blockId), nil
}

// get block , transfer
func (q *QuantaGraphene) GetRefundsInBlock(blockID int64, trustAddress string) ([]Refund, int64, error) {
	var refunds []Refund
	block, err := q.Database.GetBlock(uint32(blockID))
	if err != nil {
		return refunds, 0, err
	}
	var i, j int
	for i = 0; i < len(block.Transactions); i++ {
		for j = 0; j < len(block.Transactions[i].Operations); j++ {
			op := block.Transactions[i].Operations[j]
			if op.Type() == types.TransferOpType {
				op := op.(*types.TransferOperation)

				receiver, err := q.Database.GetObjects(op.To)
				to := &Object{}
				err = json.Unmarshal(receiver[0], &to)
				if err != nil {
					return refunds, 0, err
				}
				if to.Name == trustAddress {
					coin, err := q.Database.GetObjects(op.Amount.AssetID)
					result := &Asset{}
					err = json.Unmarshal(coin[0], &result)
					if err != nil {
						return refunds, 0, err
					}

					sender, err := q.Database.GetObjects(op.From)
					from := &Object{}
					err = json.Unmarshal(sender[0], &from)
					if err != nil {
						return refunds, 0, err
					}

					txid := strconv.Itoa(int(blockID)) + "_" + strconv.Itoa(int(op.Type()))

					newRefund := Refund{
						OperationID:        int64(op.Type()),
						SourceAddress:      from.Name,
						DestinationAddress: to.Name,
						Amount:             op.Amount.Amount,
						CoinName:           result.Symbol,
						TransactionId:      txid,
						PageTokenID:        blockID,
						LedgerID:           int32(blockID),
					}
					refunds = append(refunds, newRefund)
				}
			}
		}
	}
	return refunds, 0, nil
}

func (q *QuantaGraphene) GetBalance(assetName string, quantaAddress string) (float64, error) {
	id, err := q.Database.LookupAssetSymbols(assetName)
	var balance []*types.AssetAmount
	balance, err = q.Database.GetNamedAccountBalances(quantaAddress, id[0].ID)
	if err != nil {
		return 0, err
	}
	precision := math.Pow(10, float64(id[0].Precision))
	return float64(balance[0].Amount) / precision, nil
}

func (q *QuantaGraphene) GetAllBalances(quantaAddress string) (map[string]float64, error) {
	balance, err := q.Database.GetNamedAccountBalances(quantaAddress)
	if err != nil {
		return nil, err
	}
	balances := make(map[string]float64, len(balance))
	var i int
	for i = 0; i < len(balance); i++ {
		balances[string(i)] = float64(balance[i].Amount)
	}
	return balances, nil
}

// https://github.com/scorum/bitshares-go/blob/bbfc9bedaa1b2ddaead3eafe47237efcd9b8496d/client.go
func (q *QuantaGraphene) CreateProposeTransaction(dep *coin.Deposit) (string, error) {
	var fee types.AssetAmount
	var amount types.AssetAmount

	id, err := q.Database.LookupAssetSymbols(dep.CoinName)
	if err != nil {
		return "", err
	}
	amount.Amount = uint64(dep.Amount)
	amount.AssetID = id[0].ID

	fee.Amount = 0
	fee.AssetID = id[0].ID

	userIdSender, err := q.Database.LookupAccounts(dep.SenderAddr, 1)
	if err != nil {
		return "", err
	}
	userIdReceiver, err := q.Database.LookupAccounts(dep.QuantaAddr, 1)
	if err != nil {
		return "", err
	}

	op := types.NewTransferOperation(userIdSender[dep.SenderAddr], userIdReceiver[dep.QuantaAddr], amount, fee)

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return "", err

	}
	op.Fee.Amount = fees[0].Amount

	return q.PrepareTX(op)
}

func (q *QuantaGraphene) AssetProposeTransaction() (string, error) {
	var d = make([]types.ObjectID, 0)

	var issuer, baseObject, quoteObject types.ObjectID
	issuer.Space = 1
	issuer.Type = 2
	issuer.ID = 23

	baseObject.Space = 1
	baseObject.Type = 3
	baseObject.ID = 0

	quoteObject.Space = 1
	quoteObject.Type = 3
	quoteObject.ID = 1

	var base, quote types.AssetAmount
	base.Amount = 1
	base.AssetID = baseObject

	quote.Amount = 1
	quote.AssetID = quoteObject

	var fee types.AssetAmount
	fee.Amount = 0
	fee.AssetID = base.AssetID
	fmt.Println(base.AssetID, fee.AssetID)

	p := &types.Price{
		Base:  base,
		Quote: quote,
	}
	var u = make([]json.RawMessage, 0)
	s := &types.AssetOptions{
		MaxSupply:            1000000000000000,
		MarketFeePercent:     0,
		MaxMarketFee:         1000,
		IssuerPermissions:    79,
		Flags:                0,
		CoreExchangeRate:     *p,
		WhiteListAuthorities: d,
		BlackListAuthorities: d,
		WhiteListMarkets:     d,
		BlackListMarkets:     d,
		Description:          "My fancy description",
		Extensions:           u,
	}
	w := &types.CreateAsset{
		Issuer:             issuer,
		Symbol:             "TOKENXA",
		Precision:          5,
		CommonOptions:      *s,
		IsPredictionMarket: false,
		Extensions:         u,
	}
	w.Fee.AssetID = baseObject
	fees, _ := q.Database.GetRequiredFee([]types.Operation{w}, fee.AssetID.String())
	//fmt.Println("error = ", err)
	w.Fee.Amount = fees[0].Amount

	//fmt.Println("type = ", w.Type())
	fmt.Println("AssetId = ", w.Fee.AssetID)
	return q.PrepareTX(w)
}

func (q *QuantaGraphene) IssueAssetPropose() (string, error) {
	var e, f types.ObjectID
	e.Space = 1
	e.Type = 2
	e.ID = 16

	f.Space = 1
	f.Type = 3
	f.ID = 0

	var a types.AssetAmount
	a.Amount = 1
	a.AssetID = f

	w := &types.IssueAsset{
		To:        e,
		Amount:    a,
		Memo:      "",
		Broadcast: false,
	}
	return q.PrepareTX(w)

}

func (q *QuantaGraphene) PrepareTX(operations ...types.Operation) (string, error) {
	//fmt.Println("operation = ", operations[0].(*types.CreateAsset).Fee.AssetID)
	props, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return "", err
	}

	block, err := q.Database.GetBlock(props.LastIrreversibleBlockNum)
	if err != nil {
		return "", err
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return "", err
	}

	expiration := props.Time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.LastIrreversibleBlockNum - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}

	data, err := json.Marshal(stx)
	fmt.Println("stx = ", stx.Transaction.Operations[0])
	fmt.Println("data = ", string(data))
	fmt.Println(err)

	return string(data), err
}

func (q *QuantaGraphene) LimitOrderCancel(key string, feePayingAccount, order types.ObjectID, fee types.AssetAmount) (string, error) {
	op := &types.LimitOrderCancelOperation{
		Fee:              fee,
		FeePayingAccount: feePayingAccount,
		Order:            order,
		Extensions:       []json.RawMessage{},
	}

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return "", err
	}
	op.Fee.Amount = fees[0].Amount

	return q.PrepareTX(op)
}

func (q *QuantaGraphene) LimitOrderCreate(key string, seller types.ObjectID, fee, amToSell, minToRecive types.AssetAmount, expiration time.Duration, fillOrKill bool) (string, error) {
	props, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return "", err
	}

	op := &types.LimitOrderCreateOperation{
		Fee:          fee,
		Seller:       seller,
		AmountToSell: amToSell,
		MinToReceive: minToRecive,
		Expiration:   types.NewTime(props.Time.Add(expiration)),
		FillOrKill:   fillOrKill,
		Extensions:   []json.RawMessage{},
	}

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return "", err
	}
	op.Fee.Amount = fees[0].Amount

	return q.PrepareTX(op)
}

func (q *QuantaGraphene) DecodeTransaction(base64 string) (*coin.Deposit, error) {
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(base64), &tx)

	op := tx.Operations[0]
	if op.Type() == types.TransferOpType {
		op := op.(*types.TransferOperation)

		receiver, err := q.Database.GetObjects(op.To)
		to := &Object{}
		err = json.Unmarshal(receiver[0], &to)
		if err != nil {
			return nil, err
		}

		coinName, err := q.Database.GetObjects(op.Amount.AssetID)
		asset := &Asset{}
		err = json.Unmarshal(coinName[0], &asset)
		if err != nil {
			return nil, err
		}

		sender, err := q.Database.GetObjects(op.From)
		from := &Object{}
		err = json.Unmarshal(sender[0], &from)
		if err != nil {
			return nil, err
		}

		return &coin.Deposit{CoinName: asset.Symbol,
			QuantaAddr: to.Name,
			Amount:     int64(op.Amount.Amount),
			BlockID:    0,
		}, nil
	}
	return nil, nil
}

func ProcessGrapheneTransaction(proposed string, sigs []string) (string, error) {
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(proposed), &tx)

	tx.Transaction.Signatures = sigs
	signed, err := json.Marshal(tx)
	return string(signed), err
}

func (q *QuantaGraphene) ProcessDeposit(deposit *coin.Deposit, proposed string) error {
	txe, err := ProcessGrapheneTransaction(proposed, deposit.Signatures)
	println(txe, err)
	return db.ChangeSubmitQueue(q.Db, deposit.Tx, txe, db.DEPOSIT)
}
