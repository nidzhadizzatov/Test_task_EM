package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"subscription-service/internal/api/handlers"
	"subscription-service/internal/config"
	"subscription-service/internal/logger"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize logger
	logger := logger.NewLogger(cfg.LogLevel)
	logger.Info("Starting subscription service...")

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	logger.Info("Database connection established")

	// Initialize repository
	repo := repository.NewPostgresRepository(db)

	// Initialize service
	subscriptionService := service.NewSubscriptionService(repo)

	// Initialize handlers
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// Initialize Gin router
	r := gin.Default()

	// Add middleware for CORS and request logging
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register routes
	subscriptionHandler.RegisterRoutes(r)

	// Start server
	serverAddr := ":" + cfg.ServerPort
	logger.Infof("Server starting on port %s", cfg.ServerPort)
	logger.Info("Available endpoints:")
	logger.Info("  POST /api/v1/subscriptions - Create subscription")
	logger.Info("  GET /api/v1/subscriptions - Get all subscriptions")
	logger.Info("  GET /api/v1/subscriptions/:id - Get subscription by ID")
	logger.Info("  PUT /api/v1/subscriptions/:id - Update subscription")
	logger.Info("  DELETE /api/v1/subscriptions/:id - Delete subscription")
	logger.Info("  GET /api/v1/subscriptions/cost - Calculate total cost with filters")
	logger.Info("  GET /health - Health check")

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
