#!/bin/sh

# postgres container entrypoint runs init scripts before enabling a
# tcp listener, so must connect to the socket instead

/usr/local/bin/cloudyipam initialize \
  --db-host=/var/run/postgresql \
  --db-user=postgres \
  --db-tls=false
