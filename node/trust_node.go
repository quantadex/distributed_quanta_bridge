package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/go-errors/errors"
	"github.com/go-pg/pg"
	"github.com/op/go-logging"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/common/listener"
	"github.com/quantadex/distributed_quanta_bridge/common/logger"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"github.com/quantadex/distributed_quanta_bridge/trust/registrar_client"
	"github.com/scorum/bitshares-go/apis/database"
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
	btcKM    key_manager.KeyManager
	man      *manifest.Manifest
	q        quanta.Quanta
	eth      coin.Coin
	btc      coin.Coin
	db       kv_store.KVStore
	rDb      *db.DB
	peer     peer_contact.PeerContact
	reg      registrar_client.RegistrarContact
	cTQ      *control.CoinToQuanta
	qTC      *control.QuantaToCoin
	nodeID   int
	coinName string
	queue    queue.Queue
	listener listener.Listener
	restApi  *Server

	doneChan chan bool
}

/**
 * init
 *
 * Initialize all sub-modules. Attach to databases.
 */
func initNode(config common.Config, targetCoin coin.Coin) (*TrustNode, bool) {
	var err error
	node := &TrustNode{}
	node.doneChan = make(chan bool, 1)
	node.queue = queue.NewMemoryQueue()
	node.log, err = logger.NewLogger(strconv.Itoa(config.ListenPort))

	if err != nil {
		return nil, false
	}

	level, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		fmt.Println("Log level not parsed")
		level = logging.INFO
	}
	node.log.SetLogLevel(level)

	node.quantakM, err = key_manager.NewGrapheneKeyManager(config.ChainId)
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

	node.btcKM, err = key_manager.NewBitCoinKeyManager()
	if err != nil {
		node.log.Error("Failed to create BTC key manager")
		return nil, false
	}
	err = node.btcKM.LoadNodeKeys(config.BtcPrivateKey)
	if err != nil {
		node.log.Error("Failed to set up btc keys")
		return nil, false
	}

	needsInitialize := !kv_store.DbExists(config.KvDbName)

	node.db, err = kv_store.NewKVPGStore()
	if err != nil {
		node.log.Error("Failed to create database")
		return nil, false
	}

	println(config.DatabaseUrl)
	err = node.db.Connect(config.DatabaseUrl)
	if err != nil {
		node.log.Error("Failed to connect to database")
		return nil, false
	}

	if needsInitialize {
		node.log.Info("Initialize ledger")
		control.InitLedger(node.db)
	}

	// connect to do
	node.rDb = &db.DB{}
	info, err := pg.ParseURL(config.DatabaseUrl)
	if err != nil {
		node.log.Error(err.Error())
	}
	//node.rDb.Debug()
	node.rDb.Connect(info.Addr, info.User, info.Password, info.Database)
	db.MigrateTx(node.rDb)
	db.MigrateKv(node.rDb)
	db.MigrateXC(node.rDb)

	node.eth = targetCoin

	//node.coinName = coin.BLOCKCHAIN_ETH
	//node.coinName = config.CoinName
	err = node.eth.Attach()
	if err != nil {
		node.log.Error("Failed to attach to coin " + err.Error())
		return nil, false
	}

	// attach bitcoin
	coin, err := coin.NewBitcoinCoin(&chaincfg.RegressionNetParams, config.BtcSigners, "../blockchain/bitcoin/data")
	if err != nil {
		panic(fmt.Errorf("cannot create ethereum listener"))
	}

	err = coin.Attach()
	if err != nil {
		panic(err)
	}
	node.btc = coin

	node.q, err = quanta.NewQuantaGraphene(quanta.QuantaClientOptions{
		node.log,
		node.rDb,
		config.ChainId,
		config.IssuerAddress,
		config.NetworkUrl,
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

	pubKey, _ := node.quantakM.GetPublicKey()

	blockchain := make([]string, len(config.CoinMapping)+1)
	blockchain[0] = control.QUANTA
	i := 1
	for _, v := range config.CoinMapping {
		blockchain[i] = v
		i = i + 1
	}
	node.restApi = NewApiServer(node, blockchain, pubKey, config.ListenIp, node.db, node.rDb, fmt.Sprintf(":%d", config.ExternalListenPort), node.log)

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
func (n *TrustNode) registerNode(config common.Config) bool {
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

	btcPub, err := n.btcKM.GetPublicKey()
	chainAddress := map[string]string{
		"BTC": btcPub,
	}

	err = n.reg.RegisterNode(nodeIP, strconv.Itoa(nodePort), strconv.Itoa(config.ExternalListenPort), n.quantakM, chainAddress)
	if err != nil {
		n.log.Error("Failed to send node info to registrar " + err.Error())
		return false
	}

	go n.restApi.Start()

	// Now we sit and wait to be added to quorum
	for {
		//n.log.Info("Wait to be added to quorum")
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
func (n *TrustNode) initTrust(config common.Config) {
	n.log.Info("Trust initialized")

	coinInfo := make(map[string]*database.Asset)

	for _, v := range config.CoinMapping {
		asset, _ := n.q.GetAsset(v)
		coinInfo[v] = asset
	}

	n.qTC = control.NewQuantaToCoin(n.log,
		n.db,
		n.rDb,
		map[string]coin.Coin{coin.BLOCKCHAIN_ETH: n.eth, coin.BLOCKCHAIN_BTC: n.btc},
		n.q,
		n.man,
		config.IssuerAddress,
		config.EthereumTrustAddr,
		map[string]key_manager.KeyManager{coin.BLOCKCHAIN_ETH: n.coinkM, coin.BLOCKCHAIN_BTC: n.btcKM},
		config.CoinMapping,
		n.peer,
		n.queue,
		n.nodeID,
		coinInfo)

	n.cTQ = control.NewCoinToQuanta(n.log,
		n.db,
		n.rDb,
		n.eth,
		n.q,
		n.man,
		n.quantakM,
		n.nodeID,
		n.peer,
		n.queue,
		control.C2QOptions{
			config.EthereumTrustAddr,
			config.EthereumBlockStart,
		},
		quanta.QuantaClientOptions{
			NetworkUrl: config.NetworkUrl,
			Network:    config.ChainId,
			Issuer:     config.IssuerAddress,
		})
}

/**
 * run
 *
 * An infinite loop where we sleep 1 second. Then process trusts.
 */
func (n *TrustNode) run() {
	delayTime := time.Second
	//init := false
	for true {
		select {
		case <-time.After(delayTime):
			if n.reg.HealthCheckRequested() {
				n.reg.SendHealth("RUNNING", n.quantakM)
			}

			// handled in CLI
			//blockIDs := n.cTQ.GetNewCoinBlockIDs()
			//n.cTQ.DoLoop(blockIDs)

			blockIDs := n.qTC.GetNewCoinBlockIDs()
			//if init == false {
			//	cursor = 1911002
			//	init = true
			//}
			for _, cursor := range blockIDs {
				n.qTC.DoLoop(cursor)
			}

			// scale up time
			if len(blockIDs) == control.MAX_PROCESS_BLOCKS {
				delayTime = time.Second
			} else {
				delayTime = time.Second * 3
			}

		case <-n.doneChan:
			n.log.Infof("Exiting.")
			break
		}
	}
}

func bootstrapNode(config common.Config, targetCoin coin.Coin) *TrustNode {
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

func registerNode(config common.Config, node *TrustNode) error {
	success := node.registerNode(config)
	if !success {
		node.log.Error("Failed to register node")
		return errors.New("Failed to register node")
	}
	node.initTrust(config)
	return nil
}

func (n *TrustNode) CreateMultisig(blockchain string, accountId string) (*crypto.ForwardInput, error) {
	msig, err := n.btc.GenerateMultisig(accountId)
	if err != nil {
		return nil, err
	}

	//TODO: should we validate user?

	addr := crypto.ForwardInput{
		msig,
		common2.HexToAddress("0x0"),
		accountId,
		"",
		n.btc.Blockchain(),
	}
	err = n.rDb.AddCrosschainAddress(&addr)
	return &addr, err
}

func (n *TrustNode) Stop() {
	n.db.CloseDB()
	n.doneChan <- true
	n.qTC.Stop()
	n.listener.Stop()
	n.restApi.Stop()
}

func (node *TrustNode) StopListener() {
	node.listener.Stop()
}

func (node *TrustNode) StartListener(config common.Config) {
	go func() {
		err := node.listener.Run(config.ListenIp, config.ListenPort)
		if err != nil {
			node.log.Error("Failed to start listener")
			return
		}
	}()
}
