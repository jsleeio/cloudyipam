#!/bin/sh

set -e -u
top="$(git rev-parse --show-toplevel)"
cd "$top"

docker-compose --file testing/compose-test.yaml down \
  --volumes --remove-orphans

docker-compose --file testing/compose-test.yaml up \
  --remove-orphans --renew-anon-volumes --force-recreate
