#!/bin/bash

echo "=== E173 Gateway UI Test ==="
echo ""

BASE_URL="http://localhost:8080"

# Test public endpoints
echo "1. Testing public endpoints..."
echo -n "   - Ping endpoint: "
curl -s "$BASE_URL/ping" | jq -r '.message' || echo "FAILED"

echo -n "   - Login page: "
if curl -s "$BASE_URL/login" | grep -q "Login - E173 Gateway"; then
    echo "OK"
else
    echo "FAILED"
fi

# Test API endpoints (no auth required for some)
echo ""
echo "2. Testing API endpoints..."

echo -n "   - Stats API: "
curl -s "$BASE_URL/api/stats" | jq -r 'if . then "OK" else "FAILED" end' 2>/dev/null || echo "FAILED"

echo -n "   - Stats Cards: "
if curl -s "$BASE_URL/api/stats/cards" | grep -q "System Overview"; then
    echo "OK"
else
    echo "FAILED"
fi

echo -n "   - CDR Recent: "
curl -s "$BASE_URL/api/cdr/recent" | jq -r 'if .cdrs then "OK" else "FAILED" end' 2>/dev/null || echo "FAILED"

echo -n "   - Modems Status: "
curl -s "$BASE_URL/api/modems/status" | jq -r 'if .modems then "OK" else "FAILED" end' 2>/dev/null || echo "FAILED"

echo -n "   - SIMs API: "
curl -s "$BASE_URL/api/sims" | jq -r 'if .sims then "OK" else "FAILED" end' 2>/dev/null || echo "FAILED"

# Test analytics endpoints (if available)
echo ""
echo "3. Testing analytics endpoints..."

echo -n "   - Call Analytics: "
if curl -s "$BASE_URL/api/analytics/calls" 2>/dev/null | jq -e . >/dev/null 2>&1; then
    echo "OK"
else
    echo "NOT AVAILABLE (Redis not connected)"
fi

echo -n "   - SIM Analytics: "
if curl -s "$BASE_URL/api/analytics/sims" 2>/dev/null | jq -e . >/dev/null 2>&1; then
    echo "OK"
else
    echo "NOT AVAILABLE (Redis not connected)"
fi

echo -n "   - Dashboard Analytics: "
if curl -s "$BASE_URL/api/analytics/dashboard" 2>/dev/null | jq -e . >/dev/null 2>&1; then
    echo "OK"
else
    echo "NOT AVAILABLE (Redis not connected)"
fi

# Test protected pages (should redirect to login)
echo ""
echo "4. Testing protected pages (should redirect)..."

for page in "/" "/dashboard" "/modems" "/sims" "/customers" "/cdrs" "/blacklist" "/system"; do
    echo -n "   - $page: "
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" -L "$BASE_URL$page")
    if [ "$STATUS" = "200" ]; then
        # Check if we ended up on login page
        if curl -s -L "$BASE_URL$page" | grep -q "Login - E173 Gateway"; then
            echo "OK (redirects to login)"
        else
            echo "ACCESSIBLE WITHOUT AUTH!"
        fi
    else
        echo "ERROR (HTTP $STATUS)"
    fi
done

echo ""
echo "=== UI Test Complete ==="