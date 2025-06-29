# Auto-Reload Solution for Remote Testing

## The Problem
- Browser Tools MCP can capture screenshots but cannot navigate/reload pages
- You're working remotely via SSH and can't physically reload the browser
- We need the browser to show fresh content, not cached pages

## Solutions

### 1. Add Auto-Refresh Meta Tag (Quick Fix)
Add to dashboard template:
```html
<meta http-equiv="refresh" content="30">
```
This will auto-reload every 30 seconds.

### 2. HTMX Auto-Reload (Better)
Since we're using HTMX, add this to trigger full page reload:
```html
<div hx-get="/dashboard" hx-trigger="every 60s" hx-target="body" hx-swap="outerHTML"></div>
```

### 3. JavaScript Auto-Reload (Most Control)
```html
<script>
// Reload page every 60 seconds
setInterval(() => {
    window.location.reload();
}, 60000);
</script>
```

### 4. Force HTMX to Skip Cache
Modify our stats endpoint to include cache-busting:
```go
c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
c.Header("Pragma", "no-cache")
c.Header("Expires", "0")
```

### 5. Browser Automation Alternative
If we need full browser control, we could set up:
- Selenium WebDriver
- Playwright
- Puppeteer

These can navigate, click, reload programmatically.

## Immediate Solution
Let me add cache headers to force fresh content on HTMX updates.