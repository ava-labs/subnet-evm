// Package dstest implements parsing of [DSTest] Solidity-testing errors.
//
// [DSTest]: https://github.com/dapphub/ds-test
package dstest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/testing/dstest/internal/dstestbindings"
	"github.com/ethereum/go-ethereum/common"
)

// New returns a new `Parser` with the provided addresses already
// `Register()`ed.
func New(tests ...common.Address) *Parser {
	p := &Parser{
		tests: make(map[common.Address]bool),
	}
	for _, tt := range tests {
		p.Register(tt)
	}
	return p
}

// A Parser inspects transaction logs of `Register()`ed test addresses, parsing
// those that correspond to [DSTest] error logs.
//
// [DSTest]: https://github.com/dapphub/ds-test
type Parser struct {
	tests map[common.Address]bool
}

// Register marks the provided `Address` as being a test that inherits from the
// [DSTest contract].
//
// [DSTest contract]: https://github.com/dapphub/ds-test/blob/master/src/test.sol
func (p *Parser) Register(test common.Address) {
	p.tests[test] = true
}

// A Log represents a Solidity event emitted by the `DSTest` contract. Although
// all assertion failures result in an event, not all logged events correspond
// to failures.
type Log struct {
	unpacked map[string]any
}

// String returns `l` as a human-readable string.
//
// The format is not guaranteed to be stable and the returned value SHOULD NOT
// be parsed.
func (l Log) String() string {
	switch u := l.unpacked; len(u) {
	case 1:
		return fmt.Sprintf("%v", u["arg0"])
	case 2:
		return fmt.Sprintf("%s = %v", u["key"], u["val"])
	case 3:
		return fmt.Sprintf("%s = %v (%v decimals)", u["key"], u["val"], u["decimals"])
	default:
		// The above cases are exhaustive at the time of writing; if the default
		// is reached then they need to be updated.
		return fmt.Sprintf("%+v", map[string]any(u))
	}
}

type Logs []*Log

// String() returns `ls` as a human-readable string.
func (ls Logs) String() string {
	s := make([]string, len(ls))
	for i, l := range ls {
		s[i] = l.String()
	}
	return strings.Join(s, "\n")
}

// Parse finds all [types.Log]s emitted by test contracts in the provided
// `Transaction`, filters them to keep only those corresponding to `DSTest`
// events, and returns the unpacked data.
func (p *Parser) Parse(ctx context.Context, tx *types.Transaction, b bind.DeployBackend) (Logs, error) {
	r, err := b.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, err
	}

	var logs []*Log
	for _, l := range r.Logs {
		if !p.tests[l.Address] {
			continue
		}
		ev, err := dstestbindings.EventByID(l.Topics[0])
		if err != nil /* not found */ {
			continue
		}

		l, err := unpack(ev, l)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func unpack(ev *abi.Event, l *types.Log) (*Log, error) {
	unpacked := make(map[string]any)
	if err := dstestbindings.UnpackLogIntoMap(unpacked, ev.Name, *l); err != nil {
		return nil, err
	}
	return &Log{unpacked}, nil
}

// ParseTB is identical to [Parse] except that it reports all errors on
// [testing.TB.Fatal].
func (p *Parser) ParseTB(ctx context.Context, tb testing.TB, tx *types.Transaction, b bind.DeployBackend) Logs {
	tb.Helper()
	l, err := p.Parse(ctx, tx, b)
	if err != nil {
		tb.Fatalf("%T.Parse(): %v", p, err)
	}
	return l
}
