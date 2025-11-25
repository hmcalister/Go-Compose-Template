# Go Compose Template

An architecture template for an app using Golang. The server uses the Go standard router to handle requests. [Postgres](https://www.postgresql.org/) is used for the database. [SQLC](https://sqlc.dev/) generates type-safe queries to the database. Podman and Podman Compose allow for simple deployment and the ability to scale later by developing with containers in mind.

## Running This Project

This walkthrough is designed for Linux systems. Windows systems can follow the general thread but may need to substitute some commands.

### Secrets and Set Up

- Create the file `secrets/.env` to hold secrets for the deployment. The following key-values must be present:
    - `POSTGRES_USER`
    - `POSTGRES_PASSWORD`
    - `POSTGRES_DB`

### SQL and Generation

- Define your table schemas in `sqlc/schema.sql`. 
- Define all queries under `sqlc/queries.sql`.
- If needed, make any changes to `sqlc/sqlc.yaml`.

### Initialization

- Run `make build` to download all Go dependencies and containers.

### Run

- Run `make run` to start the project with Podman Compose. 
    - On the first run, the database container will create a new database. This may not finish before the app container complains. If the app container is killed during the first `make run`, simply wait for the database to finish (a log message will inform you that the "database system is ready to accept connections") before stopping the process and calling `make run` again. If the app container still does not start up, try altering the healthcheck in `compose.yaml` or removing it entirely.

Persistent volumes are used to ensure data is not lost if the database container is destroyed.

## Make

This template provides a Makefile to help set up and manage the various systems.

- `make sqlc` generates SQLC for use in the Go project. Schemas and queries are written in the respective files under `sqlc`. Generated code is placed under `internal/app/database`. Note that the `schema.sql` file is used to generate the database in Postgres.
- `make app` builds the application logic (a Go project under `internal/app`) and places the executable in `build`.
- `make podmanBuild` builds the app container and creates the database container. If the database container has not been built before, the file `sqlc/schema.sql` that is copied in will create a new set of tables. Note that the database is stored in a persistent volume managed by Docker (volume `database-data`) and hence restarting the container will ensure the database persists. Since `sqlc/schema.sql` is run on every start up, adding the relevant `DROP TABLE` statements to the start of the schema will allow for a fresh database with each container restart, and is useful for testing and development. A similar result can be achieved by removing the volume with `docker volume rm <volume name>`
    - `make -i dockerClean` allows for the explicit removal of the docker containers and volumes --- BEWARE! This will drop your database and containers!
- `make build` runs the above commands in order, preparing the application to be run.
- `make run` first runs `make build` then `docker compose up` to start the application --- this effectively replaces `go run .` in your debugging step. 

## Logging and Debugging

This template uses `slog` to log events in the application. By default the log level is set to `INFO` and logs are written in JSON format to a file that is mounted to the host file-system at `./logs/log`. Note that the `./logs/log` directory must be writeable by the container! Therefore, run `touch ./logs/log && chmod 777 ./logs/log` to allow writing of the log file.

An environment variable `DEBUG` is checked when the container starts up, and can be set in the `compose.yaml`. If this variable is set, logging is instead set to level `DEBUG` and logs are written as text directly to `stdout`.
