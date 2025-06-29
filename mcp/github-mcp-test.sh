#!/bin/bash

# Test the GitHub MCP server capabilities

source ../.github_config

echo "Testing GitHub MCP Server Connection..."
echo "======================================="

# Test 1: Get authenticated user
echo -e "\n1. Testing authenticated user info..."
echo '{"jsonrpc": "2.0", "method": "call_tool", "params": {"name": "get_me", "arguments": {}}, "id": 1}' | \
docker run -i --rm \
  -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN \
  ghcr.io/github/github-mcp-server stdio

# Test 2: Get repository info
echo -e "\n\n2. Testing repository access..."
echo '{"jsonrpc": "2.0", "method": "call_tool", "params": {"name": "get_repository", "arguments": {"owner": "Mrspoq", "name": "e173-go-gateway"}}, "id": 2}' | \
docker run -i --rm \
  -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN \
  ghcr.io/github/github-mcp-server stdio

# Test 3: List issues
echo -e "\n\n3. Testing issue listing..."
echo '{"jsonrpc": "2.0", "method": "call_tool", "params": {"name": "search_issues_and_pull_requests", "arguments": {"repository": "Mrspoq/e173-go-gateway", "query": "is:issue is:open"}}, "id": 3}' | \
docker run -i --rm \
  -e GITHUB_PERSONAL_ACCESS_TOKEN=$GITHUB_TOKEN \
  ghcr.io/github/github-mcp-server stdio