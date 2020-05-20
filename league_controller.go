package main

import(
  "fmt"
  "net/http"
  "strconv"

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
  Id    int     `redii:"pk" json:"id"`
  Name  string  `redii:"name" json:"name" binding:"required"`
}

// Index attempts to return a slice of all leagues in the store.
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

// Show takes an id as a url parameter and attempts to return the 
// associated resource.
func (l *LeagueController) Show(c *gin.Context) {
  // instantiate league object and get id
  var league League
  var err error
  league.Id, err = strconv.Atoi(c.Param("id"))
  if err != nil {
    errText := c.Param("id")
    c.JSON(http.StatusBadRequest, gin.H{
      "error": fmt.Sprintf("could not parse '%s' as id", errText),
    })
    return
  }

  // fill struct values from store
  if err := HashToStruct(l.store, "league", &league); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": fmt.Sprintf("no league with id '%d' found", league.Id),
    })

    return
  }

  // return struct
  c.JSON(http.StatusOK, gin.H{"league":league})
}

// Create accepts a JSON Header POST request with the name for a 
// new league and attempts to create a resource in the data store.
func (l *LeagueController) Create(c *gin.Context) {
  if(HasPermission(c)) {
    var league League
    if err := c.ShouldBindJSON(&league); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    if err := CreateHash(l.store, "league", &league); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

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
