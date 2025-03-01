package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ibLu247/wareman.git/internal/db"
	"github.com/ibLu247/wareman.git/internal/models"
)

func Healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

// Добавить склад
func AddWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.BindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный запрос"})
		return
	}
	id := uuid.New() // Генерация нового UUID
	db.Conn.Exec(context.Background(), "INSERT INTO warehouses (id, address) VALUES ($1, $2)", id, warehouse.Address)

	c.Status(http.StatusCreated)
}

// Получить список складов
func GetWarehouses(c *gin.Context) {
	values, _ := db.Conn.Query(context.Background(), "SELECT * FROM warehouses")
	var warehouses []models.Warehouse
	for values.Next() {
		var warehouse models.Warehouse
		values.Scan(&warehouse.ID, &warehouse.Address)
		warehouses = append(warehouses, warehouse)
	}

	c.JSON(http.StatusOK, warehouses)
}

// Добавить товар
func AddProduct(c *gin.Context) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный запрос"})
		return
	}
	id := uuid.New()
	characteristics, _ := json.Marshal(product.Сharacteristics)
	var sql string = "INSERT INTO products(id, name, description, characteristics, weight, barcode) VALUES($1, $2, $3, $4, $5, $6)"
	db.Conn.Exec(context.Background(), sql, id, product.Name, product.Description, characteristics, product.Weight, product.Barcode)

	c.Status(http.StatusCreated)
}

// Получить список всех товаров
func GetProducts(c *gin.Context) {
	values, _ := db.Conn.Query(context.Background(), "SELECT id, name, description, characteristics, weight, barcode FROM products")
	var products []models.Product
	for values.Next() {
		var product models.Product
		values.Scan(&product.ID, &product.Name, &product.Description, &product.Сharacteristics, &product.Weight, &product.Barcode)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// Обновить товар
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	var updatedProduct models.Product
	c.BindJSON(&updatedProduct)
	characteristicsJSON, _ := json.Marshal(updatedProduct.Сharacteristics)
	db.Conn.Exec(context.Background(), "UPDATE products SET description = $1, characteristics = $2 WHERE id = $3", updatedProduct.Description, characteristicsJSON, productID)

	c.Status(http.StatusOK)
}

// Создать инвентаризацию
func AddInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := c.BindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный запрос"})
		return
	}
	id := uuid.New()
	var sql string = "INSERT INTO inventory (id, quantity, price, discount, product_id, warehouse_id) VALUES ($1, $2, $3, $4, $5, $6)"
	db.Conn.Exec(context.Background(), sql, id, inventory.Quantity, inventory.Price, inventory.Discount, inventory.ProductID, inventory.WarehouseID)
	c.Status(http.StatusCreated)
}

// Обновить количество товара
func UpdateQuantity(c *gin.Context) {
	var inventory models.Inventory
	c.BindJSON(&inventory)
	var sql string = "UPDATE inventory SET quantity = quantity + $1 WHERE product_id = $2 AND warehouse_id = $3"
	db.Conn.Exec(context.Background(), sql, inventory.Quantity, inventory.ProductID, inventory.WarehouseID)

	c.Status(http.StatusOK)
}

func GetInventories(c *gin.Context) {
	values, _ := db.Conn.Query(context.Background(), "SELECT * FROM inventory")
	var inventories []models.Inventory
	for values.Next() {
		var inventory models.Inventory
		values.Scan(&inventory.ID, &inventory.Quantity, &inventory.Price, &inventory.Discount, &inventory.ProductID, &inventory.WarehouseID)
		inventories = append(inventories, inventory)
	}

	c.JSON(http.StatusOK, inventories)
}

// Создать скидку
func AddDiscount(c *gin.Context) {
	var inventory models.Inventory
	c.BindJSON(&inventory)
	var sql string = "UPDATE inventory SET discount = $1 WHERE product_id = $2 AND warehouse_id = $3"
	db.Conn.Exec(context.Background(), sql, inventory.Discount, inventory.ProductID, inventory.WarehouseID)

	c.Status(http.StatusOK)
}
