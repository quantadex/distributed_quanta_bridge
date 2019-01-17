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
	handlers    *mux.Router
	logger      logger.Logger
	httpService *http.Server
	kv          kv_store.KVStore
	db          *db.DB
}

func NewApiServer(kv kv_store.KVStore, db *db.DB, url string, logger logger.Logger) *Server {
	return &Server{url: url, logger: logger, kv: kv, db: db, httpService: &http.Server{Addr: url}}
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
	server.handlers.HandleFunc("/api/history", server.historyHandler)
	server.httpService.Handler = server.handlers
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
