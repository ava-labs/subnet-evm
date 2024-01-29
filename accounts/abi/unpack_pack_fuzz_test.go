// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build fuzz
// +build fuzz

package abi

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"rogchap.com/v8go"
)

const jsEncodeValueFn = `
function encodeValue(name, value) {
	var abiCoder = new ethers.AbiCoder();
	var encodedValue = abiCoder.encode([name], [value]);
	return encodedValue;
}
`

const jsEncodeCallFn = `encodeValue("%s", %s)`

var (
	codeCache          *v8go.CompilerCachedData
	precompiledScripts string
	cacheLock          sync.Mutex
)

func initializeJSVM(t *testing.T, iso *v8go.Isolate, cache *v8go.CompilerCachedData) (*v8go.Context, *v8go.UnboundScript, error) {
	ctx := v8go.NewContext(iso)

	// Inject ethers library into the runtime
	// Compile the ethers subset script
	script, err := iso.CompileUnboundScript(
		precompiledScripts,
		"ethers-lib.js",
		v8go.CompileOptions{
			CachedData: &v8go.CompilerCachedData{Bytes: codeCache.Bytes},
		})
	if err != nil {
		return nil, nil, err
	}

	return ctx, script, nil
}

// basic types

func FuzzPackUint8(f *testing.F) {
	f.Fuzz(func(t *testing.T, val uint8) {
		testPackType(t, "uint8", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackInt8(f *testing.F) {
	f.Fuzz(func(t *testing.T, val int8) {
		testPackType(t, "int8", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackUint32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val uint32) {
		testPackType(t, "uint32", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackInt32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val int32) {
		testPackType(t, "int32", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackUint64(f *testing.F) {
	f.Fuzz(func(t *testing.T, val uint64) {
		testPackType(t, "uint64", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackInt64(f *testing.F) {
	f.Fuzz(func(t *testing.T, val int64) {
		testPackType(t, "int64", val, wrapQuote(fmt.Sprintf("%d", val)))
	})
}

func FuzzPackUint256(f *testing.F) {
	f.Fuzz(func(t *testing.T, val []byte) {
		bigIntVal := new(big.Int).SetBytes(val)
		if bigIntVal.Cmp(MaxUint256) > 0 {
			t.Skip()
		}
		testPackType(t, "uint256", bigIntVal, wrapQuote(bigIntVal.String()))
	})
}

func FuzzPackInt256(f *testing.F) {
	f.Fuzz(func(t *testing.T, val []byte) {
		bigIntVal := new(big.Int).SetBytes(val)
		// skip if bigInt exceeds int256 bounds
		if bigIntVal.Cmp(MaxInt256) > 0 {
			t.Skip()
		}
		testPackType(t, "int256", bigIntVal, wrapQuote(bigIntVal.String()))
	})
}

func FuzzPackBool(f *testing.F) {
	f.Fuzz(func(t *testing.T, val bool) {
		testPackType(t, "bool", val, fmt.Sprintf("%v", val))
	})
}

func FuzzPackAddress(f *testing.F) {
	f.Fuzz(func(t *testing.T, val []byte) {
		addressVal := common.BytesToAddress(val)
		testPackType(t, "address", addressVal, wrapQuote(addressVal.String()))
	})
}

// array types

func FuzzPackArrayUint8(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 uint8, val2 uint8) {
		val := []uint8{val1, val2}
		testPackType(t, "uint8[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayInt8(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 int8, val2 int8) {
		val := []int8{val1, val2}
		testPackType(t, "int8[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayUint32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 uint32, val2 uint32) {
		val := []uint32{val1, val2}
		testPackType(t, "uint32[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayInt32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 int32, val2 int32) {
		val := []int32{val1, val2}
		testPackType(t, "int32[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayUint64(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 uint64, val2 uint64) {
		val := []uint64{val1, val2}
		testPackType(t, "uint64[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayInt64(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 int64, val2 int64) {
		val := []int64{val1, val2}
		testPackType(t, "int64[]", val, fmt.Sprintf("[%d,%d]", val1, val2))
	})
}

func FuzzPackArrayUint256(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 []byte, val2 []byte) {
		bigIntVal1 := new(big.Int).SetBytes(val1)
		bigIntVal2 := new(big.Int).SetBytes(val2)
		if bigIntVal1.Cmp(MaxUint256) > 0 || bigIntVal2.Cmp(MaxUint256) > 0 {
			t.Skip()
		}
		val := []*big.Int{bigIntVal1, bigIntVal2}
		testPackType(t, "uint256[]", val, fmt.Sprintf("['%s','%s']", bigIntVal1.String(), bigIntVal2.String()))
	})
}

func FuzzPackArrayInt256(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 []byte, val2 []byte) {
		bigIntVal1 := new(big.Int).SetBytes(val1)
		bigIntVal2 := new(big.Int).SetBytes(val2)
		if bigIntVal1.Cmp(MaxInt256) > 0 || bigIntVal2.Cmp(MaxInt256) > 0 {
			t.Skip()
		}
		val := []*big.Int{bigIntVal1, bigIntVal2}
		testPackType(t, "int256[]", val, fmt.Sprintf("['%s','%s']", bigIntVal1.String(), bigIntVal2.String()))
	})
}

func FuzzPackArrayBool(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 bool, val2 bool) {
		val := []bool{val1, val2}
		testPackType(t, "bool[]", val, fmt.Sprintf("[%v,%v]", val1, val2))
	})
}

func FuzzPackArrayAddress(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 []byte, val2 []byte) {
		addressVal1 := common.BytesToAddress(val1)
		addressVal2 := common.BytesToAddress(val2)
		val := []common.Address{addressVal1, addressVal2}
		testPackType(t, "address[]", val, fmt.Sprintf("['%s','%s']", addressVal1.String(), addressVal2.String()))
	})
}

func FuzzPackString(f *testing.F) {
	f.Fuzz(func(t *testing.T, val string) {
		t.Log("Testing string: ", []byte(val))
		testPackType(t, "string", val, strconv.Quote(val))
	})
}

func FuzzPackArrayString(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 string, val2 string) {
		val := []string{val1, val2}
		testPackType(t, "string[]", val, fmt.Sprintf("[%v,%v]", strconv.Quote(val1), strconv.Quote(val2)))
	})
}

func FuzzPackBytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, val []byte) {
		testPackType(t, "bytes", val, wrapQuote("0x"+common.Bytes2Hex(val)))
	})
}

func FuzzPackFixedBytes32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val []byte) {
		var array32 [32]byte
		copy(array32[:], val)
		testPackType(t, "bytes32", array32, wrapQuote("0x"+common.Bytes2Hex(array32[:])))
	})
}

func FuzzPackArrayBytes(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 []byte, val2 []byte) {
		val := [][]byte{val1, val2}
		testPackType(t, "bytes[]", val, fmt.Sprintf("[%v,%v]", wrapQuote("0x"+common.Bytes2Hex(val1)), wrapQuote("0x"+common.Bytes2Hex(val2))))
	})
}

func FuzzPackFixedArrayBytes32(f *testing.F) {
	f.Fuzz(func(t *testing.T, val1 []byte, val2 []byte) {
		var array32_1 [32]byte
		var array32_2 [32]byte
		copy(array32_1[:], val1)
		copy(array32_2[:], val2)
		val := [][32]byte{array32_1, array32_2}
		testPackType(t, "bytes32[]", val, fmt.Sprintf("[%v,%v]", wrapQuote("0x"+common.Bytes2Hex(array32_1[:])), wrapQuote("0x"+common.Bytes2Hex(array32_2[:]))))
	})
}

func preCacheScript() error {
	ethersBytes, err := os.ReadFile("ethers.min.js")
	if err != nil {
		panic(fmt.Errorf("Error reading ethers.min.js file: %w", err))
	}
	globalInjectionScript := "var global = (function(){ return this; }).call(null);"

	precompiledScripts = globalInjectionScript + string(ethersBytes) + jsEncodeValueFn

	iso := v8go.NewIsolate()
	defer iso.Dispose()
	script, err := iso.CompileUnboundScript(precompiledScripts, "ethers.js", v8go.CompileOptions{})
	if err != nil {
		return fmt.Errorf("Error compiling ethers.js script: %w", err)
	}

	codeCache = script.CreateCodeCache()
	return nil
}

func testPackType(t *testing.T, typeName string, val interface{}, stringVal string) {
	t.Logf("Testing type \"%s\", value: %v, stringVal: '%s'", typeName, val, stringVal)
	cacheLock.Lock()
	if codeCache == nil {
		err := preCacheScript()
		if err != nil {
			cacheLock.Unlock()
			t.Fatal(err)
		}
	}
	cacheLock.Unlock()
	require := require.New(t)
	abiType, err := NewType(typeName, "", nil)
	require.NoError(err)
	abiArg := Argument{
		Name: "test",
		Type: abiType,
	}
	abiArgs := Arguments{abiArg}
	packed, err := abiArgs.Pack(val)
	require.NoError(err)
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	// Initialize the JS VM
	ctx, ethersScript, err := initializeJSVM(t, iso, codeCache)
	defer ctx.Close()
	require.NoError(err)
	// Run the ethers script to load
	_, err = ethersScript.Run(ctx)
	require.NoError(err)
	// Run the encode script
	script := fmt.Sprintf(jsEncodeCallFn, abiType.String(), stringVal)

	t.Logf("Running script: %s", script)
	result, err := ctx.RunScript(script, "main.js")
	require.NoError(err)

	require.Equal("0x"+common.Bytes2Hex(packed), result.String())
}

func wrapQuote(s string) string {
	return fmt.Sprintf("%q", s)
}
