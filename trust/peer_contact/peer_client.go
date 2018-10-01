package peer_contact

import (
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"github.com/spf13/viper"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"encoding/json"
	"net/http"
	"bytes"
	"errors"
	"fmt"
)

type PeerClient struct {
	q queue.Queue
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
	privateKey := viper.GetString("NODE_KEY")

	if signature := crypto.SignMessage(msg.Body, privateKey); signature != nil {
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

	msg := &PeerMsgRequest{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		fmt.Printf("Unable to parse json\n")
		return nil
	}
	//println("parsed peer message")
	return &msg.Body
}
