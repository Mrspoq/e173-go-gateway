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
