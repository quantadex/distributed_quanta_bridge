package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Webhook struct {
	TrustNode *TrustNode
	EventChan chan map[string]string
}

func NewWebhook(node *TrustNode) *Webhook {
	return &Webhook{TrustNode: node, EventChan: make(chan map[string]string, 100)}
}

func postData(data, url string) error {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	fmt.Println("response = ", err)
	if err != nil {
		return err
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			fmt.Println(string(body))
		}
	}
	return nil
}

func (w *Webhook) Start() {
	for {
		select {
		case msg := <-w.EventChan:
			if quanta, ok := msg["Deposit.Successful"]; ok {
				url := w.TrustNode.rDb.GetURLForQuanta(quanta, "Deposit.Successful")
				if url == "" {
					break
				}
				postData("DepositSuccessful", url)
			}
		}
	}
}
