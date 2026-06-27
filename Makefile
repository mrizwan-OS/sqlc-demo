.PHONY: help migrate-up migrate-down migrate-create generate run build clean test test-coverage test-integration test-bench

help:
	@echo "Available commands:"
	@echo "  make migrate-up          - Apply all pending migrations"
	@echo "  make migrate-down        - Roll back one migration"
	@echo "  make migrate-create NAME - Create new migration"
	@echo "  make migrate-status      - Show migration status"
	@echo "  make generate           - Run sqlc generate"
	@echo "  make run                - Run the application"
	@echo "  make build              - Build the application"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make test               - Run all tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-bench         - Run benchmarks"
	@echo "  make test-specific TEST=TestName - Run specific test"

migrate-up:
	./migrate.sh up

migrate-down:
	./migrate.sh down

migrate-create:
	./migrate.sh create $(NAME)

migrate-status:
	./migrate.sh status

generate:
	sqlc generate

run:
	go run main.go

build:
	go build -o app main.go

clean:
	rm -f app
	go clean -cache

test:
	go test -v ./...

test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-integration:
	go test -v -run TestIntegrationSuite ./...

test-bench:
	go test -bench=. -benchmem ./...

test-specific:
	go test -v -run $(TEST) ./...
