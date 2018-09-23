package main

import "github.com/quantadex/distributed_quanta_bridge/common/manifest"

type Registry struct {
	manifest *manifest.Manifest
}

func (r *Registry) AddNode(n *NodeInfo) error {
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
	return r
}