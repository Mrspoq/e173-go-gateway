# Current Browser Testing Issues

## Problem Statement
We need to reload/navigate browser pages remotely via SSH for accurate UI testing, but Browser Tools MCP only monitors, it doesn't control.

## Options Available:

### 1. Playwright MCP Server (if not already set up)
- Wraps existing Playwright in MCP protocol
- Provides navigate, reload, click functions
- Works alongside Browser Tools MCP

### 2. Use Existing Playwright Directly
- If Playwright is already installed
- Create scripts to control browser
- Integrate with our testing workflow

### 3. Simple xdotool Solution
```bash
# Reload current page
DISPLAY=:0 xdotool key F5

# Navigate to URL
DISPLAY=:0 xdotool key ctrl+l
DISPLAY=:0 xdotool type "http://192.168.1.35:8080/dashboard"
DISPLAY=:0 xdotool key Return
```

### 4. Browser Extension Bridge
- Create extension that listens for commands
- Send commands via local API
- Browser executes navigation/reload

## Dashboard Layout Issue
Despite multiple fixes, dashboard still shows cards vertically:
- Tried: grid-cols-5
- Tried: flexbox with inline styles
- Issue: Cards still display 1 per row
- Possible causes:
  - Tailwind CSS not including grid utilities
  - Parent container constraints
  - HTMX replacing wrong element
  - CSS specificity issues

## What We Know Works:
- Server returns correct HTML with grid/flex
- Browser Tools MCP captures screenshots successfully
- HTMX updates are happening
- Individual stat cards load correctly

## What Doesn't Work:
- Cards won't display horizontally
- Can't reload browser remotely
- Can't navigate to test other pages