# TODO app

Simple TODO app made with `go`, `chi`, `pgx` and `htmx`.

App runs on http://localhost:8080.

## Installation
### Prerequisites
- [go](https://go.dev/dl/)
- [just](https://github.com/casey/just) - command runner
- docker
- `sqlx-cli` - `cargo install sqlx-cli`

### `.env` file
Create `.env` file containing:
``` sh
DATABASE_ADDRESS=localhost:5432
POSTGRES_PASSWORD=<YOUR PASSWORD>
POSTGRES_DB=quotes
DATABASE_URL=postgresql://postgres:<YOUR PASSWORD>@localhost:5432/quotes
```

### Migrating database
To setup the database, simply run
``` sh
docker-compose up -d
just migrate
```

## Running server
`just run`
