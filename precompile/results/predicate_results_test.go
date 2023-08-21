// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package results

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPredicateResults(t *testing.T) {
	require := require.New(t)
	predicateResults := NewPredicateResults()
	emptyResultsBytes, err := predicateResults.Bytes()
	require.NoError(err)

	parsedPredicateResults, err := ParsePredicateResults(emptyResultsBytes)
	require.NoError(err)
	require.Equal(predicateResults, parsedPredicateResults)
}
