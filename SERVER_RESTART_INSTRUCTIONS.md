# Server Restart Instructions

## When you return, please follow these steps to apply all UI fixes:

### 1. Stop any running servers
```bash
pkill -f server_
# or
ps aux | grep server_ | grep -v grep | awk '{print $2}' | xargs kill -9
```

### 2. Build the latest server with all fixes
```bash
cd /root/e173_go_gateway
go build -o server_latest ./cmd/server/main.go
```

### 3. Create/update the admin user (if needed)
```bash
go run scripts/create_admin_user.go
```

### 4. Start the server
```bash
./server_latest > server.log 2>&1 &
```

### 5. Verify server is running
```bash
# Check if server is listening
curl http://localhost:8080/ping

# Check logs
tail -f server.log
```

### 6. Test the UI fixes

1. **Dashboard** - http://localhost:8080/
   - Should show 5 cards (not 4)
   - No full page refresh
   - Loading indicators should work properly

2. **Authentication** - http://localhost:8080/login
   - Login with admin/admin
   - Should see "Welcome, Admin" in navigation
   - Logout button should be visible

3. **Gateways** - http://localhost:8080/gateways
   - Should show "Gateway Management" (not "Modem Management")
   - Gateway list should display properly

4. **Customers** - http://localhost:8080/customers
   - Stats should show as numbers (not JSON)
   - Customer list should load properly

5. **CDR/Blacklist** - http://localhost:8080/cdrs and http://localhost:8080/blacklist
   - Should NOT repeatedly call /api/v1/modems
   - Check browser network tab to verify

### 7. Run automated UI tests (optional)
```bash
cd /root/e173_go_gateway/mcp/browser-automation
node ui-tests.js
```

## Troubleshooting

If the server doesn't start:
1. Check for port conflicts: `lsof -i :8080`
2. Check PostgreSQL is running: `systemctl status postgresql`
3. Check Redis is running: `systemctl status redis`
4. Review server logs: `tail -100 server.log`

If templates don't load correctly:
1. Ensure all template files are in place
2. The modems/list.tmpl file should be removed (it was causing conflicts)
3. Check template loading errors in server.log

## Summary of Fixes Applied

✅ Dashboard grid layout (5 columns)
✅ Authentication with admin/admin
✅ User display in navigation
✅ Template collision fix (gateways/modems)
✅ HTMX polling cleanup
✅ Customer stats HTML response
✅ Empty state handling
✅ Loading indicators

All fixes are in place and just need a server restart to take effect!