#!/bin/sh

set -e

# Lecture 50: run migration in golang file
# comment out the migration lines in start.sh
# echo "Starting db migration"

# in docker file, the app.env file is copied to the /app directory
# source /app/app.env
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the app"

exec "$@"