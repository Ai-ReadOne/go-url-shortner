# Migrations
golang-migrate https://github.com/golang-migrate/migrate for database migrations. 

## Setup
Set up the CLI of golang-migrate to get started using any of the methods on the website. For build setup, we'll use this:
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

if not already added, you can temporarily add $GOPATH/bin to your PATH:
```
export PATH=$PATH:$(go env GOPATH)/bin
```

You can now run migration commands

## Applying migrations
You need to specify the database and point at the migration files. For example:

```
migrate -database YOUR_DATABASE_URL -path PATH_TO_YOUR_MIGRATIONS up
```

or to rollback
```
migrate -database YOUR_DATABASE_URL -path PATH_TO_YOUR_MIGRATIONS down
```

You can also specify up [N] or down [N] to indicate the number of migrations to apply.

For example, to setup

```
export POSTGRES_URL='postgres://postgres:password@localhost/db_name?sslmode=disable'
```

migrate up:

```
migrate -database ${POSTGRES_URL} -path internal/db/migrations up
```

### Note: if you do not wish to use the golang-migrate package, you can copy and past the migrations in [migrations folder](../internal/database/migrations) in your postgres terminal or PGAdmin