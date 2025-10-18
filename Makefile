# Database settings
USER=user
DBNAME=core-back-db
PASSWORD=1234
HOST=localhost
PORT=5432

.PHONY: migrate docker-up docker-down docker-restart docker-logs db-status

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  migrate - Run migrations"
	@echo "  migrate-down - Rollback one migration"
	@echo "  migrate-status - Show migration status"
	@echo "  "
	@echo "  docker-up - Start docker containers"
	@echo "  docker-down - Stop docker containers"
	@echo "  docker-restart - Restart docker containers"
	@echo "  docker-logs - Show docker logs"
	@echo "  db-status - Show database status"
	@echo "  setup - Full setup: start docker and run migrations"
	@echo "  clean - Clean everything (including volumes)"

# Run migrations
migrate:
	goose --dir ./migrations postgres "user=$(USER) dbname=$(DBNAME) password=$(PASSWORD) host=$(HOST) port=$(PORT) sslmode=disable" up-by-one

# Migrate down (rollback one migration)
migrate-down:
	goose --dir ./migrations postgres "user=$(USER) dbname=$(DBNAME) password=$(PASSWORD) host=$(HOST) port=$(PORT) sslmode=disable" down

# Show migration status
migrate-status:
	goose --dir ./migrations postgres "user=$(USER) dbname=$(DBNAME) password=$(PASSWORD) host=$(HOST) port=$(PORT) sslmode=disable" status

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-restart:
	docker-compose restart

docker-logs:
	docker-compose logs -f postgres

db-status:
	docker-compose ps postgres

# Full setup: start docker and run migrations
setup: docker-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@make migrate
	@echo "✅ Setup complete! Database is ready."
	@echo "📊 pgAdmin available at: http://localhost:5050"
	@echo "   Email: admin@admin.com"
	@echo "   Password: admin"

# Clean everything (including volumes)
clean:
	docker-compose down -v