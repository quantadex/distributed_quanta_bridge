package service

import (
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"crypto/ecdsa"
	"sync"
)

type Registry struct {
	manifest *manifest.Manifest
	listener *coin.Listener
	ownerEthereumKey *ecdsa.PrivateKey
	trustEthereumAddress common.Address
	sync.RWMutex
}

func (r *Registry) AddNode(n *msgs.NodeInfo) error {
	r.Lock()
	defer r.Unlock()
	return r.manifest.AddNode(n.NodeIp, n.NodePort,n.NodeKey)
}

func (r *Registry) ReceiveHealth(nodeKey string, state string) error {
	return r.manifest.UpdateState(nodeKey, state)
}

func (r *Registry) Manifest() *manifest.Manifest {
	r.RLock()
	defer r.RUnlock()
	return r.manifest
}

func (r *Registry) GetAddress(quantaAddr string) (string, error) {
	return Forwarder.SubmitContract(r.listener.Client.(bind.ContractBackend), r.ownerEthereumKey, r.trustEthereumAddress, quantaAddr)
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.manifest = manifest.CreateNewManifest(3)
	r.listener = &coin.Listener{}
	//r.ownerEthereumKey = viper.GetString("CREATOR_ETHEREUM_KEY")
	r.trustEthereumAddress = common.HexToAddress(viper.GetString("TRUST_ETHEREUM_ADDRESS"))

	//err := r.listener.Start()
	//
	//if err != nil {
	//	log.Panic(err)
	//}

	return r
}