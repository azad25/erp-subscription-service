# Subscription Service Makefile

.PHONY: help build run test clean proto migrate seed docker-build docker-run dev

# Variables
SERVICE_NAME = subscription-service
BINARY_NAME = subscription-service
MAIN_PATH = ./cmd/subscription-service
BUILD_DIR = bin
DOCKER_IMAGE = erp-subscription-service
DOCKER_TAG = latest

# Default target
help:
	@echo "Subscription Service - Available Commands"
	@echo ""
	@echo "🚀 Development Commands:"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  run           - Run the service"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "🧪 Testing Commands:"
	@echo "  test          - Run all tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  test-race     - Run tests with race detection"
	@echo ""
	@echo "🔧 Setup Commands:"
	@echo "  proto         - Generate protobuf files"
	@echo "  migrate       - Run database migrations"
	@echo "  seed          - Seed database with initial data"
	@echo "  deps          - Install dependencies"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo ""
	@echo "📊 Utility Commands:"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  vet           - Vet code"
	@echo "  mod-tidy      - Tidy go modules"

# Development commands
dev:
	@echo "🚀 Starting subscription service in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  Air not found, installing..." && \
		go install github.com/cosmtrek/air@latest && \
		air; \
	fi

run:
	@echo "🚀 Starting subscription service..."
	@go run $(MAIN_PATH)/main.go

build:
	@echo "🔨 Building subscription service..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)/main.go
	@echo "✅ Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✅ Cleaned"

# Testing commands
test:
	@echo "🧪 Running tests..."
	@go test ./...

test-cover:
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-race:
	@echo "🧪 Running tests with race detection..."
	@go test -race ./...

# Setup commands
proto:
	@echo "🔧 Generating protobuf files..."
	@if command -v protoc > /dev/null; then \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/*.proto; \
	else \
		echo "❌ protoc not found. Please install Protocol Buffers compiler."; \
		exit 1; \
	fi
	@echo "✅ Protobuf files generated"

migrate:
	@echo "🗄️  Running database migrations..."
	@go run cmd/migrate/main.go

seed:
	@echo "🌱 Seeding database..."
	@go run cmd/seed/main.go

deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8081:8081 -p 50051:50051 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Utility commands
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, installing..." && \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
		golangci-lint run; \
	fi

vet:
	@echo "🔍 Vetting code..."
	@go vet ./...
	@echo "✅ Code vetted"

mod-tidy:
	@echo "📦 Tidying go modules..."
	@go mod tidy
	@echo "✅ Go modules tidied"

# Cross-compilation
build-linux:
	@echo "🔨 Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)/main.go

build-windows:
	@echo "🔨 Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)/main.go

build-mac:
	@echo "🔨 Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_PATH)/main.go

# All platforms
build-all: build-linux build-windows build-mac
	@echo "✅ Built for all platforms"

# Health check
health:
	@echo "🏥 Checking service health..."
	@curl -f http://localhost:8081/health || echo "❌ Service not responding"

# Database commands
db-create:
	@echo "🗄️  Creating database..."
	@createdb erp_subscription || echo "Database already exists"

db-drop:
	@echo "🗄️  Dropping database..."
	@dropdb erp_subscription || echo "Database does not exist"

db-reset: db-drop db-create migrate seed
	@echo "✅ Database reset complete"

# Environment setup
env-setup:
	@echo "⚙️  Setting up environment..."
	@cp .env.example .env || echo "⚠️  .env.example not found"
	@echo "✅ Environment setup complete"

# Full development setup
setup: env-setup deps proto migrate seed
	@echo "✅ Full development setup complete"

# Production build
prod-build:
	@echo "🚀 Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)/main.go
	@echo "✅ Production binary built" 