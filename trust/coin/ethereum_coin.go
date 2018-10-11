package coin

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"encoding/base64"
)

const paymentTx = "paymentTx"

const abiJson = `
[
	{ "type" : "function", "name" : "paymentTx",  "constant" : false, "inputs" : [ { "name" : "tx_id", "type" : "uint64" }, { "name" : "erc20Address", "type" : "address" }, { "name" : "to", "type" : "address" }, { "name" : "amount", "type" : "uint256" } ] },
]`


type EthereumCoin struct {
	client *Listener
	maxRange int64
	networkId string
	ethereumRpc string
	abi abi.ABI
}

func (c *EthereumCoin) Attach() error {
	c.client = &Listener{NetworkID: c.networkId}
	ethereumClient, err := ethclient.Dial(c.ethereumRpc)
	if err != nil {
		return err
	}

	c.abi, err = abi.JSON(strings.NewReader(abiJson))
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
	data, err := c.abi.Pack(paymentTx, w.NodeID, common2.HexToAddress(w.DestinationAddress), w.Amount)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func (c *EthereumCoin)  DecodeRefund(encoded string) (*Withdrawal, error) {
	return nil, nil
}
