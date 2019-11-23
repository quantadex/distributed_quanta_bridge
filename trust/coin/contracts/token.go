// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// SimpleTokenABI is the input ABI used to generate the binding from.
const SimpleTokenABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"},{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isSignedBy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_v\",\"type\":\"uint8\"},{\"name\":\"_r\",\"type\":\"bytes32\"},{\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"recoverSignerVRS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hexstr\",\"type\":\"string\"}],\"name\":\"hexstrToBytes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_char\",\"type\":\"string\"}],\"name\":\"parseInt16Char\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToBytes32\",\"outputs\":[{\"name\":\"b\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"toEthereumSignedMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToString\",\"outputs\":[{\"name\":\"str\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_str\",\"type\":\"string\"},{\"name\":\"_startIndex\",\"type\":\"uint256\"},{\"name\":\"_endIndex\",\"type\":\"uint256\"}],\"name\":\"substring\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_completed_migration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"name\":\"completed\",\"type\":\"uint256\"}],\"name\":\"setCompleted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"signers\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"requiredVotes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"txIdLast\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_CANDIDATES\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"txId\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"erc20Addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"toAddr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"verified\",\"type\":\"bool[]\"}],\"name\":\"TransactionResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Fund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"initialSigners\",\"type\":\"address[]\"}],\"name\":\"assignInitialSigners\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"isSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numSigners\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"getAddCandidateVotes\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numAddCandidates\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"getRemoveCandidateVotes\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numRemoveCandidates\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txId\",\"type\":\"uint64\"},{\"name\":\"erc20Addr\",\"type\":\"address\"},{\"name\":\"toAddr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"paymentTx\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"voteAddSigner\",\"outputs\":[{\"name\":\"votesNeeded\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"voteRemoveSigner\",\"outputs\":[{\"name\":\"votesNeeded\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"quantaAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"destinationAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"trust\",\"type\":\"address\"},{\"name\":\"quantaAddr\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogForwarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogFlushed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"trust\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"quanta\",\"type\":\"string\"}],\"name\":\"LogCreated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenContractAddress\",\"type\":\"address\"}],\"name\":\"flushTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SimpleToken is an auto generated Go binding around an Ethereum contract.
type SimpleToken struct {
	SimpleTokenCaller     // Read-only binding to the contract
	SimpleTokenTransactor // Write-only binding to the contract
	SimpleTokenFilterer   // Log filterer for contract events
}

// SimpleTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleTokenSession struct {
	Contract     *SimpleToken      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleTokenCallerSession struct {
	Contract *SimpleTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SimpleTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleTokenTransactorSession struct {
	Contract     *SimpleTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SimpleTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleTokenRaw struct {
	Contract *SimpleToken // Generic contract binding to access the raw methods on
}

// SimpleTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleTokenCallerRaw struct {
	Contract *SimpleTokenCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleTokenTransactorRaw struct {
	Contract *SimpleTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleToken creates a new instance of SimpleToken, bound to a specific deployed contract.
func NewSimpleToken(address common.Address, backend bind.ContractBackend) (*SimpleToken, error) {
	contract, err := bindSimpleToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleToken{SimpleTokenCaller: SimpleTokenCaller{contract: contract}, SimpleTokenTransactor: SimpleTokenTransactor{contract: contract}, SimpleTokenFilterer: SimpleTokenFilterer{contract: contract}}, nil
}

// NewSimpleTokenCaller creates a new read-only instance of SimpleToken, bound to a specific deployed contract.
func NewSimpleTokenCaller(address common.Address, caller bind.ContractCaller) (*SimpleTokenCaller, error) {
	contract, err := bindSimpleToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenCaller{contract: contract}, nil
}

// NewSimpleTokenTransactor creates a new write-only instance of SimpleToken, bound to a specific deployed contract.
func NewSimpleTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleTokenTransactor, error) {
	contract, err := bindSimpleToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenTransactor{contract: contract}, nil
}

// NewSimpleTokenFilterer creates a new log filterer instance of SimpleToken, bound to a specific deployed contract.
func NewSimpleTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleTokenFilterer, error) {
	contract, err := bindSimpleToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenFilterer{contract: contract}, nil
}

// bindSimpleToken binds a generic wrapper to an already deployed contract.
func bindSimpleToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleToken *SimpleTokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleToken.Contract.SimpleTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleToken *SimpleTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleToken.Contract.SimpleTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleToken *SimpleTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleToken.Contract.SimpleTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleToken *SimpleTokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleToken *SimpleTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleToken *SimpleTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleToken.Contract.contract.Transact(opts, method, params...)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) INITIALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "INITIAL_SUPPLY")
	return *ret0, err
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) INITIALSUPPLY() (*big.Int, error) {
	return _SimpleToken.Contract.INITIALSUPPLY(&_SimpleToken.CallOpts)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) INITIALSUPPLY() (*big.Int, error) {
	return _SimpleToken.Contract.INITIALSUPPLY(&_SimpleToken.CallOpts)
}

// MAXCANDIDATES is a free data retrieval call binding the contract method 0xf0786096.
//
// Solidity: function MAX_CANDIDATES() constant returns(uint8)
func (_SimpleToken *SimpleTokenCaller) MAXCANDIDATES(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "MAX_CANDIDATES")
	return *ret0, err
}

// MAXCANDIDATES is a free data retrieval call binding the contract method 0xf0786096.
//
// Solidity: function MAX_CANDIDATES() constant returns(uint8)
func (_SimpleToken *SimpleTokenSession) MAXCANDIDATES() (uint8, error) {
	return _SimpleToken.Contract.MAXCANDIDATES(&_SimpleToken.CallOpts)
}

// MAXCANDIDATES is a free data retrieval call binding the contract method 0xf0786096.
//
// Solidity: function MAX_CANDIDATES() constant returns(uint8)
func (_SimpleToken *SimpleTokenCallerSession) MAXCANDIDATES() (uint8, error) {
	return _SimpleToken.Contract.MAXCANDIDATES(&_SimpleToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.Allowance(&_SimpleToken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.Allowance(&_SimpleToken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.BalanceOf(&_SimpleToken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.BalanceOf(&_SimpleToken.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_SimpleToken *SimpleTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_SimpleToken *SimpleTokenSession) Decimals() (uint8, error) {
	return _SimpleToken.Contract.Decimals(&_SimpleToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_SimpleToken *SimpleTokenCallerSession) Decimals() (uint8, error) {
	return _SimpleToken.Contract.Decimals(&_SimpleToken.CallOpts)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_SimpleToken *SimpleTokenCaller) DestinationAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "destinationAddress")
	return *ret0, err
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_SimpleToken *SimpleTokenSession) DestinationAddress() (common.Address, error) {
	return _SimpleToken.Contract.DestinationAddress(&_SimpleToken.CallOpts)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_SimpleToken *SimpleTokenCallerSession) DestinationAddress() (common.Address, error) {
	return _SimpleToken.Contract.DestinationAddress(&_SimpleToken.CallOpts)
}

// GetAddCandidateVotes is a free data retrieval call binding the contract method 0x4fee1821.
//
// Solidity: function getAddCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenCaller) GetAddCandidateVotes(opts *bind.CallOpts, candidate common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "getAddCandidateVotes", candidate)
	return *ret0, err
}

// GetAddCandidateVotes is a free data retrieval call binding the contract method 0x4fee1821.
//
// Solidity: function getAddCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenSession) GetAddCandidateVotes(candidate common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.GetAddCandidateVotes(&_SimpleToken.CallOpts, candidate)
}

// GetAddCandidateVotes is a free data retrieval call binding the contract method 0x4fee1821.
//
// Solidity: function getAddCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenCallerSession) GetAddCandidateVotes(candidate common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.GetAddCandidateVotes(&_SimpleToken.CallOpts, candidate)
}

// GetRemoveCandidateVotes is a free data retrieval call binding the contract method 0xf2e108d7.
//
// Solidity: function getRemoveCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenCaller) GetRemoveCandidateVotes(opts *bind.CallOpts, candidate common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "getRemoveCandidateVotes", candidate)
	return *ret0, err
}

// GetRemoveCandidateVotes is a free data retrieval call binding the contract method 0xf2e108d7.
//
// Solidity: function getRemoveCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenSession) GetRemoveCandidateVotes(candidate common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.GetRemoveCandidateVotes(&_SimpleToken.CallOpts, candidate)
}

// GetRemoveCandidateVotes is a free data retrieval call binding the contract method 0xf2e108d7.
//
// Solidity: function getRemoveCandidateVotes(candidate address) constant returns(count uint256)
func (_SimpleToken *SimpleTokenCallerSession) GetRemoveCandidateVotes(candidate common.Address) (*big.Int, error) {
	return _SimpleToken.Contract.GetRemoveCandidateVotes(&_SimpleToken.CallOpts, candidate)
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_SimpleToken *SimpleTokenCaller) HexstrToBytes(opts *bind.CallOpts, _hexstr string) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "hexstrToBytes", _hexstr)
	return *ret0, err
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_SimpleToken *SimpleTokenSession) HexstrToBytes(_hexstr string) ([]byte, error) {
	return _SimpleToken.Contract.HexstrToBytes(&_SimpleToken.CallOpts, _hexstr)
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_SimpleToken *SimpleTokenCallerSession) HexstrToBytes(_hexstr string) ([]byte, error) {
	return _SimpleToken.Contract.HexstrToBytes(&_SimpleToken.CallOpts, _hexstr)
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_SimpleToken *SimpleTokenCaller) IsSignedBy(opts *bind.CallOpts, _hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "isSignedBy", _hashedMsg, _sig, _addr)
	return *ret0, err
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_SimpleToken *SimpleTokenSession) IsSignedBy(_hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	return _SimpleToken.Contract.IsSignedBy(&_SimpleToken.CallOpts, _hashedMsg, _sig, _addr)
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_SimpleToken *SimpleTokenCallerSession) IsSignedBy(_hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	return _SimpleToken.Contract.IsSignedBy(&_SimpleToken.CallOpts, _hashedMsg, _sig, _addr)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_SimpleToken *SimpleTokenCaller) IsSigner(opts *bind.CallOpts, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "isSigner", signer)
	return *ret0, err
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_SimpleToken *SimpleTokenSession) IsSigner(signer common.Address) (bool, error) {
	return _SimpleToken.Contract.IsSigner(&_SimpleToken.CallOpts, signer)
}

// IsSigner is a free data retrieval call binding the contract method 0x7df73e27.
//
// Solidity: function isSigner(signer address) constant returns(bool)
func (_SimpleToken *SimpleTokenCallerSession) IsSigner(signer common.Address) (bool, error) {
	return _SimpleToken.Contract.IsSigner(&_SimpleToken.CallOpts, signer)
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) LastCompletedMigration(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "last_completed_migration")
	return *ret0, err
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) LastCompletedMigration() (*big.Int, error) {
	return _SimpleToken.Contract.LastCompletedMigration(&_SimpleToken.CallOpts)
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) LastCompletedMigration() (*big.Int, error) {
	return _SimpleToken.Contract.LastCompletedMigration(&_SimpleToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_SimpleToken *SimpleTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_SimpleToken *SimpleTokenSession) Name() (string, error) {
	return _SimpleToken.Contract.Name(&_SimpleToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_SimpleToken *SimpleTokenCallerSession) Name() (string, error) {
	return _SimpleToken.Contract.Name(&_SimpleToken.CallOpts)
}

// NumAddCandidates is a free data retrieval call binding the contract method 0x4e63af51.
//
// Solidity: function numAddCandidates() constant returns(count uint256)
func (_SimpleToken *SimpleTokenCaller) NumAddCandidates(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "numAddCandidates")
	return *ret0, err
}

// NumAddCandidates is a free data retrieval call binding the contract method 0x4e63af51.
//
// Solidity: function numAddCandidates() constant returns(count uint256)
func (_SimpleToken *SimpleTokenSession) NumAddCandidates() (*big.Int, error) {
	return _SimpleToken.Contract.NumAddCandidates(&_SimpleToken.CallOpts)
}

// NumAddCandidates is a free data retrieval call binding the contract method 0x4e63af51.
//
// Solidity: function numAddCandidates() constant returns(count uint256)
func (_SimpleToken *SimpleTokenCallerSession) NumAddCandidates() (*big.Int, error) {
	return _SimpleToken.Contract.NumAddCandidates(&_SimpleToken.CallOpts)
}

// NumRemoveCandidates is a free data retrieval call binding the contract method 0x5023d141.
//
// Solidity: function numRemoveCandidates() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) NumRemoveCandidates(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "numRemoveCandidates")
	return *ret0, err
}

// NumRemoveCandidates is a free data retrieval call binding the contract method 0x5023d141.
//
// Solidity: function numRemoveCandidates() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) NumRemoveCandidates() (*big.Int, error) {
	return _SimpleToken.Contract.NumRemoveCandidates(&_SimpleToken.CallOpts)
}

// NumRemoveCandidates is a free data retrieval call binding the contract method 0x5023d141.
//
// Solidity: function numRemoveCandidates() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) NumRemoveCandidates() (*big.Int, error) {
	return _SimpleToken.Contract.NumRemoveCandidates(&_SimpleToken.CallOpts)
}

// NumSigners is a free data retrieval call binding the contract method 0x12679fed.
//
// Solidity: function numSigners() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) NumSigners(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "numSigners")
	return *ret0, err
}

// NumSigners is a free data retrieval call binding the contract method 0x12679fed.
//
// Solidity: function numSigners() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) NumSigners() (*big.Int, error) {
	return _SimpleToken.Contract.NumSigners(&_SimpleToken.CallOpts)
}

// NumSigners is a free data retrieval call binding the contract method 0x12679fed.
//
// Solidity: function numSigners() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) NumSigners() (*big.Int, error) {
	return _SimpleToken.Contract.NumSigners(&_SimpleToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SimpleToken *SimpleTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SimpleToken *SimpleTokenSession) Owner() (common.Address, error) {
	return _SimpleToken.Contract.Owner(&_SimpleToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SimpleToken *SimpleTokenCallerSession) Owner() (common.Address, error) {
	return _SimpleToken.Contract.Owner(&_SimpleToken.CallOpts)
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) ParseInt16Char(opts *bind.CallOpts, _char string) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "parseInt16Char", _char)
	return *ret0, err
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) ParseInt16Char(_char string) (*big.Int, error) {
	return _SimpleToken.Contract.ParseInt16Char(&_SimpleToken.CallOpts, _char)
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) ParseInt16Char(_char string) (*big.Int, error) {
	return _SimpleToken.Contract.ParseInt16Char(&_SimpleToken.CallOpts, _char)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_SimpleToken *SimpleTokenCaller) QuantaAddress(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "quantaAddress")
	return *ret0, err
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_SimpleToken *SimpleTokenSession) QuantaAddress() (string, error) {
	return _SimpleToken.Contract.QuantaAddress(&_SimpleToken.CallOpts)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_SimpleToken *SimpleTokenCallerSession) QuantaAddress() (string, error) {
	return _SimpleToken.Contract.QuantaAddress(&_SimpleToken.CallOpts)
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_SimpleToken *SimpleTokenCaller) RecoverSigner(opts *bind.CallOpts, _hashedMsg [32]byte, _sig string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "recoverSigner", _hashedMsg, _sig)
	return *ret0, err
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_SimpleToken *SimpleTokenSession) RecoverSigner(_hashedMsg [32]byte, _sig string) (common.Address, error) {
	return _SimpleToken.Contract.RecoverSigner(&_SimpleToken.CallOpts, _hashedMsg, _sig)
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_SimpleToken *SimpleTokenCallerSession) RecoverSigner(_hashedMsg [32]byte, _sig string) (common.Address, error) {
	return _SimpleToken.Contract.RecoverSigner(&_SimpleToken.CallOpts, _hashedMsg, _sig)
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_SimpleToken *SimpleTokenCaller) RecoverSignerVRS(opts *bind.CallOpts, _hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "recoverSignerVRS", _hashedMsg, _v, _r, _s)
	return *ret0, err
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_SimpleToken *SimpleTokenSession) RecoverSignerVRS(_hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	return _SimpleToken.Contract.RecoverSignerVRS(&_SimpleToken.CallOpts, _hashedMsg, _v, _r, _s)
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_SimpleToken *SimpleTokenCallerSession) RecoverSignerVRS(_hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	return _SimpleToken.Contract.RecoverSignerVRS(&_SimpleToken.CallOpts, _hashedMsg, _v, _r, _s)
}

// RequiredVotes is a free data retrieval call binding the contract method 0xbd31a4d8.
//
// Solidity: function requiredVotes() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) RequiredVotes(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "requiredVotes")
	return *ret0, err
}

// RequiredVotes is a free data retrieval call binding the contract method 0xbd31a4d8.
//
// Solidity: function requiredVotes() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) RequiredVotes() (*big.Int, error) {
	return _SimpleToken.Contract.RequiredVotes(&_SimpleToken.CallOpts)
}

// RequiredVotes is a free data retrieval call binding the contract method 0xbd31a4d8.
//
// Solidity: function requiredVotes() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) RequiredVotes() (*big.Int, error) {
	return _SimpleToken.Contract.RequiredVotes(&_SimpleToken.CallOpts)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers( uint256) constant returns(address)
func (_SimpleToken *SimpleTokenCaller) Signers(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "signers", arg0)
	return *ret0, err
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers( uint256) constant returns(address)
func (_SimpleToken *SimpleTokenSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _SimpleToken.Contract.Signers(&_SimpleToken.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x2079fb9a.
//
// Solidity: function signers( uint256) constant returns(address)
func (_SimpleToken *SimpleTokenCallerSession) Signers(arg0 *big.Int) (common.Address, error) {
	return _SimpleToken.Contract.Signers(&_SimpleToken.CallOpts, arg0)
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_SimpleToken *SimpleTokenCaller) Substring(opts *bind.CallOpts, _str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "substring", _str, _startIndex, _endIndex)
	return *ret0, err
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_SimpleToken *SimpleTokenSession) Substring(_str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	return _SimpleToken.Contract.Substring(&_SimpleToken.CallOpts, _str, _startIndex, _endIndex)
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_SimpleToken *SimpleTokenCallerSession) Substring(_str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	return _SimpleToken.Contract.Substring(&_SimpleToken.CallOpts, _str, _startIndex, _endIndex)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_SimpleToken *SimpleTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_SimpleToken *SimpleTokenSession) Symbol() (string, error) {
	return _SimpleToken.Contract.Symbol(&_SimpleToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_SimpleToken *SimpleTokenCallerSession) Symbol() (string, error) {
	return _SimpleToken.Contract.Symbol(&_SimpleToken.CallOpts)
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_SimpleToken *SimpleTokenCaller) ToEthereumSignedMessage(opts *bind.CallOpts, _msg string) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "toEthereumSignedMessage", _msg)
	return *ret0, err
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_SimpleToken *SimpleTokenSession) ToEthereumSignedMessage(_msg string) ([32]byte, error) {
	return _SimpleToken.Contract.ToEthereumSignedMessage(&_SimpleToken.CallOpts, _msg)
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_SimpleToken *SimpleTokenCallerSession) ToEthereumSignedMessage(_msg string) ([32]byte, error) {
	return _SimpleToken.Contract.ToEthereumSignedMessage(&_SimpleToken.CallOpts, _msg)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_SimpleToken *SimpleTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_SimpleToken *SimpleTokenSession) TotalSupply() (*big.Int, error) {
	return _SimpleToken.Contract.TotalSupply(&_SimpleToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_SimpleToken *SimpleTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _SimpleToken.Contract.TotalSupply(&_SimpleToken.CallOpts)
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_SimpleToken *SimpleTokenCaller) TxIdLast(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "txIdLast")
	return *ret0, err
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_SimpleToken *SimpleTokenSession) TxIdLast() (uint64, error) {
	return _SimpleToken.Contract.TxIdLast(&_SimpleToken.CallOpts)
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_SimpleToken *SimpleTokenCallerSession) TxIdLast() (uint64, error) {
	return _SimpleToken.Contract.TxIdLast(&_SimpleToken.CallOpts)
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_SimpleToken *SimpleTokenCaller) UintToBytes32(opts *bind.CallOpts, _uint *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "uintToBytes32", _uint)
	return *ret0, err
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_SimpleToken *SimpleTokenSession) UintToBytes32(_uint *big.Int) ([]byte, error) {
	return _SimpleToken.Contract.UintToBytes32(&_SimpleToken.CallOpts, _uint)
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_SimpleToken *SimpleTokenCallerSession) UintToBytes32(_uint *big.Int) ([]byte, error) {
	return _SimpleToken.Contract.UintToBytes32(&_SimpleToken.CallOpts, _uint)
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_SimpleToken *SimpleTokenCaller) UintToString(opts *bind.CallOpts, _uint *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleToken.contract.Call(opts, out, "uintToString", _uint)
	return *ret0, err
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_SimpleToken *SimpleTokenSession) UintToString(_uint *big.Int) (string, error) {
	return _SimpleToken.Contract.UintToString(&_SimpleToken.CallOpts, _uint)
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_SimpleToken *SimpleTokenCallerSession) UintToString(_uint *big.Int) (string, error) {
	return _SimpleToken.Contract.UintToString(&_SimpleToken.CallOpts, _uint)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.Approve(&_SimpleToken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.Approve(&_SimpleToken.TransactOpts, _spender, _value)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_SimpleToken *SimpleTokenTransactor) AssignInitialSigners(opts *bind.TransactOpts, initialSigners []common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "assignInitialSigners", initialSigners)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_SimpleToken *SimpleTokenSession) AssignInitialSigners(initialSigners []common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.AssignInitialSigners(&_SimpleToken.TransactOpts, initialSigners)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_SimpleToken *SimpleTokenTransactorSession) AssignInitialSigners(initialSigners []common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.AssignInitialSigners(&_SimpleToken.TransactOpts, initialSigners)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.DecreaseApproval(&_SimpleToken.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.DecreaseApproval(&_SimpleToken.TransactOpts, _spender, _subtractedValue)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_SimpleToken *SimpleTokenTransactor) Flush(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "flush")
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_SimpleToken *SimpleTokenSession) Flush() (*types.Transaction, error) {
	return _SimpleToken.Contract.Flush(&_SimpleToken.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_SimpleToken *SimpleTokenTransactorSession) Flush() (*types.Transaction, error) {
	return _SimpleToken.Contract.Flush(&_SimpleToken.TransactOpts)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_SimpleToken *SimpleTokenTransactor) FlushTokens(opts *bind.TransactOpts, tokenContractAddress common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "flushTokens", tokenContractAddress)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_SimpleToken *SimpleTokenSession) FlushTokens(tokenContractAddress common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.FlushTokens(&_SimpleToken.TransactOpts, tokenContractAddress)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_SimpleToken *SimpleTokenTransactorSession) FlushTokens(tokenContractAddress common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.FlushTokens(&_SimpleToken.TransactOpts, tokenContractAddress)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.IncreaseApproval(&_SimpleToken.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.IncreaseApproval(&_SimpleToken.TransactOpts, _spender, _addedValue)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_SimpleToken *SimpleTokenTransactor) PaymentTx(opts *bind.TransactOpts, txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "paymentTx", txId, erc20Addr, toAddr, amount, v, r, s)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_SimpleToken *SimpleTokenSession) PaymentTx(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _SimpleToken.Contract.PaymentTx(&_SimpleToken.TransactOpts, txId, erc20Addr, toAddr, amount, v, r, s)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_SimpleToken *SimpleTokenTransactorSession) PaymentTx(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _SimpleToken.Contract.PaymentTx(&_SimpleToken.TransactOpts, txId, erc20Addr, toAddr, amount, v, r, s)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleToken *SimpleTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleToken *SimpleTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimpleToken.Contract.RenounceOwnership(&_SimpleToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleToken *SimpleTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimpleToken.Contract.RenounceOwnership(&_SimpleToken.TransactOpts)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_SimpleToken *SimpleTokenTransactor) SetCompleted(opts *bind.TransactOpts, completed *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "setCompleted", completed)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_SimpleToken *SimpleTokenSession) SetCompleted(completed *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.SetCompleted(&_SimpleToken.TransactOpts, completed)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_SimpleToken *SimpleTokenTransactorSession) SetCompleted(completed *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.SetCompleted(&_SimpleToken.TransactOpts, completed)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.Transfer(&_SimpleToken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.Transfer(&_SimpleToken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.TransferFrom(&_SimpleToken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_SimpleToken *SimpleTokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SimpleToken.Contract.TransferFrom(&_SimpleToken.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SimpleToken *SimpleTokenTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SimpleToken *SimpleTokenSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.TransferOwnership(&_SimpleToken.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_SimpleToken *SimpleTokenTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.TransferOwnership(&_SimpleToken.TransactOpts, _newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_SimpleToken *SimpleTokenTransactor) Upgrade(opts *bind.TransactOpts, new_address common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "upgrade", new_address)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_SimpleToken *SimpleTokenSession) Upgrade(new_address common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.Upgrade(&_SimpleToken.TransactOpts, new_address)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_SimpleToken *SimpleTokenTransactorSession) Upgrade(new_address common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.Upgrade(&_SimpleToken.TransactOpts, new_address)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenTransactor) VoteAddSigner(opts *bind.TransactOpts, candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "voteAddSigner", candidate)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenSession) VoteAddSigner(candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.VoteAddSigner(&_SimpleToken.TransactOpts, candidate)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenTransactorSession) VoteAddSigner(candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.VoteAddSigner(&_SimpleToken.TransactOpts, candidate)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenTransactor) VoteRemoveSigner(opts *bind.TransactOpts, candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.contract.Transact(opts, "voteRemoveSigner", candidate)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenSession) VoteRemoveSigner(candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.VoteRemoveSigner(&_SimpleToken.TransactOpts, candidate)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(candidate address) returns(votesNeeded uint256)
func (_SimpleToken *SimpleTokenTransactorSession) VoteRemoveSigner(candidate common.Address) (*types.Transaction, error) {
	return _SimpleToken.Contract.VoteRemoveSigner(&_SimpleToken.TransactOpts, candidate)
}

// SimpleTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SimpleToken contract.
type SimpleTokenApprovalIterator struct {
	Event *SimpleTokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenApproval represents a Approval event raised by the SimpleToken contract.
type SimpleTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_SimpleToken *SimpleTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SimpleTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenApprovalIterator{contract: _SimpleToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_SimpleToken *SimpleTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SimpleTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenApproval)
				if err := _SimpleToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenFundIterator is returned from FilterFund and is used to iterate over the raw logs and unpacked data for Fund events raised by the SimpleToken contract.
type SimpleTokenFundIterator struct {
	Event *SimpleTokenFund // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenFundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenFund)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenFund)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenFundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenFundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenFund represents a Fund event raised by the SimpleToken contract.
type SimpleTokenFund struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFund is a free log retrieval operation binding the contract event 0xda8220a878ff7a89474ccffdaa31ea1ed1ffbb0207d5051afccc4fbaf81f9bcd.
//
// Solidity: e Fund(_from indexed address, _value uint256)
func (_SimpleToken *SimpleTokenFilterer) FilterFund(opts *bind.FilterOpts, _from []common.Address) (*SimpleTokenFundIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "Fund", _fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenFundIterator{contract: _SimpleToken.contract, event: "Fund", logs: logs, sub: sub}, nil
}

// WatchFund is a free log subscription operation binding the contract event 0xda8220a878ff7a89474ccffdaa31ea1ed1ffbb0207d5051afccc4fbaf81f9bcd.
//
// Solidity: e Fund(_from indexed address, _value uint256)
func (_SimpleToken *SimpleTokenFilterer) WatchFund(opts *bind.WatchOpts, sink chan<- *SimpleTokenFund, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "Fund", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenFund)
				if err := _SimpleToken.contract.UnpackLog(event, "Fund", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenLogCreatedIterator is returned from FilterLogCreated and is used to iterate over the raw logs and unpacked data for LogCreated events raised by the SimpleToken contract.
type SimpleTokenLogCreatedIterator struct {
	Event *SimpleTokenLogCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenLogCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenLogCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenLogCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenLogCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenLogCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenLogCreated represents a LogCreated event raised by the SimpleToken contract.
type SimpleTokenLogCreated struct {
	Trust  common.Address
	Quanta string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogCreated is a free log retrieval operation binding the contract event 0xa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a.
//
// Solidity: e LogCreated(trust address, quanta string)
func (_SimpleToken *SimpleTokenFilterer) FilterLogCreated(opts *bind.FilterOpts) (*SimpleTokenLogCreatedIterator, error) {

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "LogCreated")
	if err != nil {
		return nil, err
	}
	return &SimpleTokenLogCreatedIterator{contract: _SimpleToken.contract, event: "LogCreated", logs: logs, sub: sub}, nil
}

// WatchLogCreated is a free log subscription operation binding the contract event 0xa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a.
//
// Solidity: e LogCreated(trust address, quanta string)
func (_SimpleToken *SimpleTokenFilterer) WatchLogCreated(opts *bind.WatchOpts, sink chan<- *SimpleTokenLogCreated) (event.Subscription, error) {

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "LogCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenLogCreated)
				if err := _SimpleToken.contract.UnpackLog(event, "LogCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenLogFlushedIterator is returned from FilterLogFlushed and is used to iterate over the raw logs and unpacked data for LogFlushed events raised by the SimpleToken contract.
type SimpleTokenLogFlushedIterator struct {
	Event *SimpleTokenLogFlushed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenLogFlushedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenLogFlushed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenLogFlushed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenLogFlushedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenLogFlushedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenLogFlushed represents a LogFlushed event raised by the SimpleToken contract.
type SimpleTokenLogFlushed struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogFlushed is a free log retrieval operation binding the contract event 0xa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea.
//
// Solidity: e LogFlushed(sender indexed address, amount uint256)
func (_SimpleToken *SimpleTokenFilterer) FilterLogFlushed(opts *bind.FilterOpts, sender []common.Address) (*SimpleTokenLogFlushedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "LogFlushed", senderRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenLogFlushedIterator{contract: _SimpleToken.contract, event: "LogFlushed", logs: logs, sub: sub}, nil
}

// WatchLogFlushed is a free log subscription operation binding the contract event 0xa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea.
//
// Solidity: e LogFlushed(sender indexed address, amount uint256)
func (_SimpleToken *SimpleTokenFilterer) WatchLogFlushed(opts *bind.WatchOpts, sink chan<- *SimpleTokenLogFlushed, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "LogFlushed", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenLogFlushed)
				if err := _SimpleToken.contract.UnpackLog(event, "LogFlushed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenLogForwardedIterator is returned from FilterLogForwarded and is used to iterate over the raw logs and unpacked data for LogForwarded events raised by the SimpleToken contract.
type SimpleTokenLogForwardedIterator struct {
	Event *SimpleTokenLogForwarded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenLogForwardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenLogForwarded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenLogForwarded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenLogForwardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenLogForwardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenLogForwarded represents a LogForwarded event raised by the SimpleToken contract.
type SimpleTokenLogForwarded struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogForwarded is a free log retrieval operation binding the contract event 0x5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6.
//
// Solidity: e LogForwarded(sender indexed address, amount uint256)
func (_SimpleToken *SimpleTokenFilterer) FilterLogForwarded(opts *bind.FilterOpts, sender []common.Address) (*SimpleTokenLogForwardedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "LogForwarded", senderRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenLogForwardedIterator{contract: _SimpleToken.contract, event: "LogForwarded", logs: logs, sub: sub}, nil
}

// WatchLogForwarded is a free log subscription operation binding the contract event 0x5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6.
//
// Solidity: e LogForwarded(sender indexed address, amount uint256)
func (_SimpleToken *SimpleTokenFilterer) WatchLogForwarded(opts *bind.WatchOpts, sink chan<- *SimpleTokenLogForwarded, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "LogForwarded", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenLogForwarded)
				if err := _SimpleToken.contract.UnpackLog(event, "LogForwarded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the SimpleToken contract.
type SimpleTokenOwnershipRenouncedIterator struct {
	Event *SimpleTokenOwnershipRenounced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenOwnershipRenounced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenOwnershipRenounced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenOwnershipRenounced represents a OwnershipRenounced event raised by the SimpleToken contract.
type SimpleTokenOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_SimpleToken *SimpleTokenFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*SimpleTokenOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenOwnershipRenouncedIterator{contract: _SimpleToken.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_SimpleToken *SimpleTokenFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *SimpleTokenOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenOwnershipRenounced)
				if err := _SimpleToken.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimpleToken contract.
type SimpleTokenOwnershipTransferredIterator struct {
	Event *SimpleTokenOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenOwnershipTransferred represents a OwnershipTransferred event raised by the SimpleToken contract.
type SimpleTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_SimpleToken *SimpleTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SimpleTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenOwnershipTransferredIterator{contract: _SimpleToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_SimpleToken *SimpleTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenOwnershipTransferred)
				if err := _SimpleToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenTransactionResultIterator is returned from FilterTransactionResult and is used to iterate over the raw logs and unpacked data for TransactionResult events raised by the SimpleToken contract.
type SimpleTokenTransactionResultIterator struct {
	Event *SimpleTokenTransactionResult // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenTransactionResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenTransactionResult)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenTransactionResult)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenTransactionResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenTransactionResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenTransactionResult represents a TransactionResult event raised by the SimpleToken contract.
type SimpleTokenTransactionResult struct {
	Success   bool
	TxId      uint64
	Erc20Addr common.Address
	ToAddr    common.Address
	Amount    *big.Int
	Verified  []bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransactionResult is a free log retrieval operation binding the contract event 0x060aac1656ca90ad85bf3d9126ac8709f68d4403c923998b24ea513a9e1ba3a6.
//
// Solidity: e TransactionResult(success bool, txId uint64, erc20Addr address, toAddr address, amount uint256, verified bool[])
func (_SimpleToken *SimpleTokenFilterer) FilterTransactionResult(opts *bind.FilterOpts) (*SimpleTokenTransactionResultIterator, error) {

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "TransactionResult")
	if err != nil {
		return nil, err
	}
	return &SimpleTokenTransactionResultIterator{contract: _SimpleToken.contract, event: "TransactionResult", logs: logs, sub: sub}, nil
}

// WatchTransactionResult is a free log subscription operation binding the contract event 0x060aac1656ca90ad85bf3d9126ac8709f68d4403c923998b24ea513a9e1ba3a6.
//
// Solidity: e TransactionResult(success bool, txId uint64, erc20Addr address, toAddr address, amount uint256, verified bool[])
func (_SimpleToken *SimpleTokenFilterer) WatchTransactionResult(opts *bind.WatchOpts, sink chan<- *SimpleTokenTransactionResult) (event.Subscription, error) {

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "TransactionResult")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenTransactionResult)
				if err := _SimpleToken.contract.UnpackLog(event, "TransactionResult", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SimpleToken contract.
type SimpleTokenTransferIterator struct {
	Event *SimpleTokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleTokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleTokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleTokenTransfer represents a Transfer event raised by the SimpleToken contract.
type SimpleTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_SimpleToken *SimpleTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenTransferIterator{contract: _SimpleToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_SimpleToken *SimpleTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SimpleTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SimpleToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleTokenTransfer)
				if err := _SimpleToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
