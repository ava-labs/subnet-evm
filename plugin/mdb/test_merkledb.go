// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"errors"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/ethdb"
)

func copyMemDB(db ethdb.Database) (ethdb.Database, error) {
	newDB := rawdb.NewMemoryDatabase()
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		if err := newDB.Put(iter.Key(), iter.Value()); err != nil {
			return nil, err
		}
	}

	return newDB, nil
}

func copyMerkleDB(db merkledb.MerkleDB) (merkledb.MerkleDB, error) {
	memDB := memdb.New()
	newDB, err := merkledb.New(context.Background(), memDB, NewBasicConfig())
	if err != nil {
		return nil, err
	}
	iter := db.NewIterator()
	defer iter.Release()
	for iter.Next() {
		if err := newDB.Put(iter.Key(), iter.Value()); err != nil {
			return nil, err
		}
	}
	return newDB, nil
}

func CopyMemDB(db ethdb.Database) (ethdb.Database, error) {
	mdb, ok := db.(*WithMerkleDB)
	if !ok {
		return nil, errors.New("not merkleDB")
	}
	dbCopy, err := copyMemDB(mdb.Database)
	if err != nil {
		return nil, err
	}
	mdbCopy, err := copyMerkleDB(mdb.merkleDB)
	if err != nil {
		return nil, err
	}
	return NewWithMerkleDB(dbCopy, mdbCopy), nil
}