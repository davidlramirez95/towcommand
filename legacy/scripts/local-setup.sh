#!/usr/bin/env bash
set -euo pipefail

echo "🚀 TowCommand PH — Local Development Setup"
echo "============================================"

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required"; exit 1; }
command -v pnpm >/dev/null 2>&1 || { echo "❌ pnpm is required (npm install -g pnpm)"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js 20+ is required"; exit 1; }

NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 20 ]; then
  echo "❌ Node.js 20+ required (found v$NODE_VERSION)"
  exit 1
fi

echo "✅ Prerequisites check passed"

# Copy env file if not exists
if [ ! -f .env ]; then
  cp .env.example .env
  echo "📄 Created .env from .env.example"
fi

# Start Docker services
echo "🐳 Starting Docker services..."
docker-compose up -d

# Wait for services
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check LocalStack
until docker-compose exec -T localstack awslocal sts get-caller-identity > /dev/null 2>&1; do
  echo "  Waiting for LocalStack..."
  sleep 2
done
echo "✅ LocalStack ready"

# Check Redis
until docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; do
  echo "  Waiting for Redis..."
  sleep 1
done
echo "✅ Redis ready"

# Check PostgreSQL
until docker-compose exec -T postgres pg_isready > /dev/null 2>&1; do
  echo "  Waiting for PostgreSQL..."
  sleep 1
done
echo "✅ PostgreSQL ready"

# Create DynamoDB table in LocalStack
echo "📊 Creating DynamoDB table..."
DYNAMODB_ENDPOINT=http://localhost:4566 pnpm run db:migrate 2>/dev/null || echo "  Table already exists"

# Install dependencies
echo "📦 Installing dependencies..."
pnpm install

# Build packages
echo "🔨 Building packages..."
pnpm run build

echo ""
echo "✅ Local development environment is ready!"
echo ""
echo "Available services:"
echo "  LocalStack:       http://localhost:4566"
echo "  Redis:            localhost:6379"
echo "  Redis Commander:  http://localhost:8081"
echo "  PostgreSQL:       localhost:5432"
echo ""
echo "Run 'pnpm run dev' to start development"
