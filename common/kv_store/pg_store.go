package kv_store

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/go-pg/pg"
	"github.com/go-errors/errors"
	"strings"
)

type PostgresKVStore struct {
	db *db.DB
}

func (p *PostgresKVStore) Connect(name string) error {
	p.db = &db.DB{}
	info, err := pg.ParseURL(name)
	if err != nil {
		return err
	}
	p.db.Connect(info.Network, info.User, info.Password, info.Database)
	return nil
}

func (p *PostgresKVStore) CreateTable(tableName string) error {
	return nil
}

func (p *PostgresKVStore) GetValue(tableName string, key string) (value *string, err error) {
	v := db.GetValue(p.db, tableName + "." + key)
	if v == nil {
		return nil, errors.New("Not found")
	}
	return &v.Value, nil
}

func (p *PostgresKVStore) RemoveKey(tableName string, key string) error {
	keyPath := tableName + "." + key
	return db.RemoveKey(p.db, keyPath)
}

func (p *PostgresKVStore) CloseDB() error {
	return p.db.Close()
}

func (p *PostgresKVStore) GetAllValues(tableName string) (map[string]string, error) {
	var values []db.KeyValue

	p.db.Model(values).Where("id LIKE '?.%'", tableName).Select()

	out := map[string]string{}
	for _, k := range values {
		p := strings.Split(k.Id,".")
		if len(p) >= 2 {
			out[p[1]] = k.Value
		}
	}
	return out, nil
}

func (p *PostgresKVStore) SetValue(tableName string, key string, oldValue string, newValue string) error {
	keyPath := tableName + "." + key
	return db.UpdateValue(p.db, keyPath, newValue)
}
