postgres:
	docker start postgres
# 	docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=sumanth123 -d postgres

sqlshell:
	docker exec -it postgres psql -U postgres -d xbank

# Executes a command inside postgres container shell for Creating the New Database 
createdb:
	docker exec -it postgresql createdb --username=root --owner=root xbank

# Executes a command inside postgres container for droping the Database 
dropdb:
	docker exec -it postgres dropdb --username=postgres xbank

# Run this migrate cli command after creating the xbank db
migrateup:
	migrate -path db/migration -database "postgresql://root:sumanth123@localhost:5432/xbank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:sumanth123@localhost:5432/xbank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server
