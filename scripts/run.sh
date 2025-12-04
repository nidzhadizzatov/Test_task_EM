#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")/.."

# Run database migrations
./scripts/migrate.sh

# Start the application
go run cmd/server/main.go