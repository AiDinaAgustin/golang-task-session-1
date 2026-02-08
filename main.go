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
	"tugas-session-1/service"

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

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

	// Initialize repository and handlers
	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize product with service layer and dependency injection
	productRepo := repository.NewProductRepository(database.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize transaction with service layer and dependency injection
	transactionRepo := repository.NewTransactionRepository(database.DB)
	transactionService := service.NewTransactionService(transactionRepo, productRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	handler := router.SetupRoutes(categoryHandler, productHandler, transactionHandler)

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
