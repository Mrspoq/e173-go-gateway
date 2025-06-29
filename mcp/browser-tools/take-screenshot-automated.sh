#!/bin/bash

# Function to take screenshot
take_screenshot() {
    echo "Taking screenshot..."
    
    cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>/dev/null | grep -E "(Successfully saved|error)"
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "takeScreenshot", "arguments": {}}, "id": 2}
EOF
    
    # Get latest screenshot
    LATEST=$(ls -t /root/Downloads/mcp-screenshots/*.png 2>/dev/null | head -1)
    echo "Screenshot saved: $LATEST"
}

# Function to get console logs
get_console_logs() {
    echo "Getting console logs..."
    
    cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>/dev/null | jq -r '.result.content[0].text' 2>/dev/null
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "getConsoleLogs", "arguments": {}}, "id": 2}
EOF
}

# Main execution
echo "=== Browser Tools MCP Dashboard Analysis ==="
echo ""

take_screenshot
echo ""
get_console_logs