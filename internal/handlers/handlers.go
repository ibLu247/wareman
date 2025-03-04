package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ibLu247/wareman.git/internal/db"
	"github.com/ibLu247/wareman.git/internal/models"
	"go.uber.org/zap"
)

// Функция проверяет на валидность данные в теле запроса
func ValidateJSON(c *gin.Context, value interface{}) bool {
	if err := c.BindJSON(value); err != nil {
		logger := c.MustGet("logger").(*zap.Logger)
		logger.Error("Невалидные данные", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{
			"Ошибка": "Невалидные данные",
		})
		return false
	}
	return true
}

func Healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

// Добавить склад
func AddWarehouse(c *gin.Context) {
	var warehouse models.Warehouse

	if !ValidateJSON(c, &warehouse) {
		return
	}

	id := uuid.New()

	var sql string = "INSERT INTO warehouses (id, address) VALUES ($1, $2)"
	db.Conn.Exec(context.Background(), sql, id, warehouse.Address)

	c.Status(http.StatusCreated)
}

// Получить список складов
func GetWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse
	var sql string = "SELECT * FROM warehouses"

	values, _ := db.Conn.Query(context.Background(), sql)

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

	if !ValidateJSON(c, &product) {
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
	var products []models.Product

	var sql string = "SELECT * FROM products"
	values, _ := db.Conn.Query(context.Background(), sql)

	for values.Next() {
		var product models.Product
		values.Scan(&product.ID, &product.Name, &product.Description, &product.Сharacteristics, &product.Weight, &product.Barcode)
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// Обновить товар
func UpdateProduct(c *gin.Context) {
	var updatedProduct models.Product

	productID := c.Param("id")

	if !ValidateJSON(c, &updatedProduct) {
		return
	}
	characteristicsJSON, _ := json.Marshal(updatedProduct.Сharacteristics)

	var sql string = "UPDATE products SET description = $1, characteristics = $2 WHERE id = $3"
	db.Conn.Exec(context.Background(), sql, updatedProduct.Description, characteristicsJSON, productID)

	c.Status(http.StatusOK)
}

// Создать инвентаризацию
func AddInventory(c *gin.Context) {
	var inventory models.Inventory

	if !ValidateJSON(c, &inventory) {
		return
	}
	id := uuid.New()

	var sql string = "INSERT INTO inventory (id, quantity, price, discount, discounted_price, product_id, warehouse_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	db.Conn.Exec(context.Background(), sql, id, inventory.Quantity, inventory.Price, inventory.Discount, inventory.DiscountedPrice, inventory.ProductID, inventory.WarehouseID)

	c.Status(http.StatusCreated)
}

// Поступление товара на склад
func UpdateQuantity(c *gin.Context) {
	var inventory models.Inventory

	if !ValidateJSON(c, &inventory) {
		return
	}

	var sql string = "UPDATE inventory SET quantity = quantity + $1 WHERE product_id = $2 AND warehouse_id = $3"
	db.Conn.Exec(context.Background(), sql, inventory.Quantity, inventory.ProductID, inventory.WarehouseID)

	c.Status(http.StatusOK)
}

// Создать скидку
func AddDiscount(c *gin.Context) {
	var inventory models.Inventory

	if !ValidateJSON(c, &inventory) {
		return
	}

	var price float32

	var sql string = "SELECT price FROM inventory WHERE product_id = $1 AND warehouse_id = $2"
	db.Conn.QueryRow(context.Background(), sql, inventory.ProductID, inventory.WarehouseID).Scan(&price)

	var discountedPrice float32 = price - ((price / 100) * inventory.Discount)

	sql = "UPDATE inventory SET discount = $1, discounted_price = $2 WHERE product_id = $3 AND warehouse_id = $4"
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
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	var offset int = (page - 1) * limit

	var sql string = `
						SELECT products.id, products.name, inventory.price, inventory.discounted_price 
						FROM products JOIN inventory ON inventory.product_id = products.id 
						WHERE inventory.warehouse_id = $1
						LIMIT $2 OFFSET $3
					`
	values, _ := db.Conn.Query(context.Background(), sql, warehouseID, limit, offset)

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

	var sql string = `
						SELECT products.id, products.name, products.description, products.characteristics, products.weight, products.barcode, inventory.price, inventory.discounted_price, inventory.quantity 
						FROM products JOIN inventory ON inventory.product_id = products.id 
						WHERE inventory.product_id = $1 AND inventory.warehouse_id = $2
					`
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
	if !ValidateJSON(c, &req) {
		return
	}

	var res []Res
	var price float32

	for _, item := range req {
		var sql string = "SELECT price FROM inventory WHERE warehouse_id = $1 AND product_id = $2"
		values, _ := db.Conn.Query(context.Background(), sql, warehouseID, item.ProductID)
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
	if !ValidateJSON(c, &req) {
		return
	}

	for _, item := range req {
		var currentProducts int
		var discountedPrice float32
		var price float32

		var sql string = "SELECT quantity, discounted_price, price FROM inventory WHERE warehouse_id = $1 AND product_id = $2"
		db.Conn.QueryRow(context.Background(), sql, warehouseID, item.ProductID).Scan(&currentProducts, &discountedPrice, &price)

		if currentProducts < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Недостаточно товаров на складе"})
			return
		}

		db.Conn.Exec(context.Background(), "UPDATE inventory SET quantity = quantity - $1 WHERE warehouse_id = $2 AND product_id = $3", item.Quantity, warehouseID, item.ProductID)

		// Запись аналитики при покупке
		CreateAnalytics(warehouseID, item.ProductID, item.Quantity, price, discountedPrice)
	}

	c.Status(http.StatusOK)
}

// Создать аналитику при покупке товара
func CreateAnalytics(warehouseID string, productID uuid.UUID, quantity int, price float32, discountedPrice float32) {

	var currentPrice float32 // Переменная хранит в себе либо цену без скидки либо цену со скидкой
	if discountedPrice != 0 {
		currentPrice = discountedPrice
	} else {
		currentPrice = price
	}

	var total_sum float32 = currentPrice * float32(quantity) // Переменная хранит в себе сумму на которую была совершена покупка

	id := uuid.New()

	var sql string = `
						INSERT INTO analytics (id, warehouse_id, product_id, quantity_sold_products, total_sum) 
						VALUES ($1, $2, $3, $4, $5) 
						ON CONFLICT (warehouse_id, product_id) 
						DO UPDATE SET 
							quantity_sold_products = analytics.quantity_sold_products + EXCLUDED.quantity_sold_products, 
							total_sum = analytics.total_sum + EXCLUDED.total_sum
					`
	db.Conn.Exec(context.Background(), sql, id, warehouseID, productID, quantity, total_sum)
}

// Получить аналитику по складу по каждому товару
func GetAnalytic(c *gin.Context) {
	type Analytic struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
		TotalSum  float32   `json:"total_sum"`
	}

	warehouseID := c.Param("id")

	var analytics []Analytic

	var sql string = "SELECT product_id, quantity_sold_products, total_sum FROM analytics WHERE warehouse_id = $1"
	values, _ := db.Conn.Query(context.Background(), sql, warehouseID)

	for values.Next() {
		var analytic Analytic
		values.Scan(&analytic.ProductID, &analytic.Quantity, &analytic.TotalSum)
		analytics = append(analytics, analytic)
	}

	c.JSON(http.StatusOK, analytics)
}

// Получить топ 10 складов которые сделали больше всего выручки
func GetMostWarehouses(c *gin.Context) {
	type Analytic struct {
		WarehouseID uuid.UUID `json:"warehouse_id"`
		Address     string    `json:"address"`
		TotalSum    float32   `json:"total_sum"`
	}

	var sql string = `
						SELECT warehouses.id, warehouses.address, SUM(analytics.total_sum) AS total
						FROM analytics
						JOIN warehouses ON analytics.warehouse_id = warehouses.id
						GROUP BY warehouses.id, warehouses.address
						ORDER BY total DESC
						LIMIT 10;
					`
	values, _ := db.Conn.Query(context.Background(), sql)

	var analytics []Analytic
	for values.Next() {
		var analytic Analytic
		values.Scan(&analytic.WarehouseID, &analytic.Address, &analytic.TotalSum)
		analytics = append(analytics, analytic)
	}

	c.JSON(http.StatusOK, analytics)
}
