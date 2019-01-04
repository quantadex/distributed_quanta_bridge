package types

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/scorum/bitshares-go/encoding/transaction"
	"reflect"
)

type Operation interface {
	Type() OpType
}

type Operations []Operation

func (ops *Operations) UnmarshalJSON(b []byte) (err error) {
	// unmarshal array
	var o []json.RawMessage
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	// foreach operation
	for _, op := range o {
		var kv []json.RawMessage
		if err := json.Unmarshal(op, &kv); err != nil {
			return err
		}

		if len(kv) != 2 {
			return errors.New("invalid operation format: should be name, value")
		}

		var opType uint16
		if err := json.Unmarshal(kv[0], &opType); err != nil {
			return err
		}

		val, err := unmarshalOperation(OpType(opType), kv[1])
		if err != nil {
			return err
		}

		*ops = append(*ops, val)
	}

	return nil
}

type operationTuple struct {
	Type OpType
	Data Operation
}

func (op *operationTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		op.Type,
		op.Data,
	})
}

func (ops Operations) MarshalJSON() ([]byte, error) {
	tuples := make([]*operationTuple, 0, len(ops))
	for _, op := range ops {
		if op.Type() == TransferOpType || op.Type() == CreateAssetOpType || op.Type() == IssueAssetOpType {
			tuples = append(tuples, &operationTuple{
				Type: op.Type(),
				Data: op,
			})
		}
	}
	return json.Marshal(tuples)
}

func unmarshalOperation(opType OpType, obj json.RawMessage) (Operation, error) {
	op, ok := knownOperations[opType]
	if !ok {
		// operation is unknown wrap it as an unknown operation
		val := UnknownOperation{
			kind: opType,
			Data: obj,
		}
		println("unknown operation????")
		return &val, nil
	} else {
		val := reflect.New(op).Interface()
		if err := json.Unmarshal(obj, val); err != nil {
			return nil, err
		}
		return val.(Operation), nil
	}
}

var knownOperations = map[OpType]reflect.Type{
	TransferOpType:         reflect.TypeOf(TransferOperation{}),
	LimitOrderCreateOpType: reflect.TypeOf(LimitOrderCreateOperation{}),
	LimitOrderCancelOpType: reflect.TypeOf(LimitOrderCancelOperation{}),
	CreateAssetOpType:      reflect.TypeOf(CreateAsset{}),
	IssueAssetOpType:       reflect.TypeOf(IssueAsset{}),
}

// UnknownOperation
type UnknownOperation struct {
	kind OpType
	Data json.RawMessage
}

func (op *UnknownOperation) Type() OpType { return op.kind }

// NewTransferOperation returns a new instance of TransferOperation
func NewTransferOperation(from, to ObjectID, amount, fee AssetAmount) *TransferOperation {
	op := &TransferOperation{
		From:       from,
		To:         to,
		Amount:     amount,
		Fee:        fee,
		Extensions: []json.RawMessage{},
	}

	return op
}

// TransferOperation
type TransferOperation struct {
	From       ObjectID          `json:"from"`
	To         ObjectID          `json:"to"`
	Amount     AssetAmount       `json:"amount"`
	Fee        AssetAmount       `json:"fee"`
	Memo       *Memo             `json:"memo,omitempty"`
	Extensions []json.RawMessage `json:"extensions"`
}

type Memo struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Nonce   string `json:"nonce"`
	Message string `json:"message"`
}

func (op *TransferOperation) Type() OpType { return TransferOpType }

func (op *TransferOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.From)
	enc.Encode(op.To)
	enc.Encode(op.Amount)

	//Memo?
	enc.EncodeUVarint(0)
	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

type CreateAsset struct {
	Fee                AssetAmount       `json:"fee"`
	Issuer             ObjectID          `json:"issuer"`
	Symbol             string            `json:"symbol"`
	Precision          uint8             `json:"precision"`
	CommonOptions      AssetOptions      `json:"common_options"`
	IsPredictionMarket bool              `json:"is_prediction_market"`
	Extensions         []json.RawMessage `json:"extensions"`
}

type AssetOptions struct {
	MaxSupply            int64             `json:"max_supply"`
	MarketFeePercent     uint16            `json:"market_fee_percent"`
	MaxMarketFee         int64             `json:"max_market_fee"`
	IssuerPermissions    uint16            `json:"issuer_permissions"`
	Flags                uint16            `json:"flags"`
	CoreExchangeRate     Price             `json:"core_exchange_rate"`
	WhiteListAuthorities []ObjectID        `json:"white_list_authorities"`
	BlackListAuthorities []ObjectID        `json:"black_list_authorities"`
	WhiteListMarkets     []ObjectID        `json:"white_list_markets"`
	BlackListMarkets     []ObjectID        `json:"black_list_markets"`
	Description          string            `json:"description"`
	Extensions           []json.RawMessage `json:"extensions"`
}

type BitAssetOptions struct {
	FeedLifeTimeSec              uint32   `json:"feed_lifetime_sec"`
	MinimumFeeds                 uint8    `json:"minimum_feeds"`
	ForceSettlementDelaySec      uint32   `json:"force_settlement_delay_sec"`
	ForceSettlementOffsetPercent uint16   `json:"force_settlement_offset_percent"`
	MaximumForceSettlementVolume uint16   `json:"maximum_force_settlement_volume"`
	ShortBackingAsset            ObjectID `json:"short_backing_asset"`
}

func (op *CreateAsset) Type() OpType { return CreateAssetOpType }

func (op *CreateAsset) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Issuer)
	enc.Encode(op.Symbol)
	enc.Encode(op.Precision)
	enc.Encode(op.CommonOptions.MaxSupply)
	enc.Encode(op.CommonOptions.MarketFeePercent)
	enc.Encode(op.CommonOptions.MaxMarketFee)
	enc.Encode(op.CommonOptions.IssuerPermissions)
	enc.Encode(op.CommonOptions.Flags)
	enc.Encode(op.CommonOptions.CoreExchangeRate)
	enc.EncodeUVarint(0)
	enc.EncodeUVarint(0)
	enc.EncodeUVarint(0)
	enc.EncodeUVarint(0)
	enc.Encode(op.CommonOptions.Description)
	enc.EncodeUVarint(0)

	// bitassets_opts
	enc.EncodeUVarint(0)

	enc.EncodeBool(op.IsPredictionMarket)

	//extensions
	enc.EncodeUVarint(0)

	//AssetOptions
	//enc.EncodeUVarint(0)
	//BitAssetOptions
	//enc.EncodeUVarint(0)
	return enc.Err()
}

type IssueAsset struct {
	Fee            AssetAmount       `json:"fee"`
	Issuer         ObjectID          `json:"issuer"`
	AssetToIssue   AssetAmount       `json:"asset_to_issue"`
	IssueToAccount ObjectID          `json:"issue_to_account"`
	Memo           *Memo             `json:"memo,omitempty"`
	Extensions     []json.RawMessage `json:"extensions"`
}

type MemoData struct {
	Nonce   uint64 `json:"nonce"`
	Message []byte `json:"message"`
}

func (op *IssueAsset) Type() OpType { return IssueAssetOpType }

func (op *IssueAsset) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Issuer)
	enc.Encode(op.AssetToIssue)
	enc.Encode(op.IssueToAccount)

	//Memo?
	enc.EncodeUVarint(0)

	enc.EncodeUVarint(0)
	return enc.Err()
}

// LimitOrderCreateOperation
type LimitOrderCreateOperation struct {
	Fee          AssetAmount       `json:"fee"`
	Seller       ObjectID          `json:"seller"`
	AmountToSell AssetAmount       `json:"amount_to_sell"`
	MinToReceive AssetAmount       `json:"min_to_receive"`
	Expiration   Time              `json:"expiration"`
	FillOrKill   bool              `json:"fill_or_kill"`
	Extensions   []json.RawMessage `json:"extensions"`
}

func (op *LimitOrderCreateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)

	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Seller)
	enc.Encode(op.AmountToSell)
	enc.Encode(op.MinToReceive)
	enc.Encode(op.Expiration)
	enc.EncodeBool(op.FillOrKill)

	//extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

func (op *LimitOrderCreateOperation) Type() OpType { return LimitOrderCreateOpType }

// LimitOrderCancelOpType
type LimitOrderCancelOperation struct {
	Fee              AssetAmount       `json:"fee"`
	FeePayingAccount ObjectID          `json:"fee_paying_account"`
	Order            ObjectID          `json:"order"`
	Extensions       []json.RawMessage `json:"extensions"`
}

func (op *LimitOrderCancelOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)

	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.FeePayingAccount)
	enc.Encode(op.Order)

	// extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}

func (op *LimitOrderCancelOperation) Type() OpType { return LimitOrderCancelOpType }

// FillOrderOpType
type FillOrderOperation struct {
	Order   ObjectID
	Account ObjectID
	Pays    AssetAmount
	Recives AssetAmount
	Fee     AssetAmount
	Price   Price
	IsMaker bool
}

func (op *FillOrderOperation) Type() OpType { return FillOrderOpType }
