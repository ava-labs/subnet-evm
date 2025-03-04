// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// ACP176 implements the fee logic specified here:
// https://github.com/avalanche-foundation/ACPs/blob/main/ACPs/176-dynamic-evm-gas-limit-and-price-discovery-updates/README.md
package acp176

import (
	"fmt"

	"github.com/ava-labs/avalanchego/vms/components/gas"
)

var DefaultAcp176Config = &Acp176Config{
	MinTargetPerSecond:            1_000_000,
	MinGasPrice:                   1,
	TargetToMax:                   2,
	TimeToFillCapacity:            5,
	TargetToPriceUpdateConversion: 87,
	MaxTargetChangeRate:           1024,
}

type Acp176Config struct {
	MinTargetPerSecond            uint64    // P
	MinGasPrice                   gas.Price // M
	TargetToMax                   uint64    // multiplier to convert from target per second to max per second
	TimeToFillCapacity            uint64    // in seconds
	TargetToPriceUpdateConversion gas.Gas   // 60s ~= 87 * ln(2), so choose 87 for 1 min to double fee
	MaxTargetChangeRate           uint64    // Controls the rate that the target can change per block.
}

func (c *Acp176Config) Verify() error {
	switch {
	case c.MinTargetPerSecond <= 0:
		return fmt.Errorf("MinTargetPerSecond %d cannot generate blocks", c.MinTargetPerSecond)
	case c.MinGasPrice <= 0:
		return fmt.Errorf("MinGasPrice must be nonzero")
	case c.TargetToMax <= 1:
		return fmt.Errorf("Max block size must be greater than target, current ratio %d", c.TargetToMax)
	case c.TimeToFillCapacity <= 1:
		return fmt.Errorf("Capacity must fill slower than blocks are being issued, current time %d", c.TimeToFillCapacity)
	case c.TargetToPriceUpdateConversion <= 0:
		return fmt.Errorf("TargetToPriceUpdateConversion must be positive, current value %d", c.TargetToPriceUpdateConversion)
	case c.MaxTargetChangeRate <= 0:
		return fmt.Errorf("MaxTargetChangeRate must be positive, current value %d", c.MaxTargetChangeRate)
	}

	return nil
}
