package key_manager

import (
	"errors"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"log"
	"net"
	"net/rpc"
)

const DatabaseUrl = "DatabaseUrl"
const GrapheneSeedPrefix = "GrapheneSeedPrefix"
const EthereumKeyStore = "EthereumKeyStore"

type RemoteKeyManagerService struct {
	manager map[string]KeyManager
	service net.Listener
	server  *rpc.Server
}

type SignMessage struct {
	Chain   string
	Message string
}
type SignResponse struct {
	Signed string
}

type Signer struct {
	manager map[string]KeyManager
}

type Secrets struct {
	config common.Secrets
}

func (s *Secrets) GetSecretsWithoutKeys(param string, secrets *common.Secrets) error {
	*secrets = s.config
	return nil
}

func (h *Signer) GetSigners(blockchain string, res *[]string) error {
	if km, ok := h.manager[blockchain]; ok {
		*res = km.GetSigners()
		return nil
	} else {
		return errors.New("No chain found")
	}
}

func (h *Signer) GetPublicKey(chain string, pub *string) error {
	log.Println("GetPublicKey", chain)

	var err error
	if km, ok := h.manager[chain]; ok {
		*pub, err = km.GetPublicKey()
		return err
	} else {
		return errors.New("No chain found")
	}

	return nil
}

func (h *Signer) SignTx(message SignMessage, response *SignResponse) error {
	log.Println("SignTX", message.Chain, message.Message)

	if km, ok := h.manager[message.Chain]; ok {
		res, err := km.SignTransaction(message.Message)
		response.Signed = res
		return err
	} else {
		return errors.New("No chain found")
	}
}

func NewRemoteKeyManagerService(managers map[string]KeyManager, secrets common.Secrets) *RemoteKeyManagerService {
	service := RemoteKeyManagerService{}
	server := rpc.NewServer()
	err := server.Register(&Signer{managers})

	secrets.BchPrivateKey = ""
	secrets.BtcPrivateKey = ""
	secrets.LtcPrivateKey = ""
	secrets.NodeKey = ""
	err = server.Register(&Secrets{secrets})

	if err != nil {
		log.Fatal("Problem registering.", err)
	}
	service.server = server
	return &service
}

func (r *RemoteKeyManagerService) Serve(address string) {
	println("Serving on ", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("listen(%q): %s\n", address, err)
		return
	}

	go func() {
		for {
			cxn, err := ln.Accept()
			if err != nil {
				log.Printf("listen(%q): %s\n", address, err)
				return
			}
			//log.Printf("Server accepted connection to %s from %s\n", cxn.LocalAddr(), cxn.RemoteAddr())
			go r.server.ServeConn(cxn)
		}
	}()
}
