package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibLu247/wareman.git/internal/db"
	"github.com/ibLu247/wareman.git/internal/handlers"
	"github.com/ibLu247/wareman.git/internal/logger"
	"github.com/ibLu247/wareman.git/internal/middleware"
)

func main() {
	log := logger.NewLogger()
	db.ConnectDB()

	router := gin.Default()

	router.Use(middleware.AddReqID(log))

	router.GET("/api/health", handlers.Healthcheck)

	router.POST("/api/warehouse", handlers.AddWarehouse)
	router.GET("/api/warehouses", handlers.GetWarehouses)

	router.POST("/api/product", handlers.AddProduct)
	router.GET("/api/products", handlers.GetProducts)
	router.PATCH("/api/product/:id", handlers.UpdateProduct)

	router.POST("/api/inventory", handlers.AddInventory)
	router.PATCH("/api/inventory", handlers.UpdateQuantity)
	router.PATCH("/api/inventory/discount", handlers.AddDiscount)
	router.GET("/api/inventory", handlers.GetProductsFromWarehouse)
	router.GET("/api/inventory/:id", handlers.GetProductFromWarehouse)
	router.POST("/api/inventory/:id", handlers.GetSum)
	router.POST("/api/inventory/product/:id", handlers.BuyProducts)

	router.GET("/api/analytic/:id", handlers.GetAnalytic)
	router.GET("/api/analytics", handlers.GetMostWarehouses)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db.DisconnectDB()

	_ = server.Shutdown(ctx)
}
