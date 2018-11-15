package key_manager

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/walletdb"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/bitcoin_network"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"strconv"

	"encoding/hex"
	"encoding/json"
	"bytes"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var BTC_NETWORK = bitcoin_network.NETWORKS["btc"]

type BitcoinKeyManager struct {
	wif *btcutil.WIF
}

func (b *BitcoinKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

/**
 * Loads the keys, i.e. the bitcoin walletdb.
 * See the keystore/init-bitcoin-wallet.go CLI to iniitalize one.
 *
 * :dbPath the file path to store the wallet, e.g. /tmp/mywallet.db
 */
func (b *BitcoinKeyManager) LoadNodeKeys(dbPath string) error {
	// https://github.com/btcsuite/btcwallet
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
			return err
	}
	defer db.Close()

	bucketKey := []byte("btc.keymanager.quanta")

	// read the WIF key from the wallet db
	err = walletdb.View(db, func(tx walletdb.ReadTx) error {
	    bucket := tx.ReadBucket(bucketKey)
			if bucket == nil {
				return fmt.Errorf("No such bucket '%s' in the wallet '%s'", bucketKey, dbPath)
			}

	    // Read the key back and ensure it matches.
			byteArr := bucket.Get([]byte("wif"))
			if byteArr == nil {
				return fmt.Errorf("No such key 'wif' in the wallet '%s'", dbPath)
			}

			s := string(bucket.Get([]byte("wifLen")))
			n, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("Corrupted key 'wifLen' in the wallet '%s'", dbPath)
			}
			wif_string := string(byteArr[:n])

			wif, err := btcutil.DecodeWIF(wif_string)
			if err != nil {
				return err
			}

			// println("retrieved WIF:", wif_string, "| serializePubKey:", wif.SerializePubKey())
			b.wif = wif
			return err
	})

	return err
}

func (b *BitcoinKeyManager) GetPublicKey() (string, error) {
	return crypto.PubkeyToAddress(b.wif.PrivKey.ToECDSA().PublicKey).Hex(), nil
}

func (b *BitcoinKeyManager) GetPrivateKey() (*ecdsa.PrivateKey) {
	return b.wif.PrivKey.ToECDSA()
}

func (b *BitcoinKeyManager) SignMessage(original []byte) ([]byte, error) {
	// https://godoc.org/github.com/btcsuite/btcwallet/wallet#Wallet.SignTransaction
	panic("implement me")
}

func (b *BitcoinKeyManager) SignMessageObj(original interface{}) (*string) {
	panic("implement me")
}

func (b *BitcoinKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (b *BitcoinKeyManager) SignTransaction(encoded string) (string, error) {
	// we need the address information in the encoded string
	coin := &coin.BitcoinCoin{}

	wPtr, err := coin.DecodeRefund(encoded)
	if err != nil {
		return "", err
	}

	// assert.NotNil(b.wif)

	fmt.Printf("wif=%s, amount=%d, dest=%s\n", b.wif, int64(wPtr.Amount), wPtr.DestinationAddress)

	txHash, err := FindPreviousTxHash(wPtr)

	tx, err := CreateTransaction(
		b.wif,
		wPtr.DestinationAddress,
		int64(wPtr.Amount),
		txHash)

	if err != nil {
		return "", err
	}

	// FIXME: maybe just return (Transaction)tx.signedtx?
	data, _ := json.Marshal(tx)
	fmt.Println(string(data))

	return string(data), nil
}

func FindPreviousTxHash(wPtr *coin.Withdrawal) (string, error) {
	txHash := "81b4c832d70cb56ff957589752eb4125a4cab78a25a8fc52d6a09e5bd4404d48"
	// FIXME: bitcoin requires us to chain in the previous hash...
  // so do we need to make this part of the Withdrawal object or what?
	return txHash, nil
}

func (b *BitcoinKeyManager) VerifyTransaction(encoded string) (bool, error) {
	panic("implement me")
}


type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

// TODO: maybe use https://github.com/soroushjp/go-bitcoin-multisig/blob/master/multisig/address.go:generateAddress
/** Creates a signed Bitcoin Transaction. This is a pay to address script,
 * need to modify somehow as a multisig chain...
 *
 *   >>> transaction, err := CreateTransaction("5HusYj2b2x4nroApgfvaSfKYZhRbKFH41bVyPooymbC6KfgSXdD", "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa", 91234, "81b4c832d70cb56ff957589752eb4125a4cab78a25a8fc52d6a09e5bd4404d48")
 *   >>> data, _ := json.Marshal(transaction)
 *   >>> fmt.Println(string(data))
 *	 {
 * 		    "txid": "4e8378675bcf6a389c8cfe246094aafa44249e48ab88a40e6fda3bf0f44f916a",
 *		    "source_address": "1MMMMSUb1piy2ufrSguNUdFmAcvqrQF8M5",
 *		    "destination_address": "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
 *			    "amount": 91234,
 *			    "unsignedtx": "0100000001484d40d45b9ea0d652fca8258ab7caa42541eb52975857f96fb50cd732c8b4810000000000ffffffff0162640100000000001976a914df3bd30160e6c6145baaf2c88a8844c13a00d1d588ac00000000",
 *			    "signedtx": "01000000016a914ff4f03bda6f0ea488ab489e2444faaa946024fe8c9c386acf5b6778834e000000008b483045022100904dbeddeecccf6391ac92922381ae006bf244c002f42e195daa0a9837a4b5820220461677f9dbb7d9580e268ac486cfeb4b9d87bfdd6d4e7b1be09b8e6f5cc0a70701410414e301b2328f17442c0b8310d787bf3d8a404cfbd0704f135b6ad4b2d3ee751310f981926e53a6e8c39bd7d3fefd576c543cce493cbac06388f2651d1aacbfcdffffffff0162640100000000001976a914c8e90996c7c6080ee06284600c684ed904d14c5c88ac00000000"
 *		}
 * see https://www.thepolyglotdeveloper.com/2018/03/create-sign-bitcoin-transactions-golang/
 */
func CreateTransaction(wif *btcutil.WIF, destination string, amount int64, prevTxHash string) (Transaction, error) {
	var transaction Transaction
	addresspubkey, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.MainNetParams)
	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceUtxoHash, _ := chainhash.NewHashFromStr(prevTxHash)
	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)

	destinationAddress, err := btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)

	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), &chaincfg.MainNetParams)
	if err != nil {
		return Transaction{}, err
	}

	// use PayToAddrScript for a traditional single address destination
	// destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)

	// use MultiSigScript for Quanta
	// FIXME: this is a single node multisig script... need to collect the other nodes
	signerPubKeys := make([]*btcutil.AddressPubKey, 1)
	signerPubKeys[0] = addresspubkey
	nrequired := 1
	destinationPkScript, _ := txscript.MultiSigScript(signerPubKeys, nrequired)

	// assumes the source was a traditional pay to address script
	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&sourceTxHash, 0)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	// FIXME: this probably needs to be adjusted when there is more than one signer in the multi sig...
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return Transaction{}, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return Transaction{}, err
	}
	if err := vm.Execute(); err != nil {
		return Transaction{}, err
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	return transaction, nil
}
