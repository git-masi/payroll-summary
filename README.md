# PAYROLL SUMMARY

## Getting started

Requirements:
- [Go](https://go.dev/) (see go.mod for required version)
- [Docker](https://www.docker.com/)
- `tern` [sql migration tool](https://github.com/jackc/tern)
- [SQLc](https://docs.sqlc.dev/en/stable/overview/install.html)

### Start docker

```sh
docker compose up -d
```

### Run migrations

Run all migrations:
```sh
tern migrate --migrations ./migrations --config ./migrations/tern.conf
```

See the status of the migrations:
```sh
tern status --migrations ./migrations --config ./migrations/tern.conf
```

### SQLc gen

Generate Go code from queries:
```sh
sqlc generate
```

## Run the program

The program can create workers, crews, payrolls, and earnings. You can configure everything except earnings by passing in arguments like so:

```sh
go run ./cmd/... -dsn "postgres://postgres:postgres@localhost/postgres" -num_workers 0 -num_crews 0 -should_create_payrolls false
```