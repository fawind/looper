#!/usr/bin/env bash
set -e

(cd ../../cli && go build)

../../cli/docker-test replay \
    --service notes-service \
    --compose ./docker-compose.yml \
    --test 'sleep 10 && (cd ./tests && npm test)'

