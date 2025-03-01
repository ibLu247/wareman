CREATE TABLE inventory (
    id UUID PRIMARY KEY,
    quantity INT,
    price DECIMAL(10, 2),
    discount DECIMAL(10, 2),
    product_id UUID REFERENCES products(id),
    warehouse_id UUID REFERENCES warehouses(id)
);