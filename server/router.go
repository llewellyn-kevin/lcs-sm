package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func SetRoutes() {
	router := gin.Default()
	router.Use(Authorize())

	var storeURL string
	if storeURL = os.Getenv("REDIS_URL"); storeURL == "" {
		// by default, dial localhost:6379, no authentication
		storeURL = "redis://:@:6379"
	}

	store, err := redis.DialURL(storeURL)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	authController := InitAuthController(store)
	leagueController := InitLeagueController(store)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	anonymous := router.Group("/api/v2")
	{
		anonymous.POST("/signup", authController.Signup)
		anonymous.POST("/signin", authController.Signin)

		anonymous.GET("/leagues", leagueController.Index)
		anonymous.GET("/leagues/:id", leagueController.Show)
	}

	authorized := router.Group("/api/v2")
	authorized.Use(CheckAuthenticated())
	{
		authorized.POST("/signout", authController.Signout)
	}

	admin := router.Group("/api/v2")
	admin.Use(CheckAdmin())
	{
		admin.POST("/leagues", leagueController.Create)
		admin.PUT("/leagues/:id", leagueController.Update)
		admin.DELETE("/leagues/:id", leagueController.Destroy)
	}

	/*
	  su := router.Group("/api/v2/su")
	  {
	  }
	*/

	router.Use(cors.Default())
	router.Run(":8080")
}
