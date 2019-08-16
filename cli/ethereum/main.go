package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
)

func main() {
	config, quanta, rdb, kdb, log, secrets := cli.Setup()

	// setup coin
	blackList := crypto.GetBlackListedUsersByBlockchain(config.BlackList, coin.BLOCKCHAIN_ETH)
	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc, secrets.EthereumKeyStore, config.Erc20Mapping, config.EthWithdrawMin, config.EthWithdrawFee, config.EthWithdrawGasFee, blackList)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	println("Starting ethereum deposit sync")
	depositSync := sync.NewEthereumSync(coin, config.EthereumTrustAddr, config.CoinMapping, quanta, kdb, rdb, log, config.EthereumBlockStart, config.EthFlush, config.EthMinConfirmation)
	depositSync.Run()
}
