package cli

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/op/go-logging"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"strconv"
	"syscall"
)

func Setup(local_km bool) (*common.Config, quanta.Quanta, *db.DB, kv_store.KVStore, logger.Logger, *common.Secrets) {
	viper.SetConfigType("yaml")
	var secretsFile *string
	var secrets common.Secrets
	configFile := flag.String("config", "config.yml", "configuration file")

	if local_km {
		secretsFile = flag.String("secrets", "secrets.yml", "secrets file")
	}
	flag.Parse()

	viper.SetConfigType("yaml")

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

	if local_km {
		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		secrets, err = crypto.DecryptSecretsFile(*secretsFile, string(password))

		if err != nil {
			panic(fmt.Errorf("Fatal error secrets file: %s \n", err))
		}
	} else {
		km, err := key_manager.NewRemoteKeyManager(coin.BLOCKCHAIN_QUANTA, fmt.Sprintf("%s:%d", config.KmIp, config.KmPort))
		if err != nil {
			panic(fmt.Errorf("Fatal error secrets file 1: %s \n", err))
		}
		secrets, err = km.GetSecretsWithoutKeys()
		if err != nil {
			panic(fmt.Errorf("Fatal error secrets file 2: %s \n", err))
		}
	}

	// setup logger
	log, err := logger.NewLogger(strconv.Itoa(config.ListenPort))

	if err != nil {
		panic(err)
	}

	level, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		fmt.Println("Log level not parsed")
		level = logging.INFO
	}
	log.SetLogLevel(level)

	// setup kv
	needsInitialize := !kv_store.DbExists(config.KvDbName)

	kdb, err := kv_store.NewKVPGStore()
	if err != nil {
		log.Error("Failed to create database")
		panic(err)
	}

	err = kdb.Connect(secrets.DatabaseUrl)
	if err != nil {
		log.Error("Failed to connect to database")
		panic(err)
	}

	if needsInitialize {
		log.Info("Initialize ledger")
		control.InitLedger(kdb)
	}

	// setup db
	rDb := &db.DB{}
	info, err := pg.ParseURL(secrets.DatabaseUrl)
	if err != nil {
		log.Error(err.Error())
	}
	//node.rDb.Debug()
	rDb.Connect(info.Addr, info.User, info.Password, info.Database)
	db.MigrateTx(rDb)
	db.MigrateKv(rDb)
	db.MigrateXC(rDb)
	db.MigrateW(rDb)
	db.MigrateFM(rDb)

	quanta, err := quanta.NewQuantaGraphene(quanta.QuantaClientOptions{
		log,
		rDb,
		config.ChainId,
		config.IssuerAddress,
		config.NetworkUrl,
	})

	quanta.AttachQueue(kdb)
	if err != nil {
		log.Error("Failed to create quanta")
		panic(err)
	}
	err = quanta.Attach()
	if err != nil {
		log.Error("Failed to connect to quanta")
		panic(err)
	}

	return &config, quanta, rDb, kdb, log, &secrets
}

func RunSigner(config *common.Config, secrets *common.Secrets, logger logger.Logger, port int) error {
	var err error

	quanta, err := key_manager.NewGrapheneKeyManager(config.ChainId)
	if err != nil {
		logger.Error("Failed to set up quanta keys")
		return err
	}
	err = quanta.LoadNodeKeys(secrets.NodeKey)
	if err != nil {
		logger.Error("Failed to set up node keys")
		return err
	}

	coinkM, err := key_manager.NewEthKeyManager()
	if err != nil {
		logger.Error("Failed to create key manager")
		return err
	}

	err = coinkM.LoadNodeKeys(secrets.EthereumKeyStore)
	if err != nil {
		logger.Error("Failed to set up ethereum keys")
		return err
	}

	btcKM, err := key_manager.NewBitCoinKeyManager(config.BtcRpc, config.BtcNetwork, secrets.BtcRpcUser, secrets.BtcRpcPassword, secrets.BtcSigners)
	if err != nil {
		logger.Error("Failed to set up node keys")
		return err
	}

	err = btcKM.LoadNodeKeys(secrets.BtcPrivateKey)
	if err != nil {
		logger.Error("Failed to set up btc keys")
		return err
	}

	ltcKM, err := key_manager.NewLiteCoinKeyManager(config.LtcRpc, config.LtcNetwork, secrets.LtcRpcUser, secrets.LtcRpcPassword, secrets.LtcSigners)
	if err != nil {
		logger.Error("Failed to create LTC key manager")
		return err
	}

	err = ltcKM.LoadNodeKeys(secrets.LtcPrivateKey)
	if err != nil {
		logger.Error("Failed to set up ltc keys")
		return err
	}
	//p,_:=ltcKM.GetPublicKey()
	//panic(p)

	bchKM, err := key_manager.NewBCHCoinKeyManager(config.BchRpc, config.BchNetwork, secrets.BchRpcUser, secrets.BchRpcPassword, secrets.BchSigners)
	if err != nil {
		logger.Error("Failed to create BCH key manager")
		return err
	}
	err = bchKM.LoadNodeKeys(secrets.BchPrivateKey)
	if err != nil {
		logger.Errorf("%s : Failed to set up bch keys", err)
		return err
	}

	kms := map[string]key_manager.KeyManager{
		coin.BLOCKCHAIN_QUANTA: quanta,
		coin.BLOCKCHAIN_BTC:    btcKM,
		coin.BLOCKCHAIN_ETH:    coinkM,
		coin.BLOCKCHAIN_BCH:    bchKM,
		coin.BLOCKCHAIN_LTC:    ltcKM,
	}

	service := key_manager.NewRemoteKeyManagerService(kms, *secrets)
	service.Serve(fmt.Sprintf(":%d", port))

	return nil
}
