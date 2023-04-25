DB_NAME=golangsimplebank
DB_URL=postgresql://root:secret@127.0.0.1:5432/golangsimplebank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres12-alpine-docker --network database -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12.13-alpine3.17

createdb:
	docker exec -it postgres12-alpine-docker createdb ${DB_NAME} --username=root --owner=root

dropdb:
	docker exec -it postgres12-alpine-docker dropdb -U root ${DB_NAME}

freshdb: dropalldb
	migrate -path ./migration -database "$(DB_URL)" -verbose up

dropalldb:
	migrate -path ./migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go clean -testcache && go test -v -cover -short ./...

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/claytten/golang-simplebank/internal/db/sqlc Store

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
	evans --host localhost --port 8081 -r repl

.PHONY: createdb dropdb freshdb dropalldb sqlc server test mock proto evans