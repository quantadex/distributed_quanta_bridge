package crypto

import "github.com/btcsuite/btcd/chaincfg"

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

// what happens if the public key is associated with multiple accounts?
