package db

import (
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"time"
)

type CrosschainAddress struct {
	Address    string
	QuantaAddr string
	TxHash     string
	Blockchain string
	Updated    time.Time
}

func MigrateXC(db *DB) error {
	err := db.CreateTable(&CrosschainAddress{})
	if err != nil {
		return err
	}
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

func (db *DB) RemoveCrosschainAddress(id string) error {
	_, err := db.Model(&CrosschainAddress{}).Where("id=?", id).Delete()
	return err
}

func (db *DB) AddCrosschainAddress(input *crypto.ForwardInput) error {
	tx := &CrosschainAddress{input.ContractAddress, input.QuantaAddr,
		input.TxHash, input.Blockchain, time.Now()}
	_, err := db.Model(tx).Insert()
	return err
}
