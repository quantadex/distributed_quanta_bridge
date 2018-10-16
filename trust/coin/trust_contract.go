// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package coin

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
const TrustContractABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_sig\",\"type\":\"string\"},{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isSignedBy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hashedMsg\",\"type\":\"bytes32\"},{\"name\":\"_v\",\"type\":\"uint8\"},{\"name\":\"_r\",\"type\":\"bytes32\"},{\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"recoverSignerVRS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_hexstr\",\"type\":\"string\"}],\"name\":\"hexstrToBytes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_char\",\"type\":\"string\"}],\"name\":\"parseInt16Char\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToBytes32\",\"outputs\":[{\"name\":\"b\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_msg\",\"type\":\"string\"}],\"name\":\"toEthereumSignedMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_uint\",\"type\":\"uint256\"}],\"name\":\"uintToString\",\"outputs\":[{\"name\":\"str\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_str\",\"type\":\"string\"},{\"name\":\"_startIndex\",\"type\":\"uint256\"},{\"name\":\"_endIndex\",\"type\":\"uint256\"}],\"name\":\"substring\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"last_completed_migration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"name\":\"completed\",\"type\":\"uint256\"}],\"name\":\"setCompleted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"txId\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"erc20Address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"verified\",\"type\":\"bool[]\"},{\"indexed\":false,\"name\":\"v\",\"type\":\"uint8[]\"},{\"indexed\":false,\"name\":\"r\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"s\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"name\":\"debugSigMsgLength\",\"type\":\"uint256\"}],\"name\":\"TransactionResult\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"txId\",\"type\":\"uint64\"},{\"name\":\"erc20Addr\",\"type\":\"address\"},{\"name\":\"toAddr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8[]\"},{\"name\":\"r\",\"type\":\"bytes32[]\"},{\"name\":\"s\",\"type\":\"bytes32[]\"}],\"name\":\"paymentTx\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"voteAddSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"voteRemoveSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"txId\",\"type\":\"uint64\"},{\"name\":\"erc20Addr\",\"type\":\"address\"},{\"name\":\"toAddr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"toQuantaPaymentSignatureMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

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

// ToQuantaPaymentSignatureMessage is a free data retrieval call binding the contract method 0x517839f0.
//
// Solidity: function toQuantaPaymentSignatureMessage(txId uint64, erc20Addr address, toAddr address, amount uint256) constant returns(bytes)
func (_TrustContract *TrustContractCaller) ToQuantaPaymentSignatureMessage(opts *bind.CallOpts, txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _TrustContract.contract.Call(opts, out, "toQuantaPaymentSignatureMessage", txId, erc20Addr, toAddr, amount)
	return *ret0, err
}

// ToQuantaPaymentSignatureMessage is a free data retrieval call binding the contract method 0x517839f0.
//
// Solidity: function toQuantaPaymentSignatureMessage(txId uint64, erc20Addr address, toAddr address, amount uint256) constant returns(bytes)
func (_TrustContract *TrustContractSession) ToQuantaPaymentSignatureMessage(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int) ([]byte, error) {
	return _TrustContract.Contract.ToQuantaPaymentSignatureMessage(&_TrustContract.CallOpts, txId, erc20Addr, toAddr, amount)
}

// ToQuantaPaymentSignatureMessage is a free data retrieval call binding the contract method 0x517839f0.
//
// Solidity: function toQuantaPaymentSignatureMessage(txId uint64, erc20Addr address, toAddr address, amount uint256) constant returns(bytes)
func (_TrustContract *TrustContractCallerSession) ToQuantaPaymentSignatureMessage(txId uint64, erc20Addr common.Address, toAddr common.Address, amount *big.Int) ([]byte, error) {
	return _TrustContract.Contract.ToQuantaPaymentSignatureMessage(&_TrustContract.CallOpts, txId, erc20Addr, toAddr, amount)
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
	Success           bool
	TxId              uint64
	Erc20Address      common.Address
	To                common.Address
	Amount            *big.Int
	Verified          []bool
	V                 []uint8
	R                 [][32]byte
	S                 [][32]byte
	DebugSigMsgLength *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTransactionResult is a free log retrieval operation binding the contract event 0xa43d6dd134bc846b99efc88f3b1bcc7ce753b6da962afb8c5436e9d1d81fcbb1.
//
// Solidity: e TransactionResult(success bool, txId uint64, erc20Address address, to address, amount uint256, verified bool[], v uint8[], r bytes32[], s bytes32[], debugSigMsgLength uint256)
func (_TrustContract *TrustContractFilterer) FilterTransactionResult(opts *bind.FilterOpts) (*TrustContractTransactionResultIterator, error) {

	logs, sub, err := _TrustContract.contract.FilterLogs(opts, "TransactionResult")
	if err != nil {
		return nil, err
	}
	return &TrustContractTransactionResultIterator{contract: _TrustContract.contract, event: "TransactionResult", logs: logs, sub: sub}, nil
}

// WatchTransactionResult is a free log subscription operation binding the contract event 0xa43d6dd134bc846b99efc88f3b1bcc7ce753b6da962afb8c5436e9d1d81fcbb1.
//
// Solidity: e TransactionResult(success bool, txId uint64, erc20Address address, to address, amount uint256, verified bool[], v uint8[], r bytes32[], s bytes32[], debugSigMsgLength uint256)
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
