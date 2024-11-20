include .env

.SILENT:

DB_URL=postgresql://postgres:1105@localhost:5432/$(POSTGRES_DATABASE)?sslmode=disable

start:
	go run cmd/main.go

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

local-up:
	docker compose --env-file ./.env.docker up -d


.PHONY: start migrateup migratedown