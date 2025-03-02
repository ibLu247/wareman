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

	var sql string = "INSERT INTO inventory (id, quantity, price, discount, discounted_price, product_id, warehouse_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	db.Conn.Exec(context.Background(), sql, id, inventory.Quantity, inventory.Price, inventory.Discount, inventory.DiscountedPrice, inventory.ProductID, inventory.WarehouseID)

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

// Получить список инвентаризаций (ВРЕМЕННО)
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
	var price float32
	db.Conn.QueryRow(context.Background(), "SELECT price FROM inventory").Scan(&price)

	var discountedPrice float32 = price - ((price / 100) * inventory.Discount)

	var sql string = "UPDATE inventory SET discount = $1, discounted_price = $2 WHERE product_id = $3 AND warehouse_id = $4"
	db.Conn.Exec(context.Background(), sql, inventory.Discount, discountedPrice, inventory.ProductID, inventory.WarehouseID)

	c.Status(http.StatusOK)
}

// Получить список товаров по конкретному складу с пагинацией
func GetProductsFromWarehouse(c *gin.Context) {
	type Product struct {
		ID              uuid.UUID `json:"product_id"`
		Name            string    `json:"name"`
		Price           float32   `json:"price"`
		DiscountedPrice float32   `json:"discounted_price"`
	}

	warehouseID := c.Query("warehouse_id")

	var sql string = "SELECT products.id, products.name, inventory.price, inventory.discounted_price FROM products JOIN inventory ON inventory.product_id = products.id WHERE inventory.warehouse_id = $1"
	values, _ := db.Conn.Query(context.Background(), sql, warehouseID)

	var products []Product
	for values.Next() {
		var product Product
		values.Scan(&product.ID, &product.Name, &product.Price, &product.DiscountedPrice)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// Получить всю информаицю о товаре на складе
func GetProductFromWarehouse(c *gin.Context) {
	type Product struct {
		ID              uuid.UUID         `json:"product_id"`
		Name            string            `json:"name"`
		Description     string            `json:"description"`
		Сharacteristics map[string]string `json:"characteristics"`
		Weight          float32           `json:"weight"`
		Barcode         int               `json:"barcode"`
		Price           float32           `json:"price"`
		DiscountedPrice float32           `json:"discounted_price"`
		Quantity        int               `json:"quantity"`
	}

	productID := c.Param("id")
	warehouseID := c.Query("warehouse_id")

	var sql string = "SELECT products.id, products.name, products.description, products.characteristics, products.weight, products.barcode, inventory.price, inventory.discounted_price, inventory.quantity FROM products JOIN inventory ON inventory.product_id = products.id WHERE inventory.product_id = $1 AND inventory.warehouse_id = $2"
	values, _ := db.Conn.Query(context.Background(), sql, productID, warehouseID)

	var product Product
	var characteristics string

	values.Next()
	values.Scan(&product.ID, &product.Name, &product.Description, &characteristics, &product.Weight, &product.Barcode, &product.Price, &product.DiscountedPrice, &product.Quantity)
	json.Unmarshal([]byte(characteristics), &product.Сharacteristics)

	c.JSON(http.StatusOK, product)
}

// Получить подсумировку цен списка
func GetSum(c *gin.Context) {
	type Req struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	type Res struct {
		ProductID  uuid.UUID `json:"product_id"`
		TotalPrice float32   `json:"total_price"`
	}

	warehouseID := c.Param("id")

	var req []Req
	c.BindJSON(&req)

	var res []Res
	var price float32

	for _, item := range req {
		values, _ := db.Conn.Query(context.Background(), "SELECT price FROM inventory WHERE warehouse_id = $1 AND product_id = $2", warehouseID, item.ProductID)
		values.Next()
		values.Scan(&price)
		totalPrice := float32(item.Quantity) * price
		res = append(res, Res{
			ProductID:  item.ProductID,
			TotalPrice: totalPrice,
		})
	}

	c.JSON(http.StatusOK, res)
}

// Покупка товаров
func BuyProducts(c *gin.Context) {
	type Req struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	warehouseID := c.Param("id")

	var req []Req
	c.BindJSON(&req)

	for _, item := range req {
		var currentProducts int
		var discountedPrice float32
		var price float32
		var currentPrice float32

		var sql string = "SELECT quantity, discounted_price, price FROM inventory WHERE warehouse_id = $1 AND product_id = $2"
		db.Conn.QueryRow(context.Background(), sql, warehouseID, item.ProductID).Scan(&currentProducts, &discountedPrice, &price)

		if currentProducts < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Недостаточно товаров на складе"})
			return
		}

		if discountedPrice != 0 {
			currentPrice = discountedPrice
		} else {
			currentPrice = price
		}
		db.Conn.Exec(context.Background(), "UPDATE inventory SET quantity = quantity - $1 WHERE warehouse_id = $2 AND product_id = $3", item.Quantity, warehouseID, item.ProductID)

		id := uuid.New()
		sql = `
    			INSERT INTO analytics (id, warehouse_id, product_id, quantity_sold_products, total_sum) 
    			VALUES ($1, $2, $3, $4, $5) 
				ON CONFLICT (warehouse_id, product_id) 
    			DO UPDATE SET 
        			quantity_sold_products = analytics.quantity_sold_products + EXCLUDED.quantity_sold_products, 
        			total_sum = analytics.total_sum + EXCLUDED.total_sum
			`
		db.Conn.Exec(context.Background(), sql, id, warehouseID, item.ProductID, item.Quantity, currentPrice)
	}

	c.Status(http.StatusOK)
}
