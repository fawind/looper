#!/usr/bin/env bash
set -e

(cd ../../cli && go build)

../../cli/docker-test replay \
    --service notes-service \
    --compose ./docker-compose.yml \
    --sleep 10 \
    --test 'cd ./tests && npm test'

