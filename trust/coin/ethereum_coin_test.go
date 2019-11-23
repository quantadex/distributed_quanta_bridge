package coin

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quantadex/distributed_quanta_bridge/common/test"
	"math"
	"strings"
	"testing"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/common"
	"encoding/hex"
)

func TestEthereumEncodeDecode(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		"some_long_tx_id",
		1,
		"ETH",
		"ETH",
		strings.ToLower("0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"),
		"0x000000000000000000000000000000000000000a",
		12345,
		5,
		nil,
		"",
	}
	fmt.Printf("original: %v\n", w)

	encoded, _ := coin.EncodeRefund(w)
	println(encoded)

	decoded, _ := coin.DecodeRefund(encoded)
	fmt.Printf("decoded: %v\n", decoded)

	if w.CoinName != decoded.CoinName {
		t.Error("Coin name do not match")
	}
	if w.DestinationAddress != decoded.DestinationAddress {
		t.Error("Destination do not match")
	}
	if uint64(math.Pow10(13)*float64(w.Amount)) != decoded.Amount {
		fmt.Println("decoded amount = ", decoded.Amount, decoded.DestinationAddress)
		t.Error("Amount do not match")
	}
}

//															00000000000000fb00000000000000000000000000000000000000006f7a21e074807c00893862e6a98b062c58799e07000000000000000000000000000000000000000000000000687fef1ab6b10000
//0x497483d100000000000000000000000000000000000000000000000000000000000000fb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000006f7a21e074807c00893862e6a98b062c58799e07000000000000000000000000000000000000000000000000687fef1ab6b1000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e00000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001b000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000001b0000000000000000000000000000000000000000000000000000000000000003dbea818f70fdb0b308f917385e7821bfaa955e6aa5ded8d4be7c0f965d711733008545d74392f06e23ecd7b8a5c068925b53c93794ab44077e6fee932eedf6b192b04bf88a3c538188e97f2f1d8555c589c04f83a76769a0fb7819fae570db8600000000000000000000000000000000000000000000000000000000000000033716c55cd358cecd2999a741c76f1d74684f59a7a4ab382a32db34edf15129d409d50366c6fb067e5e73429f46f62cbcf4679073875529b47dca5dd559586d16145b3a630b9d59cd67aa5a64cd48b5d29739fb0d8c58b5d63c076dfd6f1a687e

func TestEthereumEncodeDecode2(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		"some_long_tx_id",
		251,
		"ETH",
		"ETH",
		strings.ToLower("0x0"),
		"0x6f7a21e074807C00893862E6a98B062c58799e07",
		12345,
		753000,
		nil,
		"",
	}
	fmt.Printf("original: %v\n", w)

	encoded, _ := coin.EncodeRefund(w)
	println(encoded)
	msg := &EncodedMsg{}
	err := json.Unmarshal([]byte(encoded), msg)
	assert.NoError(t, err)

	km,_ := key_manager.NewEthKeyManager()
	km.LoadNodeKeys("0afd879321bf647a8ae9484e780916e681710255d733b3d5033710aa09ddfcd1")
	sign1,_ := km.SignTransaction(msg.Message)
	println(sign1)

	data := common.Hex2Bytes(sign1)
	if len(data) != 65 {
		fmt.Println("Signature is not correct length " + string(len(data)))
	}

	var r1 [32]byte
	copy(r1[0:32], data[0:32])

	var s1 [32]byte
	copy(s1[0:32], data[32:64])

	//v1 := data[64]+27

	//bb := bytes.NewBuffer(nil)
	//bb.WriteByte(v1)
	//bb.Write(r1[:])
	//bb.Write(s1[:])

	println(hex.EncodeToString(r1[:]))
	println(hex.EncodeToString(s1[:]))

	}

func TestEthereumEncodeDecodeERC20(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		"some_long_tx_id",
		1,
		",0Xc300ee2594fe0404a278f6ea81a024729843fa02",
		"ETH",
		strings.ToLower("0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"),
		"0x000000000000000000000000000000000000000a",
		12345,
		5,
		nil,
		"",
	}
	fmt.Printf("original: %v\n", w)

	network := test.ETHER_NETWORKS[test.ROPSTEN]
	listener := &Listener{NetworkID: network.NetworkId}
	ethereumClient, err := ethclient.Dial(network.Rpc)

	if err != nil {
		t.Error(err)
		return
	}

	coin.client = listener

	listener.Client = ethereumClient
	listener.Start(nil)

	encoded, _ := coin.EncodeRefund(w)
	println(encoded)

	decoded, err := coin.DecodeRefund(encoded)
	fmt.Printf("decoded: %v\n", decoded)

	if strings.ToLower(w.CoinName) != decoded.CoinName {
		t.Error("Coin name do not match")
	}
	if w.DestinationAddress != decoded.DestinationAddress {
		t.Error("Destination do not match")
	}
	if uint64(math.Pow10(4)*float64(w.Amount)) != decoded.Amount {
		t.Error("Amount do not match")
	}
}

func TestAdhocEncode(t *testing.T) {
	coin := &EthereumCoin{}

	//{ "txId": 1, "erc20Address": 0xf17f52151ebef6c7334fad080c5704d77216b732, "toAddress": 0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef, "amount": 1}
	w := Withdrawal{
		"some_long_tx_id",
		1,
		strings.ToLower(",0xf17f52151ebef6c7334fad080c5704d77216b732"),
		"ETH",
		strings.ToLower("0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef"),
		"a",
		1,
		5,
		nil,
		"",
	}
	encoded, _ := coin.EncodeRefund(w)
	println(encoded)
	fmt.Printf("original: %v\n", w)
}

func TestWithdrawal(t *testing.T) {

}
