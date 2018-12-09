# docker-test

## Usage

1. Build the project: `go build`
2. Run service in record mode: `./docker-test record --service sender --compose ./playground/docker-compose.yml`
3. Run service in replay mode: `./docker-test replay --service sender --compose ./playground/docker-compose.yml`

