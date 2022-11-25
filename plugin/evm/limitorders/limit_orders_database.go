package limitorders

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type LimitOrder struct {
	positionType      string
	id                int64
	userAddress       string
	baseAssetQuantity int
	price             float64
}

type LimitOrderDatabase interface {
	InsertLimitOrder(positionType string, userAddress string, baseAssetQuantity int, price float64) error
	GetLimitOrderByPositionTypeAndPrice(positionType string, price float64) []LimitOrder
}

type limitOrderDatabase struct {
	db *sql.DB
}

func InitializeDatabase() (LimitOrderDatabase, error) {
	if _, err := os.Stat("hubble.db"); err != nil {
		file, err := os.Create("hubble.db") // Create SQLite file
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	database, _ := sql.Open("sqlite3", "./hubble.db") // Open the created SQLite File
	err := createTable(database)                      // Create Database Tables

	lod := &limitOrderDatabase{
		db: database,
	}

	return lod, err
}

func (lod *limitOrderDatabase) InsertLimitOrder(positionType string, userAddress string, baseAssetQuantity int, price float64) error {
	err := validateInsertLimitOrderInputs(positionType, userAddress, baseAssetQuantity, price)
	if err != nil {
		fmt.Println(err)
		return err
	}
	insertStudentSQL := "INSERT INTO limit_orders(user_address, position_type, base_asset_quantity, price) VALUES (?, ?, ?, ?)"
	statement, err := lod.db.Prepare(insertStudentSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(userAddress, positionType, baseAssetQuantity, price)
	return err
}

func (lod *limitOrderDatabase) GetLimitOrderByPositionTypeAndPrice(positionType string, price float64) []LimitOrder {
	var rows = &sql.Rows{}
	var limitOrders = []LimitOrder{}
	if positionType == "short" {
		rows = getShortLimitOrderByPrice(lod.db, price)
	}
	if positionType == "long" {
		rows = getLongLimitOrderByPrice(lod.db, price)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var userAddress string
		var baseAssetQuantity int
		var price float64
		_ = rows.Scan(&id, &userAddress, &baseAssetQuantity, &price)
		limitOrder := &LimitOrder{
			id:                id,
			positionType:      positionType,
			userAddress:       userAddress,
			baseAssetQuantity: baseAssetQuantity,
			price:             price,
		}
		limitOrders = append(limitOrders, *limitOrder)
	}
	return limitOrders
}

func getShortLimitOrderByPrice(db *sql.DB, price float64) *sql.Rows {
	stmt, _ := db.Prepare("SELECT id, user_address, base_asset_quantity, price from limit_orders where position_type = ? and price < ?")
	rows, _ := stmt.Query("short", price)
	return rows
}

func getLongLimitOrderByPrice(db *sql.DB, price float64) *sql.Rows {
	stmt, _ := db.Prepare("SELECT id, user_address, base_asset_quantity, price from limit_orders where position_type = ? and price > ?")
	rows, _ := stmt.Query("long", price)
	return rows
}

func createTable(db *sql.DB) error {
	createLimitOrderTableSql := `CREATE TABLE if not exists limit_orders (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"position_type" VARCHAR(64) NOT NULL, 
    "user_address" VARCHAR(64) NOT NULL,
    "base_asset_quantity" INTEGER NOT NULL,
    "price" FLOAT NOT NULL
	);`

	statement, err := db.Prepare(createLimitOrderTableSql)
	if err != nil {
		return err
	}
	_, err = statement.Exec() // Execute SQL Statements
	return err
}

func validateInsertLimitOrderInputs(positionType string, userAddress string, baseAssetQuantity int, price float64) error {
	if positionType == "long" || positionType == "short" {
	} else {
		return errors.New("invalid position type")
	}

	if userAddress == "" {
		return errors.New("user address cannot be blank")
	}

	if baseAssetQuantity == 0 {
		return errors.New("baseAssetQuantity cannot be zero")
	}

	if price == 0 {
		return errors.New("price cannot be zero")
	}

	return nil
}
