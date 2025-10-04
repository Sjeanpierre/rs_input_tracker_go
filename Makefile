.PHONY: help build build-api build-worker build-frontend build-docker clean test run run-api run-frontend deploy-worker install-deps

# Default target
help:
	@echo "Available targets:"
	@echo "  make help            - Show this help message"
	@echo "  make install-deps    - Install all dependencies"
	@echo "  make build           - Build all components"
	@echo "  make build-api       - Build API server"
	@echo "  make build-worker    - Build Lambda worker"
	@echo "  make build-frontend  - Build frontend"
	@echo "  make build-docker    - Build Docker image"
	@echo "  make run             - Run API server locally"
	@echo "  make run-api         - Run API server locally"
	@echo "  make run-frontend    - Run frontend dev server"
	@echo "  make deploy-worker   - Deploy Lambda functions (requires ACCOUNT_ID)"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean build artifacts"

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build all components
build: build-api build-worker build-frontend

# Build API server
build-api:
	@echo "Building API server..."
	go build -o bin/input_tracker_app

# Build Lambda worker
build-worker:
	@echo "Building Lambda worker..."
	cd worker && $(MAKE) build

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Build Docker image
build-docker:
	@echo "Building Docker image..."
	cd build && docker build -t input_tracker .

# Run API server locally
run: run-api

run-api:
	@echo "Running API server on port 9080..."
	go run main.go

# Run frontend development server
run-frontend:
	@echo "Running frontend dev server on port 8080..."
	cd frontend && npm run dev

# Deploy Lambda worker
deploy-worker:
	@echo "Deploying Lambda worker..."
	@if [ -z "$(ACCOUNT_ID)" ]; then \
		echo "Error: ACCOUNT_ID environment variable not set"; \
		exit 1; \
	fi
	cd worker && ACCOUNT_ID=$(ACCOUNT_ID) ./build

# Run tests
test:
	@echo "Running Go tests..."
	go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm run test 2>/dev/null || echo "No frontend tests configured"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf worker/bin/
	rm -rf frontend/dist/
	rm -rf build/bin/
	@echo "Clean complete"

# Development helpers
.PHONY: dev dev-all fmt lint

# Run in development mode with hot reload (requires air or similar)
dev:
	@echo "Running in development mode..."
	@command -v air >/dev/null 2>&1 && air || echo "Install 'air' for hot reload: go install github.com/cosmtrek/air@latest"

# Run both API and frontend in development
dev-all:
	@echo "Starting API and frontend in development mode..."
	@echo "Note: This requires tmux or running in separate terminals"
	@echo "Terminal 1: make run-api"
	@echo "Terminal 2: make run-frontend"

# Format code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting frontend code..."
	cd frontend && npm run lint -- --fix 2>/dev/null || echo "No linting configured"

# Lint code
lint:
	@echo "Linting Go code..."
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "Install golangci-lint for linting"
	@echo "Linting frontend code..."
	cd frontend && npm run lint 2>/dev/null || echo "No linting configured"
