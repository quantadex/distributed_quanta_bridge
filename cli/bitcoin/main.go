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
	coin, err := coin.NewBitcoinCoin(config.BtcRpc, crypto.GetChainCfgByString(config.BtcNetwork), config.BtcSigners)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	time.Sleep(3 * time.Second)
	println("Starting bitcoin deposit sync")

	depositSync := sync.NewBitcoinSync(coin, config.CoinMapping, quanta, kdb, rdb, log, config.BtcBlockStart, config.BtcMinConfirmation)
	depositSync.Run()
}
