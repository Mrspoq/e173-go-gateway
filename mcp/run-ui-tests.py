#!/usr/bin/env python3
"""
E173 Gateway UI Test Runner
Uses both Browser Use MCP and Browser Tools MCP for comprehensive testing
"""

import json
import subprocess
import time
import os
from datetime import datetime

class E173UITester:
    def __init__(self):
        self.base_url = "http://192.168.1.35:8080"
        self.results_dir = "/root/e173_go_gateway/mcp/test-results"
        self.timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        self.test_report = []
        
        # Create results directory
        os.makedirs(self.results_dir, exist_ok=True)
        
    def log(self, message):
        """Log message with timestamp"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] {message}")
        self.test_report.append({
            "timestamp": timestamp,
            "message": message
        })
    
    def run_browser_use_task(self, task):
        """Execute a task using Browser Use MCP"""
        self.log(f"Browser Use: {task}")
        # In real implementation, this would send commands to the MCP server
        # For now, we'll create a placeholder
        return {
            "status": "pending",
            "task": task,
            "note": "Requires MCP server connection"
        }
    
    def run_browser_tools_command(self, tool, args=None):
        """Execute a command using Browser Tools MCP"""
        self.log(f"Browser Tools: {tool}")
        # In real implementation, this would send commands to the MCP server
        return {
            "status": "pending",
            "tool": tool,
            "args": args,
            "note": "Requires MCP server connection"
        }
    
    def test_login_page(self):
        """Test the login page functionality"""
        self.log("=== Testing Login Page ===")
        
        # Navigate to login
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/login")
        
        # Clear console logs
        self.run_browser_tools_command("wipeLogs")
        
        # Take screenshot
        self.run_browser_tools_command("takeScreenshot", {"filename": f"login_page_{self.timestamp}.png"})
        
        # Test login
        self.run_browser_use_task("Fill the username field with 'admin', password field with 'admin', and click the login button")
        
        # Check for errors
        errors = self.run_browser_tools_command("getConsoleErrors")
        
        # Verify redirect
        self.run_browser_use_task("Verify that we are now on the dashboard page")
        
        return {
            "test": "login_page",
            "status": "tested",
            "errors": errors
        }
    
    def test_dashboard(self):
        """Test the dashboard functionality"""
        self.log("=== Testing Dashboard ===")
        
        # Navigate to dashboard
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/dashboard")
        
        # Clear logs
        self.run_browser_tools_command("wipeLogs")
        
        # Verify layout
        self.run_browser_use_task("Count the number of stat cards in the top row and verify they are displayed horizontally")
        
        # Check spacing
        self.run_browser_use_task("Check if there is proper spacing between the stat cards and the panels below (Live Modem Status and Recent Call Activity)")
        
        # Monitor console
        console_logs = self.run_browser_tools_command("getConsoleLogs")
        
        # Check network
        network_logs = self.run_browser_tools_command("getNetworkLogs")
        
        # Screenshot
        self.run_browser_tools_command("takeScreenshot", {"filename": f"dashboard_{self.timestamp}.png"})
        
        return {
            "test": "dashboard",
            "status": "tested",
            "console_logs": console_logs,
            "network_logs": network_logs
        }
    
    def test_customers_page(self):
        """Test the customers page functionality"""
        self.log("=== Testing Customers Page ===")
        
        # Navigate
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/customers")
        
        # Clear logs
        self.run_browser_tools_command("wipeLogs")
        
        # Check list
        self.run_browser_use_task("Verify that the customers list is displayed")
        
        # Test create button
        self.run_browser_use_task("Click on the 'Create Customer' button and verify the form appears")
        
        # Test edit button
        self.run_browser_use_task("Click on the edit button for the first customer in the list")
        
        # Check for auth redirect issue
        self.run_browser_use_task("Verify that we stay on the customers page and are not redirected to login")
        
        # Get errors
        errors = self.run_browser_tools_command("getConsoleErrors")
        
        return {
            "test": "customers_page",
            "status": "tested",
            "errors": errors
        }
    
    def test_gateways_page(self):
        """Test the gateways page"""
        self.log("=== Testing Gateways Page ===")
        
        # Navigate
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/gateways")
        
        # Clear logs
        self.run_browser_tools_command("wipeLogs")
        
        # Check if page loads
        self.run_browser_use_task("Check if the gateways page loads properly without showing a blank page")
        
        # Get errors
        errors = self.run_browser_tools_command("getConsoleErrors")
        
        # Screenshot
        self.run_browser_tools_command("takeScreenshot", {"filename": f"gateways_{self.timestamp}.png"})
        
        return {
            "test": "gateways_page",
            "status": "tested",
            "errors": errors
        }
    
    def test_modems_page(self):
        """Test the modems page"""
        self.log("=== Testing Modems Page ===")
        
        # Navigate
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/modems")
        
        # Clear logs
        self.run_browser_tools_command("wipeLogs")
        
        # Check layout
        self.run_browser_use_task("Check if modems are displayed properly without nested boxes")
        
        # Get errors
        errors = self.run_browser_tools_command("getConsoleErrors")
        
        return {
            "test": "modems_page",
            "status": "tested",
            "errors": errors
        }
    
    def test_cdr_page(self):
        """Test the CDR page"""
        self.log("=== Testing CDR Page ===")
        
        # Navigate
        self.run_browser_use_task("Navigate to http://192.168.1.35:8080/cdrs")
        
        # Clear logs
        self.run_browser_tools_command("wipeLogs")
        
        # Check empty state
        self.run_browser_use_task("Verify that an empty table structure is shown even when there are no CDR records")
        
        # Get errors
        errors = self.run_browser_tools_command("getConsoleErrors")
        
        return {
            "test": "cdr_page",
            "status": "tested",
            "errors": errors
        }
    
    def run_all_tests(self):
        """Run all UI tests"""
        self.log("Starting E173 Gateway Comprehensive UI Tests")
        self.log(f"Timestamp: {self.timestamp}")
        self.log(f"Base URL: {self.base_url}")
        
        test_results = []
        
        # Run tests
        test_results.append(self.test_login_page())
        test_results.append(self.test_dashboard())
        test_results.append(self.test_customers_page())
        test_results.append(self.test_gateways_page())
        test_results.append(self.test_modems_page())
        test_results.append(self.test_cdr_page())
        
        # Save results
        report_file = os.path.join(self.results_dir, f"test_report_{self.timestamp}.json")
        with open(report_file, 'w') as f:
            json.dump({
                "timestamp": self.timestamp,
                "base_url": self.base_url,
                "test_results": test_results,
                "full_log": self.test_report
            }, f, indent=2)
        
        self.log(f"Test report saved to: {report_file}")
        
        # Create summary
        self.create_summary(test_results)
        
    def create_summary(self, test_results):
        """Create a summary of test results"""
        self.log("\n=== TEST SUMMARY ===")
        
        total_tests = len(test_results)
        errors_found = []
        
        for result in test_results:
            test_name = result['test']
            if result.get('errors'):
                errors_found.append(test_name)
        
        self.log(f"Total tests run: {total_tests}")
        self.log(f"Tests with errors: {len(errors_found)}")
        
        if errors_found:
            self.log("\nTests with errors:")
            for test in errors_found:
                self.log(f"  - {test}")
        
        # Create action items
        self.log("\n=== ACTION ITEMS ===")
        self.log("1. Connect MCP servers to execute actual tests")
        self.log("2. Review and fix any errors found")
        self.log("3. Update GitHub project tracker with results")
        self.log("4. Implement missing features from PRD")

if __name__ == "__main__":
    tester = E173UITester()
    tester.run_all_tests()