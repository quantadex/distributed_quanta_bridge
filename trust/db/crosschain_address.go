package db

import (
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
)

type CrosschainAddress struct {
	Address      string
	QuantaAddr   string
	TxHash       string
	Blockchain   string
	Updated time.Time
}

func MigrateXC(db *DB) error {
	err := db.CreateTable(&CrosschainAddress{})
	if err != nil {
		return err
	}
	return err
}

func GetCrosschainByBlockchain(db *DB, blockchain string) []CrosschainAddress {
	var tx []CrosschainAddress
	err := db.Model(tx).Where("blockchain=?", blockchain ).Select()
	if err != nil {
		return nil
	}
	return tx
}

func RemoveCrosschainAddress(db *DB, id string) error {
	_, err := db.Model(&CrosschainAddress{}).Where("id=?", id).Delete()
	return err
}

func AddCrosschainAddress(db *DB, input coin.ForwardInput) error {
	tx := &CrosschainAddress{ input.ContractAddress.Hex(), input.QuantaAddr,
							input.TxHash, input.Blockchain, time.Now() }
	_, err := db.Model(tx).Insert()
	return err
}