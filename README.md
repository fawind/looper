# docker-test

## CLI Usage

1. `cd cli`
2. Build the project: `go build`
3. Run service in record mode: `./docker-test record --service sender --compose ../playground/docker-compose.yml`
4. Run service in replay mode: `./docker-test replay --service sender --compose ../playground/docker-compose.yml`

