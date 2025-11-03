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

// ExampleWarpMetaData contains all meta data concerning the ExampleWarp contract.
var ExampleWarpMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendWarpMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockchainID\",\"type\":\"bytes32\"}],\"name\":\"validateGetBlockchainID\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"validateInvalidWarpBlockHash\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"validateInvalidWarpMessage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceChainID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"validateWarpBlockHash\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceChainID\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"originSenderAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"validateWarpMessage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040527302000000000000000000000000000000000000056000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015606357600080fd5b50610da4806100736000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806315f0c959146100675780635bd05f061461008357806377ca84db1461009f578063e519286f146100bb578063ee5b48eb146100d7578063f25ec06a146100f3575b600080fd5b610081600480360381019061007c91906106a8565b61010f565b005b61009d600480360381019061009891906107d4565b6101ac565b005b6100b960048036038101906100b4919061085c565b6102df565b005b6100d560048036038101906100d09190610889565b6103b6565b005b6100f160048036038101906100ec91906108dc565b610488565b005b61010d6004803603810190610108919061085c565b61052b565b005b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634213cf786040518163ffffffff1660e01b8152600401602060405180830381865afa15801561017a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019e919061093e565b81146101a957600080fd5b50565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636f825350886040518263ffffffff1660e01b8152600401610208919061097a565b600060405180830381865afa158015610225573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061024e9190610bc8565b915091508061025c57600080fd5b8582600001511461026c57600080fd5b8473ffffffffffffffffffffffffffffffffffffffff16826020015173ffffffffffffffffffffffffffffffffffffffff16146102a857600080fd5b83836040516102b8929190610c63565b6040518091039020826040015180519060200120146102d657600080fd5b50505050505050565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ce7f5929846040518263ffffffff1660e01b815260040161033b919061097a565b606060405180830381865afa158015610358573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061037c9190610ccc565b91509150801561038b57600080fd5b6000801b82600001511461039e57600080fd5b6000801b8260200151146103b157600080fd5b505050565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ce7f5929866040518263ffffffff1660e01b8152600401610412919061097a565b606060405180830381865afa15801561042f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104539190610ccc565b915091508061046157600080fd5b8382600001511461047157600080fd5b8282602001511461048157600080fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ee5b48eb83836040518363ffffffff1660e01b81526004016104e3929190610d4a565b6020604051808303816000875af1158015610502573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610526919061093e565b505050565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636f825350846040518263ffffffff1660e01b8152600401610587919061097a565b600060405180830381865afa1580156105a4573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906105cd9190610bc8565b9150915080156105dc57600080fd5b6000801b8260000151146105ef57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16826020015173ffffffffffffffffffffffffffffffffffffffff161461062c57600080fd5b60405180602001604052806000815250805190602001208260400151805190602001201461065957600080fd5b505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b61068581610672565b811461069057600080fd5b50565b6000813590506106a28161067c565b92915050565b6000602082840312156106be576106bd610668565b5b60006106cc84828501610693565b91505092915050565b600063ffffffff82169050919050565b6106ee816106d5565b81146106f957600080fd5b50565b60008135905061070b816106e5565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061073c82610711565b9050919050565b61074c81610731565b811461075757600080fd5b50565b60008135905061076981610743565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126107945761079361076f565b5b8235905067ffffffffffffffff8111156107b1576107b0610774565b5b6020830191508360018202830111156107cd576107cc610779565b5b9250929050565b6000806000806000608086880312156107f0576107ef610668565b5b60006107fe888289016106fc565b955050602061080f88828901610693565b94505060406108208882890161075a565b935050606086013567ffffffffffffffff8111156108415761084061066d565b5b61084d8882890161077e565b92509250509295509295909350565b60006020828403121561087257610871610668565b5b6000610880848285016106fc565b91505092915050565b6000806000606084860312156108a2576108a1610668565b5b60006108b0868287016106fc565b93505060206108c186828701610693565b92505060406108d286828701610693565b9150509250925092565b600080602083850312156108f3576108f2610668565b5b600083013567ffffffffffffffff8111156109115761091061066d565b5b61091d8582860161077e565b92509250509250929050565b6000815190506109388161067c565b92915050565b60006020828403121561095457610953610668565b5b600061096284828501610929565b91505092915050565b610974816106d5565b82525050565b600060208201905061098f600083018461096b565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6109e38261099a565b810181811067ffffffffffffffff82111715610a0257610a016109ab565b5b80604052505050565b6000610a1561065e565b9050610a2182826109da565b919050565b600080fd5b600081519050610a3a81610743565b92915050565b600080fd5b600067ffffffffffffffff821115610a6057610a5f6109ab565b5b610a698261099a565b9050602081019050919050565b60005b83811015610a94578082015181840152602081019050610a79565b60008484015250505050565b6000610ab3610aae84610a45565b610a0b565b905082815260208101848484011115610acf57610ace610a40565b5b610ada848285610a76565b509392505050565b600082601f830112610af757610af661076f565b5b8151610b07848260208601610aa0565b91505092915050565b600060608284031215610b2657610b25610995565b5b610b306060610a0b565b90506000610b4084828501610929565b6000830152506020610b5484828501610a2b565b602083015250604082015167ffffffffffffffff811115610b7857610b77610a26565b5b610b8484828501610ae2565b60408301525092915050565b60008115159050919050565b610ba581610b90565b8114610bb057600080fd5b50565b600081519050610bc281610b9c565b92915050565b60008060408385031215610bdf57610bde610668565b5b600083015167ffffffffffffffff811115610bfd57610bfc61066d565b5b610c0985828601610b10565b9250506020610c1a85828601610bb3565b9150509250929050565b600081905092915050565b82818337600083830152505050565b6000610c4a8385610c24565b9350610c57838584610c2f565b82840190509392505050565b6000610c70828486610c3e565b91508190509392505050565b600060408284031215610c9257610c91610995565b5b610c9c6040610a0b565b90506000610cac84828501610929565b6000830152506020610cc084828501610929565b60208301525092915050565b60008060608385031215610ce357610ce2610668565b5b6000610cf185828601610c7c565b9250506040610d0285828601610bb3565b9150509250929050565b600082825260208201905092915050565b6000610d298385610d0c565b9350610d36838584610c2f565b610d3f8361099a565b840190509392505050565b60006020820190508181036000830152610d65818486610d1d565b9050939250505056fea2646970667358221220e3f858dc0a344a2ad3c566692a2f4551730ea3e157ede8ed0446a67be11fa39364736f6c634300081e0033",
}

// ExampleWarpABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleWarpMetaData.ABI instead.
var ExampleWarpABI = ExampleWarpMetaData.ABI

// ExampleWarpBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleWarpMetaData.Bin instead.
var ExampleWarpBin = ExampleWarpMetaData.Bin

// DeployExampleWarp deploys a new Ethereum contract, binding an instance of ExampleWarp to it.
func DeployExampleWarp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExampleWarp, error) {
	parsed, err := ExampleWarpMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleWarpBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExampleWarp{ExampleWarpCaller: ExampleWarpCaller{contract: contract}, ExampleWarpTransactor: ExampleWarpTransactor{contract: contract}, ExampleWarpFilterer: ExampleWarpFilterer{contract: contract}}, nil
}

// ExampleWarp is an auto generated Go binding around an Ethereum contract.
type ExampleWarp struct {
	ExampleWarpCaller     // Read-only binding to the contract
	ExampleWarpTransactor // Write-only binding to the contract
	ExampleWarpFilterer   // Log filterer for contract events
}

// ExampleWarpCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleWarpCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleWarpTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleWarpTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleWarpFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleWarpFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleWarpSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleWarpSession struct {
	Contract     *ExampleWarp      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExampleWarpCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleWarpCallerSession struct {
	Contract *ExampleWarpCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ExampleWarpTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleWarpTransactorSession struct {
	Contract     *ExampleWarpTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ExampleWarpRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleWarpRaw struct {
	Contract *ExampleWarp // Generic contract binding to access the raw methods on
}

// ExampleWarpCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleWarpCallerRaw struct {
	Contract *ExampleWarpCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleWarpTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleWarpTransactorRaw struct {
	Contract *ExampleWarpTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExampleWarp creates a new instance of ExampleWarp, bound to a specific deployed contract.
func NewExampleWarp(address common.Address, backend bind.ContractBackend) (*ExampleWarp, error) {
	contract, err := bindExampleWarp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExampleWarp{ExampleWarpCaller: ExampleWarpCaller{contract: contract}, ExampleWarpTransactor: ExampleWarpTransactor{contract: contract}, ExampleWarpFilterer: ExampleWarpFilterer{contract: contract}}, nil
}

// NewExampleWarpCaller creates a new read-only instance of ExampleWarp, bound to a specific deployed contract.
func NewExampleWarpCaller(address common.Address, caller bind.ContractCaller) (*ExampleWarpCaller, error) {
	contract, err := bindExampleWarp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleWarpCaller{contract: contract}, nil
}

// NewExampleWarpTransactor creates a new write-only instance of ExampleWarp, bound to a specific deployed contract.
func NewExampleWarpTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleWarpTransactor, error) {
	contract, err := bindExampleWarp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleWarpTransactor{contract: contract}, nil
}

// NewExampleWarpFilterer creates a new log filterer instance of ExampleWarp, bound to a specific deployed contract.
func NewExampleWarpFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleWarpFilterer, error) {
	contract, err := bindExampleWarp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleWarpFilterer{contract: contract}, nil
}

// bindExampleWarp binds a generic wrapper to an already deployed contract.
func bindExampleWarp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleWarpMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleWarp *ExampleWarpRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleWarp.Contract.ExampleWarpCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleWarp *ExampleWarpRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleWarp.Contract.ExampleWarpTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleWarp *ExampleWarpRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleWarp.Contract.ExampleWarpTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleWarp *ExampleWarpCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleWarp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleWarp *ExampleWarpTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleWarp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleWarp *ExampleWarpTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleWarp.Contract.contract.Transact(opts, method, params...)
}

// ValidateGetBlockchainID is a free data retrieval call binding the contract method 0x15f0c959.
//
// Solidity: function validateGetBlockchainID(bytes32 blockchainID) view returns()
func (_ExampleWarp *ExampleWarpCaller) ValidateGetBlockchainID(opts *bind.CallOpts, blockchainID [32]byte) error {
	var out []interface{}
	err := _ExampleWarp.contract.Call(opts, &out, "validateGetBlockchainID", blockchainID)

	if err != nil {
		return err
	}

	return err

}

// ValidateGetBlockchainID is a free data retrieval call binding the contract method 0x15f0c959.
//
// Solidity: function validateGetBlockchainID(bytes32 blockchainID) view returns()
func (_ExampleWarp *ExampleWarpSession) ValidateGetBlockchainID(blockchainID [32]byte) error {
	return _ExampleWarp.Contract.ValidateGetBlockchainID(&_ExampleWarp.CallOpts, blockchainID)
}

// ValidateGetBlockchainID is a free data retrieval call binding the contract method 0x15f0c959.
//
// Solidity: function validateGetBlockchainID(bytes32 blockchainID) view returns()
func (_ExampleWarp *ExampleWarpCallerSession) ValidateGetBlockchainID(blockchainID [32]byte) error {
	return _ExampleWarp.Contract.ValidateGetBlockchainID(&_ExampleWarp.CallOpts, blockchainID)
}

// ValidateInvalidWarpBlockHash is a free data retrieval call binding the contract method 0x77ca84db.
//
// Solidity: function validateInvalidWarpBlockHash(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpCaller) ValidateInvalidWarpBlockHash(opts *bind.CallOpts, index uint32) error {
	var out []interface{}
	err := _ExampleWarp.contract.Call(opts, &out, "validateInvalidWarpBlockHash", index)

	if err != nil {
		return err
	}

	return err

}

// ValidateInvalidWarpBlockHash is a free data retrieval call binding the contract method 0x77ca84db.
//
// Solidity: function validateInvalidWarpBlockHash(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpSession) ValidateInvalidWarpBlockHash(index uint32) error {
	return _ExampleWarp.Contract.ValidateInvalidWarpBlockHash(&_ExampleWarp.CallOpts, index)
}

// ValidateInvalidWarpBlockHash is a free data retrieval call binding the contract method 0x77ca84db.
//
// Solidity: function validateInvalidWarpBlockHash(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpCallerSession) ValidateInvalidWarpBlockHash(index uint32) error {
	return _ExampleWarp.Contract.ValidateInvalidWarpBlockHash(&_ExampleWarp.CallOpts, index)
}

// ValidateInvalidWarpMessage is a free data retrieval call binding the contract method 0xf25ec06a.
//
// Solidity: function validateInvalidWarpMessage(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpCaller) ValidateInvalidWarpMessage(opts *bind.CallOpts, index uint32) error {
	var out []interface{}
	err := _ExampleWarp.contract.Call(opts, &out, "validateInvalidWarpMessage", index)

	if err != nil {
		return err
	}

	return err

}

// ValidateInvalidWarpMessage is a free data retrieval call binding the contract method 0xf25ec06a.
//
// Solidity: function validateInvalidWarpMessage(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpSession) ValidateInvalidWarpMessage(index uint32) error {
	return _ExampleWarp.Contract.ValidateInvalidWarpMessage(&_ExampleWarp.CallOpts, index)
}

// ValidateInvalidWarpMessage is a free data retrieval call binding the contract method 0xf25ec06a.
//
// Solidity: function validateInvalidWarpMessage(uint32 index) view returns()
func (_ExampleWarp *ExampleWarpCallerSession) ValidateInvalidWarpMessage(index uint32) error {
	return _ExampleWarp.Contract.ValidateInvalidWarpMessage(&_ExampleWarp.CallOpts, index)
}

// ValidateWarpBlockHash is a free data retrieval call binding the contract method 0xe519286f.
//
// Solidity: function validateWarpBlockHash(uint32 index, bytes32 sourceChainID, bytes32 blockHash) view returns()
func (_ExampleWarp *ExampleWarpCaller) ValidateWarpBlockHash(opts *bind.CallOpts, index uint32, sourceChainID [32]byte, blockHash [32]byte) error {
	var out []interface{}
	err := _ExampleWarp.contract.Call(opts, &out, "validateWarpBlockHash", index, sourceChainID, blockHash)

	if err != nil {
		return err
	}

	return err

}

// ValidateWarpBlockHash is a free data retrieval call binding the contract method 0xe519286f.
//
// Solidity: function validateWarpBlockHash(uint32 index, bytes32 sourceChainID, bytes32 blockHash) view returns()
func (_ExampleWarp *ExampleWarpSession) ValidateWarpBlockHash(index uint32, sourceChainID [32]byte, blockHash [32]byte) error {
	return _ExampleWarp.Contract.ValidateWarpBlockHash(&_ExampleWarp.CallOpts, index, sourceChainID, blockHash)
}

// ValidateWarpBlockHash is a free data retrieval call binding the contract method 0xe519286f.
//
// Solidity: function validateWarpBlockHash(uint32 index, bytes32 sourceChainID, bytes32 blockHash) view returns()
func (_ExampleWarp *ExampleWarpCallerSession) ValidateWarpBlockHash(index uint32, sourceChainID [32]byte, blockHash [32]byte) error {
	return _ExampleWarp.Contract.ValidateWarpBlockHash(&_ExampleWarp.CallOpts, index, sourceChainID, blockHash)
}

// ValidateWarpMessage is a free data retrieval call binding the contract method 0x5bd05f06.
//
// Solidity: function validateWarpMessage(uint32 index, bytes32 sourceChainID, address originSenderAddress, bytes payload) view returns()
func (_ExampleWarp *ExampleWarpCaller) ValidateWarpMessage(opts *bind.CallOpts, index uint32, sourceChainID [32]byte, originSenderAddress common.Address, payload []byte) error {
	var out []interface{}
	err := _ExampleWarp.contract.Call(opts, &out, "validateWarpMessage", index, sourceChainID, originSenderAddress, payload)

	if err != nil {
		return err
	}

	return err

}

// ValidateWarpMessage is a free data retrieval call binding the contract method 0x5bd05f06.
//
// Solidity: function validateWarpMessage(uint32 index, bytes32 sourceChainID, address originSenderAddress, bytes payload) view returns()
func (_ExampleWarp *ExampleWarpSession) ValidateWarpMessage(index uint32, sourceChainID [32]byte, originSenderAddress common.Address, payload []byte) error {
	return _ExampleWarp.Contract.ValidateWarpMessage(&_ExampleWarp.CallOpts, index, sourceChainID, originSenderAddress, payload)
}

// ValidateWarpMessage is a free data retrieval call binding the contract method 0x5bd05f06.
//
// Solidity: function validateWarpMessage(uint32 index, bytes32 sourceChainID, address originSenderAddress, bytes payload) view returns()
func (_ExampleWarp *ExampleWarpCallerSession) ValidateWarpMessage(index uint32, sourceChainID [32]byte, originSenderAddress common.Address, payload []byte) error {
	return _ExampleWarp.Contract.ValidateWarpMessage(&_ExampleWarp.CallOpts, index, sourceChainID, originSenderAddress, payload)
}

// SendWarpMessage is a paid mutator transaction binding the contract method 0xee5b48eb.
//
// Solidity: function sendWarpMessage(bytes payload) returns()
func (_ExampleWarp *ExampleWarpTransactor) SendWarpMessage(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _ExampleWarp.contract.Transact(opts, "sendWarpMessage", payload)
}

// SendWarpMessage is a paid mutator transaction binding the contract method 0xee5b48eb.
//
// Solidity: function sendWarpMessage(bytes payload) returns()
func (_ExampleWarp *ExampleWarpSession) SendWarpMessage(payload []byte) (*types.Transaction, error) {
	return _ExampleWarp.Contract.SendWarpMessage(&_ExampleWarp.TransactOpts, payload)
}

// SendWarpMessage is a paid mutator transaction binding the contract method 0xee5b48eb.
//
// Solidity: function sendWarpMessage(bytes payload) returns()
func (_ExampleWarp *ExampleWarpTransactorSession) SendWarpMessage(payload []byte) (*types.Transaction, error) {
	return _ExampleWarp.Contract.SendWarpMessage(&_ExampleWarp.TransactOpts, payload)
}
