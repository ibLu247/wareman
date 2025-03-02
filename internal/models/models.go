package models

import "github.com/google/uuid"

// Модель склада
type Warehouse struct {
	ID      uuid.UUID `json:"id"`
	Address string    `json:"address"`
}

// Модель товара
type Product struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Сharacteristics map[string]string `json:"characteristics"`
	Weight          float32           `json:"weight"`
	Barcode         int               `json:"barcode"`
}

// Модель инвентаризации
type Inventory struct {
	ID              uuid.UUID `json:"id"`
	Quantity        int       `json:"quantity"`
	Price           float32   `json:"price"`
	Discount        float32   `json:"discount"`
	DiscountedPrice float32   `json:"discounted_price"`
	ProductID       uuid.UUID `json:"product_id"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
}

// Модель аналитики
type Analytic struct {
	ID                   uuid.UUID `json:"id"`
	WarehouseID          uuid.UUID `json:"warehouse_id"`
	ProductID            uuid.UUID `json:"product_id"`
	QuantitySoldProducts int       `json:"quantity_sold_products"`
	TotalSum             float32   `json:"total_sum"`
}
