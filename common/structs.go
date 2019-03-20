package common

import (
	btcjson3 "github.com/bchsuite/bchd/bchjson"
	"github.com/btcsuite/btcd/btcjson"
	btcjson2 "github.com/ltcsuite/ltcd/btcjson"
)

type TransactionBitcoin struct {
	Tx       string
	RawInput []btcjson.RawTxInput
}

type TransactionLitecoin struct {
	Tx       string
	RawInput []btcjson2.RawTxInput
}

type TransactionBCH struct {
	Tx       string
	RawInput []btcjson3.RawTxInput
}
