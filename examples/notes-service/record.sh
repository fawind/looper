#!/usr/bin/env bash
set -e

(cd ../.. && go build)
(cd ./tests && npm install)

../../docker-test record \
    --service notes-service \
    --compose ./docker-compose.yml \
    --sleep 3000 \
    --test 'cd ./tests && npm test'

