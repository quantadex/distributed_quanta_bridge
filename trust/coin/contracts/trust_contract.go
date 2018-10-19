// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// TrustContractABI is the input ABI used to generate the binding from.
const TrustContractABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"},{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isSignedBy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_v\",\"type\":\"uint8\"},{\"name\":\"_r\",\"type\":\"bytes32\"},{\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"recoverSignerVRS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hexstr\",\"type\":\"string\"}],\"name\":\"hexstrToBytes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_char\",\"type\":\"string\"}],\"name\":\"parseInt16Char\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToBytes32\",\"outputs\":[{\"name\":\"b\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"toEthereumSignedMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToString\",\"outputs\":[{\"name\":\"str\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_str\",\"type\":\"string\"},{\"name\":\"_startIndex\",\"type\":\"uint256\"},{\"name\":\"_endIndex\",\"type\":\"uint256\"}],\"name\":\"substring\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_completed_migration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"name\":\"completed\",\"type\":\"uint256\"}],\"name\":\"setCompleted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"txIdLast\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"txId\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"erc20Addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"toAddr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"verified\",\"type\":\"bool[]\"}],\"name\":\"TransactionResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Fund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"initialSigners\",\"type\":\"address[]\"}],\"name\":\"assignInitialSigners\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalSigners\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txId\",\"type\":\"uint64\"},{\"name\":\"erc20Addr\",\"type\":\"address\"},{\"name\":\"toAddr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"paymentTx\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"voteAddSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"voteRemoveSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INITIAL_SUPPLY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseApproval\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TrustContract is an auto generated Go binding around an Ethereum contract.
type TrustContract struct {
	TrustContractCaller     // Read-only binding to the contract
	TrustContractTransactor // Write-only binding to the contract
	TrustContractFilterer   // Log filterer for contract events
}

// TrustContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type TrustContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TrustContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TrustContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrustContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TrustContractSession struct {
	Contract     *TrustContract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TrustContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TrustContractCallerSession struct {
	Contract *TrustContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TrustContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TrustContractTransactorSession struct {
	Contract     *TrustContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TrustContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type TrustContractRaw struct {
	Contract *TrustContract // Generic contract binding to access the raw methods on
}

// TrustContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TrustContractCallerRaw struct {
	Contract *TrustContractCaller // Generic read-only contract binding to access the raw methods on
}

// TrustContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TrustContractTransactorRaw struct {
	Contract *TrustContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrustContract creates a new instance of TrustContract, bound to a specific deployed contract.
func NewTrustContract(address common.Address, backend bind.ContractBackend) (*TrustContract, error) {
	contract, err := bindTrustContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TrustContract{TrustContractCaller: TrustContractCaller{contract: contract}, TrustContractTransactor: TrustContractTransactor{contract: contract}, TrustContractFilterer: TrustContractFilterer{contract: contract}}, nil
}

// NewTrustContractCaller creates a new read-only instance of TrustContract, bound to a specific deployed contract.
func NewTrustContractCaller(address common.Address, caller bind.ContractCaller) (*TrustContractCaller, error) {
	contract, err := bindTrustContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrustContractCaller{contract: contract}, nil
}

// NewTrustContractTransactor creates a new write-only instance of TrustContract, bound to a specific deployed contract.
func NewTrustContractTransactor(address common.Address, transactor bind.ContractTransactor) (*TrustContractTransactor, error) {
	contract, err := bindTrustContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrustContractTransactor{contract: contract}, nil
}

// NewTrustContractFilterer creates a new log filterer instance of TrustContract, bound to a specific deployed contract.
func NewTrustContractFilterer(address common.Address, filterer bind.ContractFilterer) (*TrustContractFilterer, error) {
	contract, err := bindTrustContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrustContractFilterer{contract: contract}, nil
}

// bindTrustContract binds a generic wrapper to an already deployed contract.
func bindTrustContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TrustContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustContract *TrustContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TrustContract.Contract.TrustContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustContract *TrustContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustContract.Contract.TrustContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustContract *TrustContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustContract.Contract.TrustContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrustContract *TrustContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TrustContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrustContract *TrustContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrustContract *TrustContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrustContract.Contract.contract.Transact(opts, method, params...)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_TrustContract *TrustContractCaller) INITIALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "INITIAL_SUPPLY")
	return *ret0, err
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_TrustContract *TrustContractSession) INITIALSUPPLY() (*big.Int, error) {
	return _TrustContract.Contract.INITIALSUPPLY(&_TrustContract.CallOpts)
}

// INITIALSUPPLY is a free data retrieval call binding the contract method 0x2ff2e9dc.
//
// Solidity: function INITIAL_SUPPLY() constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) INITIALSUPPLY() (*big.Int, error) {
	return _TrustContract.Contract.INITIALSUPPLY(&_TrustContract.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_TrustContract *TrustContractCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "allowance", _owner, _spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_TrustContract *TrustContractSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _TrustContract.Contract.Allowance(&_TrustContract.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(_owner address, _spender address) constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _TrustContract.Contract.Allowance(&_TrustContract.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_TrustContract *TrustContractCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_TrustContract *TrustContractSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _TrustContract.Contract.BalanceOf(&_TrustContract.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _TrustContract.Contract.BalanceOf(&_TrustContract.CallOpts, _owner)
}

// GetTotalSigners is a free data retrieval call binding the contract method 0x2cef7165.
//
// Solidity: function getTotalSigners() constant returns(uint256)
func (_TrustContract *TrustContractCaller) GetTotalSigners(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "getTotalSigners")
	return *ret0, err
}

// GetTotalSigners is a free data retrieval call binding the contract method 0x2cef7165.
//
// Solidity: function getTotalSigners() constant returns(uint256)
func (_TrustContract *TrustContractSession) GetTotalSigners() (*big.Int, error) {
	return _TrustContract.Contract.GetTotalSigners(&_TrustContract.CallOpts)
}

// GetTotalSigners is a free data retrieval call binding the contract method 0x2cef7165.
//
// Solidity: function getTotalSigners() constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) GetTotalSigners() (*big.Int, error) {
	return _TrustContract.Contract.GetTotalSigners(&_TrustContract.CallOpts)
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_TrustContract *TrustContractCaller) HexstrToBytes(opts *bind.CallOpts, _hexstr string) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "hexstrToBytes", _hexstr)
	return *ret0, err
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_TrustContract *TrustContractSession) HexstrToBytes(_hexstr string) ([]byte, error) {
	return _TrustContract.Contract.HexstrToBytes(&_TrustContract.CallOpts, _hexstr)
}

// HexstrToBytes is a free data retrieval call binding the contract method 0x1445f713.
//
// Solidity: function hexstrToBytes(_hexstr string) constant returns(bytes)
func (_TrustContract *TrustContractCallerSession) HexstrToBytes(_hexstr string) ([]byte, error) {
	return _TrustContract.Contract.HexstrToBytes(&_TrustContract.CallOpts, _hexstr)
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_TrustContract *TrustContractCaller) IsSignedBy(opts *bind.CallOpts, _hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "isSignedBy", _hashedMsg, _sig, _addr)
	return *ret0, err
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_TrustContract *TrustContractSession) IsSignedBy(_hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	return _TrustContract.Contract.IsSignedBy(&_TrustContract.CallOpts, _hashedMsg, _sig, _addr)
}

// IsSignedBy is a free data retrieval call binding the contract method 0x1052506f.
//
// Solidity: function isSignedBy(_hashedMsg bytes32, _sig string, _addr address) constant returns(bool)
func (_TrustContract *TrustContractCallerSession) IsSignedBy(_hashedMsg [32]byte, _sig string, _addr common.Address) (bool, error) {
	return _TrustContract.Contract.IsSignedBy(&_TrustContract.CallOpts, _hashedMsg, _sig, _addr)
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_TrustContract *TrustContractCaller) LastCompletedMigration(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "last_completed_migration")
	return *ret0, err
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_TrustContract *TrustContractSession) LastCompletedMigration() (*big.Int, error) {
	return _TrustContract.Contract.LastCompletedMigration(&_TrustContract.CallOpts)
}

// LastCompletedMigration is a free data retrieval call binding the contract method 0x445df0ac.
//
// Solidity: function last_completed_migration() constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) LastCompletedMigration() (*big.Int, error) {
	return _TrustContract.Contract.LastCompletedMigration(&_TrustContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TrustContract *TrustContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TrustContract *TrustContractSession) Name() (string, error) {
	return _TrustContract.Contract.Name(&_TrustContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TrustContract *TrustContractCallerSession) Name() (string, error) {
	return _TrustContract.Contract.Name(&_TrustContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TrustContract *TrustContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TrustContract *TrustContractSession) Owner() (common.Address, error) {
	return _TrustContract.Contract.Owner(&_TrustContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_TrustContract *TrustContractCallerSession) Owner() (common.Address, error) {
	return _TrustContract.Contract.Owner(&_TrustContract.CallOpts)
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_TrustContract *TrustContractCaller) ParseInt16Char(opts *bind.CallOpts, _char string) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "parseInt16Char", _char)
	return *ret0, err
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_TrustContract *TrustContractSession) ParseInt16Char(_char string) (*big.Int, error) {
	return _TrustContract.Contract.ParseInt16Char(&_TrustContract.CallOpts, _char)
}

// ParseInt16Char is a free data retrieval call binding the contract method 0x38b025b2.
//
// Solidity: function parseInt16Char(_char string) constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) ParseInt16Char(_char string) (*big.Int, error) {
	return _TrustContract.Contract.ParseInt16Char(&_TrustContract.CallOpts, _char)
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_TrustContract *TrustContractCaller) RecoverSigner(opts *bind.CallOpts, _hashedMsg [32]byte, _sig string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "recoverSigner", _hashedMsg, _sig)
	return *ret0, err
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_TrustContract *TrustContractSession) RecoverSigner(_hashedMsg [32]byte, _sig string) (common.Address, error) {
	return _TrustContract.Contract.RecoverSigner(&_TrustContract.CallOpts, _hashedMsg, _sig)
}

// RecoverSigner is a free data retrieval call binding the contract method 0xdca95419.
//
// Solidity: function recoverSigner(_hashedMsg bytes32, _sig string) constant returns(address)
func (_TrustContract *TrustContractCallerSession) RecoverSigner(_hashedMsg [32]byte, _sig string) (common.Address, error) {
	return _TrustContract.Contract.RecoverSigner(&_TrustContract.CallOpts, _hashedMsg, _sig)
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_TrustContract *TrustContractCaller) RecoverSignerVRS(opts *bind.CallOpts, _hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "recoverSignerVRS", _hashedMsg, _v, _r, _s)
	return *ret0, err
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_TrustContract *TrustContractSession) RecoverSignerVRS(_hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	return _TrustContract.Contract.RecoverSignerVRS(&_TrustContract.CallOpts, _hashedMsg, _v, _r, _s)
}

// RecoverSignerVRS is a free data retrieval call binding the contract method 0x36cae648.
//
// Solidity: function recoverSignerVRS(_hashedMsg bytes32, _v uint8, _r bytes32, _s bytes32) constant returns(address)
func (_TrustContract *TrustContractCallerSession) RecoverSignerVRS(_hashedMsg [32]byte, _v uint8, _r [32]byte, _s [32]byte) (common.Address, error) {
	return _TrustContract.Contract.RecoverSignerVRS(&_TrustContract.CallOpts, _hashedMsg, _v, _r, _s)
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_TrustContract *TrustContractCaller) Substring(opts *bind.CallOpts, _str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "substring", _str, _startIndex, _endIndex)
	return *ret0, err
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_TrustContract *TrustContractSession) Substring(_str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	return _TrustContract.Contract.Substring(&_TrustContract.CallOpts, _str, _startIndex, _endIndex)
}

// Substring is a free data retrieval call binding the contract method 0x1dcd9b55.
//
// Solidity: function substring(_str string, _startIndex uint256, _endIndex uint256) constant returns(string)
func (_TrustContract *TrustContractCallerSession) Substring(_str string, _startIndex *big.Int, _endIndex *big.Int) (string, error) {
	return _TrustContract.Contract.Substring(&_TrustContract.CallOpts, _str, _startIndex, _endIndex)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TrustContract *TrustContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TrustContract *TrustContractSession) Symbol() (string, error) {
	return _TrustContract.Contract.Symbol(&_TrustContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TrustContract *TrustContractCallerSession) Symbol() (string, error) {
	return _TrustContract.Contract.Symbol(&_TrustContract.CallOpts)
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_TrustContract *TrustContractCaller) ToEthereumSignedMessage(opts *bind.CallOpts, _msg string) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "toEthereumSignedMessage", _msg)
	return *ret0, err
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_TrustContract *TrustContractSession) ToEthereumSignedMessage(_msg string) ([32]byte, error) {
	return _TrustContract.Contract.ToEthereumSignedMessage(&_TrustContract.CallOpts, _msg)
}

// ToEthereumSignedMessage is a free data retrieval call binding the contract method 0xdae21454.
//
// Solidity: function toEthereumSignedMessage(_msg string) constant returns(bytes32)
func (_TrustContract *TrustContractCallerSession) ToEthereumSignedMessage(_msg string) ([32]byte, error) {
	return _TrustContract.Contract.ToEthereumSignedMessage(&_TrustContract.CallOpts, _msg)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TrustContract *TrustContractCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TrustContract *TrustContractSession) TotalSupply() (*big.Int, error) {
	return _TrustContract.Contract.TotalSupply(&_TrustContract.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TrustContract *TrustContractCallerSession) TotalSupply() (*big.Int, error) {
	return _TrustContract.Contract.TotalSupply(&_TrustContract.CallOpts)
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_TrustContract *TrustContractCaller) TxIdLast(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "txIdLast")
	return *ret0, err
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_TrustContract *TrustContractSession) TxIdLast() (uint64, error) {
	return _TrustContract.Contract.TxIdLast(&_TrustContract.CallOpts)
}

// TxIdLast is a free data retrieval call binding the contract method 0xdd098d0b.
//
// Solidity: function txIdLast() constant returns(uint64)
func (_TrustContract *TrustContractCallerSession) TxIdLast() (uint64, error) {
	return _TrustContract.Contract.TxIdLast(&_TrustContract.CallOpts)
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_TrustContract *TrustContractCaller) UintToBytes32(opts *bind.CallOpts, _uint *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "uintToBytes32", _uint)
	return *ret0, err
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_TrustContract *TrustContractSession) UintToBytes32(_uint *big.Int) ([]byte, error) {
	return _TrustContract.Contract.UintToBytes32(&_TrustContract.CallOpts, _uint)
}

// UintToBytes32 is a free data retrieval call binding the contract method 0x886d3db9.
//
// Solidity: function uintToBytes32(_uint uint256) constant returns(b bytes)
func (_TrustContract *TrustContractCallerSession) UintToBytes32(_uint *big.Int) ([]byte, error) {
	return _TrustContract.Contract.UintToBytes32(&_TrustContract.CallOpts, _uint)
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_TrustContract *TrustContractCaller) UintToString(opts *bind.CallOpts, _uint *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "uintToString", _uint)
	return *ret0, err
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_TrustContract *TrustContractSession) UintToString(_uint *big.Int) (string, error) {
	return _TrustContract.Contract.UintToString(&_TrustContract.CallOpts, _uint)
}

// UintToString is a free data retrieval call binding the contract method 0xe9395679.
//
// Solidity: function uintToString(_uint uint256) constant returns(str string)
func (_TrustContract *TrustContractCallerSession) UintToString(_uint *big.Int) (string, error) {
	return _TrustContract.Contract.UintToString(&_TrustContract.CallOpts, _uint)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_TrustContract *TrustContractSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.Approve(&_TrustContract.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(_spender address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.Approve(&_TrustContract.TransactOpts, _spender, _value)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_TrustContract *TrustContractTransactor) AssignInitialSigners(opts *bind.TransactOpts, initialSigners []common.Address) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "assignInitialSigners", initialSigners)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_TrustContract *TrustContractSession) AssignInitialSigners(initialSigners []common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.AssignInitialSigners(&_TrustContract.TransactOpts, initialSigners)
}

// AssignInitialSigners is a paid mutator transaction binding the contract method 0x29d8e43d.
//
// Solidity: function assignInitialSigners(initialSigners address[]) returns()
func (_TrustContract *TrustContractTransactorSession) AssignInitialSigners(initialSigners []common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.AssignInitialSigners(&_TrustContract.TransactOpts, initialSigners)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_TrustContract *TrustContractTransactor) DecreaseApproval(opts *bind.TransactOpts, _spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "decreaseApproval", _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_TrustContract *TrustContractSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.DecreaseApproval(&_TrustContract.TransactOpts, _spender, _subtractedValue)
}

// DecreaseApproval is a paid mutator transaction binding the contract method 0x66188463.
//
// Solidity: function decreaseApproval(_spender address, _subtractedValue uint256) returns(bool)
func (_TrustContract *TrustContractTransactorSession) DecreaseApproval(_spender common.Address, _subtractedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.DecreaseApproval(&_TrustContract.TransactOpts, _spender, _subtractedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_TrustContract *TrustContractTransactor) IncreaseApproval(opts *bind.TransactOpts, _spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "increaseApproval", _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_TrustContract *TrustContractSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.IncreaseApproval(&_TrustContract.TransactOpts, _spender, _addedValue)
}

// IncreaseApproval is a paid mutator transaction binding the contract method 0xd73dd623.
//
// Solidity: function increaseApproval(_spender address, _addedValue uint256) returns(bool)
func (_TrustContract *TrustContractTransactorSession) IncreaseApproval(_spender common.Address, _addedValue *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.IncreaseApproval(&_TrustContract.TransactOpts, _spender, _addedValue)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_TrustContract *TrustContractTransactor) PaymentTx(opts *bind.TransactOpts, txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "paymentTx", txId, erc20Addr, toAddr, amount, v, r, s)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_TrustContract *TrustContractSession) PaymentTx(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _TrustContract.Contract.PaymentTx(&_TrustContract.TransactOpts, txId, erc20Addr, toAddr, amount, v, r, s)
}

// PaymentTx is a paid mutator transaction binding the contract method 0x497483d1.
//
// Solidity: function paymentTx(txId uint64, erc20Addr address, toAddr address, amount uint256, v uint8[], r bytes32[], s bytes32[]) returns()
func (_TrustContract *TrustContractTransactorSession) PaymentTx(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int, v []uint8, r [][32]byte, s [][32]byte) (*types.Transaction, error) {
	return _TrustContract.Contract.PaymentTx(&_TrustContract.TransactOpts, txId, erc20Addr, toAddr, amount, v, r, s)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustContract *TrustContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustContract *TrustContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _TrustContract.Contract.RenounceOwnership(&_TrustContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrustContract *TrustContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TrustContract.Contract.RenounceOwnership(&_TrustContract.TransactOpts)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_TrustContract *TrustContractTransactor) SetCompleted(opts *bind.TransactOpts, completed *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "setCompleted", completed)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_TrustContract *TrustContractSession) SetCompleted(completed *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.SetCompleted(&_TrustContract.TransactOpts, completed)
}

// SetCompleted is a paid mutator transaction binding the contract method 0xfdacd576.
//
// Solidity: function setCompleted(completed uint256) returns()
func (_TrustContract *TrustContractTransactorSession) SetCompleted(completed *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.SetCompleted(&_TrustContract.TransactOpts, completed)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.Transfer(&_TrustContract.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.Transfer(&_TrustContract.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.TransferFrom(&_TrustContract.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(bool)
func (_TrustContract *TrustContractTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TrustContract.Contract.TransferFrom(&_TrustContract.TransactOpts, _from, _to, _value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_TrustContract *TrustContractTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_TrustContract *TrustContractSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.TransferOwnership(&_TrustContract.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_TrustContract *TrustContractTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.TransferOwnership(&_TrustContract.TransactOpts, _newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_TrustContract *TrustContractTransactor) Upgrade(opts *bind.TransactOpts, new_address common.Address) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "upgrade", new_address)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_TrustContract *TrustContractSession) Upgrade(new_address common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.Upgrade(&_TrustContract.TransactOpts, new_address)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(new_address address) returns()
func (_TrustContract *TrustContractTransactorSession) Upgrade(new_address common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.Upgrade(&_TrustContract.TransactOpts, new_address)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(signer address) returns()
func (_TrustContract *TrustContractTransactor) VoteAddSigner(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "voteAddSigner", signer)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(signer address) returns()
func (_TrustContract *TrustContractSession) VoteAddSigner(signer common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.VoteAddSigner(&_TrustContract.TransactOpts, signer)
}

// VoteAddSigner is a paid mutator transaction binding the contract method 0x6a5b3347.
//
// Solidity: function voteAddSigner(signer address) returns()
func (_TrustContract *TrustContractTransactorSession) VoteAddSigner(signer common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.VoteAddSigner(&_TrustContract.TransactOpts, signer)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(signer address) returns()
func (_TrustContract *TrustContractTransactor) VoteRemoveSigner(opts *bind.TransactOpts, signer common.Address) (*types.Transaction, error) {
	return _TrustContract.contract.Transact(opts, "voteRemoveSigner", signer)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(signer address) returns()
func (_TrustContract *TrustContractSession) VoteRemoveSigner(signer common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.VoteRemoveSigner(&_TrustContract.TransactOpts, signer)
}

// VoteRemoveSigner is a paid mutator transaction binding the contract method 0x3296fedc.
//
// Solidity: function voteRemoveSigner(signer address) returns()
func (_TrustContract *TrustContractTransactorSession) VoteRemoveSigner(signer common.Address) (*types.Transaction, error) {
	return _TrustContract.Contract.VoteRemoveSigner(&_TrustContract.TransactOpts, signer)
}

// TrustContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TrustContract contract.
type TrustContractApprovalIterator struct {
	Event *TrustContractApproval // Event containing the contract specifics and raw log

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
func (it *TrustContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractApproval)
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
		it.Event = new(TrustContractApproval)
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
func (it *TrustContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractApproval represents a Approval event raised by the TrustContract contract.
type TrustContractApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_TrustContract *TrustContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TrustContractApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TrustContractApprovalIterator{contract: _TrustContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_TrustContract *TrustContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TrustContractApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractApproval)
				if err := _TrustContract.contract.UnpackLog(event, "Approval", log); err != nil {
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

// TrustContractFundIterator is returned from FilterFund and is used to iterate over the raw logs and unpacked data for Fund events raised by the TrustContract contract.
type TrustContractFundIterator struct {
	Event *TrustContractFund // Event containing the contract specifics and raw log

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
func (it *TrustContractFundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractFund)
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
		it.Event = new(TrustContractFund)
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
func (it *TrustContractFundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractFundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractFund represents a Fund event raised by the TrustContract contract.
type TrustContractFund struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFund is a free log retrieval operation binding the contract event 0xda8220a878ff7a89474ccffdaa31ea1ed1ffbb0207d5051afccc4fbaf81f9bcd.
//
// Solidity: e Fund(_from indexed address, _value uint256)
func (_TrustContract *TrustContractFilterer) FilterFund(opts *bind.FilterOpts, _from []common.Address) (*TrustContractFundIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "Fund", _fromRule)
	if err != nil {
		return nil, err
	}
	return &TrustContractFundIterator{contract: _TrustContract.contract, event: "Fund", logs: logs, sub: sub}, nil
}

// WatchFund is a free log subscription operation binding the contract event 0xda8220a878ff7a89474ccffdaa31ea1ed1ffbb0207d5051afccc4fbaf81f9bcd.
//
// Solidity: e Fund(_from indexed address, _value uint256)
func (_TrustContract *TrustContractFilterer) WatchFund(opts *bind.WatchOpts, sink chan<- *TrustContractFund, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "Fund", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractFund)
				if err := _TrustContract.contract.UnpackLog(event, "Fund", log); err != nil {
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

// TrustContractOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the TrustContract contract.
type TrustContractOwnershipRenouncedIterator struct {
	Event *TrustContractOwnershipRenounced // Event containing the contract specifics and raw log

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
func (it *TrustContractOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractOwnershipRenounced)
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
		it.Event = new(TrustContractOwnershipRenounced)
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
func (it *TrustContractOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractOwnershipRenounced represents a OwnershipRenounced event raised by the TrustContract contract.
type TrustContractOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_TrustContract *TrustContractFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*TrustContractOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TrustContractOwnershipRenouncedIterator{contract: _TrustContract.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_TrustContract *TrustContractFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *TrustContractOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractOwnershipRenounced)
				if err := _TrustContract.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
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

// TrustContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TrustContract contract.
type TrustContractOwnershipTransferredIterator struct {
	Event *TrustContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TrustContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractOwnershipTransferred)
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
		it.Event = new(TrustContractOwnershipTransferred)
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
func (it *TrustContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractOwnershipTransferred represents a OwnershipTransferred event raised by the TrustContract contract.
type TrustContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_TrustContract *TrustContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TrustContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TrustContractOwnershipTransferredIterator{contract: _TrustContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_TrustContract *TrustContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TrustContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractOwnershipTransferred)
				if err := _TrustContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// TrustContractTransactionResultIterator is returned from FilterTransactionResult and is used to iterate over the raw logs and unpacked data for TransactionResult events raised by the TrustContract contract.
type TrustContractTransactionResultIterator struct {
	Event *TrustContractTransactionResult // Event containing the contract specifics and raw log

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
func (it *TrustContractTransactionResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractTransactionResult)
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
		it.Event = new(TrustContractTransactionResult)
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
func (it *TrustContractTransactionResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractTransactionResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractTransactionResult represents a TransactionResult event raised by the TrustContract contract.
type TrustContractTransactionResult struct {
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
func (_TrustContract *TrustContractFilterer) FilterTransactionResult(opts *bind.FilterOpts) (*TrustContractTransactionResultIterator, error) {

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "TransactionResult")
	if err != nil {
		return nil, err
	}
	return &TrustContractTransactionResultIterator{contract: _TrustContract.contract, event: "TransactionResult", logs: logs, sub: sub}, nil
}

// WatchTransactionResult is a free log subscription operation binding the contract event 0x060aac1656ca90ad85bf3d9126ac8709f68d4403c923998b24ea513a9e1ba3a6.
//
// Solidity: e TransactionResult(success bool, txId uint64, erc20Addr address, toAddr address, amount uint256, verified bool[])
func (_TrustContract *TrustContractFilterer) WatchTransactionResult(opts *bind.WatchOpts, sink chan<- *TrustContractTransactionResult) (event.Subscription, error) {

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "TransactionResult")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractTransactionResult)
				if err := _TrustContract.contract.UnpackLog(event, "TransactionResult", log); err != nil {
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

// TrustContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TrustContract contract.
type TrustContractTransferIterator struct {
	Event *TrustContractTransfer // Event containing the contract specifics and raw log

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
func (it *TrustContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrustContractTransfer)
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
		it.Event = new(TrustContractTransfer)
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
func (it *TrustContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrustContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrustContractTransfer represents a Transfer event raised by the TrustContract contract.
type TrustContractTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_TrustContract *TrustContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TrustContractTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TrustContractTransferIterator{contract: _TrustContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_TrustContract *TrustContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TrustContractTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TrustContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrustContractTransfer)
				if err := _TrustContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
