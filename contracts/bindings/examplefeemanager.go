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

// FeeConfig is an auto generated low-level Go binding around an user-defined struct.
type FeeConfig struct {
	GasLimit                 *big.Int
	TargetBlockRate          *big.Int
	MinBaseFee               *big.Int
	TargetGas                *big.Int
	BaseFeeChangeDenominator *big.Int
	MinBlockGasCost          *big.Int
	MaxBlockGasCost          *big.Int
	BlockGasCostStep         *big.Int
}

// ExampleFeeManagerMetaData contains all meta data concerning the ExampleFeeManager contract.
var ExampleFeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"enableCChainFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetBlockRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockGasCostStep\",\"type\":\"uint256\"}],\"internalType\":\"structFeeConfig\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"enableCustomFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enableWAGMIFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentFeeConfig\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetBlockRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockGasCostStep\",\"type\":\"uint256\"}],\"internalType\":\"structFeeConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfigLastChangedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052730200000000000000000000000000000000000003600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006557600080fd5b5073020000000000000000000000000000000000000333600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100ee5760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100e5919061024a565b60405180910390fd5b6100fd8161014560201b60201c565b5080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610265565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061023482610209565b9050919050565b61024481610229565b82525050565b600060208201905061025f600083018461023b565b92915050565b6114d2806102746000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806374a8f103116100975780639e05549a116100665780639e05549a14610224578063d0ebdbe714610242578063f2fde38b1461025e578063f3ae24151461027a576100f5565b806374a8f103146101b057806385c1b4ac146101cc5780638da5cb5b146101d65780639015d371146101f4576100f5565b806352965cfc116100d357806352965cfc146101645780636f0edc9d14610180578063704b6c021461018a578063715018a6146101a6576100f5565b80630aaf7043146100fa57806324d7806c1461011657806341f5772814610146575b600080fd5b610114600480360381019061010f9190610efa565b6102aa565b005b610130600480360381019061012b9190610efa565b6102be565b60405161013d9190610f42565b60405180910390f35b61014e61036b565b60405161015b9190611018565b60405180910390f35b61017e600480360381019061017991906111bb565b61045c565b005b610188610562565b005b6101a4600480360381019061019f9190610efa565b61065c565b005b6101ae610670565b005b6101ca60048036038101906101c59190610efa565b610684565b005b6101d4610698565b005b6101de610791565b6040516101eb91906111f8565b60405180910390f35b61020e60048036038101906102099190610efa565b6107ba565b60405161021b9190610f42565b60405180910390f35b61022c610868565b6040516102399190611222565b60405180910390f35b61025c60048036038101906102579190610efa565b610900565b005b61027860048036038101906102739190610efa565b610914565b005b610294600480360381019061028f9190610efa565b61099a565b6040516102a19190610f42565b60405180910390f35b6102b2610a47565b6102bb81610ace565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b815260040161031c91906111f8565b602060405180830381865afa158015610339573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061035d9190611252565b905060028114915050919050565b610373610e48565b61037b610e48565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635fbbc0d26040518163ffffffff1660e01b815260040161010060405180830381865afa1580156103e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061040d919061127f565b88600001896020018a6040018b6060018c6080018d60a0018e60c0018f60e001888152508881525088815250888152508881525088815250888152508881525050505050505050508091505090565b610465336107ba565b6104a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049b90611392565b60405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638f10b586826000015183602001518460400151856060015186608001518760a001518860c001518960e001516040518963ffffffff1660e01b815260040161052d9897969594939291906113b2565b600060405180830381600087803b15801561054757600080fd5b505af115801561055b573d6000803e3d6000fd5b5050505050565b61056b336107ba565b6105aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105a190611392565b60405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638f10b5866301312d006002633b9aca006305f5e10060306000629896806207a1206040518963ffffffff1660e01b81526004016106289897969594939291906113b2565b600060405180830381600087803b15801561064257600080fd5b505af1158015610656573d6000803e3d6000fd5b50505050565b610664610a47565b61066d81610b5e565b50565b610678610a47565b6106826000610bee565b565b61068c610a47565b61069581610cb2565b50565b6106a1336107ba565b6106e0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106d790611392565b60405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638f10b586627a120060026405d21dba0062e4e1c060246000620f4240620186a06040518963ffffffff1660e01b815260040161075d9897969594939291906113b2565b600060405180830381600087803b15801561077757600080fd5b505af115801561078b573d6000803e3d6000fd5b50505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b815260040161081891906111f8565b602060405180830381865afa158015610835573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108599190611252565b90506000811415915050919050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639e05549a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108fb9190611252565b905090565b610908610a47565b61091181610db0565b50565b61091c610a47565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361098e5760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161098591906111f8565b60405180910390fd5b61099781610bee565b50565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016109f891906111f8565b602060405180830381865afa158015610a15573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a399190611252565b905060038114915050919050565b610a4f610e40565b73ffffffffffffffffffffffffffffffffffffffff16610a6d610791565b73ffffffffffffffffffffffffffffffffffffffff1614610acc57610a90610e40565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401610ac391906111f8565b60405180910390fd5b565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b8152600401610b2991906111f8565b600060405180830381600087803b158015610b4357600080fd5b505af1158015610b57573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b8152600401610bb991906111f8565b600060405180830381600087803b158015610bd357600080fd5b505af1158015610be7573d6000803e3d6000fd5b5050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1603610d20576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d179061147c565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b8152600401610d7b91906111f8565b600060405180830381600087803b158015610d9557600080fd5b505af1158015610da9573d6000803e3d6000fd5b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b8152600401610e0b91906111f8565b600060405180830381600087803b158015610e2557600080fd5b505af1158015610e39573d6000803e3d6000fd5b5050505050565b600033905090565b60405180610100016040528060008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b6000604051905090565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610ec782610e9c565b9050919050565b610ed781610ebc565b8114610ee257600080fd5b50565b600081359050610ef481610ece565b92915050565b600060208284031215610f1057610f0f610e97565b5b6000610f1e84828501610ee5565b91505092915050565b60008115159050919050565b610f3c81610f27565b82525050565b6000602082019050610f576000830184610f33565b92915050565b6000819050919050565b610f7081610f5d565b82525050565b61010082016000820151610f8d6000850182610f67565b506020820151610fa06020850182610f67565b506040820151610fb36040850182610f67565b506060820151610fc66060850182610f67565b506080820151610fd96080850182610f67565b5060a0820151610fec60a0850182610f67565b5060c0820151610fff60c0850182610f67565b5060e082015161101260e0850182610f67565b50505050565b60006101008201905061102e6000830184610f76565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61108282611039565b810181811067ffffffffffffffff821117156110a1576110a061104a565b5b80604052505050565b60006110b4610e8d565b90506110c08282611079565b919050565b6110ce81610f5d565b81146110d957600080fd5b50565b6000813590506110eb816110c5565b92915050565b6000610100828403121561110857611107611034565b5b6111136101006110aa565b90506000611123848285016110dc565b6000830152506020611137848285016110dc565b602083015250604061114b848285016110dc565b604083015250606061115f848285016110dc565b6060830152506080611173848285016110dc565b60808301525060a0611187848285016110dc565b60a08301525060c061119b848285016110dc565b60c08301525060e06111af848285016110dc565b60e08301525092915050565b600061010082840312156111d2576111d1610e97565b5b60006111e0848285016110f1565b91505092915050565b6111f281610ebc565b82525050565b600060208201905061120d60008301846111e9565b92915050565b61121c81610f5d565b82525050565b60006020820190506112376000830184611213565b92915050565b60008151905061124c816110c5565b92915050565b60006020828403121561126857611267610e97565b5b60006112768482850161123d565b91505092915050565b600080600080600080600080610100898b0312156112a05761129f610e97565b5b60006112ae8b828c0161123d565b98505060206112bf8b828c0161123d565b97505060406112d08b828c0161123d565b96505060606112e18b828c0161123d565b95505060806112f28b828c0161123d565b94505060a06113038b828c0161123d565b93505060c06113148b828c0161123d565b92505060e06113258b828c0161123d565b9150509295985092959890939650565b600082825260208201905092915050565b7f6e6f7420656e61626c6564000000000000000000000000000000000000000000600082015250565b600061137c600b83611335565b915061138782611346565b602082019050919050565b600060208201905081810360008301526113ab8161136f565b9050919050565b6000610100820190506113c8600083018b611213565b6113d5602083018a611213565b6113e26040830189611213565b6113ef6060830188611213565b6113fc6080830187611213565b61140960a0830186611213565b61141660c0830185611213565b61142360e0830184611213565b9998505050505050505050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000611466601683611335565b915061147182611430565b602082019050919050565b6000602082019050818103600083015261149581611459565b905091905056fea2646970667358221220d0b91fea29c818847d4dcf4db121a8709dc54729cf872805314a42f97de6966864736f6c634300081e0033",
}

// ExampleFeeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleFeeManagerMetaData.ABI instead.
var ExampleFeeManagerABI = ExampleFeeManagerMetaData.ABI

// ExampleFeeManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleFeeManagerMetaData.Bin instead.
var ExampleFeeManagerBin = ExampleFeeManagerMetaData.Bin

// DeployExampleFeeManager deploys a new Ethereum contract, binding an instance of ExampleFeeManager to it.
func DeployExampleFeeManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ExampleFeeManager, error) {
	parsed, err := ExampleFeeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleFeeManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ExampleFeeManager{ExampleFeeManagerCaller: ExampleFeeManagerCaller{contract: contract}, ExampleFeeManagerTransactor: ExampleFeeManagerTransactor{contract: contract}, ExampleFeeManagerFilterer: ExampleFeeManagerFilterer{contract: contract}}, nil
}

// ExampleFeeManager is an auto generated Go binding around an Ethereum contract.
type ExampleFeeManager struct {
	ExampleFeeManagerCaller     // Read-only binding to the contract
	ExampleFeeManagerTransactor // Write-only binding to the contract
	ExampleFeeManagerFilterer   // Log filterer for contract events
}

// ExampleFeeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleFeeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleFeeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleFeeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleFeeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleFeeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleFeeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleFeeManagerSession struct {
	Contract     *ExampleFeeManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ExampleFeeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleFeeManagerCallerSession struct {
	Contract *ExampleFeeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ExampleFeeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleFeeManagerTransactorSession struct {
	Contract     *ExampleFeeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ExampleFeeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleFeeManagerRaw struct {
	Contract *ExampleFeeManager // Generic contract binding to access the raw methods on
}

// ExampleFeeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleFeeManagerCallerRaw struct {
	Contract *ExampleFeeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleFeeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleFeeManagerTransactorRaw struct {
	Contract *ExampleFeeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExampleFeeManager creates a new instance of ExampleFeeManager, bound to a specific deployed contract.
func NewExampleFeeManager(address common.Address, backend bind.ContractBackend) (*ExampleFeeManager, error) {
	contract, err := bindExampleFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExampleFeeManager{ExampleFeeManagerCaller: ExampleFeeManagerCaller{contract: contract}, ExampleFeeManagerTransactor: ExampleFeeManagerTransactor{contract: contract}, ExampleFeeManagerFilterer: ExampleFeeManagerFilterer{contract: contract}}, nil
}

// NewExampleFeeManagerCaller creates a new read-only instance of ExampleFeeManager, bound to a specific deployed contract.
func NewExampleFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*ExampleFeeManagerCaller, error) {
	contract, err := bindExampleFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleFeeManagerCaller{contract: contract}, nil
}

// NewExampleFeeManagerTransactor creates a new write-only instance of ExampleFeeManager, bound to a specific deployed contract.
func NewExampleFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleFeeManagerTransactor, error) {
	contract, err := bindExampleFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleFeeManagerTransactor{contract: contract}, nil
}

// NewExampleFeeManagerFilterer creates a new log filterer instance of ExampleFeeManager, bound to a specific deployed contract.
func NewExampleFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleFeeManagerFilterer, error) {
	contract, err := bindExampleFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleFeeManagerFilterer{contract: contract}, nil
}

// bindExampleFeeManager binds a generic wrapper to an already deployed contract.
func bindExampleFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleFeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleFeeManager *ExampleFeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleFeeManager.Contract.ExampleFeeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleFeeManager *ExampleFeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.ExampleFeeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleFeeManager *ExampleFeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.ExampleFeeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExampleFeeManager *ExampleFeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExampleFeeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExampleFeeManager *ExampleFeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExampleFeeManager *ExampleFeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.contract.Transact(opts, method, params...)
}

// GetCurrentFeeConfig is a free data retrieval call binding the contract method 0x41f57728.
//
// Solidity: function getCurrentFeeConfig() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_ExampleFeeManager *ExampleFeeManagerCaller) GetCurrentFeeConfig(opts *bind.CallOpts) (FeeConfig, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "getCurrentFeeConfig")

	if err != nil {
		return *new(FeeConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(FeeConfig)).(*FeeConfig)

	return out0, err

}

// GetCurrentFeeConfig is a free data retrieval call binding the contract method 0x41f57728.
//
// Solidity: function getCurrentFeeConfig() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_ExampleFeeManager *ExampleFeeManagerSession) GetCurrentFeeConfig() (FeeConfig, error) {
	return _ExampleFeeManager.Contract.GetCurrentFeeConfig(&_ExampleFeeManager.CallOpts)
}

// GetCurrentFeeConfig is a free data retrieval call binding the contract method 0x41f57728.
//
// Solidity: function getCurrentFeeConfig() view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) GetCurrentFeeConfig() (FeeConfig, error) {
	return _ExampleFeeManager.Contract.GetCurrentFeeConfig(&_ExampleFeeManager.CallOpts)
}

// GetFeeConfigLastChangedAt is a free data retrieval call binding the contract method 0x9e05549a.
//
// Solidity: function getFeeConfigLastChangedAt() view returns(uint256)
func (_ExampleFeeManager *ExampleFeeManagerCaller) GetFeeConfigLastChangedAt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "getFeeConfigLastChangedAt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeeConfigLastChangedAt is a free data retrieval call binding the contract method 0x9e05549a.
//
// Solidity: function getFeeConfigLastChangedAt() view returns(uint256)
func (_ExampleFeeManager *ExampleFeeManagerSession) GetFeeConfigLastChangedAt() (*big.Int, error) {
	return _ExampleFeeManager.Contract.GetFeeConfigLastChangedAt(&_ExampleFeeManager.CallOpts)
}

// GetFeeConfigLastChangedAt is a free data retrieval call binding the contract method 0x9e05549a.
//
// Solidity: function getFeeConfigLastChangedAt() view returns(uint256)
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) GetFeeConfigLastChangedAt() (*big.Int, error) {
	return _ExampleFeeManager.Contract.GetFeeConfigLastChangedAt(&_ExampleFeeManager.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCaller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerSession) IsAdmin(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsAdmin(&_ExampleFeeManager.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsAdmin(&_ExampleFeeManager.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCaller) IsEnabled(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "isEnabled", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerSession) IsEnabled(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsEnabled(&_ExampleFeeManager.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) IsEnabled(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsEnabled(&_ExampleFeeManager.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCaller) IsManager(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "isManager", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerSession) IsManager(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsManager(&_ExampleFeeManager.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) IsManager(addr common.Address) (bool, error) {
	return _ExampleFeeManager.Contract.IsManager(&_ExampleFeeManager.CallOpts, addr)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleFeeManager *ExampleFeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExampleFeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleFeeManager *ExampleFeeManagerSession) Owner() (common.Address, error) {
	return _ExampleFeeManager.Contract.Owner(&_ExampleFeeManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExampleFeeManager *ExampleFeeManagerCallerSession) Owner() (common.Address, error) {
	return _ExampleFeeManager.Contract.Owner(&_ExampleFeeManager.CallOpts)
}

// EnableCChainFees is a paid mutator transaction binding the contract method 0x85c1b4ac.
//
// Solidity: function enableCChainFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) EnableCChainFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "enableCChainFees")
}

// EnableCChainFees is a paid mutator transaction binding the contract method 0x85c1b4ac.
//
// Solidity: function enableCChainFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) EnableCChainFees() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableCChainFees(&_ExampleFeeManager.TransactOpts)
}

// EnableCChainFees is a paid mutator transaction binding the contract method 0x85c1b4ac.
//
// Solidity: function enableCChainFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) EnableCChainFees() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableCChainFees(&_ExampleFeeManager.TransactOpts)
}

// EnableCustomFees is a paid mutator transaction binding the contract method 0x52965cfc.
//
// Solidity: function enableCustomFees((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256) config) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) EnableCustomFees(opts *bind.TransactOpts, config FeeConfig) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "enableCustomFees", config)
}

// EnableCustomFees is a paid mutator transaction binding the contract method 0x52965cfc.
//
// Solidity: function enableCustomFees((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256) config) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) EnableCustomFees(config FeeConfig) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableCustomFees(&_ExampleFeeManager.TransactOpts, config)
}

// EnableCustomFees is a paid mutator transaction binding the contract method 0x52965cfc.
//
// Solidity: function enableCustomFees((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256) config) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) EnableCustomFees(config FeeConfig) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableCustomFees(&_ExampleFeeManager.TransactOpts, config)
}

// EnableWAGMIFees is a paid mutator transaction binding the contract method 0x6f0edc9d.
//
// Solidity: function enableWAGMIFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) EnableWAGMIFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "enableWAGMIFees")
}

// EnableWAGMIFees is a paid mutator transaction binding the contract method 0x6f0edc9d.
//
// Solidity: function enableWAGMIFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) EnableWAGMIFees() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableWAGMIFees(&_ExampleFeeManager.TransactOpts)
}

// EnableWAGMIFees is a paid mutator transaction binding the contract method 0x6f0edc9d.
//
// Solidity: function enableWAGMIFees() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) EnableWAGMIFees() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.EnableWAGMIFees(&_ExampleFeeManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.RenounceOwnership(&_ExampleFeeManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.RenounceOwnership(&_ExampleFeeManager.TransactOpts)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) Revoke(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "revoke", addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.Revoke(&_ExampleFeeManager.TransactOpts, addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.Revoke(&_ExampleFeeManager.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) SetAdmin(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "setAdmin", addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetAdmin(&_ExampleFeeManager.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetAdmin(&_ExampleFeeManager.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) SetEnabled(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "setEnabled", addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetEnabled(&_ExampleFeeManager.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetEnabled(&_ExampleFeeManager.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) SetManager(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "setManager", addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetManager(&_ExampleFeeManager.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.SetManager(&_ExampleFeeManager.TransactOpts, addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleFeeManager *ExampleFeeManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.TransferOwnership(&_ExampleFeeManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExampleFeeManager *ExampleFeeManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExampleFeeManager.Contract.TransferOwnership(&_ExampleFeeManager.TransactOpts, newOwner)
}

// ExampleFeeManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExampleFeeManager contract.
type ExampleFeeManagerOwnershipTransferredIterator struct {
	Event *ExampleFeeManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExampleFeeManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExampleFeeManagerOwnershipTransferred)
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
		it.Event = new(ExampleFeeManagerOwnershipTransferred)
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
func (it *ExampleFeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExampleFeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExampleFeeManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ExampleFeeManager contract.
type ExampleFeeManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleFeeManager *ExampleFeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExampleFeeManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleFeeManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExampleFeeManagerOwnershipTransferredIterator{contract: _ExampleFeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExampleFeeManager *ExampleFeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExampleFeeManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExampleFeeManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExampleFeeManagerOwnershipTransferred)
				if err := _ExampleFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ExampleFeeManager *ExampleFeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ExampleFeeManagerOwnershipTransferred, error) {
	event := new(ExampleFeeManagerOwnershipTransferred)
	if err := _ExampleFeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
