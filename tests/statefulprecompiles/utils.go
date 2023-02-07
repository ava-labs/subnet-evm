// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package statefulprecompiles

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

type precompileTest struct {
	caller      common.Address
	input       func() []byte
	suppliedGas uint64
	readOnly    bool

	config config.Config

	preCondition func(t *testing.T, state *state.StateDB)
	assertState  func(t *testing.T, state *state.StateDB)

	expectedRes []byte
	expectedErr string
}
