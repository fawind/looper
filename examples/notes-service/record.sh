#!/usr/bin/env bash
set -e

(cd ../../cli && go build)
(cd ./tests && npm install)

../../cli/docker-test record \
    --service notes-service \
    --compose ./docker-compose.yml \
    --sleep 30 \
    --test 'cd ./tests && npm test'

