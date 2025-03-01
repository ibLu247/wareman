package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ibLu247/wareman.git/internal/db"
	"github.com/ibLu247/wareman.git/internal/handlers"
)

func main() {
	db.ConnectDB()
	defer db.DisconnectDB()
	router := gin.Default()

	router.GET("/api/health", handlers.Healthcheck)
	router.POST("/api/warehouse", handlers.AddWarehouse)
	router.GET("/api/warehouses", handlers.GetWarehouses)
	router.POST("/api/product", handlers.AddProduct)
	router.GET("/api/products", handlers.GetProducts)
	router.PATCH("/api/product/:id", handlers.UpdateProduct)
	router.POST("/api/inventory", handlers.AddInventory)
	router.PATCH("/api/inventory", handlers.UpdateQuantity)
	router.PATCH("/api/inventory-discount", handlers.AddDiscount)
	router.GET("/api/inventories", handlers.GetInventories)
	router.Run()
}
