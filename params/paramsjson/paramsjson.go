// Package paramsjson provides JSON unmarshalling for `params` types that depend
// on the `modules` package. This avoids `params` depending on `modules`, even
// transitively, which would result in a circular dependency.
//
// This package doesn't export any identifiers. It should instead be blank _
// imported to register its unmarshallers, similarly to SQL drivers.
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
	params.MustRegisterJSONUnmarshalers(
		unmarshalChainConfig,
		unmarshalChainConfigWithUpgrades,
		unmarshalUpgradeConfig,
		func(data []byte, v *params.PrecompileUpgrade) error {
			return json.Unmarshal(data, (*precompileUpgrade)(v))
		},
	)
}

func unmarshalChainConfig(data []byte, v *params.ChainConfig) error {
	return unmarshalChainConfigAndUpgrades(data, v, nil, "")
}

func unmarshalChainConfigWithUpgrades(data []byte, v *params.ChainConfigWithUpgradesJSON) error {
	const fldName = "UpgradeConfig"
	_ = v.UpgradeConfig // if changing this then change the line above too

	tStruct := reflect.TypeOf(v).Elem()
	fld, ok := tStruct.FieldByName(fldName)
	if !ok {
		// If this happens then the constant `fldName` is of a different name to the actual struct field used below.
		return fmt.Errorf("BUG: %T(%v).FieldByName(%q) returned false", tStruct, tStruct, fldName)
	}
	return unmarshalChainConfigAndUpgrades(data, &v.ChainConfig, &v.UpgradeConfig, strings.Split(fld.Tag.Get("json"), ",")[0])
}

func unmarshalChainConfigAndUpgrades(data []byte, cc *params.ChainConfig, upgrades *params.UpgradeConfig, upgradesJSONField string) error {
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
