#!/bin/bash

# E173 Gateway Database Backup Script
# Creates timestamped PostgreSQL backups

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-e173_gateway}"
DB_USER="${DB_USER:-e173_user}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"

# Load database credentials from .env
if [ -f "$(dirname "$0")/../.env" ]; then
    export $(grep -v '^#' "$(dirname "$0")/../.env" | xargs)
    export PGPASSWORD="$DB_PASSWORD"
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Generate timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/e173_gateway_backup_$TIMESTAMP.sql"

echo "🔄 Starting database backup..."
echo "📅 Timestamp: $TIMESTAMP"
echo "🗄️  Database: $DB_NAME"
echo "📁 Backup file: $BACKUP_FILE"

# Create backup
echo "📊 Creating database backup..."
if pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -w > "$BACKUP_FILE"; then
    echo "✅ Backup completed successfully!"
    echo "📊 Backup size: $(du -h "$BACKUP_FILE" | cut -f1)"
    
    # Compress backup
    gzip "$BACKUP_FILE"
    echo "🗜️  Backup compressed: ${BACKUP_FILE}.gz"
    
    # Keep only last 7 days of backups
    find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete
    echo "🧹 Old backups cleaned up (keeping last 7 days)"
    
else
    echo "❌ Backup failed!"
    exit 1
fi

echo "✅ Database backup completed: ${BACKUP_FILE}.gz"
