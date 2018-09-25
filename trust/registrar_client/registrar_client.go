package registrar_client

import (
	"github.com/spf13/viper"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"fmt"
	"net/http"
	"errors"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
	"github.com/quantadex/distributed_quanta_bridge/common/msgs"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
)

type RegistrarClient struct{
	address string
	port int
	url string
	healthQueueName string
}

func (r *RegistrarClient) GetRegistrar() error {
	r.address = viper.GetString("REGISTRAR_IP")
	r.port = viper.GetInt("REGISTRAR_PORT")
	r.url = fmt.Sprintf("http://%s:%d", r.address, r.port)
	r.healthQueueName = viper.GetString("HEALTH_QUEUE")

	return nil
}

/**
 * Listen to the node's calls
 */
func (r *RegistrarClient) AttachToListener() error {
	// memory queue, not necessary
	return nil
}


func (r *RegistrarClient) RegisterNode(nodeIP string, nodePort string, km key_manager.KeyManager) error {
	msg := msgs.RegisterReq{}
	nodeKey, _ := km.GetPublicKey()
	msg.Body = msgs.NodeInfo{ nodeIP, nodePort, nodeKey }
	fmt.Printf("Send to %s %s\n", r.url, nodeKey)

	if signature := km.SignMessageObj(msg.Body); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(r.url + "/registry/api/register", "application/json", bytes.NewReader(data))
		return nil
	}
	return errors.New("unable to sign message")
}

func (r *RegistrarClient) SendHealth(nodeState string, km key_manager.KeyManager) error {
	msg := msgs.PingReq{}
	nodeKey, _ := km.GetPublicKey()
	msg.Body = msgs.PingBody{ nodeState, nodeKey }

	if signature := km.SignMessageObj(msg.Body); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(r.url + "/registry/api/health", "application/json", bytes.NewReader(data))
		return nil
	}

	return errors.New("unable to sign message")
}

func (r *RegistrarClient) GetManifest() *manifest.Manifest {
	res, err := http.Get(r.url + "/registry/api/manifest")
	if err != nil {
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	manifest, err := manifest.CreateManifestFromJSON(bodyBytes)
	if err != nil {
		return nil
	}

	return manifest
}

func (r *RegistrarClient) HealthCheckRequested() bool {
	q := queue.GetGlobalQueue()
	_, err := q.Get(queue.HEALTH_QUEUE)
	if err != nil {
		return false
	}
	return true
}
