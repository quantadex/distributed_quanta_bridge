package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"net/http"
	"strings"
)

type Server struct {
	url         string
	publicKey   string
	listenIp	string
	handlers    *mux.Router
	logger      logger.Logger
	httpService *http.Server
	kv          kv_store.KVStore
	db          *db.DB
	trustNode   *TrustNode
	coinNames    []string
}

func NewApiServer(trustNode *TrustNode, coinNames []string, publicKey string, listenIp string, kv kv_store.KVStore, db *db.DB, url string, logger logger.Logger) *Server {
	return &Server{trustNode: trustNode, coinNames: coinNames, publicKey: publicKey, listenIp: listenIp, url: url, logger: logger, kv: kv, db: db, httpService: &http.Server{Addr: url}}
}

func (server *Server) Stop() {
	server.httpService.Shutdown(context.Background())
}

func (server *Server) Start() {
	server.logger.Infof("REST API started at %s...\n", server.url)
	server.setRoute()

	if err := server.httpService.ListenAndServe(); err != nil {
		server.logger.Error("Start server failed: " + err.Error())
		return
	}
}

func (server *Server) setRoute() {
	server.handlers = mux.NewRouter()
	server.handlers.HandleFunc("/api/address/eth/{quanta}", server.addressHandler)
	server.handlers.HandleFunc("/api/address/new/{blockchain}/{quanta}", server.newAddressHandler)
	server.handlers.HandleFunc("/api/history", server.historyHandler)
	server.handlers.HandleFunc("/api/status", server.statusHandler)

	server.httpService.Handler = server.handlers
}

func (server *Server) newAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockchain := strings.ToUpper(vars["blockchain"])
	quanta := vars["quanta"]

	if blockchain == "BTC" {
		// if client, broadcast it
		if r.Header.Get("IS_PEER") != "true" {
			for k, _ := range server.trustNode.man.Nodes {
				if k != server.trustNode.nodeID {
					peer := server.trustNode.man.Nodes[k]
					url := fmt.Sprintf("http://%s:%s%s", peer.IP, peer.ExternalPort, r.RequestURI)
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
					}
					req.Header.Set("IS_PEER", "true")
					client := &http.Client{}
					_, err = client.Do(req)

				}
			}
		}

		forwardInput, err := server.trustNode.CreateMultisig(blockchain, quanta)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		json, err := json.Marshal(forwardInput)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(json)

		// else just return result
	} else {
		w.Write([]byte("Not supported"))
	}
}

func (server *Server) addressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quanta := strings.ToUpper(vars["quanta"])
	values, err := server.kv.GetAllValues(control.ETHADDR_QUANTAADDR)
	fmt.Printf("%v %v", values, err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addresses := []string{}
	for k, v := range values {
		if v == quanta {
			addresses = append(addresses, k)
		}
	}

	data, _ := json.Marshal(addresses)

	w.Write(data)
}

func (server *Server) historyHandler(w http.ResponseWriter, r *http.Request) {
	txs, err := db.QueryAllTX(server.db)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		data, _ := json.Marshal(txs)
		w.Write(data)
	}
}

func (server *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]string {}
	status["VERSION"] = Version
	status["BUILDTIME"] = BuildStamp
	status["GITHASH"] = GitHash
	status["LISTEN_IP"] = server.listenIp
	status["PUBLIC_KEY"] = server.publicKey

	for _, coinName := range server.coinNames {
		lastProcessed, valid := control.GetLastBlock(server.kv, coinName)
		if valid {
			status["CURRENTBLOCK:" + coinName] = fmt.Sprintf("%d",lastProcessed)
		}
	}

	data, _ := json.Marshal(status)
	w.Write(data)
}