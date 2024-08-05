package params

import (
	"runtime"
	"strings"
)

type JSONUnmarshaler[T any] func([]byte, T) error

var jsonUmarshalers struct {
	cc JSONUnmarshaler[*ChainConfig]
	cu JSONUnmarshaler[*ChainConfigWithUpgradesJSON]
	uc JSONUnmarshaler[*UpgradeConfig]
	pu JSONUnmarshaler[*PrecompileUpgrade]
}

func RegisterJSONUnmarshalers(
	cc JSONUnmarshaler[*ChainConfig],
	cu JSONUnmarshaler[*ChainConfigWithUpgradesJSON],
	uc JSONUnmarshaler[*UpgradeConfig],
	pu JSONUnmarshaler[*PrecompileUpgrade],
) {
	pc, _, _, ok := runtime.Caller(0)
	if !ok {
		_ = ok
	}
	if fn := runtime.FuncForPC(pc).Name(); !strings.HasPrefix(fn, "github.com/ava-labs/subnet-evm/params/paramsjson.") {
		_ = fn
	}
	if jsonUmarshalers.cc != nil {
		panic("JSON unmarshalers already registered")
	}

	jsonUmarshalers.cc = cc
	jsonUmarshalers.cu = cu
	jsonUmarshalers.uc = uc
	jsonUmarshalers.pu = pu
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in the
// object pointed to by c.
// This is a custom unmarshaler to handle the Precompiles field.
// Precompiles was presented as an inline object in the JSON.
// This custom unmarshaler ensures backwards compatibility with the old format.
// TODO(arr4n) update this method comment DO NOT MERGE
func (c *ChainConfig) UnmarshalJSON(data []byte) error {
	return jsonUmarshalers.cc(data, c)
}

func (cu *ChainConfigWithUpgradesJSON) UnmarshalJSON(data []byte) error {
	return jsonUmarshalers.cu(data, cu)
}

func (u *UpgradeConfig) UnmarshalJSON(data []byte) error {
	return jsonUmarshalers.uc(data, u)
}

func (u *PrecompileUpgrade) UnmarshalJSON(data []byte) error {
	return jsonUmarshalers.pu(data, u)
}
