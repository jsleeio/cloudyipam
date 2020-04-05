#!/bin/sh

set -e -u
top="$(git rev-parse --show-toplevel)"
cd "$top"

docker-compose \
  --file testing/compose-test.yaml \
  down \
    --volumes \
    --remove-orphans

find testing -name .terraform -type d -print0 | xargs rm -r
find testing -name terraform.d -type d -print0 | xargs rm -r
find testing -name terraform-provider-cloudyipam -delete
