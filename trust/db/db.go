package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"log"

	//"log"
	"fmt"
	"time"
)

type DB struct {
	db *pg.DB

	Drop     bool
	DebugSQL bool
}

type eventHookTest struct {
	beforeQueryMethod func(*pg.QueryEvent)
	afterQueryMethod  func(*pg.QueryEvent)
}

func (e eventHookTest) BeforeQuery(event *pg.QueryEvent) {
	e.beforeQueryMethod(event)
}

func (e eventHookTest) AfterQuery(event *pg.QueryEvent) {
	e.afterQueryMethod(event)
}

func (db *DB) Debug() {
	db.DebugSQL = true
	db.Drop = true
}

func (db *DB) Connect(addr, user, pass, database string) {

	if db.Drop {
		db2 := pg.Connect(&pg.Options{
			Addr:     addr,
			User:     user,
			Password: pass,
			Database: "postgres",
		})

		_, err := db2.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'", database))
		if err != nil {
			panic(err)
		}

		//_, err = db.Exec("DROP DATABASE IF EXISTS " + database)
		//fmt.Println("Dropping database")
		//if err != nil {
		//	panic(err)
		//}

		_, err = db2.Exec("CREATE DATABASE " + database)
		log.Println("Creating database")
		if err != nil {
			//panic(err)
		}
		db2.Close()
	}
	db.db = pg.Connect(&pg.Options{
		Addr:     addr,
		User:     user,
		Password: pass,
		Database: database,
		PoolSize:           20,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	})

	// uncomment to see raw SQL
	//hookImpl := struct{ eventHookTest }{}
	//hookImpl.beforeQueryMethod = func(event *pg.QueryEvent) {
	//	msg, _ := event.FormattedQuery()
	//	log.Println(msg)
	//}
	//hookImpl.afterQueryMethod = hookImpl.beforeQueryMethod
	//
	//db.db.AddQueryHook(hookImpl)
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
