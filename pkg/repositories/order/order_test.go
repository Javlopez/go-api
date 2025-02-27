package order

import (
	"github.com/Javlopez/go-api/pkg/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateOrder(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create sqlx database with the mock
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Create repository with the mock
	repo := &PostgresOrderRepository{DB: sqlxDB}

	// Create test order
	now := time.Now()
	order := &models.Order{
		Symbol:    "AAPL",
		Price:     150.5,
		Quantity:  10,
		OrderType: models.Buy,
		CreatedAt: now,
	}

	// Setup expectations
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.Symbol, order.Price, order.Quantity, order.OrderType, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the Create method
	err = repo.Create(order)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), order.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllOrders(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create sqlx database with the mock
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Create repository with the mock
	repo := &PostgresOrderRepository{DB: sqlxDB}

	// Create time for testing
	now := time.Now()

	// Setup expected rows
	rows := sqlmock.NewRows([]string{"id", "symbol", "price", "quantity", "order_type", "created_at"}).
		AddRow(1, "AAPL", 150.5, 10, models.Buy, now).
		AddRow(2, "MSFT", 250.75, 5, models.Sell, now)

	// Setup expectations
	mock.ExpectQuery("SELECT (.+) FROM orders").WillReturnRows(rows)

	// Call the GetAll method
	orders, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, "AAPL", orders[0].Symbol)
	assert.Equal(t, "MSFT", orders[1].Symbol)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrderNilOrder(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create sqlx database with the mock
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Create repository with the mock
	repo := &PostgresOrderRepository{DB: sqlxDB}

	// Call the Create method with nil
	err = repo.Create(nil)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
