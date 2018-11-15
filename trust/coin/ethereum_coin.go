package coin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"math/big"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"encoding/json"
)

const sign_prefix = "\x19Ethereum Signed Message:\n"

type EthereumCoin struct {
	client      *Listener
	maxRange    int64
	networkId   string
	ethereumRpc string
}

type EncodedMsg struct {
	Message string
	Tx string
	BlockNumber int64
}

func (c *EthereumCoin) Attach() error {
	c.client = &Listener{NetworkID: c.networkId}
	ethereumClient, err := ethclient.Dial(c.ethereumRpc)
	if err != nil {
		return err
	}

	c.client.Client = ethereumClient
	err = c.client.Start()
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumCoin) GetTopBlockID() (int64, error) {
	topBlockId, err := c.client.GetTopBlockNumber()
	if err != nil {
		return 0, err
	}

	return common.Min64(c.maxRange, topBlockId), nil
}

func (c *EthereumCoin) GetTxID(trustAddress common2.Address) (uint64, error) {
	return c.client.GetTxID(nil, trustAddress)
}

func (c *EthereumCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	ndeposits, err := c.client.GetNativeDeposits(blockID, trustAddress)
	if err != nil {
		return nil, err
	}
	deps, err := c.client.FilterTransferEvent(blockID, trustAddress)
	if err != nil {
		return nil, err
	}

	return append(ndeposits, deps...), nil
}

func (c *EthereumCoin) GetForwardersInBlock(blockID int64) ([]*ForwardInput, error) {
	forwarders, err := c.client.GetForwardContract(blockID)
	if err != nil {
		return nil, err
	}
	return forwarders, nil
}

//
func (c *EthereumCoin) SendWithdrawal(trustAddress common2.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {
	return c.client.SendWithDrawalToRPC(trustAddress, ownerKey, w)
}

func (c *EthereumCoin) EncodeRefund(w Withdrawal) (string, error) {
	var encoded bytes.Buffer
	var smartAddress string
	parts := strings.Split(w.CoinName, ",")

	if len(parts) == 2 {
		smartAddress = parts[1]
	} else {
		smartAddress = ""
	}

	var number = common2.Big256
	number.SetUint64(w.Amount)
	//encoded.WriteString(sign_prefix + "80")
	binary.Write(&encoded, binary.BigEndian, uint64(w.TxId))
	encoded.Write(common2.HexToAddress(strings.ToLower(smartAddress)).Bytes())
	encoded.Write(common2.HexToAddress(strings.ToLower(w.DestinationAddress)).Bytes())
	encoded.Write(abi.U256(new(big.Int).SetUint64(uint64(w.Amount))))
	//binary.Write(&encoded, binary.BigEndian, abi.U256(new(big.Int).SetUint64(uint64(w.Amount))))

	//println("# of bytes " , encoded.Len(), common2.Bytes2Hex(encoded.Bytes()))
	data, err := json.Marshal(&EncodedMsg{ common2.Bytes2Hex(encoded.Bytes()),  w.Tx, w.QuantaBlockID})
	return common2.Bytes2Hex(data), err
}

func (c *EthereumCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	decoded := common2.Hex2Bytes(encoded)
	msg := &EncodedMsg{}
	err := json.Unmarshal(decoded, msg)
	if err != nil {
		return nil, err
	}

	w := &Withdrawal{}
	w.Tx = msg.Tx
	w.QuantaBlockID = msg.BlockNumber
	decoded = common2.Hex2Bytes(msg.Message)

	//pl := len(sign_prefix)
	//header := decoded[0:pl]
	//if string(header) != sign_prefix {
	//	return nil, errors.New("Unexpected prefix")
	//}

	// skip 4 bytes to length
	pl := 0
	txIdBytes := decoded[pl : pl+8]
	txId := new(big.Int).SetBytes(txIdBytes).Uint64()
	w.TxId = txId

	pl += 8
	smartAddress := decoded[pl : pl+20]

	pl += 20
	destAddress := decoded[pl : pl+20]

	pl += 20
	amount := decoded[pl : pl+32]
	smartNumber := new(big.Int).SetBytes(smartAddress)

	if smartNumber.Cmp(big.NewInt(0)) == 0 {
		w.CoinName = "ETH"
	} else {
		w.CoinName = "," + strings.ToLower(common2.BytesToAddress(smartAddress).Hex())
	}

	w.DestinationAddress = strings.ToLower(common2.BytesToAddress(destAddress).Hex())
	w.Amount = new(big.Int).SetBytes(amount).Uint64()

	return w, nil
}
