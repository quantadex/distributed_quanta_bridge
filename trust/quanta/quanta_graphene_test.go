package quanta

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestDynamicGlobalProperties(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()
	block, err := api.GetTopBlockID("pooja")
	assert.NoError(t, err)
	fmt.Println("TopBlock = ", block)

}

func TestGetBalances(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	balance, err := api.GetBalance("QDEX", "pooja")
	assert.NoError(t, err)
	fmt.Println("pooja balance = ", balance)

	balances, err := api.GetAllBalances("quanta_foundation")
	assert.NoError(t, err)
	fmt.Println("quanta_foundation balance = ", balances)

}

func TestGetRefund(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	refund, _, err := api.GetRefundsInBlock(int64(29105), "quanta_foundation")
	assert.NoError(t, err)
	fmt.Println("refund = ", refund)
}

func TestCreateTransaction(t *testing.T) {
	api := QuantaGraphene{}
	api.Attach()

	balance, err := api.GetBalance("QDEX", "pooja")
	assert.Equal(t, balance, 99687.10431)

	dep := &coin.Deposit{
		SenderAddr: "crosschain2",
		QuantaAddr: "pooja",
		Amount:     2,
		CoinName:   "QDEX",
	}
	proposed, err := api.CreateProposeTransaction(dep)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	//km.LoadNodeKeys("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	km.LoadNodeKeys("5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx")
	km.LoadNodeKeys("5KFJnRn38wuXnpKGvkxmsyiWUuUkPXKZGvdG8aTzHCTvJMUQ4sA")

	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)

	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	// ready to submit to network
	err = api.Broadcast(submitTx)
	assert.NoError(t, err)

	newbalance, err := api.GetBalance("QDEX", "pooja")
	assert.Equal(t, newbalance, math.Floor((balance+0.00002)*100000)/100000)

	decoded, err := api.DecodeTransaction(proposed)
	assert.NoError(t, err)
	fmt.Println(decoded, err)
}
