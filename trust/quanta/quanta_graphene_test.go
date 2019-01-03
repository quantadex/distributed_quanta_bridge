package quanta

import (
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
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

	//balance, err := api.GetBalance("QDEX", "pooja")
	//assert.Equal(t, balance, 99687.10431)

	dep := &coin.Deposit{
		SenderAddr: "pooja",
		QuantaAddr: "quanta_foundation",
		Amount:     5000,
		CoinName:   "QDEX",
	}
	fmt.Println(dep)
	//proposed, err := api.CreateProposeTransaction(dep)
	//proposed, err := api.IssueAssetPropose()
	proposed, err := api.AssetProposeTransaction()
	fmt.Println("Proposed = ", proposed)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)
	fmt.Println("chainId = ", *chainID)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	km.LoadNodeKeys("5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx")
	//km.LoadNodeKeys("5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx")
	//km.LoadNodeKeys("5KFJnRn38wuXnpKGvkxmsyiWUuUkPXKZGvdG8aTzHCTvJMUQ4sA")

	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)
	fmt.Println(sig)

	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)
	fmt.Println(submitTx)

	// ready to submit to network
	err = api.Broadcast(submitTx)
	assert.NoError(t, err)
	fmt.Println(err)

	//newBalance, err := api.GetBalance("QDEX", "pooja")
	//assert.Equal(t, newBalance, math.Floor((balance+0.00002)*100000)/100000)

	//decoded, err := api.DecodeTransaction(proposed)
	//assert.NoError(t, err)
	//fmt.Println(decoded, err)
}
