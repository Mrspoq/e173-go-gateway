#!/bin/bash

echo "Setting up Browser Tools MCP for E173 Gateway UI Testing"
echo "========================================================="

# Create package.json for the project
cat > package.json << 'EOF'
{
  "name": "e173-browser-tools-mcp",
  "version": "1.0.0",
  "description": "Browser Tools MCP setup for E173 Gateway UI testing",
  "scripts": {
    "start-server": "npx @agentdeskai/browser-tools-server@latest",
    "start-mcp": "npx @agentdeskai/browser-tools-mcp@latest",
    "test-ui": "node test-ui.js"
  },
  "dependencies": {
    "@agentdeskai/browser-tools-server": "latest",
    "@agentdeskai/browser-tools-mcp": "latest"
  }
}
EOF

echo "✅ Created package.json"

# Create MCP configuration
cat > browser-tools-config.json << 'EOF'
{
  "mcpServers": {
    "browser-tools": {
      "command": "npx",
      "args": ["@agentdeskai/browser-tools-mcp@latest"],
      "env": {
        "BROWSER_TOOLS_SERVER_URL": "http://localhost:3000",
        "LOG_LEVEL": "debug"
      }
    }
  }
}
EOF

echo "✅ Created MCP configuration"

# Create startup script
cat > start-browser-tools.sh << 'EOF'
#!/bin/bash

echo "Starting Browser Tools MCP Server..."
echo "===================================="

# Start the browser tools server in the background
echo "Starting Browser Tools Server on port 3000..."
npx @agentdeskai/browser-tools-server@latest &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo "Browser Tools Server started with PID: $SERVER_PID"
echo ""
echo "⚠️  IMPORTANT: Install the Chrome Extension"
echo "   Download from: https://github.com/AgentDeskAI/browser-tools-mcp/releases"
echo ""
echo "To stop the server, run: kill $SERVER_PID"
echo ""
echo "Server is ready for UI testing!"

# Keep the script running
wait $SERVER_PID
EOF

chmod +x start-browser-tools.sh

echo "✅ Created startup script"

# Create test configuration for E173 Gateway
cat > e173-ui-test-config.js << 'EOF'
// E173 Gateway UI Test Configuration
module.exports = {
  baseURL: 'http://192.168.1.35:8080',
  credentials: {
    username: 'admin',
    password: 'admin123'
  },
  tests: [
    {
      name: 'Dashboard Layout Test',
      path: '/dashboard',
      checks: [
        { type: 'element', selector: '#stats-cards', description: 'Stats cards container' },
        { type: 'console', level: 'error', description: 'Check for console errors' },
        { type: 'network', url: '/api/stats/', description: 'Stats API calls' }
      ]
    },
    {
      name: 'Gateway Page Test',
      path: '/gateways',
      checks: [
        { type: 'element', selector: '.gateway-card', description: 'Gateway cards' },
        { type: 'button', selector: '.btn-test', description: 'Test button functionality' }
      ]
    },
    {
      name: 'Customer Management Test',
      path: '/customers',
      checks: [
        { type: 'button', selector: '#add-customer-btn', description: 'Add customer button' },
        { type: 'button', selector: '.btn-edit', description: 'Edit buttons' }
      ]
    },
    {
      name: 'CDR Page Test',
      path: '/cdrs',
      checks: [
        { type: 'element', selector: 'table', description: 'CDR table structure' },
        { type: 'element', selector: 'thead', description: 'Table headers' }
      ]
    }
  ]
};
EOF

echo "✅ Created test configuration"

# Create documentation
cat > README.md << 'EOF'
# Browser Tools MCP for E173 Gateway

## Overview
This setup provides browser automation and debugging capabilities for testing the E173 Gateway UI.

## Components
1. **Browser Tools Server**: Local Node.js server that collects browser data
2. **Browser Tools MCP**: MCP interface for Claude Code
3. **Chrome Extension**: Captures console logs, network requests, and DOM elements

## Setup Instructions

### 1. Install Chrome Extension
Download and install the Browser Tools Chrome Extension from:
https://github.com/AgentDeskAI/browser-tools-mcp/releases

### 2. Start the Server
```bash
./start-browser-tools.sh
```

### 3. Configure Claude Code
Add the following to your Claude Code MCP configuration:
```json
{
  "mcpServers": {
    "browser-tools": {
      "command": "npx",
      "args": ["@agentdeskai/browser-tools-mcp@latest"]
    }
  }
}
```

## Features
- 📸 Screenshot capture
- 📝 Console log monitoring
- 🌐 Network request tracking
- 🎯 DOM element selection
- 🔍 Developer tools integration
- 🚦 Lighthouse audits

## Testing E173 Gateway

The configuration includes tests for:
- Dashboard layout (5-column grid)
- Gateway management page
- Customer management buttons
- CDR table structure

## Usage

1. Start the Browser Tools server
2. Open Chrome and navigate to http://192.168.1.35:8080
3. Login with admin/admin123
4. The extension will capture all browser activity
5. Use Claude Code to analyze logs and debug issues

## Troubleshooting

- Ensure Chrome Extension is installed and active
- Check that the server is running on port 3000
- Verify the E173 Gateway is accessible at the configured URL
EOF

echo "✅ Created documentation"

echo ""
echo "================================"
echo "✅ Browser Tools MCP Setup Complete!"
echo "================================"
echo ""
echo "Next steps:"
echo "1. Download Chrome Extension from: https://github.com/AgentDeskAI/browser-tools-mcp/releases"
echo "2. Run: ./start-browser-tools.sh"
echo "3. Configure Claude Code with the MCP server"
echo ""