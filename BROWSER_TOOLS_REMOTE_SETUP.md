# Browser Tools MCP - Remote Setup Guide

## Your Setup:
- **Development Machine**: Linux server at 192.168.1.35
- **Client Machine**: Mac with VSCode (SSH Remote)
- **Purpose**: Automated browser testing for Claude/Cline

## ✅ Server Status:
- Browser Tools Server: Running on `0.0.0.0:3025` (accessible from network)
- MCP Server: Running and configured

## 🚀 Quick Setup Steps for Your Mac:

### Step 1: Copy Chrome Extension to Your Mac
On your Mac terminal, run:
```bash
scp -r root@192.168.1.35:/root/Cline/MCP/BrowserTools-extension/chrome-extension ~/Desktop/BrowserTools-extension
```

### Step 2: Install Extension in Chrome (on your Mac)
1. Open Chrome on your Mac
2. Go to `chrome://extensions/`
3. Enable "Developer mode" (top right toggle)
4. Click "Load unpacked"
5. Select the `~/Desktop/BrowserTools-extension` folder
6. The extension icon should appear in your Chrome toolbar

### Step 3: Test Connection
1. Open any website in Chrome
2. Open DevTools (F12 or Cmd+Option+I)
3. Look for "BrowserToolsMCP" panel in DevTools
4. The extension should automatically connect to `192.168.1.35:3025`

## 🎯 How I (Claude/Cline) Will Use This:

Once the extension is installed on your Mac's Chrome, I can automatically:

- **Take screenshots** of your application as I develop features
- **Check console logs** for JavaScript errors during testing
- **Monitor network requests** to debug API issues
- **Run accessibility audits** to ensure WCAG compliance
- **Analyze performance** to optimize loading times
- **Debug issues** without asking you to manually check

### Example Commands You Can Give Me:
- "Test the login page and fix any console errors"
- "Take a screenshot of the dashboard and optimize its performance"
- "Debug why the API calls are failing on the customer page"
- "Run an accessibility audit on all pages and fix the issues"

## ⚡ Connection Details:
- **Server IP**: 192.168.1.35
- **Port**: 3025
- **Protocol**: WebSocket + HTTP

The extension will automatically try to connect to the server. Since the server is already listening on all interfaces (0.0.0.0), no additional configuration is needed!

## 🔍 Verify It's Working:
Once installed, ask me to:
```
"Take a screenshot of google.com to test the browser tools"
```

If successful, I'll be able to capture and analyze the page automatically!
