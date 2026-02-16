DB_USER=appuser2
DB_PASSWORD=app123
DB_HOST=localhost
DB_PORT=5432
DB_NAME=School-manajement
DB_SSLMODE=disable

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_PATH=./cmd/migrations

.PHONY: m-up m-down m-force m-create run

m-up:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) up

m-down:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) down

m-force:
	migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) force $(version)

m-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
 
run:
	go run ./cmd/main.go

test: 
	go test -v ./service/products