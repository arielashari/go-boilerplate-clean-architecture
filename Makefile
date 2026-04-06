.PHONY: run build migrate-up migrate-down migrate-create seed generate tidy up down

APP_NAME=app
MIGRATE=migrate -path pkg/database/migrations -database "postgresql://postgres:postgres@localhost:5432/app_dev?sslmode=disable"

run:
	air

build:
	go build -o tmp/$(APP_NAME) ./cmd/app

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir pkg/database/migrations -seq $$name

seed:
	go run ./cmd/seeder

generate:
	sqlc generate

tidy:
	go mod tidy

test:
	go test ./... -v -race

test-cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out