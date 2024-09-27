package database

import (
	"fmt"
	"path/filepath"

	"github.com/ava-labs/avalanchego/api/metrics"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/leveldb"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/database/meterdb"
	"github.com/ava-labs/avalanchego/database/pebbledb"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/node"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/metric"
	"github.com/ethereum/go-ethereum/log"
)

var dbNamespace = "subnet_evm" + metric.NamespaceSeparator + "db"

// New returns a new database instance with the provided configuration
func New(gatherer metrics.MultiGatherer, dbConfig node.DatabaseConfig, logger logging.Logger) (database.Database, error) {
	log.Info("initializing database")
	dbRegisterer, err := metrics.MakeAndRegister(
		gatherer,
		dbNamespace,
	)
	if err != nil {
		return nil, err
	}

	var db database.Database
	// start the db
	switch dbConfig.Name {
	case leveldb.Name:
		dbPath := filepath.Join(dbConfig.Path, leveldb.Name)
		db, err = leveldb.New(dbPath, dbConfig.Config, logger, dbRegisterer)
		if err != nil {
			return nil, fmt.Errorf("couldn't create %s at %s: %w", leveldb.Name, dbPath, err)
		}
	case memdb.Name:
		db = memdb.New()
	case pebbledb.Name:
		dbPath := filepath.Join(dbConfig.Path, pebbledb.Name)
		db, err = pebbledb.New(dbPath, dbConfig.Config, logger, dbRegisterer)
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

	meterDBReg, err := metrics.MakeAndRegister(
		gatherer,
		"all",
	)
	if err != nil {
		return nil, err
	}

	db, err = meterdb.New(meterDBReg, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
