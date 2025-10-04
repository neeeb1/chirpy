# Chirpy 🐦
[![Go](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org/)  
[![Postgres](https://img.shields.io/badge/Postgres-15+-blueviolet.svg)](https://www.postgresql.org/)  

A **minimal Twitter-style backend** built with Go.  
Chirpy was built as a Boot.Dev project to demonstrate REST APIs, JSON Web Token (JWT) authentication, refresh tokens, and Postgres persistence.  

---

## ✨ Features

- User registration, login, and updates
- JWT + refresh tokens for authentication
- CRUD for “Chirps” (tweets)
- Sorting and filtering of Chrips
- Mock payment processor webhooks ("Polka") for premium account upgrades (“Chirpy Red”)
- PostgreSQL + `sqlc` for DB access
- A small number of unit tests

---

## 📦 Tech Stack

- **Go** (≥ 1.25)  
- **PostgreSQL**  
- **sqlc** for generating DB queries 
    - *https://github.com/sqlc-dev/sqlc*
- **JWT** for authorization
    - *https://github.com/golang-jwt/jwt*
- **goose** for database migrations
    - *https://github.com/pressly/goose*

---

## ⚙️ Setup & Installation

```bash
# Clone the repo
git clone https://github.com/neeeb1/chirpy.git
cd chirpy

# Download dependencies
go mod download
```



Create a `.env` file in the project root:

```env
JWT_SECRET=your_jwt_secret_here
POLKA_API_KEY="f271c81ff7084ee5b99a5091b42d486e"
POSTGRES_URL="postgres://user:password@localhost:5432/chirpy_db?sslmode=disable"
```

Polka doesn't exist, so we're using the mock API Key provided by Boot.Dev for this project.

To generate a JWT_SECRET, I recommend creating a random base64 string with openssl.

```bash
openssl rand -base64 64
```

> ⚠️ Make sure `.env` is in `.gitignore` so your secrets aren't committed!


---

## 🚀 Running the Server

```bash
go run .
```

Or build first:

```bash
go build -o chirpy
./chirpy
```

The default port is **:8080**

---

## 📖 API Reference

### 🔹 Users

**Register a new user**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

**Update a user's email or password**
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com","password":"newpass"}'
```

---

### 🔹 Auth

**Login**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

**Refresh access token**
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <refresh_token>"
```

**Revoke a refresh token**
```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer <refresh_token>"
```

---

### 🔹 Chirps

**Post a new Chirp**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello Chirpy!"}'
```

**List all Chirps**

Sort and author_id are both optional parameters. If not provided, will default  to listing all chirps in the database in ascending order by date created.

```bash
curl "http://localhost:8080/api/chirps?sort=desc&author_id=<uuid>"
```

**Get a specific Chirp by it's ID**
```bash
curl http://localhost:8080/api/chirps/<id>
```

**Delete a Chirp**
```bash
curl -X DELETE http://localhost:8080/api/chirps/<id> \
  -H "Authorization: Bearer <token>"
```

---

### 🔹 'Polka' payment processor webhook

**Upgrade User to Chirpy Red**  
```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Authorization: ApiKey <POLKA_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"event":"user.upgraded","data":{"user_id":"<uuid>"}}'
```

---

## 🛠️ Database

- PostgreSQL stores users, Chirps, tokens  
- `sqlc` generates Go code for sql queries
- Use `goose` migrations to manage schema changes

---

## 🧪 Testing

Run unit tests:

```bash
go test ./...
```