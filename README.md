# ShubServer (Webserver)

A Go backend service providing authentication, notes, and password-manager APIs, backed by PostgreSQL (Supabase-ready) with migrations and Docker support.

---

## ✨ Features

- **Auth**: Register / Login with JWT
- **Notes**: CRUD with user scoping
- **Password Manager**: Encrypted blob storage (AES-GCM), never decrypts on server
- **Health Checks**: Lightweight + detailed endpoints
- **Migrations**: Versioned SQL migrations
- **Docker-ready**: Multi-stage build + compose support

---

## 🧱 Tech Stack

- **Go** (chi router)
- **PostgreSQL** (Supabase)
- **Migrations**: `migrate/migrate`
- **Docker** + Docker Compose

---

## 📁 Project Structure

```
cmd/
  server/
    main.go
internal/
  app/
    routes.go
    server.go
  auth/
    handlers.go
    jwt.go
    routes.go
    service.go
  config/
    config.go
  db/
    db.go
  health/
    handler.go
    routes.go
  middleware/
    auth.go
  notes/
    handlers.go
    model.go
    repository.go
    routes.go
  password-manager/
    handlers.go
    model.go
    repository.go
    routes.go
    service.go
  users/
    model.go
    model_db.go
    postgres.go
    repository.go
  utils/
    errors.go
    response.go
migrations/
  001_create_extensions.*.sql
  002_create_users_table.*.sql
  003_create_notes_table.*.sql
  004_create_passwords_table.*.sql
```

---

## 🚀 Quick Start

### 1) Environment Variables

Create a `.env` file:

```env
APP_ENV=development
APP_VERSION=dev
JWT_SECRET_KEY=your-secret-key-here
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
```

---

### 2) Run with Docker (Recommended)

```bash
docker compose up --build
```

---

### 3) Run with Local Go

```bash
go mod download
go run ./cmd/server
```

---

## 🗃️ Database Migrations

Run migrations:

```bash
docker compose run migrate
```

---

## 🔌 Health Endpoints

- `GET /health` → Lightweight
- `GET /health/status` → Status summary
- `GET /health/detailed` → Full checks

---

## 🔐 Security Model (Password Manager)

- **Server never sees secrets**
- Client encrypts data using **AES-256-GCM**
- Keys derived with **Argon2id**
- Server stores only ciphertext + nonce

---

## 🧪 API Overview

### Auth
- `POST /auth/register`
- `POST /auth/login`

### Notes
- `GET /notes`
- `GET /notes/{id}`
- `POST /notes`
- `PUT /notes/{id}`
- `DELETE /notes/{id}`

### Password Manager
- `GET /passwords`
- `GET /passwords/{id}`
- `POST /passwords`
- `PUT /passwords/{id}`
- `DELETE /passwords/{id}`

---

## ⚙️ Docker Setup

### Dockerfile
- Multi-stage build
- Minimal Alpine runtime

### Compose
- `docker-compose.yml`: app + migrations
- `docker_compose.dev.yml`: local dev DB

---

## ✅ Contributing

1. Fork the repo
2. Create a branch
3. Commit changes
4. Open a PR

---

## 📄 License

MIT (add your license here)
