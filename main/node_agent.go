package main

import (
	"os"
	"github.com/quantadex/distributed_quanta_bridge/common/listener"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
)

const QUEUE_NAME = "NODE_LISTENER_QUEUE"

/**
 * This is the node listener. It runs as a separate service.
 * It is a server that gets messages sent to the node from registrar
 * and other peers and queues them.
 *
 * The queue is polled by the actual node logic
 */
func main() {
    log, err := logger.NewLogger()
    if err != nil {
        return
    }
    listener, err := listener.NewListener()
    if err != nil {
        log.Error("Failed to create listener module")
        return
    }
    err = listener.AttachQueue(os.Getenv(QUEUE_NAME))
    if err != nil {
        log.Error("Failed to attach to listener queue")
        return
    }
    err = listener.AddEndpoint("manifest", "/node/api/manifest")
    if err != nil {
        log.Error("Failed to create endpoint")
        return
    }
    err = listener.AddEndpoint("healthcheck", "/node/api/healthcheck")
    if err != nil {
        log.Error("Failed to create endpoint")
        return
    }
    err = listener.AddEndpoint("peer_msg", "/node/api/peer")
    if err != nil {
        log.Error("Failed to create endpoint")
        return
    }
    err = listener.Run(os.Getenv(NODE_IP), os.Getenv(NODE_PORT))
    if err != nil {
        log.Error("Failed to start listener")
        return
    }
}
