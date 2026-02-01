package repository

import (
	"database/sql"
	"fmt"
	"tugas-session-1/models"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll retrieves all products with their category names
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name, p.created_at, p.updated_at 
		FROM products p
		JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.CategoryName,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByID retrieves a product by ID with its category name
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name as category_name, p.created_at, p.updated_at 
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	
	var product models.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.CategoryName,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Return nil without error to indicate not found
	}
	if err != nil {
		return nil, fmt.Errorf("error querying product: %w", err)
	}

	return &product, nil
}

// Create creates a new product
func (r *ProductRepository) Create(req *models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, price, stock, category_id, created_at, updated_at
	`
	
	var product models.Product
	err := r.db.QueryRow(query, req.Name, req.Price, req.Stock, req.CategoryID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(id int, req *models.UpdateProductRequest) (*models.Product, error) {
	query := `
		UPDATE products
		SET name = $1, price = $2, stock = $3, category_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, name, price, stock, category_id, created_at, updated_at
	`
	
	var product models.Product
	err := r.db.QueryRow(query, req.Name, req.Price, req.Stock, req.CategoryID, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Return nil without error to indicate not found
	}
	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return &product, nil
}

// Delete deletes a product by ID
func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
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
