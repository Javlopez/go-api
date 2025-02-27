-- migrations/000001_create_orders_table.up.sql
-- Up: Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    price DECIMAL(12, 4) NOT NULL,
    quantity INTEGER NOT NULL,
    order_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );