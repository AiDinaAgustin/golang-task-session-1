package router

import (
	"encoding/json"
	"net/http"
	"strings"
	"tugas-session-1/handlers"
)

// SetupRoutes configures all application routes
func SetupRoutes(categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler) http.Handler {
	mux := http.NewServeMux()

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only handle exact root path
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message": "Welcome to Category & Product CRUD API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"GET /":                   "This welcome message",
				"GET /health":             "Health check",
				"GET /categories":         "Get all categories",
				"POST /categories":        "Create a category",
				"GET /categories/{id}":    "Get category by ID",
				"PUT /categories/{id}":    "Update category",
				"DELETE /categories/{id}": "Delete category",
				"GET /products":           "Get all products",
				"POST /products":          "Create a product",
				"GET /products/{id}":      "Get product by ID",
				"PUT /products/{id}":      "Update product",
				"DELETE /products/{id}":   "Delete product",
				"GET /docs":               "Swagger API Documentation",
			},
			"documentation": "/docs",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Category routes
	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAllCategories(w, r)
		case http.MethodPost:
			categoryHandler.CreateCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only handle if there's an ID in the path
		if strings.TrimPrefix(r.URL.Path, "/categories/") == "" {
			http.Error(w, "Category ID required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetCategoryByID(w, r)
		case http.MethodPut:
			categoryHandler.UpdateCategory(w, r)
		case http.MethodDelete:
			categoryHandler.DeleteCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Product routes
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case http.MethodGet:
			productHandler.GetAllProducts(w, r)
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only handle if there's an ID in the path
		if strings.TrimPrefix(r.URL.Path, "/products/") == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductByID(w, r)
		case http.MethodPut:
			productHandler.UpdateProduct(w, r)
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Swagger documentation endpoints
	mux.HandleFunc("/docs", categoryHandler.SwaggerHandler)
	mux.HandleFunc("/swagger.json", categoryHandler.SwaggerJSONHandler)

	return mux
}
