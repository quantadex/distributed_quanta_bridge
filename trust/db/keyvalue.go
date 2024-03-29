package db

import (
	"github.com/go-pg/pg"
	"time"
)

type KeyValue struct {
	Id      string
	Value   string
	Updated time.Time
}

func MigrateKv(db *DB) error {
	err := db.CreateTable(&KeyValue{})
	if err != nil {
		return err
	}

	db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Exec("ALTER TABLE transactions ADD COLUMN submit_tx_hash text")
		return err
	})

	return err
}

func GetValue(db *DB, id string) *KeyValue {
	tx := &KeyValue{}
	err := db.Model(tx).Where("id=?", id).Select()
	if err != nil {
		return nil
	}
	return tx
}

func RemoveKey(db *DB, id string) error {
	_, err := db.Model(&KeyValue{}).Where("id=?", id).Delete()
	return err
}

func UpdateValue(db *DB, id string, value string) error {
	tx := &KeyValue{id, value, time.Now()}
	_, err := db.Model(tx).OnConflict("(id) DO UPDATE").Set("value = EXCLUDED.value,updated = EXCLUDED.updated").Insert()
	return err
}
