package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tugas-session-1/database"
	"tugas-session-1/handlers"
	"tugas-session-1/repository"
	"tugas-session-1/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	if err := database.Connect(dbURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repository and handlers
	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	// Setup routes
	handler := router.SetupRoutes(categoryHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
