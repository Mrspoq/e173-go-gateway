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
