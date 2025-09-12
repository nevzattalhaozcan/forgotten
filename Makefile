.PHONY: build run test clean dev deps

DOCKER_COMPOSE_FILE = docker/docker-compose.yml

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run with live reload (requires air)
dev:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Run all tests
test:
	go test -v ./...

# Run unit tests only
test-unit:
	go test -v -short ./...

# Run integration tests only  
test-integration:
	go test -v -run Integration ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	go test -v -race ./...

# Clean test cache
test-clean:
	go clean -testcache

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install development tools
tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database migrations (example)
migrate-up:
	migrate -path internal/database/migrations -database $(DATABASE_URL) up

migrate-down:
	migrate -path internal/database/migrations -database $(DATABASE_URL) down

docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# GitHub Container Registry settings
GHCR_IMAGE = ghcr.io/nevzattalhaozcan/forgotten-app
DOCKER_TAG = latest

# Build the Docker image
docker-build:
	docker build -f docker/Dockerfile -t $(GHCR_IMAGE):$(DOCKER_TAG) .

# Login to GitHub Container Registry
ghcr-login:
	@echo "Please create a Personal Access Token with 'write:packages' permission"
	@echo "Go to: https://github.com/settings/tokens"
	docker login ghcr.io

# Build and push to GitHub Container Registry
ghcr-push: docker-build
	docker push $(GHCR_IMAGE):$(DOCKER_TAG)
	@echo "Image pushed to: $(GHCR_IMAGE):$(DOCKER_TAG)"

# Pull from GitHub Container Registry
ghcr-pull:
	docker pull $(GHCR_IMAGE):$(DOCKER_TAG)

# Update docker-compose to use GHCR image for sharing
ghcr-compose-generate:
	@echo "Creating docker-compose.ghcr.yml for sharing..."
	@sed 's|build:|# build:|g; s|context: ..|# context: ..|g; s|dockerfile: docker/Dockerfile|# dockerfile: docker/Dockerfile|g' $(DOCKER_COMPOSE_FILE) > docker/docker-compose.ghcr.yml
	@sed -i '' '/# build:/a\'$$'\n''    image: $(GHCR_IMAGE):$(DOCKER_TAG)' docker/docker-compose.ghcr.yml
	@echo "Generated docker/docker-compose.ghcr.yml for QA sharing"