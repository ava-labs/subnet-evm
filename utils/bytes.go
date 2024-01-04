// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"math/big"

	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ethereum/go-ethereum/common"
)

const BigIntBytesLength = 32

// IncrOne increments bytes value by one
func IncrOne(bytes []byte) {
	index := len(bytes) - 1
	for index >= 0 {
		if bytes[index] < 255 {
			bytes[index]++
			break
		} else {
			bytes[index] = 0
			index--
		}
	}
}

// HashSliceToBytes serializes a []common.Hash into a tightly packed byte array.
func HashSliceToBytes(hashes []common.Hash) []byte {
	bytes := make([]byte, common.HashLength*len(hashes))
	for i, hash := range hashes {
		copy(bytes[i*common.HashLength:], hash[:])
	}
	return bytes
}

// BytesToHashSlice packs [b] into a slice of hash values with zero padding
// to the right if the length of b is not a multiple of 32.
func BytesToHashSlice(b []byte) []common.Hash {
	var (
		numHashes = (len(b) + 31) / 32
		hashes    = make([]common.Hash, numHashes)
	)

	for i := range hashes {
		start := i * common.HashLength
		copy(hashes[i][:], b[start:])
	}
	return hashes
}

func PackBigInt(p *wrappers.Packer, number *big.Int) error {
	p.PackBool(number == nil)
	if p.Err == nil && number != nil {
		p.PackFixedBytes(number.FillBytes(make([]byte, BigIntBytesLength)))
	}

	return p.Err
}

func UnpackBigInt(p *wrappers.Packer) (*big.Int, error) {
	isNil := p.UnpackBool()
	if p.Err != nil || isNil {
		return nil, p.Err
	}

	number := big.NewInt(0).SetBytes(p.UnpackFixedBytes(BigIntBytesLength))
	return number, p.Err
}

func PackAddresses(p *wrappers.Packer, addresses []common.Address) error {
	p.PackBool(addresses == nil)
	if addresses == nil {
		return nil
	}
	p.PackInt(uint32(len(addresses)))
	if p.Err != nil {
		return p.Err
	}
	for _, address := range addresses {
		p.PackFixedBytes(address[:])
		if p.Err != nil {
			return p.Err
		}
	}
	return nil
}

func UnpackAddresses(p *wrappers.Packer) ([]common.Address, error) {
	isNil := p.UnpackBool()
	if isNil || p.Err != nil {
		return nil, p.Err
	}
	length := p.UnpackInt()
	if p.Err != nil {
		return nil, p.Err
	}

	addresses := make([]common.Address, 0, length)
	for i := uint32(0); i < length; i++ {
		bytes := p.UnpackFixedBytes(common.AddressLength)
		addresses = append(addresses, common.BytesToAddress(bytes))
		if p.Err != nil {
			return nil, p.Err
		}
	}

	return addresses, nil
}
