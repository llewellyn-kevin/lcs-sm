package main

import(
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

	// TEMP: Seed the database manually
	//db.AutoMigrate(&Team{})

	//db.Create(&Team{Name: "Cloud9", CurrentValue: 605})
	//db.Create(&Team{Name: "Team Liquid", CurrentValue: 580})
}