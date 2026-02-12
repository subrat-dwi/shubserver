# Auth Module

User authentication and JWT token management for ShubServer.

---

## 📋 Overview

The `auth` module handles:
- User registration with bcrypt password hashing
- User login with credential validation
- JWT token generation and verification
- Token claims management

---

## 🏗️ Architecture

### Components

| File | Responsibility |
|------|-----------------|
| `service.go` | Business logic (register, login) |
| `handlers.go` | HTTP request/response handling |
| `jwt.go` | JWT generation and verification |
| `routes.go` | Route definitions |

### Flow Diagram
```
Client Request
↓
handlers.go (HTTP layer)
↓
service.go (Business logic)
↓
users repository (Database)
↓
jwt.go (Token generation)
↓
Response with JWT Token
```
---

## 🔐 Security

### Password Handling
- **Never stored in plaintext**
- Hashed with **bcrypt** (default cost: 10)
- Compared securely during login

### JWT Tokens
- **Algorithm**: HMAC with SHA-256
- **Secret**: Loaded from `JWT_SECRET_KEY` environment variable
- **Claims**: User ID + standard JWT claims

### Environment Variables Required

```env
JWT_SECRET_KEY=your-secret-key-here  # Keep this secure!
```
### 📡 API Endpoints
Register User

```bash
POST /auth/register
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "secure-password"
}
```
Response (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Login User

```bash
POST /auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "secure-password"
}
```

Response (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 🔑 JWT Token Structure
Token Claims

```go
type Claims struct {
    UserID string                  // User's UUID
    RegisteredClaims jwt.RegisteredClaims
}
```
Registered Claims

iat (Issued At): Timestamp when token was created

exp (Expires At): Timestamp when token expires (now + 12 hours)

```go
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "iat": 1707734400,
  "exp": 1707778400
}
```

#### 🔄 Dependencies
 
| Dependency | Purpose |
|---|---|
| `github.com/golang-jwt/jwt/v5` | JWT token handling |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `users/` | User repository |
| `utils/` | Response helpers |


#### 🧪 Testing
Register Test
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Login Test
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 🔒 Security Checklist
    ✅ Passwords hashed with bcrypt
    ✅ JWT secret from environment variable
    ✅ Token expiry set (12 hours)
    ✅ HMAC-SHA256 signing
    ✅ Secure credential comparison
    ✅ No plaintext password logging

### 🚧 Future Enhancements
    Refresh token mechanism
    Rate limiting on login attempts
    Password reset / forgot password
    Email verification on registration