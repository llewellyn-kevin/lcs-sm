package main

import(
	"fmt"
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
	teamID := c.Param("team_id")
	var team models.Team
	db.First(&team, teamID)
	
	if team.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No team with ID: %v was found", teamID))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id": team.ID,
			"name": team.Name,
			"value": team.CurrentValue,
		})
	}
}

func CreateTeam(c *gin.Context) {
	name := c.Query("name")
	team := models.Team{Name: name, CurrentValue: 500}

	db.Create(&models.Team{Name: name, CurrentValue: 500})

	if db.NewRecord(team) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Team succesfully created",
			"team": gin.H{
				"id": team.ID,
				"name": team.Name,
				"value": team.CurrentValue,
			},
			"endpoints": gin.H{
				"GET": fmt.Sprintf("/v1/teams/%v", team.ID),
			},
		})
	} else {
		c.String(http.StatusInternalServerError, "Unable to create resource")
	}
}