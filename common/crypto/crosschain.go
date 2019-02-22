package crypto

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ForwardInput struct {
	ContractAddress string
	Trust           common.Address
	QuantaAddr      string
	TxHash          string
	Blockchain      string
}

type CrosschainAddress struct {
	Address    string
	QuantaAddr string
	TxHash     string
	Blockchain string
	Updated    time.Time
}
