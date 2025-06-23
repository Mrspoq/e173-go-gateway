#!/bin/bash

# E173 Gateway: Quick Rollback System
# Instantly rollback to previous working snapshot

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
COMMITS_LOG="$SNAPSHOTS_DIR/commits.log"

# Default rollback steps
ROLLBACK_STEPS="${1:-1}"

echo "🔄 E173 Gateway Quick Rollback"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if commits log exists
if [ ! -f "$COMMITS_LOG" ]; then
    echo "❌ No snapshots found!"
    echo "💡 Create your first snapshot: ./scripts/snapshot_create.sh 'Initial snapshot'"
    exit 1
fi

# Get current commit
CURRENT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Find target commit for rollback
if [ "$ROLLBACK_STEPS" -eq 1 ]; then
    echo "🔍 Finding previous snapshot..."
    TARGET_COMMIT=$(tail -2 "$COMMITS_LOG" | head -1 | cut -d'|' -f1)
else
    echo "🔍 Finding snapshot $ROLLBACK_STEPS steps back..."
    TARGET_COMMIT=$(tail -$((ROLLBACK_STEPS + 1)) "$COMMITS_LOG" | head -1 | cut -d'|' -f1)
fi

if [ -z "$TARGET_COMMIT" ]; then
    echo "❌ Cannot find snapshot $ROLLBACK_STEPS steps back"
    echo "📋 Available snapshots:"
    tail -5 "$COMMITS_LOG" | while IFS='|' read -r commit date message size records; do
        echo "   $commit - $message"
    done
    exit 1
fi

# Get target snapshot info
TARGET_LINE=$(grep "^$TARGET_COMMIT|" "$COMMITS_LOG")
if [ -z "$TARGET_LINE" ]; then
    echo "❌ Target commit $TARGET_COMMIT not found in snapshots"
    exit 1
fi

# Parse target snapshot details
IFS='|' read -r commit date message size records <<< "$TARGET_LINE"

echo "📋 Current:  $CURRENT_COMMIT"
echo "📋 Target:   $TARGET_COMMIT"
echo "📝 Message:  $message"
echo "📅 Date:     $(echo $date | cut -d'T' -f1)"
echo "📊 Size:     $size"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$CURRENT_COMMIT" = "$TARGET_COMMIT" ]; then
    echo "ℹ️  Already at target snapshot $TARGET_COMMIT"
    exit 0
fi

echo ""
echo "⚠️  ROLLBACK WARNING:"
echo "   • Rolling back from: $CURRENT_COMMIT"
echo "   • Rolling back to:   $TARGET_COMMIT ($message)"
echo "   • Steps back:        $ROLLBACK_STEPS"
echo ""
echo "   This will:"
echo "   ✓ Restore code to previous version"
echo "   ✓ Restore database to previous state"
echo "   ✓ Restart all services"
echo "   ⚠ Current changes will be preserved in emergency backup"
echo ""

read -p "Proceed with rollback? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Rollback cancelled"
    exit 0
fi

echo ""
echo "🔄 Starting rollback process..."

# Create emergency backup of current state
echo "💾 Creating emergency backup of current state..."
EMERGENCY_MESSAGE="Emergency backup before rollback to $TARGET_COMMIT"
if ! git diff-index --quiet HEAD --; then
    git add .
    git commit -m "$EMERGENCY_MESSAGE" || true
    echo "✅ Emergency backup created"
else
    echo "ℹ️  No uncommitted changes to backup"
fi

# Perform the rollback using snapshot_restore.sh
echo "🔄 Executing rollback to $TARGET_COMMIT..."
"$SCRIPT_DIR/snapshot_restore.sh" "$TARGET_COMMIT" <<< "yes"

# Verify rollback success
echo ""
echo "🔍 Verifying rollback..."

# Check if we're on the right commit
RESTORED_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
if [ "$RESTORED_COMMIT" = "$TARGET_COMMIT" ]; then
    echo "✅ Code rollback successful"
else
    echo "⚠️  Code rollback may be incomplete"
fi

# Check if server is responding
sleep 3
if curl -s http://localhost:8080/ping > /dev/null 2>&1; then
    echo "✅ Server is responding after rollback"
else
    echo "⚠️  Server may need more time to start"
fi

echo ""
echo "🎉 Rollback completed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Rolled back to: $TARGET_COMMIT"
echo "📝 Message:        $message"
echo "🕒 Rollback time:  $(date)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🔧 Next steps:"
echo "   • Verify functionality: curl http://localhost:8080/"
echo "   • Check dashboard: http://localhost:8080"
echo "   • If issues persist: ./scripts/rollback.sh 2"
echo ""
echo "🔄 To rollback this rollback:"
echo "   • Find emergency backup: git log --oneline | grep 'Emergency backup'"
echo "   • Restore to latest: ./scripts/snapshot_restore.sh <latest_commit>"
