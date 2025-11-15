// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deployerallowlisttest

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

// DeployerListTestMetaData contains all meta data concerning the DeployerListTest contract.
var DeployerListTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"log_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"log_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"log_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"name\":\"log_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"val\",\"type\":\"address\"}],\"name\":\"log_named_address\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"val\",\"type\":\"bytes\"}],\"name\":\"log_named_bytes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"val\",\"type\":\"bytes32\"}],\"name\":\"log_named_bytes32\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"log_named_decimal_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"val\",\"type\":\"int256\"}],\"name\":\"log_named_int\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"val\",\"type\":\"string\"}],\"name\":\"log_named_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"log_named_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"log_string\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"log_uint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"logs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"IS_TEST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"failed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setUp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_addDeployerThroughContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_adminAddContractAsAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_adminCanRevokeDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_deployerCanDeploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_newAddressHasNoRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_noRoleCannotDeploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_noRoleIsNotAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"step_verifySenderIsAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260016000806101000a81548160ff021916908315150217905550730200000000000000000000000000000000000000600060026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015607e57600080fd5b50612c568061008e6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80634c37cc2e116100715780634c37cc2e146100e0578063712268f1146100ea578063ba414fa6146100f4578063f26c562c14610112578063fa7626d41461011c578063ffc46bc21461013a576100a9565b80630a9254e4146100ae5780631d44d2da146100b857806328002804146100c25780632999b249146100cc57806333cb47db146100d6575b600080fd5b6100b6610144565b005b6100c061026f565b005b6100ca6103d8565b005b6100d4610480565b005b6100de610807565b005b6100e86108d7565b005b6100f2610cbe565b005b6100fc611178565b6040516101099190611a6d565b60405180910390f35b61011a611315565b005b61012461148a565b6040516101319190611a6d565b60405180910390f35b61014261149b565b005b73020000000000000000000000000000000000000060405161016590611a45565b61016f9190611ac9565b604051809103906000f08015801561018b573d6000803e3d6000fd5b50600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b730fa8ea536be85f32724d57a37758761b864161236040518263ffffffff1660e01b815260040161023b9190611ac9565b600060405180830381600087803b15801561025557600080fd5b505af1158015610269573d6000803e3d6000fd5b50505050565b610315600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1326040518263ffffffff1660e01b81526004016102cd9190611ac9565b602060405180830381865afa1580156102ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061030e9190611b1f565b6000611742565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636cd5c39b6040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561037f57600080fd5b505af1925050508015610390575060015b156103d6576103d560006040518060400160405280601a81526020017f6465706c6f79436f6e74726163742073686f756c64206661696c000000000000815250611762565b5b565b61047e600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1336040518263ffffffff1660e01b81526004016104369190611ac9565b602060405180830381865afa158015610453573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104779190611b1f565b6002611742565b565b60007302000000000000000000000000000000000000006040516104a390611a45565b6104ad9190611ac9565b604051809103906000f0801580156104c9573d6000803e3d6000fd5b5090506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600082905061059e600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016105569190611ac9565b602060405180830381865afa158015610573573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105979190611b1f565b6000611742565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02836040518263ffffffff1660e01b81526004016105f99190611ac9565b600060405180830381600087803b15801561061357600080fd5b505af1158015610627573d6000803e3d6000fd5b505050506106d1600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016106899190611ac9565b602060405180830381865afa1580156106a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106ca9190611b1f565b6002611742565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b815260040161072c9190611ac9565b600060405180830381600087803b15801561074657600080fd5b505af115801561075a573d6000803e3d6000fd5b50505050610802600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639015d371836040518263ffffffff1660e01b81526004016107bc9190611ac9565b602060405180830381865afa1580156107d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107fd9190611b78565b6117ac565b505050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506108d4600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1836040518263ffffffff1660e01b815260040161088c9190611ac9565b602060405180830381865afa1580156108a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108cd9190611b1f565b6000611742565b50565b60007302000000000000000000000000000000000000006040516108fa90611a45565b6109049190611ac9565b604051809103906000f080158015610920573d6000803e3d6000fd5b5090506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008290506109f5600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016109ad9190611ac9565b602060405180830381865afa1580156109ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109ee9190611b1f565b6000611742565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02836040518263ffffffff1660e01b8152600401610a509190611ac9565b600060405180830381600087803b158015610a6a57600080fd5b505af1158015610a7e573d6000803e3d6000fd5b50505050610b28600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610ae09190611ac9565b602060405180830381865afa158015610afd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b219190611b1f565b6002611742565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b8152600401610b839190611ac9565b600060405180830381600087803b158015610b9d57600080fd5b505af1158015610bb1573d6000803e3d6000fd5b50505050610c59600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639015d371836040518263ffffffff1660e01b8152600401610c139190611ac9565b602060405180830381865afa158015610c30573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c549190611b78565b6117ac565b8273ffffffffffffffffffffffffffffffffffffffff16636cd5c39b6040518163ffffffff1660e01b8152600401600060405180830381600087803b158015610ca157600080fd5b505af1158015610cb5573d6000803e3d6000fd5b50505050505050565b6000730200000000000000000000000000000000000000604051610ce190611a45565b610ceb9190611ac9565b604051809103906000f080158015610d07573d6000803e3d6000fd5b5090506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000829050610ddc600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610d949190611ac9565b602060405180830381865afa158015610db1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dd59190611b1f565b6000611742565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02836040518263ffffffff1660e01b8152600401610e379190611ac9565b600060405180830381600087803b158015610e5157600080fd5b505af1158015610e65573d6000803e3d6000fd5b50505050610f0f600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b8152600401610ec79190611ac9565b602060405180830381865afa158015610ee4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f089190611b1f565b6002611742565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b8152600401610f6a9190611ac9565b600060405180830381600087803b158015610f8457600080fd5b505af1158015610f98573d6000803e3d6000fd5b50505050611040600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639015d371836040518263ffffffff1660e01b8152600401610ffa9190611ac9565b602060405180830381865afa158015611017573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061103b9190611b78565b6117ac565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166374a8f103826040518263ffffffff1660e01b815260040161109b9190611ac9565b600060405180830381600087803b1580156110b557600080fd5b505af11580156110c9573d6000803e3d6000fd5b50505050611173600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1836040518263ffffffff1660e01b815260040161112b9190611ac9565b602060405180830381865afa158015611148573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061116c9190611b1f565b6000611742565b505050565b60008060019054906101000a900460ff16156111a557600060019054906101000a900460ff169050611312565b60006111af6117f2565b1561130d5760007f885cb69240a935d632d79c317109709ecfa91a80626ff3989d68f67f5b1dd12d60001c60601b60601c73ffffffffffffffffffffffffffffffffffffffff167f667f9d70ca411d70ead50d8d5c22070dafc36ad75f3dcf5e7237b22ade9aecc47f885cb69240a935d632d79c317109709ecfa91a80626ff3989d68f67f5b1dd12d60001c60601b60601c7f6661696c65640000000000000000000000000000000000000000000000000000604051602001611273929190611bbe565b604051602081830303815290604052604051602001611293929190611ca5565b6040516020818303038152906040526040516112af9190611ccd565b6000604051808303816000865af19150503d80600081146112ec576040519150601f19603f3d011682016040523d82523d6000602084013e6112f1565b606091505b50915050808060200190518101906113099190611b78565b9150505b809150505b90565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506113e2600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1836040518263ffffffff1660e01b815260040161139a9190611ac9565b602060405180830381865afa1580156113b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113db9190611b1f565b6000611742565b611487600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166324d7806c836040518263ffffffff1660e01b81526004016114409190611ac9565b602060405180830381865afa15801561145d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114819190611b78565b156117ac565b50565b60008054906101000a900460ff1681565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050611568600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1836040518263ffffffff1660e01b81526004016115209190611ac9565b602060405180830381865afa15801561153d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115619190611b1f565b6000611742565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b81526004016115c39190611ac9565b600060405180830381600087803b1580156115dd57600080fd5b505af11580156115f1573d6000803e3d6000fd5b5050505061169b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1836040518263ffffffff1660e01b81526004016116539190611ac9565b602060405180830381865afa158015611670573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116949190611b1f565b6002611742565b61173f600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166324d7806c836040518263ffffffff1660e01b81526004016116f99190611ac9565b602060405180830381865afa158015611716573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061173a9190611b78565b6117ac565b50565b61175e8282600381111561175957611758611ce4565b5b61181b565b5050565b816117a8577f280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583816040516117969190611dc5565b60405180910390a16117a7826117ac565b5b5050565b806117ef577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f506040516117de90611e46565b60405180910390a16117ee6118d2565b5b50565b60008060009050737109709ecfa91a80626ff3989d68f67f5b1dd12d3b90506000811191505090565b8082146118ce577f41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f5060405161184f90611ed8565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8826040516118869190611f53565b60405180910390a17fb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8816040516118bd9190611fcd565b60405180910390a16118cd6118d2565b5b5050565b6118da6117f2565b15611a285760007f885cb69240a935d632d79c317109709ecfa91a80626ff3989d68f67f5b1dd12d60001c60601b60601c73ffffffffffffffffffffffffffffffffffffffff167f70ca10bbd0dbfd9020a9f4b13402c16cb120705e0d1c0aeab10fa353ae586fc47f885cb69240a935d632d79c317109709ecfa91a80626ff3989d68f67f5b1dd12d60001c60601b60601c7f6661696c65640000000000000000000000000000000000000000000000000000600160001b6040516020016119a493929190611ffb565b6040516020818303038152906040526040516020016119c4929190611ca5565b6040516020818303038152906040526040516119e09190611ccd565b6000604051808303816000865af19150503d8060008114611a1d576040519150601f19603f3d011682016040523d82523d6000602084013e611a22565b606091505b50509050505b6001600060016101000a81548160ff021916908315150217905550565b610bee8061203383390190565b60008115159050919050565b611a6781611a52565b82525050565b6000602082019050611a826000830184611a5e565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611ab382611a88565b9050919050565b611ac381611aa8565b82525050565b6000602082019050611ade6000830184611aba565b92915050565b600080fd5b6000819050919050565b611afc81611ae9565b8114611b0757600080fd5b50565b600081519050611b1981611af3565b92915050565b600060208284031215611b3557611b34611ae4565b5b6000611b4384828501611b0a565b91505092915050565b611b5581611a52565b8114611b6057600080fd5b50565b600081519050611b7281611b4c565b92915050565b600060208284031215611b8e57611b8d611ae4565b5b6000611b9c84828501611b63565b91505092915050565b6000819050919050565b611bb881611ba5565b82525050565b6000604082019050611bd36000830185611aba565b611be06020830184611baf565b9392505050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b611c2e611c2982611be7565b611c13565b82525050565b600081519050919050565b600081905092915050565b60005b83811015611c68578082015181840152602081019050611c4d565b60008484015250505050565b6000611c7f82611c34565b611c898185611c3f565b9350611c99818560208601611c4a565b80840191505092915050565b6000611cb18285611c1d565b600482019150611cc18284611c74565b91508190509392505050565b6000611cd98284611c74565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600082825260208201905092915050565b7f4572726f72000000000000000000000000000000000000000000000000000000600082015250565b6000611d5a600583611d13565b9150611d6582611d24565b602082019050919050565b600081519050919050565b6000601f19601f8301169050919050565b6000611d9782611d70565b611da18185611d13565b9350611db1818560208601611c4a565b611dba81611d7b565b840191505092915050565b60006040820190508181036000830152611dde81611d4d565b90508181036020830152611df28184611d8c565b905092915050565b7f4572726f723a20417373657274696f6e204661696c6564000000000000000000600082015250565b6000611e30601783611d13565b9150611e3b82611dfa565b602082019050919050565b60006020820190508181036000830152611e5f81611e23565b9050919050565b7f4572726f723a2061203d3d2062206e6f7420736174697366696564205b75696e60008201527f745d000000000000000000000000000000000000000000000000000000000000602082015250565b6000611ec2602283611d13565b9150611ecd82611e66565b604082019050919050565b60006020820190508181036000830152611ef181611eb5565b9050919050565b7f2020202020204c65667400000000000000000000000000000000000000000000600082015250565b6000611f2e600a83611d13565b9150611f3982611ef8565b602082019050919050565b611f4d81611ae9565b82525050565b60006040820190508181036000830152611f6c81611f21565b9050611f7b6020830184611f44565b92915050565b7f2020202020526967687400000000000000000000000000000000000000000000600082015250565b6000611fb7600a83611d13565b9150611fc282611f81565b602082019050919050565b60006040820190508181036000830152611fe681611faa565b9050611ff56020830184611f44565b92915050565b60006060820190506120106000830186611aba565b61201d6020830185611baf565b61202a6040830184611baf565b94935050505056fe608060405234801561001057600080fd5b50604051610bee380380610bee833981810160405281019061003291906100dd565b80806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505061010a565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100aa8261007f565b9050919050565b6100ba8161009f565b81146100c557600080fd5b50565b6000815190506100d7816100b1565b92915050565b6000602082840312156100f3576100f261007a565b5b6000610101848285016100c8565b91505092915050565b610ad5806101196000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806374a8f1031161005b57806374a8f103146100ff5780639015d3711461011b578063d0ebdbe71461014b578063f3ae24151461016757610088565b80630aaf70431461008d57806324d7806c146100a95780636cd5c39b146100d9578063704b6c02146100e3575b600080fd5b6100a760048036038101906100a2919061086a565b610197565b005b6100c360048036038101906100be919061086a565b6101fb565b6040516100d091906108b2565b60405180910390f35b6100e16102a6565b005b6100fd60048036038101906100f8919061086a565b6102d2565b005b6101196004803603810190610114919061086a565b610336565b005b6101356004803603810190610130919061086a565b61039a565b60405161014291906108b2565b60405180910390f35b6101656004803603810190610160919061086a565b610446565b005b610181600480360381019061017c919061086a565b6104aa565b60405161018e91906108b2565b60405180910390f35b6101a0336101fb565b806101b057506101af336104aa565b5b6101ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101e69061092a565b60405180910390fd5b6101f881610555565b50565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016102579190610959565b602060405180830381865afa158015610274573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029891906109aa565b905060028114915050919050565b6040516102b2906107fb565b604051809103906000f0801580156102ce573d6000803e3d6000fd5b5050565b6102db336101fb565b806102eb57506102ea336104aa565b5b61032a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103219061092a565b60405180910390fd5b610333816105e3565b50565b61033f336101fb565b8061034f575061034e336104aa565b5b61038e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103859061092a565b60405180910390fd5b61039781610671565b50565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016103f69190610959565b602060405180830381865afa158015610413573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061043791906109aa565b90506000811415915050919050565b61044f336101fb565b8061045f575061045e336104aa565b5b61049e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104959061092a565b60405180910390fd5b6104a78161076d565b50565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eb54dae1846040518263ffffffff1660e01b81526004016105069190610959565b602060405180830381865afa158015610523573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061054791906109aa565b905060038114915050919050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630aaf7043826040518263ffffffff1660e01b81526004016105ae9190610959565b600060405180830381600087803b1580156105c857600080fd5b505af11580156105dc573d6000803e3d6000fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663704b6c02826040518263ffffffff1660e01b815260040161063c9190610959565b600060405180830381600087803b15801561065657600080fd5b505af115801561066a573d6000803e3d6000fd5b5050505050565b8073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16036106df576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106d690610a23565b60405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638c6bfb3b826040518263ffffffff1660e01b81526004016107389190610959565b600060405180830381600087803b15801561075257600080fd5b505af1158015610766573d6000803e3d6000fd5b5050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0ebdbe7826040518263ffffffff1660e01b81526004016107c69190610959565b600060405180830381600087803b1580156107e057600080fd5b505af11580156107f4573d6000803e3d6000fd5b5050505050565b605c80610a4483390190565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006108378261080c565b9050919050565b6108478161082c565b811461085257600080fd5b50565b6000813590506108648161083e565b92915050565b6000602082840312156108805761087f610807565b5b600061088e84828501610855565b91505092915050565b60008115159050919050565b6108ac81610897565b82525050565b60006020820190506108c760008301846108a3565b92915050565b600082825260208201905092915050565b7f63616e6e6f74206d6f6469667920616c6c6f77206c6973740000000000000000600082015250565b60006109146018836108cd565b915061091f826108de565b602082019050919050565b6000602082019050818103600083015261094381610907565b9050919050565b6109538161082c565b82525050565b600060208201905061096e600083018461094a565b92915050565b6000819050919050565b61098781610974565b811461099257600080fd5b50565b6000815190506109a48161097e565b92915050565b6000602082840312156109c0576109bf610807565b5b60006109ce84828501610995565b91505092915050565b7f63616e6e6f74207265766f6b65206f776e20726f6c6500000000000000000000600082015250565b6000610a0d6016836108cd565b9150610a18826109d7565b602082019050919050565b60006020820190508181036000830152610a3c81610a00565b905091905056fe6080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea2646970667358221220c801cdd35636508ed1db0208fd8ce8693222b619053b5e7ddef745b8de38516864736f6c634300081e0033a2646970667358221220dd6ccf85cadd84da4adb9358bdb8ea1cb772d271d7fb8863a428d9af48242fed64736f6c634300081e0033a26469706673582212204b5525f0bc3c1c08c72dcb5b1f7c4b7408bc24f2c09d175b66ebefa6d4a3886b64736f6c634300081e0033",
}

// DeployerListTestABI is the input ABI used to generate the binding from.
// Deprecated: Use DeployerListTestMetaData.ABI instead.
var DeployerListTestABI = DeployerListTestMetaData.ABI

// DeployerListTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeployerListTestMetaData.Bin instead.
var DeployerListTestBin = DeployerListTestMetaData.Bin

// DeployDeployerListTest deploys a new Ethereum contract, binding an instance of DeployerListTest to it.
func DeployDeployerListTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DeployerListTest, error) {
	parsed, err := DeployerListTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeployerListTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DeployerListTest{DeployerListTestCaller: DeployerListTestCaller{contract: contract}, DeployerListTestTransactor: DeployerListTestTransactor{contract: contract}, DeployerListTestFilterer: DeployerListTestFilterer{contract: contract}}, nil
}

// DeployerListTest is an auto generated Go binding around an Ethereum contract.
type DeployerListTest struct {
	DeployerListTestCaller     // Read-only binding to the contract
	DeployerListTestTransactor // Write-only binding to the contract
	DeployerListTestFilterer   // Log filterer for contract events
}

// DeployerListTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type DeployerListTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerListTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DeployerListTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerListTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeployerListTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerListTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeployerListTestSession struct {
	Contract     *DeployerListTest // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DeployerListTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeployerListTestCallerSession struct {
	Contract *DeployerListTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DeployerListTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeployerListTestTransactorSession struct {
	Contract     *DeployerListTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DeployerListTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type DeployerListTestRaw struct {
	Contract *DeployerListTest // Generic contract binding to access the raw methods on
}

// DeployerListTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeployerListTestCallerRaw struct {
	Contract *DeployerListTestCaller // Generic read-only contract binding to access the raw methods on
}

// DeployerListTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeployerListTestTransactorRaw struct {
	Contract *DeployerListTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeployerListTest creates a new instance of DeployerListTest, bound to a specific deployed contract.
func NewDeployerListTest(address common.Address, backend bind.ContractBackend) (*DeployerListTest, error) {
	contract, err := bindDeployerListTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DeployerListTest{DeployerListTestCaller: DeployerListTestCaller{contract: contract}, DeployerListTestTransactor: DeployerListTestTransactor{contract: contract}, DeployerListTestFilterer: DeployerListTestFilterer{contract: contract}}, nil
}

// NewDeployerListTestCaller creates a new read-only instance of DeployerListTest, bound to a specific deployed contract.
func NewDeployerListTestCaller(address common.Address, caller bind.ContractCaller) (*DeployerListTestCaller, error) {
	contract, err := bindDeployerListTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerListTestCaller{contract: contract}, nil
}

// NewDeployerListTestTransactor creates a new write-only instance of DeployerListTest, bound to a specific deployed contract.
func NewDeployerListTestTransactor(address common.Address, transactor bind.ContractTransactor) (*DeployerListTestTransactor, error) {
	contract, err := bindDeployerListTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerListTestTransactor{contract: contract}, nil
}

// NewDeployerListTestFilterer creates a new log filterer instance of DeployerListTest, bound to a specific deployed contract.
func NewDeployerListTestFilterer(address common.Address, filterer bind.ContractFilterer) (*DeployerListTestFilterer, error) {
	contract, err := bindDeployerListTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeployerListTestFilterer{contract: contract}, nil
}

// bindDeployerListTest binds a generic wrapper to an already deployed contract.
func bindDeployerListTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DeployerListTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployerListTest *DeployerListTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployerListTest.Contract.DeployerListTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployerListTest *DeployerListTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.Contract.DeployerListTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployerListTest *DeployerListTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployerListTest.Contract.DeployerListTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployerListTest *DeployerListTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployerListTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployerListTest *DeployerListTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployerListTest *DeployerListTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployerListTest.Contract.contract.Transact(opts, method, params...)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerListTest *DeployerListTestCaller) ISTEST(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DeployerListTest.contract.Call(opts, &out, "IS_TEST")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerListTest *DeployerListTestSession) ISTEST() (bool, error) {
	return _DeployerListTest.Contract.ISTEST(&_DeployerListTest.CallOpts)
}

// ISTEST is a free data retrieval call binding the contract method 0xfa7626d4.
//
// Solidity: function IS_TEST() view returns(bool)
func (_DeployerListTest *DeployerListTestCallerSession) ISTEST() (bool, error) {
	return _DeployerListTest.Contract.ISTEST(&_DeployerListTest.CallOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerListTest *DeployerListTestTransactor) Failed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "failed")
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerListTest *DeployerListTestSession) Failed() (*types.Transaction, error) {
	return _DeployerListTest.Contract.Failed(&_DeployerListTest.TransactOpts)
}

// Failed is a paid mutator transaction binding the contract method 0xba414fa6.
//
// Solidity: function failed() returns(bool)
func (_DeployerListTest *DeployerListTestTransactorSession) Failed() (*types.Transaction, error) {
	return _DeployerListTest.Contract.Failed(&_DeployerListTest.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerListTest *DeployerListTestTransactor) SetUp(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "setUp")
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerListTest *DeployerListTestSession) SetUp() (*types.Transaction, error) {
	return _DeployerListTest.Contract.SetUp(&_DeployerListTest.TransactOpts)
}

// SetUp is a paid mutator transaction binding the contract method 0x0a9254e4.
//
// Solidity: function setUp() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) SetUp() (*types.Transaction, error) {
	return _DeployerListTest.Contract.SetUp(&_DeployerListTest.TransactOpts)
}

// StepAddDeployerThroughContract is a paid mutator transaction binding the contract method 0x2999b249.
//
// Solidity: function step_addDeployerThroughContract() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepAddDeployerThroughContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_addDeployerThroughContract")
}

// StepAddDeployerThroughContract is a paid mutator transaction binding the contract method 0x2999b249.
//
// Solidity: function step_addDeployerThroughContract() returns()
func (_DeployerListTest *DeployerListTestSession) StepAddDeployerThroughContract() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAddDeployerThroughContract(&_DeployerListTest.TransactOpts)
}

// StepAddDeployerThroughContract is a paid mutator transaction binding the contract method 0x2999b249.
//
// Solidity: function step_addDeployerThroughContract() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepAddDeployerThroughContract() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAddDeployerThroughContract(&_DeployerListTest.TransactOpts)
}

// StepAdminAddContractAsAdmin is a paid mutator transaction binding the contract method 0xffc46bc2.
//
// Solidity: function step_adminAddContractAsAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepAdminAddContractAsAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_adminAddContractAsAdmin")
}

// StepAdminAddContractAsAdmin is a paid mutator transaction binding the contract method 0xffc46bc2.
//
// Solidity: function step_adminAddContractAsAdmin() returns()
func (_DeployerListTest *DeployerListTestSession) StepAdminAddContractAsAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAdminAddContractAsAdmin(&_DeployerListTest.TransactOpts)
}

// StepAdminAddContractAsAdmin is a paid mutator transaction binding the contract method 0xffc46bc2.
//
// Solidity: function step_adminAddContractAsAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepAdminAddContractAsAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAdminAddContractAsAdmin(&_DeployerListTest.TransactOpts)
}

// StepAdminCanRevokeDeployer is a paid mutator transaction binding the contract method 0x712268f1.
//
// Solidity: function step_adminCanRevokeDeployer() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepAdminCanRevokeDeployer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_adminCanRevokeDeployer")
}

// StepAdminCanRevokeDeployer is a paid mutator transaction binding the contract method 0x712268f1.
//
// Solidity: function step_adminCanRevokeDeployer() returns()
func (_DeployerListTest *DeployerListTestSession) StepAdminCanRevokeDeployer() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAdminCanRevokeDeployer(&_DeployerListTest.TransactOpts)
}

// StepAdminCanRevokeDeployer is a paid mutator transaction binding the contract method 0x712268f1.
//
// Solidity: function step_adminCanRevokeDeployer() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepAdminCanRevokeDeployer() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepAdminCanRevokeDeployer(&_DeployerListTest.TransactOpts)
}

// StepDeployerCanDeploy is a paid mutator transaction binding the contract method 0x4c37cc2e.
//
// Solidity: function step_deployerCanDeploy() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepDeployerCanDeploy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_deployerCanDeploy")
}

// StepDeployerCanDeploy is a paid mutator transaction binding the contract method 0x4c37cc2e.
//
// Solidity: function step_deployerCanDeploy() returns()
func (_DeployerListTest *DeployerListTestSession) StepDeployerCanDeploy() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepDeployerCanDeploy(&_DeployerListTest.TransactOpts)
}

// StepDeployerCanDeploy is a paid mutator transaction binding the contract method 0x4c37cc2e.
//
// Solidity: function step_deployerCanDeploy() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepDeployerCanDeploy() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepDeployerCanDeploy(&_DeployerListTest.TransactOpts)
}

// StepNewAddressHasNoRole is a paid mutator transaction binding the contract method 0x33cb47db.
//
// Solidity: function step_newAddressHasNoRole() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepNewAddressHasNoRole(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_newAddressHasNoRole")
}

// StepNewAddressHasNoRole is a paid mutator transaction binding the contract method 0x33cb47db.
//
// Solidity: function step_newAddressHasNoRole() returns()
func (_DeployerListTest *DeployerListTestSession) StepNewAddressHasNoRole() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNewAddressHasNoRole(&_DeployerListTest.TransactOpts)
}

// StepNewAddressHasNoRole is a paid mutator transaction binding the contract method 0x33cb47db.
//
// Solidity: function step_newAddressHasNoRole() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepNewAddressHasNoRole() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNewAddressHasNoRole(&_DeployerListTest.TransactOpts)
}

// StepNoRoleCannotDeploy is a paid mutator transaction binding the contract method 0x1d44d2da.
//
// Solidity: function step_noRoleCannotDeploy() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepNoRoleCannotDeploy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_noRoleCannotDeploy")
}

// StepNoRoleCannotDeploy is a paid mutator transaction binding the contract method 0x1d44d2da.
//
// Solidity: function step_noRoleCannotDeploy() returns()
func (_DeployerListTest *DeployerListTestSession) StepNoRoleCannotDeploy() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNoRoleCannotDeploy(&_DeployerListTest.TransactOpts)
}

// StepNoRoleCannotDeploy is a paid mutator transaction binding the contract method 0x1d44d2da.
//
// Solidity: function step_noRoleCannotDeploy() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepNoRoleCannotDeploy() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNoRoleCannotDeploy(&_DeployerListTest.TransactOpts)
}

// StepNoRoleIsNotAdmin is a paid mutator transaction binding the contract method 0xf26c562c.
//
// Solidity: function step_noRoleIsNotAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepNoRoleIsNotAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_noRoleIsNotAdmin")
}

// StepNoRoleIsNotAdmin is a paid mutator transaction binding the contract method 0xf26c562c.
//
// Solidity: function step_noRoleIsNotAdmin() returns()
func (_DeployerListTest *DeployerListTestSession) StepNoRoleIsNotAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNoRoleIsNotAdmin(&_DeployerListTest.TransactOpts)
}

// StepNoRoleIsNotAdmin is a paid mutator transaction binding the contract method 0xf26c562c.
//
// Solidity: function step_noRoleIsNotAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepNoRoleIsNotAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepNoRoleIsNotAdmin(&_DeployerListTest.TransactOpts)
}

// StepVerifySenderIsAdmin is a paid mutator transaction binding the contract method 0x28002804.
//
// Solidity: function step_verifySenderIsAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactor) StepVerifySenderIsAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployerListTest.contract.Transact(opts, "step_verifySenderIsAdmin")
}

// StepVerifySenderIsAdmin is a paid mutator transaction binding the contract method 0x28002804.
//
// Solidity: function step_verifySenderIsAdmin() returns()
func (_DeployerListTest *DeployerListTestSession) StepVerifySenderIsAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepVerifySenderIsAdmin(&_DeployerListTest.TransactOpts)
}

// StepVerifySenderIsAdmin is a paid mutator transaction binding the contract method 0x28002804.
//
// Solidity: function step_verifySenderIsAdmin() returns()
func (_DeployerListTest *DeployerListTestTransactorSession) StepVerifySenderIsAdmin() (*types.Transaction, error) {
	return _DeployerListTest.Contract.StepVerifySenderIsAdmin(&_DeployerListTest.TransactOpts)
}

// DeployerListTestLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the DeployerListTest contract.
type DeployerListTestLogIterator struct {
	Event *DeployerListTestLog // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLog)
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
		it.Event = new(DeployerListTestLog)
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
func (it *DeployerListTestLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLog represents a Log event raised by the DeployerListTest contract.
type DeployerListTestLog struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLog(opts *bind.FilterOpts) (*DeployerListTestLogIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogIterator{contract: _DeployerListTest.contract, event: "log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x41304facd9323d75b11bcdd609cb38effffdb05710f7caf0e9b16c6d9d709f50.
//
// Solidity: event log(string arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *DeployerListTestLog) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLog)
				if err := _DeployerListTest.contract.UnpackLog(event, "log", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLog(log types.Log) (*DeployerListTestLog, error) {
	event := new(DeployerListTestLog)
	if err := _DeployerListTest.contract.UnpackLog(event, "log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogAddressIterator is returned from FilterLogAddress and is used to iterate over the raw logs and unpacked data for LogAddress events raised by the DeployerListTest contract.
type DeployerListTestLogAddressIterator struct {
	Event *DeployerListTestLogAddress // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogAddress)
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
		it.Event = new(DeployerListTestLogAddress)
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
func (it *DeployerListTestLogAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogAddress represents a LogAddress event raised by the DeployerListTest contract.
type DeployerListTestLogAddress struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogAddress is a free log retrieval operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogAddress(opts *bind.FilterOpts) (*DeployerListTestLogAddressIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogAddressIterator{contract: _DeployerListTest.contract, event: "log_address", logs: logs, sub: sub}, nil
}

// WatchLogAddress is a free log subscription operation binding the contract event 0x7ae74c527414ae135fd97047b12921a5ec3911b804197855d67e25c7b75ee6f3.
//
// Solidity: event log_address(address arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogAddress(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogAddress) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogAddress)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_address", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogAddress(log types.Log) (*DeployerListTestLogAddress, error) {
	event := new(DeployerListTestLogAddress)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogBytesIterator is returned from FilterLogBytes and is used to iterate over the raw logs and unpacked data for LogBytes events raised by the DeployerListTest contract.
type DeployerListTestLogBytesIterator struct {
	Event *DeployerListTestLogBytes // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogBytes)
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
		it.Event = new(DeployerListTestLogBytes)
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
func (it *DeployerListTestLogBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogBytes represents a LogBytes event raised by the DeployerListTest contract.
type DeployerListTestLogBytes struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes is a free log retrieval operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogBytes(opts *bind.FilterOpts) (*DeployerListTestLogBytesIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogBytesIterator{contract: _DeployerListTest.contract, event: "log_bytes", logs: logs, sub: sub}, nil
}

// WatchLogBytes is a free log subscription operation binding the contract event 0x23b62ad0584d24a75f0bf3560391ef5659ec6db1269c56e11aa241d637f19b20.
//
// Solidity: event log_bytes(bytes arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogBytes(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogBytes) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogBytes)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_bytes", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogBytes(log types.Log) (*DeployerListTestLogBytes, error) {
	event := new(DeployerListTestLogBytes)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogBytes32Iterator is returned from FilterLogBytes32 and is used to iterate over the raw logs and unpacked data for LogBytes32 events raised by the DeployerListTest contract.
type DeployerListTestLogBytes32Iterator struct {
	Event *DeployerListTestLogBytes32 // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogBytes32)
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
		it.Event = new(DeployerListTestLogBytes32)
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
func (it *DeployerListTestLogBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogBytes32 represents a LogBytes32 event raised by the DeployerListTest contract.
type DeployerListTestLogBytes32 struct {
	Arg0 [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogBytes32 is a free log retrieval operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogBytes32(opts *bind.FilterOpts) (*DeployerListTestLogBytes32Iterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogBytes32Iterator{contract: _DeployerListTest.contract, event: "log_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogBytes32 is a free log subscription operation binding the contract event 0xe81699b85113eea1c73e10588b2b035e55893369632173afd43feb192fac64e3.
//
// Solidity: event log_bytes32(bytes32 arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogBytes32(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogBytes32) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogBytes32)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_bytes32", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogBytes32(log types.Log) (*DeployerListTestLogBytes32, error) {
	event := new(DeployerListTestLogBytes32)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogIntIterator is returned from FilterLogInt and is used to iterate over the raw logs and unpacked data for LogInt events raised by the DeployerListTest contract.
type DeployerListTestLogIntIterator struct {
	Event *DeployerListTestLogInt // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogInt)
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
		it.Event = new(DeployerListTestLogInt)
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
func (it *DeployerListTestLogIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogInt represents a LogInt event raised by the DeployerListTest contract.
type DeployerListTestLogInt struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogInt is a free log retrieval operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogInt(opts *bind.FilterOpts) (*DeployerListTestLogIntIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogIntIterator{contract: _DeployerListTest.contract, event: "log_int", logs: logs, sub: sub}, nil
}

// WatchLogInt is a free log subscription operation binding the contract event 0x0eb5d52624c8d28ada9fc55a8c502ed5aa3fbe2fb6e91b71b5f376882b1d2fb8.
//
// Solidity: event log_int(int256 arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogInt(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogInt) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogInt)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_int", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogInt(log types.Log) (*DeployerListTestLogInt, error) {
	event := new(DeployerListTestLogInt)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedAddressIterator is returned from FilterLogNamedAddress and is used to iterate over the raw logs and unpacked data for LogNamedAddress events raised by the DeployerListTest contract.
type DeployerListTestLogNamedAddressIterator struct {
	Event *DeployerListTestLogNamedAddress // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedAddress)
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
		it.Event = new(DeployerListTestLogNamedAddress)
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
func (it *DeployerListTestLogNamedAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedAddress represents a LogNamedAddress event raised by the DeployerListTest contract.
type DeployerListTestLogNamedAddress struct {
	Key string
	Val common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedAddress is a free log retrieval operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedAddress(opts *bind.FilterOpts) (*DeployerListTestLogNamedAddressIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedAddressIterator{contract: _DeployerListTest.contract, event: "log_named_address", logs: logs, sub: sub}, nil
}

// WatchLogNamedAddress is a free log subscription operation binding the contract event 0x9c4e8541ca8f0dc1c413f9108f66d82d3cecb1bddbce437a61caa3175c4cc96f.
//
// Solidity: event log_named_address(string key, address val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedAddress(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedAddress) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_address")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedAddress)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_address", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedAddress(log types.Log) (*DeployerListTestLogNamedAddress, error) {
	event := new(DeployerListTestLogNamedAddress)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_address", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedBytesIterator is returned from FilterLogNamedBytes and is used to iterate over the raw logs and unpacked data for LogNamedBytes events raised by the DeployerListTest contract.
type DeployerListTestLogNamedBytesIterator struct {
	Event *DeployerListTestLogNamedBytes // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedBytesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedBytes)
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
		it.Event = new(DeployerListTestLogNamedBytes)
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
func (it *DeployerListTestLogNamedBytesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedBytesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedBytes represents a LogNamedBytes event raised by the DeployerListTest contract.
type DeployerListTestLogNamedBytes struct {
	Key string
	Val []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes is a free log retrieval operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedBytes(opts *bind.FilterOpts) (*DeployerListTestLogNamedBytesIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedBytesIterator{contract: _DeployerListTest.contract, event: "log_named_bytes", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes is a free log subscription operation binding the contract event 0xd26e16cad4548705e4c9e2d94f98ee91c289085ee425594fd5635fa2964ccf18.
//
// Solidity: event log_named_bytes(string key, bytes val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedBytes(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedBytes) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_bytes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedBytes)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedBytes(log types.Log) (*DeployerListTestLogNamedBytes, error) {
	event := new(DeployerListTestLogNamedBytes)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_bytes", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedBytes32Iterator is returned from FilterLogNamedBytes32 and is used to iterate over the raw logs and unpacked data for LogNamedBytes32 events raised by the DeployerListTest contract.
type DeployerListTestLogNamedBytes32Iterator struct {
	Event *DeployerListTestLogNamedBytes32 // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedBytes32Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedBytes32)
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
		it.Event = new(DeployerListTestLogNamedBytes32)
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
func (it *DeployerListTestLogNamedBytes32Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedBytes32Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedBytes32 represents a LogNamedBytes32 event raised by the DeployerListTest contract.
type DeployerListTestLogNamedBytes32 struct {
	Key string
	Val [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedBytes32 is a free log retrieval operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedBytes32(opts *bind.FilterOpts) (*DeployerListTestLogNamedBytes32Iterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedBytes32Iterator{contract: _DeployerListTest.contract, event: "log_named_bytes32", logs: logs, sub: sub}, nil
}

// WatchLogNamedBytes32 is a free log subscription operation binding the contract event 0xafb795c9c61e4fe7468c386f925d7a5429ecad9c0495ddb8d38d690614d32f99.
//
// Solidity: event log_named_bytes32(string key, bytes32 val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedBytes32(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedBytes32) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_bytes32")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedBytes32)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedBytes32(log types.Log) (*DeployerListTestLogNamedBytes32, error) {
	event := new(DeployerListTestLogNamedBytes32)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_bytes32", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedDecimalIntIterator is returned from FilterLogNamedDecimalInt and is used to iterate over the raw logs and unpacked data for LogNamedDecimalInt events raised by the DeployerListTest contract.
type DeployerListTestLogNamedDecimalIntIterator struct {
	Event *DeployerListTestLogNamedDecimalInt // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedDecimalIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedDecimalInt)
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
		it.Event = new(DeployerListTestLogNamedDecimalInt)
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
func (it *DeployerListTestLogNamedDecimalIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedDecimalIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedDecimalInt represents a LogNamedDecimalInt event raised by the DeployerListTest contract.
type DeployerListTestLogNamedDecimalInt struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalInt is a free log retrieval operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedDecimalInt(opts *bind.FilterOpts) (*DeployerListTestLogNamedDecimalIntIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedDecimalIntIterator{contract: _DeployerListTest.contract, event: "log_named_decimal_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalInt is a free log subscription operation binding the contract event 0x5da6ce9d51151ba10c09a559ef24d520b9dac5c5b8810ae8434e4d0d86411a95.
//
// Solidity: event log_named_decimal_int(string key, int256 val, uint256 decimals)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedDecimalInt(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedDecimalInt) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_decimal_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedDecimalInt)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedDecimalInt(log types.Log) (*DeployerListTestLogNamedDecimalInt, error) {
	event := new(DeployerListTestLogNamedDecimalInt)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_decimal_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedDecimalUintIterator is returned from FilterLogNamedDecimalUint and is used to iterate over the raw logs and unpacked data for LogNamedDecimalUint events raised by the DeployerListTest contract.
type DeployerListTestLogNamedDecimalUintIterator struct {
	Event *DeployerListTestLogNamedDecimalUint // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedDecimalUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedDecimalUint)
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
		it.Event = new(DeployerListTestLogNamedDecimalUint)
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
func (it *DeployerListTestLogNamedDecimalUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedDecimalUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedDecimalUint represents a LogNamedDecimalUint event raised by the DeployerListTest contract.
type DeployerListTestLogNamedDecimalUint struct {
	Key      string
	Val      *big.Int
	Decimals *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLogNamedDecimalUint is a free log retrieval operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedDecimalUint(opts *bind.FilterOpts) (*DeployerListTestLogNamedDecimalUintIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedDecimalUintIterator{contract: _DeployerListTest.contract, event: "log_named_decimal_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedDecimalUint is a free log subscription operation binding the contract event 0xeb8ba43ced7537421946bd43e828b8b2b8428927aa8f801c13d934bf11aca57b.
//
// Solidity: event log_named_decimal_uint(string key, uint256 val, uint256 decimals)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedDecimalUint(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedDecimalUint) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_decimal_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedDecimalUint)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedDecimalUint(log types.Log) (*DeployerListTestLogNamedDecimalUint, error) {
	event := new(DeployerListTestLogNamedDecimalUint)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_decimal_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedIntIterator is returned from FilterLogNamedInt and is used to iterate over the raw logs and unpacked data for LogNamedInt events raised by the DeployerListTest contract.
type DeployerListTestLogNamedIntIterator struct {
	Event *DeployerListTestLogNamedInt // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedInt)
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
		it.Event = new(DeployerListTestLogNamedInt)
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
func (it *DeployerListTestLogNamedIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedInt represents a LogNamedInt event raised by the DeployerListTest contract.
type DeployerListTestLogNamedInt struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedInt is a free log retrieval operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedInt(opts *bind.FilterOpts) (*DeployerListTestLogNamedIntIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedIntIterator{contract: _DeployerListTest.contract, event: "log_named_int", logs: logs, sub: sub}, nil
}

// WatchLogNamedInt is a free log subscription operation binding the contract event 0x2fe632779174374378442a8e978bccfbdcc1d6b2b0d81f7e8eb776ab2286f168.
//
// Solidity: event log_named_int(string key, int256 val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedInt(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedInt) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_int")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedInt)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_int", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedInt(log types.Log) (*DeployerListTestLogNamedInt, error) {
	event := new(DeployerListTestLogNamedInt)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_int", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedStringIterator is returned from FilterLogNamedString and is used to iterate over the raw logs and unpacked data for LogNamedString events raised by the DeployerListTest contract.
type DeployerListTestLogNamedStringIterator struct {
	Event *DeployerListTestLogNamedString // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedString)
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
		it.Event = new(DeployerListTestLogNamedString)
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
func (it *DeployerListTestLogNamedStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedString represents a LogNamedString event raised by the DeployerListTest contract.
type DeployerListTestLogNamedString struct {
	Key string
	Val string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedString is a free log retrieval operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedString(opts *bind.FilterOpts) (*DeployerListTestLogNamedStringIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedStringIterator{contract: _DeployerListTest.contract, event: "log_named_string", logs: logs, sub: sub}, nil
}

// WatchLogNamedString is a free log subscription operation binding the contract event 0x280f4446b28a1372417dda658d30b95b2992b12ac9c7f378535f29a97acf3583.
//
// Solidity: event log_named_string(string key, string val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedString(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedString) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedString)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_string", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedString(log types.Log) (*DeployerListTestLogNamedString, error) {
	event := new(DeployerListTestLogNamedString)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogNamedUintIterator is returned from FilterLogNamedUint and is used to iterate over the raw logs and unpacked data for LogNamedUint events raised by the DeployerListTest contract.
type DeployerListTestLogNamedUintIterator struct {
	Event *DeployerListTestLogNamedUint // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogNamedUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogNamedUint)
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
		it.Event = new(DeployerListTestLogNamedUint)
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
func (it *DeployerListTestLogNamedUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogNamedUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogNamedUint represents a LogNamedUint event raised by the DeployerListTest contract.
type DeployerListTestLogNamedUint struct {
	Key string
	Val *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNamedUint is a free log retrieval operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogNamedUint(opts *bind.FilterOpts) (*DeployerListTestLogNamedUintIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogNamedUintIterator{contract: _DeployerListTest.contract, event: "log_named_uint", logs: logs, sub: sub}, nil
}

// WatchLogNamedUint is a free log subscription operation binding the contract event 0xb2de2fbe801a0df6c0cbddfd448ba3c41d48a040ca35c56c8196ef0fcae721a8.
//
// Solidity: event log_named_uint(string key, uint256 val)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogNamedUint(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogNamedUint) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_named_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogNamedUint)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_named_uint", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogNamedUint(log types.Log) (*DeployerListTestLogNamedUint, error) {
	event := new(DeployerListTestLogNamedUint)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_named_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogStringIterator is returned from FilterLogString and is used to iterate over the raw logs and unpacked data for LogString events raised by the DeployerListTest contract.
type DeployerListTestLogStringIterator struct {
	Event *DeployerListTestLogString // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogStringIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogString)
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
		it.Event = new(DeployerListTestLogString)
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
func (it *DeployerListTestLogStringIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogStringIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogString represents a LogString event raised by the DeployerListTest contract.
type DeployerListTestLogString struct {
	Arg0 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogString is a free log retrieval operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogString(opts *bind.FilterOpts) (*DeployerListTestLogStringIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogStringIterator{contract: _DeployerListTest.contract, event: "log_string", logs: logs, sub: sub}, nil
}

// WatchLogString is a free log subscription operation binding the contract event 0x0b2e13ff20ac7b474198655583edf70dedd2c1dc980e329c4fbb2fc0748b796b.
//
// Solidity: event log_string(string arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogString(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogString) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_string")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogString)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_string", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogString(log types.Log) (*DeployerListTestLogString, error) {
	event := new(DeployerListTestLogString)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_string", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogUintIterator is returned from FilterLogUint and is used to iterate over the raw logs and unpacked data for LogUint events raised by the DeployerListTest contract.
type DeployerListTestLogUintIterator struct {
	Event *DeployerListTestLogUint // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogUint)
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
		it.Event = new(DeployerListTestLogUint)
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
func (it *DeployerListTestLogUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogUint represents a LogUint event raised by the DeployerListTest contract.
type DeployerListTestLogUint struct {
	Arg0 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogUint is a free log retrieval operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogUint(opts *bind.FilterOpts) (*DeployerListTestLogUintIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogUintIterator{contract: _DeployerListTest.contract, event: "log_uint", logs: logs, sub: sub}, nil
}

// WatchLogUint is a free log subscription operation binding the contract event 0x2cab9790510fd8bdfbd2115288db33fec66691d476efc5427cfd4c0969301755.
//
// Solidity: event log_uint(uint256 arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogUint(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogUint) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "log_uint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogUint)
				if err := _DeployerListTest.contract.UnpackLog(event, "log_uint", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogUint(log types.Log) (*DeployerListTestLogUint, error) {
	event := new(DeployerListTestLogUint)
	if err := _DeployerListTest.contract.UnpackLog(event, "log_uint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerListTestLogsIterator is returned from FilterLogs and is used to iterate over the raw logs and unpacked data for Logs events raised by the DeployerListTest contract.
type DeployerListTestLogsIterator struct {
	Event *DeployerListTestLogs // Event containing the contract specifics and raw log

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
func (it *DeployerListTestLogsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerListTestLogs)
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
		it.Event = new(DeployerListTestLogs)
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
func (it *DeployerListTestLogsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerListTestLogsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerListTestLogs represents a Logs event raised by the DeployerListTest contract.
type DeployerListTestLogs struct {
	Arg0 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogs is a free log retrieval operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DeployerListTest *DeployerListTestFilterer) FilterLogs(opts *bind.FilterOpts) (*DeployerListTestLogsIterator, error) {

	logs, sub, err := _DeployerListTest.contract.FilterLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return &DeployerListTestLogsIterator{contract: _DeployerListTest.contract, event: "logs", logs: logs, sub: sub}, nil
}

// WatchLogs is a free log subscription operation binding the contract event 0xe7950ede0394b9f2ce4a5a1bf5a7e1852411f7e6661b4308c913c4bfd11027e4.
//
// Solidity: event logs(bytes arg0)
func (_DeployerListTest *DeployerListTestFilterer) WatchLogs(opts *bind.WatchOpts, sink chan<- *DeployerListTestLogs) (event.Subscription, error) {

	logs, sub, err := _DeployerListTest.contract.WatchLogs(opts, "logs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerListTestLogs)
				if err := _DeployerListTest.contract.UnpackLog(event, "logs", log); err != nil {
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
func (_DeployerListTest *DeployerListTestFilterer) ParseLogs(log types.Log) (*DeployerListTestLogs, error) {
	event := new(DeployerListTestLogs)
	if err := _DeployerListTest.contract.UnpackLog(event, "logs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
