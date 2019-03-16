package quanta

/**
 Links:
http://docs.bitshares.org/integration/traders/index.html#public-api
https://github.com/scorum/bitshares-go/blob/master/apis/database/api_test.go
*/
import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-errors/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/scorum/bitshares-go/apis/database"
	"github.com/scorum/bitshares-go/apis/login"
	"github.com/scorum/bitshares-go/apis/networkbroadcast"
	"github.com/scorum/bitshares-go/sign"
	"github.com/scorum/bitshares-go/transport/websocket"
	"github.com/scorum/bitshares-go/types"
	"gopkg.in/matryer/try.v1"
	"math"
	"strconv"
	"time"
)

type QuantaGraphene struct {
	QuantaClientOptions
	Database         *database.API
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

func (q *QuantaGraphene) Attach() error {
	transport, err := websocket.NewTransport(q.QuantaClientOptions.NetworkUrl)
	if err != nil {
		return err
	}
	databaseAPIID, err := login.NewAPI(transport).Database()
	if err != nil {
		return err
	}
	api := database.NewAPI(databaseAPIID, transport)
	q.Database = api

	networkAPIID, err := login.NewAPI(transport).NetworkBroadcast()
	if err != nil {
		return err
	}
	apiNetwork := networkbroadcast.NewAPI(networkAPIID, transport)
	q.NetworkBroadcast = apiNetwork
	return nil
}

func (q *QuantaGraphene) Reconnect() {
	try.Do(func(attempt int) (bool, error) {
		var err error
		err = q.Attach()
		if err != nil {
			time.Sleep(5 * time.Second)
		} else {
			q.Logger.Error(err.Error())
		}
		return true, err
	})
}

func (q *QuantaGraphene) AssetExist(issuer string, symbol string) (bool, error) {
	asset, err := q.Database.LookupAssetSymbols(symbol)
	if err != nil {
		return false, err
	}
	if len(asset) == 0 || asset[0] == nil {
		return false, nil
	}

	issuerId, err := q.LookupAccount(issuer)
	if err != nil {
		return false, err
	}
	for i := range asset {
		if types.MustParseObjectID(asset[i].Issuer) == issuerId {
			return true, nil
		}
	}
	return false, errors.New("issuer do not match")
}

func (q *QuantaGraphene) Broadcast(stx string) (*networkbroadcast.BroadcastResponse, error) {
	// broadcast here
	var err error
	var tx sign.SignedTransaction
	json.Unmarshal([]byte(stx), &tx)

	resp, err := q.NetworkBroadcast.BroadcastTransactionSynchronous(tx.Transaction)
	return resp, err
}

func (q *QuantaGraphene) AttachQueue(kv kv_store.KVStore) error {
	return nil
}

// get_dynamics
func (q *QuantaGraphene) GetTopBlockID() (int64, error) {
	res, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	blockId := res.HeadBlockNumber

	return int64(blockId - 1), nil
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
				if err != nil {
					return refunds, 0, err
				}
				if receiver == nil {
					return refunds, 0, errors.New("receiver not found")
				}

				to := &Object{}
				err = json.Unmarshal(receiver[0], &to)
				if err != nil {
					return refunds, 0, err
				}

				if to.Name == trustAddress {
					coin, err := q.Database.GetObjects(op.Amount.AssetID)
					if err != nil {
						return refunds, 0, err
					}
					if coin == nil {
						return refunds, 0, errors.New("coin not found")
					}

					coinName := &Asset{}
					err = json.Unmarshal(coin[0], &coinName)
					if err != nil {
						return refunds, 0, err
					}

					sender, err := q.Database.GetObjects(op.From)
					if err != nil {
						return refunds, 0, err
					}
					if sender == nil {
						return refunds, 0, errors.New("sender not found")
					}

					from := &Object{}
					err = json.Unmarshal(sender[0], &from)
					if err != nil {
						return refunds, 0, err
					}

					txId := strconv.Itoa(int(blockID)) + "_" + strconv.Itoa(int(op.Type()))

					var str string
					if op.Memo != nil {
						op.Memo.Message = "0x" + op.Memo.Message
						byt, err := hexutil.Decode(op.Memo.Message)
						if err != nil {
							return nil, 0, err
						}
						str = string(byt)
					}

					newRefund := Refund{
						OperationID:        int64(op.Type()),
						SourceAddress:      from.Name,
						DestinationAddress: str,
						Amount:             op.Amount.Amount,
						CoinName:           coinName.Symbol,
						TransactionId:      txId,
						PageTokenID:        blockID,
						LedgerID:           int32(blockID),
					}
					refunds = append(refunds, newRefund)
				}
			}
		}
	}
	return refunds, blockID, nil
}

func (q *QuantaGraphene) GetBalance(assetName string, quantaAddress string) (float64, error) {
	id, err := q.GetAsset(assetName)
	if err != nil {
		return 0, err
	}

	balance, err := q.Database.GetNamedAccountBalances(quantaAddress, id.ID)
	if err != nil {
		return 0, err
	}
	if len(balance) == 0 || balance[0] == nil {
		return 0, errors.New("balance not found")
	}

	precision := math.Pow(10, float64(id.Precision))
	return float64(balance[0].Amount) / precision, nil
}

func (q *QuantaGraphene) GetAllBalances(quantaAddress string) (map[string]float64, error) {
	balance, err := q.Database.GetNamedAccountBalances(quantaAddress)
	if err != nil {
		return nil, err
	}
	if len(balance) == 0 || balance[0] == nil {
		return nil, errors.New("balance not found")
	}

	balances := make(map[string]float64, len(balance))
	var i int
	for i = 0; i < len(balance); i++ {
		balances[string(i)] = float64(balance[i].Amount)
	}
	return balances, nil
}

func (q *QuantaGraphene) GetAsset(assetName string) (*database.Asset, error) {
	asset, err := q.Database.LookupAssetSymbols(assetName)
	if err != nil {
		return nil, err
	}
	if len(asset) == 0 || asset[0] == nil {
		return nil, errors.New("asset does not exist")
	}

	return asset[0], nil
}

func (q *QuantaGraphene) GetIssuer() string {
	return q.QuantaClientOptions.Issuer
}

func (q *QuantaGraphene) AccountExist(quantaAddr string) bool {
	accountMap, err := q.Database.LookupAccounts(quantaAddr, 1)
	if err != nil {
		return false
	}
	if accountMap[quantaAddr].Space == 0 && accountMap[quantaAddr].Type == 0 {
		return false
	}

	return true
}

func (q *QuantaGraphene) LookupAccount(account string) (types.ObjectID, error) {
	var accountId types.ObjectID
	accountMap, err := q.Database.LookupAccounts(account, 1)
	if err != nil {
		return accountId, err
	}
	if accountMap[account].Space == 0 && accountMap[account].Type == 0 {
		return accountId, errors.New("account not found")
	}
	return accountMap[account], nil
}

// https://github.com/scorum/bitshares-go/blob/bbfc9bedaa1b2ddaead3eafe47237efcd9b8496d/client.go
func (q *QuantaGraphene) CreateTransferProposal(dep *coin.Deposit) (string, error) {
	var fee types.AssetAmount
	var amount types.AssetAmount

	id, err := q.GetAsset(dep.CoinName)
	if err != nil {
		return "", err
	}
	amount.Amount = uint64(dep.Amount)
	amount.AssetID = id.ID

	var QDEX types.ObjectID
	QDEX.Space = 1
	QDEX.Type = 3
	QDEX.ID = 0

	fee.Amount = 0
	fee.AssetID = QDEX

	userIdSender, err := q.LookupAccount(q.QuantaClientOptions.Issuer)
	if err != nil {
		return "", err
	}

	userIdReceiver, err := q.LookupAccount(dep.QuantaAddr)
	if err != nil {
		return "", err
	}

	op := types.NewTransferOperation(userIdSender, userIdReceiver, amount, fee)

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		return "", err
	}
	if fees == nil {
		return "", errors.New("cannot calculate the fees")
	}

	op.Fee.Amount = fees[0].Amount

	return q.PrepareTX(op)
}

func (q *QuantaGraphene) CreateNewAssetProposal(issuer string, symbol string, precision uint8) (string, error) {
	issuerId, err := q.LookupAccount(issuer)
	if err != nil {
		return "", err
	}

	var baseObject, quoteObject types.ObjectID
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

	p := &types.Price{
		Base:  base,
		Quote: quote,
	}
	var d = make([]types.ObjectID, 0)
	var extensions = make([]json.RawMessage, 0)
	assetOptions := &types.AssetOptions{
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
		Description:          "QUANTA crosschain for " + symbol,
		Extensions:           extensions,
	}
	w := &types.CreateAsset{
		Issuer:             issuerId,
		Symbol:             symbol,
		Precision:          precision,
		CommonOptions:      *assetOptions,
		IsPredictionMarket: false,
		Extensions:         extensions,
	}
	w.Fee.AssetID = baseObject
	fees, err := q.Database.GetRequiredFee([]types.Operation{w}, fee.AssetID.String())
	if err != nil {
		return "", err
	}
	if fees == nil {
		return "", errors.New("cannot calculate fees")
	}
	w.Fee.Amount = fees[0].Amount
	return q.PrepareTX(w)
}

func (q *QuantaGraphene) CreateIssueAssetProposal(dep *coin.Deposit) (string, error) {
	issuerId, err := q.LookupAccount(q.QuantaClientOptions.Issuer)
	if err != nil {
		return "", err
	}
	accountId, err := q.LookupAccount(dep.QuantaAddr)
	if err != nil {
		return "", err
	}

	assetId, err := q.GetAsset(dep.CoinName)
	if err != nil {
		return "", err
	}

	var asset types.AssetAmount
	asset.Amount = uint64(dep.Amount)
	asset.AssetID = assetId.ID

	var QDEX types.ObjectID
	QDEX.Space = 1
	QDEX.Type = 3
	QDEX.ID = 0

	var fee types.AssetAmount
	fee.Amount = 0
	fee.AssetID = QDEX

	w := &types.IssueAsset{
		Fee:            fee,
		Issuer:         issuerId,
		AssetToIssue:   asset,
		IssueToAccount: accountId,
		Extensions:     []json.RawMessage{},
	}

	w.Fee.AssetID = QDEX
	fees, err := q.Database.GetRequiredFee([]types.Operation{w}, fee.AssetID.String())
	if err != nil {
		return "", err
	}
	if fees == nil {
		return "", errors.New("cannot calculate fees")
	}

	w.Fee.Amount = fees[0].Amount
	return q.PrepareTX(w)
}

func (q *QuantaGraphene) PrepareTX(operations ...types.Operation) (string, error) {
	props, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return "", err
	}

	block, err := q.Database.GetBlock(props.HeadBlockNumber)
	if err != nil {
		return "", err
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return "", err
	}

	expiration := props.Time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.HeadBlockNumber - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}

	data, err := json.Marshal(stx)
	if err != nil {
		return "", err
	}
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
		return "", err
	}
	if fees == nil {
		return "", errors.New("cannot calculate fees")
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
		return "", err
	}
	if fees == nil {
		return "", errors.New("cannot calculate fees")
	}
	op.Fee.Amount = fees[0].Amount

	return q.PrepareTX(op)
}

func (q *QuantaGraphene) DecodeTransaction(base64 string) (*coin.Deposit, error) {
	var tx sign.SignedTransaction
	err := json.Unmarshal([]byte(base64), &tx)
	if err != nil {
		return nil, err
	}

	op := tx.Operations[0]
	if op.Type() == types.TransferOpType {
		op := op.(*types.TransferOperation)

		receiver, err := q.Database.GetObjects(op.To)
		if err != nil {
			return nil, err
		}
		if receiver == nil {
			return nil, errors.New("error in decoding, receiver not found")
		}

		to := &Object{}
		err = json.Unmarshal(receiver[0], &to)
		if err != nil {
			return nil, err
		}

		coinName, err := q.Database.GetObjects(op.Amount.AssetID)
		if err != nil {
			return nil, err
		}
		if coinName == nil {
			return nil, errors.New("error in decoding, coinName not found")
		}

		asset := &Asset{}
		err = json.Unmarshal(coinName[0], &asset)
		if err != nil {
			return nil, err
		}

		sender, err := q.Database.GetObjects(op.From)
		if err != nil {
			return nil, err
		}
		if sender == nil {
			return nil, errors.New("error in decoding, sender not found")
		}

		from := &Object{}
		err = json.Unmarshal(sender[0], &from)
		if err != nil {
			return nil, err
		}

		return &coin.Deposit{CoinName: asset.Symbol,
			QuantaAddr: to.Name,
			Amount:     int64(op.Amount.Amount),
			BlockID:    0,
			Type:       types.TransferOpType,
		}, nil
	} else if op.Type() == types.IssueAssetOpType {
		op := op.(*types.IssueAsset)

		receiver, err := q.Database.GetObjects(op.IssueToAccount)
		if err != nil {
			return nil, err
		}
		if receiver == nil {
			return nil, errors.New("error in decoding, receiver not found")
		}

		to := &Object{}
		err = json.Unmarshal(receiver[0], &to)
		if err != nil {
			return nil, err
		}

		coinName, err := q.Database.GetObjects(op.AssetToIssue.AssetID)
		if err != nil {
			return nil, err
		}
		if coinName == nil {
			return nil, errors.New("error in decoding, coinName not found")
		}

		asset := &Asset{}
		err = json.Unmarshal(coinName[0], &asset)
		if err != nil {
			return nil, err
		}

		sender, err := q.Database.GetObjects(op.Issuer)
		if err != nil {
			return nil, err
		}
		if sender == nil {
			return nil, errors.New("error in decoding, sender not found")
		}

		from := &Object{}
		err = json.Unmarshal(sender[0], &from)
		if err != nil {
			return nil, err
		}
		return &coin.Deposit{CoinName: asset.Symbol,
			QuantaAddr: to.Name,
			Amount:     int64(op.AssetToIssue.Amount),
			BlockID:    0,
			Type:       types.IssueAssetOpType,
		}, nil
	} else if op.Type() == types.CreateAssetOpType {
		op := op.(*types.CreateAsset)
		return &coin.Deposit{CoinName: op.Symbol,
			Type: types.CreateAssetOpType,
		}, nil

	}

	return nil, nil
}

func ProcessGrapheneTransaction(proposed string, sigs []string) (string, error) {
	var tx sign.SignedTransaction
	err := json.Unmarshal([]byte(proposed), &tx)
	if err != nil {
		return "", err
	}
	tx.Transaction.Signatures = sigs
	signed, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return string(signed), err
}

func (q *QuantaGraphene) ProcessDeposit(deposit *coin.Deposit, proposed string) error {
	txe, err := ProcessGrapheneTransaction(proposed, deposit.Signatures)
	if err != nil {
		return err
	}
	return db.ChangeSubmitQueue(q.Db, deposit.Tx, txe, db.DEPOSIT)
}
