package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db 				*gorm.DB
	connectionError error
)

type Team struct {
	gorm.Model
	Name 		 string
	CurrentValue int
}

type Split struct {
	gorm.Model
	League string
	Season string
	Year   int
}

type Value struct {
	gorm.Model
	TeamID  int
	SplitID int
	Week 	int
}

func GetTeams(c *gin.Context) {
	var teams []Team
	db.Find(&teams)
	c.JSON(http.StatusOK, &teams)
}

func GetTeam(c *gin.Context) {
	teamID := c.Param("team_id")
	var team Team
	db.First(&team, teamID)
	
	if team.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No team with ID: %v was found", teamID))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id": team.ID,
			"name": team.Name,
			"value": team.CurrentValue,
		})
	}
}

func CreateTeam(c *gin.Context) {
	name := c.Query("name")
	team := Team{Name: name, CurrentValue: 500}

	db.Create(&Team{Name: name, CurrentValue: 500})

	if db.NewRecord(team) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Team succesfully created",
			"team": gin.H{
				"id": team.ID,
				"name": team.Name,
				"value": team.CurrentValue,
			},
			"endpoints": gin.H{
				"GET": fmt.Sprintf("/v1/teams/%v", team.ID),
			},
		})
	} else {
		c.String(http.StatusInternalServerError, "Unable to create resource")
	}
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
		v1.GET("/teams", GetTeams)
		v1.GET("/teams/:team_id", GetTeam)
		v1.POST("/teams", CreateTeam)
	}

	router.Run(":8080")
}