package main

import (
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db *gorm.DB
	connectionError error
)

type Team struct {
	gorm.Model
	Name string
	CurrentValue int
}

func GetTeam(c *gin.Context) {
	teamID := c.Param("team_id")
	var team Team
	db.First(&team, teamID)
	
	c.JSON(200, gin.H{
		"id": teamID,
		"name": team.Name,
		"value": team.CurrentValue,
	})
}

func main() {
	router := gin.Default()
	dbname := "lcs_sm.db"

	// Establish database connection
	db, connectionError = gorm.Open("sqlite3", dbname)
	if connectionError != nil {
		panic(connectionError)
	}
	defer db.Close()

	// TEMP: Seed the database manually
	//db.AutoMigrate(&Team{})

	//db.Create(&Team{Name: "Cloud9", CurrentValue: 605})
	//db.Create(&Team{Name: "Team Liquid", CurrentValue: 580})

	// Routing
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API Version 1
	v1 := router.Group("/v1")
	{
		v1.GET("/teams/:team_id", GetTeam)
	}

	router.Run(":8080")
}