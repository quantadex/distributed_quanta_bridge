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

// ERC20BasicABI is the input ABI used to generate the binding from.
const ERC20BasicABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// ERC20BasicBin is the compiled bytecode used for deploying new contracts.
const ERC20BasicBin = `0x`

// DeployERC20Basic deploys a new Ethereum contract, binding an instance of ERC20Basic to it.
func DeployERC20Basic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Basic, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20BasicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// ERC20Basic is an auto generated Go binding around an Ethereum contract.
type ERC20Basic struct {
	ERC20BasicCaller     // Read-only binding to the contract
	ERC20BasicTransactor // Write-only binding to the contract
	ERC20BasicFilterer   // Log filterer for contract events
}

// ERC20BasicCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BasicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BasicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BasicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BasicSession struct {
	Contract     *ERC20Basic       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BasicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BasicCallerSession struct {
	Contract *ERC20BasicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ERC20BasicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BasicTransactorSession struct {
	Contract     *ERC20BasicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ERC20BasicRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BasicRaw struct {
	Contract *ERC20Basic // Generic contract binding to access the raw methods on
}

// ERC20BasicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BasicCallerRaw struct {
	Contract *ERC20BasicCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BasicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BasicTransactorRaw struct {
	Contract *ERC20BasicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Basic creates a new instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20Basic(address common.Address, backend bind.ContractBackend) (*ERC20Basic, error) {
	contract, err := bindERC20Basic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// NewERC20BasicCaller creates a new read-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicCaller(address common.Address, caller bind.ContractCaller) (*ERC20BasicCaller, error) {
	contract, err := bindERC20Basic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicCaller{contract: contract}, nil
}

// NewERC20BasicTransactor creates a new write-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BasicTransactor, error) {
	contract, err := bindERC20Basic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransactor{contract: contract}, nil
}

// NewERC20BasicFilterer creates a new log filterer instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BasicFilterer, error) {
	contract, err := bindERC20Basic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicFilterer{contract: contract}, nil
}

// bindERC20Basic binds a generic wrapper to an already deployed contract.
func bindERC20Basic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Basic.Contract.ERC20BasicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20Basic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) BalanceOf(opts *bind.CallOpts, _who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Basic.contract.Call(opts, out, "balanceOf", _who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) BalanceOf(_who common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, _who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_who address) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) BalanceOf(_who common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, _who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Basic.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns(bool)
func (_ERC20Basic *ERC20BasicTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, _to, _value)
}

// ERC20BasicTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Basic contract.
type ERC20BasicTransferIterator struct {
	Event *ERC20BasicTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20BasicTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BasicTransfer)
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
		it.Event = new(ERC20BasicTransfer)
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
func (it *ERC20BasicTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BasicTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BasicTransfer represents a Transfer event raised by the ERC20Basic contract.
type ERC20BasicTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20Basic *ERC20BasicFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20BasicTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransferIterator{contract: _ERC20Basic.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_ERC20Basic *ERC20BasicFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20BasicTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BasicTransfer)
				if err := _ERC20Basic.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// QuantaForwarderABI is the input ABI used to generate the binding from.
const QuantaForwarderABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"quantaAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"tokenContractAddress\",\"type\":\"address\"}],\"name\":\"flushTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"destinationAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"trust\",\"type\":\"address\"},{\"name\":\"quantaAddr\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogForwarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogFlushed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"trust\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"quanta\",\"type\":\"string\"}],\"name\":\"LogCreated\",\"type\":\"event\"}]"

// QuantaForwarderBin is the compiled bytecode used for deploying new contracts.
const QuantaForwarderBin = `0x608060405234801561001057600080fd5b5060405161061838038061061883398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b610453806101c56000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100d45780633ef133671461015e5780636b9f96ea14610181578063ca32546914610196575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a260008054604051600160a060020a03909116913480156108fc02929091818181858888f193505050501580156100d1573d6000803e3d6000fd5b50005b3480156100e057600080fd5b506100e96101c7565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012357818101518382015260200161010b565b50505050905090810190601f1680156101505780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016a57600080fd5b5061017f600160a060020a0360043516610254565b005b34801561018d57600080fd5b5061017f6103a3565b3480156101a257600080fd5b506101ab610418565b60408051600160a060020a039092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561024c5780601f106102215761010080835404028352916020019161024c565b820191906000526020600020905b81548152906001019060200180831161022f57829003601f168201915b505050505081565b604080517f70a0823100000000000000000000000000000000000000000000000000000000815230600482018190529151839291600091600160a060020a038516916370a0823191602480830192602092919082900301818787803b1580156102bc57600080fd5b505af11580156102d0573d6000803e3d6000fd5b505050506040513d60208110156102e657600080fd5b505190508015156102f65761039d565b60008054604080517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a0392831660048201526024810185905290519186169263a9059cbb926044808401936020939083900390910190829087803b15801561036657600080fd5b505af115801561037a573d6000803e3d6000fd5b505050506040513d602081101561039057600080fd5b5051151561039d57600080fd5b50505050565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a260008054604051600160a060020a0390911691303180156108fc02929091818181858888f19350505050158015610415573d6000803e3d6000fd5b50565b600054600160a060020a0316815600a165627a7a72305820a08141b05d23dc047d89c549f09510a50f67aeb1611506427042580b59f795410029`

// DeployQuantaForwarder deploys a new Ethereum contract, binding an instance of QuantaForwarder to it.
func DeployQuantaForwarder(auth *bind.TransactOpts, backend bind.ContractBackend, trust common.Address, quantaAddr string) (common.Address, *types.Transaction, *QuantaForwarder, error) {
	parsed, err := abi.JSON(strings.NewReader(QuantaForwarderABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(QuantaForwarderBin), backend, trust, quantaAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &QuantaForwarder{QuantaForwarderCaller: QuantaForwarderCaller{contract: contract}, QuantaForwarderTransactor: QuantaForwarderTransactor{contract: contract}, QuantaForwarderFilterer: QuantaForwarderFilterer{contract: contract}}, nil
}

// QuantaForwarder is an auto generated Go binding around an Ethereum contract.
type QuantaForwarder struct {
	QuantaForwarderCaller     // Read-only binding to the contract
	QuantaForwarderTransactor // Write-only binding to the contract
	QuantaForwarderFilterer   // Log filterer for contract events
}

// QuantaForwarderCaller is an auto generated read-only Go binding around an Ethereum contract.
type QuantaForwarderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuantaForwarderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QuantaForwarderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuantaForwarderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuantaForwarderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuantaForwarderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuantaForwarderSession struct {
	Contract     *QuantaForwarder  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuantaForwarderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuantaForwarderCallerSession struct {
	Contract *QuantaForwarderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// QuantaForwarderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuantaForwarderTransactorSession struct {
	Contract     *QuantaForwarderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// QuantaForwarderRaw is an auto generated low-level Go binding around an Ethereum contract.
type QuantaForwarderRaw struct {
	Contract *QuantaForwarder // Generic contract binding to access the raw methods on
}

// QuantaForwarderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuantaForwarderCallerRaw struct {
	Contract *QuantaForwarderCaller // Generic read-only contract binding to access the raw methods on
}

// QuantaForwarderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuantaForwarderTransactorRaw struct {
	Contract *QuantaForwarderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQuantaForwarder creates a new instance of QuantaForwarder, bound to a specific deployed contract.
func NewQuantaForwarder(address common.Address, backend bind.ContractBackend) (*QuantaForwarder, error) {
	contract, err := bindQuantaForwarder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarder{QuantaForwarderCaller: QuantaForwarderCaller{contract: contract}, QuantaForwarderTransactor: QuantaForwarderTransactor{contract: contract}, QuantaForwarderFilterer: QuantaForwarderFilterer{contract: contract}}, nil
}

// NewQuantaForwarderCaller creates a new read-only instance of QuantaForwarder, bound to a specific deployed contract.
func NewQuantaForwarderCaller(address common.Address, caller bind.ContractCaller) (*QuantaForwarderCaller, error) {
	contract, err := bindQuantaForwarder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderCaller{contract: contract}, nil
}

// NewQuantaForwarderTransactor creates a new write-only instance of QuantaForwarder, bound to a specific deployed contract.
func NewQuantaForwarderTransactor(address common.Address, transactor bind.ContractTransactor) (*QuantaForwarderTransactor, error) {
	contract, err := bindQuantaForwarder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderTransactor{contract: contract}, nil
}

// NewQuantaForwarderFilterer creates a new log filterer instance of QuantaForwarder, bound to a specific deployed contract.
func NewQuantaForwarderFilterer(address common.Address, filterer bind.ContractFilterer) (*QuantaForwarderFilterer, error) {
	contract, err := bindQuantaForwarder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderFilterer{contract: contract}, nil
}

// bindQuantaForwarder binds a generic wrapper to an already deployed contract.
func bindQuantaForwarder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(QuantaForwarderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuantaForwarder *QuantaForwarderRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _QuantaForwarder.Contract.QuantaForwarderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuantaForwarder *QuantaForwarderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.QuantaForwarderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuantaForwarder *QuantaForwarderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.QuantaForwarderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuantaForwarder *QuantaForwarderCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _QuantaForwarder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuantaForwarder *QuantaForwarderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuantaForwarder *QuantaForwarderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.contract.Transact(opts, method, params...)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_QuantaForwarder *QuantaForwarderCaller) DestinationAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _QuantaForwarder.contract.Call(opts, out, "destinationAddress")
	return *ret0, err
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_QuantaForwarder *QuantaForwarderSession) DestinationAddress() (common.Address, error) {
	return _QuantaForwarder.Contract.DestinationAddress(&_QuantaForwarder.CallOpts)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() constant returns(address)
func (_QuantaForwarder *QuantaForwarderCallerSession) DestinationAddress() (common.Address, error) {
	return _QuantaForwarder.Contract.DestinationAddress(&_QuantaForwarder.CallOpts)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_QuantaForwarder *QuantaForwarderCaller) QuantaAddress(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _QuantaForwarder.contract.Call(opts, out, "quantaAddress")
	return *ret0, err
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_QuantaForwarder *QuantaForwarderSession) QuantaAddress() (string, error) {
	return _QuantaForwarder.Contract.QuantaAddress(&_QuantaForwarder.CallOpts)
}

// QuantaAddress is a free data retrieval call binding the contract method 0x3c8410a2.
//
// Solidity: function quantaAddress() constant returns(string)
func (_QuantaForwarder *QuantaForwarderCallerSession) QuantaAddress() (string, error) {
	return _QuantaForwarder.Contract.QuantaAddress(&_QuantaForwarder.CallOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_QuantaForwarder *QuantaForwarderTransactor) Flush(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuantaForwarder.contract.Transact(opts, "flush")
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_QuantaForwarder *QuantaForwarderSession) Flush() (*types.Transaction, error) {
	return _QuantaForwarder.Contract.Flush(&_QuantaForwarder.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_QuantaForwarder *QuantaForwarderTransactorSession) Flush() (*types.Transaction, error) {
	return _QuantaForwarder.Contract.Flush(&_QuantaForwarder.TransactOpts)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_QuantaForwarder *QuantaForwarderTransactor) FlushTokens(opts *bind.TransactOpts, tokenContractAddress common.Address) (*types.Transaction, error) {
	return _QuantaForwarder.contract.Transact(opts, "flushTokens", tokenContractAddress)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_QuantaForwarder *QuantaForwarderSession) FlushTokens(tokenContractAddress common.Address) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.FlushTokens(&_QuantaForwarder.TransactOpts, tokenContractAddress)
}

// FlushTokens is a paid mutator transaction binding the contract method 0x3ef13367.
//
// Solidity: function flushTokens(tokenContractAddress address) returns()
func (_QuantaForwarder *QuantaForwarderTransactorSession) FlushTokens(tokenContractAddress common.Address) (*types.Transaction, error) {
	return _QuantaForwarder.Contract.FlushTokens(&_QuantaForwarder.TransactOpts, tokenContractAddress)
}

// QuantaForwarderLogCreatedIterator is returned from FilterLogCreated and is used to iterate over the raw logs and unpacked data for LogCreated events raised by the QuantaForwarder contract.
type QuantaForwarderLogCreatedIterator struct {
	Event *QuantaForwarderLogCreated // Event containing the contract specifics and raw log

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
func (it *QuantaForwarderLogCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuantaForwarderLogCreated)
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
		it.Event = new(QuantaForwarderLogCreated)
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
func (it *QuantaForwarderLogCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuantaForwarderLogCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuantaForwarderLogCreated represents a LogCreated event raised by the QuantaForwarder contract.
type QuantaForwarderLogCreated struct {
	Trust  common.Address
	Quanta string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogCreated is a free log retrieval operation binding the contract event 0xa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a.
//
// Solidity: e LogCreated(trust address, quanta string)
func (_QuantaForwarder *QuantaForwarderFilterer) FilterLogCreated(opts *bind.FilterOpts) (*QuantaForwarderLogCreatedIterator, error) {

	logs, sub, err := _QuantaForwarder.contract.FilterLogs(opts, "LogCreated")
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderLogCreatedIterator{contract: _QuantaForwarder.contract, event: "LogCreated", logs: logs, sub: sub}, nil
}

// WatchLogCreated is a free log subscription operation binding the contract event 0xa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a.
//
// Solidity: e LogCreated(trust address, quanta string)
func (_QuantaForwarder *QuantaForwarderFilterer) WatchLogCreated(opts *bind.WatchOpts, sink chan<- *QuantaForwarderLogCreated) (event.Subscription, error) {

	logs, sub, err := _QuantaForwarder.contract.WatchLogs(opts, "LogCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuantaForwarderLogCreated)
				if err := _QuantaForwarder.contract.UnpackLog(event, "LogCreated", log); err != nil {
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

// QuantaForwarderLogFlushedIterator is returned from FilterLogFlushed and is used to iterate over the raw logs and unpacked data for LogFlushed events raised by the QuantaForwarder contract.
type QuantaForwarderLogFlushedIterator struct {
	Event *QuantaForwarderLogFlushed // Event containing the contract specifics and raw log

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
func (it *QuantaForwarderLogFlushedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuantaForwarderLogFlushed)
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
		it.Event = new(QuantaForwarderLogFlushed)
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
func (it *QuantaForwarderLogFlushedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuantaForwarderLogFlushedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuantaForwarderLogFlushed represents a LogFlushed event raised by the QuantaForwarder contract.
type QuantaForwarderLogFlushed struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogFlushed is a free log retrieval operation binding the contract event 0xa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea.
//
// Solidity: e LogFlushed(sender indexed address, amount uint256)
func (_QuantaForwarder *QuantaForwarderFilterer) FilterLogFlushed(opts *bind.FilterOpts, sender []common.Address) (*QuantaForwarderLogFlushedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _QuantaForwarder.contract.FilterLogs(opts, "LogFlushed", senderRule)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderLogFlushedIterator{contract: _QuantaForwarder.contract, event: "LogFlushed", logs: logs, sub: sub}, nil
}

// WatchLogFlushed is a free log subscription operation binding the contract event 0xa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea.
//
// Solidity: e LogFlushed(sender indexed address, amount uint256)
func (_QuantaForwarder *QuantaForwarderFilterer) WatchLogFlushed(opts *bind.WatchOpts, sink chan<- *QuantaForwarderLogFlushed, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _QuantaForwarder.contract.WatchLogs(opts, "LogFlushed", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuantaForwarderLogFlushed)
				if err := _QuantaForwarder.contract.UnpackLog(event, "LogFlushed", log); err != nil {
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

// QuantaForwarderLogForwardedIterator is returned from FilterLogForwarded and is used to iterate over the raw logs and unpacked data for LogForwarded events raised by the QuantaForwarder contract.
type QuantaForwarderLogForwardedIterator struct {
	Event *QuantaForwarderLogForwarded // Event containing the contract specifics and raw log

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
func (it *QuantaForwarderLogForwardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(QuantaForwarderLogForwarded)
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
		it.Event = new(QuantaForwarderLogForwarded)
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
func (it *QuantaForwarderLogForwardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *QuantaForwarderLogForwardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// QuantaForwarderLogForwarded represents a LogForwarded event raised by the QuantaForwarder contract.
type QuantaForwarderLogForwarded struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogForwarded is a free log retrieval operation binding the contract event 0x5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6.
//
// Solidity: e LogForwarded(sender indexed address, amount uint256)
func (_QuantaForwarder *QuantaForwarderFilterer) FilterLogForwarded(opts *bind.FilterOpts, sender []common.Address) (*QuantaForwarderLogForwardedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _QuantaForwarder.contract.FilterLogs(opts, "LogForwarded", senderRule)
	if err != nil {
		return nil, err
	}
	return &QuantaForwarderLogForwardedIterator{contract: _QuantaForwarder.contract, event: "LogForwarded", logs: logs, sub: sub}, nil
}

// WatchLogForwarded is a free log subscription operation binding the contract event 0x5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6.
//
// Solidity: e LogForwarded(sender indexed address, amount uint256)
func (_QuantaForwarder *QuantaForwarderFilterer) WatchLogForwarded(opts *bind.WatchOpts, sink chan<- *QuantaForwarderLogForwarded, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _QuantaForwarder.contract.WatchLogs(opts, "LogForwarded", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(QuantaForwarderLogForwarded)
				if err := _QuantaForwarder.contract.UnpackLog(event, "LogForwarded", log); err != nil {
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
