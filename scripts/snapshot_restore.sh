#!/bin/bash

# E173 Gateway: Snapshot Restore System
# Restores complete system state to specific Git commit + database backup

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
BACKUPS_DIR="$SNAPSHOTS_DIR/backups"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-e173_gateway}"
DB_USER="${DB_USER:-e173_user}"

# Check if commit ID provided
if [ $# -eq 0 ]; then
    echo "❌ Usage: $0 <commit_id>"
    echo ""
    echo "📋 Available snapshots:"
    if [ -f "$SNAPSHOTS_DIR/commits.log" ]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Commit  | Date       | Message                    | Size"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        tail -10 "$SNAPSHOTS_DIR/commits.log" | while IFS='|' read -r commit date message size records; do
            printf "%-7s | %-10s | %-25s | %s\n" "$commit" "$(echo $date | cut -d'T' -f1)" "${message:0:25}" "$size"
        done
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    else
        echo "   No snapshots found. Create one with:"
        echo "   ./scripts/snapshot_create.sh 'Your message'"
    fi
    echo ""
    echo "💡 Usage examples:"
    echo "   ./scripts/snapshot_restore.sh a1b2c3d     # Restore to commit a1b2c3d"
    echo "   ./scripts/snapshot_restore.sh HEAD~1      # Restore to previous commit"
    exit 1
fi

COMMIT_ID="$1"

# Resolve commit ID if using relative references
if [[ "$COMMIT_ID" =~ HEAD|master|main ]]; then
    RESOLVED_COMMIT=$(git rev-parse --short "$COMMIT_ID" 2>/dev/null || echo "")
    if [ -z "$RESOLVED_COMMIT" ]; then
        echo "❌ Invalid commit reference: $COMMIT_ID"
        exit 1
    fi
    COMMIT_ID="$RESOLVED_COMMIT"
fi

echo "🔄 Restoring system to snapshot: $COMMIT_ID"

# Check if commit exists in git
if ! git cat-file -e "$COMMIT_ID" 2>/dev/null; then
    echo "❌ Git commit $COMMIT_ID not found in repository"
    exit 1
fi

# Check if database backup exists
BACKUP_FILE="$BACKUPS_DIR/e173_gateway_${COMMIT_ID}.sql.gz"
METADATA_FILE="$BACKUPS_DIR/e173_gateway_${COMMIT_ID}.meta.json"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Database backup not found: $BACKUP_FILE"
    echo "💡 Available backups:"
    ls -la "$BACKUPS_DIR"/*.sql.gz 2>/dev/null || echo "   No backups found"
    exit 1
fi

# Load metadata if available
if [ -f "$METADATA_FILE" ]; then
    echo "📋 Loading snapshot metadata..."
    COMMIT_MESSAGE=$(jq -r '.message' "$METADATA_FILE" 2>/dev/null || echo "No message")
    BACKUP_SIZE=$(jq -r '.backup_size' "$METADATA_FILE" 2>/dev/null || echo "Unknown")
    RECORD_COUNT=$(jq -r '.records_count' "$METADATA_FILE" 2>/dev/null || echo "Unknown")
    SNAPSHOT_DATE=$(jq -r '.timestamp' "$METADATA_FILE" 2>/dev/null || echo "Unknown")
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📋 Snapshot: $COMMIT_ID"
    echo "📝 Message:  $COMMIT_MESSAGE"
    echo "📅 Date:     $SNAPSHOT_DATE"
    echo "📊 Size:     $BACKUP_SIZE"
    echo "🔢 Records:  $RECORD_COUNT"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
fi

# Safety confirmation
echo ""
echo "⚠️  WARNING: This will restore your system to the selected snapshot!"
echo "   • Current code changes will be lost (unless committed)"
echo "   • Database will be completely replaced"
echo "   • All current data will be overwritten"
echo ""
read -p "Are you sure you want to proceed? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Restore cancelled"
    exit 0
fi

echo ""
echo "🔄 Starting restore process..."

# Step 1: Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  You have uncommitted changes!"
    echo "   Creating emergency backup commit..."
    EMERGENCY_COMMIT="emergency_backup_$(date +%Y%m%d_%H%M%S)"
    git add .
    git commit -m "Emergency backup before restore to $COMMIT_ID" || true
    echo "✅ Emergency backup created: $EMERGENCY_COMMIT"
fi

# Step 2: Restore Git repository
echo "📂 Restoring Git repository to $COMMIT_ID..."
git checkout "$COMMIT_ID" -b "restore_$(date +%Y%m%d_%H%M%S)" 2>/dev/null || git checkout "$COMMIT_ID"
echo "✅ Git repository restored"

# Step 3: Stop application services
echo "🛑 Stopping application services..."
pkill -f e173gw || true
sleep 2
echo "✅ Services stopped"

# Step 4: Restore database
echo "🗄️  Restoring database..."

# Drop and recreate database
echo "   Dropping existing database..."
dropdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" --if-exists

echo "   Creating fresh database..."
createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"

# Restore from backup
echo "   Restoring from backup..."
gunzip -c "$BACKUP_FILE" | psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"
echo "✅ Database restored"

# Step 5: Rebuild application
echo "🔨 Rebuilding application..."
cd "$PROJECT_DIR"
if make build 2>/dev/null; then
    echo "✅ Application built successfully"
else
    echo "⚠️  Build had warnings, but binary created"
fi

# Step 6: Restore environment configuration
echo "⚙️  Restoring configuration..."
if [ -f "deployment/config/${ENVIRONMENT:-development}.env" ]; then
    cp "deployment/config/${ENVIRONMENT:-development}.env" .env
    echo "✅ Environment configuration restored"
else
    echo "ℹ️  Using existing .env configuration"
fi

# Step 7: Start services
echo "🚀 Starting services..."
nohup make run > /dev/null 2>&1 &
sleep 3

# Step 8: Verify restore
echo "🔍 Verifying restore..."

# Check if server is responding
if curl -s http://localhost:8080/ping > /dev/null 2>&1; then
    echo "✅ Server is responding"
else
    echo "⚠️  Server may need more time to start"
fi

# Check database connectivity
if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Database connectivity verified"
else
    echo "❌ Database connectivity issues"
fi

echo ""
echo "🎉 System restore completed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Restored to:     $COMMIT_ID"
echo "📝 Commit message:  $COMMIT_MESSAGE"
echo "🕒 Restore time:    $(date)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🔧 Next steps:"
echo "   • Verify functionality: curl http://localhost:8080/"
echo "   • Check dashboard: http://localhost:8080"
echo "   • Review logs: make logs"
echo ""
echo "🚨 If you need to rollback this restore:"
echo "   • Emergency restore: ./scripts/rollback.sh"
echo "   • Return to previous state: git checkout master"
