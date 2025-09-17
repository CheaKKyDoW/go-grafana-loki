.PHONY: build run test clean docker-build docker-run migrate-up migrate-down

.PHONY: build run test clean docker-build docker-run migrate-up migrate-down deps fmt lint mocks

# Main API commands
build:
	go build -o bin/main cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

# User Service commands
build-user-service:
	cd services/user-service && go build -o ../../bin/user-service cmd/api/main.go

run-user-service:
	cd services/user-service && go run cmd/api/main.go

test-user-service:
	cd services/user-service && go test -v ./...

# Multi-service commands
build-all:
	make build
	make build-user-service

run-all-local:
	make run &
	make run-user-service &
	wait

test-all:
	make test
	make test-user-service

# Clean build artifacts
clean:
	rm -rf bin/
	cd services/user-service && rm -rf bin/

# Docker commands
docker-build:
	docker build -t go-clean-arch-api .
	docker build -t go-clean-arch-user-service ./services/user-service

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/cleanarch?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/cleanarch?sslmode=disable" down

# Dependencies
deps:
	go mod download
	go mod tidy
	cd services/user-service && go mod download && go mod tidy

# Code quality
fmt:
	go fmt ./...
	cd services/user-service && go fmt ./...

lint:
	golangci-lint run
	cd services/user-service && golangci-lint run

# Development helpers
dev-setup:
	make deps
	make migrate-up

logs-loki:
	curl -G -s "http://localhost:3100/loki/api/v1/query" --data-urlencode 'query={service="go-clean-arch-api"}' | jq

logs-grafana:
	@echo "Grafana available at: http://localhost:3000"
	@echo "Loki datasource configured automatically"

# Generate mocks
mocks:
	mockery --all --output=mocks