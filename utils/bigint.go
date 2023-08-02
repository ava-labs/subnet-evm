package utils

import (
	"fmt"
	"math"
	"math/big"
)

func BigIntMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) == 1 {
		return big.NewInt(0).Set(x)
	} else {
		return big.NewInt(0).Set(y)
	}
}

func BigIntMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) == -1 {
		return big.NewInt(0).Set(x)
	} else {
		return big.NewInt(0).Set(y)
	}
}

// BigIntMinAbs calculates minimum of absolute values
func BigIntMinAbs(x, y *big.Int) *big.Int {
	xAbs := big.NewInt(0).Abs(x)
	yAbs := big.NewInt(0).Abs(y)
	if xAbs.Cmp(yAbs) == -1 {
		return big.NewInt(0).Set(xAbs)
	} else {
		return big.NewInt(0).Set(yAbs)
	}
}

func BigIntToDecimal(x *big.Int, scale int, decimals int) string {
	// Create big.Float from x
	f := new(big.Float).SetInt(x)

	// Create big.Float for scale and set its value
	s := new(big.Float)
	s.SetInt(big.NewInt(int64(1)))
	for i := 0; i < scale; i++ {
		s.Mul(s, big.NewFloat(10))
	}

	// Divide x by scale
	f.Quo(f, s)

	// Setting precision and converting big.Float to string
	str := fmt.Sprintf("%.*f", decimals, f)

	return str
}

func BigIntToFloat(number *big.Int, scale int8) float64 {
	float, _ := new(big.Float).Quo(new(big.Float).SetInt(number), big.NewFloat(math.Pow10(int(scale)))).Float64()
	return float
}
