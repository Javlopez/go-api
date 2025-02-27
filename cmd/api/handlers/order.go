// pkg/api/handlers/order_handler.go
package handlers

import (
	"errors"
	"fmt"
	"github.com/Javlopez/go-api/pkg/repositories/order"
	"net/http"
	"strings"

	"github.com/Javlopez/go-api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	repo order.OrderRepository
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repo order.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// CreateOrder godoc
// @Summary Create a new trade order
// @Description Create a new trade order with the provided details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.OrderRequest true "Order details"
// @Success 201 {object} models.Order
// @Failure 400 {object} models.ValidationErrorResponse "Validation failed"
// @Failure 500 {object} models.ErrorResponse "Server error"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderRequest models.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		var validationErrors models.ValidationErrorResponse

		// Check if this is a validation error
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			// Process validation errors
			for _, e := range validationErrs {
				field := strings.ToLower(e.Field())
				var message string

				// Create user-friendly error messages
				switch e.Tag() {
				case "required":
					message = fmt.Sprintf("%s is required", field)
				case "gt":
					message = fmt.Sprintf("%s must be greater than %s", field, e.Param())
				case "oneof":
					message = fmt.Sprintf("%s must be one of: %s", field, e.Param())
				default:
					message = fmt.Sprintf("%s failed validation: %s", field, e.Tag())
				}

				validationErrors.Errors = append(validationErrors.Errors, models.ValidationError{
					Field:   field,
					Message: message,
				})
			}
		}

		c.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	order := models.Order{
		Symbol:    orderRequest.Symbol,
		Price:     orderRequest.Price,
		Quantity:  orderRequest.Quantity,
		OrderType: orderRequest.OrderType,
	}

	if err := h.repo.Create(&order); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrders godoc
// @Summary Get all trade orders
// @Description Retrieve a list of all submitted trade orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} models.ErrorResponse "Server error"
// @Router /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}
