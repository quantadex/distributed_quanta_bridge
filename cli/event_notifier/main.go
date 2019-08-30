package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/webhook_process"
	"github.com/spf13/viper"
	"io/ioutil"
	"strconv"
)

func main() {
	privKey := flag.String("key", "", "private key")
	request := flag.String("request", "something.json", "request")
	stop := flag.String("stop", "", "stop")
	configFile := flag.String("config", "config.yml", "configuration file")
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

	var js []byte
	if *request != "" {
		js, err = ioutil.ReadFile(*request)
		if err != nil {
			panic(err)
		}
	}

	process := webhook_process.NewWebhookServer(":5300", log, fmt.Sprintf("http://localhost:%d", config.ExternalListenPort), *privKey, string(js))
	process.Start()

	if *stop != "" {
		process.Stop()
	}
}
