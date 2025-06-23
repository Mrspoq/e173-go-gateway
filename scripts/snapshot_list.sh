#!/bin/bash

# E173 Gateway: Snapshot List Utility
# Lists all available coordinated snapshots

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
COMMITS_LOG="$SNAPSHOTS_DIR/commits.log"
BACKUPS_DIR="$SNAPSHOTS_DIR/backups"

echo "📋 E173 Gateway Snapshots"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if commits log exists
if [ ! -f "$COMMITS_LOG" ]; then
    echo "❌ No snapshots found!"
    echo ""
    echo "💡 Create your first snapshot:"
    echo "   ./scripts/snapshot_create.sh 'Initial stable version'"
    echo ""
    exit 0
fi

# Count total snapshots
TOTAL_SNAPSHOTS=$(wc -l < "$COMMITS_LOG")
echo "📊 Total Snapshots: $TOTAL_SNAPSHOTS"

# Get current commit
CURRENT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
echo "📍 Current Commit:  $CURRENT_COMMIT"
echo ""

# Table header
printf "%-8s | %-10s | %-19s | %-35s | %-8s | %-10s\n" "Commit" "Date" "Time" "Message" "Size" "Records"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# List snapshots (most recent first)
tac "$COMMITS_LOG" | while IFS='|' read -r commit date message size records; do
    # Parse date/time
    SNAPSHOT_DATE=$(echo "$date" | cut -d'T' -f1)
    SNAPSHOT_TIME=$(echo "$date" | cut -d'T' -f2 | cut -d'+' -f1 | cut -c1-8)
    
    # Truncate message if too long
    TRUNCATED_MESSAGE="${message:0:35}"
    if [ ${#message} -gt 35 ]; then
        TRUNCATED_MESSAGE="${TRUNCATED_MESSAGE}..."
    fi
    
    # Mark current commit
    MARKER=""
    if [ "$commit" = "$CURRENT_COMMIT" ]; then
        MARKER="→"
    else
        MARKER=" "
    fi
    
    # Check if backup file exists
    BACKUP_EXISTS="✅"
    if [ ! -f "$BACKUPS_DIR/e173_gateway_${commit}.sql.gz" ]; then
        BACKUP_EXISTS="❌"
    fi
    
    printf "%s%-7s | %-10s | %-19s | %-35s | %-8s | %-10s %s\n" \
        "$MARKER" "$commit" "$SNAPSHOT_DATE" "$SNAPSHOT_TIME" \
        "$TRUNCATED_MESSAGE" "$size" "$records" "$BACKUP_EXISTS"
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🔧 Usage Examples:"
echo "   ./scripts/snapshot_restore.sh a1b2c3d     # Restore to specific snapshot"
echo "   ./scripts/rollback.sh                     # Rollback to previous snapshot"
echo "   ./scripts/rollback.sh 3                   # Rollback 3 snapshots"
echo "   ./scripts/deploy_fresh.sh a1b2c3d         # Deploy snapshot to server"
echo ""
echo "Legend:"
echo "   → Current active snapshot"
echo "   ✅ Database backup available"
echo "   ❌ Database backup missing"
