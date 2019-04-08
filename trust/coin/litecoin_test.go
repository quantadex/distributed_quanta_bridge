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

const LOCAL_RPC_HOST_LTC = "localhost:19332"

func SendLTC(address string, amount ltcutil.Amount) (string, error) {
	amountStr := fmt.Sprintf("%f", amount.ToBTC())
	fmt.Printf("Sending to %s amount of %s\n", address, amountStr)
	args := []string{
		//"-datadir=../../blockchain/bitcoin/data",
		"sendtoaddress",
		address,
		amountStr,
	}

	cmd := exec.Command("litecoin-cli", args...)
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

	cmd := exec.Command("litecoin-cli", args...)
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

	cmd := exec.Command("litecoin-cli", args...)
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

/**

address: LhG69oGwZFfF73KzQWDUshZzExousncNrG
private: 6uRzAi2FSxRVTuCVexQTzAPXJstTaqahcCrjd4T2vXob47YcyBv
public: 0488436B49B7DC34175D8074F6827D499C13E73A45783F8182C2ECB1F18C443ECE31CCC2C27375282A1539FF8519107ECF2EE02977191BF2D0A1AE86495A3CF8F9

address: LbLtvETwstn64nnubroJ1898Sfu1yeAGzD
private: 6vCii4KtPhyt5TYgwXQpMbAuYM92NJCncbH7VJ9sD6ZA5LNifDi
public: 0485A98A320E559ABF481F874521CDA6CE8A21691E4DF39B0389950C10F07A8CEAA25385B1E3D8807D4980F47082730C38E6A556E7B43F323FFCC6FC78487AEF7C
*/
func TestLTCEncodeRefund(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, []string{"047AABB69BBE1B5D9E2EFD10D0215A37AE835EAE08DFDF795E5A8411271F690CC8797CF4DEB3508844920E28A42A67D8A3F56D5B6B65401DEDB1E130F9F9908463", "04851D591308AFBE768566060C01A60A5F6AC6C78C3766559C835BEF0485628013ADC7D7E7676B0281FB83E788F4BC11E4CA597D1A53AF5F0BB90D555A28B55504"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	account := "pooja5"
	litecoin := client.(*LiteCoin)
	msig, err := litecoin.GenerateMultisig(account)

	litecoin.crosschainAddr = map[string]string{msig: account}
	fmt.Println(litecoin.crosschainAddr)

	amount, err := ltcutil.NewAmount(10.00)
	res, err := SendLTC(msig, amount)
	println(res, err)
	res, err = GenerateLTCBlock()
	println(res, err)

	addr1, err := litecoin.GenerateMultisig("crosschain2")

	log.Println("multisig: ", msig, err)

	GenerateLTCBlock()

	w := Withdrawal{
		SourceAddress:      msig,
		DestinationAddress: addr1,
		Amount:             1000,
		QuantaBlockID:      0,
	}
	tx, err := client.EncodeRefund(w)
	fmt.Println("tx = ", tx)
	assert.NoError(t, err)
	var encoded EncodedMsg
	json.Unmarshal([]byte(tx), &encoded)

	km, _ := key_manager.NewLiteCoinKeyManager(LOCAL_RPC_HOST, "regnet")

	err = km.LoadNodeKeys("92P5DpWDiuttphtXV5qrHjMnFU2nAyiR8NpyEkF5s8uAngVgBFb")
	assert.NoError(t, err)

	pub, _ := km.GetPublicKey()
	fmt.Println("loaded public key signer=", pub)

	tx_signed1, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	err = km.LoadNodeKeys("926mkZAmMowq4HaLqpNjwuJuPe3vP6iTVQnt1x9GWdwbnwQjDea")

	pub, _ = km.GetPublicKey()
	fmt.Println("loaded public key signer 2=", pub)

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
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)
}

func TestDepositsLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, nil)
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()

	_, err = client.GetDepositsInBlock(blockId, map[string]string{"QfLb3gv7pmgwsd4NBZkeZXY73i1QKiS3RW": "pooja"})
	assert.NoError(t, err)
}

func TestDecodeLTC(t *testing.T) {
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, []string{"0488436B49B7DC34175D8074F6827D499C13E73A45783F8182C2ECB1F18C443ECE31CCC2C27375282A1539FF8519107ECF2EE02977191BF2D0A1AE86495A3CF8F9", "0485A98A320E559ABF481F874521CDA6CE8A21691E4DF39B0389950C10F07A8CEAA25385B1E3D8807D4980F47082730C38E6A556E7B43F323FFCC6FC78487AEF7C"})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	litecoin := client.(*LiteCoin)
	addr1, err := litecoin.GenerateMultisig("crosschain2")
	ImportLTCAddress(addr1)
	addr2, err := litecoin.GenerateMultisig("token_sale")
	ImportLTCAddress(addr2)

	amount, err := ltcutil.NewAmount(0.1)
	res, err := SendLTC(addr1, amount)
	println(res, err)
	GenerateLTCBlock()

	crosschainAddr := make(map[string]string)
	crosschainAddr[addr1] = "pooja"
	litecoin.crosschainAddr = crosschainAddr
	//btec, err := crypto.NewGraphenePublicKeyFromString("QA5nvEN2S7Dej2C9hrLJTHNeMGeHq6uyjMdoceR74CksyApeZHWS")
	btec, err := crypto.GenerateGrapheneKeyWithSeed("pooja")
	assert.NoError(t, err)

	msig, err := litecoin.GenerateMultisig(btec)

	log.Println("multisig: ", msig, err)

	w := Withdrawal{
		SourceAddress:      msig,
		DestinationAddress: addr2,
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
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, []string{"0488436B49B7DC34175D8074F6827D499C13E73A45783F8182C2ECB1F18C443ECE31CCC2C27375282A1539FF8519107ECF2EE02977191BF2D0A1AE86495A3CF8F9", "0485A98A320E559ABF481F874521CDA6CE8A21691E4DF39B0389950C10F07A8CEAA25385B1E3D8807D4980F47082730C38E6A556E7B43F323FFCC6FC78487AEF7C"})
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
	client, err := NewLitecoinCoin(LOCAL_RPC_HOST_LTC, &chaincfg.RegressionNetParams, []string{"03AF8891DA9BBF3CED03F04BC3C17EC4D3AE61D464E9B89A6B6A1FA60E361FDEA4", "038CAFE50CA757FAD36DA592A7C2B19158C0163445BAC2DDF6A59BDDC8F5BF6AD1", "03F8C8D630BB53B2E08FB108E2A951C84E582BB3D585D2127FAE6DE43150A415AE"})
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
