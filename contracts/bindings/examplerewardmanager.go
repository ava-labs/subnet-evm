// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ava-labs/libevm"
	"github.com/ava-labs/libevm/accounts/abi"
	"github.com/ava-labs/libevm/accounts/abi/bind"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ExampleRewardManagerMetaData contains all meta data concerning the ExampleRewardManager contract.
var ExampleRewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"allowFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"areFeeRecipientsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRewardAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setRewardAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052730200000000000000000000000000000000000004600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006557600080fd5b5033600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100d95760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100d091906101f3565b60405180910390fd5b6100e8816100ee60201b60201c565b5061020e565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101dd826101b2565b9050919050565b6101ed816101d2565b82525050565b600060208201905061020860008301846101e4565b92915050565b6107f48061021d6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063bc1786281161005b578063bc178628146100db578063e915608b146100e5578063f2fde38b14610103578063f6542b2e1461011f57610088565b80630329099f1461008d5780635e00e67914610097578063715018a6146100b35780638da5cb5b146100bd575b600080fd5b61009561013d565b005b6100b160048036038101906100ac9190610696565b6101c9565b005b6100bb610261565b005b6100c5610275565b6040516100d291906106d2565b60405180910390f35b6100e361029e565b005b6100ed61032a565b6040516100fa91906106d2565b60405180910390f35b61011d60048036038101906101189190610696565b6103c2565b005b610127610448565b6040516101349190610708565b60405180910390f35b6101456104e0565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630329099f6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101af57600080fd5b505af11580156101c3573d6000803e3d6000fd5b50505050565b6101d16104e0565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635e00e679826040518263ffffffff1660e01b815260040161022c91906106d2565b600060405180830381600087803b15801561024657600080fd5b505af115801561025a573d6000803e3d6000fd5b5050505050565b6102696104e0565b6102736000610567565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6102a66104e0565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663bc1786286040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561031057600080fd5b505af1158015610324573d6000803e3d6000fd5b50505050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e915608b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610399573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103bd9190610738565b905090565b6103ca6104e0565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361043c5760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161043391906106d2565b60405180910390fd5b61044581610567565b50565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f6542b2e6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104db9190610791565b905090565b6104e861062b565b73ffffffffffffffffffffffffffffffffffffffff16610506610275565b73ffffffffffffffffffffffffffffffffffffffff16146105655761052961062b565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161055c91906106d2565b60405180910390fd5b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600033905090565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061066382610638565b9050919050565b61067381610658565b811461067e57600080fd5b50565b6000813590506106908161066a565b92915050565b6000602082840312156106ac576106ab610633565b5b60006106ba84828501610681565b91505092915050565b6106cc81610658565b82525050565b60006020820190506106e760008301846106c3565b92915050565b60008115159050919050565b610702816106ed565b82525050565b600060208201905061071d60008301846106f9565b92915050565b6000815190506107328161066a565b92915050565b60006020828403121561074e5761074d610633565b5b600061075c84828501610723565b91505092915050565b61076e816106ed565b811461077957600080fd5b50565b60008151905061078b81610765565b92915050565b6000602082840312156107a7576107a6610633565b5b60006107b58482850161077c565b9150509291505056fea264697066735822122084570f899c259febf88f28ea6fbe13fbe192de6b14f550ec974419bfb0d1de6c64736f6c634300081e0033",
}

// ExampleRewardManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleRewardManagerMetaData.ABI instead.
var ExampleRewardManagerABI = ExampleRewardManagerMetaData.ABI

// ExampleRewardManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleRewardManagerMetaData.Bin instead.
var ExampleRewardManagerBin = ExampleRewardManagerMetaData.Bin

// DeployExampleRewardManager deploys a new Ethereum contract, binding an instance of ExampleRewardManager to it.
func DeployExampleRewardManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExampleRewardManager, error) {
	parsed, err := ExampleRewardManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleRewardManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExampleRewardManager{ExampleRewardManagerCaller: ExampleRewardManagerCaller{contract: contract}, ExampleRewardManagerTransactor: ExampleRewardManagerTransactor{contract: contract}, ExampleRewardManagerFilterer: ExampleRewardManagerFilterer{contract: contract}}, nil
}

// ExampleRewardManager is an auto generated Go binding around an Ethereum contract.
type ExampleRewardManager struct {
	ExampleRewardManagerCaller     // Read-only binding to the contract
	ExampleRewardManagerTransactor // Write-only binding to the contract
	ExampleRewardManagerFilterer   // Log filterer for contract events
}

// ExampleRewardManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleRewardManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleRewardManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleRewardManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleRewardManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleRewardManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleRewardManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleRewardManagerSession struct {
	Contract     *ExampleRewardManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ExampleRewardManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleRewardManagerCallerSession struct {
	Contract *ExampleRewardManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ExampleRewardManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleRewardManagerTransactorSession struct {
	Contract     *ExampleRewardManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ExampleRewardManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleRewardManagerRaw struct {
	Contract *ExampleRewardManager // Generic contract binding to access the raw methods on
}

// ExampleRewardManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleRewardManagerCallerRaw struct {
	Contract *ExampleRewardManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleRewardManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleRewardManagerTransactorRaw struct {
	Contract *ExampleRewardManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExampleRewardManager creates a new instance of ExampleRewardManager, bound to a specific deployed contract.
func NewExampleRewardManager(address common.Address, backend bind.ContractBackend) (*ExampleRewardManager, error) {
	contract, err := bindExampleRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExampleRewardManager{ExampleRewardManagerCaller: ExampleRewardManagerCaller{contract: contract}, ExampleRewardManagerTransactor: ExampleRewardManagerTransactor{contract: contract}, ExampleRewardManagerFilterer: ExampleRewardManagerFilterer{contract: contract}}, nil
}

// NewExampleRewardManagerCaller creates a new read-only instance of ExampleRewardManager, bound to a specific deployed contract.
func NewExampleRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*ExampleRewardManagerCaller, error) {
	contract, err := bindExampleRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleRewardManagerCaller{contract: contract}, nil
}

// NewExampleRewardManagerTransactor creates a new write-only instance of ExampleRewardManager, bound to a specific deployed contract.
func NewExampleRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleRewardManagerTransactor, error) {
	contract, err := bindExampleRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleRewardManagerTransactor{contract: contract}, nil
}

// NewExampleRewardManagerFilterer creates a new log filterer instance of ExampleRewardManager, bound to a specific deployed contract.
func NewExampleRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleRewardManagerFilterer, error) {
	contract, err := bindExampleRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleRewardManagerFilterer{contract: contract}, nil
}

// bindExampleRewardManager binds a generic wrapper to an already deployed contract.
func bindExampleRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleRewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleRewardManager *ExampleRewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleRewardManager.Contract.ExampleRewardManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleRewardManager *ExampleRewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.ExampleRewardManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleRewardManager *ExampleRewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.ExampleRewardManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleRewardManager *ExampleRewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleRewardManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleRewardManager *ExampleRewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleRewardManager *ExampleRewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.contract.Transact(opts, method, params...)
}

// AreFeeRecipientsAllowed is a free data retrieval call binding the contract method 0xf6542b2e.
//
// Solidity: function areFeeRecipientsAllowed() view returns(bool)
func (_ExampleRewardManager *ExampleRewardManagerCaller) AreFeeRecipientsAllowed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ExampleRewardManager.contract.Call(opts, &out, "areFeeRecipientsAllowed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AreFeeRecipientsAllowed is a free data retrieval call binding the contract method 0xf6542b2e.
//
// Solidity: function areFeeRecipientsAllowed() view returns(bool)
func (_ExampleRewardManager *ExampleRewardManagerSession) AreFeeRecipientsAllowed() (bool, error) {
	return _ExampleRewardManager.Contract.AreFeeRecipientsAllowed(&_ExampleRewardManager.CallOpts)
}

// AreFeeRecipientsAllowed is a free data retrieval call binding the contract method 0xf6542b2e.
//
// Solidity: function areFeeRecipientsAllowed() view returns(bool)
func (_ExampleRewardManager *ExampleRewardManagerCallerSession) AreFeeRecipientsAllowed() (bool, error) {
	return _ExampleRewardManager.Contract.AreFeeRecipientsAllowed(&_ExampleRewardManager.CallOpts)
}

// CurrentRewardAddress is a free data retrieval call binding the contract method 0xe915608b.
//
// Solidity: function currentRewardAddress() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerCaller) CurrentRewardAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExampleRewardManager.contract.Call(opts, &out, "currentRewardAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CurrentRewardAddress is a free data retrieval call binding the contract method 0xe915608b.
//
// Solidity: function currentRewardAddress() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerSession) CurrentRewardAddress() (common.Address, error) {
	return _ExampleRewardManager.Contract.CurrentRewardAddress(&_ExampleRewardManager.CallOpts)
}

// CurrentRewardAddress is a free data retrieval call binding the contract method 0xe915608b.
//
// Solidity: function currentRewardAddress() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerCallerSession) CurrentRewardAddress() (common.Address, error) {
	return _ExampleRewardManager.Contract.CurrentRewardAddress(&_ExampleRewardManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExampleRewardManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerSession) Owner() (common.Address, error) {
	return _ExampleRewardManager.Contract.Owner(&_ExampleRewardManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleRewardManager *ExampleRewardManagerCallerSession) Owner() (common.Address, error) {
	return _ExampleRewardManager.Contract.Owner(&_ExampleRewardManager.CallOpts)
}

// AllowFeeRecipients is a paid mutator transaction binding the contract method 0x0329099f.
//
// Solidity: function allowFeeRecipients() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactor) AllowFeeRecipients(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleRewardManager.contract.Transact(opts, "allowFeeRecipients")
}

// AllowFeeRecipients is a paid mutator transaction binding the contract method 0x0329099f.
//
// Solidity: function allowFeeRecipients() returns()
func (_ExampleRewardManager *ExampleRewardManagerSession) AllowFeeRecipients() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.AllowFeeRecipients(&_ExampleRewardManager.TransactOpts)
}

// AllowFeeRecipients is a paid mutator transaction binding the contract method 0x0329099f.
//
// Solidity: function allowFeeRecipients() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactorSession) AllowFeeRecipients() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.AllowFeeRecipients(&_ExampleRewardManager.TransactOpts)
}

// DisableRewards is a paid mutator transaction binding the contract method 0xbc178628.
//
// Solidity: function disableRewards() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactor) DisableRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleRewardManager.contract.Transact(opts, "disableRewards")
}

// DisableRewards is a paid mutator transaction binding the contract method 0xbc178628.
//
// Solidity: function disableRewards() returns()
func (_ExampleRewardManager *ExampleRewardManagerSession) DisableRewards() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.DisableRewards(&_ExampleRewardManager.TransactOpts)
}

// DisableRewards is a paid mutator transaction binding the contract method 0xbc178628.
//
// Solidity: function disableRewards() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactorSession) DisableRewards() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.DisableRewards(&_ExampleRewardManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleRewardManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleRewardManager *ExampleRewardManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.RenounceOwnership(&_ExampleRewardManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.RenounceOwnership(&_ExampleRewardManager.TransactOpts)
}

// SetRewardAddress is a paid mutator transaction binding the contract method 0x5e00e679.
//
// Solidity: function setRewardAddress(address addr) returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactor) SetRewardAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.contract.Transact(opts, "setRewardAddress", addr)
}

// SetRewardAddress is a paid mutator transaction binding the contract method 0x5e00e679.
//
// Solidity: function setRewardAddress(address addr) returns()
func (_ExampleRewardManager *ExampleRewardManagerSession) SetRewardAddress(addr common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.SetRewardAddress(&_ExampleRewardManager.TransactOpts, addr)
}

// SetRewardAddress is a paid mutator transaction binding the contract method 0x5e00e679.
//
// Solidity: function setRewardAddress(address addr) returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactorSession) SetRewardAddress(addr common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.SetRewardAddress(&_ExampleRewardManager.TransactOpts, addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleRewardManager *ExampleRewardManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.TransferOwnership(&_ExampleRewardManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleRewardManager *ExampleRewardManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleRewardManager.Contract.TransferOwnership(&_ExampleRewardManager.TransactOpts, newOwner)
}

// ExampleRewardManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExampleRewardManager contract.
type ExampleRewardManagerOwnershipTransferredIterator struct {
	Event *ExampleRewardManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExampleRewardManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExampleRewardManagerOwnershipTransferred)
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
		it.Event = new(ExampleRewardManagerOwnershipTransferred)
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
func (it *ExampleRewardManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExampleRewardManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExampleRewardManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ExampleRewardManager contract.
type ExampleRewardManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleRewardManager *ExampleRewardManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExampleRewardManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleRewardManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExampleRewardManagerOwnershipTransferredIterator{contract: _ExampleRewardManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleRewardManager *ExampleRewardManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExampleRewardManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleRewardManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExampleRewardManagerOwnershipTransferred)
				if err := _ExampleRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleRewardManager *ExampleRewardManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ExampleRewardManagerOwnershipTransferred, error) {
	event := new(ExampleRewardManagerOwnershipTransferred)
	if err := _ExampleRewardManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
