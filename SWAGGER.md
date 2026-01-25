# Swagger API Documentation

Dokumentasi API tersedia melalui Swagger UI yang **terintegrasi langsung** dengan aplikasi.

## 🚀 Cara Akses

Jalankan server:
```bash
go run main.go
```

Lalu buka browser dan akses:
```
http://localhost:8080/docs
```

## ✨ Features

- ✅ UI interaktif dengan **Try it out** feature
- ✅ Contoh request/response untuk setiap endpoint
- ✅ Parameter dan schema lengkap
- ✅ Bisa langsung test API dari browser
- ✅ Download OpenAPI spec di `/swagger.json`

## 📡 Endpoints

- `GET /docs` - Swagger UI (buka di browser)
- `GET /swagger.json` - OpenAPI JSON specification

## 📥 Import ke Tools Lain

### Postman
1. Buka Postman
2. Import → Link → `http://localhost:8080/swagger.json`

### Insomnia
1. Buka Insomnia
2. Import → URL → `http://localhost:8080/swagger.json`

## 📝 Dokumentasi Mencakup

- ✅ Semua 5 endpoint CRUD categories
- ✅ Health check endpoint
- ✅ Request body schemas
- ✅ Response examples (success & error)
- ✅ HTTP status codes
- ✅ Validasi requirements
