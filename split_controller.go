package main

import(
	"fmt"
	"strconv"
	"net/http"

	"github.com/llewellyn-kevin/lcs-sm/models"

	"github.com/gin-gonic/gin"
)

func GetSplits(c *gin.Context) {
	var splits []models.Split
	db.Find(&splits)
	c.JSON(http.StatusOK, &splits)
}

func GetSplit(c *gin.Context) {
	splitID := c.Param("split-id")
	var split models.Split
	db.First(&split, splitID)
	
	if split.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No split with ID: %v was found", splitID))
	} else {
		c.JSON(http.StatusOK, &split)
	}
}

func CreateSplit(c *gin.Context) {
	league := c.Query("league")
	season := c.Query("season")
	year, err := strconv.Atoi(c.Query("year"))

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Unable to convert year: %v to an integer value", year))
	} else {
		split := models.Split{League: league, Season: season, Year: year}

		db.Create(&models.Split{League: league, Season: season, Year: year})

		if db.NewRecord(split) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Split succesfully created",
				"split": &split,
			})
		} else {
			c.String(http.StatusInternalServerError, "Unable to create resource")
		}
	}
}

func DeleteSplit(c *gin.Context) {
	splitID := c.Param("split-id")
	var split models.Split
	db.First(&split, splitID)

	if split.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No split with ID: %v was found", splitID))
	} else if split.DeletedAt != nil {
		c.String(http.StatusBadRequest, "Resource has already been deleted")
	} else {
		db.Delete(&split)
		c.JSON(http.StatusOK, gin.H{
			"message": "Split succesfully deleted",
			"split": &split,
		})
	}
}