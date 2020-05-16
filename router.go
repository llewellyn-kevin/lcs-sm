package main

import(
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

func SetRoutes() {
	router := gin.Default()
  router.Use(Authorize())

  store, err := redis.Dial("tcp", ":6379")
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

	// API Version 1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/teams", GetTeams)
		v1.GET("/teams/:team-id", GetTeam)
		v1.POST("/teams", CreateTeam)
		v1.DELETE("/teams/:team-id", DeleteTeam)
		v1.PUT("/teams/:team-id", UpdateTeam)

    v1.GET("/teams/:team-id/splits", SplitsByTeam)

		v1.GET("/splits", GetSplits)
		v1.GET("/splits/:split-id", GetSplit)
		v1.POST("/splits", CreateSplit)
		v1.DELETE("/splits/:split-id", DeleteSplit)

		v1.GET("/stock-values", GetStockValues)
		v1.GET("/stock-values/:stock-value-id", GetStockValue)
		v1.DELETE("/stock-values/:stock-value-id", DeleteStockValue)
		v1.PUT("/stock-values/:stock-value-id", UpdateStockValue)

    v1.GET("/splits/:split-id/teams/:team-id/stock-values", GetSplitTeamStocks)
		v1.POST("/splits/:split-id/teams/:team-id/stock-values", CreateStockValue)
	}

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

  admin := router.Group("/api/v2/admin")
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
