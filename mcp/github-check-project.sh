#!/bin/bash

# Check our E173 Gateway project status using GitHub MCP

source ../.github_config

echo "Checking E173 Gateway Project Status..."
echo "======================================"

# Check repository info and open issues
cat << 'EOF' | docker run -i --rm -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN ghcr.io/github/github-mcp-server stdio 2>/dev/null | grep -A50 '"result"' | jq -r '.result.content // .result'
{"jsonrpc": "2.0", "method": "initialize", "params": {"capabilities": {}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_me", "arguments": {}}, "id": 2}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "list_issues", "arguments": {"owner": "Mrspoq", "repo": "e173-go-gateway", "state": "open"}}, "id": 3}
EOF