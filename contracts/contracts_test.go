package contracts

import (
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/testing/evmsim"
)

//go:generate sh -c "solc --evm-version=paris --base-path=./ --include-path=./node_modules --combined-json=abi,bin contracts/**/*.sol | abigen --combined-json=- --pkg contracts | sed -E 's,github.com/ethereum/go-ethereum/(accounts|core)/,github.com/ava-labs/subnet-evm/\\1/,' > generated_test.go"

func newEVMSim(tb testing.TB, genesis params.Precompiles) *evmsim.Backend {
	tb.Helper()

	// The geth SimulatedBackend constructor doesn't allow for injection of
	// the ChainConfig, instead using a global. They have recently overhauled
	// the implementation so there's no point in sending a PR to allow for
	// injection.
	// TODO(arr4n): once we have upgraded to a geth version with the new
	// simulated.Backend, change how we inject the precompiles.
	copy := *params.TestChainConfig
	defer func() {
		params.TestChainConfig = &copy
	}()
	params.TestChainConfig.GenesisPrecompiles = genesis

	return evmsim.NewWithHexKeys(tb, keys())
}

// keys returns the hex-encoded private keys of the testing accounts; these
// identically match the accounts used in the Hardhat config.
func keys() []string {
	return []string{
		"0x56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027",
		"0x7b4198529994b0dc604278c99d153cfd069d594753d471171a1d102a10438e07",
		"0x15614556be13730e9e8d6eacc1603143e7b96987429df8726384c2ec4502ef6e",
		"0x31b571bf6894a248831ff937bb49f7754509fe93bbd2517c9c73c4144c0e97dc",
		"0x6934bef917e01692b789da754a0eae31a8536eb465e7bff752ea291dad88c675",
		"0xe700bdbdbc279b808b1ec45f8c2370e4616d3a02c336e68d85d4668e08f53cff",
		"0xbbc2865b76ba28016bc2255c7504d000e046ae01934b04c694592a6276988630",
		"0xcdbfd34f687ced8c6968854f8a99ae47712c4f4183b78dcc4a903d1bfe8cbf60",
		"0x86f78c5416151fe3546dece84fda4b4b1e36089f2dbc48496faf3a950f16157c",
		"0x750839e9dbbd2a0910efe40f50b2f3b2f2f59f5580bb4b83bd8c1201cf9a010a",
	}
}

// Convenience labels for using an account by name instead of number.
const (
	admin = iota
	_
	_
	_
	_
	_
	_
	_
	_
	allowlistOther
)
