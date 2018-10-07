package main

import (
	"testing"
	"os"
	"github.com/spf13/viper"
	"bytes"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"time"
)

var NODE_KEYS = []string {
	"ZBYEUJIWP2AXG2V6ZW4F5OTM5APW3SOTTM6YGMKO6MQSY7U3IHFJZHWQ",
	"ZAFYSHEOQIK67O6S6SD5X7PVTLULQH3WQ3AMAGOO4NHSRM5SIKWCWFZB",
	"ZC4U5P5DWNXGRUENOCOKZFHAWFKBE7JFOB2BCEKCM7BKXXKQE3DARXIJ",
}

//address:QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN weight:1
//address:QCNKL7QKKQZD63UW27JLY7LDLR6MME3WNLUJ47VP25EZH5THRPEZRSAK weight:1
//address:QCN2DWLVXNAZW6ALR6KXJWGQB4J2J5TBJVPYLQMIU2TDCXIOBID5WRU5 weight:1
//address:QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB weight:1 *** Issuer

func SetConfig(key string, port int) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var config = []byte(fmt.Sprintf(`
LISTEN_IP: 0.0.0.0
LISTEN_PORT: %d
USE_PREV_KEYS: true
KV_DB_NAME: kv_db_%d
COIN_NAME: ETH
ISSUER_ADDRESS: QAHXFPFJ33VV4C4BTXECIQCNI7CXRKA6KKG5FP3TJFNWGE7YUC4MBNFB
NODE_KEY: %s
HORIZON_URL: http://testnet-02.quantachain.io:8000/
NETWORK_PASSPHRASE: QUANTA Test Network ; September 2018
REGISTRAR_IP: localhost
REGISTRAR_PORT: 5001
ETHEREUM_NETWORK_ID: 3
ETHEREUM_RPC: testnet-02.quantachain.io:8545
`, port, port, key))

	viper.ReadConfig(bytes.NewBuffer(config))
}

func StartNodes(n int)[]*TrustNode {
	nodes := []*TrustNode{}

	for i := 0; i < n; i++ {
		os.Remove(fmt.Sprintf("./kv_db_%d.db", 5100+i))
		SetConfig(NODE_KEYS[i], 5100 + i)
		nodes = append(nodes, bootstrapNode())
	}

	return nodes
}

func DoLoop(nodes []*TrustNode) {
	for _, n := range nodes {
		n.cTQ.DoLoop()
	}
}

func TestNode(t *testing.T) {
	nodes := StartNodes(3)

	time.Sleep(time.Second)

	dummy := coin.GetDummyInstance()
	dummy.CreateNewBlock()

	// user generated ETH address from Quanta address
	// assume it is deposited to ETH address
	// 0x0f8d1c23a90795a7a738d90380ec8bb5e984ce9259b78ee7c5d1592253e4798c
	// with our system associating to QDCFARPB4ZR7VGTEL2XII5OPPUPPX2PQAYZURXRVR6Z34GNWTUHGVSXT
	dummy.AddDeposit(&coin.Deposit{"ETH",
					"QDCFARPB4ZR7VGTEL2XII5OPPUPPX2PQAYZURXRVR6Z34GNWTUHGVSXT",
					15*10000000, 1})

	DoLoop(nodes)
	DoLoop(nodes)
	DoLoop(nodes)

	time.Sleep(time.Second * 12)
}