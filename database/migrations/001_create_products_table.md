# Database Migration - Products Table

## SQL Script untuk Neon PostgreSQL

Jalankan SQL berikut di Neon SQL Editor:

```sql
-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for product name for faster searches
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
```

## Cara Menjalankan

1. Buka https://console.neon.tech
2. Pilih project Anda
3. Klik **SQL Editor**
4. Copy-paste SQL di atas
5. Klik **Run**

## Verifikasi

Setelah menjalankan SQL, verifikasi dengan:

```sql
SELECT * FROM products;
```

Table seharusnya sudah ada dan kosong (belum ada data).
