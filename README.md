# Go Compose Template

An architecture template for an app using Golang. The server uses the Chi router to handle requests. Postgres is used for the database. SQLC generates type-safe queries to the database. Docker Compose allows for a simple deployment and the ability to scale later by developing with containers in mind.

### Running This Project

This walkthrough is designed for Linux systems. Windows systems can follow the general thread but may need to substitute some commands.

- First, create the file `secrets/postgres-password.secret` and add a password to it. This password is used by the database and app containers to improve security, but the secret file itself is not tracked.
- Run `make build` (or, run each of `make sqlc`, `make app`, `make docker` separately, see below). This will download all Go dependencies and docker containers required to run the project.
- Run `make run` to start the project with Docker Compose. 
    - On the first run, the database container will create a new database. This may not finish before the app container complains. If the app container is killed during the first `make run`, simply wait for the database to finish (a log message will inform you that the "database system is ready to accept connections") before stopping the process and calling `make run` again. If the app container still does not start up, try altering the healthcheck in `compose.yaml` or removing it entirely.
- Once the containers have started up, the app container should log to stdout with the message "ready to serve".
- With a new terminal, run `curl 127.0.0.1:8080/newAuthor/author`. The app container should log `msg="new author request" authorName=author`.
- In the new terminal, run `curl 127.0.0.1:8080/allAuthors`. You should receive a response `1: author`. 

Thanks to the persistent volume, restarting the Docker Compose session should retain the data (i.e. curling the `allAuthors` end point should return the authors made previously).

Check our the Go code under `internal/app` to see how the app logic works, including the log statements seen in this example. See the SQLC schemas and queries under `sqlc/` for information on the database schema. Finally, `compose.yaml` shows how these containers are built, managed, and communicate.

### Make

This template provides a Makefile to help set up and manage the various systems. Use `make -B <target>` to force the command to be run if `make` refuses to build something.

- `make sqlc` generates SQLC for use in the Go project. Schemas and queries are written in the respective files under `sqlc`. Generated code is placed under `internal/app/database`. Note that the `schema.sql` file is used to generate the database in Postgres.
- `make app` builds the application logic (a Go project under `internal/app`) and places the executable in `build`. This can be changed to use Docker to build the application as well if need be, and is essential if the host operating system is not Linux based.
- `make docker` builds the app container (specified by the Dockerfile under `build`, copying in the compiled Go project) and creates the database container. If the database container has not been built before, the file `sqlc/schema.sql` that is copied in will create a new set of tables. Note that the database is stored in a persistent volume managed by Docker (volume `database-data`) and hence restarting the container will ensure the database persists. Since `sqlc/schema.sql` is run on every start up, adding the relevant `DROP TABLE` statements to the start of the schema will allow for a fresh database with each container restart, and is useful for testing and development. A similar result can be achieved by removing the volume with `docker volume rm <volume name>`
    - `make -i dockerClean` allows for the explicit removal of the docker containers and volumes --- BEWARE! This will drop your database and containers!
- `make build` runs the above commands in order, preparing the application to be run.
- `make run` first runs `make build` then `docker compose up` to start the application --- this effectively replaces `go run .` in your debugging step. 

### Secrets

The `secrets` directory allows for secrets (passwords, API keys, and so on) to be placed within the project without fear of being tracked with git due to the `.gitignore`. To compile this project the file `postgres-password.secret` must be present in this directory. 

### Logging and Debugging

This template uses `slog` to log events in the application. By default the log level is set to `INFO` and logs are written in JSON format to a file that is mounted to the host file-system at `./logs/log`. 

An environment variable `DEBUG` is checked when the container starts up, and can be set in the `compose.yaml`. If this variable is set, logging is instead set to level `DEBUG` and logs are written as text directly to `stdout`.

### TODO:

- Add SSL to the Postgres database.
- Prevent outside requests to Postgres database --- look into Docker Compose networking.
