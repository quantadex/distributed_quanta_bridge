package coin

import (
	"fmt"
	"math"
	"strings"
	"testing"
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

func TestEthereumEncodeDecodeERC20(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		"some_long_tx_id",
		1,
		",0Xba7573C0e805ef71ACB7f1c4a55E7b0af4169999",
		"ETH",
		strings.ToLower("0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"),
		"0x000000000000000000000000000000000000000a",
		12345,
		5,
		nil,
	}
	fmt.Printf("original: %v\n", w)

	encoded, _ := coin.EncodeRefund(w)
	println(encoded)

	decoded, _ := coin.DecodeRefund(encoded)
	fmt.Printf("decoded: %v\n", decoded)

	if strings.ToLower(w.CoinName) != decoded.CoinName {
		fmt.Println("coin name :", decoded.CoinName, decoded.Amount)
		t.Error("Coin name do not match")
	}
	if w.DestinationAddress != decoded.DestinationAddress {
		t.Error("Destination do not match")
	}
	if uint64(math.Pow10(13)*float64(w.Amount)) != decoded.Amount {
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
	}
	encoded, _ := coin.EncodeRefund(w)
	println(encoded)
	fmt.Printf("original: %v\n", w)
}

func TestWithdrawal(t *testing.T) {

}
