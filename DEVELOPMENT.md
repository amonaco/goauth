# goauth Development Guide

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+

### Setup

1. **Start the development environment:**
   ```bash
   ./scripts/dev-setup.sh
   ```
   This will start Redis and PostgreSQL with test data.

2. **Run the API:**
   ```bash
   go run main.go
   ```

3. **Test the API:**
   ```bash
   # Health check (no auth required)
   curl http://localhost:8080/healthcheck

   # Login
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": 1,
       "company_id": 1,
       "roles": ["superadmin"]
     }'

   # Use the token from login response
   TOKEN="your-token-here"

   # Refresh token
   curl -X POST http://localhost:8080/refresh \
     -H "Content-Type: application/json" \
     -d "{\"token\": \"$TOKEN\"}"

   # Get current user
   curl http://localhost:8080/me \
     -H "Authorization: Bearer $TOKEN"

   # Logout
   curl -X POST http://localhost:8080/logout \
     -H "auth-token: $TOKEN"
   ```

## Development Services

### Redis
- **URL:** `redis://localhost:6379`
- **CLI:** `docker-compose exec redis redis-cli`

### PostgreSQL
- **Connection:** `postgres://goauth_user:goauth_password@localhost:5432/goauth_db`
- **CLI:** `docker-compose exec postgres psql -U goauth_user -d goauth_db`

### Test Data

**Users:**
- `admin@example.com` - superadmin in Acme Corp
- `user@example.com` - standard user in Acme Corp
- `dev@example.com` - admin in Tech Startup, user in Acme Corp

**Companies:**
- Acme Corp (ID: 1)
- Tech Startup (ID: 2)

**Roles:**
- user (ID: 1)
- admin (ID: 2)
- superadmin (ID: 3)

## Common Commands

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f redis
docker-compose logs -f postgres
```

### Run tests
```bash
go test ./...
```

### Clean up
```bash
# Stop services
./scripts/dev-teardown.sh

# Remove volumes and data
docker-compose down -v
```

## Configuration

The API reads from `config/config.yml`. For development, this is pre-configured with:
- Redis connection string
- JWT secret (for development only)
- Session TTL (24 hours)
- Insecure cookies (enable only on localhost)

**Important:** Change `jwt_secret` and set `cookie_secure: true` in production.

## Troubleshooting

### Redis connection refused
```bash
# Check Redis is running
docker-compose exec redis redis-cli ping

# Restart Redis
docker-compose restart redis
```

### PostgreSQL connection refused
```bash
# Check PostgreSQL is running
docker-compose exec postgres psql -U goauth_user -d goauth_db -c "SELECT 1;"

# View logs
docker-compose logs postgres
```

### Migrations not applied
```bash
# Re-initialize database
docker-compose down -v
docker-compose up -d
```
