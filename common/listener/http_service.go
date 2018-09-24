package listener

import (
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"io/ioutil"
	"fmt"
)

type HttpListener struct {
	handlers *http.ServeMux
}

func (h *HttpListener) AttachQueue(queueName string) error {
	return nil
}

func (h *HttpListener) AddEndpoint(name string, route string) error {
	queue.GetGlobalQueue().CreateQueue(name)

	h.handlers.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		queue.GetGlobalQueue().Put(name, bodyBytes)
	})

	return nil
}

/**
  * Run blocks
 */
func (h *HttpListener) Run(host string, port int) error {
	var err error
	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port),  h.handlers); err != nil {
		return nil
	}
	return err
}
