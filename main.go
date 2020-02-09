package main

import (
	"github.com/gin-gonic/gin"
)

func Getteam(c *gin.Context) {
	TeamID := c.Param("team_id")

	c.JSON(200, gin.H{
		"name": "Cloud9",
		"team_id": TeamID,
		"current_value": "630",
		"week1_value": "575",
	})
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API Version 1
	v1 := router.Group("/v1")
	{
		v1.GET("/teams/:team_id", Getteam)
	}

	router.Run(":8080")
}