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

// AllowListMetaData contains all meta data concerning the AllowList contract.
var AllowListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"precompileAddr\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610d92380380610d9283398181016040528101906100329190610223565b33600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a55760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161009c919061025f565b60405180910390fd5b6100b4816100fc60201b60201c565b5080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061027a565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101f0826101c5565b9050919050565b610200816101e5565b811461020b57600080fd5b50565b60008151905061021d816101f7565b92915050565b600060208284031215610239576102386101c0565b5b60006102478482850161020e565b91505092915050565b610259816101e5565b82525050565b60006020820190506102746000830184610250565b92915050565b610b09806102896000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c80638da5cb5b116100665780638da5cb5b146101315780639015d3711461014f578063d0ebdbe71461017f578063f2fde38b1461019b578063f3ae2415146101b75761009e565b80630aaf7043146100a357806324d7806c146100bf578063704b6c02146100ef578063715018a61461010b57806374a8f10314610115575b600080fd5b6100bd60048036038101906100b89190610966565b6101e7565b005b6100d960048036038101906100d49190610966565b6101fb565b6040516100e691906109ae565b60405180910390f35b61010960048036038101906101049190610966565b6102a8565b005b6101136102bc565b005b61012f600480360381019061012a9190610966565b6102d0565b005b6101396102e4565b60405161014691906109d8565b60405180910390f35b61016960048036038101906101649190610966565b61030d565b60405161017691906109ae565b60405180910390f35b61019960048036038101906101949190610966565b6103bb565b005b6101b560048036038101906101b09190610966565b6103cf565b005b6101d160048036038101906101cc9190610966565b610455565b6040516101de91906109ae565b60405180910390f35b6101ef610502565b6101f881610589565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b815260040161025991906109d8565b602060405180830381865afa158015610276573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029a9190610a29565b905060028114915050919050565b6102b0610502565b6102b981610619565b50565b6102c4610502565b6102ce60006106a9565b565b6102d8610502565b6102e18161076d565b50565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b815260040161036b91906109d8565b602060405180830381865afa158015610388573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ac9190610a29565b90506000811415915050919050565b6103c3610502565b6103cc8161086b565b50565b6103d7610502565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036104495760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161044091906109d8565b60405180910390fd5b610452816106a9565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016104b391906109d8565b602060405180830381865afa1580156104d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104f49190610a29565b905060038114915050919050565b61050a6108fb565b73ffffffffffffffffffffffffffffffffffffffff166105286102e4565b73ffffffffffffffffffffffffffffffffffffffff16146105875761054b6108fb565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161057e91906109d8565b60405180910390fd5b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b81526004016105e491906109d8565b600060405180830381600087803b1580156105fe57600080fd5b505af1158015610612573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b815260040161067491906109d8565b600060405180830381600087803b15801561068e57600080fd5b505af11580156106a2573d6000803e3d6000fd5b5050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16036107db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107d290610ab3565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b815260040161083691906109d8565b600060405180830381600087803b15801561085057600080fd5b505af1158015610864573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b81526004016108c691906109d8565b600060405180830381600087803b1580156108e057600080fd5b505af11580156108f4573d6000803e3d6000fd5b5050505050565b600033905090565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061093382610908565b9050919050565b61094381610928565b811461094e57600080fd5b50565b6000813590506109608161093a565b92915050565b60006020828403121561097c5761097b610903565b5b600061098a84828501610951565b91505092915050565b60008115159050919050565b6109a881610993565b82525050565b60006020820190506109c3600083018461099f565b92915050565b6109d281610928565b82525050565b60006020820190506109ed60008301846109c9565b92915050565b6000819050919050565b610a06816109f3565b8114610a1157600080fd5b50565b600081519050610a23816109fd565b92915050565b600060208284031215610a3f57610a3e610903565b5b6000610a4d84828501610a14565b91505092915050565b600082825260208201905092915050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000610a9d601683610a56565b9150610aa882610a67565b602082019050919050565b60006020820190508181036000830152610acc81610a90565b905091905056fea2646970667358221220653f508c4599dbac96151ace6baf758bb7ce1fcd069eb55b4e32d4a051c77a4064736f6c634300081e0033",
}

// AllowListABI is the input ABI used to generate the binding from.
// Deprecated: Use AllowListMetaData.ABI instead.
var AllowListABI = AllowListMetaData.ABI

// AllowListBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AllowListMetaData.Bin instead.
var AllowListBin = AllowListMetaData.Bin

// DeployAllowList deploys a new Ethereum contract, binding an instance of AllowList to it.
func DeployAllowList(auth *bind.TransactOpts, backend bind.ContractBackend, precompileAddr common.Address) (common.Address, *types.Transaction, *AllowList, error) {
	parsed, err := AllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AllowListBin), backend, precompileAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AllowList{AllowListCaller: AllowListCaller{contract: contract}, AllowListTransactor: AllowListTransactor{contract: contract}, AllowListFilterer: AllowListFilterer{contract: contract}}, nil
}

// AllowList is an auto generated Go binding around an Ethereum contract.
type AllowList struct {
	AllowListCaller     // Read-only binding to the contract
	AllowListTransactor // Write-only binding to the contract
	AllowListFilterer   // Log filterer for contract events
}

// AllowListCaller is an auto generated read-only Go binding around an Ethereum contract.
type AllowListCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AllowListTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AllowListTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AllowListFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AllowListFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AllowListSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AllowListSession struct {
	Contract     *AllowList        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AllowListCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AllowListCallerSession struct {
	Contract *AllowListCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// AllowListTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AllowListTransactorSession struct {
	Contract     *AllowListTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AllowListRaw is an auto generated low-level Go binding around an Ethereum contract.
type AllowListRaw struct {
	Contract *AllowList // Generic contract binding to access the raw methods on
}

// AllowListCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AllowListCallerRaw struct {
	Contract *AllowListCaller // Generic read-only contract binding to access the raw methods on
}

// AllowListTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AllowListTransactorRaw struct {
	Contract *AllowListTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAllowList creates a new instance of AllowList, bound to a specific deployed contract.
func NewAllowList(address common.Address, backend bind.ContractBackend) (*AllowList, error) {
	contract, err := bindAllowList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AllowList{AllowListCaller: AllowListCaller{contract: contract}, AllowListTransactor: AllowListTransactor{contract: contract}, AllowListFilterer: AllowListFilterer{contract: contract}}, nil
}

// NewAllowListCaller creates a new read-only instance of AllowList, bound to a specific deployed contract.
func NewAllowListCaller(address common.Address, caller bind.ContractCaller) (*AllowListCaller, error) {
	contract, err := bindAllowList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AllowListCaller{contract: contract}, nil
}

// NewAllowListTransactor creates a new write-only instance of AllowList, bound to a specific deployed contract.
func NewAllowListTransactor(address common.Address, transactor bind.ContractTransactor) (*AllowListTransactor, error) {
	contract, err := bindAllowList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AllowListTransactor{contract: contract}, nil
}

// NewAllowListFilterer creates a new log filterer instance of AllowList, bound to a specific deployed contract.
func NewAllowListFilterer(address common.Address, filterer bind.ContractFilterer) (*AllowListFilterer, error) {
	contract, err := bindAllowList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AllowListFilterer{contract: contract}, nil
}

// bindAllowList binds a generic wrapper to an already deployed contract.
func bindAllowList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AllowListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AllowList *AllowListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AllowList.Contract.AllowListCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AllowList *AllowListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AllowList.Contract.AllowListTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AllowList *AllowListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AllowList.Contract.AllowListTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AllowList *AllowListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AllowList.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AllowList *AllowListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AllowList.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AllowList *AllowListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AllowList.Contract.contract.Transact(opts, method, params...)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_AllowList *AllowListCaller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AllowList.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_AllowList *AllowListSession) IsAdmin(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsAdmin(&_AllowList.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_AllowList *AllowListCallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsAdmin(&_AllowList.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_AllowList *AllowListCaller) IsEnabled(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AllowList.contract.Call(opts, &out, "isEnabled", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_AllowList *AllowListSession) IsEnabled(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsEnabled(&_AllowList.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_AllowList *AllowListCallerSession) IsEnabled(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsEnabled(&_AllowList.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_AllowList *AllowListCaller) IsManager(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AllowList.contract.Call(opts, &out, "isManager", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_AllowList *AllowListSession) IsManager(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsManager(&_AllowList.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_AllowList *AllowListCallerSession) IsManager(addr common.Address) (bool, error) {
	return _AllowList.Contract.IsManager(&_AllowList.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AllowList *AllowListCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AllowList.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AllowList *AllowListSession) Owner() (common.Address, error) {
	return _AllowList.Contract.Owner(&_AllowList.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AllowList *AllowListCallerSession) Owner() (common.Address, error) {
	return _AllowList.Contract.Owner(&_AllowList.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AllowList *AllowListTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AllowList *AllowListSession) RenounceOwnership() (*types.Transaction, error) {
	return _AllowList.Contract.RenounceOwnership(&_AllowList.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AllowList *AllowListTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AllowList.Contract.RenounceOwnership(&_AllowList.TransactOpts)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_AllowList *AllowListTransactor) Revoke(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "revoke", addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_AllowList *AllowListSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.Revoke(&_AllowList.TransactOpts, addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_AllowList *AllowListTransactorSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.Revoke(&_AllowList.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_AllowList *AllowListTransactor) SetAdmin(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "setAdmin", addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_AllowList *AllowListSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetAdmin(&_AllowList.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_AllowList *AllowListTransactorSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetAdmin(&_AllowList.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_AllowList *AllowListTransactor) SetEnabled(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "setEnabled", addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_AllowList *AllowListSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetEnabled(&_AllowList.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_AllowList *AllowListTransactorSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetEnabled(&_AllowList.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_AllowList *AllowListTransactor) SetManager(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "setManager", addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_AllowList *AllowListSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetManager(&_AllowList.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_AllowList *AllowListTransactorSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.SetManager(&_AllowList.TransactOpts, addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AllowList *AllowListTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AllowList.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AllowList *AllowListSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.TransferOwnership(&_AllowList.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AllowList *AllowListTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AllowList.Contract.TransferOwnership(&_AllowList.TransactOpts, newOwner)
}

// AllowListOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AllowList contract.
type AllowListOwnershipTransferredIterator struct {
	Event *AllowListOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AllowListOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AllowListOwnershipTransferred)
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
		it.Event = new(AllowListOwnershipTransferred)
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
func (it *AllowListOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AllowListOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AllowListOwnershipTransferred represents a OwnershipTransferred event raised by the AllowList contract.
type AllowListOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AllowList *AllowListFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AllowListOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AllowList.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AllowListOwnershipTransferredIterator{contract: _AllowList.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AllowList *AllowListFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AllowListOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AllowList.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AllowListOwnershipTransferred)
				if err := _AllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AllowList *AllowListFilterer) ParseOwnershipTransferred(log types.Log) (*AllowListOwnershipTransferred, error) {
	event := new(AllowListOwnershipTransferred)
	if err := _AllowList.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
