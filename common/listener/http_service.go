package listener

import (
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"io/ioutil"
	"fmt"
)

type HttpListener struct {
	handlers *http.ServeMux
	queue queue.Queue
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
	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port),  h.handlers); err != nil {
		return nil
	}
	return err
}
