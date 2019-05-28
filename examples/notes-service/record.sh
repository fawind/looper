#!/usr/bin/env bash
set -e

(cd ../.. && make)

../../docker-test record \
    --service notes-service \
    --compose ./docker-compose.yml \
    --sleep 10000 \
    --test 'cd ./tests && npm test'

