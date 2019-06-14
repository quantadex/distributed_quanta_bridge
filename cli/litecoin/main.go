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
	config, quanta, rdb, kdb, log, secrets := cli.Setup()

	// setup coin
	blackList := crypto.GetBlackListedUsersByBlockchain(config.BlackList, coin.BLOCKCHAIN_LTC)
	coin, err := coin.NewLitecoinCoin(config.LtcRpc, crypto.GetChainCfgByStringLTC(config.LtcNetwork), secrets.LtcSigners, secrets.LtcRpcUser, secrets.LtcRpcPassword, secrets.GrapheneSeedPrefix, config.LtcWithdrawMin, config.LtcWithdrawFee, blackList)
	if err != nil {
		panic(fmt.Errorf("cannot create litecoin coin"))
	}

	err = coin.Attach()
	if err != nil {
		log.Error("Failed to attach to coin " + err.Error())
		panic(err)
	}

	time.Sleep(3 * time.Second)
	println("Starting litecoin deposit sync")

	depositSync := sync.NewLitecoinSync(coin, config.CoinMapping, quanta, kdb, rdb, log, config.LtcBlockStart, config.LtcMinConfirmation)
	depositSync.Run()
}
