// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	avalanchemetrics "github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/leveldb"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/database/meterdb"
	"github.com/ava-labs/avalanchego/database/pebbledb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	avalancheNode "github.com/ava-labs/avalanchego/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// createDatabase returns a new database instance with the provided configuration
func (vm *VM) createDatabase(dbConfig avalancheNode.DatabaseConfig) (database.Database, error) {
	dbRegisterer, err := avalanchemetrics.MakeAndRegister(
		vm.ctx.Metrics,
		dbMetricsPrefix,
	)
	if err != nil {
		return nil, err
	}
	var db database.Database
	// start the db
	switch dbConfig.Name {
	case leveldb.Name:
		dbPath := filepath.Join(dbConfig.Path, leveldb.Name)
		db, err = leveldb.New(dbPath, dbConfig.Config, vm.ctx.Log, dbRegisterer)
		if err != nil {
			return nil, fmt.Errorf("couldn't create %s at %s: %w", leveldb.Name, dbPath, err)
		}
	case memdb.Name:
		db = memdb.New()
	case pebbledb.Name:
		dbPath := filepath.Join(dbConfig.Path, pebbledb.Name)
		db, err = pebbledb.New(dbPath, dbConfig.Config, vm.ctx.Log, dbRegisterer)
		if err != nil {
			return nil, fmt.Errorf("couldn't create %s at %s: %w", pebbledb.Name, dbPath, err)
		}
	default:
		return nil, fmt.Errorf(
			"db-type was %q but should have been one of {%s, %s, %s}",
			dbConfig.Name,
			leveldb.Name,
			memdb.Name,
			pebbledb.Name,
		)
	}

	if dbConfig.ReadOnly && dbConfig.Name != memdb.Name {
		db = versiondb.New(db)
	}

	meterDBReg, err := avalanchemetrics.MakeAndRegister(
		vm.ctx.Metrics,
		"meterdb",
	)
	if err != nil {
		return nil, err
	}

	db, err = meterdb.New(meterDBReg, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create meterdb: %w", err)
	}

	return db, nil
}

// useStandaloneDatabase returns true if the chain can and should use a standalone database
// other than given by [db] in Initialize()
func (vm *VM) useStandaloneDatabase(acceptedDB database.Database) (bool, error) {
	// no config provided, use default
	standaloneDBFlag := vm.config.UseStandaloneDatabase
	if standaloneDBFlag != nil {
		return standaloneDBFlag.Bool(), nil
	}

	// check if the chain can use a standalone database
	_, err := acceptedDB.Get(lastAcceptedKey)
	if err == database.ErrNotFound {
		// If there is nothing in the database, we can use the standalone database
		return true, nil
	}
	return false, err
}

// getDatabaseConfig returns the database configuration for the chain
// to be used by separate, standalone database.
func getDatabaseConfig(config Config, chainDataDir string) (avalancheNode.DatabaseConfig, error) {
	var (
		configBytes []byte
		err         error
	)
	if len(config.DatabaseConfigContent) != 0 {
		dbConfigContent := config.DatabaseConfigContent
		configBytes, err = base64.StdEncoding.DecodeString(dbConfigContent)
		if err != nil {
			return avalancheNode.DatabaseConfig{}, fmt.Errorf("unable to decode base64 content: %w", err)
		}
	} else if len(config.DatabaseConfigFile) != 0 {
		configPath := config.DatabaseConfigFile
		configBytes, err = os.ReadFile(configPath)
		if err != nil {
			return avalancheNode.DatabaseConfig{}, err
		}
	}

	dbPath := filepath.Join(chainDataDir, "db")
	if len(config.DatabasePath) != 0 {
		dbPath = config.DatabasePath
	}

	return avalancheNode.DatabaseConfig{
		Name:     config.DatabaseType,
		ReadOnly: config.DatabaseReadOnly,
		Path:     dbPath,
		Config:   configBytes,
	}, nil
}

func inspectDB(db database.Database, label string) error {
	it := db.NewIterator()
	defer it.Release()

	var (
		count  int64
		start  = time.Now()
		logged = time.Now()

		// Totals
		total common.StorageSize
	)
	// Inspect key-value database first.
	for it.Next() {
		var (
			key  = it.Key()
			size = common.StorageSize(len(key) + len(it.Value()))
		)
		total += size
		count++
		if count%1000 == 0 && time.Since(logged) > 8*time.Second {
			log.Info("Inspecting database", "label", label, "count", count, "elapsed", common.PrettyDuration(time.Since(start)))
			logged = time.Now()
		}
	}
	// Display the database statistic.
	log.Info("Database statistics", "label", label, "total", total.String(), "count", count)
	return nil
}
