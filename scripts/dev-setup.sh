#!/bin/bash
# Local development setup script

set -e

echo "🚀 Starting goauth development environment..."

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install Docker Desktop or docker-compose."
    exit 1
fi

# Start services
echo "🐳 Starting Docker services..."
docker-compose up -d

echo "⏳ Waiting for services to be healthy..."

# Wait for Redis
echo "   Waiting for Redis..."
for i in {1..30}; do
    if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
        echo "   ✅ Redis is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ❌ Redis failed to start"
        exit 1
    fi
    sleep 1
done

# Wait for PostgreSQL
echo "   Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U goauth_user -d goauth_db > /dev/null 2>&1; then
        echo "   ✅ PostgreSQL is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ❌ PostgreSQL failed to start"
        exit 1
    fi
    sleep 1
done

echo ""
echo "✅ Development environment is ready!"
echo ""
echo "📋 Services:"
echo "   - Redis:      redis://localhost:6379"
echo "   - PostgreSQL: postgres://goauth_user:goauth_password@localhost:5432/goauth_db"
echo ""
echo "👤 Test Users (PostgreSQL):"
echo "   - admin@example.com (superadmin in Acme Corp)"
echo "   - user@example.com (user in Acme Corp)"
echo "   - dev@example.com (admin in Tech Startup, user in Acme Corp)"
echo ""
echo "🔧 Next steps:"
echo "   1. Run: go run main.go"
echo "   2. Test: curl http://localhost:8080/healthcheck"
echo ""
echo "To stop services: docker-compose down"
echo "To view logs: docker-compose logs -f [service-name]"
