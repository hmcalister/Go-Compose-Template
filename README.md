# Go Compose Template

An architecture template for an app using Golang. The server uses the Go standard router to handle requests. [Postgres](https://www.postgresql.org/) is used for the database. [SQLC](https://sqlc.dev/) generates type-safe queries to the database. Podman and Podman Compose allow for simple deployment and the ability to scale later by developing with containers in mind.

## Running This Project

This walkthrough is designed for Linux systems. Windows systems can follow the general thread but may need to substitute some commands.

### Secrets and Set Up

- Create the file `.env` to hold secrets for the deployment. The following key-values must be present:
    - `POSTGRES_USER`
    - `POSTGRES_PASSWORD`
    - `POSTGRES_DB`

### SQL and Generation

- Define your table schemas in `sqlc/schema.sql`. 
- Define all queries under `sqlc/queries.sql`.
- If needed, make any changes to `sqlc/sqlc.yaml`.
- Then `sqlc generate -f sqlc/sqlc.yaml`

### Development

- Run `podman compose up --build` to build the images and start the stack.
  - Note the `--build` flag! The Go source is compiled into the image, so a source change only takes effect on a rebuild.
- Run `podman compose down` to stop. 
  - Add `-v` to also drop the volumes. BEWARE! This will drop your database.
