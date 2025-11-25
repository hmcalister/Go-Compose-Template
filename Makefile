.PHONY: build run sqlcGenerate appBuild podmanBuild podmanClean podmanRun

include secrets/.env
export

build: sqlcGenerate appBuild podmanBuild

run: build podmanRun

sqlcGenerate:
	sqlc generate -f sqlc/sqlc.yaml

appBuild:
	go build -o build/ ./cmd/main.go

podmanBuild:
	podman compose build 

podmanClean:
	podman compose down -v

podmanRun: podmanBuild
	podman compose up --force-recreate

