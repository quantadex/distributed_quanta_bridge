package main

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
)

/**
 * main
 *
 * Runs the trust node
 */
func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("node")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	config := Config {}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	node := bootstrapNode(config, coin)
	node.run()
}
