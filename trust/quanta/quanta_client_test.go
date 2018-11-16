package quanta

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func _GetClient() (Quanta, error){
	client, err := NewQuanta(QuantaClientOptions{
		Logger: logger.NewGoLogger("test"),
		Network: "QUANTA Test Network ; September 2018",
		Issuer: "QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB",
		HorizonUrl: "http://testnet-02.quantachain.io:8000/",
	})
	if err != nil {
		return nil, err
	}

	client.Attach()
	return client, nil
}

func TestGetTopID(t *testing.T) {
	client, err := _GetClient()

	if err != nil {
		t.Error(err)
		return
	}

	// quoc1 account
	accountId := "QB3WOAL55IVT6E7BVUNRW6TUVCAOPH5RJYPUUL643YMKMJSZFZGWDJU3"
	maxId, err := client.GetTopBlockID(accountId)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("MaxID=%d\n", maxId)

	refunds, _, _ := client.GetRefundsInBlock(0, accountId)
	fmt.Printf("Refunds = %v\n", refunds)
}


func TestCreateProposeTransaction(t *testing.T) {
	client, err := _GetClient()

	deposit := coin.Deposit{
		"ACME",
		"QARMQQVXLEUCTUYXVGBXOQ6BTO7EFCG42KO5RLWEMTFP4XU4BIF6ATBI",
		"QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5",
		42,
		0,  // ignored?
	}

	encoded, err := client.CreateProposeTransaction(&deposit)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)

	txe := &xdr.TransactionEnvelope{}
	err = xdr.SafeUnmarshalBase64(encoded, txe)
	assert.NotNil(t, txe)
	assert.NoError(t, err)

	ops := txe.Tx.Operations
	assert.Equal(t, 1, len(ops))

	payOp, success := ops[0].Body.GetPaymentOp()
	assert.True(t, success)
	assert.NotNil(t, payOp)

	// TODO: shouldn't this be 'ACME'?
	assert.Regexp(t, regexp.MustCompile("/ACME/"), payOp.Asset.String())
	// TODO: test for SenderAddress?
	assert.Equal(t, "QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5", payOp.Destination.Address())
	assert.Equal(t, int64(42), int64(payOp.Amount))
	// TODO: test for BlockId?
}

func TestDecodeTransaction(t *testing.T) {
	client, err := _GetClient()

	odeposit := coin.Deposit{
		"ACME",
		"QARMQQVXLEUCTUYXVGBXOQ6BTO7EFCG42KO5RLWEMTFP4XU4BIF6ATBI",
		"QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5",
		42,
		0,  // ignored?
	}

	encoded, _ := client.CreateProposeTransaction(&odeposit)

	deposit, err := client.DecodeTransaction(encoded)
	assert.NoError(t, err)
	assert.NotNil(t, deposit)

	// TODO: shouldn't this be 'ACME'?
	assert.Regexp(t, regexp.MustCompile("/ACME/"), deposit.CoinName)

	// TODO: should this really be blank?
	assert.Equal(t, "", deposit.SenderAddr)
	assert.Equal(t, "QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5", deposit.QuantaAddr)
	assert.Equal(t, int64(42), deposit.Amount)

	// TODO: should this really be 0?
	assert.Equal(t, int64(0), deposit.BlockID)
}

func TestCreateProposeTransactionInvalidAmount(t *testing.T) {
	client, _ := _GetClient()

	odeposit := coin.Deposit{
		"ACME",
		"QARMQQVXLEUCTUYXVGBXOQ6BTO7EFCG42KO5RLWEMTFP4XU4BIF6ATBI",
		"QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5",
	  999999999999999 + 1,
		0,
	}

	_, err := client.CreateProposeTransaction(&odeposit)

	assert.Error(t, err)
}

func TestCreateProposeTransactionNinesAmount(t *testing.T) {
	client, _ := _GetClient()

	odeposit := coin.Deposit{
		"ACME",
		"QARMQQVXLEUCTUYXVGBXOQ6BTO7EFCG42KO5RLWEMTFP4XU4BIF6ATBI",
		"QAJUT2FOY66CDSB6TNLOQSJHL4STHF2HFTIGTMJ7XNRNQDIKPBPG42H5",
	  999999999999999,
		0,  // ignored?
	}

	encoded, _ := client.CreateProposeTransaction(&odeposit)

	txe := &xdr.TransactionEnvelope{}
	_ = xdr.SafeUnmarshalBase64(encoded, txe)
	payOp, _ := txe.Tx.Operations[0].Body.GetPaymentOp()

	assert.Equal(t, int64(999999999999999), int64(payOp.Amount))
}
