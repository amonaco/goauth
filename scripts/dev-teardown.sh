#!/bin/bash
# Stop and remove all development containers

echo "🛑 Stopping goauth development environment..."
docker-compose down

echo "✅ Development environment stopped."
echo ""
echo "To remove volumes and data:"
echo "   docker-compose down -v"
