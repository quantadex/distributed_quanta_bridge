package main

import (
	"fmt"
	"net/http"
	"github.com/spf13/viper"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
)

type Server struct {
	url string
	registry *Registry
	handlers *http.ServeMux
}


func (server *Server) Start() {
	fmt.Printf("Server will be started at %s...\n", server.Url)
	if err := http.ListenAndServe(server.Url,  server.Handlers); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) setRoute() {
	server.handlers.HandleFunc("/health", server.receiveHealthCheck)
	server.handlers.HandleFunc("/manifest", server.manifest)
	server.handlers.HandleFunc("/register", server.register)
}

func (server *Server) register(w http.ResponseWriter, request *http.Request) {
	var msg RegisterReq
	err := json.NewDecoder(request.Body).Decode(&msg)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified := crypto.VerifyMessage(msg.Body, msg.Body.NodeKey, msg.Signature)

	if !verified {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	server.registry.AddNode(&msg.Body)
	w.WriteHeader(http.StatusOK)
}

func (server *Server) receiveHealthCheck(w http.ResponseWriter, request *http.Request) {
	var msg PingReq
	err := json.NewDecoder(request.Body).Decode(&msg)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified := crypto.VerifyMessage(msg.Body, msg.Body.NodeKey, msg.Signature)

	if !verified {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	server.registry.ReceiveHealth(msg.Body.NodeKey, msg.Body.Status)
	w.WriteHeader(http.StatusOK)
}

func (server *Server) manifest(w http.ResponseWriter, request *http.Request) {

}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	s := &Server{}
	s.setRoute()
	s.registry = NewRegistry()
	s.url = viper.GetString("server_url")

}