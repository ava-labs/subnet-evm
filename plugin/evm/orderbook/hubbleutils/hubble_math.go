package hubbleutils

import (
	"fmt"
	"math/big"

	// "github.com/ava-labs/subnet-evm/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ONE_E_6  = big.NewInt(1e6)
	ONE_E_12 = big.NewInt(1e12)
	ONE_E_18 = big.NewInt(1e18)
)

func Add1e6(a *big.Int) *big.Int {
	return Add(a, ONE_E_6)
}

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

func Abs(a *big.Int) *big.Int {
	return new(big.Int).Abs(a)
}

func RoundOff(a, b *big.Int) *big.Int {
	return Mul(Div(a, b), b)
}

func Mod(a, b *big.Int) *big.Int {
	return new(big.Int).Mod(a, b)
}

func Neg(a *big.Int) *big.Int {
	return new(big.Int).Neg(a)
}

func Scale(a *big.Int, decimals uint8) *big.Int {
	return Mul(a, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
}

func Unscale(a *big.Int, decimals uint8) *big.Int {
	return Div(a, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
}

func ECRecover(data, sign hexutil.Bytes) (common.Address, error) {
	sig := make([]byte, len(sign))
	copy(sig, sign)

	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.Ecrecover(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(common.LeftPadBytes(crypto.Keccak256(rpk[1:])[12:], 32)), nil
}
