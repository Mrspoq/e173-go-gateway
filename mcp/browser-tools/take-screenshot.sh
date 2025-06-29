#!/bin/bash

echo "Taking screenshot of current browser tab..."
echo "=========================================="

# Take a screenshot using Browser Tools MCP
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>/dev/null
{"jsonrpc": "2.0", "method": "initialize", "params": {"capabilities": {}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "takeScreenshot", "arguments": {}}, "id": 2}
EOF

echo ""
echo "Screenshot request sent!"
echo "Check Downloads/mcp-screenshots/ for the image"