package main

import(
	"github.com/llewellyn-kevin/lcs-sm/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db 				*gorm.DB
	connectionError error
)

func ConnectDB(name string) {
	// Establish database connection
	db, connectionError = gorm.Open("sqlite3", name)
	if connectionError != nil {
		panic(connectionError)
	}

	// Database Migrations
	db.AutoMigrate(&models.Team{})
	db.AutoMigrate(&models.Split{})
	db.AutoMigrate(&models.StockValue{})
}