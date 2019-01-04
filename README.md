# docker-test

## CLI Usage

1. Build the project:
```bash
cd cli
go build
```
2. Run service in record mode:
```bash
./docker-test record --service sender --compose ../playground/docker-compose.yml --test 'echo <TESTCMD>'
```
3. Run service in replay mode:
```bash
./docker-test replay --service sender --compose ../playground/docker-compose.yml --test 'echo <TESTCMD>'
```
