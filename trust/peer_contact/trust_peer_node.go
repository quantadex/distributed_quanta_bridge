package peer_contact

import (
	"time"
	"github.com/quantadex/quanta_book/consensus"
	"github.com/quantadex/quanta_book/common"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/quanta_book/consensus/cosi"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"encoding/json"
	"fmt"
)

type TrustPeerNode struct{
	man *manifest.Manifest
	peer PeerContact
	nodeID int
	q queue.Queue
}

func NewTrustPeerNode(man *manifest.Manifest, peer PeerContact, nodeID int,q queue.Queue) *TrustPeerNode {
	return &TrustPeerNode{
		man: man,
		peer: peer,
		nodeID: nodeID,
		q: q,
	}
}

func (t *TrustPeerNode) GetMsg() *cosi.CosiMessage {
	data, err := t.q.Get(queue.REFUNDMSG_QUEUE)
	if err != nil {
		//fmt.Printf("queue is empty\n")
		return nil
	}

	msg := &cosi.CosiMessage{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		fmt.Printf("Unable to parse json\n")
		return nil
	}
	return msg
}


func (t *TrustPeerNode) BroadcastMessage(callFunc string, req interface{}) {
	for k,_ := range t.man.Nodes {
		t.peer.SendMsg(t.man, k, &PeerMessage{})
	}
}

func (t *TrustPeerNode) FetchFullBlock(timeout time.Duration, req interface{}) ([]consensus.FullBlockResponse, error) {
	panic("implement me")
}

func (t *TrustPeerNode) Address() string {
	return t.man.Nodes[t.nodeID].IP
}

func (t *TrustPeerNode) Ready() bool {
	// do nothing
	return true
}

func (t *TrustPeerNode) Start() {
	// do nothing
}

func (t *TrustPeerNode) Stop() {
	// do nothing
}

func (t *TrustPeerNode) Members() []*common.Member {
	members := []*common.Member{}
	for _,v := range t.man.Nodes {
		members = append(members, &common.Member{
			Address: v.IP + ":" + v.Port,
		})
	}
	return members
}
