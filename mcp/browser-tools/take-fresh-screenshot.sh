#!/bin/bash

echo "Fresh Screenshot Workflow"
echo "========================"

# 1. Check server is running
echo "1. Checking server status..."
if curl -s http://192.168.1.35:8080/ping > /dev/null; then
    echo "   ✓ Server is running"
else
    echo "   ✗ Server is down!"
    exit 1
fi

# 2. Verify endpoint
echo "2. Checking stats endpoint..."
if curl -s http://192.168.1.35:8080/api/stats/cards | grep -q "display: flex"; then
    echo "   ✓ Endpoint returns updated HTML"
else
    echo "   ✗ Endpoint not returning expected HTML"
fi

# 3. Wait for page reload
echo "3. Please reload the browser page (F5)..."
echo "   Waiting 3 seconds..."
sleep 3

# 4. Take screenshot
echo "4. Taking screenshot..."
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>/dev/null | grep -E "(Successfully saved|error)"
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "takeScreenshot", "arguments": {}}, "id": 2}
EOF

# 5. Find latest screenshot
LATEST=$(ls -t /root/Downloads/mcp-screenshots/*.png 2>/dev/null | head -1)
echo ""
echo "Screenshot saved: $LATEST"
echo "Ready to analyze!"