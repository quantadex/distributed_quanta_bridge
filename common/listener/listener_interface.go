package listener

import "github.com/quantadex/distributed_quanta_bridge/common/queue"

/**
 * Listener
 *
 * Listener is a simple http server that listens on the given port.
 * It attaches to the specified queue interface.
 * One queue is created for each endpoint added.
 * It will listen on the specified endpoints and add the received message to the matching queue.
 * It will return 200 for any successfull POST on a known endpoint where the message was queued.
 */
type Listener interface {
	/**
	 * AttachQueue
	 *
	 * Inputs:
	 *  queueName string - the named instance of the queue service to attach to via the QueueInterface
	 *
	 * Outputs:
	 *  - nil on success
	 *  - propogate error
	 *
	 * Attaches to the named queue service via the queue interface
	 */
	AttachQueue(queueName queue.Queue) error

	/**
	 * AddEndpoint
	 *
	 * Inputs:
	 *  name string - the name of the endpoint (e.g. healthcheck)
	 *  route string - the route of endpoint (e.g. /api/v3/node/healthcheck)
	 *
	 * Outputs:
	 *  - nil on success
	 *  - propogate error
	 *
	 * Adds a route to listen to. Created a queue in the queue service with endpoint name.
	 * Messages received on the endpoint will go to that queue
	 */
	AddEndpoint(name string, route string) error

	/*
	 * Run
	 *
	 * Start the infinite listening loop
	 */
	Run(host string, port int) error

	Stop()
}

func NewListener() (Listener, error) {
	return &HttpListener{}, nil
}
