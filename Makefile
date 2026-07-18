DB_HOST ?= localhost
DB_URL = postgresql://root:secret@$(DB_HOST):5432/the_bank?sslmode=disable

network:
	docker network create the_bank_network
postgres:
	docker run --name the_bank_db --network the_bank_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
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

.PHONY: migrateup migratedown sqlc test testrace server mock postgres network
