#!/bin/bash

# GitHub Issue Creation Script
# Requires GITHUB_TOKEN environment variable

if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GITHUB_TOKEN environment variable not set"
    echo "Please export GITHUB_TOKEN=your_token_here"
    exit 1
fi

# Repository information
OWNER="Mrspoq"
REPO="e173-go-gateway"

# Function to create an issue
create_issue() {
    local title=$1
    local body=$2
    local labels=$3
    
    echo "Creating issue: $title"
    
    curl -X POST \
        -H "Authorization: token $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/$OWNER/$REPO/issues \
        -d "{\"title\":\"$title\",\"body\":\"$body\",\"labels\":$labels}"
}

# Example usage
# create_issue "Test Issue" "This is a test issue body" '["bug", "enhancement"]'