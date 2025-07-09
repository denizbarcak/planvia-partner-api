package main

import (
	"log"

	"planvia-partner-api/config"
	"planvia-partner-api/internal/database"
	"planvia-partner-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := database.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	// Initialize Fiber app
	app := fiber.New()

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Initialize handlers
	db := client.Database(cfg.DBName)
	partnerHandler := handlers.NewPartnerHandler(db)

	// Setup routes
	api := app.Group("/api")
	partners := api.Group("/partners")
	partners.Post("/register", partnerHandler.Register)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 