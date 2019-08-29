package db

import (
	"time"
)

type FailedMessage struct {
	Id          string `sql:"unique:id_event"`
	Event       string `sql:"unique:id_event"`
	Quanta      string
	FirstRetry  time.Time
	LastRetry   time.Time
	NoOfRetries int64
}

func MigrateFM(db *DB) error {
	err := db.CreateTable(&FailedMessage{})
	if err != nil {
		return err
	}
	return err
}

func (db *DB) GetAllFailedMessages() []FailedMessage {
	var tx []FailedMessage
	err := db.Model(&tx).Order("first_retry asc").Select()
	if err != nil {
		return tx
	}
	return tx
}

func (db *DB) GetFailedMessageByIdAndEvent(id, event string) (FailedMessage, error) {
	var tx FailedMessage
	err := db.Model(&tx).Where("id=? and event=?", id, event).Select()
	if err != nil {
		return tx, err
	}
	return tx, nil
}

func RemoveFailedMessage(db *DB, id, event string) error {
	_, err := db.Model(&FailedMessage{}).Where("id=? and event=?", id, event).Delete()
	return err
}

func (db *DB) AddFailedMessage(id, event, quanta string, date time.Time) error {
	tx := &FailedMessage{id, event, quanta, date, date, 1}
	_, err := db.Model(tx).Insert()
	return err
}

func (db *DB) UpdateFailedMessage(id, event string, date time.Time, count int64) error {
	tx := &FailedMessage{Id: id, Event: event}
	tx.LastRetry = date
	tx.NoOfRetries = count
	_, err := db.Model(tx).Column("last_retry", "no_of_retries").Where("Id=? and event=?", id, event).Update()
	return err
}
