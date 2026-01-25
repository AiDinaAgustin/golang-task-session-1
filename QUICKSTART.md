# Quick Start Guide - Go CRUD API

## 🚀 Langkah Cepat

### 1. Setup Database di Neon

1. Buka https://console.neon.tech
2. Pilih project Anda
3. Klik **SQL Editor**
4. Copy-paste SQL berikut:

```sql
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
```

5. Klik **Run**

### 2. Configure Environment

Edit file `.env` dan isi DATABASE_URL Anda:

```env
DATABASE_URL=postgres://user:password@host/database?sslmode=require
PORT=8080
```

### 3. Run Server

```powershell
go run main.go
```

### 4. Test API

```powershell
# Health check
Invoke-WebRequest -Uri http://localhost:8080/health

# Create category
Invoke-WebRequest -Uri http://localhost:8080/categories -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name":"Electronics","description":"Gadgets"}'

# Get all
Invoke-WebRequest -Uri http://localhost:8080/categories | Select-Object -ExpandProperty Content
```

## ✅ Done!

API CRUD Anda sudah jalan! 🎉

Dokumentasi lengkap: [README.md](file:///c:/project-frontend/project-golang/tugas-session-1/README.md)
