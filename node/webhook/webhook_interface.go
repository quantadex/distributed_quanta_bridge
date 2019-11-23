package webhook

import (
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
)

type WebhookInterface interface {
	ProcessEvents()
	GetEventsChan() chan Event
	Stop()
}

func NewWebhook(rDb *db.DB) WebhookInterface {
	log, _ := logger.NewLogger("webhook")
	return &Webhook{rDb: rDb, eventChan: make(chan Event, 100), doneChan: make(chan bool, 1), log: log}
}
