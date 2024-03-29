package main

import (
	"context"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/metric"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	url           string
	publicKey     string
	listenIp      string
	handlers      *mux.Router
	logger        logger.Logger
	httpService   *http.Server
	kv            kv_store.KVStore
	db            *db.DB
	trustNode     *TrustNode
	coinNames     []string
	coins         []coin.Coin
	addressChange *AddressConsensus
	counter       uint64
	isTest        bool
}

func NewApiServer(trustNode *TrustNode, coinNames []string, publicKey string, listenIp string, kv kv_store.KVStore, db *db.DB, url string, logger logger.Logger) *Server {
	return &Server{trustNode: trustNode, coinNames: coinNames,
		coins:     []coin.Coin{trustNode.eth, trustNode.btc, trustNode.ltc, trustNode.bch},
		publicKey: publicKey,
		listenIp:  listenIp, url: url, logger: logger,
		kv: kv, db: db, httpService: &http.Server{Addr: url},
		addressChange: NewAddressConsensus(logger, trustNode, db, kv, trustNode.config.MinBlockReuse),
		isTest:        trustNode.config.IsTest}
}

func (server *Server) Stop() {
	server.addressChange.Stop()
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
	server.handlers.HandleFunc("/api/address/{blockchain}", server.listAddressHandler)
	server.handlers.HandleFunc("/api/address/{blockchain}/{quanta}", server.addressHandler)
	server.handlers.HandleFunc("/api/address_to_quanta/{blockchain}/{address}", server.addressToQuantaHandler)
	server.handlers.HandleFunc("/api/history", server.historyHandler)
	server.handlers.HandleFunc("/api/status", server.statusHandler)
	server.handlers.HandleFunc("/api/webhook/{webhookId}", server.deleteWebhookHandler)
	server.handlers.HandleFunc("/api/webhook", server.createOrListWebhookHandler)

	server.httpService.Handler = server.handlers
}

type RequestMsg struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

func (server *Server) getQuantaAddr(msg, sig string) (string, error) {
	pubKey, err := crypto.GetPublicKey(msg, sig)
	if err != nil {
		return "", err
	}
	quanta, err := server.trustNode.q.GetAccountFromPubKey(pubKey)
	if err != nil {
		return "", err
	}
	return quanta, nil
}

func (server *Server) createOrListWebhookHandler(w http.ResponseWriter, r *http.Request) {
	sig := r.Header.Get("Signature")
	date := r.Header.Get("Date")
	msg := r.URL.String() + date

	quanta, err := server.getQuantaAddr(msg, sig)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("could not get quanta address: " + err.Error()))
		return
	}

	var data []byte

	if r.Method == "POST" {
		bodyBytes, _ := ioutil.ReadAll(r.Body)

		var reqMsg *RequestMsg
		err := json.Unmarshal(bodyBytes, &reqMsg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not unmarshal"))
			return
		}

		id := strconv.Itoa(int(server.counter))
		err = server.db.AddWebhook(id, reqMsg.URL, reqMsg.Events, quanta)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("could not add to database"))
			return
		}
		server.counter++
		value := server.db.GetWebhookById(id)
		data, _ = json.Marshal(value)
		w.Write(data)

	} else if r.Method == "GET" {
		values := server.db.GetWebhooksByQuanta(quanta)
		if values == nil {
			values = []db.Webhook{}
		}
		data, _ = json.Marshal(values)
		w.Write(data)
	}
}

func (server *Server) deleteWebhookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.ToUpper(vars["webhookId"])

	sig := r.Header.Get("Signature")
	date := r.Header.Get("Date")
	msg := r.URL.String() + date

	quanta, err := server.getQuantaAddr(msg, sig)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("could not get quanta address: " + err.Error()))
		return
	}

	err = db.RemoveWebhook(server.db, id, quanta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not remove webhook_process"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (server *Server) addressToQuantaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockchain := strings.ToUpper(vars["blockchain"])
	if !(blockchain == coin.BLOCKCHAIN_BTC || blockchain == coin.BLOCKCHAIN_ETH || blockchain == coin.BLOCKCHAIN_LTC || blockchain == coin.BLOCKCHAIN_BCH) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not a supported blockchain"))
		return
	}
	address := vars["address"]
	quanta := server.db.GetCrosschainByAddressandBlockchain(address, blockchain)
	data, err := json.Marshal(quanta.QuantaAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}

func (server *Server) generateNewAddress(blockchain string, quanta string) (*crypto.ForwardInput, error) {
	if blockchain == coin.BLOCKCHAIN_BTC || blockchain == coin.BLOCKCHAIN_LTC || blockchain == coin.BLOCKCHAIN_BCH {
		forwardInput, err := server.trustNode.CreateMultisig(blockchain, quanta)
		return forwardInput, err
	} else {
		return nil, errors.New("not supported")
	}
}

func (server *Server) listAddressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockchain := strings.ToUpper(vars["blockchain"])
	if !(blockchain == coin.BLOCKCHAIN_BTC || blockchain == coin.BLOCKCHAIN_ETH || blockchain == coin.BLOCKCHAIN_LTC || blockchain == coin.BLOCKCHAIN_BCH) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not a supported blockchain"))
		return
	}

	values := server.db.GetCrosschainByBlockchain(blockchain)
	data, _ := json.Marshal(values)
	w.Write(data)
}

func (server *Server) addressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockchain := strings.ToUpper(vars["blockchain"])
	quanta := vars["quanta"]

	if !(blockchain == coin.BLOCKCHAIN_BTC || blockchain == coin.BLOCKCHAIN_ETH || blockchain == coin.BLOCKCHAIN_LTC || blockchain == coin.BLOCKCHAIN_BCH) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not a supported blockchain"))
		return
	}

	if !server.isTest {
		accountExists := server.trustNode.q.AccountExist(quanta)
		if !accountExists {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("not a valid user"))
			return
		}
	}

	values, err := db.GetCrosschainByBlockchainAndUser(server.db, blockchain, quanta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if r.Method == "GET" {
		if values == nil {
			values = []db.CrosschainAddress{}
		}
		data, _ := json.Marshal(values)
		w.Write(data)
		return
	}

	if len(values) == 0 {
		server.logger.Infof("Request new address %v %v", blockchain, quanta)

		err = server.addressChange.GetAddress(AddressChange{blockchain, quanta, "", server.counter})
		server.counter++
		if err != nil {
			server.logger.Errorf("Could not agree on address change:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		values, err = db.GetCrosschainByBlockchainAndUser(server.db, blockchain, quanta)
		if err != nil {
			server.logger.Errorf("Could not retrieve crosschain address for %s", quanta)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//server.logger.Infof("Updated the crosschain address for account : %s to %s", quanta, addr[0].Address)
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
		offset, _ = strconv.Atoi(offsetStr)
	}
	if limitStr == "" {
		limit = 0
	} else {
		limit, _ = strconv.Atoi(limitStr)
	}

	// combine the pending tx in here
	var txs []db.Transaction
	var err error
	if user == "" {
		txs, err = db.QueryAllTX(server.db, offset, limit)
	} else {
		// filter pending if we are including w/ user.
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

func (server *Server) GetThresholdValues(blockchain string) (int64, int64, error) {
	if blockchain == coin.BLOCKCHAIN_BTC {
		return server.trustNode.config.BtcDegradedThreshold, server.trustNode.config.BtcFailureThreshold, nil
	} else if blockchain == coin.BLOCKCHAIN_LTC {
		return server.trustNode.config.LtcDegradedThreshold, server.trustNode.config.LtcFailureThreshold, nil
	} else if blockchain == coin.BLOCKCHAIN_BCH {
		return server.trustNode.config.BchDegradedThreshold, server.trustNode.config.BchFailureThreshold, nil
	} else if blockchain == coin.BLOCKCHAIN_ETH {
		return server.trustNode.config.EthDegradedThreshold, server.trustNode.config.EthFailureThreshold, nil
	} else if blockchain == control.QUANTA {
		return server.trustNode.config.QuantaDegradedThreshold, server.trustNode.config.QuantaFailureThreshold, nil
	}
	return 0, 0, errors.New("unkown blockchain")
}

func (server *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	status := make(map[string]interface{})
	status["VERSION"] = Version
	status["BUILDTIME"] = BuildStamp
	status["GITHASH"] = GitHash
	status["PUBLIC_KEY"] = server.publicKey

	totalDegraded := int64(0)
	totalFailure := int64(0)

	for _, coin := range server.coins {
		degraded, failure, err := server.GetThresholdValues(coin.Blockchain())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to get threshold values for " + coin.Blockchain() + err.Error()))
			return
		}

		res, err := metric.GetBlockchainStatus(coin, server.kv, server.db, degraded, failure)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to get blockchainstatus for " + coin.Blockchain() + err.Error()))
			return
		}
		metric.IncrFailuresAndDegraded(res.State, &totalDegraded, &totalFailure)

		status[coin.Blockchain()] = res
	}

	degraded, failure, _ := server.GetThresholdValues(control.QUANTA)
	res, err := metric.GetBlockchainStatus(server.trustNode.q, server.kv, server.db, degraded, failure)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get blockchainstatus for " + control.QUANTA + err.Error()))
		return
	}
	metric.IncrFailuresAndDegraded(res.State, &totalDegraded, &totalFailure)
	status[control.QUANTA] = res

	depStatus, err := metric.GetDepositOrWithdrawalStatus(db.DEPOSIT, server.trustNode.config.DepDegradedThreshold, server.trustNode.config.DepFailureThreshold, server.trustNode.nodeID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get deposit status" + err.Error()))
		return
	}
	metric.IncrFailuresAndDegraded(depStatus.State, &totalDegraded, &totalFailure)
	status["DEPOSIT"] = depStatus

	withdrawStatus, err := metric.GetDepositOrWithdrawalStatus(db.WITHDRAWAL, server.trustNode.config.WithdrawDegradedThreshold, server.trustNode.config.WithdrawFailureThreshold, server.trustNode.nodeID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get withdrawal status" + err.Error()))
		return
	}
	metric.IncrFailuresAndDegraded(withdrawStatus.State, &totalDegraded, &totalFailure)
	status["WITHDRAWAL"] = withdrawStatus

	if totalFailure > 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else if totalDegraded > 0 {
		w.WriteHeader(http.StatusGatewayTimeout)
	}

	status["TOTAL_DEGRADED"] = totalDegraded
	status["TOTAL_FAILURE"] = totalFailure

	data, _ := json.Marshal(status)
	w.Write(data)
}
