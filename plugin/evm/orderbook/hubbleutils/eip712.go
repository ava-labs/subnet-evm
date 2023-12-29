package hubbleutils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// EncodeForSigning - Encoding the typed data
func EncodeForSigning(typedData apitypes.TypedData) (hash common.Hash, err error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash = common.BytesToHash(crypto.Keccak256(rawData))
	return
}

var Eip712OrderTypes = apitypes.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"Order": { // has to be same as the struct name or whatever was passed when building the typed hash
		{
			Name: "orderType",
			Type: "uint8",
		},
		{
			Name: "expireAt",
			Type: "uint256",
		},
		{
			Name: "ammIndex",
			Type: "uint256",
		},
		{
			Name: "trader",
			Type: "address",
		},
		{
			Name: "baseAssetQuantity",
			Type: "int256",
		},
		{
			Name: "price",
			Type: "uint256",
		},
		{
			Name: "salt",
			Type: "uint256",
		},
		{
			Name: "reduceOnly",
			Type: "bool",
		},
		{
			Name: "postOnly",
			Type: "bool",
		},
	},
}
