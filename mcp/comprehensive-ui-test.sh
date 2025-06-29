#!/bin/bash

echo "E173 Gateway Comprehensive UI Testing Script"
echo "==========================================="
echo ""
echo "This script will test all UI components using:"
echo "- Browser Use MCP (for navigation and interaction)"
echo "- Browser Tools MCP (for monitoring and debugging)"
echo ""

# Set display for browser
export DISPLAY=:0

# Base URL
BASE_URL="http://192.168.1.35:8080"

# Test results directory
RESULTS_DIR="/root/e173_go_gateway/mcp/test-results"
mkdir -p "$RESULTS_DIR"

echo "Starting comprehensive UI tests..."
echo "Base URL: $BASE_URL"
echo "Results will be saved to: $RESULTS_DIR"
echo ""

# Function to run browser use command
run_browser_use() {
    local task="$1"
    echo "Browser Use: $task"
    cd /root/e173_go_gateway/mcp/browser-use
    source .venv/bin/activate
    # Here we would send the task to the MCP server
    # For now, we'll document what needs to be tested
}

# Function to run browser tools command
run_browser_tools() {
    local tool="$1"
    echo "Browser Tools: $tool"
    cd /root/e173_go_gateway/mcp/browser-tools
    # Here we would call the browser tools MCP
}

echo "Test Plan:"
echo "=========="
echo ""
echo "1. Login Page Tests"
echo "   - Navigate to login page"
echo "   - Check console errors"
echo "   - Test login with admin/admin"
echo "   - Verify redirect to dashboard"
echo ""
echo "2. Dashboard Tests"
echo "   - Verify 5 stat cards display horizontally"
echo "   - Check spacing between cards and panels"
echo "   - Monitor HTMX updates"
echo "   - Check console for errors"
echo ""
echo "3. Customers Page Tests"
echo "   - List customers"
echo "   - Test create customer form"
echo "   - Test edit customer (check auth redirect issue)"
echo "   - Test delete customer"
echo ""
echo "4. Gateways Page Tests"
echo "   - Verify gateway list loads"
echo "   - Test create gateway"
echo "   - Check for template errors"
echo ""
echo "5. Modems Page Tests"
echo "   - Check modem display (nested boxes issue)"
echo "   - Test modem controls"
echo ""
echo "6. CDR Page Tests"
echo "   - Verify empty table structure"
echo "   - Test with sample data if available"
echo ""
echo "7. Blacklist Page Tests"
echo "   - Test blacklist display"
echo "   - Test add/remove numbers"
echo ""

# Create test sequence file
cat > "$RESULTS_DIR/test-sequence.json" << 'EOF'
{
  "tests": [
    {
      "name": "Login Test",
      "steps": [
        {
          "tool": "browser_use",
          "action": "navigate",
          "url": "http://192.168.1.35:8080/login"
        },
        {
          "tool": "browser_tools",
          "action": "wipeLogs"
        },
        {
          "tool": "browser_tools",
          "action": "takeScreenshot",
          "save_as": "login-page.png"
        },
        {
          "tool": "browser_use",
          "action": "fill_and_submit",
          "task": "Fill username with 'admin', password with 'admin', and click login button"
        },
        {
          "tool": "browser_tools",
          "action": "getConsoleErrors"
        },
        {
          "tool": "browser_tools",
          "action": "takeScreenshot",
          "save_as": "after-login.png"
        }
      ]
    },
    {
      "name": "Dashboard Test",
      "steps": [
        {
          "tool": "browser_use",
          "action": "verify",
          "task": "Verify that 5 stat cards are displayed horizontally in one row"
        },
        {
          "tool": "browser_tools",
          "action": "getConsoleLogs"
        },
        {
          "tool": "browser_tools",
          "action": "getNetworkLogs"
        },
        {
          "tool": "browser_tools",
          "action": "takeScreenshot",
          "save_as": "dashboard-full.png"
        }
      ]
    }
  ]
}
EOF

echo ""
echo "Test sequence saved to: $RESULTS_DIR/test-sequence.json"
echo ""
echo "To run the tests:"
echo "1. Start Browser Use MCP server in one terminal"
echo "2. Start Browser Tools MCP server in another terminal"
echo "3. Execute the test sequence using both MCPs"
echo ""
echo "Next steps:"
echo "- Execute each test in sequence"
echo "- Document all errors found"
echo "- Create fixes for issues"
echo "- Update GitHub project tracker"