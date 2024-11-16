.PHONY: swag
swag:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: docs
docs:
	swag init -d ./cmd/instahelp,./internal/echohttp

.PHONY: sql
sql:
	sqlboiler psql

.PHONY: up
up:
	docker compose down
	docker compose up -d --build --force-recreate

.PHONY: stop
stop:
	docker compose stop