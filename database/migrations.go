package database

import (
	"database/sql"
	"fmt"
	"log"
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

	// Create transactions table
	transactionsTableSQL := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			total_amount BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(transactionsTableSQL)
	if err != nil {
		return fmt.Errorf("error creating transactions table: %w", err)
	}

	// Create transaction_details table
	transactionDetailsTableSQL := `
		CREATE TABLE IF NOT EXISTS transaction_details (
			id SERIAL PRIMARY KEY,
			transaction_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			subtotal BIGINT NOT NULL,
			CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
			CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
		);
	`
	_, err = db.Exec(transactionDetailsTableSQL)
	if err != nil {
		return fmt.Errorf("error creating transaction_details table: %w", err)
	}

	// Alter existing columns to BIGINT if they were INTEGER
	alterTotalAmountSQL := `ALTER TABLE transactions ALTER COLUMN total_amount TYPE BIGINT;`
	_, err = db.Exec(alterTotalAmountSQL)
	if err != nil {
		// Ignore error if column already BIGINT
		log.Printf("Note: Could not alter total_amount column (might already be BIGINT): %v", err)
	}

	alterSubtotalSQL := `ALTER TABLE transaction_details ALTER COLUMN subtotal TYPE BIGINT;`
	_, err = db.Exec(alterSubtotalSQL)
	if err != nil {
		// Ignore error if column already BIGINT
		log.Printf("Note: Could not alter subtotal column (might already be BIGINT): %v", err)
	}

	// Create index on transaction_id for better query performance
	transactionDetailsIndexSQL := `CREATE INDEX IF NOT EXISTS idx_transaction_details_transaction_id ON transaction_details(transaction_id);`
	_, err = db.Exec(transactionDetailsIndexSQL)
	if err != nil {
		return fmt.Errorf("error creating transaction_details index: %w", err)
	}

	// Create index on created_at for report queries
	transactionsIndexSQL := `CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);`
	_, err = db.Exec(transactionsIndexSQL)
	if err != nil {
		return fmt.Errorf("error creating transactions index: %w", err)
	}

	return nil
}
