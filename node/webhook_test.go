package main

import (
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/node/webhook"
	"github.com/quantadex/distributed_quanta_bridge/webhook_process"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

type WebHookTest struct {
	f func(event string)
}

func (w *WebHookTest) ProcessEvent(event string) error {
	w.f(event)
	return nil
}

func TestWebhook(t *testing.T) {
	r := StartRegistry(3, ":6000")
	nodes := StartNodes(test.GRAPHENE_ISSUER, test.GRAPHENE_TRUST, test.ETHER_NETWORKS[test.ROPSTEN], 3)
	defer func() {
		StopNodes(nodes, []int{0, 1, 2})
		StopRegistry(r)
	}()
	time.Sleep(time.Millisecond * 250)
	eventExpect := []string{"Deposit_Pending", "Deposit_Successful"}
	eventNum := 0
	testWebhook := &WebHookTest{f: func(event string) {
		// if event is not what we expect
		var e *webhook.Event
		json.Unmarshal([]byte(event), &e)
		fmt.Println(e.Name)
		assert.Equal(t, e.Name, eventExpect[eventNum])
		eventNum++
	}}

	log, _ := logger.NewLogger(strconv.Itoa(5300))
	webhook_client := webhook_process.NewWebhookServerCustom(fmt.Sprintf(":%d", 5300), log, "http://localhost:5200", testWebhook, "5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj", `{"url": "http://localhost:5300/events", "events":["Deposit_Successful", "Deposit_Pending"]}`)

	go webhook_client.Start()
	time.Sleep(time.Second * 5)

	nodes[0].webhook.GetEventsChan() <- webhook.Event{"Deposit_Pending", "pooja", "hgcdhg"}
	nodes[0].webhook.GetEventsChan() <- webhook.Event{"Deposit_Successful", "pooja", "hgcdhg"}
	time.Sleep(time.Second * 10)

	webhook_client.Stop()
}
