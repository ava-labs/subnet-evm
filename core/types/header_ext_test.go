// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ava-labs/libevm/common"
	ethtypes "github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderRLP(t *testing.T) {
	t.Parallel()

	got := testHeaderEncodeDecode(t, rlp.EncodeToBytes, rlp.DecodeBytes)

	// Golden data from original coreth implementation, before integration of
	// libevm. WARNING: changing these values can break backwards compatibility
	// with extreme consequences as block-hash calculation may break.
	const (
		wantHex     = "f90212a00100000000000000000000000000000000000000000000000000000000000000a00200000000000000000000000000000000000000000000000000000000000000940300000000000000000000000000000000000000a00400000000000000000000000000000000000000000000000000000000000000a00500000000000000000000000000000000000000000000000000000000000000a00600000000000000000000000000000000000000000000000000000000000000b901000700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008090a0b0c0da00e00000000000000000000000000000000000000000000000000000000000000880f0000000000000010151213a01400000000000000000000000000000000000000000000000000000000000000"
		wantHashHex = "2453a240c1cfa4eca66bf39db950d5bd57f5e94ffabf9d800497ace33c2a5927"
	)

	assert.Equal(t, wantHex, hex.EncodeToString(got), "Header RLP")

	header, _ := headerWithNonZeroFields()
	gotHashHex := header.Hash().Hex()
	assert.Equal(t, "0x"+wantHashHex, gotHashHex, "Header.Hash()")
}

func TestHeaderJSON(t *testing.T) {
	t.Parallel()

	// Note we ignore the returned encoded bytes because we don't
	// need to compare them to a JSON gold standard.
	_ = testHeaderEncodeDecode(t, json.Marshal, json.Unmarshal)
}

func testHeaderEncodeDecode(
	t *testing.T,
	encode func(any) ([]byte, error),
	decode func([]byte, any) error,
) (encoded []byte) {
	t.Helper()

	input, _ := headerWithNonZeroFields() // the Header carries the HeaderExtra so we can ignore it
	encoded, err := encode(input)
	require.NoError(t, err, "encode")

	gotHeader := new(Header)
	err = decode(encoded, gotHeader)
	require.NoError(t, err, "decode")
	gotExtra := GetHeaderExtra(gotHeader)

	wantHeader, wantExtra := headerWithNonZeroFields()
	wantHeader.WithdrawalsHash = nil
	assert.Equal(t, wantHeader, gotHeader)
	assert.Equal(t, wantExtra, gotExtra)

	return encoded
}

func TestHeaderWithNonZeroFields(t *testing.T) {
	t.Parallel()

	header, extra := headerWithNonZeroFields()
	t.Run("Header", func(t *testing.T) { allExportedFieldsSet(t, header) })
	t.Run("HeaderExtra", func(t *testing.T) { allExportedFieldsSet(t, extra) })
}

// headerWithNonZeroFields returns a [Header] and a [HeaderExtra],
// each with all fields set to non-zero values.
// The [HeaderExtra] extra payload is set in the [Header] via [SetHeaderExtra].
//
// NOTE: They can be used to demonstrate that RLP and JSON round-trip encoding
// can recover all fields, but not that the encoded format is correct. This is
// very important as the RLP encoding of a [Header] defines its hash.
func headerWithNonZeroFields() (*Header, *HeaderExtra) {
	header := &ethtypes.Header{
		ParentHash:       common.Hash{1},
		UncleHash:        common.Hash{2},
		Coinbase:         common.Address{3},
		Root:             common.Hash{4},
		TxHash:           common.Hash{5},
		ReceiptHash:      common.Hash{6},
		Bloom:            Bloom{7},
		Difficulty:       big.NewInt(8),
		Number:           big.NewInt(9),
		GasLimit:         10,
		GasUsed:          11,
		Time:             12,
		Extra:            []byte{13},
		MixDigest:        common.Hash{14},
		Nonce:            BlockNonce{15},
		BaseFee:          big.NewInt(16),
		WithdrawalsHash:  &common.Hash{17},
		BlobGasUsed:      ptrTo(uint64(18)),
		ExcessBlobGas:    ptrTo(uint64(19)),
		ParentBeaconRoot: &common.Hash{20},
	}
	extra := &HeaderExtra{
		BlockGasCost: big.NewInt(21),
	}
	SetHeaderExtra(header, extra)
	return header, extra
}

func allExportedFieldsSet[T interface {
	ethtypes.Header | HeaderExtra
}](t *testing.T, x *T) {
	// We don't test for nil pointers because we're only confirming that
	// test-input data is well-formed. A panic due to a dereference will be
	// reported anyway.

	v := reflect.ValueOf(*x)
	for i := range v.Type().NumField() {
		field := v.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		t.Run(field.Name, func(t *testing.T) {
			switch f := v.Field(i).Interface().(type) {
			case common.Hash:
				assertNonZero(t, f)
			case common.Address:
				assertNonZero(t, f)
			case BlockNonce:
				assertNonZero(t, f)
			case Bloom:
				assertNonZero(t, f)
			case uint64:
				assertNonZero(t, f)
			case *big.Int:
				assertNonZero(t, f)
			case *common.Hash:
				assertNonZero(t, f)
			case *uint64:
				assertNonZero(t, f)
			case []uint8:
				assert.NotEmpty(t, f)
			default:
				t.Errorf("Field %q has unsupported type %T", field.Name, f)
			}
		})
	}
}

func assertNonZero[T interface {
	common.Hash | common.Address | BlockNonce | uint64 | Bloom |
		*big.Int | *common.Hash | *uint64
}](t *testing.T, v T) {
	t.Helper()
	var zero T
	if v == zero {
		t.Errorf("must not be zero value for %T", v)
	}
}

func ptrTo[T any](x T) *T { return &x }
