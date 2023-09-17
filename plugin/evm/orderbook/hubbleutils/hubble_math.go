package hubbleutils

import (
	"math/big"
)

var (
	ONE_E_6  = big.NewInt(1e6)
	ONE_E_12 = big.NewInt(1e12)
	ONE_E_18 = big.NewInt(1e18)
)

func Mul1e6(a *big.Int) *big.Int {
	return Mul(a, ONE_E_6)
}

func Div1e6(a *big.Int) *big.Int {
	return Div(a, ONE_E_6)
}

func Mul1e18(a *big.Int) *big.Int {
	return Mul(a, ONE_E_18)
}

func Div1e18(a *big.Int) *big.Int {
	return Div(a, ONE_E_18)
}

func Add(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func Sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func Mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func Div(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}
