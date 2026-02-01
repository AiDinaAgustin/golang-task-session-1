package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tugas-session-1/models"
	"tugas-session-1/service"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new product handler with dependency injection
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetAllProducts handles GET /products
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	// Return empty array instead of null
	if products == nil {
		products = []models.Product{}
	}

	respondJSON(w, http.StatusOK, products)
}

// GetProductByID handles GET /products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if product == nil {
		respondError(w, http.StatusNotFound, "Not Found", "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	product, err := h.service.Create(&req)
	if err != nil {
		// Check if it's a validation error
		if _, ok := err.(*models.ValidationError); ok {
			respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

// UpdateProduct handles PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid product ID")
		return
	}

	var req models.UpdateProductRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	product, err := h.service.Update(id, &req)
	if err != nil {
		// Check if it's a validation error
		if _, ok := err.(*models.ValidationError); ok {
			respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if product == nil {
		respondError(w, http.StatusNotFound, "Not Found", "Product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// DeleteProduct handles DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid product ID")
		return
	}

	err = h.service.Delete(id)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Not Found", "Product not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Product deleted successfully",
	})
}
