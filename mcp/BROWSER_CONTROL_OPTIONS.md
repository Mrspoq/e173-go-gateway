# Browser Control Options for E173 Gateway Testing

## Current Situation
- **Browser Tools MCP**: Monitoring only (screenshots, logs)
- **Need**: Browser control (navigate, reload, click)

## Option 1: Playwright MCP (Recommended)
Full browser automation with MCP integration.

### Installation:
```bash
cd /root/e173_go_gateway/mcp
mkdir playwright-mcp
cd playwright-mcp

# Install Playwright MCP
npm init -y
npm install @modelcontextprotocol/server-playwright
```

### Capabilities:
- Navigate to URLs
- Reload pages
- Click elements
- Fill forms
- Take screenshots
- Execute JavaScript
- Full browser automation

## Option 2: Selenium WebDriver
Traditional browser automation.

### Installation:
```bash
# Install Selenium
pip install selenium webdriver-manager
```

## Option 3: Simple Remote Control
Using xdotool for basic control.

### Installation:
```bash
sudo apt-get install xdotool
```

### Usage:
```bash
# Reload browser
DISPLAY=:0 xdotool key F5

# Navigate to URL
DISPLAY=:0 xdotool key ctrl+l
DISPLAY=:0 xdotool type "http://192.168.1.35:8080/dashboard"
DISPLAY=:0 xdotool key Return
```

## Option 4: Browser Extension
Create a test-mode extension that provides an API endpoint for control.

## Recommendation
Install Playwright MCP for full browser control while keeping Browser Tools MCP for monitoring. This gives us:
- Full navigation control
- Screenshot capabilities
- Console log monitoring
- Network inspection
- Complete testing solution