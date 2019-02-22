package service

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	"sync"
)

type Registry struct {
	manifest             *manifest.Manifest
	listener             *coin.Listener
	ownerEthereumKey     *ecdsa.PrivateKey
	trustEthereumAddress common.Address
	sync.RWMutex
}

func (r *Registry) AddNode(n *msgs.NodeInfo) error {
	r.Lock()
	defer r.Unlock()
	return r.manifest.AddNode(n.NodeIp, n.NodePort, n.NodeExternalPort, n.NodeKey, n.ChainAddress)
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

func NewRegistry(minNodes int) *Registry {
	r := &Registry{}
	r.manifest = manifest.CreateNewManifest(minNodes)
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
