// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/load"
)

const Version = "v0.1.0"

func main() {
	config := config.BuildConfig(os.Args[1:])
	if err := load.ExecuteLoader(context.Background(), config); err != nil {
		fmt.Printf("load execution failed: %s\n", err)
		os.Exit(1)
	}
}
