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

func (db *DB) GetCrosschainByAddressandBlockchain(address, blockchain string) CrosschainAddress {
	var tx CrosschainAddress
	err := db.Model(&tx).Where("address=? and blockchain=?", address, blockchain).Select()
	if err != nil {
		return tx
	}
	return tx
}

func (db *DB) GetCrosschainAll() []CrosschainAddress {
	var tx []CrosschainAddress
	out := []CrosschainAddress{}

	err := db.Model(&tx).Order("blockchain asc", "address asc").Select()
	for _, tx := range tx {
		tx.Updated = time.Unix(0, 0)
		out = append(out, tx)
	}
	if err != nil {
		return nil
	}
	return out
}

func getDifference(mine, in []CrosschainAddress) (changed, added, deleted []CrosschainAddress) {
	mineMap := map[string]CrosschainAddress{}
	for _, r := range mine {
		mineMap[r.Address] = r
	}

	newMap := map[string]CrosschainAddress{}
	for _, r := range in {
		newMap[r.Address] = r
	}

	for k, v := range mineMap {
		if _, exist := newMap[k]; !exist {
			deleted = append(deleted, v)
		}
	}

	for k, v := range newMap {
		// exist in the old orders
		if n, exist := mineMap[k]; exist {
			if v.QuantaAddr != n.QuantaAddr || v.LastBlockNumber != n.LastBlockNumber ||
				v.TxHash != n.TxHash || v.Shared != n.Shared || v.Blockchain != n.Blockchain {
				changed = append(changed, v) // send back old one so we can cancel
			}
		} else {
			added = append(added, v)
		}
	}
	return
}

// given an input "in", repair the local database as neccessary
func (db *DB) RepairCrosschain(in []CrosschainAddress) error {
	changed, added, deleted := getDifference(db.GetCrosschainAll(), in)

	for _, tx := range append(deleted, changed...) {
		_, err := db.Model(&tx).Where("address=?", tx.Address).Delete()
		if err != nil {
			return err
		}
	}

	for _, tx := range added {
		_, err := db.Model(&tx).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) GetCrosschainByBlockchain(blockchain string) []CrosschainAddress {
	var tx []CrosschainAddress
	err := db.Model(&tx).Where("blockchain=?", blockchain).Select()
	if err != nil {
		return nil
	}
	return tx
}

func (db *DB) GetAddressCountByBlockchain(blockchain string) (int, error) {
	var tx []CrosschainAddress
	n, err := db.Model(&tx).Where("blockchain=?", blockchain).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (db *DB) GetAddressCountByBlockchain24hrs(blockchain string) (int, error) {
	t := time.Now().AddDate(0, 0, -1)
	var tx []CrosschainAddress
	n, err := db.Model(&tx).Where("blockchain=? and updated>?", blockchain, t).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
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
	lastBlock := uint64(1)
	if input.QuantaAddr == "address-pool" {
		shared = true
	}
	tx := &CrosschainAddress{input.ContractAddress, input.QuantaAddr,
		input.TxHash, input.Blockchain, shared, lastBlock, time.Now()}
	_, err := db.Model(tx).Insert()
	return err
}

func (db *DB) GetAvailableShareAddress(head_block_number int64, min_block int64) ([]CrosschainAddress, error) {
	var address []CrosschainAddress

	err := db.Model(&address).
		Where("shared=true").
		WrapWith("shared_addresses").
		Table("shared_addresses").
		Where("? - last_block_number > ?", head_block_number, min_block).Order("last_block_number asc", "address asc").Select()

	return address, err
}

func (db *DB) UpdateShareAddressDestination(address string, quantaAddr string, headBlock uint64) error {
	tx := &CrosschainAddress{Address: address}
	tx.QuantaAddr = quantaAddr
	tx.LastBlockNumber = headBlock
	_, err := db.Model(tx).Column("quanta_addr", "last_block_number").Where("Address=?", address).Update()
	return err
}

func (db *DB) UpdateCrosschainAddrBlockNumber(address string, blockNumber uint64) error {
	tx := &CrosschainAddress{Address: address}
	tx.LastBlockNumber = blockNumber
	_, err := db.Model(tx).Column("last_block_number").Where("Address=?", address).Update()
	return err
}
