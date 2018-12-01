package db

import (
	"testing"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"time"
)

func TestTransactionQuery(t *testing.T) {
	DatabaseUrl := fmt.Sprintf("postgres://postgres:@localhost/crosschain_%d",0)
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
		Tx: "123",
		CoinName: "ETH",
		DestinationAddress: "dest",
		QuantaBlockID: 0,
		Amount : 0,
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
	err = ChangeSubmitState(rDb, w.Tx, SUBMIT_QUEUE)
	if err != nil {
		t.Error(err)
	}

	// query at time t
	txs := QueryDepositByAge(rDb, time.Now().Add(-time.Second*3), []string{ SUBMIT_QUEUE})
	if len(txs) != 0 {
		t.Error("Expecting to have zero items")
	}

	// query at time t - 3
	time.Sleep(time.Second * 4)
	txs = QueryDepositByAge(rDb, time.Now().Add(-time.Second*3), []string{ SUBMIT_QUEUE})
	if len(txs) != 1 {
		t.Error("Expecting to have one items")
	}

	println(txs)

}