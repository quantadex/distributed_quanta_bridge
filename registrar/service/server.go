package service

import (
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"fmt"
	"time"
	"strings"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"bytes"
	"context"
)

type Server struct {
	url string
	registry *Registry
	handlers *http.ServeMux
	logger logger.Logger
	httpService *http.Server
}

func NewServer(registry *Registry, url string, logger logger.Logger) *Server {
	return &Server{registry: registry, url: url, logger: logger, httpService: &http.Server{ Addr: url }}
}

func SendHealthCheck(n *manifest.TrustNode) {
	url := fmt.Sprintf("http://%s:%s/node/api/healthcheck", n.IP, n.Port)
	http.Get(url)
}

func (server *Server) Stop() {
	server.httpService.Shutdown(context.Background())
}

func (server *Server) DoHealthCheck(interval int) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		for range ticker.C {
			for _, v := range server.registry.Manifest().Nodes {
				SendHealthCheck(v)
			}
		}
	}()
}

func (server *Server) Start() {
	server.logger.Infof("Server will be started at %s...\n", server.url)
	server.setRoute()

	if err := server.httpService.ListenAndServe(); err != nil {
		server.logger.Error("Start server failed: " + err.Error())
		return
	}
}

func (server *Server) setRoute() {
	server.handlers = http.NewServeMux()
	server.handlers.HandleFunc("/registry/api/health", server.receiveHealthCheck)
	server.handlers.HandleFunc("/registry/api/manifest", server.manifest)
	server.handlers.HandleFunc("/registry/api/register", server.register)
	server.handlers.HandleFunc("/registry/api/getaddr", server.getaddress)
	fs := http.FileServer(http.Dir("static"))
	server.handlers.Handle("/static/", http.StripPrefix("/static/", fs))
	server.httpService.Handler = server.handlers
}

func (server *Server) getaddress(w http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("token")
	auth := strings.Split(request.Header.Get("Authorization"),":")

	if query == "ETH" {
		server.registry.GetAddress(auth[0])
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported coin"))
	}
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

	if (server.registry.Manifest().ManifestComplete()) {
		go server.Broadcast(server.registry.Manifest())
	}
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
		//fmt.Println("Send manifest to ", url)
		http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	}
}
