// Package evmsim provides convenience extensions to geth's SimulatedBackend.
package evmsim

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind/backends"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// A Backend wraps and extends a [backends.SimulatedBackend].
type Backend struct {
	*backends.SimulatedBackend

	// AutoCommit configures whether or not to automatically call
	// [backends.SimulatedBackend.Commit] after every transaction is sent. If
	// `true`, there is only one transaction per block, but countless hours will
	// be saved as the developer tears their hair out in confusion. The default
	// value is therefore `true`.
	AutoCommit bool

	txOpts []*bind.TransactOpts
}

// SendTransaction propagates its arguments to the underlying
// [backends.SimulatedBackend.SendTransaction] and returns any error that
// occurs. If no error was returned and `b.AutoCommit` is true then it calls
// `b.Commit(true)`.
func (b *Backend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if err := b.SimulatedBackend.SendTransaction(ctx, tx); err != nil {
		return err
	}
	if b.AutoCommit {
		b.Commit(true)
	}
	return nil
}

// NewWithHexKeys is equivalent to [NewWithRawKeys], but accepting hex-encoded
// private keys.
func NewWithHexKeys(tb testing.TB, keys []string) *Backend {
	tb.Helper()

	bKeys := make([][]byte, len(keys))
	for i, k := range keys {
		b, err := hexutil.Decode(k)
		if err != nil {
			tb.Fatalf("decode hex private key: %v", err)
		}
		bKeys[i] = b
	}
	return NewWithRawKeys(tb, bKeys)
}

// NewWithNumKeys generates `n` private keys determinastically, passing them to
// [NewWithRawKeys] and returning the result.
func NewWithNumKeys(tb testing.TB, n uint) *Backend {
	tb.Helper()
	keys := make([][]byte, n)
	for i := range keys {
		keys[i] = crypto.Keccak256(keys...)
	}
	return NewWithRawKeys(tb, keys)
}

// NewWithRawKeys constructs and returns a new [Backend] with pre-allocated
// accounts corresponding to the raw private keys, which are converted using
// [crypto.ToECDSA].
func NewWithRawKeys(tb testing.TB, keys [][]byte) *Backend {
	tb.Helper()

	txOpts := make([]*bind.TransactOpts, len(keys))
	for i, k := range keys {
		priv, err := crypto.ToECDSA(k)
		if err != nil {
			tb.Fatalf("convert raw %T private key: %v", k, err)
		}

		opt, err := bind.NewKeyedTransactorWithChainID(priv, big.NewInt(1337))
		if err != nil {
			tb.Fatalf("create test-chain %T from %T: %v", opt, priv, err)
		}
		txOpts[i] = opt
	}
	return newWithTransactOpts(tb, txOpts)
}

func newWithTransactOpts(tb testing.TB, txOpts []*bind.TransactOpts) *Backend {
	tb.Helper()

	alloc := make(core.GenesisAlloc)
	for _, o := range txOpts {
		alloc[o.From] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(10), big.NewInt(18+3), nil), // 1k ether
		}
	}

	return &Backend{
		SimulatedBackend: backends.NewSimulatedBackend(alloc, 30e6),
		AutoCommit:       true,
		txOpts:           txOpts,
	}
}

// Addr returns the address of the i'th account; see [Backend.From] re
// account ordering.
func (b *Backend) Addr(i int) common.Address {
	return b.txOpts[i].From
}

// From returns the [bind.TransactOpts] of the i'th account, the order
// corresponding to the private keys provided to construct `b`.
func (b *Backend) From(i int) *bind.TransactOpts {
	cp := *b.txOpts[i]
	return &cp
}

// A Deployer deploys a Solidity contract with no constructor arguments.
type Deployer[T any] func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, T, error)

// Deploy deploys a `T` contract as the i'th account; see [Backend.From] re
// account ordering.
func Deploy[T any](tb testing.TB, b *Backend, from int, fn Deployer[T]) (common.Address, T) {
	tb.Helper()
	addr, _, contract, err := fn(b.From(from), b)
	if err != nil {
		tb.Fatalf("deploying %T: %v", contract, err)
	}
	return addr, contract
}

// A Binder binds to a Solidity contract deployed at the given address.
type Binder[T any] func(common.Address, bind.ContractBackend) (T, error)

// Bind binds to a `T` contract at the given address.
func Bind[T any](tb testing.TB, b *Backend, fn Binder[T], addr common.Address) T {
	tb.Helper()
	contract, err := fn(addr, b)
	if err != nil {
		tb.Fatalf("binding to %T at %v: %v", contract, addr, err)
	}
	return contract
}

// Call calls the bound contract method, asserting that no error occurred, and
// returns the value returned by the contract.
func Call[T any](tb testing.TB, fn func(*bind.CallOpts) (T, error), opts *bind.CallOpts) T {
	tb.Helper()
	x, err := fn(opts)
	if err != nil {
		tb.Fatalf("contract call with %T(%+v): %v", opts, opts, err)
	}
	return x
}
