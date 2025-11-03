package routers

import (
	"warehousemanagement/controllers"
	"warehousemanagement/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/auth/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/products", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is protected products route"})
	})

	protected.GET("/locations", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is protected locations route"})
	})

	protected.GET("/stock-movements", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is protected stock movements route"})
	})

	return r
}
