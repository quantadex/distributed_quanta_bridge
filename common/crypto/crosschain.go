package crypto

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"
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
	fmt.Println("map = ", blackListMap["eth"], blackListMap["ETH"])
	if blackList, ok := blackListMap[strings.ToLower(blockchain)]; ok {
		for _, users := range blackList {
			res[users] = true
		}
	}
	return res
}
