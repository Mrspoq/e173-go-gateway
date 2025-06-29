#!/bin/bash

# Load GitHub configuration
source ../.github_config

# Run the GitHub MCP server with all toolsets enabled
echo "Starting GitHub MCP server..."
echo "Repository: $GITHUB_USER/$GITHUB_REPO"
echo "Using token: ${GITHUB_TOKEN:0:10}..."

docker run -i --rm \
  -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN \
  ghcr.io/github/github-mcp-server \
  stdio \
  --toolsets all \
  --enable-command-logging \
  --log-file /tmp/github-mcp.log