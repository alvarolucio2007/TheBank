#!/bin/sh
set -e

echo "Aguardando banco de dados"
/app/wait-for.sh postgres:5432

echo "Rodando migrações"
/app/migrate -path /app/migration -database "$DB_SOURCE" up

echo "Inicializando app"
exec /app/main
