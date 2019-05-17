package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"path/filepath"
	"syscall"
)

/**
 * main
 *
 * Runs the trust node
 */
func main() {
	viper.SetConfigType("yaml")

	configFile := flag.String("config", "config.yml", "configuration file")
	secretsFile := flag.String("secrets", "secrets.yml", "secrets file")
	enableRegistry := flag.Bool("registry", false, "enables registry")

	encryptFile := flag.String("encrypt", "", "encrypt file")
	encryptOutFile := flag.String("out", "config.yml.enc", "output encrypt file")

	enableSyncAddresses := flag.String("sync_addresses", "", "sync addresses")
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

		if *enableSyncAddresses != "" {
			node, success := initNode(config, secrets, false)
			if !success {
				panic("Failed to init node")
			}
			crosschainAddresses := node.rDb.GetCrosschainByBlockchain(*enableSyncAddresses)
			for _, addr := range crosschainAddresses {
				_, err := node.CreateMultisig(*enableSyncAddresses, addr.QuantaAddr)
				if err != nil {
					panic("Could not generate multisig address")
				}
			}
		} else {
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

			node := bootstrapNode(config, secrets, false)

			err = registerNode(config, node)
			if err != nil {
				panic(err)
			}
			node.run()
		}
	}
}
