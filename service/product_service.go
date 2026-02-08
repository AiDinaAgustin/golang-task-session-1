package service

import (
	"tugas-session-1/models"
	"tugas-session-1/repository"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	GetAll(searchName string) ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(req *models.CreateProductRequest) (*models.Product, error)
	Update(id int, req *models.UpdateProductRequest) (*models.Product, error)
	Delete(id int) error
}

// productService implements ProductService interface
type productService struct {
	repo *repository.ProductRepository
}

// NewProductService creates a new product service with dependency injection
func NewProductService(repo *repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// GetAll retrieves all products, optionally filtered by name
func (s *productService) GetAll(searchName string) ([]models.Product, error) {
	return s.repo.GetAll(searchName)
}

// GetByID retrieves a product by ID
func (s *productService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// Create creates a new product with business logic
func (s *productService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// Business logic could be added here
	// For example: check if product name already exists, etc.
	
	return s.repo.Create(req)
}

// Update updates an existing product with business logic
func (s *productService) Update(id int, req *models.UpdateProductRequest) (*models.Product, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	// Business logic could be added here
	// For example: validate stock levels, price changes, etc.
	
	return s.repo.Update(id, req)
}

// Delete deletes a product
func (s *productService) Delete(id int) error {
	// Business logic could be added here
	// For example: check if product is referenced in orders, etc.
	
	return s.repo.Delete(id)
}
