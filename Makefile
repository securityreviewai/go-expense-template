.PHONY: run build tidy docker migrate-up migrate-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy

docker:
	docker build -t go-expense-boilerplate:local .

migrate-up:
	psql "$$DATABASE_URL" -f migrations/001_init.up.sql

migrate-down:
	psql "$$DATABASE_URL" -f migrations/001_init.down.sql
