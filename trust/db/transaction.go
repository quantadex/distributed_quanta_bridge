package db

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
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
const ENCODE_FAILURE = "encode_fatal"
const DUPLICATE_ASSET = "duplicate_asset"
const BAD_ADDRESS = "bad_address"
const SUBMIT_FAILURE = "submit_failure"
const PENDING = "pending"
const WAIT_FOR_CONFIRMATION = "wait_for_confirmation"
const ORPHAN = "orphan"

type Transaction struct {
	Type                string `sql:"unique:type_tx_block_hash"`
	Tx                  string `sql:"unique:type_tx_block_hash"`
	TxId                uint64
	Coin                string
	Created             time.Time
	Amount              int64
	BlockId             int64
	BlockHash           string `sql:"unique:type_tx_block_hash,notnull"`
	From                string
	To                  string
	Signed              bool `sql:",notnull"`
	IsBounced           bool `sql:",notnull"`
	SubmitState         string
	SubmitTx            string
	SubmitSigners       string
	SubmitConfirm_block int
	SubmitDate          time.Time
	SubmitTxHash        string
}

func QueryAllTX(db *DB, offset int, limit int) ([]Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Offset(offset).Limit(limit).Order("created desc").Select()
	return txs, err
}

func QueryAllTXByUser(db *DB, user string, offset int, limit int) ([]Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Offset(offset).Limit(limit).Where("\"from\" = ? OR \"to\" = ?", user, user).Order("created desc").Select()
	return txs, err
}

func QueryAllWaitForConfirmTx(db *DB, blockchain string) ([]Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Where("\"submit_state\" = ? AND \"coin\" = ?", WAIT_FOR_CONFIRMATION, blockchain).Select()
	return txs, err
}

func QueryAllWaitForConfirmTxETH(db *DB, coin string) ([]Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Where("Type=? and Submit_State=? AND (coin =? OR coin like'%0X%') ", DEPOSIT, WAIT_FOR_CONFIRMATION, coin).Select()
	return txs, err
}

func ConfirmDeposit(db *DB, dep *coin.Deposit, isBounced bool) error {
	tx := &Transaction{
		Type:        DEPOSIT,
		Tx:          dep.Tx,
		Coin:        dep.CoinName,
		Created:     time.Now(),
		Amount:      dep.Amount,
		BlockId:     dep.BlockID,
		From:        dep.SenderAddr,
		To:          dep.QuantaAddr,
		IsBounced:   isBounced,
		Signed:      false,
		SubmitState: SUBMIT_CONSENSUS,
		BlockHash:   dep.BlockHash,
	}
	_, err := db.Model(tx).OnConflict("(Type,Tx,Block_Hash) DO UPDATE").Set("Submit_State = EXCLUDED.Submit_State").Insert()

	return err
}

func WaitForConfirmation(db *DB, dep *coin.Deposit, isBounced bool) error {
	tx := &Transaction{
		Type:        DEPOSIT,
		Tx:          dep.Tx,
		Coin:        dep.CoinName,
		Created:     time.Now(),
		Amount:      dep.Amount,
		BlockId:     dep.BlockID,
		From:        dep.SenderAddr,
		To:          dep.QuantaAddr,
		IsBounced:   isBounced,
		Signed:      false,
		SubmitState: WAIT_FOR_CONFIRMATION,
		BlockHash:   dep.BlockHash,
	}
	_, err := db.Model(tx).OnConflict("(Type,Tx,Block_Hash) DO UPDATE").Set("Submit_State = EXCLUDED.Submit_State").Insert()

	return err
}

func AddPendingDeposits(db *DB, deposits []*coin.Deposit) error {
	for _, dep := range deposits {
		tx := &Transaction{
			Type:        DEPOSIT,
			Tx:          dep.Tx,
			Coin:        dep.CoinName,
			Created:     time.Now(),
			Amount:      dep.Amount,
			BlockId:     dep.BlockID,
			From:        dep.SenderAddr,
			To:          dep.QuantaAddr,
			IsBounced:   false,
			Signed:      false,
			SubmitState: PENDING,
		}
		_, err := db.Model(tx).
			Where("Type = ? AND Tx = ?", tx.Type, tx.Tx).
			OnConflict("DO NOTHING").SelectOrInsert()
		return err
	}
	return nil
}

func RemovePending(db *DB, txHash string) error {
	_, err := db.Model(&Transaction{}).Where("Tx=? AND Type=? and Submit_State=?", txHash, DEPOSIT, PENDING).Delete()
	return err
}

func SignDeposit(db *DB, dep *coin.Deposit) error {
	tx := &Transaction{Type: DEPOSIT, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Signed = true
	_, err := db.Model(tx).Column("signed").Where("Tx=? and Block_Hash=?", dep.Tx, dep.BlockHash).Update()
	return err
}

func ChangeSubmitQueue(db *DB, id string, submitTx string, typeStr string, blockHash string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = SUBMIT_QUEUE
	tx.SubmitTx = submitTx
	_, err := db.Model(tx).Column("submit_state", "submit_date", "submit_tx").Where("Tx=? and Type=? and Block_Hash=?", id, typeStr, blockHash).Returning("*").Update()
	return err
}

func ChangeSubmitState(db *DB, id string, state string, typeStr string, blockHash string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = state
	_, err := db.Model(tx).Column("submit_state", "submit_date").Where("Tx=? and Type=? and Block_Hash=?", id, typeStr, blockHash).Returning("*").Update()
	return err
}

func ChangeDepositSubmitState(db *DB, id string, state string, blocknumber int, txhash string, blockHash string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = state
	tx.SubmitConfirm_block = blocknumber
	tx.SubmitTxHash = txhash
	_, err := db.Model(tx).Column("submit_state", "submit_date", "submit_confirm_block", "submit_tx_hash").Where("Tx=? and Type=? and Block_Hash=?", id, DEPOSIT, blockHash).Returning("*").Update()
	return err
}

func ChangeWithdrawalSubmitState(db *DB, id string, state string, txid uint64, txhash string, blockHash string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.SubmitState = state
	tx.TxId = txid
	tx.SubmitTxHash = txhash
	_, err := db.Model(tx).Column("submit_state", "submit_date", "tx_id", "submit_tx_hash").Where("Tx=? and Type=? and Block_Hash=?", id, WITHDRAWAL, blockHash).Returning("*").Update()
	return err
}

func ChangeWithdrawalSubmitTx(db *DB, id string, txid uint64, submitTx string, blockHash string) error {
	tx := &Transaction{Tx: id}
	tx.SubmitDate = time.Now()
	tx.TxId = txid
	tx.SubmitTx = submitTx
	_, err := db.Model(tx).Column("submit_date", "tx_id", "submit_tx").Where("Tx=? and Type=? and Block_Hash=?", id, WITHDRAWAL, blockHash).Returning("*").Update()
	return err
}

func ConfirmWithdrawal(db *DB, dep *coin.Withdrawal) error {
	tx := &Transaction{
		Type:        WITHDRAWAL,
		Tx:          dep.Tx,
		TxId:        dep.TxId,
		Coin:        dep.CoinName,
		Created:     time.Now(),
		Amount:      int64(dep.Amount),
		BlockId:     dep.QuantaBlockID,
		From:        dep.SourceAddress,
		To:          dep.DestinationAddress,
		Signed:      false,
		SubmitState: SUBMIT_CONSENSUS,
		BlockHash:   dep.BlockHash,
	}

	_, err := db.Model(tx).
		OnConflict("(Type,Tx,Block_Hash) DO UPDATE").Set("Submit_State = EXCLUDED.Submit_State").Insert()

	return err
}

func SignWithdrawal(db *DB, dep *coin.Withdrawal) error {
	fmt.Println("Submitting withdrawal to", dep.DestinationAddress)
	tx := &Transaction{Type: WITHDRAWAL, Tx: dep.Tx}
	tx.SubmitDate = time.Now()
	tx.Signed = true
	_, err := db.Model(tx).Column("signed").Where("Tx=? and Type=? and Block_Hash=?", dep.Tx, WITHDRAWAL, dep.BlockHash).Update()
	return err
}

func MigrateTx(db *DB) error {
	err := db.CreateTable(&Transaction{})
	if err != nil {
		return err
	}

	db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Exec("ALTER TABLE transactions ADD COLUMN block_hash text")
		_, err = tx.Exec("ALTER TABLE transactions DROP CONSTRAINT transactions_type_tx_key")
		_, err = tx.Exec("ALTER TABLE transactions ADD CONSTRAINT transactions_type_tx_block_hash_key UNIQUE (type, tx, block_hash)")
		return err
	})
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

func GetAllTransaction(db *DB, txID string, txType string) ([]Transaction, error) {
	var txs []Transaction
	err := db.Model(&txs).Where("Tx=? and Type=? and Submit_State=?", txID, txType, WAIT_FOR_CONFIRMATION).Select()
	if err != nil {
		println("unable to get tx: " + err.Error())
		return nil, err
	}

	return txs, nil
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
	if err != nil {
		return err
	}
	_, err = st.Exec()

	st, err = db.db.Prepare("TRUNCATE table crosschain_addresses")
	if err != nil {
		return err
	}

	_, err = st.Exec()

	st, err = db.db.Prepare("TRUNCATE table key_values")
	if err != nil {
		return err
	}

	_, err = st.Exec()

	return err
}
