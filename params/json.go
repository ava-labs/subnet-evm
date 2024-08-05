package params

import (
	"fmt"
	"runtime"
	"strings"
)

// A JSONUnmarshaler is a type-safe function for unmarshalling a JSON buffer
// into a specific type.
type JSONUnmarshaler[T any] func([]byte, T) error

var jsonUmarshalers struct {
	cc JSONUnmarshaler[*ChainConfig]
	cu JSONUnmarshaler[*ChainConfigWithUpgradesJSON]
	uc JSONUnmarshaler[*UpgradeConfig]
	pu JSONUnmarshaler[*PrecompileUpgrade]
}

// RegisterJSONUnmarshalers registers the JSON unmarshalling functions for
// various types. This allows their unmarshalling behaviour to be injected by
// the [params/paramsjson] package, which can't be directly imported as it would
// result in a circular dependency.
//
// This function SHOULD NOT be called directly. Instead, blank import the
// [params/paramsjson] package, which registers unmarshalers in its init()
// function.
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

func unmarshalJSON[T any](u JSONUnmarshaler[T], data []byte, v T) error {
	if u == nil {
		return fmt.Errorf(`%T is nil; did you remember to import _ "github.com/ava-labs/subnet-evm/params/paramsjson"`, u)
	}
	return u(data, v)
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in the
// object pointed to by c.
// This is a custom unmarshaler to handle the Precompiles field.
// Precompiles was presented as an inline object in the JSON.
// This custom unmarshaler ensures backwards compatibility with the old format.
// TODO(arr4n) update this method comment DO NOT MERGE
func (c *ChainConfig) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(jsonUmarshalers.cc, data, c)
}

func (cu *ChainConfigWithUpgradesJSON) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(jsonUmarshalers.cu, data, cu)
}

func (u *UpgradeConfig) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(jsonUmarshalers.uc, data, u)
}

func (u *PrecompileUpgrade) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(jsonUmarshalers.pu, data, u)
}
