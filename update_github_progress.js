#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');

// Load GitHub configuration
const configPath = path.join(__dirname, '.github_config');
const config = fs.readFileSync(configPath, 'utf8')
  .split('\n')
  .reduce((acc, line) => {
    const [key, value] = line.split('=');
    if (key && value) {
      acc[key.replace('GITHUB_', '').toLowerCase()] = value;
    }
    return acc;
  }, {});

// Helper function to make API requests
function makeRequest(method, path, data = null) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: path,
      method: method,
      headers: {
        'Authorization': `token ${config.token}`,
        'Accept': 'application/vnd.github.v3+json',
        'User-Agent': 'E173-Gateway-Progress'
      }
    };

    if (data) {
      options.headers['Content-Type'] = 'application/json';
      options.headers['Content-Length'] = Buffer.byteLength(JSON.stringify(data));
    }

    const req = https.request(options, (res) => {
      let body = '';
      res.on('data', (chunk) => body += chunk);
      res.on('end', () => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(JSON.parse(body || '{}'));
        } else {
          reject(new Error(`API request failed: ${res.statusCode} ${body}`));
        }
      });
    });

    req.on('error', reject);
    if (data) {
      req.write(JSON.stringify(data));
    }
    req.end();
  });
}

async function updateProgress() {
  console.log('Updating GitHub project progress...');
  
  try {
    // Get all open issues
    const issues = await makeRequest('GET', `/repos/${config.user}/${config.repo}/issues?state=open`);
    console.log(`Found ${issues.length} open issues`);
    
    // Update specific issues based on completed work
    const updates = [
      {
        title: 'Fix import cycle preventing compilation',
        comment: '✅ COMPLETED: Fixed import cycle issue. Server now compiles successfully with all features. Created server_auth_fixed binary.',
        close: true
      },
      {
        title: 'Fix CDR and Blacklist pages calling modems API',
        comment: '🔧 IN PROGRESS: Found that API endpoints exist (/api/v1/cdr/recent/list and /api/v1/blacklist). Need to investigate why pages are calling wrong API.',
        close: false
      }
    ];
    
    // Add comments about completed work
    const completedWork = `
## Progress Update - ${new Date().toISOString()}

### ✅ Completed Tasks:
1. **Authentication System**
   - Fixed admin/admin login credentials
   - Created admin user in database with bcrypt hashing
   - Implemented user authentication display in navigation bar
   - Added logout functionality
   - Tested full authentication flow

2. **UI/UX Improvements**
   - Fixed dashboard card grid layout (4 to 5 columns)
   - Removed full page refresh on dashboard load
   - Fixed empty SIM cards container display
   - Added user authentication status in navigation

3. **Browser Automation MCP**
   - Set up Puppeteer-based browser automation server
   - Running on port 3001 for continuous UI testing
   - Added vision capabilities for screenshot analysis

4. **GitHub Project Management**
   - Created 4 project milestones
   - Created 3 epic issues
   - Created multiple sprint issues
   - Successfully integrated with GitHub API

### 🔧 In Progress:
- Investigating CDR/Blacklist pages API routing issue
- Customer stats showing JSON instead of HTML
- Gateway page routing issue

### 📝 Technical Notes:
- Server running: server_auth_fixed on port 8080
- Browser automation MCP: port 3001
- All authentication endpoints working correctly
- HTMX cleanup scripts added to prevent memory leaks
`;
    
    // Find and update the main tracking issue
    for (const issue of issues) {
      if (issue.title.includes('[EPIC] Cloud-Optimized SIP Platform')) {
        await makeRequest('POST', `/repos/${config.user}/${config.repo}/issues/${issue.number}/comments`, {
          body: completedWork
        });
        console.log(`✓ Updated epic issue #${issue.number}`);
      }
      
      // Update specific issues
      for (const update of updates) {
        if (issue.title === update.title) {
          await makeRequest('POST', `/repos/${config.user}/${config.repo}/issues/${issue.number}/comments`, {
            body: update.comment
          });
          console.log(`✓ Commented on issue #${issue.number}: ${update.title}`);
          
          if (update.close) {
            await makeRequest('PATCH', `/repos/${config.user}/${config.repo}/issues/${issue.number}`, {
              state: 'closed'
            });
            console.log(`✓ Closed issue #${issue.number}`);
          }
        }
      }
    }
    
    console.log('\nGitHub project updated successfully!');
    
  } catch (error) {
    console.error('Error updating GitHub project:', error.message);
  }
}

// Run the update
updateProgress();