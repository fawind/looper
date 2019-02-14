#!/usr/bin/env bash
set -e

(cd ../../cli && go build)

../../cli/docker-test record \
    --service notes-service \
    --compose ./docker-compose.yml \
    --test 'sleep 30 && (cd ./tests && npm test)'

