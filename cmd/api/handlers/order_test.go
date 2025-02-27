package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Javlopez/go-api/pkg/models"
)

// MockOrderRepository is a mock implementation of OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetAll() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateOrderHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Create test order request
	orderRequest := models.OrderRequest{
		Symbol:    "AAPL",
		Price:     150.5,
		Quantity:  10,
		OrderType: models.Buy,
	}

	// Setup expectations
	mockRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

	// Prepare request
	jsonData, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.POST("/api/v1/orders", handler.CreateOrder)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateOrderInvalidJSON(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Prepare invalid JSON request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.POST("/api/v1/orders", handler.CreateOrder)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrderValidationFailed(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Create invalid order request (missing required fields)
	orderRequest := models.OrderRequest{
		// Symbol is required but missing
		Price:     150.5,
		Quantity:  10,
		OrderType: models.Buy,
	}

	// Prepare request
	jsonData, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.POST("/api/v1/orders", handler.CreateOrder)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check response body contains validation error
	var response models.ValidationErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Greater(t, len(response.Errors), 0)
}

func TestCreateOrderDatabaseError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Create test order request
	orderRequest := models.OrderRequest{
		Symbol:    "AAPL",
		Price:     150.5,
		Quantity:  10,
		OrderType: models.Buy,
	}

	// Setup expectations with an error
	mockRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(errors.New("database error"))

	// Prepare request
	jsonData, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.POST("/api/v1/orders", handler.CreateOrder)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetOrdersHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Create test orders
	now := time.Now()
	orders := []models.Order{
		{ID: 1, Symbol: "AAPL", Price: 150.5, Quantity: 10, OrderType: models.Buy, CreatedAt: now},
		{ID: 2, Symbol: "MSFT", Price: 250.75, Quantity: 5, OrderType: models.Sell, CreatedAt: now},
	}

	// Setup expectations
	mockRepo.On("GetAll").Return(orders, nil)

	// Prepare request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders", nil)

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.GET("/api/v1/orders", handler.GetOrders)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains orders
	var response []models.Order
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "AAPL", response[0].Symbol)
	assert.Equal(t, "MSFT", response[1].Symbol)

	mockRepo.AssertExpectations(t)
}

func TestGetOrdersDatabaseError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(MockOrderRepository)

	// Create handler with mock repo
	handler := NewOrderHandler(mockRepo)

	// Setup expectations with an error
	mockRepo.On("GetAll").Return([]models.Order{}, errors.New("database error"))

	// Prepare request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders", nil)

	// Prepare response recorder
	w := httptest.NewRecorder()

	// Setup Gin router
	router := gin.Default()
	router.GET("/api/v1/orders", handler.GetOrders)

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}
