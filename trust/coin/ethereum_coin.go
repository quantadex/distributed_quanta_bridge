package coin

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	errors2 "github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common"
	"github.com/quantadex/distributed_quanta_bridge/common/crypto"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const sign_prefix = "\x19Ethereum Signed Message:\n"

type EthereumCoin struct {
	client         *Listener
	maxRange       int64
	networkId      string
	ethereumRpc    string
	ethereumSecret *ecdsa.PrivateKey
	erc20map       map[string]string
}

type EncodedMsg struct {
	Message     string
	Tx          string
	BlockNumber int64
	CoinName    string
}

func (c *EthereumCoin) Blockchain() string {
	return "ETH"
}

func (c *EthereumCoin) Attach() error {
	c.client = &Listener{NetworkID: c.networkId}
	ethereumClient, err := ethclient.Dial(c.ethereumRpc)
	if err != nil {
		return err
	}

	c.client.Client = ethereumClient
	err = c.client.Start(c.erc20map)
	if err != nil {
		return err
	}

	return nil
}

func (c *EthereumCoin) GetBlockInfo(hash string) (string, int64, error) {
	_, blockHash, blockNumber, _, err := c.client.GetTransactionbyHash(hash)
	var blockNum int64
	if strings.HasPrefix(strings.ToLower(blockNumber), "0x") {
		blockNum, err = strconv.ParseInt(blockNumber[2:], 16, 64)
		if err != nil {
			return "", 0, errors2.Wrap(err, "Could not convert block number")
		}
	} else {
		blockNum, err = strconv.ParseInt(blockNumber[2:], 16, 64)
		if err != nil {
			return "", 0, errors2.Wrap(err, "Could not convert block number")
		}
	}
	topBlock, err := c.GetTopBlockID()
	if err != nil {
		return "", 0, errors2.Wrap(err, "Could not get top block id")
	}
	confirm := topBlock - blockNum
	return blockHash.String(), confirm, err
}

func (c *EthereumCoin) GetTopBlockID() (int64, error) {
	topBlockId, err := c.client.GetTopBlockNumber()
	if err != nil {
		return 0, err
	}

	return common.Min64(c.maxRange, topBlockId), nil
}

func (c *EthereumCoin) GetTxID(trustAddress common2.Address) (uint64, error) {
	return c.client.GetTxID(nil, trustAddress)
}

func (c *EthereumCoin) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	ndeposits, err := c.client.GetNativeDeposits(blockID, trustAddress)
	if err != nil {
		return nil, err
	}
	deps, err := c.client.FilterTransferEvent(blockID, trustAddress)
	if err != nil {
		return nil, err
	}

	return append(ndeposits, deps...), nil
}

func (c *EthereumCoin) FlushCoin(forwarderAddr string, tokenAddr string) error {
	forwarder, err := contracts.NewQuantaForwarder(common2.HexToAddress(forwarderAddr), c.client.Client.(bind.ContractBackend))
	if err != nil {
		return err
	}

	if forwarder == nil {
		return errors.New("Unable to instantiate forwarding address for " + forwarderAddr)
	}

	auth := bind.NewKeyedTransactor(c.ethereumSecret)

	tx, err := forwarder.FlushTokens(auth, common2.HexToAddress(tokenAddr))
	if tx != nil {
		println("Flush coin ", tx.Hash().String())
	}
	return err
}

func (c *EthereumCoin) GetPendingTx(map[string]string) ([]*Deposit, error) {
	return nil, nil
}

func (c *EthereumCoin) GetForwardersInBlock(blockID int64) ([]*crypto.ForwardInput, error) {
	forwarders, err := c.client.GetForwardContract(blockID)
	if err != nil {
		return nil, err
	}
	return forwarders, nil
}

func (b *EthereumCoin) GenerateMultisig(accountId string) (string, error) {
	panic("not implemented")
}

//
func (c *EthereumCoin) SendWithdrawal(trustAddress common2.Address,
	ownerKey *ecdsa.PrivateKey,
	w *Withdrawal) (string, error) {
	return c.client.SendWithDrawalToRPC(trustAddress, ownerKey, w)
}

func (c *EthereumCoin) FillCrosschainAddress(crosschainAddr map[string]string) {

}

func (c *EthereumCoin) EncodeRefund(w Withdrawal) (string, error) {
	var encoded bytes.Buffer
	var smartAddress string

	parts := strings.Split(w.CoinName, "0X")
	var amount *big.Int
	if len(parts) == 2 {
		smartAddress = parts[1]
		erc20, err := contracts.NewSimpleToken(common2.HexToAddress(smartAddress), c.client.Client.(bind.ContractBackend))
		if err != nil {
			return "", err
		}
		dec, err := erc20.Decimals(nil)
		if err != nil {
			return "", err
		}
		amount = GrapheneToERC20(*new(big.Int).SetUint64(w.Amount), 5, int(dec))

	} else {
		smartAddress = ""
		amount = GrapheneToWei(w.Amount)
	}
	//smartAddress = w.CoinName[11:]
	//fmt.Println("Smart address = ", smartAddress)

	//encoded.WriteString(sign_prefix + "80")
	binary.Write(&encoded, binary.BigEndian, uint64(w.TxId))
	encoded.Write(common2.HexToAddress(strings.ToLower(smartAddress)).Bytes())
	encoded.Write(common2.HexToAddress(strings.ToLower(w.DestinationAddress)).Bytes())
	//encoded.Write(abi.U256(new(big.Int).SetUint64(uint64(w.Amount))))
	encoded.Write(abi.U256(amount))
	//binary.Write(&encoded, binary.BigEndian, abi.U256(new(big.Int).SetUint64(uint64(w.Amount))))

	//println("# of bytes " , encoded.Len(), common2.Bytes2Hex(encoded.Bytes()))
	data, err := json.Marshal(&EncodedMsg{common2.Bytes2Hex(encoded.Bytes()), w.Tx, w.QuantaBlockID, w.CoinName})
	//return common2.Bytes2Hex(data), err
	return string(data), err
}

func (c *EthereumCoin) DecodeRefund(encoded string) (*Withdrawal, error) {
	//decoded := common2.Hex2Bytes(encoded)
	msg := &EncodedMsg{}
	err := json.Unmarshal([]byte(encoded), msg)
	if err != nil {
		return nil, err
	}

	w := &Withdrawal{}
	w.Tx = msg.Tx
	w.QuantaBlockID = msg.BlockNumber
	decoded := common2.Hex2Bytes(msg.Message)

	//pl := len(sign_prefix)
	//header := decoded[0:pl]
	//if string(header) != sign_prefix {
	//	return nil, errors.New("Unexpected prefix")
	//}

	// skip 4 bytes to length
	pl := 0
	txIdBytes := decoded[pl : pl+8]
	txId := new(big.Int).SetBytes(txIdBytes).Uint64()
	w.TxId = txId

	pl += 8
	smartAddress := decoded[pl : pl+20]

	pl += 20
	destAddress := decoded[pl : pl+20]

	pl += 20
	amount := decoded[pl : pl+32]
	smartNumber := new(big.Int).SetBytes(smartAddress)

	if smartNumber.Cmp(big.NewInt(0)) == 0 {
		w.CoinName = "ETH"
	} else {
		w.CoinName = "," + strings.ToLower(common2.BytesToAddress(smartAddress).Hex())
	}

	w.DestinationAddress = strings.ToLower(common2.BytesToAddress(destAddress).Hex())
	w.Amount = new(big.Int).SetBytes(amount).Uint64()

	return w, nil
}

func (c *EthereumCoin) CheckValidAddress(address string) bool {
	var add [20]byte
	copy(add[:], address)
	Ma := common2.NewMixedcaseAddress(add)
	var validAddress = regexp.MustCompile(`^[0x]+[0-9a-fA-F]{40}$`)
	return Ma.ValidChecksum() && validAddress.MatchString(address)
}
