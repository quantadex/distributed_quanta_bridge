package main

import (
  "flag"
  "fmt"
  "os"
	"github.com/btcsuite/btcutil"
  "github.com/btcsuite/btcwallet/walletdb"
  _ "github.com/btcsuite/btcwallet/walletdb/bdb"
  "github.com/quantadex/distributed_quanta_bridge/trust/bitcoin_network"
)

var BTC_NETWORK = bitcoin_network.NETWORKS["btc"]

/** Initializes a Quanta Trust node BitCoin Wallet database
 * with an optional bitcoin WIF.
 *
 * This is used by the BitcoinKeyManager.
  *
 * The database contains a single bucket 'btc.keymanager.quanta',
 * with a key 'wif' containing the bitcoin WIF string (as bytes)
 * corresponding to the bitcoin address.
 * And a second key 'wifLen', indicating the number of characters of the WIF string
 *
 * Output: a json containing the created wif and address.
 *
 * Usage:
 *
 *   go build init-bitcoin-wallet.go
 *   ./init-bitcoin-wallet filename.db WIF
 *
 *   (or you may also `go run ./init-bitcoin-wallet.go ...`)
 *
 * where:
 *   filename.db - the name of the file to create
 *   WIF - the bitcoin WIF as a string, optional. If not specified, a random one is generated.
 *
 * e.g.
 *
 *   go run init-bitcoin-wallet.go foo.db
 *   go run init-bitcoin-wallet.go foo.db Kyunn4aPBkC8Xw82oRBTcomzTgM9rmQkt8i5ev3A1XkUwVL3RaT6
 */
func main() {
  flag.Parse()

  args := flag.Args()
  if len(args) == 0 || len(args) > 2 {
    fmt.Fprintf(os.Stderr, "Usage: %s [--test] filename.db WIF\n", os.Args[0])
    os.Exit(1)
  }
  filename := args[0]

  var wif *btcutil.WIF
  if len(args) == 2 {
    wifString := args[1]
    wifImport, err := BTC_NETWORK.ImportWIF(wifString)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not import WIF: %s\n", err)
      os.Exit(2)
    }
    wif = wifImport
  } else {
    // generate a new one
    wifNew, err := BTC_NETWORK.CreatePrivateKey()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not create WIF: %s\n", err)
      os.Exit(3)
    }
    wif = wifNew
  }

  if _, err := os.Stat(filename); !os.IsNotExist(err) {
    fmt.Fprintf(os.Stderr, "Wallet '%s' already exists\n", filename)
    os.Exit(4)
  }

  err := initWallet(wif, filename)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not create wallet '%s': %s\n", filename, err)
    os.Exit(5)
  }

  address, _ := BTC_NETWORK.GetAddress(wif)
  fmt.Printf("{\n  \"wif\": \"%s\",\n  \"address\": \"%s\"\n}\n", wif.String(), address.EncodeAddress())
}

func initWallet(wif *btcutil.WIF, filename string) error {
	// https://github.com/btcsuite/btcwallet
	db, err := walletdb.Create("bdb", filename)
  if err != nil {
    return err
  }

  defer db.Close()

  bucketName := "btc.keymanager.quanta"
	bucketKey := []byte(bucketName)
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
    bucket := tx.ReadWriteBucket(bucketKey)
    if bucket == nil {
      _, err = tx.CreateTopLevelBucket(bucketKey)
      if err != nil {
          return err
      }
      bucket = tx.ReadWriteBucket(bucketKey)
    }

    key := []byte("wif")
    value := []byte(wif.String())
    if err := bucket.Put(key, value); err != nil {
        return err
    }

    key = []byte("wifLen")
    value = []byte(fmt.Sprintf("%d", len(wif.String())))
    if err := bucket.Put(key, value); err != nil {
        return err
    }

    return nil
  })

  return err
}
