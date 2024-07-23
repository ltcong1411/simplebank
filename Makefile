DB_URL=postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

start_postgres:
	docker start postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path=db/migration -database="$(DB_URL)" --verbose up

migrateup1:
	migrate -path=db/migration -database="$(DB_URL)" --verbose up 1

migratedown:
	migrate -path=db/migration -database="$(DB_URL)" --verbose down

migratedown1:
	migrate -path=db/migration -database="$(DB_URL)" --verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go simplebank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis --network bank-network -p 6379:6379 -d redis:7-alpine

start_redis:
	docker start redis

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock new_migration db_schema proto evans redis start_redis