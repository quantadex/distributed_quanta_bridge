// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Forwarder

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ForwarderABI is the input ABI used to generate the binding from.
const ForwarderABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"quantaAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"destinationAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"trust\",\"type\":\"address\"},{\"name\":\"quantaAddr\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogForwarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogFlushed\",\"type\":\"event\"}]"

// ForwarderBin is the compiled bytecode used for deploying new contracts.
const ForwarderBin  = `608060405234801561001057600080fd5b5060405161041e38038061041e83398101604052805160208083015160008054600160a060020a031916600160a060020a0385161790559092018051919290916100609160019190840190610068565b505050610103565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a957805160ff19168380011785556100d6565b828001600101855582156100d6579182015b828111156100d65782518255916020019190600101906100bb565b506100e29291506100e6565b5090565b61010091905b808211156100e257600081556001016100ec565b90565b61030c806101126000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100d65780636b9f96ea14610160578063ca32546914610177575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100d3573d6000803e3d6000fd5b50005b3480156100e257600080fd5b506100eb6101b5565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012557818101518382015260200161010d565b50505050905090810190601f1680156101525780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016c57600080fd5b50610175610242565b005b34801561018357600080fd5b5061018c6102c4565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561023a5780601f1061020f5761010080835404028352916020019161023a565b820191906000526020600020905b81548152906001019060200180831161021d57829003601f168201915b505050505081565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f193505050501580156102c1573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582043822668310b3d2e90aa63e4d6df1b84d5beb5b293b898420b19f04f9269d9770029`
const ForwarderBinV2= `608060405234801561001057600080fd5b506040516104d13803806104d183398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b61030c806101c56000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100d65780636b9f96ea14610160578063ca32546914610177575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100d3573d6000803e3d6000fd5b50005b3480156100e257600080fd5b506100eb6101b5565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012557818101518382015260200161010d565b50505050905090810190601f1680156101525780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016c57600080fd5b50610175610242565b005b34801561018357600080fd5b5061018c6102c4565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561023a5780601f1061020f5761010080835404028352916020019161023a565b820191906000526020600020905b81548152906001019060200180831161021d57829003601f168201915b505050505081565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f193505050501580156102c1573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820f03161a8010bca92450c47c8afca14692b54f0d917a1bcbb6d36a9147abf885c0029`

// DeployForwarder deploys a new Ethereum contract, binding an instance of Forwarder to it.
func DeployForwarder(auth *bind.TransactOpts, backend bind.ContractBackend, trust common.Address, quantaAddr string) (common.Address, *types.Transaction, *Forwarder, error) {
	parsed, err := abi.JSON(strings.NewReader(ForwarderABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ForwarderBin), backend, trust, quantaAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Forwarder{ForwarderCaller: ForwarderCaller{contract: contract}, ForwarderTransactor: ForwarderTransactor{contract: contract}}, nil
}

// Forwarder is an auto generated Go binding around an Ethereum contract.
type Forwarder struct {
	ForwarderCaller     // Read-only binding to the contract
	ForwarderTransactor // Write-only binding to the contract
}

// ForwarderCaller is an auto generated read-only Go binding around an Ethereum contract.
type ForwarderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ForwarderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ForwarderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ForwarderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ForwarderSession struct {
	Contract     *Forwarder        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ForwarderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ForwarderCallerSession struct {
	Contract *ForwarderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ForwarderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ForwarderTransactorSession struct {
	Contract     *ForwarderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ForwarderRaw is an auto generated low-level Go binding around an Ethereum contract.
type ForwarderRaw struct {
	Contract *Forwarder // Generic contract binding to access the raw methods on
}

// ForwarderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ForwarderCallerRaw struct {
	Contract *ForwarderCaller // Generic read-only contract binding to access the raw methods on
}

// ForwarderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ForwarderTransactorRaw struct {
	Contract *ForwarderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewForwarder creates a new instance of Forwarder, bound to a specific deployed contract.
func NewForwarder(address common.Address, backend bind.ContractBackend) (*Forwarder, error) {
	contract, err := bindForwarder(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Forwarder{ForwarderCaller: ForwarderCaller{contract: contract}, ForwarderTransactor: ForwarderTransactor{contract: contract}}, nil
}

// NewForwarderCaller creates a new read-only instance of Forwarder, bound to a specific deployed contract.
func NewForwarderCaller(address common.Address, caller bind.ContractCaller) (*ForwarderCaller, error) {
	contract, err := bindForwarder(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &ForwarderCaller{contract: contract}, nil
}

// NewForwarderTransactor creates a new write-only instance of Forwarder, bound to a specific deployed contract.
func NewForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*ForwarderTransactor, error) {
	contract, err := bindForwarder(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &ForwarderTransactor{contract: contract}, nil
}

// bindForwarder binds a generic wrapper to an already deployed contract.
func bindForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ForwarderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, nil), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Forwarder *ForwarderRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Forwarder.Contract.ForwarderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Forwarder *ForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.Contract.ForwarderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Forwarder *ForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Forwarder.Contract.ForwarderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Forwarder *ForwarderCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Forwarder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Forwarder *ForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Forwarder *ForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Forwarder.Contract.contract.Transact(opts, method, params...)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_Forwarder *ForwarderCaller) DestinationAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Forwarder.contract.Call(opts, out, "destinationAddress")
	return *ret0, err
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_Forwarder *ForwarderSession) DestinationAddress() (common.Address, error) {
	return _Forwarder.Contract.DestinationAddress(&_Forwarder.CallOpts)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_Forwarder *ForwarderCallerSession) DestinationAddress() (common.Address, error) {
	return _Forwarder.Contract.DestinationAddress(&_Forwarder.CallOpts)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_Forwarder *ForwarderCaller) QuantaAddress(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Forwarder.contract.Call(opts, out, "quantaAddress")
	return *ret0, err
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_Forwarder *ForwarderSession) QuantaAddress() (string, error) {
	return _Forwarder.Contract.QuantaAddress(&_Forwarder.CallOpts)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_Forwarder *ForwarderCallerSession) QuantaAddress() (string, error) {
	return _Forwarder.Contract.QuantaAddress(&_Forwarder.CallOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Forwarder *ForwarderTransactor) Flush(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Forwarder.contract.Transact(opts, "flush")
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Forwarder *ForwarderSession) Flush() (*types.Transaction, error) {
	return _Forwarder.Contract.Flush(&_Forwarder.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Forwarder *ForwarderTransactorSession) Flush() (*types.Transaction, error) {
	return _Forwarder.Contract.Flush(&_Forwarder.TransactOpts)
}
