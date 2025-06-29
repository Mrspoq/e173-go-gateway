#!/usr/bin/env python3
"""
Check for console errors on all pages
"""

from playwright.sync_api import sync_playwright

def check_all_pages():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        
        # Track console errors
        console_errors = []
        def handle_console(msg):
            if msg.type in ['error', 'warning']:
                console_errors.append({
                    'page': page.url,
                    'type': msg.type,
                    'text': msg.text
                })
        
        page.on("console", handle_console)
        
        # Login first
        page.goto("http://192.168.1.35:8080/login")
        page.fill('input[name="username"]', "admin")
        page.fill('input[name="password"]', "admin")
        page.click('button[type="submit"]')
        page.wait_for_load_state("networkidle")
        
        # Check all pages
        pages = [
            ("/", "Dashboard"),
            ("/customers", "Customers"),
            ("/gateways", "Gateways"),
            ("/modems", "Modems"),
            ("/sims", "SIM Cards"),
            ("/cdrs", "CDR"),
            ("/blacklist", "Blacklist"),
            ("/settings-new", "Settings")
        ]
        
        for path, name in pages:
            print(f"Checking {name}...")
            page.goto(f"http://192.168.1.35:8080{path}")
            page.wait_for_load_state("networkidle")
            
            # Check for data
            content = page.content()
            if "No data" in content or "No items" in content:
                print(f"  ⚠️  {name} has no data")
        
        browser.close()
        
        # Report errors
        if console_errors:
            print("\n❌ Console Errors Found:")
            for error in console_errors:
                print(f"  {error['type']}: {error['text']}")
                print(f"    on: {error['page']}")
        else:
            print("\n✅ No console errors found!")

if __name__ == "__main__":
    check_all_pages()