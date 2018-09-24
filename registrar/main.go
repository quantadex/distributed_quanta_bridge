package registrar

import (
	"fmt"
	"net/http"
	"github.com/spf13/viper"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
)

type Server struct {
	url string
	registry *Registry
	handlers *http.ServeMux
}

func SendHealthCheck(n *manifest.TrustNode) {
	url := fmt.Sprintf("http://%s:%p/node/api/healthcheck")
	http.Get(url)
}

func (server *Server) DoHealthCheck() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(viper.GetInt("HEALTH_INTERVAL")))
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

	if err := http.ListenAndServe(server.url,  server.handlers); err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) setRoute() {
	server.handlers.HandleFunc("/registry/api/health", server.receiveHealthCheck)
	server.handlers.HandleFunc("/registry/api/manifest", server.manifest)
	server.handlers.HandleFunc("/registry/api/register", server.register)
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
	m := server.registry.Manifest()
	if data, err := json.Marshal(&m); err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
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
	s.setRoute()
	s.registry = NewRegistry()
	s.url = viper.GetString("server_url")

	s.DoHealthCheck()
	s.Start()
}