package db

import (
	"github.com/go-pg/pg"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"time"
)

type CrosschainAddress struct {
	Address         string `sql:",pk"`
	QuantaAddr      string
	TxHash          string
	Blockchain      string
	Shared          bool
	LastBlockNumber uint64 // new
	Updated         time.Time
}

func MigrateXC(db *DB) error {
	err := db.CreateTable(&CrosschainAddress{})
	if err != nil {
		return err
	}
	db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Exec("ALTER TABLE crosschain_addresses ADD COLUMN shared boolean, ADD COLUMN last_block_number bigint")
		return err
	})
	return err
}

func (db *DB) GetCrosschainByBlockchain(blockchain string) []crypto.CrosschainAddress {
	var tx []crypto.CrosschainAddress
	err := db.Model(&tx).Where("blockchain=?", blockchain).Select()
	if err != nil {
		return nil
	}
	return tx
}

func GetCrosschainByBlockchainAndUser(db *DB, blockchain string, quantaAddr string) ([]CrosschainAddress, error) {
	var tx []CrosschainAddress
	err := db.Model(&tx).Where("blockchain=? and quanta_addr=?", blockchain, quantaAddr).Select()
	return tx, err
}

func RemoveCrosschainAddress(db *DB, id string) error {
	_, err := db.Model(&CrosschainAddress{}).Where("id=?", id).Delete()
	return err
}

func (db *DB) AddCrosschainAddress(input *crypto.ForwardInput) error {
	var shared bool
	lastBlock := uint64(0)
	if input.QuantaAddr == "address-pool" {
		shared = true
	}
	tx := &CrosschainAddress{input.ContractAddress, input.QuantaAddr,
		input.TxHash, input.Blockchain, shared, lastBlock, time.Now()}
	_, err := db.Model(tx).Insert()
	return err
}

func (db *DB) GetAvailableShareAddress(head_block_number int64, min_block int64) (CrosschainAddress, error) {
	var address CrosschainAddress

	err := db.Model(&address).
		Where("shared=true").
		WrapWith("shared_addresses").
		Table("shared_addresses").
		Where("last_block_number - ? > ?", head_block_number, min_block).Order("last_block_number asc").Limit(1).Select()

	return address, err
}

func (db *DB) UpdateShareAddressDestination(address string, quantaAddr string) error {
	tx := &CrosschainAddress{Address: address}
	tx.QuantaAddr = quantaAddr
	_, err := db.Model(tx).Column("quanta_addr").Where("Address=?", address).Update()
	return err
}

func (db *DB) UpdateLastBlockNumber(address string, blockNumber uint64) error {
	tx := &CrosschainAddress{Address: address}
	tx.LastBlockNumber = blockNumber
	_, err := db.Model(tx).Column("last_block_number").Where("Address=?", address).Update()
	return err
}
