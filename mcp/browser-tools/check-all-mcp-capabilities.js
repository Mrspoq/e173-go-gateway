// Check ALL Browser Tools MCP capabilities
const tools = [
  "getConsoleLogs",
  "getConsoleErrors", 
  "getNetworkErrors",
  "getNetworkLogs",
  "takeScreenshot",
  "getSelectedElement",
  "wipeLogs",
  "runAccessibilityAudit",
  "runPerformanceAudit",
  "runSEOAudit",
  "runNextJSAudit",
  "runDebuggerMode",
  "runAuditMode",
  "runBestPracticesAudit",
  // Potential hidden tools
  "navigate",
  "reload",
  "click",
  "type",
  "goBack",
  "goForward",
  "refresh",
  "navigateTo",
  "executeScript"
];

console.log("Checking for browser control capabilities in Browser Tools MCP...");
console.log("Available tools:", tools.slice(0, 14));
console.log("\nTesting potential hidden navigation tools:", tools.slice(14));
