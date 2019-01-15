package db

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"time"
)

/*
 Table: Deposits
  tx   string
  coin string
  created datetime
  amount int64
  blockId int64
  from  string
  to    string
  status enum ('DEPOSIT', 'SIGNED')
  submit_tx  string  // must broadcast
  submit_signers string
  submit_confirm int // keep scanning for confirm
  submit_date  datetime
*/

const DEPOSIT = "deposit"
const WITHDRAWAL = "withdrawal"

const SUBMIT_CONSENSUS = "consensus"
const SUBMIT_QUEUE = "queue"
const SUBMIT_RECOVERABLE = "recoverable"
const SUBMIT_FATAL = "fatal"
const SUBMIT_SUCCESS = "success"

type Transaction struct {
	Type                string `sql:"unique:type_tx"`
	Tx                  string `sql:"unique:type_tx"`
	TxId                uint64
	Coin                string
	Created             time.Time
	Amount              int64
	BlockId             int64
	From                string
	To                  string
	Signed              bool `sql:",notnull"`
	IsBounced           bool `sql:",notnull"`
	SubmitState         string
	SubmitTx            string
	SubmitSigners       string
	SubmitConfirm_block int
	SubmitDate          time.Time
}

func ConfirmDeposit(db *DB, dep *coin.Deposit, isBounced bool) error {
	tx := &Transaction{
		Type:      DEPOSIT,
		Tx:        dep.Tx,
		Coin:      dep.CoinName,
		Created:   time.Now(),
		Amount:    dep.Amount,
		BlockId:   dep.BlockID,
		From:      dep.SenderAddr,
		To:        dep.QuantaAddr,
		IsBounced: isBounced,
		Signed:    false,
	}
	return db.Insert(tx)
}

func SignDeposit(db *DB, dep *coin.Deposit) error {
	tx := &Transaction{Type: DEPOSIT, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Signed = true
	_, err := db.Model(tx).Column("signed").Where("Tx=?", dep.Tx).Update()
	return err
}

func ChangeSubmitQueue(db *DB, id string, submitTx string, typeStr string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = SUBMIT_QUEUE
	tx.SubmitTx = submitTx
	_, err := db.Model(tx).Column("submit_state", "submit_date", "submit_tx").Where("Tx=? and Type=?", id, typeStr).Returning("*").Update()
	return err
}

func ChangeSubmitState(db *DB, id string, state string, typeStr string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = state
	_, err := db.Model(tx).Column("submit_state", "submit_date").Where("Tx=? and Type=?", id, typeStr).Returning("*").Update()
	return err
}

func ConfirmWithdrawal(db *DB, dep *coin.Withdrawal) error {
	tx := &Transaction{
		Type:    WITHDRAWAL,
		Tx:      dep.Tx,
		TxId:    dep.TxId,
		Coin:    dep.CoinName,
		Created: time.Now(),
		Amount:  int64(dep.Amount),
		BlockId: dep.QuantaBlockID,
		From:    dep.SourceAddress,
		To:      dep.DestinationAddress,
		Signed:  false,
	}

	return db.Insert(tx)
}

func SignWithdrawal(db *DB, dep *coin.Withdrawal) error {
	fmt.Println("Submitting withdrawal to", dep.DestinationAddress)
	tx := &Transaction{Type: WITHDRAWAL, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Signed = true
	_, err := db.Model(tx).Column("signed").Where("Tx=? and Type=?", dep.Tx, WITHDRAWAL).Update()
	return err
}

func MigrateTx(db *DB) error {
	err := db.CreateTable(&Transaction{})
	if err != nil {
		return err
	}
	return err
}

func GetTransaction(db *DB, txID string) (*Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Where("Tx=?", txID).Limit(1).Select()
	if err != nil {
		println("unable to get tx: " + err.Error())
		return nil, err
	}

	if len(txs) > 0 {
		return &txs[0], nil
	}
	return nil, errors.New("not found")
}

func QueryDepositByAge(db *DB, age time.Time, states []string) []Transaction {
	var txs []Transaction
	err := db.Model(&txs).Where("Created <= ? AND Type=?", age, DEPOSIT).WhereIn("submit_state IN ?", states).Select()
	if err != nil {
		println("unable to query query: " + err.Error())
	}
	return txs
}

func QueryWithdrawalByAge(db *DB, age time.Time, states []string) []Transaction {
	var txs []Transaction
	err := db.Model(&txs).Where("Created <= ? AND Type=?", age, WITHDRAWAL).WhereIn("submit_state IN ?", states).Select()
	if err != nil {
		println("unable to query query: " + err.Error())
	}
	return txs
}

func EmptyTable(db *DB) error {
	st, err := db.db.Prepare("TRUNCATE table transactions")
	_, err = st.Exec()

	st, err = db.db.Prepare("TRUNCATE table crosschain_addresses")
	_, err = st.Exec()

	st, err = db.db.Prepare("TRUNCATE table key_values")
	_, err = st.Exec()

	return err
}
