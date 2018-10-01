package kv_store

import (
	"github.com/boltdb/bolt"
	"fmt"
	"os"
	"github.com/pkg/errors"
)


const dbFile = "%s.db"

func DbFileName(name string) string {
	return fmt.Sprintf(dbFile, name)
}

func DbExists(nodeID string) bool {
	dbFile := DbFileName(nodeID)
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

type BoltStore struct {
	db  *bolt.DB
}

func (s *BoltStore) Connect(name string) error {
	var err error
	s.db, err = bolt.Open(DbFileName(name), 0600, nil)
	return err
}

func (s *BoltStore) CreateTable(tableName string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (s *BoltStore) GetValue(tableName string, key string) (*string, error) {
	var value *string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b == nil {
			return errors.New("bucket not found")
		}

		r := b.Get([]byte(key))
		if r == nil {
			return nil // silent error
		}
		s := string(r)
		value = &s

		return nil
	})

	return value, err
}

func (s *BoltStore) SetValue(tableName string, key string, oldValue string, newValue string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b == nil {
			return errors.New("bucket not found")
		}
		return b.Put([]byte(key), []byte(newValue))
	})
}


