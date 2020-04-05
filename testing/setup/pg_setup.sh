#!/bin/sh

set -e -u

db="$1"
shift
# all remaining args are treated as sql text

psql \
  --set=ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname=postgres \
  --command="DROP DATABASE IF EXISTS $db;"

psql \
  --set=ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname=postgres \
  --command="CREATE DATABASE $db;"

for sql in "$@" ; do
  psql \
    --set=ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname="$db" \
    --file="$sql"
done
