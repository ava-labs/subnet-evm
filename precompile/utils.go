// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var functionSignatureRegex = regexp.MustCompile(`[\w]+\(((([\w]+)?)|((([\w]+),)+([\w]+)))\)`)

// CalculateFunctionSelector returns the 4 byte function selector that results from [functionSignature]
// Ex. the function setBalance(addr address, balance uint256) should be passed in as the string:
// "setBalance(address,uint256)"
func CalculateFunctionSelector(functionSignature string) []byte {
	if !functionSignatureRegex.MatchString(functionSignature) {
		panic(fmt.Errorf("invalid function signature: %q", functionSignature))
	}
	hash := crypto.Keccak256([]byte(functionSignature))
	return hash[:4]
}

// deductGas checks if [suppliedGas] is sufficient against [requiredGas] and deducts [requiredGas] from [suppliedGas].
func deductGas(suppliedGas uint64, requiredGas uint64) (uint64, error) {
	if suppliedGas < requiredGas {
		return 0, vmerrs.ErrOutOfGas
	}
	return suppliedGas - requiredGas, nil
}

// packOrderedHashesWithSelector checks fullLength of given [input]
// it excludes first member since it should be the function selector
// then checks if the given [fullLength] is a multiple of member count * common.HashLength
func packOrderedHashesWithSelector(input [][]byte, fullLength int) ([]byte, error) {
	hashLen := fullLength - selectorLen
	realLen := (len(input) - 1) * common.HashLength
	if hashLen != realLen {
		return nil, fmt.Errorf("expected %d, got %d length", hashLen, realLen)
	}

	// check function selector
	if selectorLen != len(input[0]) {
		return nil, fmt.Errorf("first element of the input must be a function selector with length %d", selectorLen)
	}
	// first handle selector
	buf := make([]byte, fullLength)
	copy(buf[:selectorLen], input[0])

	// handle bytesHashes
	bytesHashes, err := packOrderedHashes(input[1:], hashLen)
	if err != nil {
		return nil, err
	}
	copy(buf[selectorLen:], bytesHashes)

	return buf, nil
}

// packOrderedHashes packs 2-d [input] array of hashes into a 1-d array.
func packOrderedHashes(input [][]byte, fullLength int) ([]byte, error) {
	realLen := len(input) * common.HashLength
	if fullLength != realLen {
		return nil, fmt.Errorf("expected %d, got %d length", fullLength, realLen)
	}

	buf := make([]byte, fullLength)
	for index, inputByte := range input {
		start := (common.HashLength * index)
		end := start + common.HashLength
		copy(buf[start:end], inputByte)
	}
	return buf, nil
}

// returnPackedElement returns packed element with common.HashLength from the [packed] at [index]
func returnPackedElement(packed []byte, index int) []byte {
	start := common.HashLength * index
	end := start + common.HashLength
	return packed[start:end]
}

func bigToHashSafe(b *big.Int) (common.Hash, error) {
	bytes := b.Bytes()
	len := len(bytes)
	if len > common.HashLength {
		return common.Hash{}, fmt.Errorf("expected %d, got %d length", common.HashLength, len)
	}
	return common.BigToHash(b), nil
}
