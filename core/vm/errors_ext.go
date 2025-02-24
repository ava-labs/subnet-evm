// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import "errors"

var (
	ErrAddrProhibited              = errors.New("prohibited address cannot be sender or created contract address")
	ErrInvalidCoinbase             = errors.New("invalid coinbase")
	ErrSenderAddressNotAllowListed = errors.New("cannot issue transaction from non-allow listed address")
)
