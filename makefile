include .env

create-db:
	docker exec -it postgres createdb --username=${LOCAL_POSTGRES_USER} --owner=${LOCAL_POSTGRES_USER} $(LOCAL_POSTGRES_DBNAME)


local-migration-up:
	@echo "Starting local migration UP"
	migrate -path db/migration -database "postgresql://${LOCAL_POSTGRES_USER}:${LOCAL_POSTGRES_SECRET}@${LOCAL_POSTGRES_HOST}/${LOCAL_POSTGRES_DBNAME}?sslmode=disable" -verbose up

local-migration-down:
	@echo "Starting local migration DOWN"
	migrate -path db/migration -database "postgresql://${LOCAL_POSTGRES_USER}:${LOCAL_POSTGRES_SECRET}@${LOCAL_POSTGRES_HOST}/${LOCAL_POSTGRES_DBNAME}?sslmode=disable" -verbose down

sqlc:
	@echo "Generating sqlc files..."
	sqlc generate

swag:
	@echo "Generating Swagger..."
	swag init --parseDependency --parseInternal -g internal/routers/router.go --output docs/app

compose:
	@echo "Running docker-compose"
	docker-compose up -d

start-server:
	go run cmd/server/main.go

start-bot:
	go run cmd/bot/main.go

setup:
	make compose
	make create-db
	make local-migration-up
	make sqlc
	make swag




.PHONY: sqlc local-migration-up local-migration-down createdb swag setup compose start-server start-bot
