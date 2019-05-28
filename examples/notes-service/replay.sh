#!/usr/bin/env bash
set -e

(cd ../.. && go build)

../../docker-test replay \
    --service notes-service \
    --compose ./docker-compose.yml \
    --sleep 3000 \
    --test 'cd ./tests && npm test'

