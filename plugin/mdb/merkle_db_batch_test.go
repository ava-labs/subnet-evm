// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/stretchr/testify/require"
)

func TestBatchOrderCommit(t *testing.T) {
	require := require.New(t)
	memDB := memdb.New()
	db, err := merkledb.New(context.Background(), memDB, NewBasicConfig())
	require.NoError(err)

	b1 := db.NewBatch()
	b2 := db.NewBatch()
	for i := 0; i < 10; i++ {
		b1.Put([]byte{byte(i)}, []byte{byte(i)})
		b2.Put([]byte{byte(i + 500)}, []byte{byte(i)})
	}
	require.NoError(b2.Write())
	require.NoError(b1.Write())
}
