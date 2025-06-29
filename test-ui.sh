#!/bin/bash

# E173 Gateway UI Testing Script
# Tests all pages and collects results

echo "==============================================="
echo "E173 Gateway UI Testing"
echo "==============================================="
echo "Date: $(date)"
echo "==============================================="

BASE_URL="http://localhost:8080"
COOKIES_FILE="/tmp/e173-cookies.txt"

# Function to test a page
test_page() {
    local name=$1
    local path=$2
    local auth=$3
    
    echo -e "\n📍 Testing $name ($path)..."
    
    # Get page status
    if [ "$auth" = "true" ]; then
        status=$(curl -s -o /dev/null -w "%{http_code}" -b "$COOKIES_FILE" "$BASE_URL$path")
    else
        status=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$path")
    fi
    
    echo "   Status: $status"
    
    # Check for common patterns
    if [ "$auth" = "true" ]; then
        content=$(curl -s -b "$COOKIES_FILE" "$BASE_URL$path")
    else
        content=$(curl -s "$BASE_URL$path")
    fi
    
    # Check for errors
    errors=$(echo "$content" | grep -i "error" | wc -l)
    if [ $errors -gt 0 ]; then
        echo "   ⚠️  Found $errors error references"
    fi
    
    # Check for 404 in content
    if echo "$content" | grep -q "404"; then
        echo "   ❌ Page contains 404 error"
    fi
    
    # Check if page has content
    if [ ${#content} -lt 100 ]; then
        echo "   ❌ Page has very little content (${#content} bytes)"
    else
        echo "   ✅ Page loaded (${#content} bytes)"
    fi
}

# First, try to login
echo -e "\n🔐 Logging in..."
curl -s -c "$COOKIES_FILE" -X POST "$BASE_URL/login" \
    -d "username=admin&password=admin" \
    -o /dev/null

# Check if login was successful by trying to access dashboard
dashboard_status=$(curl -s -o /dev/null -w "%{http_code}" -b "$COOKIES_FILE" "$BASE_URL/")
if [ "$dashboard_status" = "200" ]; then
    echo "   ✅ Login successful"
else
    echo "   ❌ Login failed (status: $dashboard_status)"
fi

# Test all pages
test_page "Login Page" "/login" false
test_page "Dashboard" "/" true
test_page "Customers" "/customers" true
test_page "Gateways" "/gateways" true
test_page "Modems" "/modems" true
test_page "SIMs" "/sims" true
test_page "CDRs" "/cdrs" true
test_page "Blacklist" "/blacklist" true

# Test API endpoints
echo -e "\n📡 Testing API Endpoints..."
api_endpoints=(
    "/api/v1/stats/cards"
    "/api/v1/modems"
    "/api/v1/simcards"
    "/api/v1/gateways"
    "/api/v1/customers"
)

for endpoint in "${api_endpoints[@]}"; do
    status=$(curl -s -o /dev/null -w "%{http_code}" -b "$COOKIES_FILE" "$BASE_URL$endpoint")
    echo "   $endpoint: $status"
done

# Clean up
rm -f "$COOKIES_FILE"

echo -e "\n==============================================="
echo "Testing complete!"
echo "==============================================="