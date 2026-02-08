package handler

import (
	"log"
	"net/http"
	"os"
	"sync"
	"tugas-session-1/database"
	"tugas-session-1/handlers"
	"tugas-session-1/repository"
	"tugas-session-1/router"
	"tugas-session-1/service"

	_ "github.com/lib/pq"
)

var (
	once        sync.Once
	mainHandler http.Handler
)

func initApp() {
	// Initialize database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARNING: DATABASE_URL environment variable is not set")
		return
	}

	err := database.Connect(dbURL)
	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v\n", err)
		return
	}

	// Run migrations
	err = database.RunMigrations(database.DB)
	if err != nil {
		log.Printf("ERROR: Failed to run migrations: %v\n", err)
	}

	// Initialize Category with service layer and dependency injection
	categoryRepo := repository.NewCategoryRepository(database.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize Product with service layer and dependency injection
	productRepo := repository.NewProductRepository(database.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize Transaction with service layer and dependency injection
	transactionRepo := repository.NewTransactionRepository(database.DB)
	transactionService := service.NewTransactionService(transactionRepo, productRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	mainHandler = router.SetupRoutes(categoryHandler, productHandler, transactionHandler)
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	
	if mainHandler == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error", "message": "App initialization failed"}`))
		return
	}

	mainHandler.ServeHTTP(w, r)
}

