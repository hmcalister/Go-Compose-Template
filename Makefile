build: sqlc app docker

run: build
	docker compose up
include secrets/.env
export

sqlc:
	sqlc generate -f sqlc/sqlc.yaml

app:
	go build -o build/ ./internal/app

docker: 
	docker compose build

dockerClean:
	docker container ls -a --filter "name=^go-compose-template" --format "{{.ID}}" | xargs docker container rm -f
	docker volume ls --filter "name=^go-compose-template" --format "{{.Name}}" | xargs docker volume rm -f

dockerRun: dockerBuild
	docker compose up

