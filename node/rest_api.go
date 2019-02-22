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
	"strconv"
	"strings"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/go-errors/errors"
	"io/ioutil"
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
	server.handlers.HandleFunc("/api/address/{blockchain}/{quanta}", server.addressHandler)
	server.handlers.HandleFunc("/api/history", server.historyHandler)
	server.handlers.HandleFunc("/api/status", server.statusHandler)

	server.httpService.Handler = server.handlers
}

func (server *Server) generateNewAddress(blockchain string, quanta string) (*crypto.ForwardInput, error) {
	if blockchain == coin.BLOCKCHAIN_BTC {
		forwardInput, err := server.trustNode.CreateMultisig(blockchain, quanta)
		return forwardInput, err
	} else {
		return nil, errors.New("not supported")
	}
}

func (server *Server) addressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockchain := strings.ToUpper(vars["blockchain"])
	quanta := vars["quanta"]

	if !(blockchain == coin.BLOCKCHAIN_BTC || blockchain == coin.BLOCKCHAIN_ETH) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not a supported blockchain"))
		return
	}

	values, err := db.GetCrosschainByBlockchainAndUser(server.db, blockchain, quanta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if len(values) == 0 && blockchain == coin.BLOCKCHAIN_BTC{
		_, err := server.generateNewAddress(blockchain, quanta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to generate address for " + blockchain))
			return
		}
		values, err = db.GetCrosschainByBlockchainAndUser(server.db, blockchain, quanta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to fetch address for " + blockchain))
			return
		}

		// if client, broadcast it
		if r.Header.Get("IS_PEER") != "true" {
			for k, _ := range server.trustNode.man.Nodes {
				if k != server.trustNode.nodeID {
					peer := server.trustNode.man.Nodes[k]
					url := fmt.Sprintf("http://%s:%s%s", peer.IP, peer.ExternalPort, r.RequestURI)
					req, err := http.NewRequest("GET", url, nil)
					println("Broadcast create address message to node", k, url)
					if err != nil {
						server.logger.Error("unable to build request: " + err.Error())
						continue
					}
					req.Header.Set("IS_PEER", "true")
					client := &http.Client{}
					res, err := client.Do(req)
					if err != nil {
						server.logger.Error("unable to broadcast: " + err.Error())
						continue
					}
					//if res.StatusCode != 200
					{
						bodyBytes, _ := ioutil.ReadAll(res.Body)
						server.logger.Errorf("Broadcast got code %s %s",res.Status, string(bodyBytes))
					}
				}
			}
		}
	}

	data, _ := json.Marshal(values)

	w.Write(data)
}

func (server *Server) historyHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	var offset, limit int
	if offsetStr == "" {
		offset = 0
	} else {
		offset,_ = strconv.Atoi(offsetStr)
	}
	if limitStr == "" {
		limit = 0
	} else {
		limit,_ = strconv.Atoi(limitStr)
	}


	var txs []db.Transaction
	var err error
	if user == "" {
		txs, err = db.QueryAllTX(server.db, offset, limit)
	} else {
		txs, err = db.QueryAllTXByUser(server.db, user, offset, limit)
	}

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