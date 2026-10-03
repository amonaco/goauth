# goauth

A lightweight Go authentication library and example API built with Go-Chi, Redis, and PostgreSQL for local development. It provides cookie- and token-based session handling, JWT support, middleware for protected routes, and a dev environment for running auth services locally.

## What this project includes

- Middleware for authenticated routes
- Session creation and revocation using Redis
- JWT generation and validation
- Secure cookie helpers
- Local Docker Compose setup for Redis and PostgreSQL
- Seeded development data for quickly testing auth flows

## Requirements

- Go 1.21+
- Docker and Docker Compose

## Quick start

Clone the repo and start the local services:

```bash
git clone https://github.com/amonaco/goauth.git
cd goauth
./scripts/dev-setup.sh
```

Then run the app:

```bash
go run main.go
```

You should be able to reach the API at:

- http://localhost:8080/healthcheck
- http://localhost:5432 for PostgreSQL
- redis://localhost:6379 for Redis

## Local development environment

The project includes a local Docker Compose stack with:

- Redis on `localhost:6379`
- PostgreSQL on `localhost:5432`
- Seeded users and roles for development

### Services and credentials

Redis:

```bash
docker-compose exec redis redis-cli
```

PostgreSQL:

```bash
docker-compose exec postgres psql -U goauth_user -d goauth_db
```

Default development credentials:

- DB user: `goauth_user`
- DB password: `goauth_password`
- DB name: `goauth_db`

### Seeded test users

- `admin@example.com` - superadmin in Acme Corp
- `user@example.com` - standard user in Acme Corp
- `dev@example.com` - admin in Tech Startup, user in Acme Corp

## Configuration

The app reads config from `config/config.yml`.

```yaml
name: 'goauth-api'
environment: 'local'
listen: '0.0.0.0:8080'
redis: 'redis://redis:6379/0'
redis_max_conn: 30
jwt_secret: 'dev-secret-change-in-production-change-in-production-key'
session_ttl: 86400
cookie_secure: false
```

Important:

- For local development, `cookie_secure: false` is acceptable.
- In production, set `jwt_secret` to a real secret and use HTTPS with `cookie_secure: true`.

## Run the app

### Start dependencies

```bash
./scripts/dev-setup.sh
```

### Start the API

```bash
go run main.go
```

### Verify the app is running

```bash
curl http://localhost:8080/healthcheck
```

Expected response:

```json
{"status":"ok"}
```

## Authentication endpoints

### Health check

```bash
curl http://localhost:8080/healthcheck
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "company_id": 1,
    "roles": ["superadmin"]
  }'
```

### Get current user

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer <token>"
```

### Refresh token

```bash
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"token":"<token>"}'
```

### Logout

```bash
curl -X POST http://localhost:8080/logout \
  -H "auth-token: <token>"
```

## Middleware usage

This library is designed to work with Go-Chi and can be applied as router middleware:

```go
import (
    "net/http"

    "github.com/amonaco/goauth/lib/middleware"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Middleware)

    r.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("authorized"))
    })

    _ = http.ListenAndServe(":8080", r)
}
```

## Project structure

```text
.
├── config/
│   └── config.yml                 # Local development configuration
├── lib/
│   ├── auth/                      # Auth helpers, JWT, cookie, session resolution
│   ├── cache/                     # Redis helpers
│   ├── config/                   # Config loading utilities
│   ├── db/                       # PostgreSQL access helpers
│   ├── middleware/               # Auth and rate-limit middleware
│   └── session/                  # Session model and Redis persistence
├── scripts/
│   ├── dev-setup.sh              # Start local dev stack
│   ├── dev-teardown.sh           # Stop local dev stack
│   ├── init.sql                  # PostgreSQL schema
│   └── seed.sql                  # Seed data for development
├── docker-compose.yml            # Redis + PostgreSQL services
├── DEVELOPMENT.md                # Historical dev notes (kept for reference)
├── go.mod
├── main.go                       # Example app entry point
├── README.md
└── LICENSE
```

## Common commands

Start local services:

```bash
./scripts/dev-setup.sh
```

Stop services:

```bash
./scripts/dev-teardown.sh
```

Reset everything including DB data:

```bash
docker-compose down -v
```

Run tests:

```bash
go test ./...
```

View logs:

```bash
docker-compose logs -f
```

## Troubleshooting

### Redis is not available

```bash
docker-compose exec redis redis-cli ping
```

### PostgreSQL is not ready

```bash
docker-compose exec postgres psql -U goauth_user -d goauth_db -c "SELECT 1;"
```

### Bootstrap data did not initialize

```bash
docker-compose down -v
docker-compose up -d
```

## Notes

This repository is currently a work-in-progress auth project intended for local development and experimentation. For production deployment, you should:

- run behind HTTPS
- set a strong JWT secret and rotation policy
- use secure cookie settings
- enforce least-privilege roles from the database
- add a real password hashing provider such as bcrypt or Argon2
- move rate limiting and session state to production-safe infrastructure

