// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customethclient

import (
	"context"
	"fmt"

	"github.com/ava-labs/subnet-evm/plugin/evm/customtypes"
)

const FujiAPIURI = "https://api.avax-test.network"

func ExampleClient() {
	ethC, err := Dial(FujiAPIURI + "/ext/bc/C/rpc")
	if err != nil {
		panic(err)
	}
	bc, err := ethC.client.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	extData := customtypes.GetHeaderExtra(bc.Header())
	// Header extra data
	fmt.Printf("Block Gas Cost: %d\n", extData.BlockGasCost)
}
