package contracts

import (
	"context"
	"testing"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/testing/dstest"
	"github.com/ava-labs/subnet-evm/testing/evmsim"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	_ "github.com/ava-labs/subnet-evm/precompile/registry"
)

func TestAllowList(t *testing.T) {
	// This is a demonstration of a Go implementation of the Hardhat tests:
	// https://github.com/ava-labs/subnet-evm/blob/dc1d78da/contracts/test/tx_allow_list.ts

	ctx := context.Background()

	sim := newEVMSim(t, params.Precompiles{
		txallowlist.ConfigKey: txallowlist.NewConfig(
			new(uint64),
			[]common.Address{common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")},
			nil, nil,
		),
	})

	allow := evmsim.Bind(t, sim, NewIAllowList, txallowlist.ContractAddress)
	allowSess := &IAllowListSession{
		Contract:     allow,
		TransactOpts: *sim.From(admin),
	}

	sutAddr, sut := evmsim.Deploy(t, sim, admin, DeployExampleTxAllowListTest)
	sutSess := &ExampleTxAllowListTestSession{
		Contract:     sut,
		TransactOpts: *sim.From(admin),
	}
	parser := dstest.New(sutAddr)

	_, err := allowSess.SetAdmin(sutAddr)
	require.NoErrorf(t, err, "%T.SetAdmin(%T address)", allow, sut)

	_, err = sutSess.SetUp()
	require.NoErrorf(t, err, "%T.SetUp()", sut)

	// TODO: This table of steps is purely to demonstrate a *direct*
	// reimplementation of the Hardhat tests in Go. I (arr4n) believe we should
	// refactor the tests before a complete translation, primarily to reduce the
	// number of calls that have to happen here. Also note that the original
	// tests use a `beforeEach()` whereas the above preamble is equivalent to a
	// `beforeAll()`.
	for _, step := range []struct {
		name string
		fn   (func() (*types.Transaction, error))
	}{
		{"should add contract deployer as admin", sutSess.StepContractOwnerIsAdmin},
		{"precompile should see admin address has admin role", sutSess.StepPrecompileHasDeployerAsAdmin},
		{"precompile should see test address has no role", sutSess.StepNewAddressHasNoRole},
	} {
		t.Run(step.name, func(t *testing.T) {
			tx, err := step.fn()
			require.NoError(t, err, "running step")

			// TODO(arr4n): DSTest can be used for general logging, so the
			// following pattern is only valid when we don't use it as such.
			// Failing assertions, however, set a private `failed` boolean,
			// which we need to expose. Alternatively, they also hook into HEVM
			// cheatcodes if they're implemented; having these would be useful
			// for testing in general.
			failures := parser.ParseTB(ctx, t, tx, sim)
			if len(failures) > 0 {
				t.Errorf("Assertion failed:\n%s", failures)
			}
		})
	}
}
