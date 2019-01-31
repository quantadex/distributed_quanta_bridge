package common

import "github.com/btcsuite/btcd/btcjson"

type TransactionBitcoin struct {
	Tx string
	RawInput []btcjson.RawTxInput
}
