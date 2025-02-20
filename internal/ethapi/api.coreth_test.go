// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ethapi

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/coreth/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

type testSuggestPriceOptionsBackend struct {
	Backend // embed the interface to avoid implementing unused methods

	estimateBaseFee  *big.Int
	suggestGasTipCap *big.Int

	cfg PriceOptionConfig
}

func (b *testSuggestPriceOptionsBackend) EstimateBaseFee(context.Context) (*big.Int, error) {
	return b.estimateBaseFee, nil
}

func (b *testSuggestPriceOptionsBackend) SuggestGasTipCap(context.Context) (*big.Int, error) {
	return b.suggestGasTipCap, nil
}

func (b *testSuggestPriceOptionsBackend) PriceOptionsConfig() PriceOptionConfig {
	return b.cfg
}

func TestSuggestPriceOptions(t *testing.T) {
	testCfg := PriceOptionConfig{
		SlowFeePercentage: 95,
		FastFeePercentage: 105,
		MaxBaseFee:        100 * params.GWei,
		MaxTip:            20 * params.GWei,
	}
	slowFeeNumerator := testCfg.SlowFeePercentage
	fastFeeNumerator := testCfg.FastFeePercentage
	maxNormalGasTip := testCfg.MaxTip
	maxNormalBaseFee := testCfg.MaxBaseFee

	tests := []struct {
		name             string
		estimateBaseFee  *big.Int
		suggestGasTipCap *big.Int
		cfg              PriceOptionConfig
		want             *PriceOptions
	}{
		{
			name:             "nil_base_fee",
			estimateBaseFee:  nil,
			suggestGasTipCap: common.Big1,
			want:             nil,
		},
		{
			name:             "nil_tip_cap",
			estimateBaseFee:  common.Big1,
			suggestGasTipCap: nil,
			want:             nil,
		},
		{
			name:             "minimum_values",
			estimateBaseFee:  bigMinBaseFee,
			suggestGasTipCap: bigMinGasTip,
			cfg:              testCfg,
			want: &PriceOptions{
				Slow: newPrice(
					minGasTip,
					uint64(minBaseFee+minGasTip),
				),
				Normal: newPrice(
					minGasTip,
					uint64(minBaseFee+minGasTip),
				),
				Fast: newPrice(
					minGasTip,
					(fastFeeNumerator*uint64(minBaseFee)/feeDenominator)+(fastFeeNumerator*uint64(minGasTip)/feeDenominator),
				),
			},
		},
		{
			name:             "maximum_values_1_slow_perc_2_fast_perc",
			estimateBaseFee:  new(big.Int).SetUint64(maxNormalBaseFee),
			suggestGasTipCap: new(big.Int).SetUint64(maxNormalGasTip),
			cfg: PriceOptionConfig{
				SlowFeePercentage: 100,
				FastFeePercentage: 200,
				MaxBaseFee:        100 * params.GWei,
				MaxTip:            20 * params.GWei,
			},
			want: &PriceOptions{
				Slow: newPrice(
					20*params.GWei,
					120*params.GWei,
				),
				Normal: newPrice(
					20*params.GWei,
					120*params.GWei,
				),
				Fast: newPrice(
					40*params.GWei,
					240*params.GWei,
				),
			},
		},
		{
			name:             "maximum_values",
			cfg:              testCfg,
			estimateBaseFee:  new(big.Int).SetUint64(maxNormalBaseFee),
			suggestGasTipCap: new(big.Int).SetUint64(maxNormalGasTip),
			want: &PriceOptions{
				Slow: newPrice(
					(slowFeeNumerator*maxNormalGasTip)/feeDenominator,
					(slowFeeNumerator*maxNormalBaseFee)/feeDenominator+(slowFeeNumerator*maxNormalGasTip)/feeDenominator,
				),
				Normal: newPrice(
					maxNormalGasTip,
					maxNormalBaseFee+maxNormalGasTip,
				),
				Fast: newPrice(
					(fastFeeNumerator*maxNormalGasTip)/feeDenominator,
					(fastFeeNumerator*maxNormalBaseFee)/feeDenominator+(fastFeeNumerator*maxNormalGasTip)/feeDenominator,
				),
			},
		},
		{
			name:             "double_maximum_values",
			estimateBaseFee:  big.NewInt(2 * int64(maxNormalBaseFee)),
			suggestGasTipCap: big.NewInt(2 * int64(maxNormalGasTip)),
			cfg:              testCfg,
			want: &PriceOptions{
				Slow: newPrice(
					(slowFeeNumerator*maxNormalGasTip)/feeDenominator,
					(slowFeeNumerator*maxNormalBaseFee)/feeDenominator+(slowFeeNumerator*maxNormalGasTip)/feeDenominator,
				),
				Normal: newPrice(
					maxNormalGasTip,
					maxNormalBaseFee+maxNormalGasTip,
				),
				Fast: newPrice(
					(fastFeeNumerator*2*maxNormalGasTip)/feeDenominator,
					(fastFeeNumerator*2*maxNormalBaseFee)/feeDenominator+(fastFeeNumerator*2*maxNormalGasTip)/feeDenominator,
				),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			backend := &testSuggestPriceOptionsBackend{
				estimateBaseFee:  test.estimateBaseFee,
				suggestGasTipCap: test.suggestGasTipCap,
				cfg:              test.cfg,
			}
			api := NewEthereumAPI(backend)

			got, err := api.SuggestPriceOptions(context.Background())
			require.NoError(err)
			require.Equal(test.want, got)
		})
	}
}

func newPrice(gasTip, gasFee uint64) *Price {
	return &Price{
		GasTip: (*hexutil.Big)(new(big.Int).SetUint64(gasTip)),
		GasFee: (*hexutil.Big)(new(big.Int).SetUint64(gasFee)),
	}
}
