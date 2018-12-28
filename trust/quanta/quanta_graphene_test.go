package quanta

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/scorum/bitshares-go/apis/database"
	"github.com/scorum/bitshares-go/apis/login"
	"github.com/scorum/bitshares-go/transport/websocket"
	"github.com/stretchr/testify/require"
	"testing"
)

const url = "ws://testnet-01.quantachain.io:8090"

func getAPI(t *testing.T) *database.API {
	transport, err := websocket.NewTransport(url)
	require.NoError(t, err)

	// request access to the database api
	databaseAPIID, err := login.NewAPI(transport).Database()
	require.NoError(t, err)
	fmt.Println(databaseAPIID)
	api := database.NewAPI(databaseAPIID, transport)
	fmt.Println(api)
	return api
}

func TestDynamicGlobalProperties(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()
	block, _ := api.GetTopBlockID("pooja")
	fmt.Println("TopBlock = ", block)

	//decoded, err := api.DecodeTransaction("c4cbaa5b1127ac692e23c04566196f3c858f0bec")
	//fmt.Println(decoded, err)

}

func TestGetRecentTransactionByID(t *testing.T) {
	databaseAPI := getAPI(t)
	//trx, err := databaseAPI.GetRecentTransactionByID(0)
	trx, err := databaseAPI.GetTransaction(29105, 0)
	fmt.Println("trx = ", trx)
	require.NoError(t, err)
	require.NotNil(t, trx)
}

func TestGetBalances(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	balance, _ := api.GetBalance("QDEX", "pooja")
	fmt.Println("pooja balance = ", balance)

	balances, _ := api.GetAllBalances("quanta_foundation", "QDEX", "ETH", "USD")
	fmt.Println("quanta_foundation balance = ", balances)

}

func TestGetRefund(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	refund, _, _ := api.GetRefundsInBlock(int64(29105), "quanta_foundation")
	fmt.Println("refund = ", refund)
}

func TestCreateTransaction(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	dep := &coin.Deposit{
		SenderAddr: "pooja",
		QuantaAddr: "quanta_foundation",
		Amount:     6000,
		CoinName:   "QDEX",
	}
	result, err := api.CreateProposeTransaction(dep)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
