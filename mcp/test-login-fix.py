#!/usr/bin/env python3
"""
Quick test to verify login fix
"""

from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page()
    
    # Test login
    print("Testing login redirect fix...")
    page.goto("http://192.168.1.35:8080/login")
    page.fill('input[name="username"]', "admin")
    page.fill('input[name="password"]', "admin")
    page.click('button[type="submit"]')
    page.wait_for_load_state("networkidle")
    
    # Check URL
    current_url = page.url
    print(f"After login URL: {current_url}")
    
    # Check if we're on dashboard
    if current_url.endswith("/") or "dashboard" in page.content().lower():
        print("✅ Login redirect fixed!")
        
        # Check stat cards
        stat_cards = page.locator('.bg-white.dark\\:bg-gray-800.rounded-lg.shadow.p-6').count()
        print(f"✅ Found {stat_cards} stat cards")
    else:
        print("❌ Login redirect still broken")
    
    browser.close()