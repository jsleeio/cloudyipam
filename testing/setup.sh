#!/bin/sh

set -e

cd "$SETUP_LOCATION"

ls -lR

./pg_setup.sh cloudyipam "$SETUP_CLOUDYIPAM_SQL_LOCATION" cloudyipam-test-user.sql
./pg_setup.sh terraform terraform-test-user.sql
