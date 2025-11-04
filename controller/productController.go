package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	model "warehousemanagement/model"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	DB *sql.DB
}

func (p *ProductController) GetProducts(c *gin.Context) {
	rows, err := p.DB.Query("SELECT id, sku_name, quantity FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var prod model.Product
		if err := rows.Scan(&prod.ID, &prod.SKUName, &prod.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, prod)
	}

	c.JSON(http.StatusOK, products)
}

func (p *ProductController) CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := p.DB.Exec(
		"INSERT INTO products (sku_name, quantity) VALUES (?, ?)",
		product.SKUName, product.Quantity,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}

	product.ID = int(id)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (p *ProductController) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	res, err := p.DB.Exec(
		"UPDATE products SET sku_name=?, quantity=? WHERE id=?",
		product.SKUName, product.Quantity, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.ID = id

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}
