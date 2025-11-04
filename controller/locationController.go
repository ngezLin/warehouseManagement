package controllers

import (
	"database/sql"
	"net/http"

	"warehousemanagement/model"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	DB *sql.DB
}

func (lc *LocationController) GetLocations(c *gin.Context) {
	rows, err := lc.DB.Query("SELECT id, code, name, capacity FROM locations")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	locations := []model.Location{}
	for rows.Next() {
		var loc model.Location
		if err := rows.Scan(&loc.ID, &loc.Code, &loc.Name, &loc.Capacity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		locations = append(locations, loc)
	}

	c.JSON(http.StatusOK, gin.H{"data": locations})
}

func (lc *LocationController) CreateLocation(c *gin.Context) {
	var loc model.Location
	if err := c.ShouldBindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := lc.DB.Exec("INSERT INTO locations (code, name, capacity) VALUES (?, ?, ?)",
		loc.Code, loc.Name, loc.Capacity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	loc.ID = id

	c.JSON(http.StatusCreated, gin.H{
		"message": "Location created successfully",
		"data":    loc,
	})
}
