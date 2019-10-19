#!/usr/bin/env bash
set -e

(cd ../.. && make)

../../looper record \
    --service notes-service \
    --compose ./docker-compose.yml \
    --test 'cd ./tests && npm test'

