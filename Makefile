DB_NAME=golang-simplebank
DB_URL=postgresql://<username>:<password>@<url_database>:<post_database>/${DB_NAME}?sslmode=disable

createdb:
	docker exec -it postgres12-alpine createdb ${DB_NAME} --username=zerocool --owner=zerocool

dropdb:
	docker exec -it postgres12-alpine dropdb -U zerocool ${DB_NAME}

freshdb: dropalldb
	migrate -path ./migration -database "$(DB_URL)" -verbose up

dropalldb:
	migrate -path ./migration -database "$(DB_URL)" -verbose down

.PHONY: createdb dropdb freshdb dropalldb