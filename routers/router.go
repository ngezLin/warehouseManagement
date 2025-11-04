package routers

import (
	"warehousemanagement/config"
	controller "warehousemanagement/controller"
	"warehousemanagement/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/auth/login", controller.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	productController := controller.ProductController{DB: config.DB}
	protected.GET("/products", productController.GetProducts)
	protected.POST("/products", productController.CreateProduct)
	protected.PUT("/products/:id", productController.UpdateProduct)

	locationController := controller.LocationController{DB: config.DB}
	protected.GET("/locations", locationController.GetLocations)
	protected.POST("/locations", locationController.CreateLocation)

	stockController := controller.StockController{DB: config.DB}
	protected.GET("/stock-movements", stockController.GetStockMovements)
	protected.POST("/stock-movements", stockController.CreateStockMovement)

	return r
}
