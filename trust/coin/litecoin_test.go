package coin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcutil"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
	"log"
	"os/exec"
	"testing"
)

func SendLTC(address string, amount ltcutil.Amount) (string, error) {
	amountStr := fmt.Sprintf("%f", amount.ToBTC())
	fmt.Printf("Sending to %s amount of %s\n", address, amountStr)
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"sendtoaddress",
		address,
		amountStr,
	}

	cmd := exec.Command("../../../../litecoin/./src/litecoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("err is here 1")
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func ImportLTCAddress(address string) {
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"importaddress",
		address,
	}

	cmd := exec.Command("../../../../litecoin/./src/litecoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("error is here 2")
		println("err", err.Error(), stderr.String())
	}
}

func GenerateLTCBlock() (string, error) {
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"generate",
		"1",
	}

	cmd := exec.Command("../../../../litecoin/./src/litecoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func TestLTCEncodeRefund(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.TestNet4Params, []string{"040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1", "049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	litecoin := client.(*LiteCoin)
	addr1, err := litecoin.GenerateMultisig("crosschain2")
	ImportLTCAddress(addr1)
	addr2, err := litecoin.GenerateMultisig("token_sale")
	fmt.Println("addresses = ", addr1, addr2)
	ImportLTCAddress(addr2)

	crosschainAddr := make(map[string]string)
	crosschainAddr[addr1] = "pooja"
	crosschainAddr[addr2] = "pooja"
	litecoin.crosschainAddr = crosschainAddr
	fmt.Println(litecoin.crosschainAddr)

	amount, err := ltcutil.NewAmount(0.01)
	res, err := SendLTC(addr1, amount)
	println(res, err)
	res, err = SendLTC(addr2, amount)
	println(res, err)
	res, err = GenerateBlock()
	println(res, err)

	//btec, err := crypto.NewGraphenePublicKeyFromString("QA5nvEN2S7Dej2C9hrLJTHNeMGeHq6uyjMdoceR74CksyApeZHWS")
	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)
	msig, err := litecoin.GenerateMultisig(btec)

	log.Println("multisig: ", msig, err)

	GenerateBlock()

	w := Withdrawal{
		SourceAddress:      "mn3SFr4mQctRqDbQDqwFMFDfzyzmhe6vxn",
		DestinationAddress: "2N8MWCvyNQzJqbcsdusUeLYTdgbQbepbuks",
		Amount:             1000,
		QuantaBlockID:      0,
	}
	tx, err := client.EncodeRefund(w)
	fmt.Println("tx = ", tx)
	assert.NoError(t, err)
	var encoded EncodedMsg
	json.Unmarshal([]byte(tx), &encoded)

	km, _ := key_manager.NewBitCoinKeyManager(LOCAL_RPC_HOST, "regnet")

	err = km.LoadNodeKeys("cNxQax7BfpbikeuCebPGCgTefTah5h1XhVDfaotVdFmXtaLCWLd9")
	assert.NoError(t, err)

	tx_signed1, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	err = km.LoadNodeKeys("cUixT9PYjTtNzcVjF8sB7iM9JeEf8tLHm9Wjgo972x8opCRNTasS")
	tx_signed2, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	fmt.Println(tx)
	fmt.Println(tx_signed1)
	fmt.Println(tx_signed2)

	w.Signatures = []string{tx_signed1, tx_signed2}
	hash, err := litecoin.SendWithdrawal(common.HexToAddress("0x0"), nil, &w)

	assert.NoError(t, err)
	fmt.Println("hash", hash, err)
}

func TestTopBlockIdLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)
}

func TestDepositsLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	//blockId, err := client.GetTopBlockID()

	_, err = client.GetDepositsInBlock(103, map[string]string{"QdRXeFHACBC9tDNc9AobaHQqMrNxgqZ1WT": "pooja"})
	assert.NoError(t, err)
}

func TestDecodeLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.TestNet4Params, []string{"03AF8891DA9BBF3CED03F04BC3C17EC4D3AE61D464E9B89A6B6A1FA60E361FDEA4", "038CAFE50CA757FAD36DA592A7C2B19158C0163445BAC2DDF6A59BDDC8F5BF6AD1", "03F8C8D630BB53B2E08FB108E2A951C84E582BB3D585D2127FAE6DE43150A415AE"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	litecoin := client.(*LiteCoin)

	crosschainAddr := make(map[string]string)
	crosschainAddr["mn3SFr4mQctRqDbQDqwFMFDfzyzmhe6vxn"] = "pooja"
	litecoin.crosschainAddr = crosschainAddr
	//btec, err := crypto.NewGraphenePublicKeyFromString("QA5nvEN2S7Dej2C9hrLJTHNeMGeHq6uyjMdoceR74CksyApeZHWS")
	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)

	msig, err := litecoin.GenerateMultisig(btec)

	log.Println("multisig: ", msig, err)

	w := Withdrawal{
		SourceAddress:      "mn3SFr4mQctRqDbQDqwFMFDfzyzmhe6vxn",
		DestinationAddress: "2NGYCnkuo62kL1QpAzV3bRaf747bSM8suQm",
		Amount:             1000,
		Tx:                 "4418603_0",
		QuantaBlockID:      0,
	}
	tx, err := client.EncodeRefund(w)
	fmt.Println("Encoded = ", tx)
	assert.NoError(t, err)

	_, err = client.DecodeRefund(tx)
	assert.NoError(t, err)
}

func TestEncodeWithMultipleInputsLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.TestNet4Params, []string{"03AF8891DA9BBF3CED03F04BC3C17EC4D3AE61D464E9B89A6B6A1FA60E361FDEA4", "038CAFE50CA757FAD36DA592A7C2B19158C0163445BAC2DDF6A59BDDC8F5BF6AD1", "03F8C8D630BB53B2E08FB108E2A951C84E582BB3D585D2127FAE6DE43150A415AE"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	litecoin := client.(*LiteCoin)

	addr1, err := litecoin.GenerateMultisig("crosschain2")
	addr2, err := litecoin.GenerateMultisig("token_sale")
	println(addr1, addr2)

	crosschainAddr := make(map[string]string)
	crosschainAddr[addr1] = "pooja"
	crosschainAddr[addr2] = "pooja"
	litecoin.crosschainAddr = crosschainAddr

	amount, err := ltcutil.NewAmount(0.01)
	res, err := SendLTC(addr1, amount)
	println(res, err)
	res, err = SendLTC(addr2, amount)
	println(res, err)
	res, err = GenerateLTCBlock()
	println(res, err)

	w := Withdrawal{
		SourceAddress:      "n2PNkvCSkkSKvgqLsQXAQACFETQwKvc16X",
		DestinationAddress: "2NGYCnkuo62kL1QpAzV3bRaf747bSM8suQm",
		Amount:             3000,
		Tx:                 "4418603_0",
		QuantaBlockID:      0,
	}

	_, err = client.EncodeRefund(w)
	assert.NoError(t, err)
}

/**
 * These are the public keys on testnet, and it failed to generate a key for some instances, fixed by adding more to the seed
 */
func TestGenerateMultisigLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST, &chaincfg.TestNet4Params, []string{"03AF8891DA9BBF3CED03F04BC3C17EC4D3AE61D464E9B89A6B6A1FA60E361FDEA4", "038CAFE50CA757FAD36DA592A7C2B19158C0163445BAC2DDF6A59BDDC8F5BF6AD1", "03F8C8D630BB53B2E08FB108E2A951C84E582BB3D585D2127FAE6DE43150A415AE"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)
	litecoin := client.(*LiteCoin)

	addr1, err := litecoin.GenerateMultisig("crosschain2")
	println(addr1, err)
	assert.NoError(t, err)

	addr2, err := litecoin.GenerateMultisig("token_sale")
	println(addr2, err)
	assert.NoError(t, err)
}
