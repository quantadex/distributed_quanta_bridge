package main

import (
    "github.com/quantadex/distributed_quanta_bridge/common/listener"
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
    "strconv"
)

/**
 * This is the node listener. It runs as a separate service.
 * It is a server that gets messages sent to the node from registrar
 * and other peers and queues them.
 *
 * The queue is polled by the actual node logic
 */
func createNodeListener(q queue.Queue, listenIp string, listenPort int) listener.Listener {
    log, err := logger.NewLogger(strconv.Itoa(listenPort))
    if err != nil {
        return nil
    }
    listener, err := listener.NewListener()
    if err != nil {
        log.Error("Failed to create listener module")
        return nil
    }
    err = listener.AttachQueue(q)
    if err != nil {
        log.Error("Failed to attach to listener queue")
        return nil
    }

    // manifest update from registry
    err = listener.AddEndpoint(queue.MANIFEST_QUEUE, "/node/api/manifest")
    if err != nil {
       log.Error("Failed to create endpoint")
       return nil
    }

    err = listener.AddEndpoint(queue.HEALTH_QUEUE, "/node/api/healthcheck")
    if err != nil {
        log.Error("Failed to create endpoint")
        return nil
    }
    err = listener.AddEndpoint(queue.PEERMSG_QUEUE, "/node/api/peer")
    if err != nil {
        log.Error("Failed to create endpoint")
        return nil
    }

    err = listener.AddEndpoint(queue.REFUNDMSG_QUEUE, "/node/api/refund")
    if err != nil {
        log.Error("Failed to create endpoint")
        return nil
    }

    return listener
}
