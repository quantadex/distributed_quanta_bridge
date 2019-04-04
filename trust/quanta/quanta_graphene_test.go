package quanta

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const url = "ws://testnet-01.quantachain.io:8090"

func TestDynamicGlobalProperties(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "crosschain2"
	api.NetworkUrl = url
	api.Attach()
	_, err := api.GetTopBlockID()
	assert.NoError(t, err)
}

func TestGetBalances(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "crosschain2"
	api.NetworkUrl = url
	api.Attach()

	_, err := api.GetBalance("QDEX", "pooja")
	assert.NoError(t, err)

	_, err = api.GetAllBalances("quanta_foundation")
	assert.NoError(t, err)

}

func TestGetRefund(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "crosschain2"
	api.NetworkUrl = url
	api.Attach()

	_, _, err := api.GetRefundsInBlock(int64(29105), "quanta_foundation")
	assert.NoError(t, err)
}

func TestMultipleSignatures(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "crosschain2"
	api.NetworkUrl = url
	api.Attach()

	dep := &coin.Deposit{
		SenderAddr: "crosschain2",
		QuantaAddr: "pooja",
		Amount:     6000,
		CoinName:   "QDEX",
	}

	proposed, err := api.CreateTransferProposal(dep)
	//proposed, err := api.CreateNewAssetProposal("crosschain2", "TESTISSUE2", 5)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	km.LoadNodeKeys("5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx")
	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)

	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	// ready to submit to network
	_, err = api.Broadcast(submitTx)

	km.LoadNodeKeys("5KFJnRn38wuXnpKGvkxmsyiWUuUkPXKZGvdG8aTzHCTvJMUQ4sA")
	sig, err = km.SignTransaction(proposed)
	assert.NoError(t, err)

	submitTx, err = ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	// ready to submit to network
	_, err = api.Broadcast(submitTx)
	assert.NoError(t, err)
}

func TestCreateTransaction(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "pooja"
	api.NetworkUrl = url
	api.Attach()

	dep := &coin.Deposit{
		SenderAddr: "pooja",
		QuantaAddr: "crosschain2",
		Amount:     6000,
		CoinName:   "TESTISSUE2",
	}
	proposed, err := api.CreateTransferProposal(dep)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	km.LoadNodeKeys("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)

	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	// ready to submit to network
	_, err = api.Broadcast(submitTx)
	assert.NoError(t, err)
}

func TestDecodeTransaction(t *testing.T) {
	api := &QuantaGraphene{}
	api.Issuer = "crosschain2"
	api.NetworkUrl = url
	api.Attach()

	dep := &coin.Deposit{
		SenderAddr: "pooja",
		QuantaAddr: "crosschain2",
		Amount:     6000,
		CoinName:   "QDEX",
	}
	proposed, err := api.CreateTransferProposal(dep)
	assert.NoError(t, err)

	_, err = api.DecodeTransaction(proposed)
	assert.NoError(t, err)
}

func TestCreateAsset(t *testing.T) {
	api := QuantaGraphene{}
	api.NetworkUrl = url
	api.Issuer = "pooja"
	api.Attach()

	randomKey, _ := crypto.GenerateKey()
	randomAddr := crypto.PubkeyToAddress(randomKey.PublicKey)
	token := "SIMPLETOKEN" + strings.ToUpper(randomAddr.String())

	proposed, err := api.CreateNewAssetProposal(api.Issuer, token, 5)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	km.LoadNodeKeys("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)

	_, err = ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	_, err = api.DecodeTransaction(proposed)
	assert.NoError(t, err)

	// ready to submit to network
	//_, err = api.Broadcast(submitTx)
	//assert.NoError(t, err)

}

func TestIssueAsset(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "pooja"
	api.NetworkUrl = url
	api.Attach()

	dep := &coin.Deposit{
		SenderAddr: "pooja",
		QuantaAddr: "crosschain2",
		Amount:     10000,
		CoinName:   "ETHERTEST5",
	}

	proposed, err := api.CreateIssueAssetProposal(dep)
	assert.NoError(t, err)

	chainID, err := api.Database.GetChainID()
	assert.NoError(t, err)

	// sign it with keymanager
	km, err := key_manager.NewGrapheneKeyManager(*chainID)
	assert.NoError(t, err)

	km.LoadNodeKeys("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	sig, err := km.SignTransaction(proposed)
	assert.NoError(t, err)

	submitTx, err := ProcessGrapheneTransaction(proposed, []string{sig})
	assert.NoError(t, err)

	_, err = api.Broadcast(submitTx)
	assert.NoError(t, err)

}

func TestRandomMissRefund(t *testing.T) {
	api := QuantaGraphene{}
	api.Issuer = "tokensale"
	api.NetworkUrl = url
	api.Attach()

	for i := 0; i < 2; i++ {
		// call 100 within 100 random time
		r := rand.Intn(100)
		refunds, _, err := api.GetRefundsInBlock(int64(2351843), "tokensale")
		if err != nil {
			fmt.Println("error ", err)
		} else {
			fmt.Printf("Number of refunds # %d r=%d ms\n", len(refunds), r)
			assert.Equal(t, 1, len(refunds))
		}

		time.Sleep(time.Millisecond * time.Duration(r))
	}
}
