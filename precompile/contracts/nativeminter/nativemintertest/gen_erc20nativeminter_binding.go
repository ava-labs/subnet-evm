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

// ERC20NativeMinterMetaData contains all meta data concerning the ERC20NativeMinter contract.
var ERC20NativeMinterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nativeMinterPrecompile\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"initSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Mintdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"mintdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051612898380380612898833981810160405281019061003291906105ea565b81336040518060400160405280601681526020017f45524332304e61746976654d696e746572546f6b656e000000000000000000008152506040518060400160405280600481526020017f584d504c0000000000000000000000000000000000000000000000000000000081525081600390816100af9190610870565b5080600490816100bf9190610870565b505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036101345760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161012b9190610951565b60405180910390fd5b610143816101de60201b60201c565b5080600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505081600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101d733826102a460201b60201c565b5050610a30565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036103165760006040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161030d9190610951565b60405180910390fd5b6103286000838361032c60201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361037e578060026000828254610372919061099b565b92505081905550610451565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561040a578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610401939291906109de565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361049a57806002600082825403925050819055506104e7565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105449190610a15565b60405180910390a3505050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061058182610556565b9050919050565b61059181610576565b811461059c57600080fd5b50565b6000815190506105ae81610588565b92915050565b6000819050919050565b6105c7816105b4565b81146105d257600080fd5b50565b6000815190506105e4816105be565b92915050565b6000806040838503121561060157610600610551565b5b600061060f8582860161059f565b9250506020610620858286016105d5565b9150509250929050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806106ab57607f821691505b6020821081036106be576106bd610664565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026107267fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826106e9565b61073086836106e9565b95508019841693508086168417925050509392505050565b6000819050919050565b600061076d610768610763846105b4565b610748565b6105b4565b9050919050565b6000819050919050565b61078783610752565b61079b61079382610774565b8484546106f6565b825550505050565b600090565b6107b06107a3565b6107bb81848461077e565b505050565b5b818110156107df576107d46000826107a8565b6001810190506107c1565b5050565b601f821115610824576107f5816106c4565b6107fe846106d9565b8101602085101561080d578190505b610821610819856106d9565b8301826107c0565b50505b505050565b600082821c905092915050565b600061084760001984600802610829565b1980831691505092915050565b60006108608383610836565b9150826002028217905092915050565b6108798261062a565b67ffffffffffffffff81111561089257610891610635565b5b61089c8254610693565b6108a78282856107e3565b600060209050601f8311600181146108da57600084156108c8578287015190505b6108d28582610854565b86555061093a565b601f1984166108e8866106c4565b60005b82811015610910578489015182556001820191506020850194506020810190506108eb565b8683101561092d5784890151610929601f891682610836565b8355505b6001600288020188555050505b505050505050565b61094b81610576565b82525050565b60006020820190506109666000830184610942565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006109a6826105b4565b91506109b1836105b4565b92508282019050808211156109c9576109c861096c565b5b92915050565b6109d8816105b4565b82525050565b60006060820190506109f36000830186610942565b610a0060208301856109cf565b610a0d60408301846109cf565b949350505050565b6000602082019050610a2a60008301846109cf565b92915050565b611e5980610a3f6000396000f3fe6080604052600436106101405760003560e01c8063715018a6116100b6578063a9059cbb1161006f578063a9059cbb1461045a578063d0e30db014610497578063d0ebdbe7146104a1578063dd62ed3e146104ca578063f2fde38b14610507578063f3ae24151461053057610140565b8063715018a61461035e57806374a8f103146103755780638da5cb5b1461039e5780639015d371146103c957806395d89b41146104065780639dc29fac1461043157610140565b806323b872dd1161010857806323b872dd1461022a57806324d7806c14610267578063313ce567146102a457806340c10f19146102cf578063704b6c02146102f857806370a082311461032157610140565b80630356b6cd1461014557806306fdde031461016e578063095ea7b3146101995780630aaf7043146101d657806318160ddd146101ff575b600080fd5b34801561015157600080fd5b5061016c600480360381019061016791906118e8565b61056d565b005b34801561017a57600080fd5b50610183610657565b60405161019091906119a5565b60405180910390f35b3480156101a557600080fd5b506101c060048036038101906101bb9190611a25565b6106e9565b6040516101cd9190611a80565b60405180910390f35b3480156101e257600080fd5b506101fd60048036038101906101f89190611a9b565b61070c565b005b34801561020b57600080fd5b50610214610770565b6040516102219190611ad7565b60405180910390f35b34801561023657600080fd5b50610251600480360381019061024c9190611af2565b61077a565b60405161025e9190611a80565b60405180910390f35b34801561027357600080fd5b5061028e60048036038101906102899190611a9b565b6107a9565b60405161029b9190611a80565b60405180910390f35b3480156102b057600080fd5b506102b9610856565b6040516102c69190611b61565b60405180910390f35b3480156102db57600080fd5b506102f660048036038101906102f19190611a25565b61085f565b005b34801561030457600080fd5b5061031f600480360381019061031a9190611a9b565b610875565b005b34801561032d57600080fd5b5061034860048036038101906103439190611a9b565b6108d9565b6040516103559190611ad7565b60405180910390f35b34801561036a57600080fd5b50610373610921565b005b34801561038157600080fd5b5061039c60048036038101906103979190611a9b565b610935565b005b3480156103aa57600080fd5b506103b3610999565b6040516103c09190611b8b565b60405180910390f35b3480156103d557600080fd5b506103f060048036038101906103eb9190611a9b565b6109c3565b6040516103fd9190611a80565b60405180910390f35b34801561041257600080fd5b5061041b610a71565b60405161042891906119a5565b60405180910390f35b34801561043d57600080fd5b5061045860048036038101906104539190611a25565b610b03565b005b34801561046657600080fd5b50610481600480360381019061047c9190611a25565b610b19565b60405161048e9190611a80565b60405180910390f35b61049f610b3c565b005b3480156104ad57600080fd5b506104c860048036038101906104c39190611a9b565b610bf1565b005b3480156104d657600080fd5b506104f160048036038101906104ec9190611ba6565b610c55565b6040516104fe9190611ad7565b60405180910390f35b34801561051357600080fd5b5061052e60048036038101906105299190611a9b565b610cdc565b005b34801561053c57600080fd5b5061055760048036038101906105529190611a9b565b610d62565b6040516105649190611a80565b60405180910390f35b6105773382610e0f565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634f5aaaba33836040518363ffffffff1660e01b81526004016105d4929190611be6565b600060405180830381600087803b1580156105ee57600080fd5b505af1158015610602573d6000803e3d6000fd5b505050503373ffffffffffffffffffffffffffffffffffffffff167f25bedde6c8ebd3a89b719a16299dbfe271c7bffa42fe1ac1a52e15ab0cb767e68260405161064c9190611ad7565b60405180910390a250565b60606003805461066690611c3e565b80601f016020809104026020016040519081016040528092919081815260200182805461069290611c3e565b80156106df5780601f106106b4576101008083540402835291602001916106df565b820191906000526020600020905b8154815290600101906020018083116106c257829003601f168201915b5050505050905090565b6000806106f4610e91565b9050610701818585610e99565b600191505092915050565b610715336107a9565b80610725575061072433610d62565b5b610764576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161075b90611cbb565b60405180910390fd5b61076d81610eab565b50565b6000600254905090565b600080610785610e91565b9050610792858285610f3b565b61079d858585610fd0565b60019150509392505050565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016108079190611b8b565b602060405180830381865afa158015610824573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108489190611cf0565b905060028114915050919050565b60006012905090565b6108676110c4565b610871828261114b565b5050565b61087e336107a9565b8061088e575061088d33610d62565b5b6108cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108c490611cbb565b60405180910390fd5b6108d6816111cd565b50565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6109296110c4565b610933600061125d565b565b61093e336107a9565b8061094e575061094d33610d62565b5b61098d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161098490611cbb565b60405180910390fd5b61099681611323565b50565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610a219190611b8b565b602060405180830381865afa158015610a3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a629190611cf0565b90506000811415915050919050565b606060048054610a8090611c3e565b80601f0160208091040260200160405190810160405280929190818152602001828054610aac90611c3e565b8015610af95780601f10610ace57610100808354040283529160200191610af9565b820191906000526020600020905b815481529060010190602001808311610adc57829003601f168201915b5050505050905090565b610b0b6110c4565b610b158282610e0f565b5050565b600080610b24610e91565b9050610b31818585610fd0565b600191505092915050565b73010000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f19350505050158015610b96573d6000803e3d6000fd5b50610ba1333461114b565b3373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c34604051610be79190611ad7565b60405180910390a2565b610bfa336107a9565b80610c0a5750610c0933610d62565b5b610c49576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c4090611cbb565b60405180910390fd5b610c5281611421565b50565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b610ce46110c4565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610d565760006040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610d4d9190611b8b565b60405180910390fd5b610d5f8161125d565b50565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610dc09190611b8b565b602060405180830381865afa158015610ddd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e019190611cf0565b905060038114915050919050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610e815760006040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610e789190611b8b565b60405180910390fd5b610e8d826000836114b1565b5050565b600033905090565b610ea683838360016116d6565b505050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b8152600401610f069190611b8b565b600060405180830381600087803b158015610f2057600080fd5b505af1158015610f34573d6000803e3d6000fd5b5050505050565b6000610f478484610c55565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015610fca5781811015610fba578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610fb193929190611d1d565b60405180910390fd5b610fc9848484840360006116d6565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036110425760006040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016110399190611b8b565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036110b45760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016110ab9190611b8b565b60405180910390fd5b6110bf8383836114b1565b505050565b6110cc610e91565b73ffffffffffffffffffffffffffffffffffffffff166110ea610999565b73ffffffffffffffffffffffffffffffffffffffff16146111495761110d610e91565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016111409190611b8b565b60405180910390fd5b565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036111bd5760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016111b49190611b8b565b60405180910390fd5b6111c9600083836114b1565b5050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b81526004016112289190611b8b565b600060405180830381600087803b15801561124257600080fd5b505af1158015611256573d6000803e3d6000fd5b5050505050565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1603611391576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161138890611da0565b60405180910390fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b81526004016113ec9190611b8b565b600060405180830381600087803b15801561140657600080fd5b505af115801561141a573d6000803e3d6000fd5b5050505050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b815260040161147c9190611b8b565b600060405180830381600087803b15801561149657600080fd5b505af11580156114aa573d6000803e3d6000fd5b5050505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036115035780600260008282546114f79190611def565b925050819055506115d6565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561158f578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161158693929190611d1d565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361161f578060026000828254039250508190555061166c565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516116c99190611ad7565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036117485760006040517fe602df0500000000000000000000000000000000000000000000000000000000815260040161173f9190611b8b565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036117ba5760006040517f94280d620000000000000000000000000000000000000000000000000000000081526004016117b19190611b8b565b60405180910390fd5b81600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555080156118a7578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161189e9190611ad7565b60405180910390a35b50505050565b600080fd5b6000819050919050565b6118c5816118b2565b81146118d057600080fd5b50565b6000813590506118e2816118bc565b92915050565b6000602082840312156118fe576118fd6118ad565b5b600061190c848285016118d3565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561194f578082015181840152602081019050611934565b60008484015250505050565b6000601f19601f8301169050919050565b600061197782611915565b6119818185611920565b9350611991818560208601611931565b61199a8161195b565b840191505092915050565b600060208201905081810360008301526119bf818461196c565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006119f2826119c7565b9050919050565b611a02816119e7565b8114611a0d57600080fd5b50565b600081359050611a1f816119f9565b92915050565b60008060408385031215611a3c57611a3b6118ad565b5b6000611a4a85828601611a10565b9250506020611a5b858286016118d3565b9150509250929050565b60008115159050919050565b611a7a81611a65565b82525050565b6000602082019050611a956000830184611a71565b92915050565b600060208284031215611ab157611ab06118ad565b5b6000611abf84828501611a10565b91505092915050565b611ad1816118b2565b82525050565b6000602082019050611aec6000830184611ac8565b92915050565b600080600060608486031215611b0b57611b0a6118ad565b5b6000611b1986828701611a10565b9350506020611b2a86828701611a10565b9250506040611b3b868287016118d3565b9150509250925092565b600060ff82169050919050565b611b5b81611b45565b82525050565b6000602082019050611b766000830184611b52565b92915050565b611b85816119e7565b82525050565b6000602082019050611ba06000830184611b7c565b92915050565b60008060408385031215611bbd57611bbc6118ad565b5b6000611bcb85828601611a10565b9250506020611bdc85828601611a10565b9150509250929050565b6000604082019050611bfb6000830185611b7c565b611c086020830184611ac8565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680611c5657607f821691505b602082108103611c6957611c68611c0f565b5b50919050565b7f63616e6e6f74206d6f6469667920616c6c6f77206c6973740000000000000000600082015250565b6000611ca5601883611920565b9150611cb082611c6f565b602082019050919050565b60006020820190508181036000830152611cd481611c98565b9050919050565b600081519050611cea816118bc565b92915050565b600060208284031215611d0657611d056118ad565b5b6000611d1484828501611cdb565b91505092915050565b6000606082019050611d326000830186611b7c565b611d3f6020830185611ac8565b611d4c6040830184611ac8565b949350505050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000611d8a601683611920565b9150611d9582611d54565b602082019050919050565b60006020820190508181036000830152611db981611d7d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611dfa826118b2565b9150611e05836118b2565b9250828201905080821115611e1d57611e1c611dc0565b5b9291505056fea2646970667358221220cb094dd1471d39812ad1f5f9b68491f879c9d6a8e916a805b85cb7e3b7f2a6d564736f6c634300081e0033",
}

// ERC20NativeMinterABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20NativeMinterMetaData.ABI instead.
var ERC20NativeMinterABI = ERC20NativeMinterMetaData.ABI

// ERC20NativeMinterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20NativeMinterMetaData.Bin instead.
var ERC20NativeMinterBin = ERC20NativeMinterMetaData.Bin

// DeployERC20NativeMinter deploys a new Ethereum contract, binding an instance of ERC20NativeMinter to it.
func DeployERC20NativeMinter(auth *bind.TransactOpts, backend bind.ContractBackend, nativeMinterPrecompile common.Address, initSupply *big.Int) (common.Address, *types.Transaction, *ERC20NativeMinter, error) {
	parsed, err := ERC20NativeMinterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20NativeMinterBin), backend, nativeMinterPrecompile, initSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20NativeMinter{ERC20NativeMinterCaller: ERC20NativeMinterCaller{contract: contract}, ERC20NativeMinterTransactor: ERC20NativeMinterTransactor{contract: contract}, ERC20NativeMinterFilterer: ERC20NativeMinterFilterer{contract: contract}}, nil
}

// ERC20NativeMinter is an auto generated Go binding around an Ethereum contract.
type ERC20NativeMinter struct {
	ERC20NativeMinterCaller     // Read-only binding to the contract
	ERC20NativeMinterTransactor // Write-only binding to the contract
	ERC20NativeMinterFilterer   // Log filterer for contract events
}

// ERC20NativeMinterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20NativeMinterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20NativeMinterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20NativeMinterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20NativeMinterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20NativeMinterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20NativeMinterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20NativeMinterSession struct {
	Contract     *ERC20NativeMinter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ERC20NativeMinterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20NativeMinterCallerSession struct {
	Contract *ERC20NativeMinterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ERC20NativeMinterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20NativeMinterTransactorSession struct {
	Contract     *ERC20NativeMinterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ERC20NativeMinterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20NativeMinterRaw struct {
	Contract *ERC20NativeMinter // Generic contract binding to access the raw methods on
}

// ERC20NativeMinterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20NativeMinterCallerRaw struct {
	Contract *ERC20NativeMinterCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20NativeMinterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20NativeMinterTransactorRaw struct {
	Contract *ERC20NativeMinterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20NativeMinter creates a new instance of ERC20NativeMinter, bound to a specific deployed contract.
func NewERC20NativeMinter(address common.Address, backend bind.ContractBackend) (*ERC20NativeMinter, error) {
	contract, err := bindERC20NativeMinter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinter{ERC20NativeMinterCaller: ERC20NativeMinterCaller{contract: contract}, ERC20NativeMinterTransactor: ERC20NativeMinterTransactor{contract: contract}, ERC20NativeMinterFilterer: ERC20NativeMinterFilterer{contract: contract}}, nil
}

// NewERC20NativeMinterCaller creates a new read-only instance of ERC20NativeMinter, bound to a specific deployed contract.
func NewERC20NativeMinterCaller(address common.Address, caller bind.ContractCaller) (*ERC20NativeMinterCaller, error) {
	contract, err := bindERC20NativeMinter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterCaller{contract: contract}, nil
}

// NewERC20NativeMinterTransactor creates a new write-only instance of ERC20NativeMinter, bound to a specific deployed contract.
func NewERC20NativeMinterTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20NativeMinterTransactor, error) {
	contract, err := bindERC20NativeMinter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterTransactor{contract: contract}, nil
}

// NewERC20NativeMinterFilterer creates a new log filterer instance of ERC20NativeMinter, bound to a specific deployed contract.
func NewERC20NativeMinterFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20NativeMinterFilterer, error) {
	contract, err := bindERC20NativeMinter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterFilterer{contract: contract}, nil
}

// bindERC20NativeMinter binds a generic wrapper to an already deployed contract.
func bindERC20NativeMinter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20NativeMinterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20NativeMinter *ERC20NativeMinterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20NativeMinter.Contract.ERC20NativeMinterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20NativeMinter *ERC20NativeMinterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.ERC20NativeMinterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20NativeMinter *ERC20NativeMinterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.ERC20NativeMinterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20NativeMinter *ERC20NativeMinterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20NativeMinter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20NativeMinter *ERC20NativeMinterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20NativeMinter *ERC20NativeMinterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20NativeMinter.Contract.Allowance(&_ERC20NativeMinter.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20NativeMinter.Contract.Allowance(&_ERC20NativeMinter.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20NativeMinter.Contract.BalanceOf(&_ERC20NativeMinter.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20NativeMinter.Contract.BalanceOf(&_ERC20NativeMinter.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Decimals() (uint8, error) {
	return _ERC20NativeMinter.Contract.Decimals(&_ERC20NativeMinter.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) Decimals() (uint8, error) {
	return _ERC20NativeMinter.Contract.Decimals(&_ERC20NativeMinter.CallOpts)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) IsAdmin(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "isAdmin", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) IsAdmin(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsAdmin(&_ERC20NativeMinter.CallOpts, addr)
}

// IsAdmin is a free data retrieval call binding the contract method 0x24d7806c.
//
// Solidity: function isAdmin(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) IsAdmin(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsAdmin(&_ERC20NativeMinter.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) IsEnabled(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "isEnabled", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) IsEnabled(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsEnabled(&_ERC20NativeMinter.CallOpts, addr)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) IsEnabled(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsEnabled(&_ERC20NativeMinter.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) IsManager(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "isManager", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) IsManager(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsManager(&_ERC20NativeMinter.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address addr) view returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) IsManager(addr common.Address) (bool, error) {
	return _ERC20NativeMinter.Contract.IsManager(&_ERC20NativeMinter.CallOpts, addr)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Name() (string, error) {
	return _ERC20NativeMinter.Contract.Name(&_ERC20NativeMinter.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) Name() (string, error) {
	return _ERC20NativeMinter.Contract.Name(&_ERC20NativeMinter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Owner() (common.Address, error) {
	return _ERC20NativeMinter.Contract.Owner(&_ERC20NativeMinter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) Owner() (common.Address, error) {
	return _ERC20NativeMinter.Contract.Owner(&_ERC20NativeMinter.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Symbol() (string, error) {
	return _ERC20NativeMinter.Contract.Symbol(&_ERC20NativeMinter.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) Symbol() (string, error) {
	return _ERC20NativeMinter.Contract.Symbol(&_ERC20NativeMinter.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20NativeMinter.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterSession) TotalSupply() (*big.Int, error) {
	return _ERC20NativeMinter.Contract.TotalSupply(&_ERC20NativeMinter.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20NativeMinter *ERC20NativeMinterCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20NativeMinter.Contract.TotalSupply(&_ERC20NativeMinter.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Approve(&_ERC20NativeMinter.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Approve(&_ERC20NativeMinter.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "burn", from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Burn(&_ERC20NativeMinter.TransactOpts, from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Burn(&_ERC20NativeMinter.TransactOpts, from, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) Deposit() (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Deposit(&_ERC20NativeMinter.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Deposit() (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Deposit(&_ERC20NativeMinter.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Mint(&_ERC20NativeMinter.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Mint(&_ERC20NativeMinter.TransactOpts, to, amount)
}

// Mintdraw is a paid mutator transaction binding the contract method 0x0356b6cd.
//
// Solidity: function mintdraw(uint256 wad) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Mintdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "mintdraw", wad)
}

// Mintdraw is a paid mutator transaction binding the contract method 0x0356b6cd.
//
// Solidity: function mintdraw(uint256 wad) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) Mintdraw(wad *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Mintdraw(&_ERC20NativeMinter.TransactOpts, wad)
}

// Mintdraw is a paid mutator transaction binding the contract method 0x0356b6cd.
//
// Solidity: function mintdraw(uint256 wad) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Mintdraw(wad *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Mintdraw(&_ERC20NativeMinter.TransactOpts, wad)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.RenounceOwnership(&_ERC20NativeMinter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.RenounceOwnership(&_ERC20NativeMinter.TransactOpts)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Revoke(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "revoke", addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Revoke(&_ERC20NativeMinter.TransactOpts, addr)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Revoke(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Revoke(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) SetAdmin(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "setAdmin", addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetAdmin(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) SetAdmin(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetAdmin(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) SetEnabled(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "setEnabled", addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetEnabled(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetEnabled is a paid mutator transaction binding the contract method 0x0aaf7043.
//
// Solidity: function setEnabled(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) SetEnabled(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetEnabled(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) SetManager(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "setManager", addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetManager(&_ERC20NativeMinter.TransactOpts, addr)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(address addr) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) SetManager(addr common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.SetManager(&_ERC20NativeMinter.TransactOpts, addr)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Transfer(&_ERC20NativeMinter.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.Transfer(&_ERC20NativeMinter.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.TransferFrom(&_ERC20NativeMinter.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.TransferFrom(&_ERC20NativeMinter.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20NativeMinter *ERC20NativeMinterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.TransferOwnership(&_ERC20NativeMinter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20NativeMinter *ERC20NativeMinterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20NativeMinter.Contract.TransferOwnership(&_ERC20NativeMinter.TransactOpts, newOwner)
}

// ERC20NativeMinterApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20NativeMinter contract.
type ERC20NativeMinterApprovalIterator struct {
	Event *ERC20NativeMinterApproval // Event containing the contract specifics and raw log

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
func (it *ERC20NativeMinterApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20NativeMinterApproval)
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
		it.Event = new(ERC20NativeMinterApproval)
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
func (it *ERC20NativeMinterApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20NativeMinterApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20NativeMinterApproval represents a Approval event raised by the ERC20NativeMinter contract.
type ERC20NativeMinterApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20NativeMinterApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterApprovalIterator{contract: _ERC20NativeMinter.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20NativeMinterApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20NativeMinterApproval)
				if err := _ERC20NativeMinter.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) ParseApproval(log types.Log) (*ERC20NativeMinterApproval, error) {
	event := new(ERC20NativeMinterApproval)
	if err := _ERC20NativeMinter.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20NativeMinterDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the ERC20NativeMinter contract.
type ERC20NativeMinterDepositIterator struct {
	Event *ERC20NativeMinterDeposit // Event containing the contract specifics and raw log

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
func (it *ERC20NativeMinterDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20NativeMinterDeposit)
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
		it.Event = new(ERC20NativeMinterDeposit)
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
func (it *ERC20NativeMinterDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20NativeMinterDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20NativeMinterDeposit represents a Deposit event raised by the ERC20NativeMinter contract.
type ERC20NativeMinterDeposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*ERC20NativeMinterDepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterDepositIterator{contract: _ERC20NativeMinter.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *ERC20NativeMinterDeposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20NativeMinterDeposit)
				if err := _ERC20NativeMinter.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) ParseDeposit(log types.Log) (*ERC20NativeMinterDeposit, error) {
	event := new(ERC20NativeMinterDeposit)
	if err := _ERC20NativeMinter.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20NativeMinterMintdrawalIterator is returned from FilterMintdrawal and is used to iterate over the raw logs and unpacked data for Mintdrawal events raised by the ERC20NativeMinter contract.
type ERC20NativeMinterMintdrawalIterator struct {
	Event *ERC20NativeMinterMintdrawal // Event containing the contract specifics and raw log

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
func (it *ERC20NativeMinterMintdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20NativeMinterMintdrawal)
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
		it.Event = new(ERC20NativeMinterMintdrawal)
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
func (it *ERC20NativeMinterMintdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20NativeMinterMintdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20NativeMinterMintdrawal represents a Mintdrawal event raised by the ERC20NativeMinter contract.
type ERC20NativeMinterMintdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMintdrawal is a free log retrieval operation binding the contract event 0x25bedde6c8ebd3a89b719a16299dbfe271c7bffa42fe1ac1a52e15ab0cb767e6.
//
// Solidity: event Mintdrawal(address indexed src, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) FilterMintdrawal(opts *bind.FilterOpts, src []common.Address) (*ERC20NativeMinterMintdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.FilterLogs(opts, "Mintdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterMintdrawalIterator{contract: _ERC20NativeMinter.contract, event: "Mintdrawal", logs: logs, sub: sub}, nil
}

// WatchMintdrawal is a free log subscription operation binding the contract event 0x25bedde6c8ebd3a89b719a16299dbfe271c7bffa42fe1ac1a52e15ab0cb767e6.
//
// Solidity: event Mintdrawal(address indexed src, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) WatchMintdrawal(opts *bind.WatchOpts, sink chan<- *ERC20NativeMinterMintdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.WatchLogs(opts, "Mintdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20NativeMinterMintdrawal)
				if err := _ERC20NativeMinter.contract.UnpackLog(event, "Mintdrawal", log); err != nil {
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

// ParseMintdrawal is a log parse operation binding the contract event 0x25bedde6c8ebd3a89b719a16299dbfe271c7bffa42fe1ac1a52e15ab0cb767e6.
//
// Solidity: event Mintdrawal(address indexed src, uint256 wad)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) ParseMintdrawal(log types.Log) (*ERC20NativeMinterMintdrawal, error) {
	event := new(ERC20NativeMinterMintdrawal)
	if err := _ERC20NativeMinter.contract.UnpackLog(event, "Mintdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20NativeMinterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ERC20NativeMinter contract.
type ERC20NativeMinterOwnershipTransferredIterator struct {
	Event *ERC20NativeMinterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ERC20NativeMinterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20NativeMinterOwnershipTransferred)
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
		it.Event = new(ERC20NativeMinterOwnershipTransferred)
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
func (it *ERC20NativeMinterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20NativeMinterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20NativeMinterOwnershipTransferred represents a OwnershipTransferred event raised by the ERC20NativeMinter contract.
type ERC20NativeMinterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ERC20NativeMinterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterOwnershipTransferredIterator{contract: _ERC20NativeMinter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ERC20NativeMinterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20NativeMinterOwnershipTransferred)
				if err := _ERC20NativeMinter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) ParseOwnershipTransferred(log types.Log) (*ERC20NativeMinterOwnershipTransferred, error) {
	event := new(ERC20NativeMinterOwnershipTransferred)
	if err := _ERC20NativeMinter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20NativeMinterTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20NativeMinter contract.
type ERC20NativeMinterTransferIterator struct {
	Event *ERC20NativeMinterTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20NativeMinterTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20NativeMinterTransfer)
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
		it.Event = new(ERC20NativeMinterTransfer)
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
func (it *ERC20NativeMinterTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20NativeMinterTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20NativeMinterTransfer represents a Transfer event raised by the ERC20NativeMinter contract.
type ERC20NativeMinterTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20NativeMinterTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20NativeMinterTransferIterator{contract: _ERC20NativeMinter.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20NativeMinterTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20NativeMinter.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20NativeMinterTransfer)
				if err := _ERC20NativeMinter.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20NativeMinter *ERC20NativeMinterFilterer) ParseTransfer(log types.Log) (*ERC20NativeMinterTransfer, error) {
	event := new(ERC20NativeMinterTransfer)
	if err := _ERC20NativeMinter.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
