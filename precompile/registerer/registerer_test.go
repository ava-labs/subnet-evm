// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package registerer

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestRegisterModule(t *testing.T) {
	data := make([]contract.Module, 0)
	// test that the module is registered in sorted order
	module1 := contract.Module{
		Address: common.BigToAddress(big.NewInt(1)),
	}
	data = insertSortedByAddress(data, module1)

	require.Equal(t, []contract.Module{module1}, data)

	module0 := contract.Module{
		Address: common.BigToAddress(big.NewInt(0)),
	}

	data = insertSortedByAddress(data, module0)
	require.Equal(t, []contract.Module{module0, module1}, data)

	module3 := contract.Module{
		Address: common.BigToAddress(big.NewInt(3)),
	}

	data = insertSortedByAddress(data, module3)
	require.Equal(t, []contract.Module{module0, module1, module3}, data)

	module2 := contract.Module{
		Address: common.BigToAddress(big.NewInt(2)),
	}

	data = insertSortedByAddress(data, module2)
	require.Equal(t, []contract.Module{module0, module1, module2, module3}, data)
}
