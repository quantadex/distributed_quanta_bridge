package coin

import (
	"context"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"fmt"
	"strings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
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

	println("Created RPC listener with id " + id.String(), l.NetworkID)

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
			return nil, nil
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

func (l *Listener) GetTopBlockNumber() (int64, error){
	header, err := l.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}

func (l *Listener) GetNativeDeposits(blockNumber int64, toAddress string) ([]*Deposit, error) {
	blocks, err := l.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	if blocks == nil {
		return nil, errors.Wrap(err, "Block not found " + err.Error())
	}

	filterAddress := common.HexToAddress(toAddress)
	events := []*Deposit{}
	for _, tx := range blocks.Transactions() {
		if filterAddress.Hex() != tx.To().Hex() {
			continue
		}
		if tx.Value().Cmp(big.NewInt(0)) != 0 {
			events = append(events, &Deposit{
				CoinName: "ETH",
				Amount: WeiToStellar(tx.Value().Int64()),
			})
		}
	}

	return events, nil
}


func (l *Listener) FilterTransferEvent(blockNumber int64, toAddress string) ([]*Deposit, error)  {
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
	filterAddress := common.HexToAddress(toAddress)

	fmt.Printf("Number of log events %d transferHash=%s\n", len(logsEvent), logTransferSigHash.Hex())
	events := []*Deposit{}
	for _, vLog := range logsEvent {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d %s\n", vLog.Index, vLog.Topics[0].Hex())
		fmt.Println(vLog.TxHash.Hex())
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():

			fmt.Printf("Log Name: Transfer %s\n", vLog.Address.Hex())

			var transferEvent LogTransfer

			err := contractAbi.Unpack(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				continue
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			if filterAddress == transferEvent.To {
				fmt.Printf("From: %s\n", transferEvent.From.Hex())
				fmt.Printf("To: %s\n", transferEvent.To.Hex())
				fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())

				events = append(events, &Deposit{
					CoinName: vLog.Address.Hex(),
					Amount: WeiToStellar(transferEvent.Tokens.Int64()),
				})
			}

		}
	}

	return events, nil
}


func (l *Listener) GetForwardContract(blockNumber int64) ([]*ForwardInput, error) {
	blocks, err := l.GetBlock(blockNumber)
	if err != nil {
		return nil, err
	}

	ABI, err := abi.JSON(strings.NewReader(Forwarder.ForwarderABI))
	if err != nil {
		return nil, err
	}

	data, err := ABI.Pack("", common.HexToAddress("0xe0006458963c3773b051e767c5c63fee24cd7ff9"),"QQQWEQWE")
	if err != nil {
		return nil, err
	}
	println("input", common.Bytes2Hex(data))

	if blocks == nil {
		return nil, errors.New("Block not found ")
	}


	events := []*ForwardInput{}
	for _, tx := range blocks.Transactions() {
		data := common.Bytes2Hex(tx.Data())
		println(data)

		// matches our forwarding contract
		if strings.HasPrefix(data, Forwarder.ForwarderBin) {
			remain := strings.TrimPrefix(data, Forwarder.ForwarderBin)

			input := &ForwardInput{}
			err = ABI.Unpack(input, "", common.Hex2Bytes(remain))
			if err != nil {
				println("Cannot unpack ", err)
				continue
			}

			events = append(events, input)
		}
	}

	return events, nil
}

