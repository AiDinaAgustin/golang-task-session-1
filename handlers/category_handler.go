package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tugas-session-1/models"
	"tugas-session-1/repository"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
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

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error JSON response
func respondError(w http.ResponseWriter, status int, err string, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   err,
		Message: message,
	})
}

// GetAllCategories handles GET /categories
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	// Return empty array instead of null
	if categories == nil {
		categories = []models.Category{}
	}

	respondJSON(w, http.StatusOK, categories)
}

// GetCategoryByID handles GET /categories/{id}
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	category, err := h.repo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if category == nil {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	category, err := h.repo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, category)
}

// UpdateCategory handles PUT /categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	var req models.UpdateCategoryRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	category, err := h.repo.Update(id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if category == nil {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}

	respondJSON(w, http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid category ID")
		return
	}

	err = h.repo.Delete(id)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Not Found", "Category not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Category deleted successfully",
	})
}

// extractIDFromPath extracts the ID from the URL path
// Expected path format: /categories/{id}
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

// SwaggerHandler serves the Swagger UI
func (h *CategoryHandler) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Category CRUD API - Swagger Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
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
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// SwaggerJSONHandler serves the swagger.json file
func (h *CategoryHandler) SwaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "swagger.json")
}
