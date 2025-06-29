# Browser Tools MCP Setup

## Overview
Browser Tools MCP is a powerful browser monitoring and interaction tool that enables AI-powered applications to capture and analyze browser data through a Chrome extension.

## Components Installed

### 1. MCP Server
- **Status**: ✅ Installed and configured
- **Command**: `npx @agentdeskai/browser-tools-mcp@latest`
- **Location**: Running from /root/Cline/MCP
- **Configuration**: Added to cline_mcp_settings.json

### 2. Browser Tools Server (Middleware)
- **Status**: ✅ Running
- **Port**: 3025
- **URL**: http://localhost:3025
- **Purpose**: Acts as middleware between Chrome extension and MCP server

### 3. Chrome Extension
- **Status**: ✅ Downloaded and extracted
- **Version**: v1.2.0
- **Location**: `/root/Cline/MCP/BrowserTools-extension/chrome-extension/`
- **Download URL**: https://github.com/AgentDeskAI/browser-tools-mcp/releases/download/v1.2.0/BrowserTools-1.2.0-extension.zip

## Available Tools

Once the Chrome extension is installed and connected, the following tools will be available:

1. **getConsoleLogs** - Check browser console logs
2. **getConsoleErrors** - Check browser console errors
3. **getNetworkErrors** - Check network error logs
4. **getNetworkLogs** - Check all network logs
5. **takeScreenshot** - Take a screenshot of the current browser tab
6. **getSelectedElement** - Get the selected element from the browser
7. **wipeLogs** - Clear all browser logs from memory
8. **runAccessibilityAudit** - Run accessibility audit on current page
9. **runPerformanceAudit** - Run performance audit on current page
10. **runSEOAudit** - Run SEO audit on current page
11. **runNextJSAudit** - Run NextJS-specific audit
12. **runDebuggerMode** - Run debugging tools in sequence
13. **runAuditMode** - Run all audit tools in sequence
14. **runBestPracticesAudit** - Run best practices audit

## Next Steps

1. ✅ Chrome extension downloaded
2. ✅ Extension unzipped to `/root/Cline/MCP/BrowserTools-extension/chrome-extension/`
3. Install the extension in Chrome:
   - Open Chrome and navigate to `chrome://extensions/`
   - Enable "Developer mode"
   - Click "Load unpacked"
   - Select the unzipped extension folder
4. Open Chrome DevTools and access the BrowserToolsMCP panel
5. The extension should automatically connect to the running servers

## Troubleshooting

If connection issues occur:
- Ensure only ONE Chrome DevTools panel is open
- Restart the browser-tools-server
- Check that both servers are running
- Verify the Chrome extension is enabled

## Configuration

The MCP server has been configured in:
`/root/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`

```json
{
  "mcpServers": {
    "github.com/AgentDeskAI/browser-tools-mcp": {
      "command": "npx",
      "args": ["@agentdeskai/browser-tools-mcp@latest"],
      "disabled": false,
      "autoApprove": []
    }
  }
}
