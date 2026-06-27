#!/bin/bash

echo "🧹 Cleaning database..."

# Stop PostgreSQL
pg_ctl -D $PREFIX/var/lib/postgresql stop

# Remove database files
rm -rf $PREFIX/var/lib/postgresql

# Reinitialize
initdb $PREFIX/var/lib/postgresql

# Start PostgreSQL
pg_ctl -D $PREFIX/var/lib/postgresql start

# Wait for PostgreSQL to start
sleep 2

# Create database
createdb mydb

# Run migrations
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://u0_a283@localhost:5432/mydb?sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations
goose up

echo "✅ Database cleaned and migrations applied!"
