#!/bin/sh

set -e -u
top="$(git rev-parse --show-toplevel)"
cd "$top"

docker run \
  --interactive \
  --tty \
  --rm \
  -p 5432:5432 \
  --env POSTGRES_PASSWORD="$RANDOM" \
  --volume $(pwd)/database/cloudyipam.sql:/docker-entrypoint-initdb.d/01_cloudyipam.sql \
  --volume $(pwd)/testing/setup/cloudyipam-test-user.sql:/docker-entrypoint-initdb.d/02_testuser.sql \
  postgres:12-alpine \
  "$@"
