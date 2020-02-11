package main

import(
	"fmt"
	"strconv"
	"net/http"

	"github.com/llewellyn-kevin/lcs-sm/models"

	"github.com/gin-gonic/gin"
)

// TODO: Currently returns values with null Team and Split member vars. Fix.
func GetStockValues(c *gin.Context) {
	var stockValues []models.StockValue
	db.Find(&stockValues)
	for _, sv := range stockValues {
		db.First(&sv.Team, sv.TeamID)
		db.First(&sv.Split, sv.SplitID)
	}
	c.JSON(http.StatusOK, &stockValues)
}

func GetStockValue(c *gin.Context) {
	stockValueID := c.Param("stock-value-id")
	var stockValue models.StockValue
	db.First(&stockValue, stockValueID)
	db.First(&stockValue.Team, stockValue.TeamID)
	db.First(&stockValue.Split, stockValue.SplitID)
	
	if stockValue.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No stockValue with ID: %v was found", stockValueID))
	} else {
		c.JSON(http.StatusOK, &stockValue)
	}
}

func CreateStockValue(c *gin.Context) {
	value, valueErr := strconv.Atoi(c.Query("value"))
	week, weekErr := strconv.Atoi(c.Query("week"))
	splitID, splitErr := strconv.Atoi(c.Param("split-id"))
	teamID, teamErr := strconv.Atoi(c.Param("team-id"))

	var split models.Split
	var team models.Team

	db.First(&split, splitID)
	db.First(&team, teamID)

	switch {
	case valueErr != nil:
		c.String(http.StatusBadRequest, fmt.Sprintf("Unable to convert parameter value to int. %v", valueErr))
	case weekErr != nil:
		c.String(http.StatusBadRequest, fmt.Sprintf("Unable to convert parameter week to int. %v", weekErr))
	case splitErr != nil:
		c.String(http.StatusBadRequest, fmt.Sprintf("Unable to convert parameter split to int. %v", splitErr))
	case teamErr != nil:
		c.String(http.StatusBadRequest, fmt.Sprintf("Unable to convert parameter team to int. %v", teamErr))
	default:
		stockValue := models.StockValue{
			Value:   value, 
			Week:    week, 
			SplitID: splitID, 
			Split:   split, 
			TeamID:  teamID, 
			Team:    team,
		}

		db.Create(&stockValue)

		if !db.NewRecord(stockValue) {
			c.JSON(http.StatusOK, gin.H{
				"message": "StockValue succesfully created",
				"stockValue": &stockValue,
			})
		} else {
			c.String(http.StatusInternalServerError, "Failed to create resource")
		}
	}
}

func DeleteStockValue(c *gin.Context) {
	stockValueID := c.Param("stockValue-id")
	var stockValue models.StockValue
	db.First(&stockValue, stockValueID)

	if stockValue.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No stockValue with ID: %v was found", stockValueID))
	} else if stockValue.DeletedAt != nil {
		c.String(http.StatusBadRequest, "Resource has already been deleted")
	} else {
		db.Delete(&stockValue)
		c.JSON(http.StatusOK, gin.H{
			"message": "StockValue succesfully deleted",
			"stockValue": &stockValue,
		})
	}
}

func UpdateStockValue(c *gin.Context) {
	stockValueID := c.Param("stock-value-id")
	value, err := strconv.Atoi(c.Query("new-value"))
	var stockValue models.StockValue
	db.First(&stockValue, stockValueID)

	if stockValue.ID == 0 {
		c.String(http.StatusNotFound, fmt.Sprintf("No stockValue with ID: %v was found", stockValueID))
	} else if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not convert %v to integer value", c.Query("new-value")))
	} else {
		stockValue.Value = value
		db.Save(stockValue)

		c.JSON(http.StatusOK, gin.H{
			"message": "StockValue value updated",
			"stockValue": &stockValue,
		})
	}
}