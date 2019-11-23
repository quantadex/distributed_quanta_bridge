package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"math"
	"net/http"
	"time"
)

type Webhook struct {
	rDb       *db.DB
	eventChan chan Event
	doneChan  chan bool
	log       logger.Logger
}

type Event struct {
	Name   string
	Quanta string
	TxId   string
}

func (w *Webhook) postData(webhookUrl, event string) error {
	body := bytes.NewBuffer([]byte(event))
	req, err := http.NewRequest("POST", webhookUrl, body)
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (w *Webhook) postDataAndCheckFailed(event Event, url string) {
	eventByte, err := json.Marshal(event)
	if err != nil {
		w.log.Info("Could not unmarshal event")
		return
	}
	err = w.postData(url, string(eventByte))
	if err != nil {
		w.log.Infof("Failed message: %s for %s error: %s", event.Name, event.Quanta, err.Error())

		failed, err := w.rDb.GetFailedMessageByIdAndEvent(event.TxId, event.Name)
		if err != nil {
			w.rDb.AddFailedMessage(event.TxId, event.Name, event.Quanta, time.Now())
		} else {
			w.rDb.UpdateFailedMessage(event.TxId, event.Name, time.Now(), failed.NoOfRetries+1)
		}

	} else {
		w.log.Infof("Message successfully sent to %s", event.Quanta)
		db.RemoveFailedMessage(w.rDb, event.TxId, event.Name)
	}
}

func (w *Webhook) ProcessEvents() {
	doneFlag := false
	for {
		select {
		case msg := <-w.eventChan:
			webhooks, err := w.rDb.GetWebhookByQuantaAndEvent(msg.Quanta, msg.Name)
			if err != nil {
				break
			}
			for _, wh := range webhooks {
				w.postDataAndCheckFailed(msg, wh.URL)
			}

		case <-time.After(time.Second * 1):
			failed := w.rDb.GetAllFailedMessages()

			for _, f := range failed {
				x := float64(f.NoOfRetries)
				if time.Since(f.FirstRetry).Hours() > 72 {
					db.RemoveFailedMessage(w.rDb, f.Id, f.Event)

				} else if time.Since(f.LastRetry).Seconds()-math.Pow(x, 2) >= 0 {
					w.log.Infof("Retrying message: %s for %s", f.Event, f.Quanta)
					w.eventChan <- Event{f.Event, f.Quanta, f.Id}
				}
			}

		case <-w.doneChan:
			doneFlag = true
			break

		}
		if doneFlag {
			break
		}
	}
}

func (w *Webhook) GetEventsChan() chan Event {
	return w.eventChan
}

func (w *Webhook) Stop() {
	w.doneChan <- true
}
