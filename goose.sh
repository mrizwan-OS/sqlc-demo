#!/bin/bash
# Goose migration helper

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://u0_a283@localhost:5432/mydb?sslmode=disable"
export GOOSE_MIGRATION_DIR=migrations

# Pass all arguments to goose
goose "$@"
