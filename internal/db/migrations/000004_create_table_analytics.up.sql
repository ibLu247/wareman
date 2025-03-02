CREATE TABLE analytics (
    id UUID PRIMARY KEY,
    warehouse_id UUID REFERENCES warehouses(id),
    product_id UUID REFERENCES products(id),
    quantity_sold_products INT,
    total_sum DECIMAL(10, 2)
);

CREATE UNIQUE INDEX idx_analytics_warehouse_product
ON analytics (warehouse_id, product_id);