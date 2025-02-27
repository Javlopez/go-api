package order

import "github.com/Javlopez/go-api/pkg/models"

// OrderRepository interface for order operations
type OrderRepository interface {
	Create(order *models.Order) error
	GetAll() ([]models.Order, error)
	Close() error
}
