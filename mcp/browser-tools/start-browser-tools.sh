#!/bin/bash

echo "Starting Browser Tools MCP Server..."
echo "===================================="

# Start the browser tools server in the background
echo "Starting Browser Tools Server on port 3000..."
npx @agentdeskai/browser-tools-server@latest &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo "Browser Tools Server started with PID: $SERVER_PID"
echo ""
echo "⚠️  IMPORTANT: Install the Chrome Extension"
echo "   Download from: https://github.com/AgentDeskAI/browser-tools-mcp/releases"
echo ""
echo "To stop the server, run: kill $SERVER_PID"
echo ""
echo "Server is ready for UI testing!"

# Keep the script running
wait $SERVER_PID
