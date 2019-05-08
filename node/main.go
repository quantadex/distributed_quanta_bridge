package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
)

/**
 * main
 *
 * Runs the trust node
 */
func main() {

	configFile := flag.String("config", "config.yml", "configuration file")
	secretsFile := flag.String("secrets", "secrets.yml", "secrets file")
	enableRegistry := flag.Bool("registry", false, "enables registry")

	encryptFile := flag.String("encrypt", "", "encrypt file")
	encryptOutFile := flag.String("out", "config.yml.enc", "output encrypt file")
	flag.Parse()

	if *encryptFile != "" {
		data, err := ioutil.ReadFile(*encryptFile)
		if err != nil {
			panic(err)
		}

		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		err = crypto.EncryptSecretsFile(string(password), data, *encryptOutFile)
		if err != nil {
			panic(err)
		}

	} else {
		flag.Parse()
		viper.SetConfigType("yaml")

		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		secrets, err := crypto.DecryptSecretsFile(*secretsFile, string(password))

		path, err := filepath.Abs(filepath.Dir(*configFile))
		if err != nil {
			panic(fmt.Errorf("Could not find file path: %s \n", err))
		}

		portNumber := flag.Int("port", 0, "overrides port")

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

		coin, err := coin.NewEthereumCoin(config.EthereumNetworkId, config.EthereumRpc, secrets.EthereumKeyStore, config.Erc20Mapping)
		if err != nil {
			panic(fmt.Errorf("cannot create ethereum listener"))
		}

		node := bootstrapNode(config, coin, secrets)
		err = registerNode(config, node)
		if err != nil {
			panic(err)
		}

		node.run()
	}

}
