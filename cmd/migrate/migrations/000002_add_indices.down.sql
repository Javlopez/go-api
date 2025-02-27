-- migrations/000002_add_indices.down.sql
-- Down: Remove indices
DROP INDEX IF EXISTS idx_orders_symbol;
DROP INDEX IF EXISTS idx_orders_created_at;