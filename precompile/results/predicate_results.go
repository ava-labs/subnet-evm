// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package results

import (
	"fmt"
	"sync"

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
		c.RegisterType(PredicateResults{}),
		Codec.RegisterCodec(Version, c),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}

// TODO: add testing and comments
type PredicateResults struct {
	lock sync.RWMutex

	results map[common.Hash]map[common.Address][]byte `serialize:"true"`
}

func NewPredicateResults() *PredicateResults {
	return &PredicateResults{
		results: make(map[common.Hash]map[common.Address][]byte),
	}
}

func ParsePredicateResults(b []byte) (*PredicateResults, error) {
	res := new(PredicateResults)
	parsedVersion, err := Codec.Unmarshal(b, res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicate results: %w", err)
	}
	if parsedVersion != Version {
		return nil, fmt.Errorf("invalid version (found %d, expected %d)", parsedVersion, Version)
	}
	return res, nil
}

func (p *PredicateResults) GetPredicateResults(txHash common.Hash, address common.Address) []byte {
	p.lock.RLock()
	defer p.lock.RUnlock()

	txResults, ok := p.results[txHash]
	if !ok {
		return nil
	}
	return txResults[address]
}

func (p *PredicateResults) SetTxPredicateResults(txHash common.Hash, txResults map[common.Address][]byte) {
	// If there are no results to add, no need to grab the lock.
	if len(txResults) == 0 {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.results[txHash] = txResults
}

func (p *PredicateResults) DeleteTxPredicateResults(txHash common.Hash) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.results, txHash)
}

func (p *PredicateResults) Bytes() ([]byte, error) {
	return Codec.Marshal(Version, p)
}
