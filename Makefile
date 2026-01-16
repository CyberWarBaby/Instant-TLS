.PHONY: dev run-api run-web build-cli migrate-up migrate-down docker-up docker-down clean

# Default target
all: dev

# Start everything for development
dev: docker-up
	@echo "Starting development environment..."
	@make -j3 run-api run-web wait-ready

wait-ready:
	@echo "Waiting for services to be ready..."
	@sleep 3
	@echo "✅ Development environment is ready!"
	@echo "   API: http://localhost:8081"
	@echo "   Web: http://localhost:3000"

# Docker
docker-up:
	@echo "Starting PostgreSQL..."
	@docker-compose -f infra/docker-compose.yml up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@make migrate-up

docker-down:
	@docker-compose -f infra/docker-compose.yml down

# API
run-api:
	@echo "Starting API server..."
	@cd apps/api && go run .

# Web
run-web:
	@echo "Starting Next.js dashboard..."
	@cd apps/web && npm install && npm run dev

# CLI
build-cli:
	@echo "Building CLI..."
	@cd cli && go build -o ../bin/instanttls .
	@echo "✅ CLI built to ./bin/instanttls"

install-cli: build-cli
	@echo "Installing CLI..."
	@sudo cp ./bin/instanttls /usr/local/bin/
	@echo "✅ CLI installed to /usr/local/bin/instanttls"

# Migrations
migrate-up:
	@echo "Running migrations..."
	@cd apps/api && go run . migrate up

migrate-down:
	@echo "Rolling back migrations..."
	@cd apps/api && go run . migrate down

# Seed demo data
seed:
	@echo "Seeding demo data..."
	@cd apps/api && go run . seed

# Clean
clean:
	@rm -rf bin/
	@rm -rf apps/web/.next
	@rm -rf apps/web/node_modules
	@docker-compose -f infra/docker-compose.yml down -v

# Test
test:
	@echo "Running API tests..."
	@cd apps/api && go test ./...
	@echo "Running CLI tests..."
	@cd cli && go test ./...

# Build for release
build-all: build-cli
	@echo "Building for all platforms..."
	@cd cli && GOOS=darwin GOARCH=amd64 go build -o ../bin/instanttls-darwin-amd64 .
	@cd cli && GOOS=darwin GOARCH=arm64 go build -o ../bin/instanttls-darwin-arm64 .
	@cd cli && GOOS=linux GOARCH=amd64 go build -o ../bin/instanttls-linux-amd64 .
	@cd cli && GOOS=linux GOARCH=arm64 go build -o ../bin/instanttls-linux-arm64 .
	@cd cli && GOOS=windows GOARCH=amd64 go build -o ../bin/instanttls-windows-amd64.exe .
	@echo "✅ All binaries built in ./bin/"

# Generate checksums
checksums:
	@cd bin && sha256sum instanttls-* > checksums.txt
	@echo "✅ Checksums generated in ./bin/checksums.txt"
