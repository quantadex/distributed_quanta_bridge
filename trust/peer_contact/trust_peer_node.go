package peer_contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/listener"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/quantadex/quanta_book/common"
	"github.com/quantadex/quanta_book/consensus"
	"github.com/quantadex/quanta_book/consensus/cosi"
	"net/http"
	"time"
)

type TrustPeerNode struct {
	man       *manifest.Manifest
	peer      PeerContact
	nodeID    int
	q         queue.Queue
	queueName string
	apiUrl    string
	km        key_manager.KeyManager
}

func NewTrustPeerNode(man *manifest.Manifest, peer PeerContact, nodeID int, q queue.Queue, queueName string, apiUrl string, km key_manager.KeyManager) *TrustPeerNode {
	//fmt.Printf("setup peer node\n")
	//for _, p := range man.Nodes {
	//	println(p.Port)
	//}
	return &TrustPeerNode{
		man:       man,
		peer:      peer,
		nodeID:    nodeID,
		q:         q,
		queueName: queueName,
		apiUrl:    apiUrl,
		km:        km,
	}
}

func (t *TrustPeerNode) GetMsg() *cosi.CosiMessage {
	data, err := t.q.Get(t.queueName)
	if err != nil {
		//fmt.Printf("queue is empty\n")
		return nil
	}
	listenerData := data.(listener.ListenerData)

	msg := &cosi.CosiMessage{}
	err = json.Unmarshal(listenerData.Body, msg)
	if err != nil {
		fmt.Printf("Unable to parse json\n")
		return nil
	}

	return msg
}

func (t *TrustPeerNode) SendMsg(destinationNodeID int, msg interface{}) error {
	peer := t.man.Nodes[destinationNodeID]
	url := fmt.Sprintf("http://%s:%s%s", peer.IP, peer.Port, t.apiUrl)
	fmt.Println("Send msg peer " + url)

	data, err := json.Marshal(&msg)
	if err != nil {
		return errors.New("unable to marshall")
	}

	_, err = t.km.SignMessage(data)
	if err != nil {
		return err
	}
	//println(string(signature))

	_, err = http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		println("Error! " + err.Error())
	}
	return err

	//}
	//return errors.New("unable to sign message")
}

func (t *TrustPeerNode) BroadcastMessage(callFunc string, req interface{}) {
	for k, _ := range t.man.Nodes {
		if k != t.nodeID {
			t.SendMsg(k, req)
		}
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
	for _, v := range t.man.Nodes {
		members = append(members, &common.Member{
			Address: v.IP + ":" + v.Port,
		})
	}
	return members
}
