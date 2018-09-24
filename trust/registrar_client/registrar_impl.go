package registrar_client

import (
	"github.com/spf13/viper"
	"github.com/quantadex/distributed_quanta_bridge/common/manifest"
)

type RegistrarClient struct{
	address string
	port int
}

func (r *RegistrarClient) GetRegistrar() error {
	r.address = viper.GetString("REGISTRAR_IP")
	r.port = viper.GetInt("REGISTRAR_PORT")

	return nil
}

/**
 * Listen to the node's calls
 */
func (r *RegistrarClient) AttachToListener() error {
	panic("implement me")
}

func (r *RegistrarClient) RegisterNode(nodeIP string, nodePort string, nodeKey string) error {
	panic("implement me")
}

func (r *RegistrarClient) SendHealth(nodeState string) error {
	panic("implement me")
}

func (r *RegistrarClient) GetManifest() *manifest.Manifest {
	panic("implement me")
}

func (r *RegistrarClient) HealthCheckRequested() bool {
	panic("implement me")
}
