# Screenshot Guide for UI Testing

## For Best Screenshots:

1. **Close or minimize Chrome DevTools** (F12 to toggle)
2. **Ensure Browser Tools extension is connected**:
   - Open DevTools briefly
   - Go to "Browser Tools" tab
   - Click "Connect" if needed
   - Close DevTools again

3. **Navigate to the page** you want to screenshot

4. **To capture full page**:
   - Dashboard: http://192.168.1.35:8080/dashboard
   - Scroll to show all 5 stat cards (if they exist)
   - Or take multiple screenshots at different scroll positions

## What I Need to See:

### Dashboard Page:
- [ ] All stat cards (should be 5 in one row)
- [ ] The layout below the cards
- [ ] Any console errors in Browser Tools

### Other Pages to Test:
- [ ] Gateway page (/gateways)
- [ ] Customer page (/customers)
- [ ] CDR page (/cdrs)

## Quick Screenshot Command:
```bash
cd /root/e173_go_gateway/mcp/browser-tools
./take-screenshot.sh
```

## If Extension Disconnects:
1. Refresh the page in Chrome
2. Open DevTools → Browser Tools tab
3. Click "Connect"
4. Try screenshot again