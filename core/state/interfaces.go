package state

import (
	"reflect"

	"github.com/ava-labs/subnet-evm/core/state/snapshot"
	"github.com/ethereum/go-ethereum/common"
)

type SnapshotTree interface {
	Snapshot(common.Hash) snapshot.Snapshot
	UpdateWithBlockHash(
		root common.Hash,
		parent common.Hash,
		blockHash common.Hash,
		parentHash common.Hash,
		destructs map[common.Hash]struct{},
		accounts map[common.Hash][]byte,
		storage map[common.Hash]map[common.Hash][]byte,
	) error
}

// https://vitaneri.com/posts/check-for-nil-interface-in-go
func checkNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
