package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"sync"
)

type Registry struct {
	manifest             *manifest.Manifest
	listener             *coin.Listener
	ownerEthereumKey     *ecdsa.PrivateKey
	trustEthereumAddress common.Address
	sync.RWMutex
	path string
}

func (r *Registry) AddNode(n *msgs.NodeInfo) error {
	r.Lock()
	defer r.Unlock()

	for _, v := range r.manifest.Nodes {
		if v.PubKey == n.NodeKey {
			v.IP = n.NodeIp
		}
	}
	err := r.manifest.AddNode(n.NodeIp, n.NodePort, n.NodeExternalPort, n.NodeKey, n.ChainAddress)
	r.SaveManifest(n)
	return err
}

func (r *Registry) SaveManifest(n *msgs.NodeInfo) {
	b, _ := r.manifest.GetJSON()
	err := ioutil.WriteFile(r.path, b, 0644)
	if err != nil {
		fmt.Println("Error while writing to the file :", err)
	}
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

func NewRegistry(minNodes int, path string) *Registry {
	r := &Registry{}
	filePath := path + "/manifest.yml"
	r.path = filePath
	os.Remove(filePath)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			r.manifest = manifest.CreateNewManifest(minNodes)
		}
	} else {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil
		} else {
			var man *manifest.Manifest
			err = json.Unmarshal(data, &man)
			if err != nil {
				return nil
			}
			r.manifest = man
		}
	}

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
