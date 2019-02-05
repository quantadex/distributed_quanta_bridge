package coin

import (
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"bytes"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"github.com/btcsuite/btcd/wire"
	"os/exec"
	"github.com/pkg/errors"
	"strings"
)

type BitcoinCoin struct {
	Client *rpcclient.Client
	chaincfg *chaincfg.Params
	command string
}

func (c *BitcoinCoin) Blockchain() string {
	return "BTC"
}

func (b *BitcoinCoin) Attach() error {
	b.chaincfg = &chaincfg.RegressionNetParams
	b.command = "-datadir=blockchain/bitcoin/data/"
	var  err error
	b.Client, err = rpcclient.New(&rpcclient.ConnConfig{ Host: "localhost:18332",
															Endpoint:"http",
																User: "user",
																Pass: "123",
																DisableTLS: true,
																HTTPPostMode: true,
																	}, nil)
	return err
}

func (b *BitcoinCoin) GenerateMultisig(addresses []btcutil.Address) (string, error){
	addr := []string{}
	for _, a := range addresses {
		addr = append(addr, a.String())
	}

	addrx, err := b.Client.AddMultisigAddress(2, addresses, "")
	fmt.Println("result ", addrx)

	scriptBytes, err := hex.DecodeString(addrx.RedeemScript)
	if err != nil {
		return "", err
	}

	res,err := b.Client.DecodeScript(scriptBytes)

	if err != nil {
		return "", err
	}

	return res.P2sh, err
}

func (b *BitcoinCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	panic("not needed for bitcoin")
}

func (b *BitcoinCoin) GetTopBlockID() (int64, error) {
	//	b.client.GetBestBlock()
	panic("implement me")

}

func (b *BitcoinCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	panic("implement me")
}

func (b *BitcoinCoin) GetForwardersInBlock(blockID int64) ([]*ForwardInput, error) {
	panic("not needed for bitcoin")
}

func (b *BitcoinCoin) CombineSignatures(signs []string) (string, error) {
	sigsByte,err := json.Marshal(signs)
	args := []string {
		"-datadir=blockchain/bitcoin/data/",
		"combinerawtransaction",
		string(sigsByte),
	}

	cmd := exec.Command("bitcoin-cli", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		println("err", err.Error(), stderr.String())
	}

	return out.String(), err
}

func (b *BitcoinCoin) SendWithdrawal(trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {

	combined, err := b.CombineSignatures(w.Signatures)
	if err != nil {
		return "", errors.Wrap(err, "Could not combined sigs")
	}

	dataBytes, err := hex.DecodeString(strings.TrimSpace(combined))
	if err != nil {
		return "", errors.Wrap(err, "Could not decode combined sig")
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	err = tx.Deserialize(bytes.NewBuffer(dataBytes))

	if err != nil {
		return "", errors.Wrap(err, "Could not deserialize combined")
	}

	hash, err := b.Client.SendRawTransaction(tx, false)
	if hash == nil {
		return "", err
	}

	return hash.String(), err
}

// TODO: inspect all unspent for addresses that matches the pattern for our multisig
// gather enough input and create a refund
func (b *BitcoinCoin) EncodeRefund(w Withdrawal) (string, error) {
	fmt.Printf("Encode refund %v\n", w)

	sourceAddr, err := btcutil.DecodeAddress(w.SourceAddress, b.chaincfg)
	if err != nil {
		return "", err
	}

	addr, err := btcutil.DecodeAddress(w.DestinationAddress, b.chaincfg)
	if err != nil {
		return "", err
	}
	amount, err := btcutil.NewAmount(float64(w.Amount)/1e5)
	println(amount.ToBTC())
	if err != nil {
		return "", err
	}

	// get latest hash
	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return "", err
	}

	inputs := []btcjson.TransactionInput{}
	unspentFound := []btcjson.ListUnspentResult{}
	rawInput := []btcjson.RawTxInput{}

	for _, e := range unspent {
		if e.Address == w.SourceAddress {
			inputs = append(inputs, btcjson.TransactionInput{ Txid: e.TxID, Vout: e.Vout })
			unspentFound = append(unspentFound, e)
			rawInput = append(rawInput, btcjson.RawTxInput{
				Txid: e.TxID,
				Vout: e.Vout,
				RedeemScript: e.RedeemScript,
				ScriptPubKey: e.ScriptPubKey,
			})
			break
			//redeemBytes, err := hex.DecodeString(e.RedeemScript)
			//decoded, err := b.Client.DecodeScript(redeemBytes)
			//if err != nil {
			//
			//}
			//fmt.Printf("%v\n", decoded)
			//addr, err = btcutil.DecodeAddress(decoded.Addresses[0], b.chaincfg)
			//println("ADDR: ", addr.String(), addr.EncodeAddress(), hex.EncodeToString(addr.ScriptAddress()))
			//addrObj, err := b.Client.ValidateAddress(addr)
			//fmt.Printf("pub: %s %v\n", addrObj.PubKey, addrObj)

			// TODO: attempt to match the pubkey for 3/4 of all the keys

			//addrBytes, err := hex.DecodeString(addrObj.PubKey)
			//addr, err = btcutil.NewAddressPubKey(addrBytes, b.chaincfg)
			//println("ADDR: ", addr.String(), addr.EncodeAddress())
			//
			//addr, err = btcutil.DecodeAddress("2NENNHR9Y9fpKzjKYobbdbwap7xno7sbf2E", b.chaincfg)
			//println("ADDR: ", addr.String(), addr.EncodeAddress(), hex.EncodeToString(addr.ScriptAddress()))
			//wif, err := btcutil.NewWIF(nil,nil false)
			//wif.SerializePubKey()
			//
			////addrBytes, err = hex.DecodeString("76a9140025fee4b761c245cba21e9993fd7d86261977a188ac")
			//addr, err = btcutil.NewAddressPubKey(addr.ScriptAddress(), b.chaincfg)
			//println("ADDR: ", addr.String(), addr.EncodeAddress())
		}
	}


	if len(inputs) == 0 {
		return "", errors.New("No unspent input found")
	}

	fee := 0.00001
	remain, err := btcutil.NewAmount(unspentFound[0].Amount - amount.ToBTC() - fee)
	if err != nil {
		return "", errors.Wrap(err, "Unable to create new amount")
	}

	tx, err := b.Client.CreateRawTransaction(inputs, map[btcutil.Address]btcutil.Amount {
		addr : amount,
		sourceAddr: remain,
	}, nil)

	if err != nil {
		return "", errors.Wrap(err, "Create Raw tx failure")
	}

	var buf bytes.Buffer
	err = tx.Serialize(&buf)

	if err != nil {
		return "", err
	}

	res := common2.TransactionBitcoin {
		Tx: hex.EncodeToString(buf.Bytes()),
		RawInput: rawInput,
	}

	resJson, err := json.Marshal(res)
	return string(resJson), err
}

func (b *BitcoinCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	panic("implement me")
}


