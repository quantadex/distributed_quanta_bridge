package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
)

func main() {
	config, quanta, rdb, kdb, log := cli.Setup()

	// setup coin
	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc, config.EthereumKeyStore)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	println("Starting ethereum deposit sync")
	depositSync := sync.NewEthereumSync(coin, config.EthereumTrustAddr, config.CoinMapping, quanta, kdb, rdb, log, config.EthereumBlockStart, config.EthFlush)
	depositSync.Run()
}
