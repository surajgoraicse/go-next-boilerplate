# Load .env file :
-include .env
export

.PHONY: build run clear-logs sqlc docker-build docker-up docker-down logs fmt lint clean \
test  swagger db-init db-reset reset-db migrate-up migrate-down migrate-create migrate-force clear-swagger

MAIN := cmd/main.go
SQLC_FILE := internal/db/sqlc.yml
SQLC := $(shell command -v sqlc 2>/dev/null)
PROJECT_NAME := boilerplate
DOCKER_POSTGRES_SERVICE_NAME := postgres

DB_URL:="postgresql://$(DB_USER):$(DB_PASSWORD)@${DOCKER_POSTGRES_SERVICE_NAME}:5432/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# build the server
build:
	@mkdir -p bin
	@echo "Building server..."
	@go build -o bin/main ${MAIN}

fmt:
	gofmt -s -w .

lint:
	golangci-lint run ./...

# run the development server
run:
	@echo "Running server..."
	@go run ${MAIN}

# clear logs files
clear-logs:
	@echo "Clearing logs..."
	@rm -rf logs

test:
	go clean -testcache
	CGO_ENABLED=1 go test -race -count=1 -p=1 ./...

# -------------------- generates swagger docs --------------------
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init \
		-o ./swagger --parseDependency --parseInternal \
		-g ${MAIN}
	@command -v npx > /dev/null 2>&1 && \
		npx prettier --write swagger/swagger.json swagger/swagger.yaml || \
		echo "npx not found, skipping prettier"

# clear swagger docs  :
clear-swagger:
	@echo "Clearing swagger docs..."
	@rm -rf swagger/

# -------------------- Database -------------------- 
sqlc:
	@echo "Generating sqlc queries..."
	@if [ -n "$(SQLC)" ]; then \
		$(SQLC) generate -f ${SQLC_FILE}; \
	else \
		go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate -f ${SQLC_FILE}; \
	fi
	@echo "SQLC queries generated!"

# Database initialization - run only consolidated baseline migration
db-init:
	@echo "Initializing database (via migrate up)..."
	$(MAKE) migrate-up
	@echo "Database initialization complete!"


# Full database reset (drops DB)
reset-db:
	@echo "Resetting entire database (all schemas)..."
	@echo "Terminating active connections..."
	docker compose exec -T db psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${DB_NAME}' AND pid <> pg_backend_pid();"
	docker compose exec -T db psql -U postgres -c "DROP DATABASE IF EXISTS \"${DB_NAME}\";"
	docker compose exec -T db psql -U postgres -c "CREATE DATABASE \"${DB_NAME}\";"
	$(MAKE) db-init
	@echo "Full database reset complete."

# backword compatibility alias
db-reset: reset-db

# Migrations  
migrate-up:
	@echo "Applying migrations (up)"
	docker compose run --rm migrate -path /migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back latest migration"
	docker compose run --rm migrate -path /migrations -database "$(DB_URL)" down 1

migrate-create:
	@if [ -z "$(name)" ]; then \
	  echo "Please provide 'name' var, e.g. 'make migrate-create name=add_thing'"; exit 1; \
	fi
	docker compose run --rm migrate create -ext sql -dir /migrations -seq $(name)

migrate-force:
	@if [ -z "$(version)" ]; then \
	  echo "Please provide 'version' var, e.g. 'make migrate-force version=3'"; exit 1; \
	fi
	docker compose run --rm migrate -path /migrations -database "$(DB_URL)" force $(version)


# -------------------- Docker : --------------------
docker-db-up:
	@echo "Starting docker containers..."
	@docker compose up -d postgres


docker-up:
	@echo "Starting docker containers..."
	@docker compose up -d


docker-down:
	@echo "Stopping docker containers..."
	@docker compose down

docker-build:
	docker compose build


clean:
	docker compose down -v --remove-orphans

logs:
	docker compose logs -f

# -------------------- Observability : --------------------
observability-up:
	@echo "Starting Prometheus and Grafana..."
	@docker compose up -d prometheus grafana

observability-down:
	@echo "Stopping Prometheus and Grafana..."
	@docker compose stop prometheus grafana

