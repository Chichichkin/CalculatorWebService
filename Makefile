# Makefile for Calculator Web Service

# Variables
SERVICE_NAME=calculator-service
DOCKER_IMAGE=calculator-service:latest
COMPOSE_FILE=docker-compose.yml
MEMORY_SERVICE_PORT=8080
FILE_SERVICE_PORT=8081

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-run compose-up compose-down test-api health

# Default target
help: ## Show this help message
	@echo "$(BLUE)Calculator Web Service - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
build: ## Build the Go application
	@echo "$(BLUE)Building Go application...$(NC)"
	go build -o bin/calculator ./cmd/main.go
	@echo "$(GREEN)Build completed!$(NC)"

run: ## Run the application locally
	@echo "$(BLUE)Running application locally...$(NC)"
	go run ./cmd/main.go

test: ## Run Go tests
	@echo "$(BLUE)Running Go tests...$(NC)"
	go test -v ./...

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -f storage.txt
	@echo "$(GREEN)Clean completed!$(NC)"

# Docker commands
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(NC)"

# Docker Compose commands (primary way to run)
compose-up: ## Start all services with Docker Compose
	@echo "$(BLUE)Starting services with Docker Compose...$(NC)"
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started!$(NC)"
	@echo "$(YELLOW)Memory storage: http://localhost:$(MEMORY_SERVICE_PORT)$(NC)"
	@echo "$(YELLOW)File storage: http://localhost:$(FILE_SERVICE_PORT)$(NC)"

compose-up-monitoring: compose-up ## Start with monitoring (depends on compose-up)
	@echo "$(BLUE)Starting monitoring...$(NC)"
	docker-compose -f $(COMPOSE_FILE) --profile monitoring up -d
	@echo "$(YELLOW)Prometheus: http://localhost:9090$(NC)"

compose-down: ## Stop all services
	@echo "$(BLUE)Stopping all services...$(NC)"
	docker-compose -f $(COMPOSE_FILE) down
	@echo "$(GREEN)All services stopped$(NC)"

compose-logs: ## Show logs from all services
	docker-compose -f $(COMPOSE_FILE) logs -f

compose-restart: compose-down compose-up ## Restart all services

# Testing commands
test-api: ## Test API endpoints (default: memory service)
	@echo "$(BLUE)Testing API endpoints on port $(MEMORY_SERVICE_PORT)...$(NC)"
	@$(MAKE) _test-service PORT=$(MEMORY_SERVICE_PORT)

test-file-storage: ## Test file storage service
	@echo "$(BLUE)Testing file storage service on port $(FILE_SERVICE_PORT)...$(NC)"
	@$(MAKE) _test-service PORT=$(FILE_SERVICE_PORT)

_test-service: ## Internal: Test service on specific port
	@echo "$(YELLOW)Testing service on port $(PORT)...$(NC)"
	@echo "$(YELLOW)Testing addition...$(NC)"
	curl -s -X POST http://localhost:$(PORT)/calculate/addition \
		-H "Content-Type: application/json" \
		-d '{"operand1": 5, "operand2": 3}' | jq .
	@echo "$(YELLOW)Testing recent calculations...$(NC)"
	curl -s -X GET http://localhost:$(PORT)/calculate/recent | jq .
	@echo "$(GREEN)Service test completed!$(NC)"

test-all: ## Run all tests
	@echo "$(BLUE)Running all tests...$(NC)"
	@$(MAKE) test-api
	@$(MAKE) test-file-storage
	@echo "$(GREEN)All tests completed!$(NC)"

# Health check
health: ## Check service health
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo "$(YELLOW)Memory service (port $(MEMORY_SERVICE_PORT)):$(NC)"
	@curl -s http://localhost:$(MEMORY_SERVICE_PORT)/health > /dev/null && \
	 echo "$(GREEN)✓ Healthy$(NC)" || echo "$(RED)✗ Unhealthy$(NC)"
	@echo "$(YELLOW)File service (port $(FILE_SERVICE_PORT)):$(NC)"
	@curl -s http://localhost:$(FILE_SERVICE_PORT)/health > /dev/null && \
	 echo "$(GREEN)✓ Healthy$(NC)" || echo "$(RED)✗ Unhealthy$(NC)"

# Development setup
dev-setup: ## Setup development environment
	go mod download
	go mod tidy
	@echo "$(GREEN)Development environment ready!$(NC)"

# Quick start (alias for compose-up)
quick-start: compose-up ## Quick start services
	@echo "$(YELLOW)Run 'make test-api' to test endpoints$(NC)"
