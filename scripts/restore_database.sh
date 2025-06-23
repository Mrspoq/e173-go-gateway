#!/bin/bash

# E173 Gateway Database Restore Script
# Restores PostgreSQL database from backup

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-e173_gateway}"
DB_USER="${DB_USER:-e173_user}"

# Check if backup file provided
if [ $# -eq 0 ]; then
    echo "❌ Usage: $0 <backup_file.sql.gz>"
    echo "📁 Available backups:"
    ls -la ./backups/*.sql.gz 2>/dev/null || echo "   No backups found in ./backups/"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "🔄 Starting database restore..."
echo "📁 Backup file: $BACKUP_FILE"
echo "🗄️  Target database: $DB_NAME"

# Confirm restore (destructive operation)
echo "⚠️  WARNING: This will OVERWRITE the current database!"
read -p "Are you sure? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Restore cancelled"
    exit 0
fi

# Drop existing database and recreate
echo "🗑️  Dropping existing database..."
dropdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" --if-exists

echo "🏗️  Creating fresh database..."
createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"

# Restore from backup
echo "📥 Restoring from backup..."
if [[ "$BACKUP_FILE" == *.gz ]]; then
    gunzip -c "$BACKUP_FILE" | psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"
else
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" < "$BACKUP_FILE"
fi

echo "✅ Database restore completed successfully!"
echo "🔧 You may need to restart your application to reconnect to the database"
