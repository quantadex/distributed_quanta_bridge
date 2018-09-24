package listener

import (
	"net/http"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
)

type HttpListener struct {
	handlers *http.ServeMux
	q queue.Queue
}

func (h *HttpListener) AttachQueue(queueName string) error {
	panic("implement me")
}

func (h *HttpListener) AddEndpoint(name string, route string) error {
	h.handlers.HandleFunc(route, func(w http.ResponseWriter, request *http.Request) {

	})
	return nil
}

func (h *HttpListener) Run(host string, port int) error {
	var err error
	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port),  h.handlers); err != nil {
		return nil
	}
	return err
}
