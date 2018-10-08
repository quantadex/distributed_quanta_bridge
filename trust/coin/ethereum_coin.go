package coin

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type EthereumCoin struct {
	client *Listener
}

func (c *EthereumCoin) Attach() error {
	client := &Listener{NetworkID: viper.GetString("ETHEREUM_NETWORK_ID")}
	ethereumClient, err := ethclient.Dial(viper.GetString("ETHEREUM_RPC"))
	if err != nil {
		return err
	}

	client.Client = ethereumClient
	return client.Start()
}

func (c *EthereumCoin) GetTopBlockID() (int64, error) {
	return c.client.GetTopBlockNumber()
}

func (c *EthereumCoin) GetDepositsInBlock(blockID int64, trustAddress string) ([]*Deposit, error) {
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

