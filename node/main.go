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
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"encoding/json"
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
	bounceTx := flag.String("bounce", "", "bounce tx")
	retryTx := flag.String("retry", "", "retry tx")
	repair := flag.Bool("repair", false, "repair")

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
			fmt.Println("Synchronize addresses")
			crosschainAddresses := node.rDb.GetCrosschainByBlockchain(*enableSyncAddresses)
			for _, addr := range crosschainAddresses {
				fmt.Println("process ", addr.Address, addr.Blockchain, addr.QuantaAddr)
				_, err := node.CreateMultisig(*enableSyncAddresses, addr.QuantaAddr)
				if err != nil {
					panic("Could not generate multisig address")
				}
			}
		} else if (*bounceTx != "") {
			node, success := initNode(config, secrets, false)
			if !success {
				panic("Failed to init node")
			}
			fmt.Println("Bounce TX=", *bounceTx)
			tx, err := db.GetTransaction(node.rDb, *bounceTx)
			if err != nil {
				panic(fmt.Errorf("failed to get tx file: %s \n", err))
			}

			if tx.Type == db.DEPOSIT {
				fmt.Println("Can't bounce deposit")
				return
			}

			if tx.SubmitState != db.SUBMIT_SUCCESS {
				fmt.Println("marking as a bounce")
				refund := quanta.Refund{
					TransactionId: tx.Tx,
					CoinName: tx.Coin,
					LedgerID: int32(tx.BlockId),
					Amount: uint64(tx.Amount),
					SourceAddress: tx.From,
					DestinationAddress: tx.To,
					BlockHash: tx.BlockHash,
				}
				node.qTC.BounceTx(&refund, db.AMOUNT_TOO_SMALL, true)
			} else {
				fmt.Println("Tx already processed successfully.")
			}
		} else if (*retryTx != "") {
			node, success := initNode(config, secrets, false)
			if !success {
				panic("Failed to init node")
			}
			fmt.Println("Retry TX=", *bounceTx)
			tx, err := db.GetTransaction(node.rDb, *bounceTx)
			if err != nil {
				panic(fmt.Errorf("failed to get tx file: %s \n", err))
			}
			if tx.SubmitState != db.SUBMIT_SUCCESS {
				fmt.Println("marking as a retry")
				if tx.Type == db.DEPOSIT {
					db.ChangeDepositSubmitState(node.rDb, tx.Tx, db.SUBMIT_CONSENSUS, tx.SubmitConfirm_block, tx.SubmitTxHash, tx.BlockHash)
				} else {
					db.ChangeWithdrawalSubmitState(node.rDb, tx.Tx, db.SUBMIT_CONSENSUS, tx.TxId, tx.SubmitTxHash, tx.BlockHash)
				}
			}
		} else if *repair {
			node, success := initNode(config, secrets, false)
			if !success {
				panic("Failed to init node")
			}

			for _, blockchain := range []string{coin.BLOCKCHAIN_BTC, coin.BLOCKCHAIN_BCH, coin.BLOCKCHAIN_LTC, coin.BLOCKCHAIN_ETH} {
				fmt.Println("Repairing ", blockchain)
				req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/address/%s",config.RegistrarIp, config.ExternalListenPort, blockchain), nil)
				req.Header.Set("User-Agent", "AmazonAPIGateway_wya99cec1d")
				client := &http.Client{}
				res, _ := client.Do(req)

				if err != nil {
					panic(err.Error())
				}
				bodyBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err.Error())
				}

				var addresses []db.CrosschainAddress
				err = json.Unmarshal(bodyBytes, &addresses)
				if err != nil {
					panic(err.Error())
				}
				for _, addr := range addresses {
					fmt.Println("process ", addr.Address, addr.Blockchain, addr.QuantaAddr)
					addr, err := node.CreateMultisig(addr.Blockchain, addr.QuantaAddr)
					if err != nil {
						panic("Could not generate multisig address: " +  err.Error())
					}

					err = node.rDb.AddCrosschainAddress(addr)
					if err != nil {
					} else {
						println("Added ", addr.ContractAddress, " to db.")
					}
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
