package main

import(
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gomodule/redigo/redis"
)

// LeagueController contains a reference to the data store containing
// league info, as well as method handlers for league related
// requests.
type LeagueController struct {
  store redis.Conn
}

// InitLeagueController instantiates a LeagueController object
// with a reference to the data store.
func InitLeagueController(store redis.Conn) *LeagueController {
  return &LeagueController{store}
}

// League contains a representation of the leagues from the data 
// store.
type League struct {
  Id    int     `json:"id"`
  Name  string  `json:"name" binding:"required"`
}

// Index returns a JSON response containing the name and id of every
// league.
// 
// Returns an InternalServerError response if there is an issue
// communicating with the store.
func (l *LeagueController) Index(c *gin.Context) {
  // Grab a list of all leagues
  var leagues []League
  ids, err := redis.Ints(l.store.Do("LRANGE", "leagues", 0, -1))
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // Instantiate league objects to be returned
  for _, v := range ids {
    leaguekey := fmt.Sprintf("league:%d", v)
    if name, err := redis.String(l.store.Do("HGET", leaguekey, "name")); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    } else {
      leagues = append(leagues, League{v, name})
    }
  }

  c.JSON(http.StatusOK, gin.H{"leagues":leagues})
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
    league.Id = id
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
    c.JSON(http.StatusOK, gin.H{
      "success":"new league created",
      "league": league,
    })
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
