CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    characteristics JSONB,
    weight FLOAT,
    barcode INT UNIQUE
);