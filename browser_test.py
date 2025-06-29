#!/usr/bin/env python3
"""
E173 Gateway Browser Testing Script
Tests all pages for UI issues, console errors, and 404s
"""

import os
import time
import json
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from webdriver_manager.chrome import ChromeDriverManager

class E173BrowserTest:
    def __init__(self):
        self.base_url = "http://localhost:8080"
        self.test_results = {
            "timestamp": datetime.now().isoformat(),
            "pages": {},
            "summary": {
                "total_errors": 0,
                "total_404s": 0,
                "total_console_errors": 0,
                "pages_with_issues": []
            }
        }
        self.setup_driver()
        
    def setup_driver(self):
        """Setup Chrome driver with appropriate options"""
        chrome_options = Options()
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--disable-dev-shm-usage')
        chrome_options.add_argument('--disable-gpu')
        chrome_options.add_argument('--window-size=1920,1080')
        
        # Enable logging
        chrome_options.set_capability('goog:loggingPrefs', {
            'browser': 'ALL',
            'performance': 'ALL'
        })
        
        # Create screenshots directory
        os.makedirs('screenshots', exist_ok=True)
        
        try:
            self.driver = webdriver.Chrome(
                service=Service(ChromeDriverManager().install()),
                options=chrome_options
            )
        except Exception as e:
            print(f"Failed to setup Chrome driver: {e}")
            print("Trying with system Chrome...")
            chrome_options.binary_location = "/usr/bin/google-chrome"
            self.driver = webdriver.Chrome(options=chrome_options)
            
    def capture_page_state(self, page_name):
        """Capture current page state including console logs and network errors"""
        state = {
            "url": self.driver.current_url,
            "title": self.driver.title,
            "console_errors": [],
            "network_errors": [],
            "javascript_errors": [],
            "missing_resources": []
        }
        
        # Get browser logs
        try:
            browser_logs = self.driver.get_log('browser')
            for log in browser_logs:
                if log['level'] == 'SEVERE':
                    state['console_errors'].append({
                        'message': log['message'],
                        'timestamp': log['timestamp']
                    })
        except Exception as e:
            print(f"Could not get browser logs: {e}")
        
        # Get performance logs to check for 404s
        try:
            performance_logs = self.driver.get_log('performance')
            for log in performance_logs:
                message = json.loads(log['message'])
                method = message.get('message', {}).get('method', '')
                
                if method == 'Network.responseReceived':
                    response = message['message']['params']['response']
                    if response['status'] == 404:
                        state['missing_resources'].append({
                            'url': response['url'],
                            'type': response.get('mimeType', 'unknown')
                        })
                    elif response['status'] >= 400:
                        state['network_errors'].append({
                            'url': response['url'],
                            'status': response['status'],
                            'statusText': response.get('statusText', '')
                        })
        except Exception as e:
            print(f"Could not get performance logs: {e}")
        
        # Check for JavaScript errors
        try:
            js_errors = self.driver.execute_script("""
                return window.jsErrors || [];
            """)
            if js_errors:
                state['javascript_errors'].extend(js_errors)
        except Exception as e:
            print(f"Could not check JS errors: {e}")
        
        # Take screenshot
        screenshot_path = f"screenshots/{page_name}_{datetime.now().strftime('%Y%m%d_%H%M%S')}.png"
        try:
            self.driver.save_screenshot(screenshot_path)
            state['screenshot'] = screenshot_path
        except Exception as e:
            print(f"Could not take screenshot: {e}")
            state['screenshot'] = None
        
        return state
    
    def test_login(self):
        """Test login functionality"""
        print("Testing login page...")
        
        # Navigate to login
        self.driver.get(f"{self.base_url}/login")
        time.sleep(2)
        
        # Capture login page state
        login_state = self.capture_page_state("login")
        self.test_results['pages']['login_page'] = login_state
        
        try:
            # Find and fill login form
            username_field = WebDriverWait(self.driver, 10).until(
                EC.presence_of_element_located((By.NAME, "username"))
            )
            password_field = self.driver.find_element(By.NAME, "password")
            
            username_field.send_keys("admin")
            password_field.send_keys("admin")
            
            # Submit form
            submit_button = self.driver.find_element(By.CSS_SELECTOR, "button[type='submit']")
            submit_button.click()
            
            # Wait for redirect
            time.sleep(3)
            
            # Capture post-login state
            post_login_state = self.capture_page_state("post_login")
            self.test_results['pages']['post_login'] = post_login_state
            
            return True
        except Exception as e:
            print(f"Login failed: {e}")
            self.test_results['pages']['login_error'] = str(e)
            return False
    
    def test_page(self, path, name):
        """Test a specific page"""
        print(f"Testing {name} page...")
        
        try:
            self.driver.get(f"{self.base_url}{path}")
            time.sleep(3)  # Wait for page to fully load
            
            # Capture page state
            page_state = self.capture_page_state(name)
            self.test_results['pages'][name] = page_state
            
            # Update summary
            if page_state['console_errors']:
                self.test_results['summary']['total_console_errors'] += len(page_state['console_errors'])
                self.test_results['summary']['pages_with_issues'].append(f"{name} (console errors)")
                
            if page_state['missing_resources']:
                self.test_results['summary']['total_404s'] += len(page_state['missing_resources'])
                self.test_results['summary']['pages_with_issues'].append(f"{name} (404s)")
                
            if page_state['network_errors']:
                self.test_results['summary']['total_errors'] += len(page_state['network_errors'])
                self.test_results['summary']['pages_with_issues'].append(f"{name} (network errors)")
                
            if page_state['javascript_errors']:
                self.test_results['summary']['total_errors'] += len(page_state['javascript_errors'])
                self.test_results['summary']['pages_with_issues'].append(f"{name} (JS errors)")
                
        except Exception as e:
            print(f"Error testing {name}: {e}")
            self.test_results['pages'][name] = {"error": str(e)}
            self.test_results['summary']['pages_with_issues'].append(f"{name} (load error)")
    
    def run_tests(self):
        """Run all tests"""
        print("Starting E173 Gateway browser tests...")
        
        # Add JS error handler
        self.driver.execute_script("""
            window.jsErrors = [];
            window.addEventListener('error', function(e) {
                window.jsErrors.push({
                    message: e.message,
                    source: e.filename,
                    line: e.lineno,
                    column: e.colno,
                    error: e.error ? e.error.toString() : ''
                });
            });
        """)
        
        # Test login
        if not self.test_login():
            print("Login failed, cannot continue tests")
            return
        
        # Test all pages
        pages_to_test = [
            ('/', 'dashboard'),
            ('/customers', 'customers'),
            ('/gateways', 'gateways'),
            ('/modems', 'modems'),
            ('/sims', 'sims'),
            ('/cdrs', 'cdrs'),
            ('/blacklist', 'blacklist')
        ]
        
        for path, name in pages_to_test:
            self.test_page(path, name)
            time.sleep(1)  # Small delay between pages
        
        # Clean up summary
        self.test_results['summary']['pages_with_issues'] = list(set(self.test_results['summary']['pages_with_issues']))
        
        # Generate report
        self.generate_report()
        
    def generate_report(self):
        """Generate test report"""
        print("\n" + "="*80)
        print("E173 GATEWAY BROWSER TEST REPORT")
        print("="*80)
        print(f"Test completed at: {self.test_results['timestamp']}")
        print(f"\nSUMMARY:")
        print(f"- Total console errors: {self.test_results['summary']['total_console_errors']}")
        print(f"- Total 404 errors: {self.test_results['summary']['total_404s']}")
        print(f"- Total network errors: {self.test_results['summary']['total_errors']}")
        print(f"- Pages with issues: {len(self.test_results['summary']['pages_with_issues'])}")
        
        print("\n\nDETAILED RESULTS BY PAGE:")
        print("-"*80)
        
        for page_name, page_data in self.test_results['pages'].items():
            print(f"\n{page_name.upper()}:")
            
            if isinstance(page_data, dict) and 'error' not in page_data:
                print(f"  URL: {page_data.get('url', 'N/A')}")
                print(f"  Title: {page_data.get('title', 'N/A')}")
                print(f"  Screenshot: {page_data.get('screenshot', 'N/A')}")
                
                if page_data.get('console_errors'):
                    print(f"  Console Errors ({len(page_data['console_errors'])}):")
                    for err in page_data['console_errors']:
                        print(f"    - {err['message']}")
                
                if page_data.get('missing_resources'):
                    print(f"  404 Errors ({len(page_data['missing_resources'])}):")
                    for res in page_data['missing_resources']:
                        print(f"    - {res['url']} ({res['type']})")
                
                if page_data.get('network_errors'):
                    print(f"  Network Errors ({len(page_data['network_errors'])}):")
                    for err in page_data['network_errors']:
                        print(f"    - {err['url']} (Status: {err['status']})")
                
                if page_data.get('javascript_errors'):
                    print(f"  JavaScript Errors ({len(page_data['javascript_errors'])}):")
                    for err in page_data['javascript_errors']:
                        print(f"    - {err['message']} at {err.get('source', 'unknown')}:{err.get('line', '?')}")
            else:
                print(f"  Error: {page_data}")
        
        # Save detailed report
        report_path = f"browser_test_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        with open(report_path, 'w') as f:
            json.dump(self.test_results, f, indent=2)
        print(f"\n\nDetailed report saved to: {report_path}")
        
        # Recommendations
        print("\n\nRECOMMENDATIONS:")
        print("-"*80)
        
        if self.test_results['summary']['total_404s'] > 0:
            print("1. Fix missing resources (404 errors):")
            print("   - Check that all CSS, JS, and image files exist")
            print("   - Verify correct paths in templates")
            print("   - Ensure static file serving is properly configured")
        
        if self.test_results['summary']['total_console_errors'] > 0:
            print("2. Address console errors:")
            print("   - Fix JavaScript syntax errors")
            print("   - Handle undefined variables")
            print("   - Add proper error handling")
        
        if self.test_results['summary']['total_errors'] > 0:
            print("3. Fix network errors:")
            print("   - Check API endpoints return proper status codes")
            print("   - Handle authentication properly")
            print("   - Verify CORS configuration if needed")
        
        if not any([self.test_results['summary']['total_404s'], 
                   self.test_results['summary']['total_console_errors'],
                   self.test_results['summary']['total_errors']]):
            print("No issues found! All pages are loading correctly.")
    
    def cleanup(self):
        """Clean up resources"""
        if hasattr(self, 'driver'):
            self.driver.quit()

if __name__ == "__main__":
    # Check if display is available
    if not os.environ.get('DISPLAY'):
        print("Setting up virtual display...")
        os.environ['DISPLAY'] = ':99'
        os.system('Xvfb :99 -screen 0 1920x1080x24 &')
        time.sleep(2)
    
    tester = E173BrowserTest()
    try:
        tester.run_tests()
    except Exception as e:
        print(f"Test failed with error: {e}")
    finally:
        tester.cleanup()