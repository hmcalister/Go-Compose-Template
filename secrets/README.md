# Secrets

House your secrets here, such as the password to the postgres database. To compile this project you MUST create a file `postgres-password.secret` in this directory (ideally with a secure password). See `compose.yaml` on how this secret is passed to the database container (to set the password) and app container (to use the password when making a connection to the database).