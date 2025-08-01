// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ava-labs/libevm/common"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalConfig(t *testing.T) {
	tests := []struct {
		name        string
		givenJSON   []byte
		expected    Config
		expectedErr bool
	}{
		{
			"string durations parsed",
			[]byte(`{"api-max-duration": "1m", "continuous-profiler-frequency": "2m"}`),
			Config{APIMaxDuration: Duration{1 * time.Minute}, ContinuousProfilerFrequency: Duration{2 * time.Minute}},
			false,
		},
		{
			"integer durations parsed",
			[]byte(fmt.Sprintf(`{"api-max-duration": "%v", "continuous-profiler-frequency": "%v"}`, 1*time.Minute, 2*time.Minute)),
			Config{APIMaxDuration: Duration{1 * time.Minute}, ContinuousProfilerFrequency: Duration{2 * time.Minute}},
			false,
		},
		{
			"nanosecond durations parsed",
			[]byte(`{"api-max-duration": 5000000000, "continuous-profiler-frequency": 5000000000}`),
			Config{APIMaxDuration: Duration{5 * time.Second}, ContinuousProfilerFrequency: Duration{5 * time.Second}},
			false,
		},
		{
			"bad durations",
			[]byte(`{"api-max-duration": "bad-duration"}`),
			Config{},
			true,
		},

		{
			"tx pool configurations",
			[]byte(`{"tx-pool-price-limit": 1, "tx-pool-price-bump": 2, "tx-pool-account-slots": 3, "tx-pool-global-slots": 4, "tx-pool-account-queue": 5, "tx-pool-global-queue": 6}`),
			Config{
				TxPoolPriceLimit:   1,
				TxPoolPriceBump:    2,
				TxPoolAccountSlots: 3,
				TxPoolGlobalSlots:  4,
				TxPoolAccountQueue: 5,
				TxPoolGlobalQueue:  6,
			},
			false,
		},

		{
			"state sync enabled",
			[]byte(`{"state-sync-enabled":true}`),
			Config{StateSyncEnabled: true},
			false,
		},
		{
			"state sync sources",
			[]byte(`{"state-sync-ids": "NodeID-CaBYJ9kzHvrQFiYWowMkJGAQKGMJqZoat"}`),
			Config{StateSyncIDs: "NodeID-CaBYJ9kzHvrQFiYWowMkJGAQKGMJqZoat"},
			false,
		},
		{
			"empty transaction history ",
			[]byte(`{}`),
			Config{TransactionHistory: 0},
			false,
		},
		{
			"zero transaction history",
			[]byte(`{"transaction-history": 0}`),
			func() Config {
				return Config{TransactionHistory: 0}
			}(),
			false,
		},
		{
			"1 transaction history",
			[]byte(`{"transaction-history": 1}`),
			func() Config {
				return Config{TransactionHistory: 1}
			}(),
			false,
		},
		{
			"-1 transaction history",
			[]byte(`{"transaction-history": -1}`),
			Config{},
			true,
		},
		{
			"deprecated tx lookup limit",
			[]byte(`{"tx-lookup-limit": 1}`),
			Config{TransactionHistory: 1, TxLookupLimit: 1},
			false,
		},
		{
			"allow unprotected tx hashes",
			[]byte(`{"allow-unprotected-tx-hashes": ["0x803351deb6d745e91545a6a3e1c0ea3e9a6a02a1a4193b70edfcd2f40f71a01c"]}`),
			Config{AllowUnprotectedTxHashes: []common.Hash{common.HexToHash("0x803351deb6d745e91545a6a3e1c0ea3e9a6a02a1a4193b70edfcd2f40f71a01c")}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tmp Config
			err := json.Unmarshal(tt.givenJSON, &tmp)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tmp.Deprecate()
				assert.Equal(t, tt.expected, tmp)
			}
		})
	}
}
