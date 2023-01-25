package utils

import "math/big"


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
