package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"time"
)

func main() {
	config, quanta, rdb, kdb, log := cli.Setup()

	// setup coin
	coin, err := coin.NewBCHCoin(config.BchRpc, crypto.GetChainCfgByStringBCH(config.BchNetwork), config.BchSigners)
	if err != nil {
		panic(fmt.Errorf("cannot create bch coin"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	time.Sleep(3 * time.Second)
	println("Starting bch deposit sync")

	depositSync := sync.NewBCHSync(coin, config.CoinMapping, quanta, kdb, rdb, log, config.BchBlockStart, config.BchMinConfirmation)
	depositSync.Run()
}
