package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/h2non/gock"
	"github.com/magiconair/properties/assert"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/scorum/bitshares-go/sign"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestWebhook(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
		gock.Off()
	}()
	time.Sleep(time.Millisecond * 250)

	msg := "https://httpbin.org/post" + time.Now().String()
	sig, err := SignMessage("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj", msg)
	assert.Equal(t, err, nil)

	client := &http.Client{}

	//post
	body := bytes.NewBuffer([]byte(`{"url": "https://httpbin.org/post", "events":["Deposit.Successful"]}`))
	req, err := http.NewRequest("POST", "http://localhost:5200/api/webhook", body)
	req.Header.Set("Signature", sig)
	req.Header.Set("Msg", msg)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Println("post data ", resp.StatusCode, string(bodyText))

	//get
	req, err = http.NewRequest("GET", "http://localhost:5200/api/webhook", nil)
	req.Header.Set("Signature", sig)
	req.Header.Set("Msg", msg)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	fmt.Println("get data ", resp.StatusCode, string(bodyText))

	nodes[0].Webhook.EventChan <- map[string]string{"Deposit.Successful": "pooja"}

	time.Sleep(time.Second * 5)

	//delete
	req, err = http.NewRequest("DELETE", "http://localhost:5200/api/webhook/0", nil)
	req.Header.Set("Signature", sig)
	req.Header.Set("Msg", msg)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete data ", resp.StatusCode)
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
