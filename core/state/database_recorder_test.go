// (c) 2020-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"fmt"
	"testing"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBasicReplay(t *testing.T) {
	require := require.New(t)
	require.NoError(nil)
	cacheConfig := &trie.Config{
		Cache: 256,
	}

	db1 := rawdb.NewMemoryDatabase()
	trieDB1 := trie.NewDatabaseWithConfig(db1, cacheConfig)
	stateDB1 := NewDatabaseWithNodeDB(db1, trieDB1)

	r := NewRecording()
	r.RegisterType((*Database)(nil))
	r.RegisterType((*Trie)(nil))

	rdb := &RecordingDatabase{Database: stateDB1, r: r}
	tr, err := rdb.OpenTrie(types.EmptyRootHash)
	require.NoError(err)
	err = tr.UpdateStorage(common.Address{0x01}, []byte{0x01}, []byte{0x01})
	require.NoError(err)
	tr, err = rdb.OpenTrie(types.EmptyRootHash)
	require.NoError(err)
	err = tr.UpdateStorage(common.Address{0x01}, []byte{0x01}, []byte{0x02})
	require.NoError(err)
	v, err := tr.GetStorage(common.Address{0x01}, []byte{0x01})
	fmt.Println("got: ", v)
	require.NoError(err)
	err = tr.UpdateAccount(common.Address{0x01}, &types.StateAccount{
		Nonce: 1,
	})
	require.NoError(err)

	root, ns := tr.Commit(false)
	fmt.Println("root: ", root, "ns: ", ns)

	// fmt.Println(id)
	// fmt.Println(r.ifaces)

	// op0 := op{
	// 	receiver: id,
	// 	method:   "OpenTrie",
	// }
	// op0.args = append(op0.args, types.EmptyRootHash)
	// op0.expected = append(op0.expected, &register{})
	// op0.expected = append(op0.expected, nil)

	// r.ops = append(r.ops, op0)
	// r.Replay()

	for _, o := range r.Ops {
		fmt.Println(o)
	}
}
