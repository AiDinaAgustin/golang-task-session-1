package service

import (
	"tugas-session-1/models"
	"tugas-session-1/repository"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetAll() ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(req *models.CreateCategoryRequest) (*models.Category, error)
	Update(id int, req *models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id int) error
}

// categoryService implements CategoryService interface
type categoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService creates a new category service with dependency injection
func NewCategoryService(repo *repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

// GetAll retrieves all categories
func (s *categoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a category by ID
func (s *categoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// Create creates a new category with business logic
func (s *categoryService) Create(req *models.CreateCategoryRequest) (*models.Category, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	return s.repo.Create(req)
}

// Update updates an existing category with business logic
func (s *categoryService) Update(id int, req *models.UpdateCategoryRequest) (*models.Category, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}
	
	return s.repo.Update(id, req)
}

// Delete deletes a category
func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
