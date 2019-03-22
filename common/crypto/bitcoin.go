package crypto

import (
	chaincfg3 "github.com/bchsuite/bchd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg"
	chaincfg2 "github.com/ltcsuite/ltcd/chaincfg"
)

/**
 * Checks whether the set of addresses match with the
 * list of addresses in the multisig
 *
 * Input: pksh
 */
func CheckKeysInMultisig() (bool, error) {
	panic("")
}

func BTCMultisig2QUANTAPubKey() {

}

func GetChainCfgByString(network string) *chaincfg.Params {
	if network == "regnet" {
		return &chaincfg.RegressionNetParams
	} else if network == "testnet" {
		return &chaincfg.TestNet3Params
	} else if network == "mainnet" {
		return &chaincfg.MainNetParams
	} else {
		panic("invalid chain cfg")
	}
}

func GetChainCfgByStringLTC(network string) *chaincfg2.Params {
	if network == "regnet" {
		return &chaincfg2.RegressionNetParams
	} else if network == "testnet" {
		return &chaincfg2.TestNet4Params
	} else if network == "mainnet" {
		return &chaincfg2.MainNetParams
	} else {
		panic("invalid chain cfg")
	}
}

func GetChainCfgByStringBCH(network string) *chaincfg3.Params {
	if network == "regnet" {
		return &chaincfg3.RegressionNetParams
	} else if network == "testnet" {
		return &chaincfg3.TestNet3Params
	} else if network == "mainnet" {
		return &chaincfg3.MainNetParams
	} else {
		panic("invalid chain cfg")
	}
}

// what happens if the public key is associated with multiple accounts?
