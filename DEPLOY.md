# Quick Deploy - Vercel

## 🚀 Deploy dalam 3 Langkah

### 1. Push ke GitHub
```bash
git add .
git commit -m "Ready for Vercel"
git push origin main
```

### 2. Import di Vercel

1. Buka https://vercel.com
2. **New Project** → Import dari GitHub
3. Pilih repository `golang-task-session-1`

### 3. Configure & Deploy

**Environment Variables:**
- `DATABASE_URL` = your Neon connection string

**Klik Deploy** → Done! 🎉

## 🔗 Setelah Deploy

API akan live di: `https://your-project.vercel.app`

**Test:**
```bash
curl https://your-project.vercel.app/health
curl https://your-project.vercel.app/categories
```

**Swagger Docs:**
```
https://your-project.vercel.app/docs
```

## 📚 Panduan Lengkap

Lihat [`VERCEL.md`](./VERCEL.md) untuk:
- Deploy via CLI
- Troubleshooting
- Custom domain
- Monitoring

---

**Auto Deploy:** Setiap push ke `main` otomatis deploy! ✨
