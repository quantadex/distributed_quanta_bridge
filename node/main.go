package main

import (
	"flag"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
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
	config := Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	enableRegistry := flag.Bool("registry", false, "enables registry")
	flag.Parse()

	if *enableRegistry {
		// start registrar if we need to
		logger, _ := logger.NewLogger("registrar")
		registrarUrl := fmt.Sprintf("%s:%d", config.RegistrarIp, config.RegistrarPort)
		s := service.NewServer(service.NewRegistry(), registrarUrl, logger)
		s.DoHealthCheck(5)
		go s.Start()
	}

	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	node := bootstrapNode(config, coin)
	node.run()
}
