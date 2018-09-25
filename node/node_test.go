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

func SetConfig(port int) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	// any approach to require this configuration into your program.
	var config = []byte(fmt.Sprintf(`
LISTEN_IP: 0.0.0.0
LISTEN_PORT: %d
USE_PREV_KEYS: true
KV_DB_NAME: kv_db_%d
NODE_KEY: SAHHDUDOPJBEBPXA2ZTBL7QWYBQHAXH2Q3HWAEZEUTG6CEKX67UKAIDY

REGISTRAR_IP: localhost
REGISTRAR_PORT: 5001
`, port, port))

	viper.ReadConfig(bytes.NewBuffer(config))
}

func TestNode(t *testing.T) {
	os.Remove("./node/kv_db.db")

	SetConfig(5100)
	n1 := bootstrapNode()

	SetConfig(5101)
	n2 := bootstrapNode()

	//go n1.run()
	//go n2.run()
	time.Sleep(time.Second)

	dummy := coin.GetDummyInstance()
	dummy.CreateNewBlock()
	dummy.AddDeposit(&coin.Deposit{"TEST", "123",100, 1})

	// run 1 epoch
	n1.cTQ.DoLoop()
	n2.cTQ.DoLoop()

	time.Sleep(time.Second * 12)
}