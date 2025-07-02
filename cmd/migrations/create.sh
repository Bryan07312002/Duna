#!/bin/sh

# this path should not have "/" at the end
MIGRATIONS_DIR="./internal/database/migrations"

# Check if a name argument was provided
if [ $# -eq 0 ]; then
    echo "Error: Please provide a name for the folder"
    echo "Usage: $0 <name>"
    exit 1
fi

# Get current timestamp in YYYYMMDD_HHMMSS format
timestamp=$(date +"%s")

# Create folder name by combining timestamp and provided name
folder_name="${MIGRATIONS_DIR}/${timestamp}-$1"

# Create the folder
mkdir -p "$folder_name"

# Create the up.sql and down.sql files inside the folder
touch "$folder_name/up.sql" "$folder_name/down.sql"

echo "Created folder '$folder_name' with up.sql and down.sql files"
