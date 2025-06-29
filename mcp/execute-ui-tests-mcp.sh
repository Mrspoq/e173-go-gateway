#!/bin/bash

echo "E173 Gateway UI Testing with MCP Servers"
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Base paths
BROWSER_USE_PATH="/root/e173_go_gateway/mcp/browser-use"
BROWSER_TOOLS_PATH="/root/e173_go_gateway/mcp/browser-tools"
RESULTS_PATH="/root/e173_go_gateway/mcp/test-results"

# Create results directory
mkdir -p "$RESULTS_PATH"

echo -e "${YELLOW}Starting MCP servers...${NC}"
echo ""

# Start Browser Use MCP in background
echo "1. Starting Browser Use MCP server..."
cd "$BROWSER_USE_PATH"
source .venv/bin/activate
export DISPLAY=:0
export MCP_HEADLESS=false
nohup mcp-server-browser-use > "$RESULTS_PATH/browser-use-mcp.log" 2>&1 &
BROWSER_USE_PID=$!
echo -e "${GREEN}Browser Use MCP started (PID: $BROWSER_USE_PID)${NC}"

# Give it time to start
sleep 3

# Start Browser Tools MCP in background
echo "2. Starting Browser Tools MCP server..."
cd "$BROWSER_TOOLS_PATH"
nohup npx @agentdeskai/browser-tools-mcp@latest > "$RESULTS_PATH/browser-tools-mcp.log" 2>&1 &
BROWSER_TOOLS_PID=$!
echo -e "${GREEN}Browser Tools MCP started (PID: $BROWSER_TOOLS_PID)${NC}"

# Give it time to start
sleep 3

echo ""
echo -e "${YELLOW}MCP servers are running!${NC}"
echo ""
echo "Logs:"
echo "- Browser Use: $RESULTS_PATH/browser-use-mcp.log"
echo "- Browser Tools: $RESULTS_PATH/browser-tools-mcp.log"
echo ""

# Create test execution plan
cat > "$RESULTS_PATH/test-plan.md" << 'EOF'
# E173 Gateway UI Test Execution Plan

## Test 1: Login Page
1. Browser Use: Navigate to login page
2. Browser Tools: Clear console logs
3. Browser Tools: Take screenshot of login page
4. Browser Use: Enter username 'admin' and password 'admin'
5. Browser Use: Click login button
6. Browser Tools: Check for console errors
7. Browser Tools: Take screenshot after login

## Test 2: Dashboard
1. Browser Use: Verify we're on dashboard
2. Browser Use: Count stat cards
3. Browser Use: Check spacing between elements
4. Browser Tools: Monitor console logs
5. Browser Tools: Check network requests
6. Browser Tools: Take full page screenshot

## Test 3: Customers Page
1. Browser Use: Navigate to customers
2. Browser Tools: Monitor for errors
3. Browser Use: Click create customer
4. Browser Use: Test edit button
5. Browser Tools: Document any auth issues

## Test 4: Gateways Page
1. Browser Use: Navigate to gateways
2. Browser Tools: Check for blank page
3. Browser Tools: Look for template errors

## Test 5: Modems Page
1. Browser Use: Navigate to modems
2. Browser Use: Check for nested boxes issue
3. Browser Tools: Take screenshot

## Test 6: CDR Page
1. Browser Use: Navigate to CDR
2. Browser Use: Verify empty table shows
3. Browser Tools: Check console

## Test 7: Blacklist Page
1. Browser Use: Navigate to blacklist
2. Browser Use: Test add/remove functionality
3. Browser Tools: Final error check
EOF

echo -e "${GREEN}Test plan created at: $RESULTS_PATH/test-plan.md${NC}"
echo ""

# Create MCP command examples
cat > "$RESULTS_PATH/mcp-commands.json" << 'EOF'
{
  "login_test": {
    "step1": {
      "server": "browser-use",
      "method": "tools/call",
      "params": {
        "name": "run_browser_agent",
        "arguments": {
          "task": "Navigate to http://192.168.1.35:8080/login"
        }
      }
    },
    "step2": {
      "server": "browser-tools",
      "method": "tools/call",
      "params": {
        "name": "wipeLogs",
        "arguments": {}
      }
    },
    "step3": {
      "server": "browser-tools",
      "method": "tools/call",
      "params": {
        "name": "takeScreenshot",
        "arguments": {}
      }
    },
    "step4": {
      "server": "browser-use",
      "method": "tools/call",
      "params": {
        "name": "run_browser_agent",
        "arguments": {
          "task": "Fill the username field with 'admin', password field with 'admin', and click the login button"
        }
      }
    }
  },
  "dashboard_test": {
    "step1": {
      "server": "browser-use",
      "method": "tools/call",
      "params": {
        "name": "run_browser_agent",
        "arguments": {
          "task": "Count the number of stat cards in the top row. They should be 5 cards displayed horizontally. Also check if there is proper spacing between the stat cards row and the two panels below (Live Modem Status and Recent Call Activity)."
        }
      }
    }
  }
}
EOF

echo -e "${GREEN}MCP command examples saved to: $RESULTS_PATH/mcp-commands.json${NC}"
echo ""

echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Use Claude or another MCP client to send commands to both servers"
echo "2. Execute the test plan step by step"
echo "3. Document all findings"
echo "4. Fix issues found"
echo "5. Update GitHub project tracker"
echo ""

echo "To stop the MCP servers later:"
echo "kill $BROWSER_USE_PID $BROWSER_TOOLS_PID"
echo ""

# Save PIDs for later
echo "$BROWSER_USE_PID" > "$RESULTS_PATH/browser-use.pid"
echo "$BROWSER_TOOLS_PID" > "$RESULTS_PATH/browser-tools.pid"

echo -e "${GREEN}Ready to start testing!${NC}"