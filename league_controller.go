package main

import(
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

type LeagueController struct {
  store redis.Conn
}

func InitLeagueController(store redis.Conn) *LeagueController {
  return &LeagueController{store}
}

type League struct {
  Name string
}

func (l *LeagueController) Index(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
}

func (l *LeagueController) Show(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
}

func (l *LeagueController) Create(c *gin.Context) {
  if(HasPermission(c)) {
    c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
  }
}

func (l *LeagueController) Update(c *gin.Context) {
  if(HasPermission(c)) {
    c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
  }
}

func (l *LeagueController) Destroy(c *gin.Context) {
  if(HasPermission(c)) {
    c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
  }
}
