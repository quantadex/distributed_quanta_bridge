package coin

import (
	"testing"
	"github.com/ethereum/go-ethereum/ethclient"
	"fmt"
)

// https://goethereumbook.org/event-read-erc20/

// https://rinkeby.etherscan.io/tx/0xfe33cb4894c3acacab726b563f4efd22acd3d519b0583ce537176349a8fa9fcb#eventlog

const server = "testnet-02.quantachain.io:8545"
const toAddress = "0x0e742968b2804afa7b8729716fd0f695cbcd8631"
const blockNumber = 806159

//func TestNode(t *testing.T) {
//	client := &Listener{NetworkID: ROPSTEN_NETWORK_ID}
//	ethereumClient, err := ethclient.Dial("http://" + server)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	client.Client = ethereumClient
//	client.Start()
//
//	events, err := client.FilterTransferEvent(blockNumber, toAddress)
//	if err != nil {
//		println("block not found??")
//		return
//	}
//	if len(events) != 1 {
//		t.Error("Expecting 1 token transfer")
//		return
//	}
//	fmt.Printf("%v\n", events)
//
//	deposits, err := client.GetNativeDeposits(blockNumber, "0xc95293836769874de62d4c0f75fff0f05d39bc28")
//	if err != nil || len(deposits) != 1 {
//		t.Error("Expecting 1 eth transfer %d", len(deposits))
//		return
//	}
//
//	fmt.Printf("%v\n", deposits)
//}

func TestForwardScan(t *testing.T) {
	//to := "0xe0006458963c3773b051e767c5c63fee24cd7ff9"

	client := &Listener{NetworkID: ROPSTEN_NETWORK_ID}
	ethereumClient, err := ethclient.Dial("https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6")
	if err != nil {
		t.Error(err)
		return
	}

	client.Client = ethereumClient
	client.Start()
	contracts, err := client.GetForwardContract(4186074)
	if err != nil {
		println("err... " + err.Error())
		t.Error(err)
		return
	}

	fmt.Printf("%d %v %v\n", len(contracts), contracts[0].QuantaAddr, contracts[0].ContractAddress)
}
