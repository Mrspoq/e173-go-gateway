#!/bin/bash

# Script to clean sensitive data from Git history

echo "Starting Git history cleanup..."

# Create backup
echo "Creating backup..."
cp -r .git .git.backup

# Define files with secrets
SECRET_FILES=(
    "scripts/.env"
    ".env"
    "MCP_API_KEYS_BACKUP.md"
    ".claude/settings.local.json"
    "mcp/github-mcp-config.json"
    "scripts/create_admin_user.go"
    "scripts/add_sample_gateways.sql"
    "scripts/add_sample_gateways_fixed.sql"
    "documentation/plan.md"
    "FINAL_REPORT_2025-06-29.md"
    "mcp/browser-use/test-all-pages.py"
)

# Remove files with secrets from history
for file in "${SECRET_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "Removing $file from history..."
        git filter-branch --force --index-filter \
            "git rm --cached --ignore-unmatch $file" \
            --prune-empty --tag-name-filter cat -- --all
    fi
done

# Clean up refs
git for-each-ref --format="%(refname)" refs/original/ | xargs -n 1 git update-ref -d

# Garbage collect
git gc --prune=now --aggressive

echo "Git history cleanup complete!"
echo "Don't forget to force push with: git push --force origin master"