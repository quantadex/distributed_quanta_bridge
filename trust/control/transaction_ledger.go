package control

/**
 * transaction_ledger
 *
 * A set of functions for managing the coin/quanta transaction ledger
 * based on a key_value store.
 */

import (
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"strconv"
)

const (
	QUANTA_CONFIRMED   = "quanta_confirmed"
	COIN_CONFIRMED     = "coin_confirmed"
	LAST_BLOCK         = "last_block"
	ETHADDR_QUANTAADDR = "ethaddr2quantaaddr"
	CONFIRMED          = "C"
	SIGNED             = "S"
	ERROR              = "E"
	NOT_FOUND          = "NF"
)

/**
 * initLedger
 *
 * Is passed an attached KVStore.
 * Creates tables for keeping coin and quanta state
 */
func InitLedger(kv kv_store.KVStore) error {
	err := kv.CreateTable(QUANTA_CONFIRMED)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(COIN_CONFIRMED)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(LAST_BLOCK)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(ETHADDR_QUANTAADDR)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG_SIGNED)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG_FAILED_UNRECOVERABLE)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG_RETRY)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG_FAILED_RECOVERABLE)
	if err != nil {
		return errors.New("Failed to create table")
	}
	err = kv.CreateTable(kv_store.ETH_TX_LOG_SUBMITTED)
	if err != nil {
		return errors.New("Failed to create table")
	}
	return nil
}

/**
 * getKeyName
 *
 * Converts a coinName, destination address and block ID into a single unique string
 * that will be used as a key
 */
func getKeyName(coinName string, dstAddress string, blockID int64) string {
	return fmt.Sprintf("%s-%s-%09d", coinName, dstAddress, blockID)
}

/**
 * getState
 *
 * Returns the state for a given key in a given table.
 * Only states of (ERROR, NOT_FOUND, CONFIRMED, SIGNED) are possible.
 */
func getState(db kv_store.KVStore, table string, k string) string {
	v, err := db.GetValue(table, k)
	if err != nil {
		return ERROR
	}
	if v == nil {
		return NOT_FOUND
	}
	return *v
}

/**
 * confirmTx
 *
 * If a given key does not exist, it will be inserted into the table in state CONFIRMED.
 * Returns true. In all other cases returns false.
 */
func confirmTx(db kv_store.KVStore, table string, k string) bool {
	err := db.SetValue(table, k, "", CONFIRMED)
	if err != nil {
		return false
	}
	return true
}

/**
 * signTx
 *
 * If a given key exists with state CONFIRMED will update state to SIGNED.
 * Returns true. In all other cases false.
 */
func signTx(db kv_store.KVStore, table string, k string) bool {
	err := db.SetValue(table, k, CONFIRMED, SIGNED)
	if err != nil {
		return false
	}
	return true
}

/**
 * getLastBlock
 *
 * Returns the last processed block for a coin.
 * Valid is true if succeeded. False otherwise.
 */
func GetLastBlock(db kv_store.KVStore, coinName string) (int64, bool) {
	v, err := db.GetValue(LAST_BLOCK, coinName)
	if err != nil {
		return 0, false
	}
	if v == nil {
		return 0, true
	}
	i, err := strconv.ParseInt(*v, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

/**
 * setLastBlock
 *
 * Updates the last processed block. Only if new value is greater than previous value.
 * Returns true if succeeded in update.
 */
func setLastBlock(db kv_store.KVStore, coinName string, newVal int64) bool {
	prevBlock, valid := GetLastBlock(db, coinName)
	if !valid {
		return false
	}
	//QUANTA Page id is not ascending :(
	//if newVal < prevBlock {
	//    return false
	//}
	err := db.SetValue(LAST_BLOCK, coinName, strconv.FormatInt(prevBlock, 10), strconv.FormatInt(newVal, 10))
	if err != nil {
		println("Bucket is not found.")
		return false
	}
	return true
}
