package main

import (
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
)

type Registry struct {
	manifest *manifest.Manifest
}

func (r *Registry) AddNode(n *msgs.NodeInfo) error {
	return r.manifest.AddNode(n.NodeIp, n.NodePort,n.NodeKey)
}

func (r *Registry) ReceiveHealth(nodeKey string, state string) error {
	return r.manifest.UpdateState(nodeKey, state)
}

func (r *Registry) Manifest() *manifest.Manifest {
	return r.manifest
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.manifest = manifest.CreateNewManifest(3)
	return r
}