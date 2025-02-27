// test/integration/api_test.go
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Javlopez/go-api/cmd/api/handlers"
	"github.com/Javlopez/go-api/pkg/models"
	"github.com/Javlopez/go-api/pkg/repositories/order"
	"github.com/Javlopez/go-api/pkg/testutils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pgContainer *testutils.PostgresContainer
	testRepo    order.OrderRepository
	router      *gin.Engine
)

func TestMain(m *testing.M) {
	// Setup
	ctx := context.Background()

	// Start Postgres container
	var err error
	pgContainer, err = testutils.NewPostgresContainer(ctx)
	if err != nil {
		fmt.Printf("Failed to start Postgres container: %v\n", err)
		os.Exit(1)
	}

	// Setup database schema
	if err := pgContainer.SetupOrdersTable(); err != nil {
		fmt.Printf("Failed to set up test database: %v\n", err)
		pgContainer.Terminate(ctx)
		os.Exit(1)
	}

	// Initialize repository
	testRepo = &order.PostgresOrderRepository{DB: pgContainer.DB}

	// Configure router
	router = setupRouter()

	// Run tests
	code := m.Run()

	// Teardown
	pgContainer.Close()
	pgContainer.Terminate(ctx)

	// Exit
	os.Exit(code)
}

// setupRouter configures the test router
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(testRepo)

	// Set up routes
	api := r.Group("/api/v1")
	{
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.GetOrders)
	}

	return r
}

// TestCreateAndGetOrders tests creating an order and retrieving it
func TestCreateAndGetOrders(t *testing.T) {
	// Clean up any existing data first
	pgContainer.CleanupData()

	// Test data
	orderRequest := models.OrderRequest{
		Symbol:    "AAPL",
		Price:     150.5,
		Quantity:  10,
		OrderType: models.Buy,
	}

	// 1. Create an order
	jsonData, _ := json.Marshal(orderRequest)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	// Parse the response
	var createdOrder models.Order
	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)
	require.NoError(t, err)
	assert.Equal(t, "AAPL", createdOrder.Symbol)
	assert.Equal(t, 150.5, createdOrder.Price)
	assert.Equal(t, 10, createdOrder.Quantity)
	assert.Equal(t, models.Buy, createdOrder.OrderType)
	assert.Greater(t, createdOrder.ID, int64(0))

	req, _ = http.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var orders []models.Order
	err = json.Unmarshal(w.Body.Bytes(), &orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, createdOrder.ID, orders[0].ID)
	assert.Equal(t, "AAPL", orders[0].Symbol)
}

// TestCreateOrderValidation tests validation on order creation
func TestCreateOrderValidation(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		request     models.OrderRequest
		expectedErr string
	}{
		{
			name: "Missing Symbol",
			request: models.OrderRequest{
				Price:     150.5,
				Quantity:  10,
				OrderType: models.Buy,
			},
			expectedErr: "symbol",
		},
		{
			name: "Invalid Price",
			request: models.OrderRequest{
				Symbol:    "AAPL",
				Price:     -10, // Invalid: price must be > 0
				Quantity:  10,
				OrderType: models.Buy,
			},
			expectedErr: "price",
		},
		{
			name: "Invalid Quantity",
			request: models.OrderRequest{
				Symbol:    "AAPL",
				Price:     150.5,
				Quantity:  0, // Invalid: quantity must be > 0
				OrderType: models.Buy,
			},
			expectedErr: "quantity",
		},
		{
			name: "Invalid Order Type",
			request: models.OrderRequest{
				Symbol:    "AAPL",
				Price:     150.5,
				Quantity:  10,
				OrderType: "INVALID", // Invalid: must be BUY or SELL
			},
			expectedErr: "ordertype",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Send request
			jsonData, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check response
			require.Equal(t, http.StatusBadRequest, w.Code)

			// Parse the validation errors
			var response models.ValidationErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			foundError := false
			for _, validationErr := range response.Errors {
				if validationErr.Field == tc.expectedErr {
					foundError = true
					break
				}
			}
			assert.True(t, foundError, "Expected validation error for field: %s", tc.expectedErr)
		})
	}
}
