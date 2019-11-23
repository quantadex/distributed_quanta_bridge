package db

import (
	"errors"
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

func (db *DB) GetWebhookByQuantaAndEvent(quanta, event string) ([]Webhook, error) {
	var tx, res []Webhook
	err := db.Model(&tx).Where("quanta=?", quanta).Select()
	if err != nil {
		return tx, err
	}
	for _, t := range tx {
		for _, e := range t.Events {
			if e == event {
				res = append(res, t)
			}
		}
	}
	return res, nil
}

func RemoveWebhook(db *DB, id, quanta string) error {
	_, err := db.Model(&Webhook{}).Where("id=? and quanta=?", id, quanta).Delete()
	return err
}

// de-dup by url - 1 registration per url
func (db *DB) AddWebhook(id string, url string, events []string, quanta string) error {
	exists, err := db.Model(&Webhook{}).Where("url=?", url).Exists()
	if exists {
		return errors.New("only one registration per url")
	}
	tx := &Webhook{id, url, events, quanta}
	_, err = db.Model(tx).Insert()
	return err
}
