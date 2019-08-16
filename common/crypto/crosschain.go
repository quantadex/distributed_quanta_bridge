package crypto

import (
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

type ForwardInput struct {
	ContractAddress string
	Trust           common.Address
	QuantaAddr      string
	TxHash          string
	Blockchain      string
}

//type CrosschainAddress struct {
//	Address    string
//	QuantaAddr string
//	TxHash     string
//	Blockchain string
//	Updated    time.Time
//}

func GetBlackListedUsersByBlockchain(blackListMap map[string][]string, blockchain string) map[string]bool {
	res := make(map[string]bool)
	if blackList, ok := blackListMap[strings.ToLower(blockchain)]; ok {
		for _, users := range blackList {
			res[users] = true
		}
	}
	return res
}
