package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
)

/**
 * main
 *
 * Runs the trust node
 */
func main() {
	viper.SetConfigType("yaml")
	configFile := flag.String("config", "config.yml", "configuration file")

	path, err := filepath.Abs(filepath.Dir(*configFile))
	if err != nil {
		panic(fmt.Errorf("Could not find file path: %s \n", err))
	}

	enableRegistry := flag.Bool("registry", false, "enables registry")
	portNumber := flag.Int("port", 0, "overrides port")
	flag.Parse()

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	err = viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := common.Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if *enableRegistry {
		// start registrar if we need to
		logger, _ := logger.NewLogger("registrar")
		registrarUrl := fmt.Sprintf(":%d", config.RegistrarPort)
		s := service.NewServer(service.NewRegistry(config.MinNodes, path), registrarUrl, logger)
		s.DoHealthCheck(5)
		go s.Start()
	}

	if *portNumber != 0 {
		config.ListenPort = *portNumber
	}

	coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc, config.EthereumKeyStore, config.Erc20Mapping)
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	node := bootstrapNode(config, coin)
	err = registerNode(config, node)
	if err != nil {
		panic(err)
	}

	node.run()
}
