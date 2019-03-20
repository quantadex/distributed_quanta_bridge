package coin

import (
	"github.com/bchsuite/bchd/bchjson"
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchd/chaincfg/chainhash"
	"github.com/bchsuite/bchd/rpcclient"
	"github.com/bchsuite/bchd/wire"
	"github.com/bchsuite/bchutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"

	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
)

type BCH struct {
	Client         *rpcclient.Client
	rpcHost        string
	chaincfg       *chaincfg.Params
	signers        []bchutil.Address
	crosschainAddr map[string]string
	fee            float64
}

func (b *BCH) Blockchain() string {
	return BLOCKCHAIN_BTC
}

func (b *BCH) Attach() error {
	var err error
	b.Client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.rpcHost,
		Endpoint:     "http",
		User:         "user",
		Pass:         "123",
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	b.fee = 0.00001
	return err
}

func (b *BCH) GetTopBlockID() (int64, error) {
	blockId, err := b.Client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockId, err
}

func (b *BCH) GenerateMultisig(accountId string) (string, error) {
	addr := []bchutil.Address{}
	addr = append(addr, b.signers...)
	btcAddressStr, err := crypto.GenerateGrapheneKeyWithSeed(accountId)
	if err != nil {
		return "", err
	}

	graphenePK, err := crypto.NewGraphenePublicKeyFromString(btcAddressStr)
	if err != nil {
		return "", err
	}

	btcAddress, err := crypto.GetBCHAddressFromGraphene(graphenePK)
	if err != nil {
		return "", err
	}

	addr = append(addr, btcAddress)

	addrx, err := b.Client.AddMultisigAddress(len(addr)-1, addr, "")
	//fmt.Println("result ", addrx)

	if err != nil {
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

	err = b.Client.ImportAddressRescan(res.P2sh, false)
	if err != nil {
		return "", errors.Wrap(err, "Unable to import address "+res.P2sh)
	}

	return res.P2sh, err
}

func (b *BCH) GetTxID(trustAddress common.Address) (uint64, error) {
	return 0, nil
}

func (b *BCH) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	blockHash, err := b.Client.GetBlockHash(blockID)
	if err != nil {
		return nil, err
	}
	block, err := b.Client.GetBlock(blockHash)

	events := []*Deposit{}

	for _, tx := range block.Transactions {
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
				return events, errors.Wrap(err, "failed to build hash")
			}
			prevTran, err := b.Client.GetRawTransactionVerbose(prevTranHash)
			if err != nil {
				return events, errors.Wrap(err, "failed to getraw for vin")
			}

			prevVout := prevTran.Vout[vin.Vout]
			fromAddress := strings.Join(prevVout.ScriptPubKey.Addresses, ",")
			vinLookup[fromAddress] = true
			vinAddresses = append(vinAddresses, fromAddress)
		}

		fromAddr := strings.Join(vinAddresses, ",")

		for _, vout := range currentTx.Vout {
			toAddr := strings.Join(vout.ScriptPubKey.Addresses, ",")

			if fromAddr == toAddr || fromAddr == "" {
				//println("Ignoring tx when from and to the same ", toAddr)
				continue
			}

			amount, err := bchutil.NewAmount(vout.Value)
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
	//msg,_ := json.Marshal(events)
	//fmt.Printf("events = %v\n", string(msg))
	return events, nil
}

func (b *BCH) FlushCoin(forwarder string, address string) error {
	panic("implement me")
}

func (b *BCH) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	panic("implement me")
}

func (b *BCH) CombineSignatures(signs []string) (string, error) {
	sigsByte, err := json.Marshal(signs)
	args := []string{
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

func (b *BCH) SendWithdrawal(trustAddress common.Address,
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

func (b *BCH) FillCrosschainAddress(crosschainAddr map[string]string) {
	b.crosschainAddr = crosschainAddr
}

func (b *BCH) GetUnspentInputs(destAddress bchutil.Address, amount bchutil.Amount) (bchutil.Amount, []bchjson.TransactionInput, []bchjson.ListUnspentResult, []bchjson.RawTxInput, error) {
	// get latest hash
	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return 0, nil, nil, nil, errors.Wrap(err, "No unspent input found")
	}

	inputs := []bchjson.TransactionInput{}
	unspentFound := []bchjson.ListUnspentResult{}
	rawInput := []bchjson.RawTxInput{}
	totalAmount, _ := bchutil.NewAmount(0)

	amountWithFee := amount.ToBCH() + (b.fee * 50)

	for _, e := range unspent {
		if _, ok := b.crosschainAddr[e.Address]; ok {
			//unspentAddr, err := btcutil.DecodeAddress(e.Address, b.chaincfg)
			//if unspentAddr.String() == destAddress.String() {
			//	return nil, nil, nil, errors.New("We don't expect destination address to be same as ")
			//}

			inputs = append(inputs, bchjson.TransactionInput{Txid: e.TxID, Vout: e.Vout})
			unspentFound = append(unspentFound, e)
			rawInput = append(rawInput, bchjson.RawTxInput{
				Txid:         e.TxID,
				Vout:         e.Vout,
				RedeemScript: e.RedeemScript,
				ScriptPubKey: e.ScriptPubKey,
			})
			unspentAmount, err := bchutil.NewAmount(e.Amount)
			if err != nil {
				return 0, nil, nil, nil, err
			}
			totalAmount += unspentAmount

			// we have enough coins, let's get out
			if totalAmount.ToBCH() >= amountWithFee {
				break
			}
		}
	}
	return totalAmount, inputs, unspentFound, rawInput, nil
}

func (b *BCH) EncodeRefund(w Withdrawal) (string, error) {
	destinationAddr, err := bchutil.DecodeAddress(w.DestinationAddress, b.chaincfg)
	if err != nil {
		return "", err
	}
	amount, err := bchutil.NewAmount(float64(w.Amount) / 1e5)
	println(amount.ToBCH())
	if err != nil {
		return "", err
	}

	totalAmount, inputs, unspentFound, rawInput, err := b.GetUnspentInputs(destinationAddr, amount)

	if err != nil {
		return "", err
	}

	if len(inputs) == 0 {
		return "", errors.New("No unspent input found")
	}

	fee := b.fee * float64(len(inputs))

	remain, err := bchutil.NewAmount(totalAmount.ToBCH() - amount.ToBCH() - fee)
	if err != nil {
		return "", errors.Wrap(err, "Unable to create new amount")
	}

	sourceAddrRefund, err := bchutil.DecodeAddress(unspentFound[0].Address, b.chaincfg)
	if err != nil {
		return "", err
	}

	tx, err := b.Client.CreateRawTransaction(inputs, map[bchutil.Address]bchutil.Amount{
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

	res := common2.TransactionBCH{
		Tx:       hex.EncodeToString(buf.Bytes()),
		RawInput: rawInput,
	}

	resJson, err := json.Marshal(res)
	data, err := json.Marshal(&EncodedMsg{string(resJson), w.Tx, w.QuantaBlockID, w.CoinName})
	return string(data), err
}

func (b *BCH) DecodeRefund(encoded string) (*Withdrawal, error) {
	var decoded EncodedMsg
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		return nil, err
	}

	var res common2.TransactionBitcoin
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
		amount, err := bchutil.NewAmount(vout.Value)
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

func (b *BCH) CheckValidAddress(address string) bool {
	if len(address) > 35 || len(address) < 26 {
		return false
	}
	_, err := bchutil.DecodeAddress(address, b.chaincfg)
	if err != nil {
		return false
	}
	var validAddress = regexp.MustCompile(`^[1-9a-zA-NP-Z]`)
	return validAddress.MatchString(address)
}
