#!/bin/sh
set -e

echo "Running migrations"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_URL" up

echo "Initializing app"
exec /app/main
