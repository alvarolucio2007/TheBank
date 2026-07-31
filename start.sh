#!/bin/sh
set -e

source app/app.env

echo "Running migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" up

echo "Initializing app"
exec /app/main
