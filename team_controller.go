package main

import(
	"fmt"
	"strconv"
	"net/http"

	"github.com/llewellyn-kevin/lcs-sm/models"

	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	var teams []models.Team
	db.Find(&teams)
	c.JSON(http.StatusOK, &teams)
}

func GetTeam(c *gin.Context) {
	teamID := c.Param("team-id")
	var team models.Team
	var stockValues []models.StockValue

	db.First(&team, teamID)
	db.Where("team_id = ?", teamID).Find(&stockValues)
	team.StockValues = stockValues

	if team.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No team with ID: %v was found", teamID))
	} else {
		c.JSON(http.StatusOK, &team)
	}
}

func SplitsByTeam(c *gin.Context) {
  teamID := c.Param("team-id")
  //var team models.Team
  //var splits []models.Split

  rows, err := db.Table("stock_values").Where("team_id = ?", teamID).Select("split_id").Group("split_id").Rows()
  if err != nil {
    c.JSON(500, "Internal Server Error")
  } else {
    c.JSON(http.StatusOK, rows)
  }
}

func CreateTeam(c *gin.Context) {
	name := c.Query("name")
	team := models.Team{Name: name, CurrentValue: 500}

	db.Create(&models.Team{Name: name, CurrentValue: 500})

	if db.NewRecord(team) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Team succesfully created",
			"team": &team,
		})
	} else {
		c.String(http.StatusInternalServerError, "Unable to create resource")
	}
}

func DeleteTeam(c *gin.Context) {
	teamID := c.Param("team-id")
	var team models.Team
	db.First(&team, teamID)

	if team.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No team with ID: %v was found", teamID))
	} else if team.DeletedAt != nil {
		c.String(http.StatusBadRequest, "Resource has already been deleted")
	} else {
		db.Delete(&team)
		c.JSON(http.StatusOK, gin.H{
			"message": "Team succesfully deleted",
			"team": &team,
		})
	}
}

func UpdateTeam(c *gin.Context) {
	teamID := c.Param("team-id")
	value, err := strconv.Atoi(c.Query("new-value"))
	var team models.Team
	db.First(&team, teamID)

	if team.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No team with ID: %v was found", teamID))
	} else if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not convert %v to integer value", c.Query("new-value")))
	} else {
		team.CurrentValue = value
		db.Save(team)

		c.JSON(http.StatusOK, gin.H{
			"message": "Team value updated",
			"team": &team,
		})
	}
}
