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

To run all migrations:
```sh
tern migrate --migrations ./migrations --config ./migrations/tern.conf
```

To see the status of the migrations:
```sh
tern status --migrations ./migrations --config ./migrations/tern.conf
```