-- migrations/000002_add_indices.down.sql
-- Up: Add indices for faster queries
CREATE INDEX IF NOT EXISTS idx_orders_symbol ON orders(symbol);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);