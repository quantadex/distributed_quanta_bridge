package coin

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"strings"
	common2 "github.com/ethereum/go-ethereum/common"
	"encoding/base64"
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
)

type EthereumCoin struct {
	client *Listener
	maxRange int64
	networkId string
	ethereumRpc string
}

func (c *EthereumCoin) Attach() error {
	c.client = &Listener{NetworkID: c.networkId}
	ethereumClient, err := ethclient.Dial(c.ethereumRpc)
	if err != nil {
		return err
	}

	c.client.Client = ethereumClient
	return c.client.Start()
}

func (c *EthereumCoin) GetTopBlockID() (int64, error) {
	topBlockId, err := c.client.GetTopBlockNumber()
	if err != nil {
		return 0, err
	}

	return common.Min64(c.maxRange, topBlockId), nil
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

func (c *EthereumCoin) SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error {
	panic("implement me")
}

func (c *EthereumCoin) EncodeRefund(w Withdrawal) (string, error) {
	var encoded bytes.Buffer
	var smartAddress string
	parts := strings.Split(w.CoinName,",")

	if len(parts) == 2 {
		smartAddress = parts[1]
	} else {
		smartAddress = ""
	}

	var number = common2.Big256
	number.SetInt64(w.Amount)
	encoded.Write(common2.HexToAddress(strings.ToLower(smartAddress)).Bytes())
	encoded.Write(common2.HexToAddress(strings.ToLower(w.DestinationAddress)).Bytes())
	encoded.Write(abi.U256(new(big.Int).SetUint64(uint64(w.Amount))))

	return base64.StdEncoding.EncodeToString(encoded.Bytes()), nil
}

func (c *EthereumCoin)  DecodeRefund(encoded string) (*Withdrawal, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	w := &Withdrawal{}
	println(len(decoded))

	smartAddress := decoded[0:20]
	destAddress := decoded[20:40]
	amount := decoded[40: 40 + 32]
	smartNumber := new(big.Int).SetBytes(smartAddress)

	if smartNumber.Cmp(big.NewInt(0)) == 0 {
		w.CoinName = "ETH"
	} else {
		w.CoinName = "," + strings.ToLower(common2.BytesToAddress(smartAddress).Hex())
	}

	w.DestinationAddress = strings.ToLower(common2.BytesToAddress(destAddress).Hex())
	w.Amount = new(big.Int).SetBytes(amount).Int64()

	return w, nil
}
