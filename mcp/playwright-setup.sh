#!/bin/bash

echo "Setting up Playwright MCP for Browser Control"
echo "============================================="

cd /root/e173_go_gateway/mcp

# Create playwright directory
mkdir -p playwright
cd playwright

# Initialize npm project
npm init -y

# Install official Microsoft Playwright MCP
echo "Installing Playwright MCP..."
npm install @modelcontextprotocol/server-playwright

# Create configuration
cat > playwright-mcp-config.json << 'EOF'
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-playwright"]
    }
  }
}
EOF

# Create test script
cat > test-playwright.sh << 'EOF'
#!/bin/bash

echo "Testing Playwright MCP..."

# Initialize and list tools
cat << 'EOT' | npx @modelcontextprotocol/server-playwright 2>/dev/null | jq -r '.result.tools[] | .name' 2>/dev/null
{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "1.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}, "id": 1}
{"jsonrpc": "2.0", "method": "tools/list", "params": {}, "id": 2}
EOT
EOF

chmod +x test-playwright.sh

echo ""
echo "✅ Playwright MCP setup complete!"
echo ""
echo "With Playwright MCP, you can:"
echo "- navigate: Go to any URL"
echo "- reload: Refresh the current page"
echo "- click: Click on elements"
echo "- fill: Fill in forms"
echo "- screenshot: Take screenshots"
echo "- execute: Run JavaScript"
echo ""
echo "Usage:"
echo "1. Run the server: npx @modelcontextprotocol/server-playwright"
echo "2. Use both Browser Tools MCP (monitoring) + Playwright MCP (control)"