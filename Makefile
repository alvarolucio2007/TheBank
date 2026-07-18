DB_HOST ?= localhost
DB_URL = postgresql://root:secret@$(DB_HOST):5432/the_bank?sslmode=disable
NETWORK_NAME=the_bank_network
ENTRY_PORT=8080:8080
POSTGRES_PORT=5432:5432


#To build and create docker app FROM SCRATCH: Run docker_build->network->postgres->postgres_create_db->migrateup->docker_run
network:
	docker network create $(NETWORK_NAME) 
docker_build:
	docker build -t thebank:latest .
docker_run:
	docker run --name thebank --network $(NETWORK_NAME) -p $(ENTRY_PORT) -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@the_bank_db:5432/the_bank?sslmode=disable" thebank:latest
postgres:
	docker run --name the_bank_db --network the_bank_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
postgres_create_db:
	docker exec -it the_bank_db psql -U root -c "CREATE DATABASE the_bank;"
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
testrace:
	go test -v -cover -race ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/alvarolucio2007/TheBank/db/sqlc Store

.PHONY: migrateup migratedown sqlc test testrace server mock postgres network docker_build docker_run
