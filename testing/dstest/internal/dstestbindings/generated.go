// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dstestbindings

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = interfaces.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DSTestMetaData contains all meta data concerning the DSTest contract.
var DSTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"log_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"log_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"log_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"log_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"val\",\"type\":\"address\"}],\"name\":\"log_named_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"val\",\"type\":\"bytes\"}],\"name\":\"log_named_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"val\",\"type\":\"bytes32\"}],\"name\":\"log_named_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"}],\"name\":\"log_named_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"val\",\"type\":\"string\"}],\"name\":\"log_named_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"log_named_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"log_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"logs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"IS_TEST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"failed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DSTestABI is the input ABI used to generate the binding from.
// Deprecated: Use DSTestMetaData.ABI instead.
var DSTestABI = DSTestMetaData.ABI

// DSTest is an auto generated Go binding around an Ethereum contract.
type DSTest struct {
	DSTestCaller     // Read-only binding to the contract
	DSTestTransactor // Write-only binding to the contract
	DSTestFilterer   // Log filterer for contract events
}

// DSTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type DSTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DSTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DSTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DSTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DSTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DSTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DSTestSession struct {
	Contract     *DSTest           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DSTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DSTestCallerSession struct {
	Contract *DSTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DSTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DSTestTransactorSession struct {
	Contract     *DSTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DSTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type DSTestRaw struct {
	Contract *DSTest // Generic contract binding to access the raw methods on
}

// DSTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DSTestCallerRaw struct {
	Contract *DSTestCaller // Generic read-only contract binding to access the raw methods on
}

// DSTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DSTestTransactorRaw struct {
	Contract *DSTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDSTest creates a new instance of DSTest, bound to a specific deployed contract.
func NewDSTest(address common.Address, backend bind.ContractBackend) (*DSTest, error) {
	contract, err := bindDSTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DSTest{DSTestCaller: DSTestCaller{contract: contract}, DSTestTransactor: DSTestTransactor{contract: contract}, DSTestFilterer: DSTestFilterer{contract: contract}}, nil
}

// NewDSTestCaller creates a new read-only instance of DSTest, bound to a specific deployed contract.
func NewDSTestCaller(address common.Address, caller bind.ContractCaller) (*DSTestCaller, error) {
	contract, err := bindDSTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DSTestCaller{contract: contract}, nil
}

// NewDSTestTransactor creates a new write-only instance of DSTest, bound to a specific deployed contract.
func NewDSTestTransactor(address common.Address, transactor bind.ContractTransactor) (*DSTestTransactor, error) {
	contract, err := bindDSTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DSTestTransactor{contract: contract}, nil
}

// NewDSTestFilterer creates a new log filterer instance of DSTest, bound to a specific deployed contract.
func NewDSTestFilterer(address common.Address, filterer bind.ContractFilterer) (*DSTestFilterer, error) {
	contract, err := bindDSTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DSTestFilterer{contract: contract}, nil
}

// bindDSTest binds a generic wrapper to an already deployed contract.
func bindDSTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DSTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DSTest *DSTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DSTest.Contract.DSTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DSTest *DSTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DSTest.Contract.DSTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DSTest *DSTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DSTest.Contract.DSTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DSTest *DSTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DSTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DSTest *DSTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DSTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DSTest *DSTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DSTest.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DSTest *DSTestCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DSTest.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DSTest *DSTestSession) ISTEST() (bool, error) {
	return _DSTest.Contract.ISTEST(&_DSTest.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DSTest *DSTestCallerSession) ISTEST() (bool, error) {
	return _DSTest.Contract.ISTEST(&_DSTest.CallOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DSTest *DSTestTransactor) Failed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DSTest.contract.Transact(opts, "failed")
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DSTest *DSTestSession) Failed() (*types.Transaction, error) {
	return _DSTest.Contract.Failed(&_DSTest.TransactOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DSTest *DSTestTransactorSession) Failed() (*types.Transaction, error) {
	return _DSTest.Contract.Failed(&_DSTest.TransactOpts)
}

// DSTestLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the DSTest contract.
type DSTestLogIterator struct {
	Event *DSTestLog // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLog)
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
		it.Event = new(DSTestLog)
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
func (it *DSTestLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLog represents a Log event raised by the DSTest contract.
type DSTestLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DSTest *DSTestFilterer) FilterLog(opts *bind.FilterOpts) (*DSTestLogIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &DSTestLogIterator{contract: _DSTest.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DSTest *DSTestFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *DSTestLog) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLog)
				if err := _DSTest.contract.UnpackLog(event, "log", log); err != nil {
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

// ParseLog is a log parse operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DSTest *DSTestFilterer) ParseLog(log types.Log) (*DSTestLog, error) {
	event := new(DSTestLog)
	if err := _DSTest.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the DSTest contract.
type DSTestLogAddressIterator struct {
	Event *DSTestLogAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogAddress)
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
		it.Event = new(DSTestLogAddress)
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
func (it *DSTestLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogAddress represents a LogAddress event raised by the DSTest contract.
type DSTestLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DSTest *DSTestFilterer) FilterLogAddress(opts *bind.FilterOpts) (*DSTestLogAddressIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &DSTestLogAddressIterator{contract: _DSTest.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DSTest *DSTestFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *DSTestLogAddress) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogAddress)
				if err := _DSTest.contract.UnpackLog(event, "log_address", log); err != nil {
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

// ParseLogAddress is a log parse operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DSTest *DSTestFilterer) ParseLogAddress(log types.Log) (*DSTestLogAddress, error) {
	event := new(DSTestLogAddress)
	if err := _DSTest.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the DSTest contract.
type DSTestLogBytesIterator struct {
	Event *DSTestLogBytes // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogBytes)
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
		it.Event = new(DSTestLogBytes)
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
func (it *DSTestLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogBytes represents a LogBytes event raised by the DSTest contract.
type DSTestLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DSTest *DSTestFilterer) FilterLogBytes(opts *bind.FilterOpts) (*DSTestLogBytesIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &DSTestLogBytesIterator{contract: _DSTest.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DSTest *DSTestFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *DSTestLogBytes) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogBytes)
				if err := _DSTest.contract.UnpackLog(event, "log_bytes", log); err != nil {
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

// ParseLogBytes is a log parse operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DSTest *DSTestFilterer) ParseLogBytes(log types.Log) (*DSTestLogBytes, error) {
	event := new(DSTestLogBytes)
	if err := _DSTest.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the DSTest contract.
type DSTestLogBytes32Iterator struct {
	Event *DSTestLogBytes32 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogBytes32)
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
		it.Event = new(DSTestLogBytes32)
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
func (it *DSTestLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogBytes32 represents a LogBytes32 event raised by the DSTest contract.
type DSTestLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DSTest *DSTestFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*DSTestLogBytes32Iterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &DSTestLogBytes32Iterator{contract: _DSTest.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DSTest *DSTestFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *DSTestLogBytes32) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogBytes32)
				if err := _DSTest.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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

// ParseLogBytes32 is a log parse operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DSTest *DSTestFilterer) ParseLogBytes32(log types.Log) (*DSTestLogBytes32, error) {
	event := new(DSTestLogBytes32)
	if err := _DSTest.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the DSTest contract.
type DSTestLogIntIterator struct {
	Event *DSTestLogInt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogInt)
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
		it.Event = new(DSTestLogInt)
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
func (it *DSTestLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogInt represents a LogInt event raised by the DSTest contract.
type DSTestLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DSTest *DSTestFilterer) FilterLogInt(opts *bind.FilterOpts) (*DSTestLogIntIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &DSTestLogIntIterator{contract: _DSTest.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DSTest *DSTestFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *DSTestLogInt) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogInt)
				if err := _DSTest.contract.UnpackLog(event, "log_int", log); err != nil {
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

// ParseLogInt is a log parse operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DSTest *DSTestFilterer) ParseLogInt(log types.Log) (*DSTestLogInt, error) {
	event := new(DSTestLogInt)
	if err := _DSTest.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the DSTest contract.
type DSTestLogNamedAddressIterator struct {
	Event *DSTestLogNamedAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedAddress)
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
		it.Event = new(DSTestLogNamedAddress)
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
func (it *DSTestLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedAddress represents a LogNamedAddress event raised by the DSTest contract.
type DSTestLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DSTest *DSTestFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*DSTestLogNamedAddressIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedAddressIterator{contract: _DSTest.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DSTest *DSTestFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedAddress)
				if err := _DSTest.contract.UnpackLog(event, "log_named_address", log); err != nil {
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

// ParseLogNamedAddress is a log parse operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DSTest *DSTestFilterer) ParseLogNamedAddress(log types.Log) (*DSTestLogNamedAddress, error) {
	event := new(DSTestLogNamedAddress)
	if err := _DSTest.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the DSTest contract.
type DSTestLogNamedBytesIterator struct {
	Event *DSTestLogNamedBytes // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedBytes)
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
		it.Event = new(DSTestLogNamedBytes)
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
func (it *DSTestLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedBytes represents a LogNamedBytes event raised by the DSTest contract.
type DSTestLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DSTest *DSTestFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*DSTestLogNamedBytesIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedBytesIterator{contract: _DSTest.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DSTest *DSTestFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedBytes)
				if err := _DSTest.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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

// ParseLogNamedBytes is a log parse operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DSTest *DSTestFilterer) ParseLogNamedBytes(log types.Log) (*DSTestLogNamedBytes, error) {
	event := new(DSTestLogNamedBytes)
	if err := _DSTest.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the DSTest contract.
type DSTestLogNamedBytes32Iterator struct {
	Event *DSTestLogNamedBytes32 // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedBytes32)
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
		it.Event = new(DSTestLogNamedBytes32)
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
func (it *DSTestLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedBytes32 represents a LogNamedBytes32 event raised by the DSTest contract.
type DSTestLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DSTest *DSTestFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*DSTestLogNamedBytes32Iterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedBytes32Iterator{contract: _DSTest.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DSTest *DSTestFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedBytes32)
				if err := _DSTest.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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

// ParseLogNamedBytes32 is a log parse operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DSTest *DSTestFilterer) ParseLogNamedBytes32(log types.Log) (*DSTestLogNamedBytes32, error) {
	event := new(DSTestLogNamedBytes32)
	if err := _DSTest.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the DSTest contract.
type DSTestLogNamedDecimalIntIterator struct {
	Event *DSTestLogNamedDecimalInt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedDecimalInt)
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
		it.Event = new(DSTestLogNamedDecimalInt)
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
func (it *DSTestLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the DSTest contract.
type DSTestLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*DSTestLogNamedDecimalIntIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedDecimalIntIterator{contract: _DSTest.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedDecimalInt)
				if err := _DSTest.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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

// ParseLogNamedDecimalInt is a log parse operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) ParseLogNamedDecimalInt(log types.Log) (*DSTestLogNamedDecimalInt, error) {
	event := new(DSTestLogNamedDecimalInt)
	if err := _DSTest.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the DSTest contract.
type DSTestLogNamedDecimalUintIterator struct {
	Event *DSTestLogNamedDecimalUint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedDecimalUint)
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
		it.Event = new(DSTestLogNamedDecimalUint)
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
func (it *DSTestLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the DSTest contract.
type DSTestLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*DSTestLogNamedDecimalUintIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedDecimalUintIterator{contract: _DSTest.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedDecimalUint)
				if err := _DSTest.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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

// ParseLogNamedDecimalUint is a log parse operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DSTest *DSTestFilterer) ParseLogNamedDecimalUint(log types.Log) (*DSTestLogNamedDecimalUint, error) {
	event := new(DSTestLogNamedDecimalUint)
	if err := _DSTest.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the DSTest contract.
type DSTestLogNamedIntIterator struct {
	Event *DSTestLogNamedInt // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedInt)
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
		it.Event = new(DSTestLogNamedInt)
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
func (it *DSTestLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedInt represents a LogNamedInt event raised by the DSTest contract.
type DSTestLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DSTest *DSTestFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*DSTestLogNamedIntIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedIntIterator{contract: _DSTest.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DSTest *DSTestFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedInt)
				if err := _DSTest.contract.UnpackLog(event, "log_named_int", log); err != nil {
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

// ParseLogNamedInt is a log parse operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DSTest *DSTestFilterer) ParseLogNamedInt(log types.Log) (*DSTestLogNamedInt, error) {
	event := new(DSTestLogNamedInt)
	if err := _DSTest.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the DSTest contract.
type DSTestLogNamedStringIterator struct {
	Event *DSTestLogNamedString // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedString)
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
		it.Event = new(DSTestLogNamedString)
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
func (it *DSTestLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedString represents a LogNamedString event raised by the DSTest contract.
type DSTestLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DSTest *DSTestFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*DSTestLogNamedStringIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedStringIterator{contract: _DSTest.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DSTest *DSTestFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedString) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedString)
				if err := _DSTest.contract.UnpackLog(event, "log_named_string", log); err != nil {
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

// ParseLogNamedString is a log parse operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DSTest *DSTestFilterer) ParseLogNamedString(log types.Log) (*DSTestLogNamedString, error) {
	event := new(DSTestLogNamedString)
	if err := _DSTest.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the DSTest contract.
type DSTestLogNamedUintIterator struct {
	Event *DSTestLogNamedUint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogNamedUint)
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
		it.Event = new(DSTestLogNamedUint)
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
func (it *DSTestLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogNamedUint represents a LogNamedUint event raised by the DSTest contract.
type DSTestLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DSTest *DSTestFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*DSTestLogNamedUintIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &DSTestLogNamedUintIterator{contract: _DSTest.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DSTest *DSTestFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *DSTestLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogNamedUint)
				if err := _DSTest.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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

// ParseLogNamedUint is a log parse operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DSTest *DSTestFilterer) ParseLogNamedUint(log types.Log) (*DSTestLogNamedUint, error) {
	event := new(DSTestLogNamedUint)
	if err := _DSTest.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the DSTest contract.
type DSTestLogStringIterator struct {
	Event *DSTestLogString // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogString)
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
		it.Event = new(DSTestLogString)
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
func (it *DSTestLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogString represents a LogString event raised by the DSTest contract.
type DSTestLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DSTest *DSTestFilterer) FilterLogString(opts *bind.FilterOpts) (*DSTestLogStringIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &DSTestLogStringIterator{contract: _DSTest.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DSTest *DSTestFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *DSTestLogString) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogString)
				if err := _DSTest.contract.UnpackLog(event, "log_string", log); err != nil {
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

// ParseLogString is a log parse operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DSTest *DSTestFilterer) ParseLogString(log types.Log) (*DSTestLogString, error) {
	event := new(DSTestLogString)
	if err := _DSTest.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the DSTest contract.
type DSTestLogUintIterator struct {
	Event *DSTestLogUint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogUint)
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
		it.Event = new(DSTestLogUint)
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
func (it *DSTestLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogUint represents a LogUint event raised by the DSTest contract.
type DSTestLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DSTest *DSTestFilterer) FilterLogUint(opts *bind.FilterOpts) (*DSTestLogUintIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &DSTestLogUintIterator{contract: _DSTest.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DSTest *DSTestFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *DSTestLogUint) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogUint)
				if err := _DSTest.contract.UnpackLog(event, "log_uint", log); err != nil {
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

// ParseLogUint is a log parse operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DSTest *DSTestFilterer) ParseLogUint(log types.Log) (*DSTestLogUint, error) {
	event := new(DSTestLogUint)
	if err := _DSTest.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DSTestLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the DSTest contract.
type DSTestLogsIterator struct {
	Event *DSTestLogs // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DSTestLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DSTestLogs)
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
		it.Event = new(DSTestLogs)
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
func (it *DSTestLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DSTestLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DSTestLogs represents a Logs event raised by the DSTest contract.
type DSTestLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DSTest *DSTestFilterer) FilterLogs(opts *bind.FilterOpts) (*DSTestLogsIterator, error) {

	logs, sub, err := _DSTest.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &DSTestLogsIterator{contract: _DSTest.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DSTest *DSTestFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *DSTestLogs) (event.Subscription, error) {

	logs, sub, err := _DSTest.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DSTestLogs)
				if err := _DSTest.contract.UnpackLog(event, "logs", log); err != nil {
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

// ParseLogs is a log parse operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DSTest *DSTestFilterer) ParseLogs(log types.Log) (*DSTestLogs, error) {
	event := new(DSTestLogs)
	if err := _DSTest.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

