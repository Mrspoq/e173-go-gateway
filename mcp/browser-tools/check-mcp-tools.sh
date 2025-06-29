#!/bin/bash

echo "Checking Browser Tools MCP available tools..."
echo "============================================"

# Send MCP protocol commands to discover available tools
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>/dev/null | grep -A100 '"tools"' | jq '.'
{"jsonrpc": "2.0", "method": "initialize", "params": {"capabilities": {}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/list", "params": {}, "id": 2}
EOF