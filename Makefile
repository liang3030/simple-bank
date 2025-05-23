DB_URL=postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres --network bank-network -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

aws-migrateup:
	migrate -path db/migration -database "postgresql://root:${secret}@simple-bank.cnyo21on7ezz.eu-central-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "${DB_URL}" -verbose down

migratedown1:
	migrate -path db/migration -database "${DB_URL}" -verbose down 1

migrateup1:
	migrate -path db/migration -database "${DB_URL}" -verbose up 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockDB -destination db/mock/store.go github.com/liang3030/simple-bank/db/sqlc IStore
	mockgen -package mockWK -destination worker/mock/distributor.go github.com/liang3030/simple-bank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto \
       --go_out=pb/ --go_opt=paths=source_relative \
       --go-grpc_out=pb/ --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=pb/ --grpc-gateway_opt=paths=source_relative \
			 --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank\
       proto/*.proto
			 statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:8.0-M03-alpine

.PHONY: postgres createdb dropdb test migrateup migratedown sqlc test server mock aws-migrateup proto evans redis new_migration