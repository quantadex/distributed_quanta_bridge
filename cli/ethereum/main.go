package main

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"github.com/quantadex/distributed_quanta_bridge/cli"
)

func main() {
	config, quanta, rdb, kdb, log := cli.Setup()

	// setup coin
	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	println("Starting ethereum deposit sync")
	depositSync := sync.NewEthereumSync(coin, config.EthereumTrustAddr, config.CoinName, quanta, kdb,rdb, log, config.EthereumBlockStart)
	depositSync.Run()
}