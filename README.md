# Category CRUD API - Go

RESTful API untuk manajemen kategori menggunakan Go (vanilla, tanpa framework) dan PostgreSQL (Neon).

## 📋 Fitur

- ✅ CRUD lengkap untuk kategori
- ✅ Validasi input
- ✅ Error handling yang proper
- ✅ CORS enabled
- ✅ Health check endpoint

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL (Neon)
- **Dependencies**: 
  - `lib/pq` - PostgreSQL driver
  - `godotenv` - Environment variable loader

## 📦 Struktur Project

```
tugas-session-1/
├── database/
│   ├── db.go           # Database connection
│   └── schema.sql      # Database schema
├── handlers/
│   └── category_handler.go  # HTTP handlers
├── models/
│   └── category.go     # Data models
├── repository/
│   └── category_repository.go  # Data access layer
├── router/
│   └── router.go       # Route configuration
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── .env.example        # Environment template
└── README.md           # This file
```

## 🚀 Setup

### 1. Clone atau navigate ke project directory

```bash
cd c:\project-frontend\project-golang\tugas-session-1
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup database

Jalankan SQL schema di Neon database Anda secara manual:

**Opsi 1: Neon Console (Paling Mudah)**
1. Buka https://console.neon.tech
2. Pilih project dan database Anda
3. Klik **SQL Editor**
4. Copy-paste isi dari `database/schema.sql` dan klik Run

**Opsi 2: psql (Command Line)**
```bash
psql "your_neon_connection_string" -f database/schema.sql
```


### 4. Configure environment variables

Copy `.env.example` ke `.env` dan isi dengan connection string Neon Anda:

```bash
cp .env.example .env
```

Edit `.env`:

```env
DATABASE_URL=postgres://username:password@your-neon-host/dbname?sslmode=require
PORT=8080
```

### 5. Run the application

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

### 1. Get All Categories

```bash
curl http://localhost:8080/categories
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": "2024-01-25T10:00:00Z",
    "updated_at": "2024-01-25T10:00:00Z"
  }
]
```

### 2. Get Category by ID

```bash
curl http://localhost:8080/categories/1
```

**Response:**
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices and gadgets",
  "created_at": "2024-01-25T10:00:00Z",
  "updated_at": "2024-01-25T10:00:00Z"
}
```

### 3. Create Category

```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Books\",\"description\":\"Books and publications\"}"
```

**Response:** `201 Created`
```json
{
  "id": 2,
  "name": "Books",
  "description": "Books and publications",
  "created_at": "2024-01-25T10:05:00Z",
  "updated_at": "2024-01-25T10:05:00Z"
}
```

### 4. Update Category

```bash
curl -X PUT http://localhost:8080/categories/1 \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Updated Electronics\",\"description\":\"Updated description\"}"
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Updated Electronics",
  "description": "Updated description",
  "created_at": "2024-01-25T10:00:00Z",
  "updated_at": "2024-01-25T10:10:00Z"
}
```

### 5. Delete Category

```bash
curl -X DELETE http://localhost:8080/categories/1
```

**Response:** `200 OK`
```json
{
  "message": "Category deleted successfully"
}
```

### 6. Health Check

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

## ⚠️ Error Responses

### 400 Bad Request
```json
{
  "error": "Bad Request",
  "message": "name is required"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Category not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "error details..."
}
```

## 🧪 Testing

Gunakan Postman, Insomnia, atau curl untuk testing API.

**Contoh testing dengan PowerShell:**

```powershell
# Create
Invoke-WebRequest -Uri http://localhost:8080/categories -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"name":"Test","description":"Test category"}'

# Get All
Invoke-WebRequest -Uri http://localhost:8080/categories

# Get by ID
Invoke-WebRequest -Uri http://localhost:8080/categories/1

# Update
Invoke-WebRequest -Uri http://localhost:8080/categories/1 -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"name":"Updated","description":"Updated desc"}'

# Delete
Invoke-WebRequest -Uri http://localhost:8080/categories/1 -Method DELETE
```

## 📝 Model Category

```go
type Category struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## 🔒 Validasi

- `name` wajib diisi (required)
- `description` optional

## 📄 License

MIT
