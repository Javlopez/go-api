// pkg/models/order.go
package models

import (
	"time"
)

// OrderType represents the type of order
type OrderType string

const (
	Buy  OrderType = "BUY"
	Sell OrderType = "SELL"
)

// Order represents a trade order
type Order struct {
	ID        int64     `json:"id" db:"id"`
	Symbol    string    `json:"symbol" db:"symbol"`
	Price     float64   `json:"price" db:"price"`
	Quantity  int       `json:"quantity" db:"quantity"`
	OrderType OrderType `json:"order_type" db:"order_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// OrderRequest represents the order creation request
type OrderRequest struct {
	Symbol    string    `json:"symbol" binding:"required" example:"AAPL"`
	Price     float64   `json:"price" binding:"required,gt=0" example:"150.50"`
	Quantity  int       `json:"quantity" binding:"required,gt=0" example:"10"`
	OrderType OrderType `json:"order_type" binding:"required,oneof=BUY SELL" example:"BUY"`
}
