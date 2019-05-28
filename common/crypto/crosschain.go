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

func GetBlackListedUsersByBlockcahin(blackListMap map[string][]string, blockchain string) map[string]bool {
	res := make(map[string]bool)
	if blackList, ok := blackListMap[blockchain]; ok {
		for _, users := range blackList {
			res[users] = true
		}
	}
	return res
}
