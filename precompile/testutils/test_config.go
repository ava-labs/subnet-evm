// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package testutils

import (
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
)

// PrecompileConfigTest is a test case for precompile configs
type ConfigVerifyTest struct {
	Config        precompileconfig.Config
	ExpectedError string
}
