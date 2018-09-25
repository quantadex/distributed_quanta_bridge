package main

import (
    "github.com/quantadex/distributed_quanta_bridge/common/listener"
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/spf13/viper"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
)

/**
 * This is the node listener. It runs as a separate service.
 * It is a server that gets messages sent to the node from registrar
 * and other peers and queues them.
 *
 * The queue is polled by the actual node logic
 */
func nodeAgent(q queue.Queue) {
    log, err := logger.NewLogger(viper.GetString(LISTEN_PORT))
    if err != nil {
        return
    }
    listener, err := listener.NewListener()
    if err != nil {
        log.Error("Failed to create listener module")
        return
    }
    err = listener.AttachQueue(q)
    if err != nil {
        log.Error("Failed to attach to listener queue")
        return
    }

    // manifest update from registry
    err = listener.AddEndpoint(queue.MANIFEST_QUEUE, "/node/api/manifest")
    if err != nil {
       log.Error("Failed to create endpoint")
       return
    }

    err = listener.AddEndpoint(queue.HEALTH_QUEUE, "/node/api/healthcheck")
    if err != nil {
        log.Error("Failed to create endpoint")
        return
    }
    err = listener.AddEndpoint(queue.PEERMSG_QUEUE, "/node/api/peer")
    if err != nil {
        log.Error("Failed to create endpoint")
        return
    }
    err = listener.Run(viper.GetString("LISTEN_IP"), viper.GetInt("LISTEN_PORT"))
    if err != nil {
        log.Error("Failed to start listener")
        return
    }
}
