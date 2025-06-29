#!/bin/bash

echo "=================================================="
echo "E173 Gateway UI Testing with Browser Tools MCP"
echo "=================================================="
echo ""

# Check if Chrome is installed
if ! which google-chrome > /dev/null; then
    echo "❌ Google Chrome is not installed!"
    echo "Run: ./install-chrome.sh"
    exit 1
fi

echo "✅ Chrome installed: $(google-chrome --version)"
echo ""

# Function to download Chrome extension
download_extension() {
    echo "📥 Downloading Browser Tools Chrome Extension..."
    mkdir -p chrome-extension
    cd chrome-extension
    
    # Get latest release
    LATEST_RELEASE=$(curl -s https://api.github.com/repos/AgentDeskAI/browser-tools-mcp/releases/latest)
    DOWNLOAD_URL=$(echo $LATEST_RELEASE | grep -o '"browser_url":"[^"]*' | cut -d'"' -f4 | head -1)
    
    if [ -z "$DOWNLOAD_URL" ]; then
        echo "⚠️  Could not find extension download URL"
        echo "Please download manually from: https://github.com/AgentDeskAI/browser-tools-mcp/releases"
    else
        wget -O browser-tools-extension.zip "$DOWNLOAD_URL"
        unzip -o browser-tools-extension.zip
        echo "✅ Extension downloaded to ./chrome-extension/"
    fi
    
    cd ..
}

# Check if extension needs to be downloaded
if [ ! -d "chrome-extension" ]; then
    download_extension
fi

echo ""
echo "Starting Browser Tools Server..."
echo "================================"

# Kill any existing browser tools server
pkill -f "browser-tools-server" 2>/dev/null

# Start the browser tools server
npx @agentdeskai/browser-tools-server@latest &
SERVER_PID=$!

echo "✅ Browser Tools Server started (PID: $SERVER_PID)"
echo ""

# Wait for server to be ready
sleep 3

# Create a Chrome launch script
cat > launch-chrome-with-extension.sh << 'EOF'
#!/bin/bash

# Launch Chrome with the Browser Tools extension
google-chrome \
    --load-extension="$(pwd)/chrome-extension" \
    --no-first-run \
    --disable-default-apps \
    --disable-popup-blocking \
    --disable-translate \
    --disable-sync \
    --no-default-browser-check \
    --window-size=1920,1080 \
    --window-position=0,0 \
    http://192.168.1.35:8080 &

CHROME_PID=$!
echo "✅ Chrome launched (PID: $CHROME_PID)"
echo ""
echo "Please complete the following steps:"
echo "1. Login to E173 Gateway (admin/admin123)"
echo "2. Open Chrome DevTools (F12)"
echo "3. Navigate to 'Browser Tools' tab"
echo "4. Click 'Connect' to start monitoring"
echo ""
echo "The extension will capture:"
echo "- Console logs and errors"
echo "- Network requests"
echo "- DOM elements"
echo "- Screenshots"
echo ""
EOF

chmod +x launch-chrome-with-extension.sh

echo "=================================================="
echo "🚀 UI Testing Environment Ready!"
echo "=================================================="
echo ""
echo "Options:"
echo "1. Launch Chrome with extension: ./launch-chrome-with-extension.sh"
echo "2. Run automated tests: npm run test-ui"
echo "3. Start MCP server: npx @agentdeskai/browser-tools-mcp@latest"
echo ""
echo "Current Status:"
echo "- Browser Tools Server: Running on http://localhost:3000"
echo "- E173 Gateway: http://192.168.1.35:8080"
echo ""
echo "To stop the server: kill $SERVER_PID"
echo ""

# Create test report script
cat > generate-test-report.sh << 'EOF'
#!/bin/bash

echo "Generating UI Test Report..."
echo "============================"

REPORT_FILE="ui-test-report-$(date +%Y%m%d-%H%M%S).md"

cat > $REPORT_FILE << 'REPORT'
# E173 Gateway UI Test Report

Generated: $(date)

## Test Environment
- Browser: Google Chrome
- URL: http://192.168.1.35:8080
- Tool: Browser Tools MCP

## Test Results

### 1. Dashboard Layout
- [ ] Shows 5 cards in one row
- [ ] Cards load via HTMX
- [ ] No console errors
- [ ] Stats update every 5 seconds

### 2. Gateway Management
- [ ] Gateway list displays correctly
- [ ] Add gateway button works
- [ ] Test connection button functional
- [ ] No authentication errors

### 3. Modems Page
- [ ] No nested boxes
- [ ] Modem cards display correctly
- [ ] Status indicators working

### 4. Customer Management
- [ ] Add customer button navigates correctly
- [ ] Edit buttons work
- [ ] No redirect to login
- [ ] Forms load properly

### 5. CDR Display
- [ ] Table structure shows when empty
- [ ] Headers display correctly
- [ ] Pagination controls visible

### 6. Authentication
- [ ] User info displays in navbar
- [ ] Session persists across pages
- [ ] Logout works correctly

## Console Errors
[Paste any console errors here]

## Network Issues
[Paste any failed API calls here]

## Screenshots
[Reference to captured screenshots]

## Recommendations
[Add specific fixes needed]

REPORT

echo "✅ Report template created: $REPORT_FILE"
echo "Please fill in the test results after running tests"
EOF

chmod +x generate-test-report.sh

echo "To generate a test report: ./generate-test-report.sh"
echo ""

# Keep the script running
trap "kill $SERVER_PID 2>/dev/null; exit" INT TERM
wait $SERVER_PID