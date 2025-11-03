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

// ExampleTxAllowListMetaData contains all meta data concerning the ExampleTxAllowList contract.
var ExampleTxAllowListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"deployContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5073020000000000000000000000000000000000000233600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100995760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161009091906101f5565b60405180910390fd5b6100a8816100f060201b60201c565b5080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610210565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101df826101b4565b9050919050565b6101ef816101d4565b82525050565b600060208201905061020a60008301846101e6565b92915050565b610bb28061021f6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806374a8f1031161007157806374a8f1031461012a5780638da5cb5b146101465780639015d37114610164578063d0ebdbe714610194578063f2fde38b146101b0578063f3ae2415146101cc576100a9565b80630aaf7043146100ae57806324d7806c146100ca5780636cd5c39b146100fa578063704b6c0214610104578063715018a614610120575b600080fd5b6100c860048036038101906100c391906109b3565b6101fc565b005b6100e460048036038101906100df91906109b3565b610210565b6040516100f191906109fb565b60405180910390f35b6101026102bd565b005b61011e600480360381019061011991906109b3565b6102e9565b005b6101286102fd565b005b610144600480360381019061013f91906109b3565b610311565b005b61014e610325565b60405161015b9190610a25565b60405180910390f35b61017e600480360381019061017991906109b3565b61034e565b60405161018b91906109fb565b60405180910390f35b6101ae60048036038101906101a991906109b3565b6103fc565b005b6101ca60048036038101906101c591906109b3565b610410565b005b6101e660048036038101906101e191906109b3565b610496565b6040516101f391906109fb565b60405180910390f35b610204610543565b61020d816105ca565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b815260040161026e9190610a25565b602060405180830381865afa15801561028b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102af9190610a76565b905060028114915050919050565b6040516102c990610944565b604051809103906000f0801580156102e5573d6000803e3d6000fd5b5050565b6102f1610543565b6102fa8161065a565b50565b610305610543565b61030f60006106ea565b565b610319610543565b610322816107ae565b50565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016103ac9190610a25565b602060405180830381865afa1580156103c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ed9190610a76565b90506000811415915050919050565b610404610543565b61040d816108ac565b50565b610418610543565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361048a5760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016104819190610a25565b60405180910390fd5b610493816106ea565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016104f49190610a25565b602060405180830381865afa158015610511573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105359190610a76565b905060038114915050919050565b61054b61093c565b73ffffffffffffffffffffffffffffffffffffffff16610569610325565b73ffffffffffffffffffffffffffffffffffffffff16146105c85761058c61093c565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105bf9190610a25565b60405180910390fd5b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b81526004016106259190610a25565b600060405180830381600087803b15801561063f57600080fd5b505af1158015610653573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b81526004016106b59190610a25565b600060405180830381600087803b1580156106cf57600080fd5b505af11580156106e3573d6000803e3d6000fd5b5050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff160361081c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161081390610b00565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b81526004016108779190610a25565b600060405180830381600087803b15801561089157600080fd5b505af11580156108a5573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b81526004016109079190610a25565b600060405180830381600087803b15801561092157600080fd5b505af1158015610935573d6000803e3d6000fd5b5050505050565b600033905090565b605c80610b2183390190565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061098082610955565b9050919050565b61099081610975565b811461099b57600080fd5b50565b6000813590506109ad81610987565b92915050565b6000602082840312156109c9576109c8610950565b5b60006109d78482850161099e565b91505092915050565b60008115159050919050565b6109f5816109e0565b82525050565b6000602082019050610a1060008301846109ec565b92915050565b610a1f81610975565b82525050565b6000602082019050610a3a6000830184610a16565b92915050565b6000819050919050565b610a5381610a40565b8114610a5e57600080fd5b50565b600081519050610a7081610a4a565b92915050565b600060208284031215610a8c57610a8b610950565b5b6000610a9a84828501610a61565b91505092915050565b600082825260208201905092915050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000610aea601683610aa3565b9150610af582610ab4565b602082019050919050565b60006020820190508181036000830152610b1981610add565b905091905056fe6080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea2646970667358221220fca91eea49d727340485b03bb9b7aa12355fe29a602dfe104931a95578c9621a64736f6c634300081e0033a2646970667358221220c104aa73c4f174654b1af61a544610e7a121fd1721188a794b30522eabc6f2ca64736f6c634300081e0033",
}

// ExampleTxAllowListABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleTxAllowListMetaData.ABI instead.
var ExampleTxAllowListABI = ExampleTxAllowListMetaData.ABI

// ExampleTxAllowListBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleTxAllowListMetaData.Bin instead.
var ExampleTxAllowListBin = ExampleTxAllowListMetaData.Bin

// DeployExampleTxAllowList deploys a new Ethereum contract, binding an instance of ExampleTxAllowList to it.
func DeployExampleTxAllowList(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExampleTxAllowList, error) {
	parsed, err := ExampleTxAllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleTxAllowListBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExampleTxAllowList{ExampleTxAllowListCaller: ExampleTxAllowListCaller{contract: contract}, ExampleTxAllowListTransactor: ExampleTxAllowListTransactor{contract: contract}, ExampleTxAllowListFilterer: ExampleTxAllowListFilterer{contract: contract}}, nil
}

// ExampleTxAllowList is an auto generated Go binding around an Ethereum contract.
type ExampleTxAllowList struct {
	ExampleTxAllowListCaller     // Read-only binding to the contract
	ExampleTxAllowListTransactor // Write-only binding to the contract
	ExampleTxAllowListFilterer   // Log filterer for contract events
}

// ExampleTxAllowListCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleTxAllowListCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleTxAllowListTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleTxAllowListTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleTxAllowListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleTxAllowListFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleTxAllowListSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleTxAllowListSession struct {
	Contract     *ExampleTxAllowList // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ExampleTxAllowListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleTxAllowListCallerSession struct {
	Contract *ExampleTxAllowListCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ExampleTxAllowListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleTxAllowListTransactorSession struct {
	Contract     *ExampleTxAllowListTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ExampleTxAllowListRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleTxAllowListRaw struct {
	Contract *ExampleTxAllowList // Generic contract binding to access the raw methods on
}

// ExampleTxAllowListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleTxAllowListCallerRaw struct {
	Contract *ExampleTxAllowListCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleTxAllowListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleTxAllowListTransactorRaw struct {
	Contract *ExampleTxAllowListTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExampleTxAllowList creates a new instance of ExampleTxAllowList, bound to a specific deployed contract.
func NewExampleTxAllowList(address common.Address, backend bind.ContractBackend) (*ExampleTxAllowList, error) {
	contract, err := bindExampleTxAllowList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExampleTxAllowList{ExampleTxAllowListCaller: ExampleTxAllowListCaller{contract: contract}, ExampleTxAllowListTransactor: ExampleTxAllowListTransactor{contract: contract}, ExampleTxAllowListFilterer: ExampleTxAllowListFilterer{contract: contract}}, nil
}

// NewExampleTxAllowListCaller creates a new read-only instance of ExampleTxAllowList, bound to a specific deployed contract.
func NewExampleTxAllowListCaller(address common.Address, caller bind.ContractCaller) (*ExampleTxAllowListCaller, error) {
	contract, err := bindExampleTxAllowList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleTxAllowListCaller{contract: contract}, nil
}

// NewExampleTxAllowListTransactor creates a new write-only instance of ExampleTxAllowList, bound to a specific deployed contract.
func NewExampleTxAllowListTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleTxAllowListTransactor, error) {
	contract, err := bindExampleTxAllowList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleTxAllowListTransactor{contract: contract}, nil
}

// NewExampleTxAllowListFilterer creates a new log filterer instance of ExampleTxAllowList, bound to a specific deployed contract.
func NewExampleTxAllowListFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleTxAllowListFilterer, error) {
	contract, err := bindExampleTxAllowList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleTxAllowListFilterer{contract: contract}, nil
}

// bindExampleTxAllowList binds a generic wrapper to an already deployed contract.
func bindExampleTxAllowList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleTxAllowListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleTxAllowList *ExampleTxAllowListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleTxAllowList.Contract.ExampleTxAllowListCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleTxAllowList *ExampleTxAllowListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.ExampleTxAllowListTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleTxAllowList *ExampleTxAllowListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.ExampleTxAllowListTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleTxAllowList *ExampleTxAllowListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleTxAllowList.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleTxAllowList *ExampleTxAllowListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleTxAllowList *ExampleTxAllowListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.contract.Transact(opts, method, params...)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCaller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleTxAllowList.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListSession) IsAdmin(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsAdmin(&_ExampleTxAllowList.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsAdmin(&_ExampleTxAllowList.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCaller) IsEnabled(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleTxAllowList.contract.Call(opts, &out, "isEnabled", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListSession) IsEnabled(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsEnabled(&_ExampleTxAllowList.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCallerSession) IsEnabled(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsEnabled(&_ExampleTxAllowList.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCaller) IsManager(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleTxAllowList.contract.Call(opts, &out, "isManager", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListSession) IsManager(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsManager(&_ExampleTxAllowList.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleTxAllowList *ExampleTxAllowListCallerSession) IsManager(addr common.Address) (bool, error) {
	return _ExampleTxAllowList.Contract.IsManager(&_ExampleTxAllowList.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleTxAllowList *ExampleTxAllowListCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExampleTxAllowList.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleTxAllowList *ExampleTxAllowListSession) Owner() (common.Address, error) {
	return _ExampleTxAllowList.Contract.Owner(&_ExampleTxAllowList.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleTxAllowList *ExampleTxAllowListCallerSession) Owner() (common.Address, error) {
	return _ExampleTxAllowList.Contract.Owner(&_ExampleTxAllowList.CallOpts)
}

// DeployContract is a paid mutator transaction binding the contract method 0x6cd5c39b.
//
// Solidity: function deployContract() returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) DeployContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "deployContract")
}

// DeployContract is a paid mutator transaction binding the contract method 0x6cd5c39b.
//
// Solidity: function deployContract() returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) DeployContract() (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.DeployContract(&_ExampleTxAllowList.TransactOpts)
}

// DeployContract is a paid mutator transaction binding the contract method 0x6cd5c39b.
//
// Solidity: function deployContract() returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) DeployContract() (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.DeployContract(&_ExampleTxAllowList.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.RenounceOwnership(&_ExampleTxAllowList.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.RenounceOwnership(&_ExampleTxAllowList.TransactOpts)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) Revoke(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "revoke", addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.Revoke(&_ExampleTxAllowList.TransactOpts, addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.Revoke(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) SetAdmin(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "setAdmin", addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetAdmin(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetAdmin(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) SetEnabled(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "setEnabled", addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetEnabled(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetEnabled(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) SetManager(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "setManager", addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetManager(&_ExampleTxAllowList.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.SetManager(&_ExampleTxAllowList.TransactOpts, addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleTxAllowList *ExampleTxAllowListSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.TransferOwnership(&_ExampleTxAllowList.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleTxAllowList *ExampleTxAllowListTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleTxAllowList.Contract.TransferOwnership(&_ExampleTxAllowList.TransactOpts, newOwner)
}

// ExampleTxAllowListOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExampleTxAllowList contract.
type ExampleTxAllowListOwnershipTransferredIterator struct {
	Event *ExampleTxAllowListOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExampleTxAllowListOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExampleTxAllowListOwnershipTransferred)
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
		it.Event = new(ExampleTxAllowListOwnershipTransferred)
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
func (it *ExampleTxAllowListOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExampleTxAllowListOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExampleTxAllowListOwnershipTransferred represents a OwnershipTransferred event raised by the ExampleTxAllowList contract.
type ExampleTxAllowListOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleTxAllowList *ExampleTxAllowListFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExampleTxAllowListOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleTxAllowList.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExampleTxAllowListOwnershipTransferredIterator{contract: _ExampleTxAllowList.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleTxAllowList *ExampleTxAllowListFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExampleTxAllowListOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleTxAllowList.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExampleTxAllowListOwnershipTransferred)
				if err := _ExampleTxAllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ExampleTxAllowList *ExampleTxAllowListFilterer) ParseOwnershipTransferred(log types.Log) (*ExampleTxAllowListOwnershipTransferred, error) {
	event := new(ExampleTxAllowListOwnershipTransferred)
	if err := _ExampleTxAllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
