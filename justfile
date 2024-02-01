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

serve *run-args:
    go run -buildvcs . serve {{run-args}}

gen:
    go generate ./...

migrate *args:
    go run -tags 'mysql' \
        github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0 \
        -source=file://db/migrations/ \
        -database=mysql://"$DSN" \
        {{args}}

oapigen *args:
    go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest {{args}}

oapigen-gen *args:
    just oapigen -generate=types -package=api -o=api/types.gen.go {{args}} api/description/api.yml
    just oapigen -generate=chi-server -package=api -o=api/server.gen.go {{args}} api/description/api.yml


jet *args:
    go run github.com/go-jet/jet/v2/cmd/jet@latest \
        -source=mysql \
        -path=./jetgen \
        -dbname=asd \
        -dsn="$DSN" \
        {{args}}

gqlgen *args:
    go run github.com/99designs/gqlgen {{args}}
