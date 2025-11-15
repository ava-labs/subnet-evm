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

// ERC20NativeMinterMetaData contains all meta data concerning the ERC20NativeMinter contract.
var ERC20NativeMinterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Mintdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"mintdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"revoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052730200000000000000000000000000000000000001600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006557600080fd5b5060405161270c38038061270c833981810160405281019061008791906105c8565b730200000000000000000000000000000000000001336040518060400160405280601681526020017f45524332304e61746976654d696e746572546f6b656e000000000000000000008152506040518060400160405280600481526020017f584d504c000000000000000000000000000000000000000000000000000000008152508160039081610118919061083b565b508060049081610128919061083b565b505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361019d5760006040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610194919061094e565b60405180910390fd5b6101ac8161021260201b60201c565b5080600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061020c6102006102d860201b60201c565b826102e060201b60201c565b50610a2d565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036103525760006040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610349919061094e565b60405180910390fd5b6103646000838361036860201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036103ba5780600260008282546103ae9190610998565b9250508190555061048d565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610446578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161043d939291906109db565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036104d65780600260008282540392505081905550610523565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105809190610a12565b60405180910390a3505050565b600080fd5b6000819050919050565b6105a581610592565b81146105b057600080fd5b50565b6000815190506105c28161059c565b92915050565b6000602082840312156105de576105dd61058d565b5b60006105ec848285016105b3565b91505092915050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061067657607f821691505b6020821081036106895761068861062f565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026106f17fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826106b4565b6106fb86836106b4565b95508019841693508086168417925050509392505050565b6000819050919050565b600061073861073361072e84610592565b610713565b610592565b9050919050565b6000819050919050565b6107528361071d565b61076661075e8261073f565b8484546106c1565b825550505050565b600090565b61077b61076e565b610786818484610749565b505050565b5b818110156107aa5761079f600082610773565b60018101905061078c565b5050565b601f8211156107ef576107c08161068f565b6107c9846106a4565b810160208510156107d8578190505b6107ec6107e4856106a4565b83018261078b565b50505b505050565b600082821c905092915050565b6000610812600019846008026107f4565b1980831691505092915050565b600061082b8383610801565b9150826002028217905092915050565b610844826105f5565b67ffffffffffffffff81111561085d5761085c610600565b5b610867825461065e565b6108728282856107ae565b600060209050601f8311600181146108a55760008415610893578287015190505b61089d858261081f565b865550610905565b601f1984166108b38661068f565b60005b828110156108db578489015182556001820191506020850194506020810190506108b6565b868310156108f857848901516108f4601f891682610801565b8355505b6001600288020188555050505b505050505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006109388261090d565b9050919050565b6109488161092d565b82525050565b6000602082019050610963600083018461093f565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006109a382610592565b91506109ae83610592565b92508282019050808211156109c6576109c5610969565b5b92915050565b6109d581610592565b82525050565b60006060820190506109f0600083018661093f565b6109fd60208301856109cc565b610a0a60408301846109cc565b949350505050565b6000602082019050610a2760008301846109cc565b92915050565b611cd080610a3c6000396000f3fe6080604052600436106101405760003560e01c8063715018a6116100b6578063a9059cbb1161006f578063a9059cbb1461045a578063d0e30db014610497578063d0ebdbe7146104a1578063dd62ed3e146104ca578063f2fde38b14610507578063f3ae24151461053057610140565b8063715018a61461035e57806374a8f103146103755780638da5cb5b1461039e5780639015d371146103c957806395d89b41146104065780639dc29fac1461043157610140565b806323b872dd1161010857806323b872dd1461022a57806324d7806c14610267578063313ce567146102a457806340c10f19146102cf578063704b6c02146102f857806370a082311461032157610140565b80630356b6cd1461014557806306fdde031461016e578063095ea7b3146101995780630aaf7043146101d657806318160ddd146101ff575b600080fd5b34801561015157600080fd5b5061016c600480360381019061016791906117cb565b61056d565b005b34801561017a57600080fd5b5061018361066c565b6040516101909190611888565b60405180910390f35b3480156101a557600080fd5b506101c060048036038101906101bb9190611908565b6106fe565b6040516101cd9190611963565b60405180910390f35b3480156101e257600080fd5b506101fd60048036038101906101f8919061197e565b610721565b005b34801561020b57600080fd5b50610214610735565b60405161022191906119ba565b60405180910390f35b34801561023657600080fd5b50610251600480360381019061024c91906119d5565b61073f565b60405161025e9190611963565b60405180910390f35b34801561027357600080fd5b5061028e6004803603810190610289919061197e565b61076e565b60405161029b9190611963565b60405180910390f35b3480156102b057600080fd5b506102b961081b565b6040516102c69190611a44565b60405180910390f35b3480156102db57600080fd5b506102f660048036038101906102f19190611908565b610824565b005b34801561030457600080fd5b5061031f600480360381019061031a919061197e565b61083a565b005b34801561032d57600080fd5b506103486004803603810190610343919061197e565b61084e565b60405161035591906119ba565b60405180910390f35b34801561036a57600080fd5b50610373610896565b005b34801561038157600080fd5b5061039c6004803603810190610397919061197e565b6108aa565b005b3480156103aa57600080fd5b506103b36108be565b6040516103c09190611a6e565b60405180910390f35b3480156103d557600080fd5b506103f060048036038101906103eb919061197e565b6108e8565b6040516103fd9190611963565b60405180910390f35b34801561041257600080fd5b5061041b610996565b6040516104289190611888565b60405180910390f35b34801561043d57600080fd5b5061045860048036038101906104539190611908565b610a28565b005b34801561046657600080fd5b50610481600480360381019061047c9190611908565b610a3e565b60405161048e9190611963565b60405180910390f35b61049f610a61565b005b3480156104ad57600080fd5b506104c860048036038101906104c3919061197e565b610b24565b005b3480156104d657600080fd5b506104f160048036038101906104ec9190611a89565b610b38565b6040516104fe91906119ba565b60405180910390f35b34801561051357600080fd5b5061052e6004803603810190610529919061197e565b610bbf565b005b34801561053c57600080fd5b506105576004803603810190610552919061197e565b610c45565b6040516105649190611963565b60405180910390f35b61057e610578610cf2565b82610cfa565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634f5aaaba6105c4610cf2565b836040518363ffffffff1660e01b81526004016105e2929190611ac9565b600060405180830381600087803b1580156105fc57600080fd5b505af1158015610610573d6000803e3d6000fd5b5050505061061c610cf2565b73ffffffffffffffffffffffffffffffffffffffff167f25bedde6c8ebd3a89b719a16299dbfe271c7bffa42fe1ac1a52e15ab0cb767e68260405161066191906119ba565b60405180910390a250565b60606003805461067b90611b21565b80601f01602080910402602001604051908101604052809291908181526020018280546106a790611b21565b80156106f45780601f106106c9576101008083540402835291602001916106f4565b820191906000526020600020905b8154815290600101906020018083116106d757829003601f168201915b5050505050905090565b600080610709610cf2565b9050610716818585610d7c565b600191505092915050565b610729610d8e565b61073281610e15565b50565b6000600254905090565b60008061074a610cf2565b9050610757858285610ea5565b610762858585610f3a565b60019150509392505050565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016107cc9190611a6e565b602060405180830381865afa1580156107e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061080d9190611b67565b905060028114915050919050565b60006012905090565b61082c610d8e565b610836828261102e565b5050565b610842610d8e565b61084b816110b0565b50565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b61089e610d8e565b6108a86000611140565b565b6108b2610d8e565b6108bb81611206565b50565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016109469190611a6e565b602060405180830381865afa158015610963573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109879190611b67565b90506000811415915050919050565b6060600480546109a590611b21565b80601f01602080910402602001604051908101604052809291908181526020018280546109d190611b21565b8015610a1e5780601f106109f357610100808354040283529160200191610a1e565b820191906000526020600020905b815481529060010190602001808311610a0157829003601f168201915b5050505050905090565b610a30610d8e565b610a3a8282610cfa565b5050565b600080610a49610cf2565b9050610a56818585610f3a565b600191505092915050565b73010000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166108fc349081150290604051600060405180830381858888f19350505050158015610abb573d6000803e3d6000fd5b50610acd610ac7610cf2565b3461102e565b610ad5610cf2565b73ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c34604051610b1a91906119ba565b60405180910390a2565b610b2c610d8e565b610b3581611304565b50565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b610bc7610d8e565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610c395760006040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610c309190611a6e565b60405180910390fd5b610c4281611140565b50565b600080600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610ca39190611a6e565b602060405180830381865afa158015610cc0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ce49190611b67565b905060038114915050919050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610d6c5760006040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610d639190611a6e565b60405180910390fd5b610d7882600083611394565b5050565b610d8983838360016115b9565b505050565b610d96610cf2565b73ffffffffffffffffffffffffffffffffffffffff16610db46108be565b73ffffffffffffffffffffffffffffffffffffffff1614610e1357610dd7610cf2565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401610e0a9190611a6e565b60405180910390fd5b565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b8152600401610e709190611a6e565b600060405180830381600087803b158015610e8a57600080fd5b505af1158015610e9e573d6000803e3d6000fd5b5050505050565b6000610eb18484610b38565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015610f345781811015610f24578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610f1b93929190611b94565b60405180910390fd5b610f33848484840360006115b9565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610fac5760006040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610fa39190611a6e565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361101e5760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016110159190611a6e565b60405180910390fd5b611029838383611394565b505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036110a05760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016110979190611a6e565b60405180910390fd5b6110ac60008383611394565b5050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b815260040161110b9190611a6e565b600060405180830381600087803b15801561112557600080fd5b505af1158015611139573d6000803e3d6000fd5b5050505050565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1603611274576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161126b90611c17565b60405180910390fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b81526004016112cf9190611a6e565b600060405180830381600087803b1580156112e957600080fd5b505af11580156112fd573d6000803e3d6000fd5b5050505050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b815260040161135f9190611a6e565b600060405180830381600087803b15801561137957600080fd5b505af115801561138d573d6000803e3d6000fd5b5050505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036113e65780600260008282546113da9190611c66565b925050819055506114b9565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015611472578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161146993929190611b94565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611502578060026000828254039250508190555061154f565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516115ac91906119ba565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff160361162b5760006040517fe602df050000000000000000000000000000000000000000000000000000000081526004016116229190611a6e565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361169d5760006040517f94280d620000000000000000000000000000000000000000000000000000000081526004016116949190611a6e565b60405180910390fd5b81600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550801561178a578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161178191906119ba565b60405180910390a35b50505050565b600080fd5b6000819050919050565b6117a881611795565b81146117b357600080fd5b50565b6000813590506117c58161179f565b92915050565b6000602082840312156117e1576117e0611790565b5b60006117ef848285016117b6565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015611832578082015181840152602081019050611817565b60008484015250505050565b6000601f19601f8301169050919050565b600061185a826117f8565b6118648185611803565b9350611874818560208601611814565b61187d8161183e565b840191505092915050565b600060208201905081810360008301526118a2818461184f565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006118d5826118aa565b9050919050565b6118e5816118ca565b81146118f057600080fd5b50565b600081359050611902816118dc565b92915050565b6000806040838503121561191f5761191e611790565b5b600061192d858286016118f3565b925050602061193e858286016117b6565b9150509250929050565b60008115159050919050565b61195d81611948565b82525050565b60006020820190506119786000830184611954565b92915050565b60006020828403121561199457611993611790565b5b60006119a2848285016118f3565b91505092915050565b6119b481611795565b82525050565b60006020820190506119cf60008301846119ab565b92915050565b6000806000606084860312156119ee576119ed611790565b5b60006119fc868287016118f3565b9350506020611a0d868287016118f3565b9250506040611a1e868287016117b6565b9150509250925092565b600060ff82169050919050565b611a3e81611a28565b82525050565b6000602082019050611a596000830184611a35565b92915050565b611a68816118ca565b82525050565b6000602082019050611a836000830184611a5f565b92915050565b60008060408385031215611aa057611a9f611790565b5b6000611aae858286016118f3565b9250506020611abf858286016118f3565b9150509250929050565b6000604082019050611ade6000830185611a5f565b611aeb60208301846119ab565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680611b3957607f821691505b602082108103611b4c57611b4b611af2565b5b50919050565b600081519050611b618161179f565b92915050565b600060208284031215611b7d57611b7c611790565b5b6000611b8b84828501611b52565b91505092915050565b6000606082019050611ba96000830186611a5f565b611bb660208301856119ab565b611bc360408301846119ab565b949350505050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000611c01601683611803565b9150611c0c82611bcb565b602082019050919050565b60006020820190508181036000830152611c3081611bf4565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611c7182611795565b9150611c7c83611795565b9250828201905080821115611c9457611c93611c37565b5b9291505056fea26469706673582212205b03965b7e7a93066f2f299f420eabaf901243558086de2b6d6882835d4bb40264736f6c634300081e0033",
}

// ERC20NativeMinterABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20NativeMinterMetaData.ABI instead.
var ERC20NativeMinterABI = ERC20NativeMinterMetaData.ABI

// ERC20NativeMinterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20NativeMinterMetaData.Bin instead.
var ERC20NativeMinterBin = ERC20NativeMinterMetaData.Bin

// DeployERC20NativeMinter deploys a new Ethereum contract, binding an instance of ERC20NativeMinter to it.
func DeployERC20NativeMinter(auth *bind.TransactOpts, backend bind.ContractBackend, initSupply *big.Int) (common.Address, *types.Transaction, *ERC20NativeMinter, error) {
	parsed, err := ERC20NativeMinterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20NativeMinterBin), backend, initSupply)
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
