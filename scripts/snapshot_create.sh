#!/bin/bash

# E173 Gateway: Coordinated Snapshot Creator
# Creates Git commit + database backup with matching reference IDs

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
BACKUPS_DIR="$SNAPSHOTS_DIR/backups"
COMMITS_LOG="$SNAPSHOTS_DIR/commits.log"

# Ensure directories exist
mkdir -p "$BACKUPS_DIR"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-e173_gateway}"
DB_USER="${DB_USER:-e173_user}"

echo "🚀 Creating coordinated snapshot..."

# Check if we have uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  You have uncommitted changes. Please commit or stash them first."
    echo "   Use: git add . && git commit -m 'Your changes'"
    exit 1
fi

# Get commit message from argument
COMMIT_MESSAGE="$1"
if [ -z "$COMMIT_MESSAGE" ]; then
    echo "❌ Usage: $0 '<commit message>'"
    echo "   Example: $0 'Added customer management system'"
    exit 1
fi

# Create Git commit first
echo "📝 Creating Git commit..."
git add .
git commit -m "$COMMIT_MESSAGE" || {
    echo "ℹ️  No changes to commit, using existing HEAD commit"
}

# Get the commit ID
COMMIT_ID=$(git rev-parse --short HEAD)
echo "📋 Commit ID: $COMMIT_ID"

# Create database backup with commit ID in filename
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUPS_DIR/e173_gateway_${COMMIT_ID}.sql"
BACKUP_FILE_GZ="${BACKUP_FILE}.gz"

# Load database credentials from .env
if [ -f "$PROJECT_DIR/.env" ]; then
    export $(grep -v '^#' "$PROJECT_DIR/.env" | xargs)
    export PGPASSWORD="$DB_PASSWORD"
fi

echo "🗄️  Creating database backup..."
if pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"; then
    # Compress backup
    gzip "$BACKUP_FILE"
    BACKUP_SIZE=$(du -h "$BACKUP_FILE_GZ" | cut -f1)
    echo "✅ Database backup created: ${BACKUP_FILE_GZ} (${BACKUP_SIZE})"
else
    echo "❌ Database backup failed!"
    exit 1
fi

# Count records for metadata
RECORD_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT SUM(n_tup_ins + n_tup_upd) 
    FROM pg_stat_user_tables;" 2>/dev/null | tr -d ' ' || echo "0")

# Create metadata
METADATA_FILE="$BACKUPS_DIR/e173_gateway_${COMMIT_ID}.meta.json"
cat > "$METADATA_FILE" << EOF
{
  "commit_id": "$COMMIT_ID",
  "timestamp": "$(date -Iseconds)",
  "message": "$COMMIT_MESSAGE",
  "author": "$(git config user.name) <$(git config user.email)>",
  "database_backup": "e173_gateway_${COMMIT_ID}.sql.gz",
  "backup_size": "$BACKUP_SIZE",
  "records_count": $RECORD_COUNT,
  "git_remote": "$(git remote get-url origin 2>/dev/null || echo 'none')",
  "environment": "${ENVIRONMENT:-development}",
  "dependencies": {
    "go_version": "$(go version | cut -d' ' -f3 2>/dev/null || echo 'unknown')",
    "postgresql_version": "$(psql --version | cut -d' ' -f3 2>/dev/null || echo 'unknown')"
  },
  "checksums": {
    "database": "$(sha256sum "$BACKUP_FILE_GZ" | cut -d' ' -f1)"
  }
}
EOF

# Log the snapshot
echo "$COMMIT_ID|$(date -Iseconds)|$COMMIT_MESSAGE|$BACKUP_SIZE|$RECORD_COUNT" >> "$COMMITS_LOG"

# Display summary
echo ""
echo "🎉 Coordinated snapshot created successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Commit ID:       $COMMIT_ID"
echo "📝 Message:         $COMMIT_MESSAGE"
echo "🗄️  Database Backup: e173_gateway_${COMMIT_ID}.sql.gz"
echo "📊 Backup Size:     $BACKUP_SIZE"
echo "🔢 Records:         $RECORD_COUNT"
echo "🕒 Timestamp:       $(date)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🔧 Usage:"
echo "   • Restore snapshot:    ./scripts/snapshot_restore.sh $COMMIT_ID"
echo "   • Deploy to server:    ./scripts/deploy_fresh.sh $COMMIT_ID [server]"
echo "   • List all snapshots:  ./scripts/snapshot_list.sh"
echo ""

# Optional: Push to GitHub
if [ "$2" = "--push" ] || [ "$2" = "-p" ]; then
    echo "🚀 Pushing to GitHub..."
    git push origin master
    echo "✅ Pushed to GitHub repository"
fi
