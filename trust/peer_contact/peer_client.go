package peer_contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/common/listener"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"net/http"
)

type PeerClient struct {
	q          queue.Queue
	privateKey string
}

func (p *PeerClient) AttachQueue(queue queue.Queue) error {
	p.q = queue
	return nil
}

func (p *PeerClient) SendMsg(m *manifest.Manifest, destinationNodeID int, peermsg *PeerMessage) error {
	peer := m.Nodes[destinationNodeID]
	url := fmt.Sprintf("http://%s:%s/node/api/peer", peer.IP, peer.Port)
	//fmt.Println("Send to peer " + url)

	msg := &PeerMsgRequest{}
	msg.Body = *peermsg

	if signature := crypto.SignMessage(msg.Body, p.privateKey); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(url, "application/json", bytes.NewReader(data))
		return nil
	}
	return errors.New("unable to sign message")
}

func (p *PeerClient) GetMsg() *PeerMessage {
	data, err := p.q.Get(queue.PEERMSG_QUEUE)
	if err != nil {
		//fmt.Printf("queue is empty\n")
		return nil
	}
	listenerData := data.(listener.ListenerData)

	msg := &PeerMsgRequest{}
	err = json.Unmarshal(listenerData.Body, msg)
	if err != nil {
		fmt.Printf("Unable to parse json\n")
		return nil
	}

	// verify signature here.

	//println("parsed peer message")
	return &msg.Body
}
