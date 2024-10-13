#!/bin/sh

set -e

echo "Starting db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the app"

exec "$@"