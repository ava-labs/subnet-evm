// Package paramsjson provides JSON unmarshalling for `params` types that depend
// on the `modules` package. This avoids `params` depending on `modules`, even
// transitively, which would result in a circular dependency.
//
// Typically there is no need to call this package directly. It should instead
// be blank _ imported to register its unmarshallers, similarly to SQL drivers.
package paramsjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/modules"
)

func init() {
	// The parameters of RegisterJSONUmarshalers are all generic but the
	// specific types are unambiguous so we don't need to specify them here.
	params.RegisterJSONUnmarshalers(Unmarshal, Unmarshal, Unmarshal, Unmarshal)
}

// An Unmarshaler can have JSON unmarshalled into it by [Unmarshal].
type Unmarshaler interface {
	*params.ChainConfig | *params.ChainConfigWithUpgradesJSON | *params.UpgradeConfig | *params.PrecompileUpgrade
}

// Unmarshal is a drop-in replacement for [json.Unmarshal].
func Unmarshal[T Unmarshaler](data []byte, v T) error {
	switch v := any(v).(type) {
	case *params.ChainConfig:
		return unmarshalChainConfig(data, v, nil, "")

	case *params.ChainConfigWithUpgradesJSON:
		const fldName = "UpgradeConfig"

		tStruct := reflect.TypeOf(v).Elem()
		fld, ok := tStruct.FieldByName(fldName)
		if !ok {
			// If this happens then the constant `fldName` is of a different name to the actual struct field used below.
			return fmt.Errorf("BUG: %T(%v).FieldByName(%q) returned false", tStruct, tStruct, fldName)
		}
		return unmarshalChainConfig(data, &v.ChainConfig, &v.UpgradeConfig, strings.Split(fld.Tag.Get("json"), ",")[0])

	case *params.UpgradeConfig:
		return unmarshalUpgradeConfig(data, v)

	case *params.PrecompileUpgrade:
		return json.Unmarshal(data, (*precompileUpgrade)(v))

	default:
		// If this happens then the Unmarshaler interface has been modified but the above cases haven't been.
		return fmt.Errorf("unsupported type %T", v)
	}
}

func unmarshalChainConfig(data []byte, cc *params.ChainConfig, upgrades *params.UpgradeConfig, upgradesJSONField string) error {
	type withoutMethods *params.ChainConfig // circumvents UnmarshalJSON() method, which always returns an error
	if err := json.Unmarshal(data, withoutMethods(cc)); err != nil {
		return err
	}

	byField := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &byField); err != nil {
		return err
	}

	if cc.GenesisPrecompiles == nil {
		cc.GenesisPrecompiles = make(params.Precompiles)
	}
	for fld, buf := range byField {
		switch mod, ok := modules.GetPrecompileModule(fld); {
		case ok:
			conf := mod.MakeConfig()
			if err := json.Unmarshal(buf, conf); err != nil {
				return err
			}
			cc.GenesisPrecompiles[mod.ConfigKey] = conf

		case fld == upgradesJSONField && upgrades != nil:
			if err := unmarshalUpgradeConfig(buf, upgrades); err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalUpgradeConfig(data []byte, uc *params.UpgradeConfig) error {
	byField := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &byField); err != nil {
		return err
	}

	precompileT := reflect.TypeOf([]params.PrecompileUpgrade{})

	config := reflect.ValueOf(uc).Elem()
	for i := 0; i < config.NumField(); i++ {
		fld := config.Type().FieldByIndex([]int{i})
		jsonFld := strings.Split(fld.Tag.Get("json"), ",")[0]
		if _, ok := byField[jsonFld]; !ok {
			continue
		}

		var jsonInto any
		switch fldVal := config.Field(i); {
		case fld.Type == precompileT:
			var out []precompileUpgrade
			jsonInto = &out
			defer func() {
				uc.PrecompileUpgrades = *(*[]params.PrecompileUpgrade)(unsafe.Pointer(&out))
			}()

		case fld.Type.Kind() == reflect.Slice:
			jsonInto = fldVal.Addr().Interface()

		case fld.Type.Kind() == reflect.Pointer:
			if fldVal.IsNil() {
				fldVal.Set(reflect.New(fld.Type.Elem()))
			}
			jsonInto = fldVal.Interface()

		default:
			return fmt.Errorf("unsupported field %T.%s", uc, fld.Name)
		}

		if err := json.Unmarshal(byField[jsonFld], jsonInto); err != nil {
			return fmt.Errorf("json.Unmarshal field %q: %v", jsonFld, err)
		}
	}

	return nil
}

type precompileUpgrade params.PrecompileUpgrade

var _ json.Unmarshaler = (*precompileUpgrade)(nil)

func (u *precompileUpgrade) UnmarshalJSON(data []byte) error {
	byField := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &byField); err != nil {
		return err
	}
	if n := len(byField); n != 1 {
		return fmt.Errorf("unmarshalling %T; got %d JSON fields; MUST be exactly one (name of precompile module)", &params.PrecompileUpgrade{}, n)
	}

	for key, value := range byField {
		mod, ok := modules.GetPrecompileModule(key)
		if !ok {
			return fmt.Errorf("unknown precompile config: %s", key)
		}
		config := mod.MakeConfig()
		if err := json.Unmarshal(value, config); err != nil {
			return err
		}
		u.Config = config
	}
	return nil
}