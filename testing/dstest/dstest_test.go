package dstest

import (
	"context"
	"fmt"
	"testing"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/testing/evmsim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate sh -c "solc --evm-version=paris --base-path=./ --include-path=./internal --combined-json=abi,bin FakeTest.t.sol | abigen --pkg dstest --combined-json=- | sed -E 's,github.com/ethereum/go-ethereum/(accounts|core)/,github.com/ava-labs/subnet-evm/\\1/,' > generated_test.go"

func TestParseLogs(t *testing.T) {
	ctx := context.Background()
	sim := evmsim.NewWithNumKeys(t, 2)

	addr, fake := evmsim.Deploy(t, sim, 0, DeployFakeTest)
	sut := New(addr)
	session := &FakeTestSession{
		Contract:     fake,
		TransactOpts: *sim.From(0),
	}

	t.Run("inherit DSTest IS_TEST constant", func(t *testing.T) {
		got := evmsim.Call(t, fake.ISTEST, nil)
		if !got {
			t.Errorf("%T.ISTEST() = false; want true", fake)
		}
	})

	tests := []struct {
		name      string
		tx        func() (*types.Transaction, error)
		want      Logs
		wantAsStr string
	}{
		{
			name: "string (foo)",
			tx: func() (*types.Transaction, error) {
				return session.LogString("foo")
			},
			want: Logs{
				{map[string]any{"arg0": "foo"}},
			},
			wantAsStr: "foo",
		},
		{
			name: "string (bar)",
			tx: func() (*types.Transaction, error) {
				return session.LogString("bar")
			},
			want: Logs{
				{map[string]any{"arg0": "bar"}},
			},
			wantAsStr: "bar",
		},
		{
			name: "named address",
			tx: func() (*types.Transaction, error) {
				return session.LogNamedAddress("Gary", sim.Addr(1))
			},
			want: Logs{
				{map[string]any{
					"key": "Gary",
					"val": sim.Addr(1),
				}},
			},
			wantAsStr: fmt.Sprintf("Gary = %v", sim.Addr(1)),
		},
		{
			name: "non-DSTest log",
			tx:   session.LogNonDSTest,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tt.tx()
			require.NoErrorf(t, err, "bad test setup; sending %T", tx)

			got, err := sut.ParseLogs(ctx, tx, sim)
			require.NoErrorf(t, err, "%T.ParseLogs()", sut)
			assert.Equalf(t, tt.want, got, "%T.ParseLogs() raw values", sut)
			assert.Equalf(t, tt.wantAsStr, got.String(), "%T.ParseLogs() as string", sut)
		})
	}
}
