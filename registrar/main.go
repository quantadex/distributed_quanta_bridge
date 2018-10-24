package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/registrar/service"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	logger, _ := logger.NewLogger("registrar")
	s := service.NewServer(service.NewRegistry(), viper.GetString("server_url"), logger)
	s.DoHealthCheck(5)
	s.Start()
}