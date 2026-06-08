#!/bin/sh

# Wait for database
if [ "$SQL_HOST" ]; then
    echo "Waiting for postgres..."
    while ! nc -z "$SQL_HOST" "$SQL_PORT"; do
        sleep 0.1
    done
    echo "PostgreSQL started"
fi

# Run migrations
/app/server migrate

exec "$@"
