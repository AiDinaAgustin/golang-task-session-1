# Deploy ke Vercel

Panduan deploy Go CRUD API ke Vercel menggunakan serverless functions.

## 📦 Persiapan

Struktur file yang dibutuhkan:
```
tugas-session-1/
├── api/
│   └── index.go          # Serverless function handler
├── vercel.json            # Konfigurasi Vercel
├── swagger.json           # API documentation
└── ... (file lainnya)
```

✅ **Sudah disediakan!**

## 🚀 Cara Deploy

### Option 1: Via Vercel Dashboard (Paling Mudah)

1. **Buka** https://vercel.com
2. **Sign up/Login** dengan GitHub
3. **Import Project** → Pilih repository `golang-task-session-1`
4. **Configure:**
   - Framework: `Other`
   - Root Directory: `./`
5. **Environment Variables:**
   - Klik **Add** → `DATABASE_URL` = your Neon connection string
6. **Deploy!**

### Option 2: Via Vercel CLI

```bash
# Install Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
cd c:\project-frontend\project-golang\tugas-session-1
vercel

# Follow prompts:
# - Setup and deploy: Y
# - Scope: pilih akun Anda
# - Link to existing project: N
# - Project name: go-crud-api (atau terserah)
# - Directory: ./ (enter)
# - Override settings: N
```

### Set Environment Variables (CLI)

```bash
vercel env add DATABASE_URL
# Paste your Neon connection string

# Deploy production
vercel --prod
```

## 🔗 Setelah Deploy

API akan tersedia di:
```
https://your-project.vercel.app
```

### Test Endpoints:

```bash
# Health check
curl https://your-project.vercel.app/health

# Get categories
curl https://your-project.vercel.app/categories

# Swagger docs (buka di browser)
https://your-project.vercel.app/docs
```

## 📝 Catatan Penting

### 1. Database Connection

Pastikan Neon connection string sudah benar di environment variables Vercel:
```
postgres://user:password@host/database?sslmode=require
```

### 2. Path Prefix

Semua endpoint bisa diakses dengan atau tanpa prefix `/api`:
- ✅ `/categories` 
- ✅ `/api/categories`
- ✅ `/health`
- ✅ `/api/health`

### 3. Cold Start

Serverless functions di Vercel punya "cold start" - request pertama bisa sedikit lambat. Ini normal untuk serverless architecture.

### 4. Connection Pooling

Code sudah optimized untuk Vercel dengan:
```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
```

## 🔄 Auto Deploy

Setiap push ke GitHub branch `main` akan otomatis deploy ke Vercel (kalau disetup via dashboard).

## 🐛 Troubleshooting

### Build Error

Pastikan `go.mod` dan `go.sum` sudah di-commit:
```bash
git add go.mod go.sum
git commit -m "Add go modules"
git push
```

### Database Connection Error

Cek environment variable `DATABASE_URL` di Vercel dashboard:
- Settings → Environment Variables
- Pastikan valuenya benar dan include `?sslmode=require`

### 500 Internal Server Error

Cek logs di Vercel dashboard:
- Deployments → klik deployment terakhir → View Function Logs

## 📊 Monitoring

Vercel dashboard menyediakan:
- ✅ Function logs
- ✅ Analytics
- ✅ Performance metrics

Access di: https://vercel.com/dashboard

## ✨ Custom Domain

1. Vercel Dashboard → Project → Settings → Domains
2. Add domain Anda
3. Update DNS sesuai instruksi
4. Done! SSL otomatis

## 🎉 Selesai!

API Anda sudah live di Vercel dengan:
- ✅ Auto-scaling
- ✅ Global CDN
- ✅ Free SSL
- ✅ Auto deploy from GitHub
- ✅ Swagger documentation

Enjoy! 🚀
