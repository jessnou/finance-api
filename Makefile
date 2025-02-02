.PHONY: run migrate-up migrate-down

include .env

export $(shell sed 's/=.*//g' .env)

run:
	docker-compose up --build -d
	docker logs -f finance_app

migrate-up:
	docker exec finance_app goose -dir migrations postgres \
		"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

migrate-down:
	docker exec finance_app goose -dir migrations postgres \
		"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" down
