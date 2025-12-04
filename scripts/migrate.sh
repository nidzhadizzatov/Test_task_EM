#!/bin/bash

# This script is used to run database migrations.

set -e

# Define the database connection parameters
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="your_username"
DB_PASSWORD="your_password"
DB_NAME="your_database"

# Run the migrations
echo "Running migrations..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f ../db/migrations/0001_init.up.sql

echo "Migrations completed successfully."