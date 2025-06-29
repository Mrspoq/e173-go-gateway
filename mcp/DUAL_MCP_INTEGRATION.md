# Dual MCP Integration Guide: Browser Use + Browser Tools

## Overview
We now have two complementary MCP servers for comprehensive browser testing:

1. **Browser Tools MCP** (`/mcp/browser-tools/`) - Monitoring & Analysis
2. **Browser Use MCP** (`/mcp/browser-use/`) - Control & Automation

## Architecture

```
┌─────────────────────┐     ┌──────────────────────┐
│  Browser Use MCP    │     │  Browser Tools MCP   │
├─────────────────────┤     ├──────────────────────┤
│ • Navigate pages    │     │ • Console logs       │
│ • Click elements    │     │ • Network monitoring │
│ • Fill forms        │     │ • Error detection    │
│ • AI-driven control │     │ • Screenshots        │
└─────────────────────┘     └──────────────────────┘
         │                           │
         └───────────┬───────────────┘
                     │
              ┌──────▼──────┐
              │   Browser   │
              │ (Chrome)    │
              └─────────────┘
```

## Combined Usage Examples

### Example 1: Test Dashboard with Full Monitoring
```bash
# Step 1: Browser Use navigates to dashboard
{
  "tool": "run_browser_agent",
  "arguments": {
    "task": "Navigate to http://192.168.1.35:8080/dashboard"
  }
}

# Step 2: Browser Tools captures state
{
  "tool": "takeScreenshot",
  "arguments": {}
}

# Step 3: Browser Tools checks console
{
  "tool": "getConsoleLogs",
  "arguments": {}
}

# Step 4: Browser Use verifies layout
{
  "tool": "run_browser_agent",
  "arguments": {
    "task": "Verify that 5 stat cards are displayed horizontally and count the visible panels"
  }
}
```

### Example 2: Form Testing with Error Monitoring
```bash
# Browser Tools: Start monitoring console
{
  "tool": "wipeLogs",
  "arguments": {}
}

# Browser Use: Fill and submit form
{
  "tool": "run_browser_agent",
  "arguments": {
    "task": "Fill the customer creation form with test data and submit"
  }
}

# Browser Tools: Check for errors
{
  "tool": "getConsoleErrors",
  "arguments": {}
}
```

## Capabilities Comparison

| Feature | Browser Tools MCP | Browser Use MCP |
|---------|------------------|-----------------|
| Navigate to URL | ❌ | ✅ |
| Click elements | ❌ | ✅ |
| Fill forms | ❌ | ✅ |
| Take screenshots | ✅ | ✅ (via vision) |
| Console logs | ✅ | ❌ |
| Network monitoring | ✅ | ❌ |
| Error detection | ✅ | ❌ |
| Performance audit | ✅ | ❌ |
| AI interpretation | ❌ | ✅ |
| Natural language tasks | ❌ | ✅ |

## Best Practices

### 1. Use Browser Use for Actions
- Navigation
- Clicking buttons
- Filling forms
- Complex interactions

### 2. Use Browser Tools for Monitoring
- Console errors
- Network failures
- Performance issues
- Visual verification

### 3. Combine for Complete Testing
```
Browser Use (navigate) → Browser Tools (monitor) → Browser Use (interact) → Browser Tools (verify)
```

## Quick Reference

### Browser Use MCP
```bash
cd /root/e173_go_gateway/mcp/browser-use
source .venv/bin/activate
mcp-server-browser-use
```

**Main Tools:**
- `run_browser_agent` - Execute browser tasks
- `run_deep_research` - Web research

### Browser Tools MCP
```bash
cd /root/e173_go_gateway/mcp/browser-tools
npx @agentdeskai/browser-tools-mcp@latest
```

**Main Tools:**
- `takeScreenshot` - Capture page
- `getConsoleLogs` - Get console output
- `getConsoleErrors` - Get errors only
- `getNetworkLogs` - Network requests
- `runPerformanceAudit` - Performance check

## Testing Workflow for E173 Gateway

1. **Setup Phase**
   - Start both MCP servers
   - Open browser on target machine
   - Ensure DISPLAY=:0 is set

2. **Navigation Phase** (Browser Use)
   - Navigate to target page
   - Wait for page load

3. **Monitoring Phase** (Browser Tools)
   - Capture initial state
   - Clear console logs
   - Start monitoring

4. **Interaction Phase** (Browser Use)
   - Perform user actions
   - Fill forms
   - Click buttons

5. **Verification Phase** (Browser Tools)
   - Check console errors
   - Verify network calls
   - Take final screenshot

6. **Analysis Phase** (Both)
   - Browser Use: AI analysis of results
   - Browser Tools: Technical verification

## Limitations & Workarounds

### Current Limitations
- Browser Tools cannot reload pages
- Browser Use requires LLM configuration
- Both need same browser instance

### Workarounds
- Use Browser Use for all navigation
- Configure Ollama for free local LLM
- Keep browser open between tests

## Future Enhancements
- Unified MCP orchestrator
- Automatic test generation
- CI/CD integration
- Test result reporting