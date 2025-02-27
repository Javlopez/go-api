package order

import (
	"errors"
	"github.com/Javlopez/go-api/pkg/models"
	"github.com/jmoiron/sqlx"
	"time"
)

// PostgresOrderRepository is an implementation of OrderRepository
type PostgresOrderRepository struct {
	DB *sqlx.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sqlx.DB) (OrderRepository, error) {
	return &PostgresOrderRepository{DB: db}, nil
}

// Create inserts a new order
func (r *PostgresOrderRepository) Create(order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	// Set created_at to current time
	order.CreatedAt = time.Now()

	query := `
		INSERT INTO orders (symbol, price, quantity, order_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	return r.DB.QueryRow(
		query,
		order.Symbol,
		order.Price,
		order.Quantity,
		order.OrderType,
		order.CreatedAt,
	).Scan(&order.ID)
}

// GetAll retrieves all orders
func (r *PostgresOrderRepository) GetAll() ([]models.Order, error) {
	orders := []models.Order{}
	query := `
		SELECT id, symbol, price, quantity, order_type, created_at
		FROM orders
		ORDER BY created_at DESC
	`

	err := r.DB.Select(&orders, query)
	return orders, err
}

// Close closes the database connection
func (r *PostgresOrderRepository) Close() error {
	return r.DB.Close()
}
