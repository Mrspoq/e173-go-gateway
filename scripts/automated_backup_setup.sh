#!/bin/bash

# E173 Gateway Automated Backup Setup
# Sets up daily automated database backups via cron

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "🔧 Setting up automated database backups..."

# Create backups directory
mkdir -p "$PROJECT_DIR/backups"
echo "📁 Created backups directory: $PROJECT_DIR/backups"

# Create cron job for daily backups at 2 AM
CRON_JOB="0 2 * * * cd $PROJECT_DIR && ./scripts/backup_database.sh"

# Check if cron job already exists
if crontab -l 2>/dev/null | grep -q "$PROJECT_DIR/scripts/backup_database.sh"; then
    echo "⚠️  Automated backup already configured"
else
    # Add cron job
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo "✅ Daily backup scheduled at 2:00 AM"
fi

# Create backup configuration file
cat > "$PROJECT_DIR/.backup.env" << EOF
# E173 Gateway Backup Configuration
# Source this file before running backup scripts

export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=e173_gateway
export DB_USER=e173_user
export BACKUP_DIR=$PROJECT_DIR/backups

# Backup retention (days)
export BACKUP_RETENTION_DAYS=7
EOF

echo "⚙️  Created backup configuration: $PROJECT_DIR/.backup.env"

# Test backup (optional)
echo ""
echo "🧪 Would you like to test the backup now? (y/n)"
read -r test_backup

if [[ "$test_backup" =~ ^[Yy]$ ]]; then
    echo "🔄 Running test backup..."
    source "$PROJECT_DIR/.backup.env"
    "$SCRIPT_DIR/backup_database.sh"
fi

echo ""
echo "✅ Automated backup setup complete!"
echo ""
echo "📋 Backup Summary:"
echo "   • Daily backups at 2:00 AM"
echo "   • Backups stored in: $PROJECT_DIR/backups"
echo "   • Retention: 7 days"
echo ""
echo "🔧 Manual Commands:"
echo "   • Create backup: ./scripts/backup_database.sh"
echo "   • Restore backup: ./scripts/restore_database.sh <backup_file>"
echo "   • List backups: ls -la ./backups/"
