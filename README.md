# dj_shorturl-go

A URL shortener with custom authentication token system, ported from Django to Go.

## Features

- Custom authentication token (HMAC-SHA256 + blake2b, educational)
- URL shortening with UUID-based identifiers
- SQLite (local dev) or PostgreSQL (production)
- CGO-free SQLite support via cross-compilation

## Quick Start (local dev with SQLite)

### 0. Init (clean slate)

```cmd
del C:\Users\farhan\Documents\src\golang\dj_shorturl-go\dj_shorturl-go.db*
```

### 1. Start server

```cmd
cd C:\Users\farhan\Documents\src\golang\dj_shorturl-go
set DB_DRIVER=sqlite3
set SQL_DATABASE=dj_shorturl-go.db
set AUTH_SECRET=dev-auth-secret-change-in-production
set PASSWORD_SECRET=dev-password-secret-change-in-production
set SERVER_PORT=8000
dj_shorturl-go.exe
```

Or in one line (cmd.exe from anywhere):

```cmd
cmd /c "cd /d C:\Users\farhan\Documents\src\golang\dj_shorturl-go & set DB_DRIVER=sqlite3& set SQL_DATABASE=dj_shorturl-go.db& set AUTH_SECRET=dev-auth-secret-change-in-production& set PASSWORD_SECRET=dev-password-secret-change-in-production& set SERVER_PORT=8000& dj_shorturl-go.exe"
```

> **Production:** Change `AUTH_SECRET` and `PASSWORD_SECRET` to random secure values. Anyone with these secrets can forge tokens and passwords.

## API Endpoints

All examples use `curl.exe`.

### Create User (POST /user/signup, no auth required)

```powershell
curl.exe -s -X POST http://localhost:8000/user/signup -H "Content-Type: application/json" -d '{\"id\":\"alice\",\"password_hash\":\"mypassword\"}'
```

Response:
```json
{"id":"alice","password_hash":"ok***hash"}
```

### Login (PUT /user/signup, no auth required)

```powershell
curl.exe -s -X PUT http://localhost:8000/user/signup -H "Content-Type: application/json" -d '{\"id\":\"alice\",\"password_hash\":\"mypassword\"}'
```

Response:
```json
{"token":"<base64-encoded-token>"}
```

### Create Short URL (POST /shorter/url, auth required)

```powershell
curl.exe -s -X POST http://localhost:8000/shorter/url -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{\"url\":\"https://example.com\"}'
```

Response:
```json
{"id":"<uuid>","url":"https://example.com"}
```

### Follow Redirect (GET /shorter/url/{uuid}, no auth required)

```powershell
curl.exe -v http://localhost:8000/shorter/url/<uuid>
```

Returns 302 Found redirecting to the original URL.

## Full Flow

```powershell
# 1. Signup
curl.exe -s -X POST http://localhost:8000/user/signup -H "Content-Type: application/json" -d '{\"id\":\"alice\",\"password_hash\":\"mypassword\"}'

# 2. Login — save the token from response
curl.exe -s -X PUT http://localhost:8000/user/signup -H "Content-Type: application/json" -d '{\"id\":\"alice\",\"password_hash\":\"mypassword\"}'

# 3. Create short URL — replace <token> with value from step 2
curl.exe -s -X POST http://localhost:8000/shorter/url -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{\"url\":\"https://example.com\"}'

# 4. Follow redirect — replace <uuid> with id from step 3
curl.exe -v http://localhost:8000/shorter/url/<uuid>
```

## Full Flow (PowerShell - with variables)

```powershell
# 1. Signup
$body = '{"id":"alice","password_hash":"mypassword"}'
curl.exe -s -X POST http://localhost:8000/user/signup -H "Content-Type: application/json" -d $body

# 2. Login
$body = '{"id":"alice","password_hash":"mypassword"}'
$response = curl.exe -s -X PUT http://localhost:8000/user/signup -H "Content-Type: application/json" -d $body
$token = ($response | ConvertFrom-Json).token

# 3. Create short URL
$body = '{"url":"https://example.com"}'
$response = curl.exe -s -X POST http://localhost:8000/shorter/url -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d $body
$uuid = ($response | ConvertFrom-Json).id

# 4. Follow redirect
curl.exe -v "http://localhost:8000/shorter/url/$uuid"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | `postgres` | Database driver (`sqlite3` or `postgres`) |
| `SERVER_PORT` | `8000` | HTTP server port |
| `SQL_HOST` | `localhost` | PostgreSQL host |
| `SQL_PORT` | `5432` | PostgreSQL port |
| `SQL_USER` | `surl_u` | DB user |
| `SQL_PASSWORD` | `surl_p` | DB password |
| `SQL_DATABASE` | `surl_d` | DB name (or SQLite file path when `DB_DRIVER=sqlite3`) |
| `AUTH_SECRET` | - | Secret for token signing |
| `PASSWORD_SECRET` | - | Secret for password hashing |
| `AUTH_TTL_MINUTES` | `12000` | Token time-to-live |

## Build

```bash
go build -o dj_shorturl-go.exe ./cmd/server
```

For CGO-enabled SQLite support, ensure gcc is in PATH and set `CGO_ENABLED=1`:
```bash
set CGO_ENABLED=1
go build -o dj_shorturl-go.exe ./cmd/server
```

## Original Project

Built by [ftamy9](https://github.com/ftamy9).  
Original Django version: [ftamy9/dj_shorturl](https://github.com/ftamy9/dj_shorturl)  
Go port: [ftamy9/dj_shorturl-go](https://github.com/ftamy9/dj_shorturl-go)
