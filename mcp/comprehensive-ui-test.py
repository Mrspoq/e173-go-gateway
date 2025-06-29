#!/usr/bin/env python3
"""
Comprehensive UI Test for E173 Gateway
Tests all pages and captures issues
"""

import os
import json
from datetime import datetime
from playwright.sync_api import sync_playwright

class E173UITester:
    def __init__(self):
        self.base_url = "http://192.168.1.35:8080"
        self.results = {
            "timestamp": datetime.now().isoformat(),
            "tests": [],
            "errors": [],
            "console_logs": []
        }
        
    def test_all_pages(self):
        with sync_playwright() as p:
            browser = p.chromium.launch(headless=True)
            context = browser.new_context()
            page = context.new_page()
            
            # Set up monitoring
            page.on("console", lambda msg: self.results["console_logs"].append({
                "type": msg.type,
                "text": msg.text,
                "url": page.url
            }))
            
            page.on("pageerror", lambda error: self.results["errors"].append({
                "error": str(error),
                "url": page.url
            }))
            
            # Test login
            print("Testing Login...")
            self.test_login(page)
            
            # Test all pages
            pages_to_test = [
                ("/dashboard", "Dashboard"),
                ("/customers", "Customers"),
                ("/gateways", "Gateways"),
                ("/modems", "Modems"),
                ("/sims", "SIM Cards"),
                ("/cdrs", "CDR"),
                ("/blacklist", "Blacklist"),
                ("/settings-new", "Settings")
            ]
            
            for path, name in pages_to_test:
                print(f"Testing {name}...")
                self.test_page(page, path, name)
            
            browser.close()
            
            # Save results
            with open("/root/e173_go_gateway/mcp/test-results/comprehensive_test_results.json", "w") as f:
                json.dump(self.results, f, indent=2)
            
            self.print_summary()
    
    def test_login(self, page):
        page.goto(f"{self.base_url}/login")
        page.fill('input[name="username"]', "admin")
        page.fill('input[name="password"]', "admin")
        page.click('button[type="submit"]')
        page.wait_for_load_state("networkidle")
        
        self.results["tests"].append({
            "page": "Login",
            "status": "success" if "dashboard" in page.url else "failed",
            "url": page.url
        })
    
    def test_page(self, page, path, name):
        try:
            page.goto(f"{self.base_url}{path}")
            page.wait_for_load_state("networkidle")
            
            # Take screenshot
            screenshot_path = f"/root/e173_go_gateway/mcp/browser-use/screenshots/{name.lower().replace(' ', '_')}.png"
            page.screenshot(path=screenshot_path)
            
            # Check for common issues
            issues = []
            
            # Check if page loaded
            if page.title() == "" or "404" in page.title():
                issues.append("Page not found or empty")
            
            # Check for empty content
            content = page.content()
            if "No data" in content or "No items" in content:
                issues.append("Empty data display")
            
            # Check for forms
            forms = page.locator("form").count()
            buttons = page.locator("button").count()
            
            self.results["tests"].append({
                "page": name,
                "path": path,
                "status": "issues" if issues else "success",
                "issues": issues,
                "forms": forms,
                "buttons": buttons,
                "screenshot": screenshot_path
            })
            
        except Exception as e:
            self.results["tests"].append({
                "page": name,
                "path": path,
                "status": "error",
                "error": str(e)
            })
    
    def print_summary(self):
        print("\n" + "="*50)
        print("UI TEST SUMMARY")
        print("="*50)
        
        for test in self.results["tests"]:
            status_icon = "✅" if test["status"] == "success" else "❌"
            print(f"{status_icon} {test['page']}: {test['status']}")
            if "issues" in test and test["issues"]:
                for issue in test["issues"]:
                    print(f"   - {issue}")
        
        print(f"\nConsole logs: {len(self.results['console_logs'])}")
        print(f"Page errors: {len(self.results['errors'])}")

if __name__ == "__main__":
    tester = E173UITester()
    tester.test_all_pages()