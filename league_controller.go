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
  var leagues []League
  ids, err := GetList(l.store, "league")
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  for _, v := range ids {
    var league League
    league.Id = v
    if err := GetHash(l.store, "league", &league); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    leagues = append(leagues, league)
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
  if err := GetHash(l.store, "league", &league); err != nil {
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

// Update accepts a JSON Header POST request with a name and id
// for a given league and attempts to update the resource in the
// data store.
func (l *LeagueController) Update(c *gin.Context) {
  if(HasPermission(c)) {
    var league League
    if err := c.ShouldBindJSON(&league); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    var err error
    league.Id, err = strconv.Atoi(c.Param("id"))
    if err != nil {
      errText := c.Param("id")
      c.JSON(http.StatusBadRequest, gin.H{
        "error": fmt.Sprintf("could not parse '%s' as id", errText),
      })
      return
    }

    if err := UpdateHash(l.store, "league", &league); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, gin.H{
      "success":"league updated",
      "league": league,
    })
  }
}

func (l *LeagueController) Destroy(c *gin.Context) {
  if(HasPermission(c)) {
    c.JSON(http.StatusOK, gin.H{"success":"not implemented"})
  }
}
