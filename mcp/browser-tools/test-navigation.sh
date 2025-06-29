#!/bin/bash

echo "Testing Browser Tools MCP Navigation Capabilities"
echo "================================================"

# Test 1: Try to reload/refresh
echo -e "\n1. Testing reload/refresh..."
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>&1 | grep -E "(result|error)" | tail -5
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "reload", "arguments": {}}, "id": 2}
EOF

# Test 2: Try navigate
echo -e "\n2. Testing navigate..."
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>&1 | grep -E "(result|error)" | tail -5
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "navigate", "arguments": {"url": "http://192.168.1.35:8080/dashboard"}}, "id": 2}
EOF

# Test 3: Try executeScript
echo -e "\n3. Testing executeScript..."
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>&1 | grep -E "(result|error)" | tail -5
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "executeScript", "arguments": {"script": "location.reload()"}}, "id": 2}
EOF

# Test 4: Check debugger mode capabilities
echo -e "\n4. Testing debugger mode (might have navigation)..."
cat << 'EOF' | npx @agentdeskai/browser-tools-mcp@latest 2>&1 | grep -E "(result|error)" | tail -5
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "clientInfo": {"name": "claude", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "runDebuggerMode", "arguments": {}}, "id": 2}
EOF