package webhook_process

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/gorilla/mux"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/scorum/bitshares-go/sign"
	"io/ioutil"
	"net/http"
	"time"
)

type WebhookClient struct {
	url         string
	handlers    *mux.Router
	logger      logger.Logger
	httpService *http.Server
	facade      ProcessInterface
	restApiUrl  string
	privateKey  string
	request     string
	id          string
}

func NewWebhookServer(url string, logger logger.Logger, restApiUrl string) *WebhookClient {
	return &WebhookClient{url: url, logger: logger,
		httpService: &http.Server{Addr: url},
		facade:      NewFacade(), restApiUrl: restApiUrl}
}

func NewWebhookServerCustom(url string, logger logger.Logger, restApiUrl string, facade ProcessInterface, privateKey string, request string) *WebhookClient {
	return &WebhookClient{url: url, logger: logger,
		httpService: &http.Server{Addr: url},
		facade:      facade, restApiUrl: restApiUrl, privateKey: privateKey, request: request}
}

func (server *WebhookClient) Start() {
	// Register the webhook
	err := server.registerWebhook()
	if err != nil {
		server.logger.Error("Could not register webhook: " + err.Error())
	}

	server.logger.Infof("Webhook server started at %s...\n", server.url)
	server.setRoute()

	if err := server.httpService.ListenAndServe(); err != nil {
		server.logger.Error("Start server failed: " + err.Error())
		return
	}
}

type Response struct {
	Id     string
	Url    string
	Events []string
	Quanta string
}

func SignMessage(wif, msg string) (string, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return "", err
	}

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	digest := sha256.Sum256(bData.Bytes())

	sig := sign.SignBufferSha256(digest[:], w.PrivKey.ToECDSA())
	sigHex := hex.EncodeToString(sig)
	return sigHex, err
}

func (server *WebhookClient) registerWebhook() error {
	date := time.Now().String()
	msg := "/api/webhook" + date
	sig, _ := SignMessage(server.privateKey, msg)
	body := bytes.NewBuffer([]byte(server.request))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/webhook", server.restApiUrl), body)
	req.Header.Set("Signature", sig)
	req.Header.Set("Date", date)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response *Response
	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		return err
	}
	server.id = response.Id

	return nil
}

func (server *WebhookClient) unregisterWebhook() error {
	date := time.Now().String()
	msg := "/api/webhook/" + server.id + date
	sig, _ := SignMessage(server.privateKey, msg)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/webhook/%s", server.restApiUrl, server.id), nil)
	req.Header.Set("Signature", sig)
	req.Header.Set("Date", date)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (server *WebhookClient) Stop() {
	err := server.unregisterWebhook()
	if err != nil {
		server.logger.Error("Could not unregister webhook: " + err.Error())
	}
	server.httpService.Shutdown(context.Background())
}

func (server *WebhookClient) setRoute() {
	server.handlers = mux.NewRouter()
	server.handlers.HandleFunc("/events", server.post)

	server.httpService.Handler = server.handlers
}

func (server *WebhookClient) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	server.facade.ProcessEvent(string(bodyBytes))
}
