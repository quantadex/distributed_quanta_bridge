package registrar_client

import (
	"github.com/spf13/viper"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
	"fmt"
	"net/http"
	"github.com/quantadex/distributed_quanta_bridge/registrar"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"errors"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"github.com/quantadex/distributed_quanta_bridge/common/queue"
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


func (r *RegistrarClient) RegisterNode(nodeIP string, nodePort string, nodeKey string) error {
	msg := registrar.RegisterReq{}
	msg.Body = registrar.NodeInfo{ nodeIP, nodePort, nodeKey }

	if signature := crypto.SignMessage(msg.Body, nodeKey); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(r.url + "/registry/api/register", "application/json", bytes.NewReader(data))
	}
	return errors.New("unable to sign message")
}

func (r *RegistrarClient) SendHealth(nodeState string, nodeKey string) error {
	msg := registrar.PingReq{}
	msg.Body = registrar.PingBody{ nodeState, nodeKey }

	if signature := crypto.SignMessage(msg.Body, nodeKey); signature != nil {
		msg.Signature = *signature

		data, err := json.Marshal(&msg)
		if err != nil {
			return errors.New("unable to marshall")
		}
		http.Post(r.url + "/registry/api/health", "application/json", bytes.NewReader(data))
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
