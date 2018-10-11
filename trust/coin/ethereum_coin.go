package coin

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common"
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

