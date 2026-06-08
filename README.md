# dj_shorturl - Go Port

A URL shortener with custom authentication token system, ported from Django to Go.

## Features

- Custom authentication token (not JWT, not as secure - educational)
- URL shortening with UUID-based identifiers
- PostgreSQL database
- Docker Compose setup

## API Endpoints

All endpoints return JSON. Replace `SERVER-OR-LOCAL-IP` with your server.

### Create User (no auth required)

```bash
curl -X POST -d '{ "id": "farhan_10", "password_hash": "123" }' \
  -H "Content-Type: application/json" \
  http://SERVER-OR-LOCAL-IP:80/user/signup
```

### Login (get token)

```bash
curl -i -X PUT -d '{ "id": "farhan_10", "password_hash": "123" }' \
  -H "Content-Type: application/json" \
  http://SERVER-OR-LOCAL-IP:80/user/signup
```

### Create Short URL (requires auth)

```bash
curl -i -X POST -d '{ "url": "http://icanhazip.com" }' \
  -H "Authorization: Basic <token>" \
  -H "Content-Type: application/json" \
  http://SERVER-OR-LOCAL-IP:80/shorter/url
```

### Follow Redirect (no auth required)

```bash
curl -iL http://SERVER-OR-LOCAL-IP:80/shorter/url/<uuid>
```

## Run Locally

```bash
docker-compose up -d --build
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8000` | HTTP server port |
| `SQL_HOST` | `localhost` | PostgreSQL host |
| `SQL_PORT` | `5432` | PostgreSQL port |
| `SQL_USER` | `surl_u` | DB user |
| `SQL_PASSWORD` | `surl_p` | DB password |
| `SQL_DATABASE` | `surl_d` | DB name |
| `AUTH_SECRET` | - | Secret for token signing |
| `PASSWORD_SECRET` | - | Secret for password hashing |
| `AUTH_TTL_MINUTES` | `12000` | Token time-to-live |

## Original Project

Built by [ftamy9](https://github.com/ftamy9). Find the original Django version at [ftamy9/dj_shorturl](https://github.com/ftamy9/dj_shorturl).
