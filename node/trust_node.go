package main

import (
	"github.com/go-errors/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/listener"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/quantadex/distributed_quanta_bridge/trust/registrar_client"
	"strconv"
	"time"
)

const (
	USE_PREV_KEYS = "USE_PREV_KEYS"
	KV_DB_NAME    = "KV_DB_NAME"
	COIN_NAME     = "COIN_NAME"
	LISTEN_IP     = "LISTEN_IP"
	LISTEN_PORT   = "LISTEN_PORT"
)

/**
 * trustNode
 *
 * The top-most object holding all state for the trust node.
 */
type TrustNode struct {
	log      logger.Logger
	quantakM key_manager.KeyManager
	coinkM   key_manager.KeyManager
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
	queue    queue.Queue
	listener listener.Listener

	doneChan chan bool
}

type Config struct {
	ListenIp           string
	ListenPort         int
	UsePrevKeys        bool
	KvDbName           string
	CoinName           string
	IssuerAddress      string
	NodeKey            string
	HorizonUrl         string
	NetworkPassphrase  string
	RegistrarIp        string
	RegistrarPort      int
	EthereumNetworkId  string
	EthereumBlockStart int64
	EthereumRpc        string
	EthereumTrustAddr  string
	EthereumKeyStore   string
}

/**
 * init
 *
 * Initialize all sub-modules. Attach to databases.
 */
func initNode(config Config, targetCoin coin.Coin) (*TrustNode, bool) {
	var err error
	node := &TrustNode{}
	node.doneChan = make(chan bool, 1)
	node.queue = queue.NewMemoryQueue()
	node.log, err = logger.NewLogger(strconv.Itoa(config.ListenPort))
	if err != nil {
		return nil, false
	}

	node.quantakM, err = key_manager.NewKeyManager(config.NetworkPassphrase)
	if err != nil {
		node.log.Error("Failed to create key manager")
		return nil, false
	}
	reuseKeys := config.UsePrevKeys
	if reuseKeys == true {
		err = node.quantakM.LoadNodeKeys(config.NodeKey)
	} else {
		err = node.quantakM.CreateNodeKeys()
	}
	if err != nil {
		node.log.Error("Failed to set up node keys")
		return nil, false
	}

	node.coinkM, err = key_manager.NewEthKeyManager()
	if err != nil {
		node.log.Error("Failed to create key manager")
		return nil, false
	}

	err = node.coinkM.LoadNodeKeys(config.EthereumKeyStore)
	if err != nil {
		node.log.Error("Failed to set up ethereum keys")
		return nil, false
	}

	needsInitialize := !kv_store.DbExists(config.KvDbName)

	node.db, err = kv_store.NewKVStore()
	if err != nil {
		node.log.Error("Failed to create database")
		return nil, false
	}

	err = node.db.Connect(config.KvDbName)
	if err != nil {
		node.log.Error("Failed to connect to database")
		return nil, false
	}

	if needsInitialize {
		node.log.Info("Initialize ledger")
		control.InitLedger(node.db)
	}

	node.c = targetCoin

	node.coinName = config.CoinName
	err = node.c.Attach()
	if err != nil {
		node.log.Error("Failed to attach to coin " + err.Error())
		return nil, false
	}

	node.q, err = quanta.NewQuanta(quanta.QuantaClientOptions{
		node.log,
		config.NetworkPassphrase,
		config.IssuerAddress,
		config.HorizonUrl,
	})

	node.q.AttachQueue(node.db)
	if err != nil {
		node.log.Error("Failed to create quanta")
		return nil, false
	}
	err = node.q.Attach()
	if err != nil {
		node.log.Error("Failed to connect to quanta")
		return nil, false
	}

	node.peer, err = peer_contact.NewPeerContact(config.NodeKey)
	if err != nil {
		node.log.Error("Failed to create peer interface")
		return nil, false
	}

	err = node.peer.AttachQueue(node.queue)
	if err != nil {
		node.log.Error("Failed to attach to peer listener")
		return nil, false
	}

	node.reg, err = registrar_client.NewRegistrar(config.RegistrarIp, config.RegistrarPort)
	if err != nil {
		node.log.Error("Failed to create registrar interface")
		return nil, false
	}
	err = node.reg.GetRegistrar()
	if err != nil {
		node.log.Error("Failed to get to registrar")
		return nil, false
	}
	err = node.reg.AttachQueue(node.queue)
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
func (n *TrustNode) registerNode(config Config) bool {
	nodeIP := config.ListenIp
	nodePort := config.ListenPort
	if nodeIP == "" {
		n.log.Error("Node IP and port not set")
		return false
	}

	nodeKey, err := n.quantakM.GetPublicKey()
	if err != nil {
		n.log.Error("Failed to get public key")
		return false
	}

	err = n.reg.RegisterNode(nodeIP, strconv.Itoa(nodePort), n.quantakM)
	if err != nil {
		n.log.Error("Failed to send node info to registrar " + err.Error())
		return false
	}

	// Now we sit and wait to be added to quorum
	for {
		n.log.Info("Wait to be added to quorum")
		time.Sleep(time.Second)
		if n.reg.HealthCheckRequested() {
			err = n.reg.SendHealth("READY", n.quantakM)
			if err != nil {
				n.log.Error("Failed to send health status to registrar")
				return false
			}
		}
		man := n.reg.GetManifest()
		if man != nil {
			// OVERRIDE WITH OUR OWN
			// man.ContractAddress = viper.GetString("TRUST_ETHEREUM_ADDR")
			n.man = man
			n.nodeID, err = n.man.FindNode(nodeIP, strconv.Itoa(nodePort), nodeKey)
			if err != nil {
				n.log.Error("Node was not added to manifest")
				return false
			}
			n.log.Info("Added to quorum")
			return true
		}
	}
}

/**
 * initTrusts
 *
 * Once we are part of a quorum we can create the trusts.
 */
func (n *TrustNode) initTrust(config Config) {
	n.log.Info("Trust initialized")
	n.qTC = control.NewQuantaToCoin(n.log,
		n.db,
		n.c,
		n.q,
		n.man,
		config.IssuerAddress,
		config.EthereumTrustAddr,
		n.coinkM,
		n.coinName,
		n.peer,
		n.queue,
		n.nodeID)

	n.cTQ = control.NewCoinToQuanta(n.log,
		n.db,
		n.c,
		n.q,
		n.man,
		n.quantakM,
		n.coinName,
		n.nodeID,
		n.peer,
		control.C2QOptions{
			config.EthereumTrustAddr,
			config.EthereumBlockStart,
		})
}

/**
 * run
 *
 * An infinite loop where we sleep 1 second. Then process trusts.
 */
func (n *TrustNode) run() {
	for true {
		select {
		case <-time.After(time.Second):
			if n.reg.HealthCheckRequested() {
				n.reg.SendHealth("RUNNING", n.quantakM)
			}
			blockIDs := n.cTQ.GetNewCoinBlockIDs()
			n.cTQ.DoLoop(blockIDs)

			cursor, _ := control.GetLastBlock(n.db, control.QUANTA)
			n.qTC.DoLoop(cursor)
		case <-n.doneChan:
			n.log.Infof("Exiting.")
			break
		}
	}
}

func bootstrapNode(config Config, targetCoin coin.Coin) *TrustNode {
	node, success := initNode(config, targetCoin)
	if !success {
		panic("Failed to init node")
		return nil
	}

	// start our http listener
	node.listener = createNodeListener(node.queue, config.ListenIp, config.ListenPort)

	go func() {
		err := node.listener.Run(config.ListenIp, config.ListenPort)
		if err != nil {
			node.log.Error("Failed to start listener")
			return
		}
	}()

	return node
}

func registerNode(config Config, node *TrustNode) error {
	success := node.registerNode(config)
	if !success {
		node.log.Error("Failed to register node")
		return errors.New("Failed to register node")
	}
	node.initTrust(config)
	return nil
}

func (n *TrustNode) Stop() {
	n.doneChan <- true
	n.listener.Stop()
}
