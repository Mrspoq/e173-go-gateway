# Browser Tools MCP Setup Instructions

## Overview
Browser Tools MCP enables Claude Code to interact with and debug the E173 Gateway UI through Chrome browser automation.

## Installation Steps

### 1. Chrome Extension Installation
The Browser Tools Chrome Extension is required for capturing browser data.

**Manual Installation:**
1. Go to: https://github.com/AgentDeskAI/browser-tools-mcp/releases
2. Download the latest `.zip` file (e.g., `browser-tools-extension-v1.2.0.zip`)
3. Extract the ZIP file
4. Open Chrome and go to `chrome://extensions/`
5. Enable "Developer mode" (top right)
6. Click "Load unpacked"
7. Select the extracted extension folder

### 2. Start Browser Tools Server
```bash
cd /root/e173_go_gateway/mcp/browser-tools
./run-ui-tests.sh
```

This will:
- Start the Browser Tools server on port 3000
- Prepare the testing environment
- Show status and options

### 3. Launch Chrome with Extension
```bash
./launch-chrome-with-extension.sh
```

### 4. Connect Extension to Server
1. Open Chrome DevTools (F12)
2. Find the "Browser Tools" tab
3. Click "Connect" button
4. You should see "Connected to server"

### 5. Configure Claude Code
Add to your MCP configuration:
```json
{
  "mcpServers": {
    "browser-tools": {
      "command": "npx",
      "args": ["@agentdeskai/browser-tools-mcp@latest"],
      "env": {
        "BROWSER_TOOLS_SERVER_URL": "http://localhost:3000"
      }
    }
  }
}
```

## Using Browser Tools MCP

### Available Tools
Once connected, Claude Code can use these tools:

1. **screenshot** - Capture current page
2. **console_logs** - Get browser console output
3. **network_logs** - Monitor API calls
4. **click** - Click elements
5. **type** - Enter text
6. **navigate** - Go to URLs
7. **element_info** - Get element details

### Testing E173 Gateway

#### Test Dashboard Layout
```
1. Navigate to http://192.168.1.35:8080/dashboard
2. Take screenshot
3. Check console for errors
4. Verify grid-cols-5 on #stats-cards
5. Count child divs (should be 5)
```

#### Test Gateway Page
```
1. Navigate to /gateways
2. Check for blank page
3. Look for template errors in console
4. Test "Add Gateway" button
5. Test connection test buttons
```

#### Test Customer Management
```
1. Navigate to /customers
2. Click "Add Customer" button
3. Verify navigation (no login redirect)
4. Test edit buttons
5. Check form submission
```

## Debugging Tips

### Console Errors
Browser Tools captures all console output including:
- JavaScript errors
- HTMX events
- API response errors
- Template rendering issues

### Network Monitoring
Monitor failed requests:
- 401 Unauthorized (auth issues)
- 404 Not Found (wrong endpoints)
- 500 Server Error (backend issues)

### Common Issues

1. **Blank Pages**
   - Check console for template errors
   - Verify CurrentUser data structure
   - Look for missing closing tags

2. **Button Not Working**
   - Check if hx-get/hx-post is correct
   - Verify endpoint exists
   - Look for JavaScript errors

3. **Layout Issues**
   - Take screenshot for visual check
   - Inspect element classes
   - Check for CSS conflicts

## Automated Testing

Run comprehensive tests:
```bash
node comprehensive-ui-test.js
```

Generate test report:
```bash
./generate-test-report.sh
```

## Troubleshooting

### Extension Not Connecting
1. Ensure server is running (port 3000)
2. Check Chrome extension is loaded
3. Verify no firewall blocking
4. Try refreshing the page

### No Console Logs
1. Make sure DevTools is open
2. Browser Tools tab selected
3. Click "Connect" if disconnected
4. Check extension permissions

### Server Issues
```bash
# Check if server is running
ps aux | grep browser-tools-server

# Restart server
pkill -f browser-tools-server
./run-ui-tests.sh
```

## Summary
With Browser Tools MCP, you can:
- See exactly what's happening in the browser
- Debug HTMX requests and responses  
- Identify console errors immediately
- Test UI interactions programmatically
- Generate comprehensive test reports

This gives Claude Code full visibility into UI issues, making it much easier to identify and fix problems accurately.