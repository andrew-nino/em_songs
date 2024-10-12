.PHONY: run

run: 
	go run cmd/songs/main.go cmd/songs/logger.go -path_to_config .env

build:
	docker compose  build $(c)

up:
	docker compose  up -d $(c)

start:
	docker compose  start $(c)

down:
	docker compose  down $(c)

destroy:
	docker compose  down -v $(c)

stop:
	docker compose  stop $(c)

restart:
	docker compose  stop $(c)

	docker compose  up -d $(c)

logs:
	docker compose  logs --tail=100 -f $(c)

logs-app:
	docker compose  logs --tail=100 -f app

ps:
	docker compose  ps

login-app:
	docker compose  exec app sh

db-psql:
	docker compose  exec postgres psql -Upostgres