#!/bin/sh
set -e

echo "Running migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" up

echo "Initializing app"
exec /app/main
