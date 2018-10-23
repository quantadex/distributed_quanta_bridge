package listener

import (
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"io/ioutil"
	"fmt"
	"context"
)

type HttpListener struct {
	handlers *http.ServeMux
	queue queue.Queue
	server *http.Server
}

func (h *HttpListener) AttachQueue(queue queue.Queue) error {
	h.queue = queue
	return nil
}

func (h *HttpListener) AddEndpoint(name string, route string) error {
	h.queue.CreateQueue(name)

	if h.handlers == nil{
		h.handlers = http.NewServeMux()
	}

	h.handlers.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("http data on " + route)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		h.queue.Put(name, bodyBytes)
	})

	return nil
}

/**
  * Run blocks
 */
func (h *HttpListener) Run(host string, port int) error {
	var err error
	h.server = &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: h.handlers}

	if err = h.server.ListenAndServe(); err != nil {
		return nil
	}
	return err
}

func (h *HttpListener) Stop() {
	h.server.Shutdown(context.Background())
}