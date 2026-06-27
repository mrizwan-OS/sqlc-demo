#!/bin/bash
# Migration management script

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://u0_a283@localhost:5432/mydb?sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations

case "$1" in
    up)
        echo "🔄 Applying all pending migrations..."
        goose up
        ;;
    down)
        echo "⏪ Rolling back one migration..."
        goose down
        ;;
    reset)
        echo "🔄 Resetting database..."
        goose reset
        ;;
    status)
        echo "📊 Migration status:"
        goose status
        ;;
    create)
        if [ -z "$2" ]; then
            echo "❌ Please provide a migration name"
            echo "Usage: ./migrate.sh create add_something"
            exit 1
        fi
        echo "📝 Creating migration: $2"
        goose create "$2" sql
        ;;
    redo)
        echo "🔄 Redoing last migration..."
        goose redo
        ;;
    *)
        echo "Usage: ./migrate.sh {up|down|reset|status|create <name>|redo}"
        exit 1
        ;;
esac
