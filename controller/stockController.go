package controllers

import (
	"database/sql"
	"net/http"
	"time"
	"warehousemanagement/model"

	"github.com/gin-gonic/gin"
)

type StockController struct {
	DB *sql.DB
}

func (s *StockController) GetStockMovements(c *gin.Context) {
	rows, err := s.DB.Query(`
		SELECT id, product_id, location_id, type, quantity, created_at 
		FROM stock_movements ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var movements []model.StockMovement
	for rows.Next() {
		var m model.StockMovement
		var createdAtStr string

		err := rows.Scan(&m.ID, &m.ProductID, &m.LocationID, &m.Type, &m.Quantity, &createdAtStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err == nil {
			m.CreatedAt = parsedTime
		}

		movements = append(movements, m)
	}

	c.JSON(http.StatusOK, gin.H{"stock_movements": movements})
}

func (s *StockController) CreateStockMovement(c *gin.Context) {
	var movement model.StockMovement
	if err := c.ShouldBindJSON(&movement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if movement.Type != "IN" && movement.Type != "OUT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be IN or OUT"})
		return
	}

	var currentStock int
	err := s.DB.QueryRow(`SELECT quantity FROM products WHERE id = ?`, movement.ProductID).Scan(&currentStock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id"})
		return
	}

	var capacity int
	err = s.DB.QueryRow(`SELECT capacity FROM locations WHERE id = ?`, movement.LocationID).Scan(&capacity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location_id"})
		return
	}

	if movement.Type == "OUT" && movement.Quantity > currentStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock OUT exceeds available stock"})
		return
	}

	if movement.Type == "IN" && currentStock+movement.Quantity > capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock IN exceeds location capacity"})
		return
	}

	movement.CreatedAt = time.Now()
	_, err = s.DB.Exec(`
		INSERT INTO stock_movements (product_id, location_id, type, quantity, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, movement.ProductID, movement.LocationID, movement.Type, movement.Quantity, movement.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if movement.Type == "IN" {
		currentStock += movement.Quantity
	} else {
		currentStock -= movement.Quantity
	}

	_, err = s.DB.Exec(`UPDATE products SET quantity = ? WHERE id = ?`, currentStock, movement.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Stock movement recorded successfully",
		"movement": movement,
	})
}
