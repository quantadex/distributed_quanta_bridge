package coin

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	crypto2 "github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"math/big"
	"strings"
	"time"
)

const abiCode = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"tokenOwner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Approval","type":"event"}]`

func (l *Listener) Start() error {
	l.log = log.DefaultLogger.WithField("service", "EthereumListener")
	l.log.Logger.Info("Ethereum listner started")

	//blockNumber, err := l.Storage.GetEthereumBlockToProcess()
	//if err != nil {
	//	err = errors.Wrap(err, "Error getting ethereum block to process from DB")
	//	l.log.Error(err.Error())
	//	return err
	//}

	// Check if connected to correct network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	id, err := l.Client.NetworkID(ctx)
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum network ID")
		l.log.Error(err.Error())
		return err
	}

	l.log.Logger.Infof("Created RPC listener with id %s", l.NetworkID)

	if id.String() != l.NetworkID {
		return errors.Errorf("Invalid network ID (have=%s, want=%s)", id.String(), l.NetworkID)
	}

	//go l.processBlocks(blockNumber)
	return nil
}

func (l *Listener) processBlocks(blockNumber int64) {
	if blockNumber == 0 {
		l.log.Info("Starting from the latest block")
	} else {
		l.log.Infof("Starting from block %d", blockNumber)
	}

	// Time when last new block has been seen
	lastBlockSeen := time.Now()
	noBlockWarningLogged := false

	for {
		block, err := l.GetBlock(blockNumber)
		if err != nil {
			l.log.WithFields(log.F{"err": err, "blockNumber": blockNumber}).Error("Error getting block")
			time.Sleep(1 * time.Second)
			continue
		}

		// Block doesn't exist yet
		if block == nil {
			if time.Since(lastBlockSeen) > 3*time.Minute && !noBlockWarningLogged {
				l.log.Warn("No new block in more than 3 minutes")
				noBlockWarningLogged = true
			}

			time.Sleep(1 * time.Second)
			continue
		}

		// Reset counter when new block appears
		lastBlockSeen = time.Now()
		noBlockWarningLogged = false

		//if block.NumberU64() == 0 {
		//	l.log.Error("Etheruem node is not synced yet. Unable to process blocks")
		//	time.Sleep(30 * time.Second)
		//	continue
		//}

		err = l.processBlock(block)
		if err != nil {
			l.log.WithFields(log.F{"err": err, "blockNumber": block.NumberU64()}).Error("Error processing block")
			time.Sleep(1 * time.Second)
			continue
		}

		// Persist block number
		err = l.Storage.SaveLastProcessedEthereumBlock(blockNumber)
		if err != nil {
			l.log.WithField("err", err).Error("Error saving last processed block")
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		blockNumber = block.Number().Int64() + 1
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Listener) GetBlock(blockNumber int64) (*types.Block, error) {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, err := l.Client.BlockByNumber(ctx, big.NewInt(blockNumber))
	if err != nil {
		if err.Error() == "not found" {
			return nil, err
		}
		err = errors.Wrap(err, "Error getting block from geth")
		l.log.WithField("block", blockNumber).Error(err)
		return nil, err
	}

	return block, nil
}

func (l *Listener) processBlock(block *types.Block) error {
	transactions := block.Transactions()
	blockTime := time.Unix(block.Time().Int64(), 0)

	localLog := l.log.WithFields(log.F{
		"blockNumber":  block.NumberU64(),
		"blockTime":    blockTime,
		"transactions": len(transactions),
	})
	localLog.Info("Processing block")

	for _, transaction := range transactions {
		to := transaction.To()
		if to == nil {
			// Contract creation
			continue
		}

		tx := Transaction{
			Hash:     transaction.Hash().Hex(),
			ValueWei: transaction.Value(),
			To:       to.Hex(),
		}
		err := l.TransactionHandler(tx)
		if err != nil {
			return errors.Wrap(err, "Error processing transaction")
		}
	}

	localLog.Info("Processed block")

	return nil
}

func (l *Listener) GetTopBlockNumber() (int64, error) {
	header, err := l.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}

func (l *Listener) GetNativeDeposits(blockNumber int64, toAddress map[string]string) ([]*Deposit, error) {
	blocks, err := l.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	if blocks == nil {
		return nil, errors.Wrap(err, "Block not found ")
	}

	events := []*Deposit{}
	for _, tx := range blocks.Transactions() {
		if tx.To() == nil {
			continue
		}

		if quantaAddr, ok := toAddress[strings.ToLower(tx.To().Hex())]; ok {
			receipt, err := l.Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				println("Cannot get Receipt ", err.Error())
			} else if receipt.Status == types.ReceiptStatusFailed {
				continue
			}
			if tx.Value().Cmp(big.NewInt(0)) != 0 {
				events = append(events, &Deposit{
					QuantaAddr: quantaAddr,
					CoinName:   "ETH",
					SenderAddr: tx.To().Hex(),
					Amount:     WeiToGraphene(*tx.Value()),
					BlockID:    blockNumber,
					Tx:         tx.Hash().Hex(),
				})
			}
		}
	}

	return events, nil
}

func (l *Listener) FilterTransferEvent(blockNumber int64, toAddress map[string]string) ([]*Deposit, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(blockNumber),
		ToBlock:   big.NewInt(blockNumber),
		//Addresses: []common.Address{
		//	common.HexToAddress(contractAddress),
		//},
	}

	logsEvent, err := l.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(abiCode)))
	if err != nil {
		return nil, err
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	//fmt.Printf("Number of log events %d transferHash=%s\n", len(logsEvent), logTransferSigHash.Hex())
	events := []*Deposit{}
	for _, vLog := range logsEvent {
		//fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		//fmt.Printf("Log Index: %d %s\n", vLog.Index, vLog.Topics[0].Hex())
		//fmt.Println(vLog.TxHash.Hex())

		if len(vLog.Topics) < 1 {
			continue
		}

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():

			//fmt.Printf("Log Name: Transfer %s\n", vLog.Address.Hex())

			var transferEvent LogTransfer

			err := contractAbi.Unpack(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				continue
			}
			if len(vLog.Topics) < 3 {
				//fmt.Println("not enough topics")
				continue
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			//fmt.Printf("From: %s\n", transferEvent.From.Hex())

			if quantaAddr, ok := toAddress[strings.ToLower(transferEvent.To.Hex())]; ok {
				fmt.Printf("To: %s Tok=%s\n", transferEvent.To.Hex(), transferEvent.Tokens.String())

				erc20, err := contracts.NewSimpleToken(vLog.Address, l.Client.(bind.ContractBackend))
				if err != nil {

				}

				name, err := erc20.Name(nil)

				dec, err := erc20.Decimals(nil)
				if err != nil {
					dec = 0
				}

				events = append(events, &Deposit{
					QuantaAddr: quantaAddr,
					BlockID:    int64(vLog.BlockNumber),
					CoinName:   name + vLog.Address.Hex(),
					SenderAddr: transferEvent.To.Hex(),
					Amount:     Erc20AmountToGraphene(*transferEvent.Tokens, dec),
					Tx:         vLog.TxHash.Hex(),
				})
			}

		}
	}

	return events, nil
}

func (l *Listener) GetForwardContract(blockNumber int64) ([]*crypto2.ForwardInput, error) {
	blocks, err := l.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	ABI, err := abi.JSON(strings.NewReader(contracts.QuantaForwarderABI))
	if err != nil {
		return nil, err
	}

	if blocks == nil {
		return nil, errors.New("Block not found ")
	}

	events := []*crypto2.ForwardInput{}
	for _, tx := range blocks.Transactions() {
		data := common.Bytes2Hex(tx.Data())
		//println(data)

		// matches our forwarding contract
		if strings.HasPrefix(data, Forwarder.ForwarderBin) || strings.HasPrefix(data, Forwarder.ForwarderBinV2) || strings.HasPrefix(data, Forwarder.ForwarderBinV3) {
			var remain string
			if strings.HasPrefix(data, Forwarder.ForwarderBin) {
				remain = strings.TrimPrefix(data, Forwarder.ForwarderBin)
			} else if strings.HasPrefix(data, Forwarder.ForwarderBinV2) {
				remain = strings.TrimPrefix(data, Forwarder.ForwarderBinV2)
			} else {
				remain = strings.TrimPrefix(data, Forwarder.ForwarderBinV3)
			}

			input := &crypto2.ForwardInput{}
			vals, err := ABI.Constructor.Inputs.UnpackValues(common.Hex2Bytes(remain))
			if err != nil {
				println("Cannot unpack ", err.Error())
				continue
			}
			if len(vals) != 2 {
				println("Values should be 2")
				continue
			}

			tr, err := l.Client.TransactionReceipt(context.Background(), tx.Hash())

			if err != nil {
				println("Cannot get receipt ", err.Error())
				continue
			}

			input.Blockchain = BLOCKCHAIN_ETH
			input.ContractAddress = tr.ContractAddress.Hex()
			input.Trust = vals[0].(common.Address)
			input.QuantaAddr = vals[1].(string)

			events = append(events, input)
		}
	}

	return events, nil
}

func (l *Listener) SendWithDrawalToRPC(trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {

	return l.SendWithdrawal(l.Client.(bind.ContractBackend), trustAddress, ownerKey, w)
}

func (l *Listener) SendWithdrawal(conn bind.ContractBackend,
	trustAddress common.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {

	auth := bind.NewKeyedTransactor(ownerKey)
	auth.GasLimit = 900000
	contract, err := contracts.NewTrustContract(trustAddress, conn)

	if err != nil {
		return "", err
	}

	var smartAddress common.Address
	parts := strings.Split(w.CoinName, "0X")
	if len(parts) > 1 {
		smartAddress = common.HexToAddress(parts[1])
	}
	//smartAddress = common.HexToAddress(w.CoinName[11:])
	toAddr := common.HexToAddress(w.DestinationAddress)
	amount := GrapheneToWei(w.Amount)
	fmt.Printf("Sending from %s\n", auth.From.Hex())
	fmt.Printf("Submit to contract=%s txId=%d erc20=%s to=%s amount=%s\n", trustAddress.Hex(), w.TxId, smartAddress.Hex(), toAddr.Hex(), amount.String())

	var r [][32]byte
	var s [][32]byte
	var v []uint8

	fmt.Printf("signatures (%d) %v\n", len(w.Signatures), w.Signatures)

	for _, signature := range w.Signatures {
		data := common.Hex2Bytes(signature)
		if len(data) != 65 {
			fmt.Println("Signature is not correct length " + string(len(data)))
			continue
		}

		var r1 [32]byte
		copy(r1[0:32], data[0:32])
		r = append(r, r1)

		var s1 [32]byte
		copy(s1[0:32], data[32:64])
		s = append(s, s1)

		v = append(v, data[64]+27)
	}
	tx, err := contract.PaymentTx(auth, w.TxId, smartAddress, toAddr, amount, v, r, s)
	if err != nil {
		return "", err
	}

	var receipt *types.Receipt
	timeBefore := time.Now()
	for receipt == nil {
		receipt, err = l.Client.TransactionReceipt(context.Background(), tx.Hash())
	}
	if err != nil {
		return tx.Hash().Hex(), errors.New("could not find receipt")
	}
	timeTaken := time.Since(timeBefore)
	fmt.Printf("Successfully submitted transaction %s, receipt status = %d, took %s sec", tx.Hash().Hex(), receipt.Status, timeTaken.String())
	fmt.Println()
	if receipt.Status == types.ReceiptStatusFailed {
		return tx.Hash().Hex(), errors.New("transaction failed")
	}
	return tx.Hash().Hex(), nil
}

func (l *Listener) GetTxID(conn bind.ContractBackend, trustAddress common.Address) (uint64, error) {

	if conn == nil {
		conn = l.Client.(bind.ContractBackend)
	}

	println("Geting txid from trustaddr", trustAddress.Hex())
	contract, err := contracts.NewTrustContract(trustAddress, conn)

	if err != nil {
		return 0, err
	}

	txId, _ := contract.TxIdLast(nil)
	println("Last TX ID= ", txId)
	return txId, nil
}
