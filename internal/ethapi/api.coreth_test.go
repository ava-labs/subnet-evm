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
					(testCfg.FastFeePercentage*uint64(minBaseFee)/feeDenominator)+(testCfg.FastFeePercentage*uint64(minGasTip)/feeDenominator),
				),
			},
		},
		{
			name:             "maximum_values_1_slow_perc_2_fast_perc",
			estimateBaseFee:  new(big.Int).SetUint64(testCfg.MaxBaseFee),
			suggestGasTipCap: new(big.Int).SetUint64(testCfg.MaxTip),
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
			estimateBaseFee:  new(big.Int).SetUint64(testCfg.MaxBaseFee),
			suggestGasTipCap: new(big.Int).SetUint64(testCfg.MaxTip),
			want: &PriceOptions{
				Slow: newPrice(
					((testCfg.SlowFeePercentage * testCfg.MaxTip) / feeDenominator),
					((testCfg.SlowFeePercentage*testCfg.MaxBaseFee)/feeDenominator)+((testCfg.SlowFeePercentage*testCfg.MaxTip)/feeDenominator),
				),
				Normal: newPrice(
					testCfg.MaxTip,
					testCfg.MaxBaseFee+testCfg.MaxTip,
				),
				Fast: newPrice(
					((testCfg.FastFeePercentage * testCfg.MaxTip) / feeDenominator),
					((testCfg.FastFeePercentage*testCfg.MaxBaseFee)/feeDenominator)+((testCfg.FastFeePercentage*testCfg.MaxTip)/feeDenominator),
				),
			},
		},
		{
			name:             "double_maximum_values",
			estimateBaseFee:  big.NewInt(2 * int64(testCfg.MaxBaseFee)),
			suggestGasTipCap: big.NewInt(2 * int64(testCfg.MaxTip)),
			cfg:              testCfg,
			want: &PriceOptions{
				Slow: newPrice(
					((testCfg.SlowFeePercentage * testCfg.MaxTip) / feeDenominator),
					((testCfg.SlowFeePercentage*testCfg.MaxBaseFee)/feeDenominator)+((testCfg.SlowFeePercentage*testCfg.MaxTip)/feeDenominator),
				),
				Normal: newPrice(
					testCfg.MaxTip,
					testCfg.MaxBaseFee+testCfg.MaxTip,
				),
				Fast: newPrice(
					((testCfg.FastFeePercentage * testCfg.MaxTip * 2) / feeDenominator),
					((testCfg.FastFeePercentage*testCfg.MaxBaseFee*2)/feeDenominator)+((testCfg.FastFeePercentage*testCfg.MaxTip*2)/feeDenominator),
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
