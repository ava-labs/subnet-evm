package ticks

import (
	"errors"
	"fmt"
	"math/big"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ethereum/go-ethereum/common"
)

func GetPrevTick(bibliophile b.BibliophileClient, input GetPrevTickInput) (*big.Int, error) {
	if input.Tick.Sign() == 0 {
		return nil, errors.New("tick price cannot be zero")
	}
	if input.IsBid {
		currentTick := bibliophile.GetBidsHead(input.Amm)
		if input.Tick.Cmp(currentTick) >= 0 {
			return nil, fmt.Errorf("tick %d is greater than or equal to bidsHead %d", input.Tick, currentTick)
		}
		for {
			nextTick := bibliophile.GetNextBidPrice(input.Amm, currentTick)
			if nextTick.Cmp(input.Tick) <= 0 {
				return currentTick, nil
			}
			currentTick = nextTick
		}
	}
	currentTick := bibliophile.GetAsksHead(input.Amm)
	if currentTick.Sign() == 0 {
		return nil, errors.New("asksHead is zero")
	}
	if input.Tick.Cmp(currentTick) <= 0 {
		return nil, fmt.Errorf("tick %d is less than or equal to asksHead %d", input.Tick, currentTick)
	}
	for {
		nextTick := bibliophile.GetNextAskPrice(input.Amm, currentTick)
		if nextTick.Cmp(input.Tick) >= 0 || nextTick.Sign() == 0 {
			return currentTick, nil
		}
		currentTick = nextTick
	}
}

func SampleImpactBid(bibliophile b.BibliophileClient, ammAddress common.Address) *big.Int {
	impactMarginNotional := bibliophile.GetImpactMarginNotional(ammAddress)
	if impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	return _sampleImpactBid(bibliophile, ammAddress, impactMarginNotional)
}

func _sampleImpactBid(bibliophile b.BibliophileClient, ammAddress common.Address, _impactMarginNotional *big.Int) *big.Int {
	if _impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	impactMarginNotional := new(big.Int).Mul(_impactMarginNotional, big.NewInt(1e12))
	accNotional := big.NewInt(0) // 18 decimals
	accBaseQ := big.NewInt(0)    // 18 decimals
	tick := bibliophile.GetBidsHead(ammAddress)
	for tick.Sign() != 0 {
		amount := bibliophile.GetBidSize(ammAddress, tick)
		accumulator := new(big.Int).Add(accNotional, hu.Div1e6(big.NewInt(0).Mul(amount, tick)))
		if accumulator.Cmp(impactMarginNotional) >= 0 {
			break
		}
		accNotional = accumulator
		accBaseQ.Add(accBaseQ, amount)
		tick = bibliophile.GetNextBidPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0)
	}
	baseQAtTick := new(big.Int).Div(hu.Mul1e6(new(big.Int).Sub(impactMarginNotional, accNotional)), tick)
	return new(big.Int).Div(hu.Mul1e6(impactMarginNotional), new(big.Int).Add(baseQAtTick, accBaseQ)) // return value is in 6 decimals
}

func SampleImpactAsk(bibliophile b.BibliophileClient, ammAddress common.Address) *big.Int {
	impactMarginNotional := bibliophile.GetImpactMarginNotional(ammAddress)
	if impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	return _sampleImpactAsk(bibliophile, ammAddress, impactMarginNotional)
}

func _sampleImpactAsk(bibliophile b.BibliophileClient, ammAddress common.Address, _impactMarginNotional *big.Int) *big.Int {
	if _impactMarginNotional.Sign() == 0 {
		return big.NewInt(0)
	}
	impactMarginNotional := new(big.Int).Mul(_impactMarginNotional, big.NewInt(1e12))
	tick := bibliophile.GetAsksHead(ammAddress)
	accNotional := big.NewInt(0) // 18 decimals
	accBaseQ := big.NewInt(0)    // 18 decimals
	for tick.Sign() != 0 {
		amount := bibliophile.GetAskSize(ammAddress, tick)
		accumulator := new(big.Int).Add(accNotional, hu.Div1e6(big.NewInt(0).Mul(amount, tick)))
		if accumulator.Cmp(impactMarginNotional) >= 0 {
			break
		}
		accNotional = accumulator
		accBaseQ.Add(accBaseQ, amount)
		tick = bibliophile.GetNextAskPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0)
	}
	baseQAtTick := new(big.Int).Div(hu.Mul1e6(new(big.Int).Sub(impactMarginNotional, accNotional)), tick)
	return new(big.Int).Div(hu.Mul1e6(impactMarginNotional), new(big.Int).Add(baseQAtTick, accBaseQ)) // return value is in 6 decimals
}

func GetBaseQuote(bibliophile b.BibliophileClient, ammAddress common.Address, quoteAssetQuantity *big.Int) *big.Int {
	if quoteAssetQuantity.Sign() > 0 { // get the qoute to long quoteQuantity dollars
		return _sampleImpactAsk(bibliophile, ammAddress, quoteAssetQuantity)
	}
	// get the qoute to short quoteQuantity dollars
	return _sampleImpactBid(bibliophile, ammAddress, new(big.Int).Neg(quoteAssetQuantity))
}

func GetQuote(bibliophile b.BibliophileClient, ammAddress common.Address, baseAssetQuantity *big.Int) *big.Int {
	if baseAssetQuantity.Sign() > 0 {
		return _sampleAsk(bibliophile, ammAddress, baseAssetQuantity)
	}
	return _sampleBid(bibliophile, ammAddress, new(big.Int).Neg(baseAssetQuantity))
}

func _sampleAsk(bibliophile b.BibliophileClient, ammAddress common.Address, baseAssetQuantity *big.Int) *big.Int {
	if baseAssetQuantity.Sign() <= 0 {
		return big.NewInt(0)
	}
	tick := bibliophile.GetAsksHead(ammAddress)
	accNotional := big.NewInt(0) // 18 decimals
	accBaseQ := big.NewInt(0)    // 18 decimals
	for tick.Sign() != 0 {
		amount := bibliophile.GetAskSize(ammAddress, tick)
		accumulator := hu.Add(accBaseQ, amount)
		if accumulator.Cmp(baseAssetQuantity) >= 0 {
			break
		}
		accNotional.Add(accNotional, hu.Div1e6(hu.Mul(amount, tick)))
		accBaseQ = accumulator
		tick = bibliophile.GetNextAskPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0) // insufficient liquidity
	}
	notionalAtTick := hu.Div1e6(hu.Mul(hu.Sub(baseAssetQuantity, accBaseQ), tick))
	return hu.Div(hu.Mul1e6(hu.Add(accNotional, notionalAtTick)), baseAssetQuantity) // return value is in 6 decimals
}

func _sampleBid(bibliophile b.BibliophileClient, ammAddress common.Address, baseAssetQuantity *big.Int) *big.Int {
	if baseAssetQuantity.Sign() <= 0 {
		return big.NewInt(0)
	}
	tick := bibliophile.GetBidsHead(ammAddress)
	accNotional := big.NewInt(0) // 18 decimals
	accBaseQ := big.NewInt(0)    // 18 decimals
	for tick.Sign() != 0 {
		amount := bibliophile.GetBidSize(ammAddress, tick)
		accumulator := hu.Add(accBaseQ, amount)
		if accumulator.Cmp(baseAssetQuantity) >= 0 {
			break
		}
		accNotional.Add(accNotional, hu.Div1e6(hu.Mul(amount, tick)))
		accBaseQ = accumulator
		tick = bibliophile.GetNextBidPrice(ammAddress, tick)
	}
	if tick.Sign() == 0 {
		return big.NewInt(0) // insufficient liquidity
	}
	notionalAtTick := hu.Div1e6(hu.Mul(hu.Sub(baseAssetQuantity, accBaseQ), tick))
	return hu.Div(hu.Mul1e6(hu.Add(accNotional, notionalAtTick)), baseAssetQuantity) // return value is in 6 decimals
}
