package db

import (
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
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

const DEPOSIT="deposit"
const WITHDRAWAL="withdrawal"

const CONFIRMED="confirmed"
const SIGNED="signed"

type Transaction struct {
	Type   string // deposit | withdrawal
	Tx     string
	Coin   string
	Created time.Time
	Amount int64
	BlockId int64
	From string
	To string
	Status string
	SubmitTx string
	SubmitSigners string
	SubmitConfirm_block int
	SubmitDate time.Time
}

func ConfirmDeposit(db *DB, dep *coin.Deposit) error {
	tx := &Transaction{
		Type: DEPOSIT,
		Tx: dep.Tx,
		Coin: dep.CoinName,
		Created: time.Now(),
		Amount: dep.Amount,
		BlockId: dep.BlockID,
		From: dep.SenderAddr,
		To: dep.QuantaAddr,
		Status: CONFIRMED,
	}

	return db.Insert(tx)
}

func SignDeposit(db *DB, dep *coin.Deposit) error {
	tx := &Transaction{Type: DEPOSIT, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Status = SIGNED
	_, err := db.Model(tx).Column("Status","SubmitDate").Where("Tx",dep.Tx).Returning("*").Update()
	return err
}


func ConfirmWithdrawal(db *DB, dep *coin.Withdrawal) error {
	tx := &Transaction{
		Type: WITHDRAWAL,
		Tx: dep.Tx,
		Coin: dep.CoinName,
		Created: time.Now(),
		Amount: int64(dep.Amount),
		BlockId: dep.QuantaBlockID,
		From: "",
		To: dep.DestinationAddress,
		Status: CONFIRMED,
	}

	return db.Insert(tx)
}

func SignWithdrawal(db *DB, dep *coin.Withdrawal) error {
	tx := &Transaction{Type: WITHDRAWAL, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Status = SIGNED
	_, err := db.Model(tx).Column("Status","SubmitDate").Where("Tx",dep.Tx).Returning("*").Update()
	return err
}

func Migrate(db *DB) error {
	err := db.CreateTable(&Transaction{})
	if err != nil {
		return err
	}
	return err
}