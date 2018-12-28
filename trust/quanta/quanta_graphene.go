package quanta

/**
 Links:
http://docs.bitshares.org/integration/traders/index.html#public-api
https://github.com/scorum/bitshares-go/blob/master/apis/database/api_test.go
*/
import (
	"encoding/json"
	"strconv"

	//"encoding/json"
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
	//"strconv"
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

const url1 = "ws://testnet-01.quantachain.io:8090"

func (q *QuantaGraphene) Attach() error {
	transport, err := websocket.NewTransport(url1)
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

// get_dynamics
func (q *QuantaGraphene) GetTopBlockID(accountId string) (uint32, error) {
	res, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	blockId := res.HeadBlockNumber

	return blockId, nil
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
					}
					fmt.Println("transaction id = ", txid)
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

func (q *QuantaGraphene) GetAllBalances(quantaAddress string, assets ...string) ([]float64, error) {
	id, err := q.Database.LookupAssetSymbols(assets...)
	fmt.Println("id = ", id)
	ids := make([]types.ObjectID, len(id))
	var i = 0
	fmt.Println(len(ids), len(id))
	for i = 0; i < len(id); i++ {
		ids[i] = id[i].ID
	}
	var balance []*types.AssetAmount
	balance, err = q.Database.GetNamedAccountBalances(quantaAddress, ids...)
	if err != nil {
		return nil, err
	}
	result := make([]float64, len(balance))
	for i = 0; i < len(balance); i++ {
		precision := math.Pow(10, float64(id[i].Precision))
		result[i] = float64(balance[i].Amount) / precision
	}
	return result[:], nil
}

// https://github.com/scorum/bitshares-go/blob/bbfc9bedaa1b2ddaead3eafe47237efcd9b8496d/client.go
func (q *QuantaGraphene) CreateProposeTransaction(dep *coin.Deposit) ([]string, error) {
	var s []string
	var fee types.AssetAmount
	var amount types.AssetAmount

	id, err := q.Database.LookupAssetSymbols(dep.CoinName)
	if err != nil {
		return s, err
	}
	amount.Amount = uint64(dep.Amount)
	amount.AssetID = id[0].ID

	fee.Amount = 0
	fee.AssetID = id[0].ID

	userIdSender, err := q.Database.LookupAccounts(dep.SenderAddr, 1)
	if err != nil {
		return s, err
	}
	userIdReceiver, err := q.Database.LookupAccounts(dep.QuantaAddr, 1)
	if err != nil {
		return s, err
	}

	fmt.Println("something = ", userIdSender[dep.SenderAddr], userIdReceiver[dep.QuantaAddr])

	op := types.NewTransferOperation(userIdSender[dep.SenderAddr], userIdReceiver[dep.QuantaAddr], amount, fee)
	fmt.Println("transaction created = ", op)

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return s, err

	}
	op.Fee.Amount = fees[0].Amount

	wifs := make([]string, 1)
	wifs[0] = "5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj"

	stx, txe, err := q.SignFunc(wifs, op)
	if err != nil {
		fmt.Println("error in signFunc")
		return nil, err
	}

	fmt.Println("signed transaction = ", stx.RefBlockNum)
	fmt.Println("txe = ", txe)

	err = q.NetworkBroadcast.BroadcastTransaction(stx.Transaction)
	if err != nil {
		fmt.Println("error in broadcast")
		return nil, err
	}
	return txe, nil
}

func (q *QuantaGraphene) LimitOrderCreate(key string, feePayingAccount, order types.ObjectID, fee types.AssetAmount) ([]string, error) {
	op := &types.LimitOrderCancelOperation{
		Fee:              fee,
		FeePayingAccount: feePayingAccount,
		Order:            order,
		Extensions:       []json.RawMessage{},
	}
	var s []string

	fees, err := q.Database.GetRequiredFee([]types.Operation{op}, fee.AssetID.String())
	if err != nil {
		log.Println(err)
		return s, err
	}
	op.Fee.Amount = fees[0].Amount

	stx, txe, err := q.SignFunc([]string{key}, op)
	if err != nil {
		return s, err
	}
	return txe, q.NetworkBroadcast.BroadcastTransaction(stx.Transaction)
}

func (q *QuantaGraphene) SignFunc(wifs []string, operations ...types.Operation) (*sign.SignedTransaction, []string, error) {
	var s []string
	props, err := q.Database.GetDynamicGlobalProperties()
	if err != nil {
		return nil, s, err
	}

	block, err := q.Database.GetBlock(props.LastIrreversibleBlockNum)
	if err != nil {
		return nil, s, err
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return nil, s, err
	}

	chainId, err := q.Database.GetChainID()

	expiration := props.Time.Add(10 * time.Minute)
	stx := sign.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.LastIrreversibleBlockNum - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}
	txe, err := stx.Sign(wifs, *chainId)
	if err != nil {
		return nil, nil, err
	}
	return stx, txe, nil
}

/*
func (q *QuantaGraphene) DecodeTransaction(base64 string) (*coin.Deposit, error) {
    txe := types.Transaction{}
    //err := xdr.SafeUnmarshalBase64(base64, txe)

    ops := txe.Operations
    if len(ops) != 1 {
        return nil, errors.New("no operations found")
    }

    paymentOp:= ops[0].Details()
    coinName, err := q.Bitshare.GetObjects(paymentOp.Amount.AssetID)
    result := &Asset{}
    err = json.Unmarshal(coinName[0], &result)
    if err != nil {
        return nil, err
    }

    sender, err := q.Bitshare.GetObjects(paymentOp.From)
    from := &Object{}
    err = json.Unmarshal(sender[0], &from)
    if err != nil {
        return nil, err
    }

    receiver, err := q.Bitshare.GetObjects(paymentOp.To)
    to := &Object{}
    err = json.Unmarshal(receiver[0], &to)
    if err != nil {
        return nil, err
    }

    return &coin.Deposit{CoinName: result.Symbol,
        QuantaAddr: to.Name,
        Amount:     int64(paymentOp.Amount.Amount),
        BlockID:    0,
    }, nil

}

func ProcessTransaction(network string, base64 string, sigs []string) (string, error) {
    return "", nil
}


func (q *QuantaGraphene) ProcessDeposit(deposit *coin.Deposit, proposed string) error {
    txe, err := ProcessTransaction("", proposed, deposit.Signatures)
    println(txe, err)
    return db.ChangeSubmitQueue(q.Db, deposit.Tx, txe, "")
}

/*
func (q *QuantaGraphene) AttachQueue(kv kv_store.KVStore) error {
    panic("implement me")
}


*/
