package main

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	config, quanta, rdb, kdb, log := cli.Setup()

	// setup coin
	coin, err := coin.NewBitcoinCoin(&chaincfg.RegressionNetParams,nil)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	println("Starting bitcoin deposit sync")
	depositSync := sync.NewBitcoinSync(coin, config.CoinName, quanta, kdb,rdb, log, config.BtcBlockStart)
	depositSync.Run()
}