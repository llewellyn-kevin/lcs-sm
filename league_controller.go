package main

import(
  "fmt"
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
  Name string `json:"name" binding:"required"`
}

func (l *LeagueController) Index(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
}

func (l *LeagueController) Show(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
}

// Create takes a JSON Post request with the name for a new league, 
// adds the leage to the data store, and returns an OK response.
//
// A InternalServerError response indicates a problem communicating
// with the store. A BadRequest response indicates a problem with 
// the request.
func (l *LeagueController) Create(c *gin.Context) {
  if(HasPermission(c)) {
    // Unmarshall request
    var league League
    if err := c.ShouldBindJSON(&league); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    // Get next league id and generate a key for new league
    id, err := redis.Int(l.store.Do("INCR", "next-league-id"))
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    leaguekey := fmt.Sprintf("league:%d", id)

    // Create league hash
    if _, err := l.store.Do("HMSET", leaguekey, "name", league.Name); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    // Add leagueid to the list of leagues
    if _, err := l.store.Do("RPUSH", "leagues", id); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    // Send success response
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