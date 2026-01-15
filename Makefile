.PHONY: help build build-frontend build-backend test run clean cleanup update docker-build docker-run dev-frontend dev-backend install-deps

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=mucal
VERSION?=dev
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
WEB_DIR=web
DIST_DIR=$(WEB_DIR)/dist
DOCKER_IMAGE=mucal
DOCKER_TAG?=latest

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

help: ## Show this help message
	@echo "$(COLOR_BOLD)μCal Makefile$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Available targets:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(COLOR_BLUE)%-20s$(COLOR_RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

install-deps: ## Install all dependencies (Go modules + npm packages)
	@echo "$(COLOR_GREEN)Installing Go dependencies...$(COLOR_RESET)"
	go mod download
	@echo "$(COLOR_GREEN)Installing frontend dependencies...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm ci

build: build-frontend build-backend ## Build complete application (frontend + backend)

build-frontend: ## Build frontend (Svelte + Vite)
	@echo "$(COLOR_GREEN)Building frontend...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm run build
	@echo "$(COLOR_GREEN)Frontend built successfully in $(DIST_DIR)$(COLOR_RESET)"

build-backend: build-frontend ## Build Go backend with embedded frontend
	@echo "$(COLOR_GREEN)Building backend...$(COLOR_RESET)"
	go build -ldflags="-X 'github.com/mano/mucal/internal/version.Version=$(VERSION)'" -o $(BINARY_NAME) ./cmd/mucal
	@echo "$(COLOR_GREEN)Backend built successfully: $(BINARY_NAME)$(COLOR_RESET)"

test: ## Run tests (Go + Frontend checks)
	@echo "$(COLOR_GREEN)Running Go tests...$(COLOR_RESET)"
	go test -v ./...
	@echo "$(COLOR_GREEN)Running frontend checks...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm run check

run: build ## Build and run the application locally
	@echo "$(COLOR_GREEN)Starting μCal...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Make sure you have a config.yaml file!$(COLOR_RESET)"
	./$(BINARY_NAME)

dev-frontend: ## Run frontend development server (Vite)
	@echo "$(COLOR_GREEN)Starting frontend dev server...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm run dev

dev-backend: ## Run backend in development mode (requires built frontend)
	@echo "$(COLOR_GREEN)Starting backend in dev mode...$(COLOR_RESET)"
	go run ./cmd/mucal

clean: ## Remove build artifacts (binary and frontend dist)
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)
	@echo "$(COLOR_GREEN)Build artifacts cleaned$(COLOR_RESET)"

cleanup: ## Remove all artifacts, caches, and dependencies
	@echo "$(COLOR_YELLOW)Removing all artifacts, caches, and dependencies...$(COLOR_RESET)"
	@echo "Removing binary..."
	rm -f $(BINARY_NAME)
	@echo "Removing frontend build artifacts..."
	rm -rf $(DIST_DIR)
	@echo "Removing frontend dependencies..."
	rm -rf $(WEB_DIR)/node_modules
	@echo "Removing Go build cache..."
	go clean -cache -modcache -testcache
	@echo "Removing Go vendor directory (if exists)..."
	rm -rf vendor
	@echo "$(COLOR_GREEN)Complete cleanup finished$(COLOR_RESET)"

update: ## Update all dependencies to latest versions
	@echo "$(COLOR_GREEN)Updating Go dependencies...$(COLOR_RESET)"
	go get -u ./...
	go mod tidy
	@echo "$(COLOR_GREEN)Updating frontend dependencies...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm update
	cd $(WEB_DIR) && npm install svelte@latest vite@latest @sveltejs/vite-plugin-svelte@latest
	cd $(WEB_DIR) && npm audit fix || true
	@echo "$(COLOR_GREEN)All dependencies updated!$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Run 'make test' to verify everything works$(COLOR_RESET)"

docker-build: ## Build Docker image
	@echo "$(COLOR_GREEN)Building Docker image...$(COLOR_RESET)"
	docker build --build-arg VERSION=$(VERSION) -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(COLOR_GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(COLOR_RESET)"

docker-run: ## Run application in Docker
	@echo "$(COLOR_GREEN)Running μCal in Docker...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Make sure you have config.yaml and secrets/ directory!$(COLOR_RESET)"
	docker run --rm -it \
		-p 8080:8080 \
		-v $(PWD)/config.yaml:/config/config.yaml:ro \
		-v $(PWD)/secrets:/secrets:ro \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop all running μCal containers
	@echo "$(COLOR_YELLOW)Stopping all μCal containers...$(COLOR_RESET)"
	docker ps -a | grep $(DOCKER_IMAGE) | awk '{print $$1}' | xargs -r docker stop
	docker ps -a | grep $(DOCKER_IMAGE) | awk '{print $$1}' | xargs -r docker rm

fmt: ## Format code (Go + Frontend)
	@echo "$(COLOR_GREEN)Formatting Go code...$(COLOR_RESET)"
	go fmt ./...
	@echo "$(COLOR_GREEN)Formatting frontend code...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm run format || echo "$(COLOR_YELLOW)Frontend formatting not configured$(COLOR_RESET)"

lint: ## Run linters
	@echo "$(COLOR_GREEN)Running Go linters...$(COLOR_RESET)"
	go vet ./...
	@echo "$(COLOR_GREEN)Running frontend checks...$(COLOR_RESET)"
	cd $(WEB_DIR) && npm run check

all: cleanup install-deps build test ## Complete build from scratch (cleanup, install, build, test)
	@echo "$(COLOR_GREEN)Complete build successful!$(COLOR_RESET)"

.PHONY: check-config
check-config: ## Verify config.yaml exists
	@if [ ! -f config.yaml ]; then \
		echo "$(COLOR_YELLOW)Warning: config.yaml not found. Copy config.example.yaml to config.yaml$(COLOR_RESET)"; \
		exit 1; \
	fi

quick: build-backend ## Quick build (assumes frontend is already built)
	@echo "$(COLOR_GREEN)Quick build complete$(COLOR_RESET)"

version: ## Show current version
	@echo "Version: $(VERSION)"
	@if [ -f $(BINARY_NAME) ]; then \
		./$(BINARY_NAME) -help 2>&1 | grep -i version || echo "Binary version: $(VERSION)"; \
	fi
