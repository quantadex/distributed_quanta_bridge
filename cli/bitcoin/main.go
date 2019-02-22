package main

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/cli"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control/sync"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
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

	fmt.Printf("BTC signers %v\n", config.BtcSigners)
	time.Sleep(3 * time.Second)
	println("Starting bitcoin deposit sync")

	depositSync := sync.NewBitcoinSync(coin, config.CoinMapping, quanta, kdb, rdb, log, config.BtcBlockStart)
	depositSync.Run()
}
