package coin

import (
	"testing"
	"fmt"
	"strings"
)

func TestEthereumEncodeDecode(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		1,
		"ETH",
		strings.ToLower("0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"),
		1,
		12345,
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
	if w.Amount != decoded.Amount {
		t.Error("Amount do not match")
	}
}


func TestEthereumEncodeDecodeERC20(t *testing.T) {
	coin := &EthereumCoin{}

	// test native ETH
	w := Withdrawal{
		1,
		strings.ToLower(",0xba7573C0e805ef71ACB7f1c4a55E7b0af4169999"),
		strings.ToLower("0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"),
		1,
		12345,
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
	if w.Amount != decoded.Amount {
		t.Error("Amount do not match")
	}
}