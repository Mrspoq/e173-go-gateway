#!/bin/bash

# Test the GitHub MCP server with proper MCP protocol

source ../.github_config

echo "Testing GitHub MCP Server with proper protocol..."
echo "================================================"

# First, let's initialize and list available tools
echo -e "\n1. Initializing MCP connection..."
cat << 'EOF' | docker run -i --rm -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN ghcr.io/github/github-mcp-server stdio
{"jsonrpc": "2.0", "method": "initialize", "params": {"capabilities": {}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/list", "params": {}, "id": 2}
EOF