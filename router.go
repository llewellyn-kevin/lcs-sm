package main

import(
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetRoutes() {
	router := gin.Default()

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API Version 1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/teams", GetTeams)
		v1.GET("/teams/:team-id", GetTeam)
		v1.POST("/teams", CreateTeam)
		v1.DELETE("/teams/:team-id", DeleteTeam)
		v1.PUT("/teams/:team-id", UpdateTeam)

		v1.GET("/splits", GetSplits)
		v1.GET("/splits/:split-id", GetSplit)
		v1.POST("/splits", CreateSplit)
		v1.DELETE("/splits/:split-id", DeleteSplit)

		v1.GET("/stock-values", GetStockValues)
		v1.GET("/stock-values/:stock-value-id", GetStockValue)
		v1.DELETE("/stock-values/:stock-value-id", DeleteStockValue)
		v1.PUT("/stock-values/:stock-value-id", UpdateStockValue)

		v1.POST("/splits/:split-id/teams/:team-id/stock-values", CreateStockValue)
	}

	router.Use(cors.Default())
	router.Run(":8080")
}