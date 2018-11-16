package key_manager

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"math/big"
	"testing"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/stretchr/testify/assert"
)


func TestWithdrawalTX(t *testing.T) {
	km, _ := NewEthKeyManager()
	err := km.LoadNodeKeys("file://../../keystore/key--7cd737655dff6f95d55b711975d2a4ace32d256e")
	if err != nil {
		t.Error(err)
		return
	}

	key, _ := crypto.GenerateKey()
	dest := bind.NewKeyedTransactor(key)

	// simulate users creating the contract
	userKey, _ := crypto.GenerateKey()
	userAuth := bind.NewKeyedTransactor(userKey)

	sim := backends.NewSimulatedBackend(core.GenesisAlloc{
		userAuth.From: {Balance: big.NewInt(10000000000)}}, 5000000)

	w := &coin.Withdrawal{
		TxId:               1,
		CoinName:           "ETH",
		DestinationAddress: dest.From.Hex(),
		QuantaBlockID:      1,
		Amount:             10000,
		Signatures:         nil,
	}

	mycoin := &coin.EthereumCoin{}
	encoded, _ := mycoin.EncodeRefund(*w)
	fmt.Printf("encoded: %s\n", encoded)
	signed, _ := km.SignTransaction(encoded)
	fmt.Printf("signed: %s\n", signed)

	w.Signatures = []string{signed, signed, signed}

	client := &coin.Listener{NetworkID: coin.ROPSTEN_NETWORK_ID}
	tx, err := client.SendWithdrawal(sim, userAuth.From, userKey, w)

	if err != nil {
		println("ERR: ", err.Error())
	}
	println(tx)
}

func TestWithdrawalGanacheTX(t *testing.T) {
	w := &coin.Withdrawal{
		TxId:               10000,
		CoinName:           "ETH",
		DestinationAddress: "0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef",
		QuantaBlockID:      1,
		Amount:             12345,
		Signatures:         nil,
	}

	network := test.ETHER_NETWORKS[test.ROPSTEN]
	coin, _ := coin.NewEthereumCoin(network.NetworkId, network.Rpc)
	err := coin.Attach()
	assert.NoError(t, err)
	defer coin.Detach()
	encoded, _ := coin.EncodeRefund(*w)
	println(encoded)

	km, _ := NewEthKeyManager()
	err = km.LoadNodeKeys(test.ROPSTEN_TRUST.NodeSecrets[0])
	if err != nil {
		t.Error(err)
		return
	}

	println("encoded", encoded)
	signed, _ := km.SignTransaction(encoded)
	println("signed", signed)

	w.Signatures = []string{signed}

	// 0xcfed223fab2a41b5a5a5f9aaae2d1e882cb6fe2d <-- geth
	tx, err := coin.SendWithdrawal(common.HexToAddress(test.ROPSTEN_TRUST.TrustContract), km.GetPrivateKey(), w)
	if err != nil {
		println("ERR: ", err.Error())
	}
	println(tx)
}
