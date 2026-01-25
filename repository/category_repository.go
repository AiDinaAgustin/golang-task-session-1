package repository

import (
	"database/sql"
	"fmt"
	"tugas-session-1/models"
)

type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY id`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`
	
	var category models.Category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Return nil without error to indicate not found
	}
	if err != nil {
		return nil, fmt.Errorf("error querying category: %w", err)
	}

	return &category, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(req *models.CreateCategoryRequest) (*models.Category, error) {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	
	var category models.Category
	err := r.db.QueryRow(query, req.Name, req.Description).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error creating category: %w", err)
	}

	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(id int, req *models.UpdateCategoryRequest) (*models.Category, error) {
	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	
	var category models.Category
	err := r.db.QueryRow(query, req.Name, req.Description, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Return nil without error to indicate not found
	}
	if err != nil {
		return nil, fmt.Errorf("error updating category: %w", err)
	}

	return &category, nil
}

// Delete deletes a category by ID
func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Indicate that no rows were deleted
	}

	return nil
}
