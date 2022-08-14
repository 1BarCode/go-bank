postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin123 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root go_bank

dropdb:
	docker exec -it postgres12 dropdb go_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:admin123@localhost:5432/go_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:admin123@localhost:5432/go_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:admin123@localhost:5432/go_bank?sslmode=disable" -verbose down		

migratedown1:
	migrate -path db/migration -database "postgresql://root:admin123@localhost:5432/go_bank?sslmode=disable" -verbose down 1		

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/1BarCode/go-bank/db/sqlc Store	

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock