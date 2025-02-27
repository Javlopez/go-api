package api

import (
	"github.com/Javlopez/go-api/cmd/api/handlers"
	_ "github.com/Javlopez/go-api/docs"
	"github.com/Javlopez/go-api/pkg/repositories/order"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the Gin router
func SetupRouter(orderRepo order.OrderRepository) *gin.Engine {
	router := gin.Default()

	// Set up CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize API group
	api := router.Group("/api/v1")
	{
		// Initialize handlers
		orderHandler := handlers.NewOrderHandler(orderRepo)

		// Order routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.GetOrders)
	}

	url := ginSwagger.URL("/docs/doc.json") // The URL pointing to API definition
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
	return router
}
