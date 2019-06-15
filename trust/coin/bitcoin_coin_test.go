package coin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/key_manager"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const LOCAL_RPC_HOST = "localhost:18332"

/*
cSJ2vqDoT9p6PXqNdNzaLNtMKePVjEzvoAryUUN7qAAB4njLKMXa 2018-10-29T23:34:06Z
reserve=1 # addr=2N39mAkxmLNnnL9WYecjkkTTtHTVQ3RtfZx hdkeypath=m/0'/0'/5'


cUgyLdmWgiMZcnCgTnmV1ag2evz5Eid6HQVacqfXPpWUbxgJcGt6 2018-10-29T23:34:06Z
reserve=1 # addr=2NF63kkxcegxtMuTKartK4tsyXsoHxhRvpN hdkeypath=m/0'/0'/13'

*/

func TestCheckHash(t *testing.T) {
	scriptBytes, _ := hex.DecodeString("004730440220304a3f60b7f5510e80b086cee9e88e38672c6031c3c2905b39bf5b180ba463b602205a14fcffee23eda01193e4dd281be474b058cc5d58b523f6bbc882469e05647c01473044022036224d12535cb02d597e1e1d0a6baccca9c27531464935110bab23c7a40f0cdd02206de86e931083484f6c23da41455863e6ff0c6b1359f8edf5f54364b83deea39c0147522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e52ae")
	addr, _ := btcutil.NewAddressScriptHash(scriptBytes, &chaincfg.RegressionNetParams)
	println(addr.EncodeAddress(), addr.String())
}

func TestBitcoinEncodeRefund(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"}, "user", "123", "", 0.00075, 0.00025, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	bitcoin := client.(*BitcoinCoin)
	blockId, err := bitcoin.GetTopBlockID()
	assert.NoError(t, err)

	if blockId < 101 {
		bitcoin.Client.Generate(101)
	}
	addr1, err := bitcoin.GenerateMultisig("aaa1")
	assert.NoError(t, err)
	addr2, err := bitcoin.GenerateMultisig("2")
	assert.NoError(t, err)
	println(addr1, addr2)

	crosschainAddr := make(map[string]string)
	crosschainAddr[addr1] = "pooja"
	crosschainAddr[addr2] = "pooja"
	bitcoin.crosschainAddr = crosschainAddr
	fmt.Println(bitcoin.crosschainAddr)

	amount, err := btcutil.NewAmount(1.0)
	bitcoinAddr, err := btcutil.DecodeAddress(addr1, &chaincfg.RegressionNetParams)
	bitcoin.Client.SendToAddress(bitcoinAddr, amount)

	bitcoinAddr, err = btcutil.DecodeAddress(addr2, &chaincfg.RegressionNetParams)
	bitcoin.Client.SendToAddress(bitcoinAddr, amount)

	bitcoin.Client.Generate(1)

	w := Withdrawal{
		SourceAddress:      addr2,
		DestinationAddress: addr1,
		Amount:             100,
		QuantaBlockID:      0,
	}
	fee, totalFee, err := bitcoin.estimateFee(2, 2)
	fmt.Printf("fee %f %f %v\n", fee, totalFee, err)

	tx, err := client.EncodeRefund(w)
	fmt.Println("tx = ", tx, err)
	assert.NoError(t, err)

	var encoded EncodedMsg
	json.Unmarshal([]byte(tx), &encoded)

	km, _ := key_manager.NewBitCoinKeyManager(LOCAL_RPC_HOST, "regnet", "user", "123")

	err = km.LoadNodeKeys("92REaZhgcw6FF2rz8EnY1HMtBvgh3qh4gs9PxnccPrju6ZCFetk")
	assert.NoError(t, err)

	tx_signed1, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	err = km.LoadNodeKeys("923EhimzuuHQvRaRWhTbKtocZSaKjvXkc32jbBiT5NPkCVGKYmf")
	tx_signed2, err := km.SignTransaction(encoded.Message)
	assert.NoError(t, err)

	fmt.Println(tx)
	fmt.Println(tx_signed1)
	fmt.Println(tx_signed2)

	w.Signatures = []string{tx_signed1, tx_signed2}
	hash, err := bitcoin.SendWithdrawal(common2.HexToAddress("0x0"), nil, &w)

	assert.NoError(t, err)
	fmt.Println("hash", hash, err)
}

func TestTopBlockId(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, nil, "user", "123", "", 0.00002, 0.00001, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()
	assert.NoError(t, err)
	fmt.Println(blockId)
}

func TestDeposits(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, nil, "user", "123", "", 0.00002, 0.00001, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	blockId, err := client.GetTopBlockID()

	_, err = client.GetDepositsInBlock(blockId, nil)
	assert.NoError(t, err)
}

func TestDecode(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"}, "user", "123", "", 0.00075, 0.00025, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	bitcoin := client.(*BitcoinCoin)
	blockId, err := bitcoin.GetTopBlockID()
	assert.NoError(t, err)
	if blockId < 101 {
		bitcoin.Client.Generate(101)
	}
	addr1, err := bitcoin.GenerateMultisig("aaa1")
	assert.NoError(t, err)

	amount, err := btcutil.NewAmount(0.02)
	assert.NoError(t, err)

	bchAddr, err := btcutil.DecodeAddress(addr1, &chaincfg.RegressionNetParams)
	assert.NoError(t, err)

	_, err = bitcoin.Client.SendToAddress(bchAddr, amount)
	assert.NoError(t, err)

	bitcoin.Client.Generate(1)

	crosschainAddr := make(map[string]string)

	crosschainAddr[addr1] = "pooja"
	bitcoin.crosschainAddr = crosschainAddr

	w := Withdrawal{
		SourceAddress:      addr1,
		DestinationAddress: "2N3Zj2iCe2YuZD7sXRLD6yvAHiz318NTiae",
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

func TestEncodeWithMultipleInputs(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.RegressionNetParams, []string{"049C8C4647E016C502766C6F5C40CFD37EE86CD02972274CA50DA16D72016CAB5812F867F27C268923E5DE3ADCB268CC8A29B96D0D8972841F286BA6D9CCF61360", "040C9B0D5324CBAF4F40A215C1D87DF1BEB51A0345E0384942FE0D60F8D796F7B7200CC5B70DDCF101E7804EFA26A0CE6EC6622C2FE90BCFD2DA2482006C455FF1"}, "user", "123", "", 0.00075, 0.00025, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)

	bitcoin := client.(*BitcoinCoin)
	blockId, err := bitcoin.GetTopBlockID()
	assert.NoError(t, err)
	if blockId < 101 {
		bitcoin.Client.Generate(101)
	}

	addr1, err := bitcoin.GenerateMultisig("aaa1")
	assert.NoError(t, err)
	addr2, err := bitcoin.GenerateMultisig("2")
	assert.NoError(t, err)
	addr3, err := bitcoin.GenerateMultisig("crosschain2")
	assert.NoError(t, err)

	crosschainAddr := make(map[string]string)

	crosschainAddr[addr1] = "pooja"
	crosschainAddr[addr2] = "pooja"
	bitcoin.crosschainAddr = crosschainAddr

	amount, err := btcutil.NewAmount(0.02)
	bchAddr, err := btcutil.DecodeAddress(addr1, &chaincfg.RegressionNetParams)
	assert.NoError(t, err)

	_, err = bitcoin.Client.SendToAddress(bchAddr, amount)
	assert.NoError(t, err)

	bchAddr, err = btcutil.DecodeAddress(addr1, &chaincfg.RegressionNetParams)
	_, err = bitcoin.Client.SendToAddress(bchAddr, amount)
	assert.NoError(t, err)
	bitcoin.Client.Generate(1)

	w := Withdrawal{
		SourceAddress:      addr2,
		DestinationAddress: addr3,
		Amount:             1000,
		Tx:                 "4418603_0",
		QuantaBlockID:      0,
	}

	encoded, err := client.EncodeRefund(w)
	fmt.Println(encoded)
	assert.NoError(t, err)
}

/**
 * These are the public keys on testnet, and it failed to generate a key for some instances, fixed by adding more to the seed
 */
func TestGenerateMultisig(t *testing.T) {
	client, err := NewBitcoinCoin(LOCAL_RPC_HOST, &chaincfg.TestNet3Params, []string{"03AF8891DA9BBF3CED03F04BC3C17EC4D3AE61D464E9B89A6B6A1FA60E361FDEA4", "038CAFE50CA757FAD36DA592A7C2B19158C0163445BAC2DDF6A59BDDC8F5BF6AD1", "03F8C8D630BB53B2E08FB108E2A951C84E582BB3D585D2127FAE6DE43150A415AE"}, "user", "123", "", 0.00002, 0.00001, map[string]bool{})
	assert.NoError(t, err)

	err = client.Attach()
	assert.NoError(t, err)
	bitcoin := client.(*BitcoinCoin)

	start := time.Now()
	for i :=0; i < 100; i++ {
		addr1, err := bitcoin.GenerateMultisig("crosschainx_" + string(i))
		println(addr1, err)
		assert.NoError(t, err)
	}
	end := time.Now()
	fmt.Printf("Time it took=%f seconds", end.Sub(start).Seconds())
}

//not needed. testdeocde already exists
//func TestDecodeRefund(t *testing.T) {
//	client, err := NewBitcoinCoin(LOCAL_RPC_HOST,&chaincfg.RegressionNetParams, []string{"2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", "2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR"})
//	assert.NoError(t, err)
//
//	err = client.Attach()
//	assert.NoError(t, err)
//
//	bitcoin := client.(*BitcoinCoin)
//	addr1, err := bitcoin.GenerateMultisig("aaa1")
//	addr2, err := bitcoin.GenerateMultisig("2")
//	println(addr1, addr2)
//
//	crosschainAddr := make(map[string]string)
//	crosschainAddr[addr1] = "pooja"
//	crosschainAddr[addr2] = "pooja"
//	bitcoin.crosschainAddr = crosschainAddr
//
//	var res common.TransactionBitcoin
//	res.Tx = "02000000084d8f3b5409f01500e832767d1e208739057c196666430b38135acf208c5ad61501000000fd1c01004730440220636998698dac64a2b266b9272f42ba8d97b116612fd4b572d09b2a44a73f2acf0220108a5987b05f56fc635c02ba6182966d8e948d03bbf24ebf9fe642ad938305580147304402203f65b0321c1a692a4a3c9a1c0fdeccdbc4041bfb384d8f2e7e1350f8439380e202202b8d09242d0018532d11e0748217c13c1e8e87ada47d39b1113425fa49679ca4014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104b06e6b02c2d8d89075d19b65a8aa12075a6ca5e96d191aff73ce5242f428e0d977ff4303a70938ca2211298cd1df1b463aaaa3609bdbaa46e3948698887c90ea53aeffffffff7ee7efdf2725d1fe9007dd0afb3ab99aa04e3df3af93b3ccda1ad4226d6fb40f00000000fd1c0100473044022073fc161a63da95a40b25825707a733ba7b36b2325512a6441c4a9ec20944de5302207f1ecc6e023057c9ded4916458dfa05f44ba4b102e4fc98bf6fa1c72fb534c8c0147304402200b092108d7da7fb9f602cfc4b4591b796415ac242bbbc18ca8924e1731ce67350220299fe3d5af26d12788dbf40259169e4bd8ca7e91b273b5d9cb72338f00e13709014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35e827898e61645fde069c2f40b6b1c10025758ba5021641ba2df796050fa9dc8253aeffffffffcab3781449ccdccc2cace24b5bb57b8ca99fde293059111758ae739410c6e51900000000fd1c0100473044022016d354e3aa88985e6294dde09760498e177da7e1dc0fae3c87d22b1638180d5702202828d9bd0d9238e309f59d248b08636fa16ffea66787c80d39247038cc5a1cf60147304402207a0f394cd111edd05c2963bf5fbcef423dd16a70b5b46dd1c08913680f25b76e0220408f6f5c811ef717e1d731a167f72d95164def522c71c523a2e78dc1a45e0817014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104b06e6b02c2d8d89075d19b65a8aa12075a6ca5e96d191aff73ce5242f428e0d977ff4303a70938ca2211298cd1df1b463aaaa3609bdbaa46e3948698887c90ea53aeffffffffdc5e4a50bdb85fb4d6ab0f15da7908a269546757efd378921f8ddb3ea8d0976e01000000fd1c0100473044022061737c60d1432fd6e6c2587ebf236bfe7fab29f70463ff3f78beeed746de26480220503781b3c3d3a433af75d62d4d477dbbe45bcdbbcd22fcc4fcc0ee694d019b37014730440220666da18508dbdabfed040d63c76c2e49efbbd7a566a01bb1ab346a678cac24e0022033b0ae0ab133a75c413b24c5494b2b518d628e75bafae5a8446a3d214bf8dfe2014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35e827898e61645fde069c2f40b6b1c10025758ba5021641ba2df796050fa9dc8253aeffffffffec112e2cc19f3a7d16843b0dcee45bef7f50631f38277733308259ec6d204bfb00000000fd1c0100473044022019b99fedf1a79e912074d40ecafd70c2cfa061af43ee9a87d8552bb21ef462c102200dc989f01e99d3ac4b480f7b60c6815f75d7ae28343509722f985d7f585e1fb301473044022003375eeccf9830b4a689c5c97e433a040e3d1eb209cdca69be65578fe9a35ae9022059fea97c91ce86a9efd36ae46178b91749415a8af0f3f140fae7850209a4c918014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104b06e6b02c2d8d89075d19b65a8aa12075a6ca5e96d191aff73ce5242f428e0d977ff4303a70938ca2211298cd1df1b463aaaa3609bdbaa46e3948698887c90ea53aeffffffffefd8029acdec3332340dba61f5292f07c75c0603763b5b92d30b60da73f8d6ab00000000fd1c0100473044022020fd3cbf8beb216da931cdce80669174155d7cb1cc083680ae588e21e5c820ad02201d3a2c001da43b8d7fca4b63b131276461de9682b4775f1cc48177447ca7b4ab0147304402200b98e2479e1c85cd3bb3cc18730e5e142e9566e17cf52f99506d2f84860926d80220497f67fb5273fce80ca62fe668c55c7d8e2d67d58adb2b7c3527b3e3813ae582014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35e827898e61645fde069c2f40b6b1c10025758ba5021641ba2df796050fa9dc8253aefffffffff2f336f911e12a08a74c56dd3a126ff02cde44fbc633fbbcc88f5b70022a791900000000fd1c010047304402206a54118e4a7919164db0bbc2f0a0a69bb386478e8f46a73f9994739c6ba5c74f02200e7553e2750a2d4121692fd06517761b4f67d9ddf1c7c692ea122cfcd98c6ca401473044022005120e80343a0881485ab59c2cf4b4c602086462db79fdb436d6b15c60662d3a02205a6c7f8d599d915fc1e6f42c5c86811d098cae2ff276d0ba21001512b9a21d10014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104b06e6b02c2d8d89075d19b65a8aa12075a6ca5e96d191aff73ce5242f428e0d977ff4303a70938ca2211298cd1df1b463aaaa3609bdbaa46e3948698887c90ea53aefffffffff7246e1dd3c7f6401e8f3e10863ee2b2a846d58ef3c5c300cf369cf96de7522800000000fd1c0100473044022077850b910589bfe395df60441bdf3755079e442906b064cb64626a3eec05ab5a0220208ad5023513e150f20a8a7980899b11d4946a39c299be5f50a0284f68f45fc70147304402203260c3248d748d6cee239cc44a01ec29cec46db313aadb29241a0ecff62d2fc202204c8f2def3aadfe58447a1516c2b6694747a20a541735fbb612cb132e216fe45c014c89522103c19460f565d12512ee584685bd8d97eb24d79a2acdf0c5b6af0b24ba29ceba0b210333d415aed3103f49346a3898efa137c42e933bbb16b3e4b56f7751670d2e0b6e4104b06e6b02c2d8d89075d19b65a8aa12075a6ca5e96d191aff73ce5242f428e0d977ff4303a70938ca2211298cd1df1b463aaaa3609bdbaa46e3948698887c90ea53aeffffffff022d360a000000000017a914a04c457920deea71e09aaf929c92ae3966aae25087405489000000000017a914bdf488ab1de8747560ff66dbb31eb91f5d66f59c8700000000"
//	res.RawInput = nil
//
//	resStr, err := json.Marshal(res)
//	data, err := json.Marshal(&EncodedMsg{string(resStr), "", 0, ""})
//
//	_, err = client.DecodeRefund(string(data))
//	assert.NoError(t, err)
//}
