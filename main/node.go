package main

import (
    "time"
    "github.com/quantadex/distributed_quanta_bridge/common/logger"
    "github.com/quantadex/distributed_quanta_bridge/common/kv_store"
    "github.com/quantadex/distributed_quanta_bridge/common/manifest"
    "github.com/quantadex/distributed_quanta_bridge/trust/quanta"
    "github.com/quantadex/distributed_quanta_bridge/trust/coin"
    "github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
    "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
    "github.com/quantadex/distributed_quanta_bridge/trust/control"
    "github.com/quantadex/distributed_quanta_bridge/trust/registrar_client"
    "github.com/spf13/viper"
    "fmt"
)

const (
    USE_PREV_KEYS = "USE_PREV_KEYS"
    DB_NAME = "KV_STORE_NAME"
    COIN_NAME = "COIN_NAME"
    LISTEN_IP = "LISTEN_IP"
    LISTEN_PORT = "LISTEN_PORT"
)

/**
 * trustNode
 *
 * The top-most object holding all state for the trust node.
 */
type TrustNode struct {
    log      logger.Logger
    kM       key_manager.KeyManager
    man      *manifest.Manifest
    q        quanta.Quanta
    c        coin.Coin
    db       kv_store.KVStore
    peer     peer_contact.PeerContact
    reg      registrar_client.RegistrarContact
    cTQ      *control.CoinToQuanta
    qTC      *control.QuantaToCoin
    nodeID   int
    coinName string
}

/**
 * init
 *
 * Initialize all sub-modules. Attach to databases.
 */
func initNode() (*TrustNode, bool) {
    var err error
    node := &TrustNode{}
    node.log, err = logger.NewLogger()
    if err != nil {
        return nil, false
    }

    node.kM, err = key_manager.NewKeyManager()
    if err != nil {
        node.log.Error("Failed to create key manager")
        return nil, false
    }
    reuseKeys := viper.GetBool(USE_PREV_KEYS)
    if reuseKeys == true {
       err = node.kM.LoadNodeKeys()
    } else {
       err = node.kM.CreateNodeKeys()
    }
    if err != nil {
        node.log.Error("Failed to set up node keys")
        return nil, false
    }

    node.db, err = kv_store.NewKVStore()
    if err != nil {
        node.log.Error("Failed to create database")
        return nil, false
    }
    err = node.db.Connect(viper.GetString(DB_NAME))
    if err != nil {
        node.log.Error("Failed to connect to database")
        return nil, false
    }

    node.c, err = coin.NewCoin()
    if err != nil {
        node.log.Error("Failed to create new coin")
        return nil, false
    }
    node.coinName = viper.GetString(COIN_NAME)
    err = node.c.Attach(node.coinName)
    if err != nil {
        node.log.Error("Failed to attach to coin")
        return nil, false
    }

    node.q, err = quanta.NewQuanta()
    if err != nil {
        node.log.Error("Failed to create quanta")
        return nil, false
    }
    err = node.q.Attach()
    if err != nil {
        node.log.Error("Failed to connect to quanta")
        return nil, false
    }

    node.peer, err = peer_contact.NewPeerContact()
    if err != nil {
        node.log.Error("Failed to create peer interface")
        return nil, false
    }
    err = node.peer.AttachToListener()
    if err != nil {
        node.log.Error("Failed to attach to peer listener")
        return nil, false
    }

    node.reg, err = registrar_client.NewRegistrar()
    if err != nil {
        node.log.Error("Failed to create registrar interface")
        return nil, false
    }
    err = node.reg.GetRegistrar()
    if err != nil {
        node.log.Error("Failed to get to registrar")
        return nil, false
    }
    err = node.reg.AttachToListener()
    if err != nil {
        node.log.Error("Failed to attach to reg listener")
        return nil, false
    }
    return node, true
}

/**
 * registerNode
 *
 * Sends registration message to registrar and waits to be added to a quorum.
 * Checks for and responds to health messages.
 * When a quorum has been created of which this node is a part of, the registrar
 * will send it the manifest. Upon receiving a manifest this function returns
 */
func (n *TrustNode) registerNode() bool {
    nodeIP := viper.GetString(LISTEN_IP)
    nodePort := viper.GetString(LISTEN_PORT)
    if nodeIP == "" || nodePort == "" {
        n.log.Error("Node IP and port not set")
        return false
    }

    nodeKey, err := n.kM.GetPublicKey()
    if err != nil {
        n.log.Error("Failed to get public key")
        return false
    }

    err = n.reg.RegisterNode(nodeIP, nodePort, nodeKey)
    if err != nil {
        n.log.Error("Failed to send node info to registrar")
        return false
    }

    // Now we sit and wait to be added to quorum
    for {
        time.Sleep(time.Second)
        if n.reg.HealthCheckRequested() {
            err = n.reg.SendHealth("READY", nodeKey)
            if err != nil {
                n.log.Error("Failed to send health status to registrar")
                return false
            }
        }
        man := n.reg.GetManifest()
        if man != nil {
            n.man = man
            n.nodeID, err = n.man.FindNode(nodeIP, nodePort, nodeKey)
            if err != nil {
                n.log.Error("Node was not added to manifest")
                return false
            }
            return true
        }
    }
}

/**
 * initTrusts
 *
 * Once we are part of a quorum we can create the trusts.
 */
func (n *TrustNode) initTrust() {
    node := &TrustNode{}
    node.qTC = control.NewQuantaToCoin( n.log,
                                        n.db,
                                        n.c,
                                        n.q,
                                        n.man.QuantaAddress,
                                        n.man.ContractCallSite,
                                        n.kM,
                                        n.coinName,
                                        n.nodeID)

    node.cTQ = control.NewCoinToQuanta( n.log,
                                        n.db,
                                        n.c,
                                        n.q,
                                        n.man,
                                        n.kM,
                                        n.coinName,
                                        n.nodeID,
                                        n.peer)
}

/**
 * run
 *
 * An infinite loop where we sleep 1 second. Then process trusts.
 */
func (n *TrustNode) run() {
    nodeKey, err := n.kM.GetPublicKey()
    if err != nil {
        n.log.Error("Failed to get public key")
    }

    for true {
        if n.reg.HealthCheckRequested() {
            n.reg.SendHealth("RUNNING", nodeKey)
        }
        n.cTQ.DoLoop()
        n.qTC.DoLoop()
        time.Sleep(time.Second)
    }
}

/**
 * main
 *
 * Runs the trust node
 */
func main() {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    node, success := initNode()
    if !success {
        return
    }
    success = node.registerNode()
    if !success {
        return
    }
    node.initTrust()

    go nodeAgent()
    node.run()
}
