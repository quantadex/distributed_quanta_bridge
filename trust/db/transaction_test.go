package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionQuery(t *testing.T) {
	DatabaseUrl := fmt.Sprintf("postgres://postgres:@localhost/crosschain_%d", 1)
	rDb := &DB{}
	info, err := pg.ParseURL(DatabaseUrl)
	if err != nil {
		t.Error(err)
	}
	rDb.Connect(info.Network, info.User, info.Password, info.Database)
	rDb.Debug()
	MigrateTx(rDb)
	EmptyTable(rDb)

	w := &coin.Withdrawal{
		Tx:                 "123",
		CoinName:           "ETH",
		DestinationAddress: "dest",
		QuantaBlockID:      0,
		Amount:             0,
	}

	err = ConfirmWithdrawal(rDb, w)

	if err != nil {
		t.Error(err)
	}

	err = SignWithdrawal(rDb, w)
	if err != nil {
		t.Error(err)
	}

	// updated at time 0
	err = ChangeSubmitState(rDb, w.Tx, SUBMIT_QUEUE, WITHDRAWAL)
	if err != nil {
		t.Error(err)
	}

	// query at time t
	txs := QueryWithdrawalByAge(rDb, time.Now().Add(-time.Second*3), []string{SUBMIT_QUEUE})
	if len(txs) != 0 {
		t.Error("Expecting to have zero items")
	}

	// query at time t - 3
	time.Sleep(time.Second * 4)
	txs = QueryWithdrawalByAge(rDb, time.Now().Add(-time.Second*3), []string{SUBMIT_QUEUE})
	if len(txs) != 1 {
		assert.Equal(t, 1, len(txs))
		t.Error("Expecting to have one items")
	}

	println(txs)

}
