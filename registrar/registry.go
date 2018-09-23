package main

import "github.com/quantadex/distributed_quanta_bridge/common/manifest"

type Registry struct {
	nodes []*NodeInfo
}

func (r *Registry) AddNode(n *NodeInfo) {

}

func (r *Registry) ReceiveHealth(nodeKey string, status string) {

}

func (r *Registry) Manifest() *manifest.Manifest {

}

func NewRegistry() *Registry {
	r := &Registry{}
	return r
}