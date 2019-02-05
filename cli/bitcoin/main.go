package main

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/btcsuite/btcd/chaincfg"
	"time"
)

func main() {
	config, quanta, rdb, kdb, log := cli.Setup()

	// setup coin
	coin, err := coin.NewBitcoinCoin(&chaincfg.RegressionNetParams,config.BtcSigners)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	fmt.Printf("BTC signers %v\n", config.BtcSigners)
	time.Sleep(3 * time.Second)
	println("Starting bitcoin deposit sync")

	depositSync := sync.NewBitcoinSync(coin, coin.Blockchain(), quanta, kdb,rdb, log, config.BtcBlockStart)
	depositSync.Run()
}