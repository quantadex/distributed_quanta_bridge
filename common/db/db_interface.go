package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
)

type DB interface {
	Debug()

	Connect(addr, user, pass, database string)

	CreateTable(model interface{}) error

	Insert(model ...interface{}) error

	Model(model ...interface{}) *orm.Query

	Begin() (*pg.Tx, error)

	RunInTransaction(fn func(*pg.Tx) error) error

	Close() error

	GetCrosschainByBlockchain(blockchain string) []crypto.CrosschainAddress

	AddCrosschainAddress(input *crypto.ForwardInput) error
}
