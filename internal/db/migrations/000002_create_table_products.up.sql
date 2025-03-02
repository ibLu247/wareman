CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    characteristics JSONB,
    weight DECIMAL(10, 2),
    barcode INT UNIQUE
);