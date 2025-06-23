#!/bin/bash

# E173 Gateway: Snapshot Cleanup Utility
# Cleans up old snapshots to save disk space

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
COMMITS_LOG="$SNAPSHOTS_DIR/commits.log"
BACKUPS_DIR="$SNAPSHOTS_DIR/backups"

# Default retention (days)
RETENTION_DAYS="${1:-30}"

echo "🧹 E173 Gateway Snapshot Cleanup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📅 Retention: $RETENTION_DAYS days"

# Check if snapshots exist
if [ ! -f "$COMMITS_LOG" ]; then
    echo "ℹ️  No snapshots found to clean up"
    exit 0
fi

# Find old snapshots
CUTOFF_DATE=$(date -d "$RETENTION_DAYS days ago" +%Y-%m-%d)
echo "🗓️  Cutoff date: $CUTOFF_DATE"

# Count snapshots before cleanup
TOTAL_BEFORE=$(wc -l < "$COMMITS_LOG" 2>/dev/null || echo "0")
echo "📊 Snapshots before cleanup: $TOTAL_BEFORE"

# Create temporary files
TEMP_LOG=$(mktemp)
OLD_SNAPSHOTS=$(mktemp)

# Process commits log
while IFS='|' read -r commit date message size records; do
    SNAPSHOT_DATE=$(echo "$date" | cut -d'T' -f1)
    
    if [[ "$SNAPSHOT_DATE" < "$CUTOFF_DATE" ]]; then
        echo "$commit|$date|$message|$size|$records" >> "$OLD_SNAPSHOTS"
    else
        echo "$commit|$date|$message|$size|$records" >> "$TEMP_LOG"
    fi
done < "$COMMITS_LOG"

# Count old snapshots
OLD_COUNT=$(wc -l < "$OLD_SNAPSHOTS" 2>/dev/null || echo "0")

if [ "$OLD_COUNT" -eq 0 ]; then
    echo "✅ No old snapshots to clean up"
    rm -f "$TEMP_LOG" "$OLD_SNAPSHOTS"
    exit 0
fi

echo "🗑️  Found $OLD_COUNT old snapshots to remove"
echo ""
echo "📋 Snapshots to be removed:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

TOTAL_SIZE=0
while IFS='|' read -r commit date message size records; do
    echo "   $commit - $(echo $date | cut -d'T' -f1) - $message"
    
    # Calculate size (convert to bytes for sum)
    if [[ "$size" =~ ([0-9.]+)([KMG]B?) ]]; then
        SIZE_NUM="${BASH_REMATCH[1]}"
        SIZE_UNIT="${BASH_REMATCH[2]}"
        case "$SIZE_UNIT" in
            "KB") SIZE_BYTES=$(echo "$SIZE_NUM * 1024" | bc -l) ;;
            "MB") SIZE_BYTES=$(echo "$SIZE_NUM * 1024 * 1024" | bc -l) ;;
            "GB") SIZE_BYTES=$(echo "$SIZE_NUM * 1024 * 1024 * 1024" | bc -l) ;;
            *) SIZE_BYTES="$SIZE_NUM" ;;
        esac
        TOTAL_SIZE=$(echo "$TOTAL_SIZE + $SIZE_BYTES" | bc -l)
    fi
done < "$OLD_SNAPSHOTS"

# Convert total size back to human readable
if command -v numfmt > /dev/null 2>&1; then
    TOTAL_SIZE_HR=$(echo "$TOTAL_SIZE" | numfmt --to=iec-i --suffix=B)
else
    TOTAL_SIZE_HR="${TOTAL_SIZE%.*} bytes"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💾 Total space to reclaim: ~$TOTAL_SIZE_HR"
echo ""

read -p "Proceed with cleanup? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Cleanup cancelled"
    rm -f "$TEMP_LOG" "$OLD_SNAPSHOTS"
    exit 0
fi

echo ""
echo "🧹 Cleaning up old snapshots..."

# Remove old backup files and metadata
REMOVED_COUNT=0
REMOVED_SIZE=0

while IFS='|' read -r commit date message size records; do
    BACKUP_FILE="$BACKUPS_DIR/e173_gateway_${commit}.sql.gz"
    METADATA_FILE="$BACKUPS_DIR/e173_gateway_${commit}.meta.json"
    
    # Remove backup file
    if [ -f "$BACKUP_FILE" ]; then
        FILE_SIZE=$(stat -f%z "$BACKUP_FILE" 2>/dev/null || stat -c%s "$BACKUP_FILE" 2>/dev/null || echo "0")
        rm -f "$BACKUP_FILE"
        REMOVED_SIZE=$((REMOVED_SIZE + FILE_SIZE))
        echo "   🗑️  Removed: $(basename "$BACKUP_FILE")"
    fi
    
    # Remove metadata file
    if [ -f "$METADATA_FILE" ]; then
        rm -f "$METADATA_FILE"
        echo "   🗑️  Removed: $(basename "$METADATA_FILE")"
    fi
    
    REMOVED_COUNT=$((REMOVED_COUNT + 1))
done < "$OLD_SNAPSHOTS"

# Update commits log
mv "$TEMP_LOG" "$COMMITS_LOG"

# Convert removed size to human readable
if command -v numfmt > /dev/null 2>&1; then
    REMOVED_SIZE_HR=$(echo "$REMOVED_SIZE" | numfmt --to=iec-i --suffix=B)
else
    REMOVED_SIZE_HR="$REMOVED_SIZE bytes"
fi

# Count snapshots after cleanup
TOTAL_AFTER=$(wc -l < "$COMMITS_LOG" 2>/dev/null || echo "0")

echo ""
echo "✅ Cleanup completed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🗑️  Removed snapshots: $REMOVED_COUNT"
echo "💾 Space reclaimed:    $REMOVED_SIZE_HR"
echo "📊 Snapshots before:   $TOTAL_BEFORE"
echo "📊 Snapshots after:    $TOTAL_AFTER"
echo "📅 Retention policy:   $RETENTION_DAYS days"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Cleanup temp files
rm -f "$OLD_SNAPSHOTS"

echo ""
echo "💡 To change retention policy:"
echo "   $0 7     # Keep last 7 days"
echo "   $0 90    # Keep last 90 days"
