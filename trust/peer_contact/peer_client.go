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

}

func (p *PeerClient) AttachToListener() error {
	return nil
}

func (p *PeerClient) SendMsg(m *manifest.Manifest, destinationNodeID int, peermsg *PeerMessage) error {
	peer := m.Nodes[destinationNodeID]
	url := fmt.Sprintf("http://%s:%s", peer.IP, peer.Port)

	msg := &PeerMsgRequest{}
	msg.Body = *peermsg
	privateKey := viper.GetString("NODE_KEY")

	if signature := crypto.SignMessage(msg.Body, privateKey); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(url + "/node/api/peer", "application/json", bytes.NewReader(data))
	}
	return errors.New("unable to sign message")
}

func (p *PeerClient) GetMsg() *PeerMessage {
	q := queue.GetGlobalQueue()
	data, err := q.Get(queue.PEERMSG_QUEUE)
	if err != nil {
		return nil
	}

	msg := &PeerMessage{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		fmt.Printf("Unable to parse json")
		return nil
	}
	return msg
}
