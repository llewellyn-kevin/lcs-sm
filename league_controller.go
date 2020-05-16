package main

import(
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
  if(HasPermission(c)) {

  }
}

func (l *LeagueController) Show(c *gin.Context) {
  if(HasPermission(c)) {

  }
}

func (l *LeagueController) Create(c *gin.Context) {
  if(HasPermission(c)) {

  }
}

func (l *LeagueController) Update(c *gin.Context) {
  if(HasPermission(c)) {

  }
}

func (l *LeagueController) Destroy(c *gin.Context) {
  if(HasPermission(c)) {

  }
}
