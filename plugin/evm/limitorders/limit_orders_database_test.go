package limitorders

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeDatabaseFirstTime(t *testing.T) {
	lod, err := InitializeDatabase()
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	assert.NotNil(t, lod)
	assert.Nil(t, err)

	_, err = os.Stat(dbName)
	assert.Nil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	rows, err := db.Query("SELECT * FROM limit_orders")
	assert.Nil(t, err)
	assert.False(t, rows.Next())
}

func TestInitializeDatabaseAfterInitializationAlreadyDone(t *testing.T) {
	InitializeDatabase()
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	dbFileInfo1, _ := os.Stat(dbName)

	_, err := InitializeDatabase()
	assert.Nil(t, err)

	dbFileInfo2, err := os.Stat(dbName)
	assert.Nil(t, err)
	assert.Equal(t, dbFileInfo1.Size(), dbFileInfo2.Size())
	assert.Equal(t, dbFileInfo1.ModTime(), dbFileInfo2.ModTime())
}

func TestInsertLimitOrderFailureWhenPositionTypeIsWrong(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := ""
	baseAssetQuantity := 10
	price := 10.14
	salt := "123"
	signature := []byte("signature")
	positionType := "neutral"
	err := lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	assert.NotNil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, base_asset_quantity, price from limit_orders where user_address = ?")
	rows, _ := stmt.Query(userAddress)
	assert.False(t, rows.Next())
}
func TestInsertLimitOrderFailureWhenUserAddressIsBlank(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := ""
	baseAssetQuantity := 10
	price := 10.14
	positionType := "long"
	salt := "123"
	signature := []byte("signature")
	err := lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	assert.NotNil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, base_asset_quantity, price from limit_orders where user_address = ?")
	rows, _ := stmt.Query(userAddress)
	assert.False(t, rows.Next())
}

func TestInsertLimitOrderFailureWhenBaseAssetQuantityIsZero(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 0
	price := 10.14
	positionType := "long"
	salt := "123"
	signature := []byte("signature")
	err := lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	assert.NotNil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, base_asset_quantity, price from limit_orders where user_address = ?")
	rows, _ := stmt.Query(userAddress)
	assert.False(t, rows.Next())
}

func TestInsertLimitOrderFailureWhenPriceIsZero(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 10
	price := 0.0
	positionType := "long"
	salt := "123"
	signature := []byte("signature")
	err := lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	assert.NotNil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, base_asset_quantity, price from limit_orders where user_address = ?")
	rows, _ := stmt.Query(userAddress)
	assert.False(t, rows.Next())
}

func TestInsertLimitOrderSuccess(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 10
	price := 10.14
	positionType := "long"
	salt := "123"
	signature := []byte("signature")
	err := lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	assert.Nil(t, err)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, position_type, base_asset_quantity, price, status from limit_orders where user_address = ?")
	rows, _ := stmt.Query(userAddress)
	defer rows.Close()
	for rows.Next() {
		var queriedId int
		var queriedPositionType string
		var queriedBaseAssetQuantity int
		var queriedPrice float64
		var queriedStatus string
		_ = rows.Scan(&queriedId, &queriedPositionType, &queriedBaseAssetQuantity, &queriedPrice, &queriedStatus)
		assert.Equal(t, 1, queriedId)
		assert.Equal(t, positionType, queriedPositionType)
		assert.Equal(t, baseAssetQuantity, queriedBaseAssetQuantity)
		assert.Equal(t, price, queriedPrice)
		assert.Equal(t, "open", queriedStatus)
	}
	positionType = "short"
	err = lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, "1", signature)
	assert.Nil(t, err)
	stmt, _ = db.Prepare("SELECT id, user_address, base_asset_quantity, price, status from limit_orders where position_type = ?")
	rows, _ = stmt.Query(userAddress)
	defer rows.Close()
	for rows.Next() {
		var queriedId int
		var queriedUserAddress string
		var queriedBaseAssetQuantity int
		var queriedPrice float64
		var queriedStatus string
		_ = rows.Scan(&queriedId, &queriedUserAddress, &queriedBaseAssetQuantity, &queriedPrice, &queriedStatus)
		assert.Equal(t, 1, queriedId)
		assert.Equal(t, userAddress, queriedUserAddress)
		assert.Equal(t, baseAssetQuantity, queriedBaseAssetQuantity)
		assert.Equal(t, price, queriedPrice)
		assert.Equal(t, "open", queriedStatus)
	}

}

func TestGetLimitOrderByPositionTypeAndPriceWhenShortOrders(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 10
	price1 := 10.14
	price2 := 11.14
	price3 := 12.14
	positionType := "short"
	signature := []byte("signature")
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price1, "1", signature)
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price2, "2", signature)
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price3, "3", signature)
	orders := lod.GetLimitOrderByPositionTypeAndPrice("short", 11.14)
	assert.Equal(t, 2, len(orders))
	for i := 0; i < len(orders); i++ {
		assert.Equal(t, orders[i].UserAddress, userAddress)
		assert.Equal(t, orders[i].BaseAssetQuantity, baseAssetQuantity)
		assert.Equal(t, orders[i].PositionType, positionType)
		assert.Equal(t, orders[i].Status, "open")
	}
	assert.Equal(t, price1, orders[0].Price)
	assert.Equal(t, price2, orders[1].Price)
}

func TestGetLimitOrderByPositionTypeAndPriceWhenLongOrders(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 10
	price1 := 10.14
	price2 := 11.14
	price3 := 12.14
	positionType := "long"
	signature := []byte("signature")
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price1, "1", signature)
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price2, "2", signature)
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price3, "3", signature)
	orders := lod.GetLimitOrderByPositionTypeAndPrice("long", 11.14)
	assert.Equal(t, 2, len(orders))
	for i := 0; i < len(orders); i++ {
		assert.Equal(t, orders[i].UserAddress, userAddress)
		assert.Equal(t, orders[i].BaseAssetQuantity, baseAssetQuantity)
		assert.Equal(t, orders[i].PositionType, positionType)
		assert.Equal(t, orders[i].Status, "open")
	}
	assert.Equal(t, price2, orders[0].Price)
	assert.Equal(t, price3, orders[1].Price)
}

func TestUpdateLimitOrderStatus(t *testing.T) {
	dbName := fmt.Sprintf("./hubble%d.db", os.Getpid()) // so that every node has a different database
	defer os.Remove(dbName)
	lod, _ := InitializeDatabase()
	userAddress := "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"
	baseAssetQuantity := 10
	price := 10.50
	positionType := "long"
	signature := []byte("signature")
	salt := "1"
	lod.InsertLimitOrder(positionType, userAddress, baseAssetQuantity, price, salt, signature)
	newStatus := "fulfilled"
	lod.UpdateLimitOrderStatus(userAddress, salt, newStatus)

	db, _ := sql.Open("sqlite3", dbName) // Open the created SQLite File
	stmt, _ := db.Prepare("SELECT id, position_type, base_asset_quantity, price, status from limit_orders where user_address = ? and salt = ?")
	rows, _ := stmt.Query(userAddress, salt)
	defer rows.Close()
	for rows.Next() {
		var queriedId int
		var queriedPositionType string
		var queriedBaseAssetQuantity int
		var queriedPrice float64
		var queriedStatus string
		_ = rows.Scan(&queriedId, &queriedPositionType, &queriedBaseAssetQuantity, &queriedPrice, &queriedStatus)
		assert.Equal(t, 1, queriedId)
		assert.Equal(t, positionType, queriedPositionType)
		assert.Equal(t, baseAssetQuantity, queriedBaseAssetQuantity)
		assert.Equal(t, price, queriedPrice)
		assert.Equal(t, "fulfilled", queriedStatus)
	}
}
