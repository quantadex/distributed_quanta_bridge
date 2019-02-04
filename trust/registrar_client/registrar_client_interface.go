package registrar_client

import (
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/common/queue"
)

/**
 * RegistrarContact
 *
 * This is the module through which a trust node communicates with the registrar.
 * All communication is async. Which means data is sent and we only look for a status 200 to ACK receipt.
 * Data is then sent by the registrar to the node's listener (which is a separate service).
 * The listener inserts all received data into the queue-service.
 * The RegistrarContract attaches to the queue-service and returns any data that has already been received when requested.
 *
 */
type RegistrarContact interface {
    /**
     * GetRegistrar
     *
     * Look up the IP (REGISTRAR_IP) and port (REGISTRAR_PORT) environment variables.
     * Stash these in the local object.
     * These will be used to send requests to the registar.
     * Return error if environment variables not found.
     */
    GetRegistrar() error

    /**
     * AttachToListener
     *
     * Connect to the queue-service for the node's listener.
     * The Queue's name is in env variable (NODE_LISTENER_QUEUE)
     * Stash the Queue object in the local object
     * Return error if no variable or propogate error from Connect()
     */
    AttachQueue(queue queue.Queue) error

    /**
     * RegisterNode
     *
     * Sends your node's info to the registrar's node_registry endpoint
     * POST /register
     * Return error if failed to send or did not get status OK
     */
    RegisterNode(nodeIP string, nodePort string, km key_manager.KeyManager, keys map[string]string) error

    /**
     * SendHealth
     *
     * Send the given status to the registrar's node_healthcheck endpoint
     * POST /health
     * Return error if failed to send or did not get status OK
     */
    SendHealth(nodeState string, km key_manager.KeyManager) error

    /**
     * GetManifest
     *
     * Returns the manifest provided by the registrar.
     * Checks the listener queue-service for the queue corresponding to the "manifest" endpoint.
     * GET /manifest
     * If item is available pulls it off and returns.
     * Otherwise returns nil
     */
    GetManifest() *manifest.Manifest

    /**
     * HealthChechRequested
     *
     * Returns true if the registrar has sent a healthcheck request.
     * Checks the listener queue-service for the queue corresponding to the "healthcheck" endpoint.
     * If items were available (drain the queue) and return true.
     * Otherwise return false
     */
    HealthCheckRequested() bool
}

func NewRegistrar(address string, port int) (RegistrarContact, error) {
    return &RegistrarClient{address:address, port: port}, nil
}
