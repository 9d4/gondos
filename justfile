set export
set dotenv-load := true

CGO_ENABLED := "0"
TZ := "UTC"

default:
    just --list

build *build-flags:
    go build -buildvcs {{build-flags}} .

run *run-args:
    go run -buildvcs . {{run-args}}

gen:
    go generate ./...

migrate *args:
    go run -tags 'mysql' \
        github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0 \
        -source=file://db/migrations/ \
        -database=mysql://"$DSN" \
        {{args}}

gqlgen *args:
    go run github.com/99designs/gqlgen {{args}}
