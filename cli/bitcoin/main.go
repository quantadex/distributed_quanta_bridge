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
	config, quanta, rdb, kdb, log, secrets := cli.Setup(false)

	// setup coin
	blackList := crypto.GetBlackListedUsersByBlockchain(config.BlackList, coin.BLOCKCHAIN_BTC)
	coin, err := coin.NewBitcoinCoin(config.BtcRpc, crypto.GetChainCfgByString(config.BtcNetwork), secrets.BtcSigners, secrets.BtcRpcUser, secrets.BtcRpcPassword, secrets.GrapheneSeedPrefix, config.BtcWithdrawMin, config.BtcWithdrawFee, blackList)
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
