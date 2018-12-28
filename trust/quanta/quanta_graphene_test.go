package quanta

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"testing"
)

func TestDynamicGlobalProperties(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()
	block, _ := api.GetTopBlockID("pooja")
	fmt.Println("TopBlock = ", block)

}

func TestGetBalances(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	balance, _ := api.GetBalance("QDEX", "pooja")
	fmt.Println("pooja balance = ", balance)

	balances, _ := api.GetAllBalances("quanta_foundation")
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
		Amount:     1,
		CoinName:   "QDEX",
	}
	proposed, err := api.CreateProposeTransaction(dep)
	if err != nil {
		fmt.Println(err)
	}
	chainID, err := api.Database.GetChainID()
	if err != nil {
		fmt.Println(err)
	}

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	km.LoadNodeKeys("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	sig, err := km.SignTransaction(proposed)
	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})

	// ready to submit to network
	err = api.Broadcast(submitTx)

	decoded, err := api.DecodeTransaction(proposed)
	fmt.Println(decoded, err)
}
