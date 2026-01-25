package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
	db   *sql.DB
)

// Category represents a category entity
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func initApp() {
	// Initialize database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Health check
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/health", healthHandler)

	// Categories endpoints
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/api/categories", categoriesHandler)
	http.HandleFunc("/categories/", categoryByIDHandler)
	http.HandleFunc("/api/categories/", categoryByIDHandler)

	// Swagger documentation
	http.HandleFunc("/docs", swaggerUIHandler)
	http.HandleFunc("/api/docs", swaggerUIHandler)
	http.HandleFunc("/swagger.json", swaggerJSONHandler)
	http.HandleFunc("/api/swagger.json", swaggerJSONHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllCategories(w, r)
	case http.MethodPost:
		createCategory(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method Not Allowed", "")
	}
}

func categoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getCategoryByID(w, r)
	case http.MethodPut:
		updateCategory(w, r)
	case http.MethodDelete:
		deleteCategory(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method Not Allowed", "")
	}
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, description, created_at, updated_at FROM categories ORDER BY id")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
			return
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []Category{}
	}

	respondJSON(w, http.StatusOK, categories)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	var cat Category
	err = db.QueryRow("SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1", id).
		Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, cat)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Bad Request", "name is required")
		return
	}

	var cat Category
	err := db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at",
		req.Name, req.Description,
	).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, cat)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Bad Request", "name is required")
		return
	}

	var cat Category
	err = db.QueryRow(
		"UPDATE categories SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, description, created_at, updated_at",
		req.Name, req.Description, id,
	).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, cat)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	result, err := db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Category deleted successfully",
	})
}

func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Category CRUD API - Swagger Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                plugins: [SwaggerUIBundle.plugins.DownloadUrl],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Embed swagger spec directly
	spec := `{
  "openapi": "3.0.3",
  "info": {
    "title": "Category CRUD API",
    "description": "RESTful API untuk manajemen kategori menggunakan Go vanilla dengan PostgreSQL (Neon)",
    "version": "1.0.0"
  },
  "servers": [{"url": "/", "description": "Current server"}],
  "tags": [
    {"name": "Categories", "description": "Operasi CRUD untuk kategori"},
    {"name": "Health", "description": "Health check endpoint"}
  ],
  "paths": {
    "/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "API is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string", "example": "healthy"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/categories": {
      "get": {
        "tags": ["Categories"],
        "summary": "Get all categories",
        "responses": {
          "200": {
            "description": "Successful operation",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"$ref": "#/components/schemas/Category"}
                }
              }
            }
          }
        }
      },
      "post": {
        "tags": ["Categories"],
        "summary": "Create a new category",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateCategoryRequest"},
              "example": {"name": "Electronics", "description": "Electronic devices"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Category created",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Category"}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Error"}
              }
            }
          }
        }
      }
    },
    "/categories/{id}": {
      "get": {
        "tags": ["Categories"],
        "summary": "Get category by ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "integer", "example": 1}
          }
        ],
        "responses": {
          "200": {
            "description": "Successful operation",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Category"}
              }
            }
          },
          "404": {
            "description": "Category not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Error"}
              }
            }
          }
        }
      },
      "put": {
        "tags": ["Categories"],
        "summary": "Update category",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "integer", "example": 1}
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/UpdateCategoryRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Category updated",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Category"}
              }
            }
          },
          "404": {
            "description": "Category not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Error"}
              }
            }
          }
        }
      },
      "delete": {
        "tags": ["Categories"],
        "summary": "Delete category",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "schema": {"type": "integer", "example": 1}
          }
        ],
        "responses": {
          "200": {
            "description": "Category deleted",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "message": {"type": "string", "example": "Category deleted successfully"}
                  }
                }
              }
            }
          },
          "404": {
            "description": "Category not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/Error"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Category": {
        "type": "object",
        "properties": {
          "id": {"type": "integer", "example": 1},
          "name": {"type": "string", "example": "Electronics"},
          "description": {"type": "string", "example": "Electronic devices"},
          "created_at": {"type": "string", "format": "date-time"},
          "updated_at": {"type": "string", "format": "date-time"}
        }
      },
      "CreateCategoryRequest": {
        "type": "object",
        "required": ["name"],
        "properties": {
          "name": {"type": "string", "example": "Electronics"},
          "description": {"type": "string", "example": "Electronic devices"}
        }
      },
      "UpdateCategoryRequest": {
        "type": "object",
        "required": ["name"],
        "properties": {
          "name": {"type": "string", "example": "Updated Electronics"},
          "description": {"type": "string", "example": "Updated description"}
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "error": {"type": "string", "example": "Bad Request"},
          "message": {"type": "string", "example": "name is required"}
        }
      }
    }
  }
}`
	
	w.Write([]byte(spec))
}

// Helper functions
func extractIDFromPath(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, sql.ErrNoRows
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, err string, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   err,
		Message: message,
	})
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	http.DefaultServeMux.ServeHTTP(w, r)
}
