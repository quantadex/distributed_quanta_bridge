package coin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gcash/bchd/btcjson"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/chaincfg/chainhash"
	"github.com/gcash/bchd/rpcclient"
	"github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	"github.com/pkg/errors"
	common2 "github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"time"

	//"os/exec"
	"strings"
)

const BLOCKCHAIN_BCH = "BCH"

type BCH struct {
	Client             *rpcclient.Client
	rpcHost            string
	chaincfg           *chaincfg.Params
	signers            []bchutil.Address
	crosschainAddr     map[string]string
	fee                float64
	rpcUser            string
	rpcPassword        string
	grapheneSeedPrefix string
	BchWithdrawMin     float64
	BchWithdrawFee     float64
	BlackList          map[string]bool
	issuerAddr         string
}

func (b *BCH) Blockchain() string {
	return BLOCKCHAIN_BCH
}

func (b *BCH) Attach() error {
	var err error
	b.Client, err = rpcclient.New(&rpcclient.ConnConfig{Host: b.rpcHost,
		Endpoint:     "http",
		User:         b.rpcUser,
		Pass:         b.rpcPassword,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	b.fee = 0.00001

	if err != nil {
		return errors.Wrap(err, "Could not attach the client for BCH")
	}

	//err = crypto.ValidateNetwork(b.Client, "Bitcoin ABC")
	return err
}

func (b *BCH) SetIssuerAddress(address string) {
	b.issuerAddr = address
}

func (b *BCH) CheckValidAmount(amount uint64) bool {
	if amount < uint64(b.BchWithdrawMin*CONST_PRECISION) {
		return false
	}
	return true
}

func (b *BCH) GetBlockTime(blockId int64) (time.Time, error) {
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

func (b *BCH) GetTopBlockID() (int64, error) {
	blockId, err := b.Client.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockId, err
}

func (b *BCH) GetBlockInfo(hash string) (string, int64, error) {
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

func (b *BCH) GenerateMultisig(accountId string) (string, error) {
	fmt.Println("BCH::GenerateMultisig ", accountId, b.signers, b.grapheneSeedPrefix)

	addr := []bchutil.Address{}
	addr = append(addr, b.signers...)
	btcAddressStr, err := crypto.GenerateGrapheneKeyWithSeed(accountId, b.grapheneSeedPrefix)
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

	addrx, err := b.Client.AddMultisigAddress(len(addr)-1, addr, "", b.chaincfg)

	if err != nil {
		return "", err
	}

	res := b.chaincfg.CashAddressPrefix + ":" + addrx.String()
	err = b.Client.ImportAddressRescan(res, "", false)
	if err != nil {
		return "", errors.Wrap(err, "Unable to import address "+res)
	}

	return res, err
}

func (b *BCH) GetTxID(trustAddress common.Address) (uint64, error) {
	return 0, nil
}

func (b *BCH) GetFromAddress(txHash *chainhash.Hash) ([]string, error) {
	currentTx, err := b.Client.GetRawTransactionVerbose(txHash)
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
		prevTran, err := b.Client.GetTransaction(prevTranHash)
		//will return error for an external deposit
		if err != nil {
			//returning "" to differentiate between external and internal deposits
			return nil, nil
		}
		tx, err := hex.DecodeString(prevTran.Hex)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode hex string")
		}
		decodedTx, err := b.Client.DecodeRawTransaction(tx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode raw transaction")
		}

		prevVout := decodedTx.Vout[vin.Vout]
		fromAddress := strings.Join(prevVout.ScriptPubKey.Addresses, ",")
		vinLookup[fromAddress] = true
		vinAddresses = append(vinAddresses, fromAddress)
	}

	//fromAddr := strings.Join(vinAddresses, ",")
	return vinAddresses, nil
}

func (b *BCH) GetPendingTx(watchMap map[string]string) ([]*Deposit, error) {
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

		isCrosschain := false
		for _, addr := range fromAddr {
			_, isCrosschain = watchMap[addr]
			if isCrosschain {
				break
			}
		}

		if fromAddr != nil && isCrosschain {
			fmt.Println("Skipping deposit in pending as it is the remaining amount")
			continue
		}

		if quantaAddr, ok := watchMap[toAddr]; ok {
			amount, err := bchutil.NewAmount(e.Amount)
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
	return events, nil
}

func (b *BCH) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
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

				isCrosschain := false
				for _, addr := range fromAddr {
					_, isCrosschain = trustAddress[addr]
					if isCrosschain {
						break
					}
				}

				if fromAddr != nil && isCrosschain {
					fmt.Println("Skipping deposit as it is the remaining amount")
					//println("Ignoring tx when from and to the same ", toAddr)
					continue
				}

				amount, err := bchutil.NewAmount(vout.Value)
				if err != nil {
					return nil, errors.Wrap(err, "unable to create new amount")
				}

				if quantaAddr, ok := trustAddress[toAddr]; ok {
					events = append(events, &Deposit{
						SenderAddr: strings.Join(fromAddr, ","),
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
	return events, nil
}

func (b *BCH) FlushCoin(forwarder string, address string) error {
	panic("implement me")
}

func (b *BCH) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	panic("implement me")
}

func (b *BCH) CombineSignatures(signs []string) (string, error) {

	// encode the input as rawmessage
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

func (b *BCH) GetUnspentInputs(destAddress bchutil.Address, amount bchutil.Amount) (bchutil.Amount, []btcjson.TransactionInput, []btcjson.ListUnspentResult, []btcjson.RawTxInput, error) {
	// get latest hash
	unspent, err := b.Client.ListUnspent()
	if err != nil {
		return 0, nil, nil, nil, errors.Wrap(err, "No unspent input found")
	}

	inputs := []btcjson.TransactionInput{}
	unspentFound := []btcjson.ListUnspentResult{}
	rawInput := []btcjson.RawTxInput{}
	totalAmount, _ := bchutil.NewAmount(0)

	amountWithFee := amount.ToBCH() + (b.fee * 50)

	for _, e := range unspent {
		if _, ok := b.crosschainAddr[e.Address]; ok {

			inputs = append(inputs, btcjson.TransactionInput{Txid: e.TxID, Vout: e.Vout})
			unspentFound = append(unspentFound, e)
			rawInput = append(rawInput, btcjson.RawTxInput{
				Txid:         e.TxID,
				Vout:         e.Vout,
				RedeemScript: e.RedeemScript,
				ScriptPubKey: e.ScriptPubKey,
				Amount:       e.Amount,
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
	if _, ok := b.BlackList[w.SourceAddress]; ok {
		return "", errors.New("BlackListed user: " + w.SourceAddress)
	}
	destinationAddr, err := bchutil.DecodeAddress(w.DestinationAddress, b.chaincfg)
	if err != nil {
		return "", err
	}
	amountMinusFee := w.Amount - uint64(b.BchWithdrawFee*CONST_PRECISION)
	amount, err := bchutil.NewAmount(float64(amountMinusFee) / CONST_PRECISION)
	if err != nil {
		return "", err
	}

	totalAmount, inputs, _, rawInput, err := b.GetUnspentInputs(destinationAddr, amount)

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

	if b.issuerAddr == "" {
		return "", errors.New("Issuer address not set for bch")
	}

	sourceAddrRefund, err := bchutil.DecodeAddress(b.issuerAddr, b.chaincfg)
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
	data, err := json.Marshal(&EncodedMsg{string(resJson), w.Tx, w.QuantaBlockID, w.CoinName, w.DestinationAddress})
	return string(data), err
}

func (b *BCH) DecodeRefund(encoded string) (*Withdrawal, error) {
	fmt.Println("encoded string = ", encoded)
	var decoded EncodedMsg
	err := json.Unmarshal([]byte(encoded), &decoded)
	if err != nil {
		return nil, err
	}

	var res common2.TransactionBCH
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
	fmt.Println("Transaction = ", decodedTx.Txid)

	vinLookup := map[string]bool{}
	vinAddresses := []string{}
	for _, vin := range decodedTx.Vin {
		if vin.Txid == "" {
			continue
		}
		fmt.Println("vin txid = ", vin.Txid)

		prevTranHash, err := chainhash.NewHashFromStr(vin.Txid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build hash")
		}
		fmt.Println("PrevTransaction hash = ", prevTranHash.String())

		prevTran, err := b.Client.GetRawTransactionVerbose(prevTranHash)
		fmt.Println("Previous transaction = ", prevTran)
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

		_, isCrosschain := b.crosschainAddr[destAddr]
		if vinLookup[destAddr] || fromAddr == "" || isCrosschain {
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
	_, err := bchutil.DecodeAddress(address, b.chaincfg)
	if err != nil {
		return false
	}
	return true
}
