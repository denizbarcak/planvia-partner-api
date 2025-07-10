package main

import (
	"context"
	"log"
	"time"

	"planvia-partner-api/config"
	"planvia-partner-api/internal/database"
	"planvia-partner-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := database.ConnectDB(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Initialize handlers
	db := client.Database(cfg.DBName)
	partnerHandler := handlers.NewPartnerHandler(db)

	// Setup routes
	api := app.Group("/api")
	partners := api.Group("/partners")
	
	// Partner routes
	partners.Post("/register", partnerHandler.Register)
	partners.Post("/login", partnerHandler.Login)

	// Start server
	port := ":" + cfg.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 