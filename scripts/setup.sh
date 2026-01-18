#!/bin/bash
# Quick start script for InstantTLS development

set -e

echo "🚀 Starting InstantTLS development environment..."
echo ""

# Check dependencies
echo "Checking dependencies..."

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+"
    exit 1
fi
echo "✅ Go $(go version | awk '{print $3}')"

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+"
    exit 1
fi
echo "✅ Node.js $(node --version)"

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker"
    exit 1
fi
echo "✅ Docker $(docker --version | awk '{print $3}' | tr -d ',')"

echo ""

# Copy env if not exists
if [ ! -f ".env" ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
fi

# Start Postgres
echo "Starting PostgreSQL..."
docker-compose -f infra/docker-compose.yml up -d
sleep 3

# Run migrations
echo "Running migrations..."
cd apps/api
go mod download
go run . migrate up
go run . seed
cd ../..

echo ""
echo "✅ Database ready with demo user:"
echo "   Email: demo@instanttls.dev"
echo "   Password: demo1234"
echo ""

# Install web dependencies
echo "Installing web dependencies..."
cd apps/web
npm install
cd ../..

echo ""
echo "🎉 Setup complete!"
echo ""
echo "To start development:"
echo ""
echo "  Terminal 1 (API):"
echo "    cd apps/api && go run ."
echo ""
echo "  Terminal 2 (Web):"
echo "    cd apps/web && npm run dev"
echo ""
echo "  Terminal 3 (CLI):"
echo "    cd cli && go build -o ../bin/instanttls . && ../bin/instanttls --help"
echo ""
