package coin

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ltcsuite/ltcd/btcjson"
	chaincfg2 "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
	"github.com/pkg/errors"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"

	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

const BLOCKCHAIN_LTC = "LTC"

type LiteCoin struct {
	Client         *rpcclient.Client
	rpcHost        string
	chaincfg       *chaincfg2.Params
	signers        []ltcutil.Address
	crosschainAddr map[string]string
	fee            float64
}

func (b *LiteCoin) Blockchain() string {
	return BLOCKCHAIN_LTC
}

func (b *LiteCoin) Attach() error {
	var err error
	b.Client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.rpcHost,
		Endpoint:     "http",
		User:         "user",
		Pass:         "123",
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "Could not attach the client for LTC")
	}

	b.fee = 0.00001

	err = crypto.ValidateNetwork(b.Client, "Litecoin")

	return err
}

func (b *LiteCoin) GetTopBlockID() (int64, error) {
	blockId, err := b.Client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockId, err
}

func (b *LiteCoin) GetPendingTx(map[string]string) ([]*Deposit, error) {
	return nil, nil
}

func (b *LiteCoin) GenerateMultisig(accountId string) (string, error) {
	addr := []ltcutil.Address{}
	addr = append(addr, b.signers...)
	btcAddressStr, err := crypto.GenerateGrapheneKeyWithSeed(accountId)
	if err != nil {
		return "", err
	}

	graphenePK, err := crypto.NewGraphenePublicKeyFromString(btcAddressStr)
	if err != nil {
		return "", err
	}

	btcAddress, err := crypto.GetLitecoinAddressFromGraphene(graphenePK)
	if err != nil {
		return "", err
	}

	addr = append(addr, btcAddress)

	addrx, err := b.Client.AddMultisigAddress(len(addr)-1, addr, "")

	if err != nil {
		fmt.Println("error = ", err)
		return "", err
	}

	scriptBytes, err := hex.DecodeString(addrx.RedeemScript)
	if err != nil {
		return "", err
	}

	res, err := b.Client.DecodeScript(scriptBytes)

	if err != nil {
		return "", err
	}

	err = b.Client.ImportAddressRescan(res.P2sh, "", false)
	if err != nil {
		return "", errors.Wrap(err, "Unable to import address "+res.P2sh)
	}

	return res.P2sh, err
}

func (b *LiteCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	//panic("implement me")
	return 0, nil
}

func (b *LiteCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	blockHash, err := b.Client.GetBlockHash(blockID)
	if err != nil {
		return nil, err
	}

	block, err := b.Client.GetBlock(blockHash)

	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unspent")
	}

	txidMap := make(map[string]string)
	for _, e := range unspent {
		txidMap[e.TxID] = e.Address
	}

	events := []*Deposit{}

	for _, tx := range block.Transactions {
		if _, ok := txidMap[tx.TxHash().String()]; ok {
			txHash := tx.TxHash()
			currentTx, err := b.Client.GetRawTransactionVerbose(&txHash)
			if err != nil {
				return nil, errors.Wrap(err, "failed to getraw for currentTx")
			}

			vinLookup := map[string]bool{}
			vinAddresses := []string{}
			for _, vin := range currentTx.Vin {
				if vin.Txid == "" {
					continue
				}

				prevTranHash, err := chainhash.NewHashFromStr(vin.Txid)
				if err != nil {
					return nil, errors.Wrap(err, "failed to build hash")
				}
				prevTran, err := b.Client.GetRawTransactionVerbose(prevTranHash)
				if err != nil {
					return nil, errors.Wrap(err, "failed to getraw for vin")
				}

				prevVout := prevTran.Vout[vin.Vout]
				fromAddress := strings.Join(prevVout.ScriptPubKey.Addresses, ",")
				vinLookup[fromAddress] = true
				vinAddresses = append(vinAddresses, fromAddress)
			}

			fromAddr := strings.Join(vinAddresses, ",")

			if err != nil {
				return nil, errors.Wrap(err, "failed to from address for currentTx")
			}

			for _, vout := range currentTx.Vout {
				toAddr := strings.Join(vout.ScriptPubKey.Addresses, ",")

				if fromAddr == toAddr || fromAddr == "" {
					//println("Ignoring tx when from and to the same ", toAddr)
					continue
				}

				amount, err := ltcutil.NewAmount(vout.Value)
				if err != nil {
					return nil, errors.Wrap(err, "unable to create new amount")
				}

				if quantaAddr, ok := trustAddress[toAddr]; ok {
					events = append(events, &Deposit{
						SenderAddr: fromAddr,
						QuantaAddr: quantaAddr,
						CoinName:   b.Blockchain(),
						Amount:     int64(amount),
						BlockID:    blockID,
						Tx:         txHash.String(),
					})
				}
			}
		}
	}
	msg, _ := json.Marshal(events)
	fmt.Printf("events = %v\n", string(msg))
	return events, nil
}

func (c *LiteCoin) FlushCoin(forwarder string, address string) error {
	panic("not implemented")
}

func (b *LiteCoin) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	panic("not needed for bitcoin")
}

func (b *LiteCoin) CombineSignatures(signs []string) (string, error) {
	sigsByte, err := json.Marshal(signs)
	args := []string{
		"combinerawtransaction",
		string(sigsByte),
	}

	cmd := exec.Command("litecoin-cli", args...)
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

func (b *LiteCoin) SendWithdrawal(trustAddress common.Address,
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

func (b *LiteCoin) FillCrosschainAddress(crosschainAddr map[string]string) {
	b.crosschainAddr = crosschainAddr
}

// GetUnspentInputs retrieves a list of unspent addresses that meets or exceed the amount, and returns a list of unspent data
// amount is 8 precision / satoshi / int64
func (b *LiteCoin) GetUnspentInputs(amount ltcutil.Amount) (ltcutil.Amount, []btcjson.TransactionInput, []btcjson.ListUnspentResult, []btcjson.RawTxInput, error) {
	// get latest hash
	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return 0, nil, nil, nil, errors.Wrap(err, "No unspent input found")
	}

	inputs := []btcjson.TransactionInput{}
	unspentFound := []btcjson.ListUnspentResult{}
	rawInput := []btcjson.RawTxInput{}
	totalAmount, _ := ltcutil.NewAmount(0)

	amountWithFee := amount.ToBTC() + (b.fee * 50)

	for _, e := range unspent {
		if e.Address == "QMmJjkHxUhLxRd73jk75qbFDuh72D22jHT" || e.Address == "QU5bHWP9os3CwMAbTVkJQQ4K8df43Vb4ru" {
			fmt.Println("address = ", e.Address)
			fmt.Println("map = ", b.crosschainAddr)
		}
		if _, ok := b.crosschainAddr[e.Address]; ok {
			//unspentAddr, err := btcutil.DecodeAddress(e.Address, b.chaincfg)
			//if unspentAddr.String() == destAddress.String() {
			//	return nil, nil, nil, errors.New("We don't expect destination address to be same as ")
			//}

			inputs = append(inputs, btcjson.TransactionInput{Txid: e.TxID, Vout: e.Vout})
			unspentFound = append(unspentFound, e)
			rawInput = append(rawInput, btcjson.RawTxInput{
				Txid:         e.TxID,
				Vout:         e.Vout,
				RedeemScript: e.RedeemScript,
				ScriptPubKey: e.ScriptPubKey,
			})
			unspentAmount, err := ltcutil.NewAmount(e.Amount)
			if err != nil {
				return 0, nil, nil, nil, err
			}
			totalAmount += unspentAmount

			// we have enough coins, let's get out
			if totalAmount.ToBTC() >= amountWithFee {
				break
			}
		}
	}
	return totalAmount, inputs, unspentFound, rawInput, nil
}

// TODO: inspect all unspent for addresses that matches the pattern for our multisig
// gather enough input and create a refund
// encoding withdrawal with precision of 5.
// must convert to our system precision
func (b *LiteCoin) EncodeRefund(w Withdrawal) (string, error) {
	fmt.Printf("Encode refund %v\n", w)

	destinationAddr, err := ltcutil.DecodeAddress(w.DestinationAddress, b.chaincfg)
	if err != nil {
		return "", err
	}
	amount, err := ltcutil.NewAmount(float64(w.Amount) / 1e5)
	println(destinationAddr.String(), amount.ToBTC())
	if err != nil {
		return "", err
	}

	totalAmount, inputs, unspentFound, rawInput, err := b.GetUnspentInputs(amount)

	if err != nil {
		return "", err
	}

	if len(inputs) == 0 {
		return "", errors.New("No unspent input found")
	}

	fee := b.fee * float64(len(inputs))

	remain, err := ltcutil.NewAmount(totalAmount.ToBTC() - amount.ToBTC() - fee)
	if err != nil {
		return "", errors.Wrap(err, "Unable to create new amount")
	}

	sourceAddrRefund, err := ltcutil.DecodeAddress(unspentFound[0].Address, b.chaincfg)
	if err != nil {
		return "", err
	}

	tx, err := b.Client.CreateRawTransaction(inputs, map[ltcutil.Address]ltcutil.Amount{
		destinationAddr:  amount,
		sourceAddrRefund: remain,
	}, nil)

	if err != nil {
		return "", errors.Wrap(err, "Create Raw tx failure")
	}

	var buf bytes.Buffer
	err = tx.Serialize(&buf)

	if err != nil {
		return "", err
	}

	res := common2.TransactionLitecoin{
		Tx:       hex.EncodeToString(buf.Bytes()),
		RawInput: rawInput,
	}

	resJson, err := json.Marshal(res)
	data, err := json.Marshal(&EncodedMsg{string(resJson), w.Tx, w.QuantaBlockID, w.CoinName})
	return string(data), err
}

// DecodeRefund decodes the encodemsg object. If there are multiple withdrawals, then it means the refund did not original from
// our encode refund, and we fail immediately.
// TODO: implement test by composing tx from sendmany from bitcoin CLI
func (b *LiteCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	var decoded EncodedMsg
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		return nil, err
	}

	var res common2.TransactionLitecoin
	err = json.Unmarshal([]byte(decoded.Message), &res)
	if err != nil {
		return nil, err
	}

	tx, err := hex.DecodeString(res.Tx)
	if err != nil {
		return nil, err
	}

	var msgTx wire.MsgTx
	err = msgTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		return nil, err
	}

	decodedTx, err := b.Client.DecodeRawTransaction(tx)
	if err != nil {
		return nil, err
	}

	vinLookup := map[string]bool{}
	vinAddresses := []string{}
	for _, vin := range decodedTx.Vin {
		if vin.Txid == "" {
			continue
		}

		prevTranHash, err := chainhash.NewHashFromStr(vin.Txid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build hash")
		}
		prevTran, err := b.Client.GetRawTransactionVerbose(prevTranHash)
		if err != nil {
			return nil, errors.Wrap(err, "failed to getraw for vin")
		}

		prevVout := prevTran.Vout[vin.Vout]
		fromAddress := strings.Join(prevVout.ScriptPubKey.Addresses, ",")
		vinLookup[fromAddress] = true
		vinAddresses = append(vinAddresses, fromAddress)
	}

	fromAddr := strings.Join(vinAddresses, ",")
	allWithdrawals := []*Withdrawal{}
	for _, vout := range decodedTx.Vout {
		destAddr := strings.Join(vout.ScriptPubKey.Addresses, ",")

		if vinLookup[destAddr] || fromAddr == "" {
			//println("Ignoring tx when from and to the same ", toAddr)
			continue
		}
		amount, err := ltcutil.NewAmount(vout.Value)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create new amount")
		}
		w := &Withdrawal{
			Amount:             uint64(amount),
			SourceAddress:      fromAddr,
			DestinationAddress: destAddr,
			Tx:                 decoded.Tx,
			QuantaBlockID:      decoded.BlockNumber,
		}
		allWithdrawals = append(allWithdrawals, w)
	}

	if len(allWithdrawals) != 1 {
		return nil, errors.New("Expect to have only 1 withdrawal")
	}
	return allWithdrawals[0], nil
}

func (b *LiteCoin) CheckValidAddress(address string) bool {
	_, err := ltcutil.DecodeAddress(address, b.chaincfg)
	if err != nil {
		return false
	}
	return true
}