.PHONY: swagger swagger-concat swagger-server swagger-install setup build clean test help mockgen

# Default target
all: setup build

# Setup environment and install dependencies
setup:
	@echo "Setting up development environment..."
	@echo "Installing system dependencies..."
	sudo apt-get update
	sudo apt-get install -y build-essential pkg-config librdkafka-dev
	@echo "Setting CGO_ENABLED=1..."
	export CGO_ENABLED=1
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing Swagger..."
	$(MAKE) swagger-install
	@echo "Installing mockgen..."
	go install github.com/golang/mock/mockgen@v1.6.0
	@echo "Generating mocks..."
	$(MAKE) mockgen
	@echo "Generating Swagger documentation..."
	$(MAKE) swagger

# Swagger commands
swagger:
	@$(MAKE) swagger-concat
	@$(MAKE) swagger-server

swagger-concat:
	@echo "Concatenating Swagger files..."
	@rm -rf docs/api/tmp
	@mkdir -p docs/api/tmp
	./swagger mixin --output=docs/api/tmp/tmp.yml --format=yaml --keep-spec-order \
		docs/api/configs/main.yml docs/api/paths/*
	./swagger flatten docs/api/tmp/tmp.yml --output=docs/api/swagger.yaml --format=yaml
	./swagger flatten docs/api/tmp/tmp.yml --output=docs/api/swagger.json --format=json
	@rm -f docs/api/tmp/tmp.yml

swagger-server:
	@echo "Generating Swagger server code..."
	@mkdir ./internal/generated || true
	@rm -rf ./internal/generated/api_models
	./swagger generate model \
		--allow-template-override \
		--spec=docs/api/swagger.yaml \
		--target=internal/generated \
		--model-package=api_models

swagger-install:
	@echo "Installing Swagger CLI..."
	if [ ! -f "./swagger" ]; then \
		wget "https://github.com/go-swagger/go-swagger/releases/download/v0.30.3/swagger_`go env GOOS`_`go env GOARCH`" -O ./swagger && \
		chmod +x ./swagger; \
	fi

new-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir docs/db/migrations -seq $$migration_name

db-migrate-up:
	@go run main.go db-migrate-up

mockgen:
	@echo "Generating mocks..."
	@mkdir -p internal/mocks/domain
	mockgen -source=internal/domain/user/user.go -destination=internal/mocks/domain/user_mock.go -package=mocks
	mockgen -source=internal/domain/notification/notification.go -destination=internal/mocks/domain/notification_mock.go -package=mocks
	mockgen -source=internal/domain/outbox/outbox.go -destination=internal/mocks/domain/outbox_mock.go -package=mocks

# Build the project
build:
	@echo "Building project..."
	go build ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	go clean -modcache
	rm -f $(BINARY)
	@echo "Cleaning Swagger generated files..."
	rm -rf ./internal/generated/api_models

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Help command
help:
	@echo "Available commands:"
	@echo "  make setup    - Install dependencies and setup environment"
	@echo "  make build    - Build the project"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make swagger  - Generate Swagger documentation and code"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Note: Before running the project, make sure to:"
	@echo "1. Run 'make setup' to install dependencies"
	@echo "2. Set CGO_ENABLED=1 in your environment"
	@echo "3. Have librdkafka-dev installed on your system"
	@echo "4. Swagger CLI will be installed automatically during setup"