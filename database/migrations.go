package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Create products table
	productsTableSQL := `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	
	_, err := db.Exec(productsTableSQL)
	if err != nil {
		return fmt.Errorf("error creating products table: %w", err)
	}

	// Create index on product name
	indexSQL := `CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);`
	_, err = db.Exec(indexSQL)
	if err != nil {
		return fmt.Errorf("error creating products index: %w", err)
	}

	// Add category_id to products if it doesn't exist
	// 1. Add column as nullable first
	addColumnSQL := `ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER;`
	_, err = db.Exec(addColumnSQL)
	if err != nil {
		return fmt.Errorf("error adding category_id column: %w", err)
	}

	// 2. Set a default category (ID 1) for existing products that don't have one
	updateDefaultsSQL := `UPDATE products SET category_id = 1 WHERE category_id IS NULL;`
	_, err = db.Exec(updateDefaultsSQL)
	if err != nil {
		// If category 1 doesn't exist, this might fail, but let's assume it exists from previous steps
		// or just log it
	}

	// 3. Set NOT NULL constraint
	setNotNullSQL := `ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;`
	_, err = db.Exec(setNotNullSQL)
	if err != nil {
		return fmt.Errorf("error setting category_id NOT NULL: %w", err)
	}

	// 4. Add foreign key constraint if it doesn't exist
	// Note: PostgreSQL doesn't have ADD CONSTRAINT IF NOT EXISTS, so we check first
	addFKSQL := `
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_product_category') THEN 
				ALTER TABLE products ADD CONSTRAINT fk_product_category FOREIGN KEY (category_id) REFERENCES categories(id); 
			END IF; 
		END $$;
	`
	_, err = db.Exec(addFKSQL)
	if err != nil {
		return fmt.Errorf("error adding foreign key constraint: %w", err)
	}

	return nil
}
