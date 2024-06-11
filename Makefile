.PHONY: prepare migrate test dev rollback stop

export GOOSE_DRIVER=postgres
export GOOSE_MIGRATION_DIR=./migrations
export GOOSE_DBSTRING=postgres://user:postgres@localhost:5432/pismo?sslmode=disable

prepare:
	go install github.com/pressly/goose/v3/cmd/goose@v3
	docker-compose -f docker-compose-db.yml up --build 

migrate:
	goose up

# Test
test:
	goose up
	go test ./tests/* -v
	docker-compose -f docker-compose-db.yml down

# Development build
dev:
	docker-compose -f docker-compose.yml up --build 

stop:
	docker-compose -f docker-compose.yml down

rollback:
	goose down
