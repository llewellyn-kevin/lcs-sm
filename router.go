package main

import(
	"github.com/gin-gonic/gin"
)

func SetRoutes() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API Version 1
	v1 := router.Group("/v1")
	{
		v1.GET("/teams", GetTeams)
		v1.GET("/teams/:team-id", GetTeam)
		v1.POST("/teams", CreateTeam)
		v1.DELETE("/teams/:team-id", DeleteTeam)
		v1.PUT("/teams/:team-id", UpdateTeam)
	}

	router.Run(":8080")
}