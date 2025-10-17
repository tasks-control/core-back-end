# Database settings
USER=XXXXX
DBNAME=functions
PASSWORD=XXXX
HOST=localhost
PORT=5434

.PHONY: migrate

migrate:
	goose --dir ./migrations postgres "user=$(USER) dbname=$(DBNAME) password=$(PASSWORD) host=$(HOST) port=$(PORT)" up-by-one