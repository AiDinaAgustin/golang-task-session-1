package models

import "time"

// Product represents a product entity
type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID int     `json:"category_id"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID int     `json:"category_id"`
}

// Validate validates the create product request
func (r *CreateProductRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if r.Price < 0 {
		return &ValidationError{Field: "price", Message: "price must be greater than or equal to 0"}
	}
	if r.Stock < 0 {
		return &ValidationError{Field: "stock", Message: "stock must be greater than or equal to 0"}
	}
	if r.CategoryID <= 0 {
		return &ValidationError{Field: "category_id", Message: "category_id is required"}
	}
	return nil
}

// Validate validates the update product request
func (r *UpdateProductRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if r.Price < 0 {
		return &ValidationError{Field: "price", Message: "price must be greater than or equal to 0"}
	}
	if r.Stock < 0 {
		return &ValidationError{Field: "stock", Message: "stock must be greater than or equal to 0"}
	}
	if r.CategoryID <= 0 {
		return &ValidationError{Field: "category_id", Message: "category_id is required"}
	}
	return nil
}
