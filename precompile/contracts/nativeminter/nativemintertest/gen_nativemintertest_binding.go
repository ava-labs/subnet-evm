// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nativemintertest

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ava-labs/libevm"
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
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

// NativeMinterTestMetaData contains all meta data concerning the NativeMinterTest contract.
var NativeMinterTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nativeMinterPrecompile\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mintNativeCoin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610355380380610355833981810160405281019061003291906100db565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610108565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100a88261007d565b9050919050565b6100b88161009d565b81146100c357600080fd5b50565b6000815190506100d5816100af565b92915050565b6000602082840312156100f1576100f0610078565b5b60006100ff848285016100c6565b91505092915050565b61023e806101176000396000f3fe6080604052600436106100225760003560e01c80634f5aaaba1461002e57610029565b3661002957005b600080fd5b34801561003a57600080fd5b5061005560048036038101906100509190610181565b610057565b005b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634f5aaaba83836040518363ffffffff1660e01b81526004016100b29291906101df565b600060405180830381600087803b1580156100cc57600080fd5b505af11580156100e0573d6000803e3d6000fd5b505050505050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610118826100ed565b9050919050565b6101288161010d565b811461013357600080fd5b50565b6000813590506101458161011f565b92915050565b6000819050919050565b61015e8161014b565b811461016957600080fd5b50565b60008135905061017b81610155565b92915050565b60008060408385031215610198576101976100e8565b5b60006101a685828601610136565b92505060206101b78582860161016c565b9150509250929050565b6101ca8161010d565b82525050565b6101d98161014b565b82525050565b60006040820190506101f460008301856101c1565b61020160208301846101d0565b939250505056fea2646970667358221220d8815e6e2ef5f9f594d8ca3800c6e2022c21c6ba5da974f0b81a27b74cea21f364736f6c634300081e0033",
}

// NativeMinterTestABI is the input ABI used to generate the binding from.
// Deprecated: Use NativeMinterTestMetaData.ABI instead.
var NativeMinterTestABI = NativeMinterTestMetaData.ABI

// NativeMinterTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NativeMinterTestMetaData.Bin instead.
var NativeMinterTestBin = NativeMinterTestMetaData.Bin

// DeployNativeMinterTest deploys a new Ethereum contract, binding an instance of NativeMinterTest to it.
func DeployNativeMinterTest(auth *bind.TransactOpts, backend bind.ContractBackend, nativeMinterPrecompile common.Address) (common.Address, *types.Transaction, *NativeMinterTest, error) {
	parsed, err := NativeMinterTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NativeMinterTestBin), backend, nativeMinterPrecompile)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NativeMinterTest{NativeMinterTestCaller: NativeMinterTestCaller{contract: contract}, NativeMinterTestTransactor: NativeMinterTestTransactor{contract: contract}, NativeMinterTestFilterer: NativeMinterTestFilterer{contract: contract}}, nil
}

// NativeMinterTest is an auto generated Go binding around an Ethereum contract.
type NativeMinterTest struct {
	NativeMinterTestCaller     // Read-only binding to the contract
	NativeMinterTestTransactor // Write-only binding to the contract
	NativeMinterTestFilterer   // Log filterer for contract events
}

// NativeMinterTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type NativeMinterTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeMinterTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NativeMinterTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeMinterTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NativeMinterTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NativeMinterTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NativeMinterTestSession struct {
	Contract     *NativeMinterTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NativeMinterTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NativeMinterTestCallerSession struct {
	Contract *NativeMinterTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// NativeMinterTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NativeMinterTestTransactorSession struct {
	Contract     *NativeMinterTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// NativeMinterTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type NativeMinterTestRaw struct {
	Contract *NativeMinterTest // Generic contract binding to access the raw methods on
}

// NativeMinterTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NativeMinterTestCallerRaw struct {
	Contract *NativeMinterTestCaller // Generic read-only contract binding to access the raw methods on
}

// NativeMinterTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NativeMinterTestTransactorRaw struct {
	Contract *NativeMinterTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNativeMinterTest creates a new instance of NativeMinterTest, bound to a specific deployed contract.
func NewNativeMinterTest(address common.Address, backend bind.ContractBackend) (*NativeMinterTest, error) {
	contract, err := bindNativeMinterTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NativeMinterTest{NativeMinterTestCaller: NativeMinterTestCaller{contract: contract}, NativeMinterTestTransactor: NativeMinterTestTransactor{contract: contract}, NativeMinterTestFilterer: NativeMinterTestFilterer{contract: contract}}, nil
}

// NewNativeMinterTestCaller creates a new read-only instance of NativeMinterTest, bound to a specific deployed contract.
func NewNativeMinterTestCaller(address common.Address, caller bind.ContractCaller) (*NativeMinterTestCaller, error) {
	contract, err := bindNativeMinterTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NativeMinterTestCaller{contract: contract}, nil
}

// NewNativeMinterTestTransactor creates a new write-only instance of NativeMinterTest, bound to a specific deployed contract.
func NewNativeMinterTestTransactor(address common.Address, transactor bind.ContractTransactor) (*NativeMinterTestTransactor, error) {
	contract, err := bindNativeMinterTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NativeMinterTestTransactor{contract: contract}, nil
}

// NewNativeMinterTestFilterer creates a new log filterer instance of NativeMinterTest, bound to a specific deployed contract.
func NewNativeMinterTestFilterer(address common.Address, filterer bind.ContractFilterer) (*NativeMinterTestFilterer, error) {
	contract, err := bindNativeMinterTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NativeMinterTestFilterer{contract: contract}, nil
}

// bindNativeMinterTest binds a generic wrapper to an already deployed contract.
func bindNativeMinterTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NativeMinterTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NativeMinterTest *NativeMinterTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeMinterTest.Contract.NativeMinterTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NativeMinterTest *NativeMinterTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.NativeMinterTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NativeMinterTest *NativeMinterTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.NativeMinterTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NativeMinterTest *NativeMinterTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NativeMinterTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NativeMinterTest *NativeMinterTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NativeMinterTest *NativeMinterTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.contract.Transact(opts, method, params...)
}

// MintNativeCoin is a paid mutator transaction binding the contract method 0x4f5aaaba.
//
// Solidity: function mintNativeCoin(address addr, uint256 amount) returns()
func (_NativeMinterTest *NativeMinterTestTransactor) MintNativeCoin(opts *bind.TransactOpts, addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeMinterTest.contract.Transact(opts, "mintNativeCoin", addr, amount)
}

// MintNativeCoin is a paid mutator transaction binding the contract method 0x4f5aaaba.
//
// Solidity: function mintNativeCoin(address addr, uint256 amount) returns()
func (_NativeMinterTest *NativeMinterTestSession) MintNativeCoin(addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.MintNativeCoin(&_NativeMinterTest.TransactOpts, addr, amount)
}

// MintNativeCoin is a paid mutator transaction binding the contract method 0x4f5aaaba.
//
// Solidity: function mintNativeCoin(address addr, uint256 amount) returns()
func (_NativeMinterTest *NativeMinterTestTransactorSession) MintNativeCoin(addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NativeMinterTest.Contract.MintNativeCoin(&_NativeMinterTest.TransactOpts, addr, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NativeMinterTest *NativeMinterTestTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NativeMinterTest.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NativeMinterTest *NativeMinterTestSession) Receive() (*types.Transaction, error) {
	return _NativeMinterTest.Contract.Receive(&_NativeMinterTest.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NativeMinterTest *NativeMinterTestTransactorSession) Receive() (*types.Transaction, error) {
	return _NativeMinterTest.Contract.Receive(&_NativeMinterTest.TransactOpts)
}
