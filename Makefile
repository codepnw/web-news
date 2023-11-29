include .env

postgresinit:
	docker run --name ${PG_DB} -p ${PG_PORT}:5432 -e POSTGRES_USER=${PG_USER} -e POSTGRES_PASSWORD=${PG_PASS} -d postgres:15.4

postgres:
	docker exec -it ${PG_DB} psql

createdb:
	docker exec -it ${PG_DB} createdb --username=${PG_USER} --owner=${PG_USER} ${PG_DB}

dropdb:
	docker exec -it ${PG_DB} dropdb ${PG_DB}

run:
	air -c .air.toml

migrate:
	go run cmd/main.go -migrate=true