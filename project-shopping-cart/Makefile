include .env
export

APP_NAME := project-shopping-cart
CMD_DIR  := ./cmd/api
BIN_DIR  := ./bin
CONN_STRING = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATION_DIRS = internal/db/migrations

run:
	go run $(CMD_DIR)/main.go

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

clean:
	rm -rf $(BIN_DIR)

tidy:
	go mod tidy

test:
	go test ./... -v

test/run:
	go test ./... -v -run $(name)

# Import database
importdb:
	docker exec -i postgres-db psql -U $(DB_USER) -d $(DB_NAME) < ./backupdb-master-golang.sql

# Export database
exportdb:
	docker exec -i postgres-db pg_dump -U $(DB_USER) -d $(DB_NAME) > ./backupdb-master-golang.sql

# Create a new migration (make migrate-create NAME=profiles)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)

# Run all pending migration (make migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# Rollback N migrations
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)

# Force migration version (use with caution example: make migrate-force VERSION=1) 
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)

# Drop everything (include schema migration)
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop

# Apply specific migration version (make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

# Generate sqlc code
sqlc:
	sqlc generate


.PHONY: run build clean tidy test importdb exportdb server migrate-create migrate-up migrate-down migrate-force migrate-drop migrate-goto migrate-down-n sqlc

