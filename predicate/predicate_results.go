// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package predicate

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/codec/linearcodec"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ethereum/go-ethereum/common"
)

const (
	Version        = uint16(0)
	MaxResultsSize = units.MiB
)

var Codec codec.Manager

func init() {
	Codec = codec.NewManager(MaxResultsSize)

	c := linearcodec.NewDefault()
	errs := wrappers.Errs{}
	errs.Add(
		c.RegisterType(Results{}),
		Codec.RegisterCodec(Version, c),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}

// TxResults is a map of results for each precompile address to the resulting byte array.
type TxResults map[common.Address][]byte

// Results encodes the precompile predicate results included in a block on a per transaction basis.
// Results is not thread-safe.
type Results struct {
	Results map[common.Hash]TxResults `serialize:"true"`
}

// NewResults returns an empty predicate results.
func NewResults() *Results {
	return &Results{
		Results: make(map[common.Hash]TxResults),
	}
}

func NewResultsFromMap(results map[common.Hash]TxResults) *Results {
	return &Results{
		Results: results,
	}
}

// ParseResults parses [b] into predicate results.
func ParseResults(b []byte) (*Results, error) {
	res := new(Results)
	parsedVersion, err := Codec.Unmarshal(b, res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicate results: %w", err)
	}
	if parsedVersion != Version {
		return nil, fmt.Errorf("invalid version (found %d, expected %d)", parsedVersion, Version)
	}
	return res, nil
}

// GetResults returns the byte array results for [txHash] from precompile [address] if available.
func (p *Results) GetResults(txHash common.Hash, address common.Address) []byte {
	txResults, ok := p.Results[txHash]
	if !ok {
		return nil
	}
	return txResults[address]
}

// SetTxResults sets the predicate results for the given [txHash]. Overrides results if present.
func (p *Results) SetTxResults(txHash common.Hash, txResults TxResults) {
	// If there are no tx results, don't store an entry in the map
	if len(txResults) == 0 {
		delete(p.Results, txHash)
		return
	}
	p.Results[txHash] = txResults
}

// DeleteTxResults deletes the predicate results for the given [txHash].
func (p *Results) DeleteTxResults(txHash common.Hash) {
	delete(p.Results, txHash)
}

// Bytes marshals the current state of predicate results
func (p *Results) Bytes() ([]byte, error) {
	return Codec.Marshal(Version, p)
}

func (p *Results) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("PredicateResults: (Size = %d)", len(p.Results)))
	for txHash, results := range p.Results {
		for address, result := range results {
			sb.WriteString(fmt.Sprintf("\n%s    %s: %x", txHash, address, result))
		}
	}

	return sb.String()
}
