-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name for faster queries
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Sample data (optional)
-- INSERT INTO categories (name, description) VALUES
-- ('Electronics', 'Electronic devices and gadgets'),
-- ('Books', 'Books and publications'),
-- ('Clothing', 'Apparel and fashion items');
