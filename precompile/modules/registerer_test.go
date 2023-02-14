// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestRegisterModule(t *testing.T) {
	data := make([]Module, 0)
	// test that the module is registered in sorted order
	module1 := Module{
		Address: common.BigToAddress(big.NewInt(1)),
	}
	data = insertSortedByAddress(data, module1)

	require.Equal(t, []Module{module1}, data)

	module0 := Module{
		Address: common.BigToAddress(big.NewInt(0)),
	}

	data = insertSortedByAddress(data, module0)
	require.Equal(t, []Module{module0, module1}, data)

	module3 := Module{
		Address: common.BigToAddress(big.NewInt(3)),
	}

	data = insertSortedByAddress(data, module3)
	require.Equal(t, []Module{module0, module1, module3}, data)

	module2 := Module{
		Address: common.BigToAddress(big.NewInt(2)),
	}

	data = insertSortedByAddress(data, module2)
	require.Equal(t, []Module{module0, module1, module2, module3}, data)
}
