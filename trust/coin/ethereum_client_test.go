package coin

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCheckDepositNode(t *testing.T) {
	network := test.ETHER_NETWORKS[test.ROPSTEN]
	client := &Listener{NetworkID: network.NetworkId}
	ethereumClient, err := ethclient.Dial(network.Rpc)
	if err != nil {
		t.Error(err)
		return
	}

	client.Client = ethereumClient
	client.Start()

	const blockNumber = 4356013

	// https://ropsten.etherscan.io/tx/0x51a1018c6b2afd7bffb52d05178c3b66c9336e8c2dfeca0110f405cc41613492
	events, err := client.FilterTransferEvent(blockNumber, map[string]string{"0x527370def157bd1113db9448bd05e6402ffb5a0d": ""})
	if err != nil {
		println("block not found??")
		return
	}
	if len(events) != 1 {
		t.Error("Expecting 1 token transfer")
		return
	}
	fmt.Printf("%v\n", events)

	// https://ropsten.etherscan.io/tx/0x1af4e1eeaa9c2f823d7c2b37973cf8829896ea217b36bdcc1049ce1ff19504f2
	// expecting .15 ETH to 0xb59e4b94e4ed7331ee0520e9377967614ca2dc98
	deposits, err := client.GetNativeDeposits(4327101, map[string]string{strings.ToLower("0xb59e4b94e4ed7331ee0520e9377967614ca2dc98"): ""})
	if err != nil || len(deposits) != 1 {
		t.Error("Expecting 1 eth transfer", len(deposits))
		return
	}
}

func TestForwardScan(t *testing.T) {
	//to := "0xe0006458963c3773b051e767c5c63fee24cd7ff9"
	network := test.ETHER_NETWORKS[test.ROPSTEN]
	client := &Listener{NetworkID: network.NetworkId}
	ethereumClient, err := ethclient.Dial(network.Rpc)

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

func TestWithdrawalTX(t *testing.T) {
	km, _ := key_manager.NewEthKeyManager()
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

	w := &Withdrawal{
		TxId:               1,
		CoinName:           "ETH",
		DestinationAddress: dest.From.Hex(),
		QuantaBlockID:      1,
		Amount:             10000,
		Signatures:         nil,
	}

	coin := &EthereumCoin{}
	encoded, _ := coin.EncodeRefund(*w)
	println(encoded)
	signed, _ := km.SignTransaction(encoded)
	println("signed", signed)

	w.Signatures = []string{signed, signed, signed}

	client := &Listener{NetworkID: ROPSTEN_NETWORK_ID}
	tx, err := client.SendWithdrawal(sim, userAuth.From, userKey, w)

	if err != nil {
		println("ERR: ", err.Error())
	}
	println(tx)
}

/*

ENCODED:
00000000000000010000000000000000000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef0000000000000000000000000000000000000000000000000000000000003039
00000000000000010000000000000000000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef0000000000000000000000000000000000000000000000000000000000003039

SIGNED
6e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b901
6e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b901
6e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b901

0x6e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b901
v[hex] = 0x1c
r[hex] = 0x6e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e
s[hex] = 0x190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b9

0x497483d100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef000000000000000000000000000000000000000000000000000000000000303900000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000016e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e0000000000000000000000000000000000000000000000000000000000000001190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b9
0x497483d100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f17f52151ebef6c7334fad080c5704d77216b732000000000000000000000000000000000000000000000000000000000000303900000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000000168be485def688da668ed01391bc9185a15f74cbf3ea94c15b4ee27576dd8a71c00000000000000000000000000000000000000000000000000000000000000015d84e165130eb90f98f4ddd247107fc96ecab25fd85882fde086958a7e0b03b9
0x497483d100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef000000000000000000000000000000000000000000000000000000000000303900000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000016e5c5c2890328e7cf9c744fa8e6d90478a1fd7144fbe0fec6992cd5f8c26864e0000000000000000000000000000000000000000000000000000000000000001190f0f012abf3d8f222576a95622a0a9904a460a551b6b3e3671aecd1832f2b9
*/

func TestWithdrawalGanacheTX(t *testing.T) {
	w := &Withdrawal{
		TxId:               10000,
		CoinName:           "ETH",
		DestinationAddress: "0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef",
		QuantaBlockID:      1,
		Amount:             12345,
		Signatures:         nil,
	}

	network := test.ETHER_NETWORKS[test.ROPSTEN]
	coin, _ := NewEthereumCoin(network.NetworkId, network.Rpc)
	coin.Attach()
	encoded, _ := coin.EncodeRefund(*w)
	println(encoded)

	km, _ := key_manager.NewEthKeyManager()
	err := km.LoadNodeKeys(test.ROPSTEN_TRUST.NodeSecrets[0])
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

/**
 * Test that we can connect to ganache
 */
func TestGanacheTX(t *testing.T) {
	network := test.ETHER_NETWORKS[test.LOCAL]

	coin, err := NewEthereumCoin(network.NetworkId, network.Rpc)
	if err != nil {
		t.Error(err)
		return
	}
	err = coin.Attach()
	if err != nil {
		t.Error(err)
		return
	}
	blockId, err := coin.GetTopBlockID()
	println(blockId)

	// try to get out of bound
	coin.GetDepositsInBlock(50000, map[string]string{})
}

var cmd *exec.Cmd

func setupEthereum() {
	fmt.Println("Spinning up GETH")
	cmd = exec.Command("./run_ethereum.sh")
	cmd.Start()
}

func teardownEthereum() {
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("failed to kill process: %v", err)
	}
}

// https://www.philosophicalhacker.com/post/integration-tests-in-go/
func TestMain(m *testing.M) {
	if !testing.Short() {
		setupEthereum()
	}
	result := m.Run()

	if !testing.Short() {
		teardownEthereum()
	}
	os.Exit(result)
}
