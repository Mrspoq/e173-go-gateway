#!/usr/bin/env python3

import requests
import json
import time
import base64
from datetime import datetime

class BrowserAutomation:
    def __init__(self, base_url="http://localhost:3001"):
        self.base_url = base_url
        self.app_url = "http://localhost:8080"
        
    def navigate(self, url):
        """Navigate to a URL"""
        response = requests.post(f"{self.base_url}/navigate", 
                                json={"url": url})
        return response.json()
    
    def screenshot(self, filename=None, full_page=False):
        """Take a screenshot"""
        if not filename:
            filename = f"screenshot_{datetime.now().strftime('%Y%m%d_%H%M%S')}.png"
        response = requests.post(f"{self.base_url}/screenshot", 
                                json={"filename": filename, "fullPage": full_page})
        return response.json()
    
    def click(self, selector):
        """Click an element"""
        response = requests.post(f"{self.base_url}/click", 
                                json={"selector": selector})
        return response.json()
    
    def type_text(self, selector, text):
        """Type text into an element"""
        response = requests.post(f"{self.base_url}/type", 
                                json={"selector": selector, "text": text})
        return response.json()
    
    def get_text(self, selector):
        """Get text from an element"""
        response = requests.post(f"{self.base_url}/getText", 
                                json={"selector": selector})
        return response.json()
    
    def exists(self, selector):
        """Check if element exists"""
        response = requests.post(f"{self.base_url}/exists", 
                                json={"selector": selector})
        return response.json()
    
    def wait_for(self, selector, timeout=5000):
        """Wait for element to appear"""
        response = requests.post(f"{self.base_url}/waitFor", 
                                json={"selector": selector, "timeout": timeout})
        return response.json()
    
    def analyze_ui(self):
        """Analyze UI for issues"""
        response = requests.post(f"{self.base_url}/analyzeUI")
        return response.json()
    
    def test_login(self):
        """Test login functionality"""
        print("\n🔐 Testing Login...")
        
        # Navigate to login page
        print("- Navigating to login page...")
        self.navigate(f"{self.app_url}/login")
        time.sleep(1)
        
        # Take screenshot
        self.screenshot("login_page.png")
        
        # Type credentials
        print("- Entering credentials...")
        self.type_text("#username", "admin")
        self.type_text("#password", "admin")
        
        # Click login button
        print("- Clicking login button...")
        self.click('button[type="submit"]')
        time.sleep(2)
        
        # Check if we're redirected
        result = self.analyze_ui()
        print(f"- Current page: {result['analysis']['title']}")
        
        # Take screenshot after login
        self.screenshot("after_login.png")
        
        return result
    
    def test_all_pages(self):
        """Test all main pages"""
        pages = [
            ("/", "Dashboard"),
            ("/gateways", "Gateways"),
            ("/modems", "Modems"),
            ("/sims", "SIM Cards"),
            ("/customers", "Customers"),
            ("/cdrs", "Call Records"),
            ("/blacklist", "Blacklist"),
            ("/settings", "Settings")
        ]
        
        results = {}
        
        for path, name in pages:
            print(f"\n📄 Testing {name} page...")
            self.navigate(f"{self.app_url}{path}")
            time.sleep(2)
            
            # Take screenshot
            screenshot = self.screenshot(f"{name.lower().replace(' ', '_')}.png")
            
            # Analyze UI
            analysis = self.analyze_ui()
            results[name] = analysis
            
            if analysis['success'] and analysis['analysis']['issues']:
                print(f"  ⚠️  Found {len(analysis['analysis']['issues'])} issues:")
                for issue in analysis['analysis']['issues']:
                    print(f"    - {issue['type']}: {issue['message']}")
            else:
                print(f"  ✅ No issues found")
        
        return results

def main():
    print("🚀 Starting E173 Gateway UI Testing")
    print("=" * 50)
    
    # Wait for MCP server to be ready
    time.sleep(2)
    
    browser = BrowserAutomation()
    
    try:
        # Test login
        login_result = browser.test_login()
        
        if login_result['success']:
            print("\n✅ Login successful!")
            
            # Test all pages
            print("\n🔍 Testing all pages...")
            results = browser.test_all_pages()
            
            # Generate report
            print("\n📊 Test Summary")
            print("=" * 50)
            
            total_issues = 0
            for page, result in results.items():
                if result['success']:
                    issues = result['analysis']['issues']
                    total_issues += len(issues)
                    print(f"\n{page}:")
                    print(f"  - Issues: {len(issues)}")
                    print(f"  - Stats: {result['analysis']['stats']}")
            
            print(f"\n🎯 Total issues found: {total_issues}")
            
            # Save detailed report
            with open('/root/e173_go_gateway/UI_TEST_REPORT.json', 'w') as f:
                json.dump({
                    'timestamp': datetime.now().isoformat(),
                    'results': results,
                    'total_issues': total_issues
                }, f, indent=2)
            
            print("\n📝 Detailed report saved to UI_TEST_REPORT.json")
            
        else:
            print("\n❌ Login failed!")
            print(login_result)
            
    except Exception as e:
        print(f"\n❌ Error during testing: {e}")

if __name__ == "__main__":
    main()