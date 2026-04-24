# Load .env file :
-include .env
export

.PHONY: build run clear-logs swagger sqlc validate

MAIN="cmd/main.go"
SQLC_FILE="db/sqlc.yaml"
PROJECT_NAME="boilerplate"

# build the server
build:
	@echo "Building server..."
	@go build -o bin/main ${MAIN}

# run the development server
run:
	@echo "Running server..."
	@go run ${MAIN}

# clear logs files
clear-logs:
	@echo "Clearing logs..."
	@rm -rf logs

# -------------------- generates swagger docs --------------------
swagger:
	@echo "Generating swagger docs..."
	swag init -o ./swagger --parseDependency --parseInternal -g $(MAIN)

# clear swagger docs  : 
clear-swagger:
	@echo "Clearing swagger docs..."
	@rm -rf swagger

# -------------------- Database -------------------- 
sqlc:
	@echo "Generating sqlc queries..."
	sqlc generate -f ${SQLC_FILE}

# Database initialization - run only consolidated baseline migration
db-init:
	@echo "Initializing database (via migrate up)..."
	$(MAKE) migrate-up
	@echo "Database initialization complete!"


# Full database reset 
reset-db:
	@echo "Dropping database schema..."
	docker compose exec postgres psql -U $(PROJECT_NAME) -d coderz -c "DROP SCHEMA IF EXISTS coderz CASCADE;"
	docker compose run --rm migrate -path /migrations -database "${DB_DSN}" drop -f
	@echo "Re-applying all migrations..."
	$(MAKE) migrate-up
	@echo "Database reset complete!"

# Migrations  
migrate-up:
	@echo "Applying migrations (up)"
	docker compose run --rm migrate -path /migrations -database "${DB_URL}" up

migrate-down:
	@echo "Rolling back latest migration"
	docker compose run --rm migrate -path /migrations -database "${DB_URL}" down 1

migrate-force:
	@echo "Forcing migration version (e.g., version=1)"
	docker compose run --rm migrate -path /migrations -database "${DB_URL}" force $(version)


# -------------------- Docker : --------------------
docker-up:
	@echo "Starting docker containers..."
	@docker compose up -d


docker-down:
	@echo "Stopping docker containers..."
	@docker compose down