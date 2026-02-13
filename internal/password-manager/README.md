# Password Manager Module

Encrypted password storage service with AES-256-GCM encryption. The server **never decrypts** passwords—it only stores encrypted blobs.

---

## 📋 Overview

The `password-manager` module provides:
- Secure password storage with client-side encryption
- AES-256-GCM encryption with 12-byte nonce
- Argon2id key derivation (client-side)
- CRUD operations for encrypted passwords
- Full-text search on encrypted metadata
- User-scoped password isolation

---

## 🏗️ Architecture

### Components

| File | Responsibility |
|------|-----------------|
| `model.go` | Password data structure |
| `service.go` | Business logic & validation |
| `repository.go` | Database operations |
| `handlers.go` | HTTP request/response handling |
| `routes.go` | Route definitions |

### Flow Diagram
    Client (Encrypts locally)
    ↓
    handlers.go (HTTP layer)
    ↓
    service.go (Validation)
    ↓
    repository.go (Database)
    ↓
    PostgreSQL (Stores ciphertext only)

---

## 🔐 Security Model

### Core Principle
> **The server never sees secrets.**

### Encryption Flow

#### Client Side (What happens on user's device)
    The client encrypts the password using a key derived from the user's password and stores the encrypted blob. The client also stores the nonce (12 bytes) and the salt (16 bytes) for future decryption.

#### Server Side (What happens on the server)
    The server receives the encrypted blob and stores it in the database. The server does **not** decrypt the password. It only stores the encrypted blob and the associated metadata.

    The key derivation function (Argon2id) is performed on the client side, ensuring that the key is derived from the user's password and not from any other source.

---

## 📦 Data Structure

### Password
```go
type Password struct {
    ID             uuid.UUID  // Unique identifier
    UserID         uuid.UUID  // User ownership
    Name           string     // Service name (e.g., "Gmail")
    Username       string     // Account identifier (e.g., "user@gmail.com")
    Ciphertext     []byte     // AES-256-GCM encrypted password
    Nonce          []byte     // 12-byte GCM nonce (never reuse!)
    EncryptVersion int        // For future algorithm migrations
    CreatedAt      time.Time  // Creation timestamp
    UpdatedAt      time.Time  // Last modification timestamp
}
```

### PasswordItem
```go
type PasswordItem struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Username  string `json:"username"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```
---

### 📡 API Endpoints
All endpoints require authentication (JWT token in Authorization: Bearer <token> header).

**Create Password**

```bash
POST /passwords
Content-Type: application/json
Authorization: Bearer <jwt_token>
{
  "name": "Gmail",
  "username": "user@gmail.com",
  "password": "base64_encoded_ciphertext",
  "nonce": "base64_encoded_12_bytes"
}
```

Response (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Gmail",
  "username": "user@gmail.com",
  "created_at": "2024-02-12T10:30:00Z",
  "updated_at": "2024-02-12T10:30:00Z"
}
```

**List Passwords**
```bash
GET /passwords
Authorization: Bearer <jwt_token>
```

Response (200 OK):
```json
{
  "passwords": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Gmail",
      "username": "user@gmail.com",
      "created_at": "2024-02-12T10:30:00Z",
      "updated_at": "2024-02-12T10:30:00Z"
    }
  ]
}
```

**Get Single Password**
```bash
GET /passwords/{id}
Authorization: Bearer <jwt_token>
```
Response (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Gmail",
  "username": "user@gmail.com",
  "password": "base64_encoded_ciphertext",
  "nonce": "base64_encoded_12_bytes",
  "created_at": "2024-02-12T10:30:00Z",
  "updated_at": "2024-02-12T10:30:00Z"
}
```

**Update Single Password**
```bash
PUT /passwords/{id}
Content-Type: application/json
Authorization: Bearer <jwt_token>
{
  "name": "Gmail Updated",
  "username": "user2@gmail.com",
  "password": "new_base64_ciphertext",
  "nonce": "new_base64_nonce_12_bytes"
}
```
Response (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Gmail Updated",
  "username": "user2@gmail.com",
  "created_at": "2024-02-12T10:30:00Z",
  "updated_at": "2024-02-12T10:30:00Z"
}
```

**Delete Single Password**
```bash
DELETE /passwords/{id}
Authorization: Bearer <jwt_token>
```

**Search Passwords**
```bash
GET /passwords/search?search=gmail
Authorization: Bearer <jwt_token>
```
Response (200 OK):
```json
{
  "passwords": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Gmail",
      "username": "user@gmail.com",
      "created_at": "2024-02-12T10:30:00Z",
      "updated_at": "2024-02-12T10:30:00Z"
    }
  ]
}
```

---

### 🔒 Security Checklist
    ✅ AES-256-GCM encryption (authenticated)
    ✅ 12-byte nonce (GCM standard)
    ✅ User-scoped access (can't access other users' passwords)
    ✅ Database constraints enforced
    ✅ No plaintext storage
    ✅ Automatic timestamp updates
    ✅ UUID-based IDs (unpredictable)

### 🔄 Dependencies
| Dependency | Purpose |
|---|---|
| `github.com/google/uuid` | UUID generation and parsing |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/go-chi/chi/v5` | HTTP routing |
| `utils` | Response helpers |
| `middleware` | Authentication middleware |


### 🚀 Example Usage

```bash
# Create a password
curl -X POST http://localhost:8080/passwords \
    -H "Authorization: Bearer <jwt_token>" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Gmail",
        "username": "user@gmail.com",
        "password": "base64_ciphertext",
        "nonce": "base64_nonce_12bytes"
    }'
```

```bash
# List all passwords
curl -H "Authorization: Bearer <jwt_token>" \
    http://localhost:8080/passwords
```

```bash
# Get a single password
curl -H "Authorization: Bearer <jwt_token>" \
    http://localhost:8080/passwords/{id}
```

```bash
# Update a password
curl -X PUT http://localhost:8080/passwords/{id} \
    -H "Authorization: Bearer <jwt_token>" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Gmail Updated",
        "username": "user2@gmail.com",
        "password": "new_ciphertext",
        "nonce": "new_nonce"
    }'
```

```bash
# Delete a password
curl -X DELETE http://localhost:8080/passwords/{id} \
    -H "Authorization: Bearer <jwt_token>"
```

---
