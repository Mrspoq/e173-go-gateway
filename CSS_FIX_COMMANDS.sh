#!/bin/bash

# CSS Fix Script - Execute immediately when you wake up
echo "🎨 Fixing CSS loading issues for E173 Gateway..."

cd /root/e173_go_gateway

# 1. Rebuild CSS assets
echo "📦 Rebuilding Tailwind CSS..."
npm install
npm run build-css

# 2. Verify CSS file exists and has content
echo "✅ Verifying CSS bundle..."
if [ -f "web/static/bundle.css" ]; then
    SIZE=$(wc -c < web/static/bundle.css)
    echo "CSS bundle size: $SIZE bytes"
    if [ $SIZE -gt 1000 ]; then
        echo "✅ CSS bundle looks good!"
    else
        echo "❌ CSS bundle too small, regenerating..."
        npm run build-css --verbose
    fi
else
    echo "❌ CSS bundle missing, creating..."
    npm run build-css
fi

# 3. Check if server can serve static files
echo "🌐 Testing static file serving..."
if pgrep -f "e173gw" > /dev/null; then
    echo "Server is running, testing CSS endpoint..."
    curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/static/bundle.css
else
    echo "Starting server for testing..."
    make run &
    sleep 3
    curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/static/bundle.css
fi

# 4. Backup and fix potential path issues
echo "🔧 Checking static file configuration..."

# Create backup of main.go
cp cmd/server/main.go cmd/server/main.go.backup

# Check if static route is correctly configured
if grep -q 'router.Static("/static", "./web/static")' cmd/server/main.go; then
    echo "✅ Static route configuration looks correct"
else
    echo "❌ Static route might be misconfigured"
    echo "Expected: router.Static(\"/static\", \"./web/static\")"
    grep -n "Static" cmd/server/main.go
fi

# 5. Test with curl and show response
echo "🧪 Full CSS test..."
echo "CSS Response Headers:"
curl -s -I http://localhost:8080/static/bundle.css

echo "CSS Content Preview:"
curl -s http://localhost:8080/static/bundle.css | head -20

# 6. Check template CSS link
echo "🔍 Checking template CSS references..."
grep -n "bundle.css" templates/base.tmpl

echo "🎯 CSS fix complete! Check browser at http://localhost:8080"
echo "If styles still not loading, check browser dev tools console for errors."
