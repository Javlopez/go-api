// ======== cmd/api/main.go ========
package main

import (
	"fmt"
	"github.com/Javlopez/go-api/cmd/api"
	"github.com/Javlopez/go-api/pkg/database"
	"github.com/Javlopez/go-api/pkg/repositories/order"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title Trade Orders API
// @version 1.0
// @description A simple API for managing trade orders
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize config
	cfg := database.NewConfig()

	// Initialize database
	db := database.New(cfg)

	// Connect to database
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	//defer db.Close()

	// Initialize repository
	orderRepo, err := order.NewOrderRepository(dbConnection)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer orderRepo.Close()

	// Initialize router
	router := api.SetupRouter(orderRepo)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
