package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	//"log"
	"fmt"
)

type DB struct {
	db *pg.DB

	Drop     bool
	DebugSQL bool
}

func (db *DB) Debug() {
	db.DebugSQL = true
	db.Drop = true

	//if db.DebugSQL {
	//	db.db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
	//		log.Println(event.FormattedQuery())
	//	})
	//}
}

func (db *DB) Connect(addr, user, pass, database string) {
	if db.Drop {
		db := pg.Connect(&pg.Options{
			Addr:     addr,
			User:     user,
			Password: pass,
			Database: "postgres",
		})

		_, err := db.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'", database))
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("DROP DATABASE IF EXISTS " + database)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("CREATE DATABASE " + database)
		if err != nil {
			panic(err)
		}
		db.Close()
	}
	db.db = pg.Connect(&pg.Options{
		Addr:     addr,
		User:     user,
		Password: pass,
		Database: database,
	})
}

func (db *DB) CreateTable(model interface{}) error {
	return db.db.CreateTable(model, &orm.CreateTableOptions{
		IfNotExists: true,
	})
}

func (db *DB) Insert(model ...interface{}) error {
	return db.db.Insert(model...)
}

func (db *DB) Model(model ...interface{}) *orm.Query {
	return db.db.Model(model...)
}

func (db *DB) Begin() (*pg.Tx, error) {
	return db.db.Begin()
}

func (db *DB) RunInTransaction(fn func(*pg.Tx) error) error {
	return db.db.RunInTransaction(fn)
}

func (db *DB) Close() error {
	return db.db.Close()
}
