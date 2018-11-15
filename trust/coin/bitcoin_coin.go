package coin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/big"
	// "io/ioutil"
	// "path/filepath"
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/btcsuite/btcd/rpcclient"

	// "github.com/btcsuite/btcd/wire"
	// "github.com/btcsuite/btcutil"
)

type BitcoinCoin struct {
	client *rpcclient.Client
}

func (b *BitcoinCoin) Attach() error {
	_, err := b.GetClient("localhost:4444", "test", "test")

	return err
}

func(b *BitcoinCoin) Detach() error {
	if (b.client != nil) {
		b.client.Shutdown()
		b.client.WaitForShutdown()
		b.client = nil
	}
	return nil
}

/**
 * :rpchost e.g. 'localhost:8334'
 * :rpcuser e.g. 'myrpcuser'
 * :rpcpass e.g. 'myrpcpass'
 */
func (b *BitcoinCoin) GetClient(rpchost string, rpcuser string, rpcpass string) (*rpcclient.Client, error) {
	// based on https://github.com/btcsuite/btcd/blob/master/rpcclient/examples/btcdwebsockets/main.go
	// https://godoc.org/github.com/btcsuite/btcd/rpcclient
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	// ntfnHandlers := rpcclient.NotificationHandlers{
	// 	OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
	// 		log.Printf("Block connected: %v (%d) %v",
	// 			header.BlockHash(), height, header.Timestamp)
	// 	},
	// 	OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
	// 		log.Printf("Block disconnected: %v (%d) %v",
	// 			header.BlockHash(), height, header.Timestamp)
	// 	},
	// }

	// Connect to local btcd RPC server using websockets.
	// btcdHomeDir := btcutil.AppDataDir("btcd", false)
	// certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	// if err != nil {
	// 	return nil, err
	// }
	connCfg := &rpcclient.ConnConfig{
		Host:         rpchost,
		// Endpoint:     "ws",
		User:         rpcuser,
		Pass:         rpcpass,
		HTTPPostMode: true,
		DisableTLS:   true,
		// Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, nil)  // &ntfnHandlers)
	if err != nil {
		return nil, err
	}

	// Register for block connect and disconnect notifications.
	// if err := client.NotifyBlocks(); err != nil {
	// 	return nil, err
	// }
	// log.Println("NotifyBlocks: Registration Complete")

	blockCount, err := client.GetBlockCount()
	if err != nil {
		return nil, err
	}
	log.Printf("BitCoin attached, block count: %d", blockCount)

	b.client = client
	return client, nil
}

func (b *BitcoinCoin) GetTopBlockID() (int64, error) {
	blockCount, err := b.client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockCount, nil
}

func (b *BitcoinCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	panic("not supported")
}

func (b *BitcoinCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	panic("implement me")
}

func (b *BitcoinCoin) GetForwardersInBlock(blockID int64) ([]*ForwardInput, error) {
	panic("implement me")
}

func (b *BitcoinCoin) SendWithdrawal(trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {
	panic("implement me")
}

func (b *BitcoinCoin) EncodeRefund(w Withdrawal) (string, error) {
	var encoded bytes.Buffer

	binary.Write(&encoded, binary.BigEndian, uint64(w.TxId))

	binary.Write(&encoded, binary.BigEndian, uint64(len(w.DestinationAddress)))
	encoded.Write([]byte(w.DestinationAddress))

	binary.Write(&encoded, binary.BigEndian, uint64(w.Amount))

	data, err := json.Marshal(&EncodedMsg{ common2.Bytes2Hex(encoded.Bytes()), w.QuantaBlockID})
	return common2.Bytes2Hex(data), err
}

func (b *BitcoinCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	decoded := common2.Hex2Bytes(encoded)
	msg := &EncodedMsg{}
	err := json.Unmarshal(decoded, msg)
	if err != nil {
		return nil, err
	}

	w := &Withdrawal{}
	w.QuantaBlockID = msg.BlockNumber
	decoded = common2.Hex2Bytes(msg.Message)

	// skip 4 bytes to length
	pl := 0

	bytesBuf := decoded[pl : pl+8]
	txId := new(big.Int).SetBytes(bytesBuf).Uint64()
	pl += 8

	// number of bytes header of the destination address
	bytesBuf = decoded[pl : pl+8]
	n := int(new(big.Int).SetBytes(bytesBuf).Uint64())
	pl += 8
	// get exactly n bytes (the destination address)
	destAddress := string(decoded[pl : pl+n])
	pl += n

	bytesBuf = decoded[pl : pl+8]
	amount := new(big.Int).SetBytes(bytesBuf).Uint64()
	pl += 8

	w.CoinName = "BTC"
	w.TxId = txId
	w.DestinationAddress = destAddress
	w.Amount = amount

	return w, nil
}
