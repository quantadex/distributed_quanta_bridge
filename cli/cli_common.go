package cli

import (
	"github.com/spf13/viper"
	"flag"
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"strconv"
	"github.com/op/go-logging"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/go-pg/pg"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
)

func Setup() (*common.Config, quanta.Quanta, *db.DB, kv_store.KVStore, logger.Logger){
	viper.SetConfigType("yaml")
	configFile := flag.String("config", "config.yml", "configuration file")
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

	err = kdb.Connect(config.DatabaseUrl)
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
	info, err := pg.ParseURL(config.DatabaseUrl)
	if err != nil {
		log.Error(err.Error())
	}
	//node.rDb.Debug()
	rDb.Connect(info.Addr, info.User, info.Password, info.Database)
	db.MigrateTx(rDb)
	db.MigrateKv(rDb)
	db.MigrateXC(rDb)

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

	return &config, quanta, rDb, kdb, log
}