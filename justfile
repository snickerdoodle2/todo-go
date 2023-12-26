set dotenv-load

watch:
    watchexec -i "templates/**" -r just run

build:
    go build

run:
    go run .

migrate:
    sqlx migrate run
