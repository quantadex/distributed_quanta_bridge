package main

import (
	"fmt"
	"net/http"
	"github.com/spf13/viper"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"bytes"
)

type Server struct {
	url string
	registry *Registry
	handlers *http.ServeMux
	logger logger.Logger
}

func SendHealthCheck(n *manifest.TrustNode) {
	url := fmt.Sprintf("http://%s:%s/node/api/healthcheck", n.IP, n.Port)
	http.Get(url)
}

func (server *Server) DoHealthCheck() {
	ticker := time.NewTicker(time.Second * time.Duration(viper.GetInt("HEALTH_INTERVAL")))
	go func() {
		for range ticker.C {
			for _, v := range server.registry.manifest.Nodes {
				SendHealthCheck(v)
			}
		}
	}()
}

func (server *Server) Start() {
	fmt.Printf("Server will be started at %s...\n", server.url)
	server.setRoute()

	if err := http.ListenAndServe(server.url,  server.handlers); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) setRoute() {
	server.handlers = http.NewServeMux()
	server.handlers.HandleFunc("/registry/api/health", server.receiveHealthCheck)
	server.handlers.HandleFunc("/registry/api/manifest", server.manifest)
	server.handlers.HandleFunc("/registry/api/register", server.register)
}

func (server *Server) register(w http.ResponseWriter, request *http.Request) {
	var msg msgs.RegisterReq
	err := json.NewDecoder(request.Body).Decode(&msg)

	if err != nil {
		server.logger.Error("Cannot decode register msg")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verified := crypto.VerifyMessage(msg.Body, msg.Body.NodeKey, msg.Signature)

	if !verified {
		w.WriteHeader(http.StatusUnauthorized)
		server.logger.Error("Message fail signature verification")
		return
	}

	err = server.registry.AddNode(&msg.Body)
	if err != nil {
		server.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK) // fail silently
	} else {
		server.logger.Info(fmt.Sprintf("Node %s:%s added to registry", msg.Body.NodeIp, msg.Body.NodePort))
		w.WriteHeader(http.StatusOK)
	}

	go server.Broadcast(server.registry.Manifest())
}

func (server *Server) receiveHealthCheck(w http.ResponseWriter, request *http.Request) {
	var msg msgs.PingReq
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
	m := server.registry.Manifest()
	if data, err := json.Marshal(&m); err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (server *Server) Broadcast(m *manifest.Manifest) {
	for _, n := range m.Nodes {
		url := fmt.Sprintf("http://%s:%s/node/api/manifest", n.IP, n.Port)
		jsonBytes, _ := m.GetJSON()
		fmt.Println("Send manifest to ", url)
		http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	s := &Server{}
	s.registry = NewRegistry()
	s.url = viper.GetString("server_url")
	s.logger, _ = logger.NewLogger("registrar")
	s.DoHealthCheck()
	s.Start()
}