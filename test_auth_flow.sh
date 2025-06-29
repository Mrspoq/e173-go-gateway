#!/bin/bash

echo "Testing Authentication Flow"
echo "============================"
echo ""

# Test 1: Access dashboard without authentication (should redirect to login)
echo "1. Testing access to dashboard without authentication:"
curl -s -o /dev/null -w "%{http_code} %{redirect_url}\n" http://localhost:8080/

echo ""
echo "2. Testing login page access:"
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/login

echo ""
echo "3. Testing login with credentials (using form data):"
# This would normally return a session cookie
curl -s -X POST -d "username=admin&password=admin" \
     -c cookies.txt \
     -w "\nHTTP Status: %{http_code}\n" \
     http://localhost:8080/login

echo ""
echo "4. Testing access to dashboard with session cookie:"
curl -s -b cookies.txt -o /dev/null -w "%{http_code}\n" http://localhost:8080/

echo ""
echo "5. Testing logout:"
curl -s -X POST -b cookies.txt -w "%{http_code}\n" http://localhost:8080/logout

# Clean up
rm -f cookies.txt

echo ""
echo "Test complete!"