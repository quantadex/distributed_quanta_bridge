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
ISSUER_ADDRESS: QAXCTOY7IOZ3OEIY434CWTXBVGNZYVNTMVZSYRYQE5VDCH64BOJ2XYFM
NODE_KEY: ZBLOHXXVNJVEU7NHNWFHFIZ5PMD6TQQIY6MOTZI4GULH633D2XAXPYUT
HORIZON_URL: https://horizon-testnet.stellar.org
NETWORK_PASSPHRASE: QUANTA NETWORK
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
	dummy.AddDeposit(&coin.Deposit{"TEST", "QDCFARPB4ZR7VGTEL2XII5OPPUPPX2PQAYZURXRVR6Z34GNWTUHGVSXT",12345, 1})

	// run 1 epoch
	n1.cTQ.DoLoop()
	n2.cTQ.DoLoop()

	// run 2 epoch
	n1.cTQ.DoLoop()
	n2.cTQ.DoLoop()

	time.Sleep(time.Second * 12)
}