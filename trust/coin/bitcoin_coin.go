package coin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"math"
	"regexp"
	"strings"
	"time"
)

type BitcoinCoin struct {
	Client             *rpcclient.Client
	rpcHost            string
	chaincfg           *chaincfg.Params
	signers            []btcutil.Address
	crosschainAddr     map[string]string
	maxFee             float64
	rpcUser            string
	rpcPassword        string
	grapheneSeedPrefix string
	BtcWithdrawMin     float64
	BtcWithdrawFee     float64
}

const BLOCKCHAIN_BTC = "BTC"

func (c *BitcoinCoin) Blockchain() string {
	return BLOCKCHAIN_BTC
}

func (b *BitcoinCoin) Attach() error {
	var err error
	b.Client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.rpcHost,
		Endpoint:     "http",
		User:         b.rpcUser,
		Pass:         b.rpcPassword,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	b.maxFee = 0.0002
	if err != nil {
		return errors.Wrap(err, "Could not attach the client for BTC")
	}

	//err = crypto.ValidateNetwork(b.Client, "Satoshi")
	return err
}

type FeeResult struct {
	FeeRate float64 `json:"feerate"`
	Blocks  int     `json:"blocks"`
}

func (b *BitcoinCoin) estimateFee(inputs, outputs int) (float64, float64, error) {
	totalBytes := float64(350.0 + (180.0 * inputs) + (34.0 * outputs) + 10.0)

	numBlocks, err := json.Marshal(int(5))
	if err != nil {
		return 0, 0, err
	}
	mode, err := json.Marshal("ECONOMICAL")
	if err != nil {
		return 0, 0, err
	}
	rawParams := []json.RawMessage{numBlocks, mode}
	res, err := b.Client.RawRequest("estimatesmartfee", rawParams)

	// decode result to string
	var result FeeResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, 0, err
	}

	if err != nil {
		return 0, 0, err
	}

	// testnet is set to zero? override with our minimum
	feeRateMin := math.Max(result.FeeRate, 0.00001)
	return result.FeeRate, feeRateMin * (totalBytes / 1000.0), nil
}

func (b *BitcoinCoin) CheckValidAmount(amount uint64) bool {
	if amount < uint64(b.BtcWithdrawMin*CONST_PRECISION) {
		return false
	}
	return true
}

func (b *BitcoinCoin) GetBlockTime(blockId int64) (time.Time, error) {
	var t time.Time
	blockHash, err := b.Client.GetBlockHash(blockId)
	if err != nil {
		return t, err
	}

	block, err := b.Client.GetBlockVerbose(blockHash)
	if err != nil {
		return t, err
	}

	return time.Unix(block.Time, 0), err
}

func (b *BitcoinCoin) GetTopBlockID() (int64, error) {
	blockId, err := b.Client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockId, err
}

func (b *BitcoinCoin) GenerateMultisig(accountId string) (string, error) {
	addr := []btcutil.Address{}
	addr = append(addr, b.signers...)
	btcAddressStr, err := crypto.GenerateGrapheneKeyWithSeed(accountId, b.grapheneSeedPrefix)
	if err != nil {
		return "", err
	}

	graphenePK, err := crypto.NewGraphenePublicKeyFromString(btcAddressStr)
	if err != nil {
		return "", err
	}

	btcAddress, err := crypto.GetBitcoinAddressFromGraphene(graphenePK)
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

	err = b.Client.ImportAddressRescan(res.P2sh, "", false)
	if err != nil {
		return "", errors.Wrap(err, "Unable to import address "+res.P2sh)
	}

	return res.P2sh, err
}

func (b *BitcoinCoin) GetTxID(trustAddress common.Address) (uint64, error) {
	//panic("implement me")
	return 0, nil
}

func (b *BitcoinCoin) GetFromAddress(txHash *chainhash.Hash) (string, error) {
	currentTx, err := b.Client.GetRawTransactionVerbose(txHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to getraw for currentTx")
	}

	vinLookup := map[string]bool{}
	vinAddresses := []string{}
	for _, vin := range currentTx.Vin {
		if vin.Txid == "" {
			continue
		}

		prevTranHash, err := chainhash.NewHashFromStr(vin.Txid)
		if err != nil {
			return "", errors.Wrap(err, "failed to build hash")
		}
		prevTran, err := b.Client.GetTransaction(prevTranHash)
		//will return error for an external deposit
		if err != nil {
			//returning "" to differentiate between external and internal deposits
			return "", nil
		}
		tx, err := hex.DecodeString(prevTran.Hex)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode hex string")
		}
		decodedTx, err := b.Client.DecodeRawTransaction(tx)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode raw transaction")
		}

		prevVout := decodedTx.Vout[vin.Vout]
		fromAddress := strings.Join(prevVout.ScriptPubKey.Addresses, ",")
		vinLookup[fromAddress] = true
		vinAddresses = append(vinAddresses, fromAddress)
	}

	fromAddr := strings.Join(vinAddresses, ",")
	fmt.Println("from address = ", fromAddr)
	return fromAddr, nil
}

func (b *BitcoinCoin) GetPendingTx(watchMap map[string]string) ([]*Deposit, error) {
	results, err := b.Client.ListUnspentMinMax(0, 0)
	if err != nil {
		return nil, err
	}

	events := []*Deposit{}

	for _, e := range results {
		toAddr := e.Address

		txHash, err := chainhash.NewHashFromStr(e.TxID)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get chainhash")
		}
		fromAddr, err := b.GetFromAddress(txHash)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get from address")
		}

		if fromAddr != "" && fromAddr == toAddr {
			fmt.Println("Skipping deposit in pending as it is the remaining amount")
			continue
		}

		//amount, err := btcutil.NewAmount(e.Amount)
		//if err != nil {
		//	return nil, errors.Wrap(err, "unable to create new amount")
		//}

		if quantaAddr, ok := watchMap[toAddr]; ok {
			amount, err := btcutil.NewAmount(e.Amount)
			if err != nil {
				return nil, errors.Wrap(err, "unable to create new amount")
			}
			events = append(events, &Deposit{
				QuantaAddr: quantaAddr,
				CoinName:   b.Blockchain(),
				Amount:     int64(amount),
				Tx:         txHash.String(),
			})
		}

	}
	//msg, _ := json.Marshal(events)
	//fmt.Printf("pending events = %v\n", string(msg))
	return events, nil
}

func (b *BitcoinCoin) GetBlockInfo(hash string) (string, int64, error) {
	var res *btcjson.GetBlockVerboseResult
	chainhash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return "", 0, err
	}

	res, err = b.Client.GetBlockVerbose(chainhash)
	if err != nil {
		return "", 0, err
	}
	return res.Hash, res.Confirmations, nil
}

func (b *BitcoinCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	blockHash, err := b.Client.GetBlockHash(blockID)
	if err != nil {
		return nil, err
	}

	block, err := b.Client.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

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

			fromAddr, err := b.GetFromAddress(&txHash)
			if err != nil {
				return nil, errors.Wrap(err, "failed to from address for currentTx")
			}

			for _, vout := range currentTx.Vout {
				toAddr := strings.Join(vout.ScriptPubKey.Addresses, ",")

				if fromAddr != "" && fromAddr == toAddr {
					fmt.Println("Skipping deposit as it is the remaining amount")
					//println("Ignoring tx when from and to the same ", toAddr)
					continue
				}

				amount, err := btcutil.NewAmount(vout.Value)
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
						BlockHash:  blockHash.String(),
					})
				}
			}
		}
	}
	//msg,_ := json.Marshal(events)
	//fmt.Printf("events = %v\n", string(msg))
	return events, nil
}

func (c *BitcoinCoin) FlushCoin(forwarder string, address string) error {
	panic("not implemented")
}

func (b *BitcoinCoin) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	panic("not needed for bitcoin")
}

func (b *BitcoinCoin) CombineSignatures(signs []string) (string, error) {
	marshalledParam, err := json.Marshal(signs)
	if err != nil {
		return "", err
	}
	rawMessage := json.RawMessage(marshalledParam)
	rawParams := []json.RawMessage{rawMessage}

	res, err := b.Client.RawRequest("combinerawtransaction", rawParams)
	if err != nil {
		return "", nil
	}

	// decode result to string
	var combinedtx string
	err = json.Unmarshal(res, &combinedtx)
	if err != nil {
		return "", err
	}

	return combinedtx, nil
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

func (b *BitcoinCoin) FillCrosschainAddress(crosschainAddr map[string]string) {
	b.crosschainAddr = crosschainAddr
}

// GetUnspentInputs retrieves a list of unspent addresses that meets or exceed the amount, and returns a list of unspent data
// amount is 8 precision / satoshi / int64
func (b *BitcoinCoin) GetUnspentInputs(destAddress btcutil.Address, amount btcutil.Amount) (btcutil.Amount, []btcjson.TransactionInput, []btcjson.ListUnspentResult, []btcjson.RawTxInput, error) {
	// get latest hash
	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return 0, nil, nil, nil, errors.Wrap(err, "No unspent input found")
	}

	inputs := []btcjson.TransactionInput{}
	unspentFound := []btcjson.ListUnspentResult{}
	rawInput := []btcjson.RawTxInput{}
	totalAmount, _ := btcutil.NewAmount(0)

	_, estimateFees, err := b.estimateFee(2, 2)
	if err != nil {
		return 0, nil, nil, nil, errors.Wrap(err, "Unable to estimate fee")
	}

	amountWithFee := amount.ToBTC() + estimateFees

	for _, e := range unspent {
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
			unspentAmount, err := btcutil.NewAmount(e.Amount)
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
func (b *BitcoinCoin) EncodeRefund(w Withdrawal) (string, error) {
	fmt.Printf("Encode refund %v\n", w)

	//sourceAddr, err := btcutil.DecodeAddress(w.SourceAddress, b.chaincfg)
	//if err != nil {
	//	return "", err
	//}

	destinationAddr, err := btcutil.DecodeAddress(w.DestinationAddress, b.chaincfg)
	if err != nil {
		return "", errors.Wrap(err, "decode address")
	}
	amountMinusFee := w.Amount - uint64(b.BtcWithdrawFee*CONST_PRECISION)
	amount, err := btcutil.NewAmount(float64(amountMinusFee) / CONST_PRECISION)
	println(amount.ToBTC())
	if err != nil {
		return "", errors.Wrap(err, "convert amount problem")
	}

	totalAmount, inputs, unspentFound, rawInput, err := b.GetUnspentInputs(destinationAddr, amount)

	if err != nil {
		return "", errors.Wrap(err, "Unable to get unspent")
	}

	if len(inputs) == 0 {
		return "", errors.New("No unspent input found")
	}

	feeRate, fees, err := b.estimateFee(len(inputs), 2)

	if err != nil {
		return "", errors.Wrap(err, "Unable to estimate fee")
	}
	if fees > b.maxFee {
		return "", errors.New("Fee is too high")
	}

	fmt.Printf("Fee calculated %f feeRate=%f remain=%f\n", fees, feeRate, totalAmount.ToBTC()-amount.ToBTC())

	remain, err := btcutil.NewAmount(totalAmount.ToBTC() - amount.ToBTC() - fees)
	if err != nil {
		return "", errors.Wrap(err, "Unable to create new amount")
	}

	sourceAddrRefund, err := btcutil.DecodeAddress(unspentFound[0].Address, b.chaincfg)
	if err != nil {
		return "", err
	}

	fmt.Printf("total=%f amount=%f, remain=%f, fees=%f acual=%f \n", totalAmount.ToBTC(), amount.ToBTC(), remain.ToBTC(), fees, totalAmount.ToBTC()-amount.ToBTC()-remain.ToBTC())

	tx, err := b.Client.CreateRawTransaction(inputs, map[btcutil.Address]btcutil.Amount{
		destinationAddr:  amount,
		sourceAddrRefund: remain,
	}, nil)

	if err != nil {
		return "", errors.Wrap(err, "Create Raw tx failure")
	}

	var buf bytes.Buffer
	err = tx.Serialize(&buf)

	if err != nil {
		return "", errors.Wrap(err, "Serialize failure")
	}

	res := common2.TransactionBitcoin{
		Tx:       hex.EncodeToString(buf.Bytes()),
		RawInput: rawInput,
	}

	resJson, err := json.Marshal(res)
	data, err := json.Marshal(&EncodedMsg{string(resJson), w.Tx, w.QuantaBlockID, w.CoinName, w.DestinationAddress})
	return string(data), err
}

// DecodeRefund decodes the encodemsg object. If there are multiple withdrawals, then it means the refund did not original from
// our encode refund, and we fail immediately.
// TODO: implement test by composing tx from sendmany from bitcoin CLI
func (b *BitcoinCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
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
		amount, err := btcutil.NewAmount(vout.Value)
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

func (b *BitcoinCoin) CheckValidAddress(address string) bool {
	if len(address) > 35 || len(address) < 26 {
		return false
	}
	_, err := btcutil.DecodeAddress(address, b.chaincfg)
	if err != nil {
		return false
	}
	var validAddress = regexp.MustCompile(`^[1-9a-zA-NP-Z]`)
	return validAddress.MatchString(address)
}
