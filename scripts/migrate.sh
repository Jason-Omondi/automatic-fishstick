#!/bin/bash
# Helper script to run database migrations
# Usage: ./migrate.sh up
# Usage: ./migrate.sh down

cd "$(dirname "$0")/.."

DB_TYPE=${DB_TYPE:-mysql}
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-ecomgo}

if [ "$DB_TYPE" = "postgres" ]; then
  DB_PORT=${DB_PORT:-5432}
fi

echo "Running migrations for $DB_TYPE database..."

migrate -path internal/migrations/sql \
  -database "$(get_db_url)" \
  $1

function get_db_url() {
  if [ "$DB_TYPE" = "mysql" ]; then
    echo "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME"
  elif [ "$DB_TYPE" = "postgres" ]; then
    echo "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
  fi
}
