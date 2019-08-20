package db

import (
	"github.com/go-pg/pg"
)

type Webhook struct {
	Id     string `sql:",pk"`
	URL    string
	Events []string
	Quanta string
}

func MigrateW(db *DB) error {
	err := db.CreateTable(&Webhook{})
	if err != nil {
		return err
	}

	db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Exec("ALTER TABLE webhooks ADD COLUMN quanta text")
		return err
	})
	return err
}

func (db *DB) GetAllWebhooks() []Webhook {
	var tx []Webhook
	err := db.Model(&tx).Select()
	if err != nil {
		return tx
	}
	return tx
}

func (db *DB) GetWebhooksByQuanta(quanta string) []Webhook {
	var tx []Webhook
	err := db.Model(&tx).Where("quanta=?", quanta).Select()
	if err != nil {
		return tx
	}
	return tx
}

func (db *DB) GetWebhookById(id string) Webhook {
	var tx Webhook
	err := db.Model(&tx).Where("id=?", id).Select()
	if err != nil {
		return tx
	}
	return tx
}

func (db *DB) GetURLForQuanta(quanta, event string) string {
	var tx Webhook
	err := db.Model(&tx).Where("quanta=?", quanta).Select()
	if err != nil {
		return ""
	}
	for _, e := range tx.Events {
		if e == event {
			return tx.URL
		}
	}
	return ""
}

func RemoveWebhook(db *DB, id, quanta string) error {
	_, err := db.Model(&Webhook{}).Where("id=? and quanta=?", id, quanta).Delete()
	return err
}

func (db *DB) AddWebhook(id string, url string, events []string, quanta string) error {
	tx := &Webhook{id, url, events, quanta}
	_, err := db.Model(tx).Insert()
	return err
}
