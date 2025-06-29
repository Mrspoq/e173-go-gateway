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

const token = config.token;
const owner = config.user;
const repo = config.repo;

// Helper function to make API requests
function makeRequest(method, path, data = null) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: path,
      method: method,
      headers: {
        'Authorization': `token ${token}`,
        'Accept': 'application/vnd.github.v3+json',
        'User-Agent': 'E173-Gateway-Setup'
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
          resolve(JSON.parse(body));
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

async function setupGitHubProject() {
  console.log('Setting up GitHub Project for E173 Gateway...');
  
  try {
    // Create milestones
    console.log('\nCreating milestones...');
    const milestones = [
      {
        title: 'Phase 1: Core Platform',
        description: 'CSS fixes, SIP server, database, basic filtering',
        due_on: '2025-07-15T00:00:00Z'
      },
      {
        title: 'Phase 2: Voice Recognition',
        description: 'Dual-direction voice detection, spam classification, SIM monitoring',
        due_on: '2025-08-01T00:00:00Z'
      },
      {
        title: 'Phase 3: AI Integration',
        description: 'AI voice agents, spam monetization, automated responses',
        due_on: '2025-08-15T00:00:00Z'
      },
      {
        title: 'Phase 4: Production Deploy',
        description: 'Cloud deployment, monitoring, multi-gateway management',
        due_on: '2025-09-01T00:00:00Z'
      }
    ];

    const createdMilestones = {};
    for (const milestone of milestones) {
      try {
        const result = await makeRequest('POST', `/repos/${owner}/${repo}/milestones`, milestone);
        createdMilestones[milestone.title] = result.number;
        console.log(`✓ Created milestone: ${milestone.title}`);
      } catch (error) {
        console.log(`✗ Failed to create milestone ${milestone.title}: ${error.message}`);
      }
    }

    // Create issues
    console.log('\nCreating issues...');
    const issues = [
      {
        title: '[EPIC] Cloud-Optimized SIP Platform',
        body: `## Epic Description
Build a scalable, cloud-ready SIP gateway platform for managing ~200 Huawei E173 USB modems.

## Success Criteria
- Asterisk integration with custom dialplan
- High-performance Go backend
- Real-time modem monitoring
- Call routing and management

## Sub-tasks
- [ ] Complete Asterisk configuration
- [ ] Implement modem management API
- [ ] Create SIP routing logic
- [ ] Add real-time monitoring`,
        labels: ['epic'],
        milestone: createdMilestones['Phase 1: Core Platform']
      },
      {
        title: '[EPIC] Multi-Gateway Voice Management',
        body: `## Epic Description
Implement distributed gateway architecture for load balancing and redundancy.

## Success Criteria
- Multiple Asterisk servers coordination
- Load balancing across gateways
- Failover mechanisms
- Centralized management

## Sub-tasks
- [ ] Gateway discovery service
- [ ] Load balancing algorithm
- [ ] Health monitoring
- [ ] Failover implementation`,
        labels: ['epic'],
        milestone: createdMilestones['Phase 2: Voice Recognition']
      },
      {
        title: '[EPIC] Voice Recognition & AI Integration',
        body: `## Epic Description
Integrate advanced voice recognition and AI-powered spam detection.

## Success Criteria
- Real-time voice transcription
- AI spam classification
- Automated response system
- WhatsApp verification

## Sub-tasks
- [ ] Voice recognition integration
- [ ] AI model deployment
- [ ] Response automation
- [ ] WhatsApp API integration`,
        labels: ['epic'],
        milestone: createdMilestones['Phase 3: AI Integration']
      },
      {
        title: 'Fix CDR and Blacklist pages calling modems API',
        body: `## Description
CDR and Blacklist pages are incorrectly calling /api/v1/modems repeatedly.

## Current Behavior
- Pages load but immediately start calling modems API
- This causes unnecessary load and incorrect data display

## Expected Behavior
- CDR page should show call records
- Blacklist page should show blocked numbers

## Tasks
- [ ] Fix template routing in CDR page
- [ ] Fix template routing in Blacklist page
- [ ] Test both pages function correctly`,
        labels: ['bug', 'frontend'],
        milestone: createdMilestones['Phase 1: Core Platform']
      },
      {
        title: 'Add customer SIP account management',
        body: `## Description
Implement SIP account management features for customers.

## Features
- Create/edit SIP accounts
- Manage credentials
- Set call routing rules
- Monitor usage

## Tasks
- [ ] Database schema for SIP accounts
- [ ] API endpoints
- [ ] UI components
- [ ] Integration with Asterisk`,
        labels: ['feature', 'backend'],
        milestone: createdMilestones['Phase 1: Core Platform']
      },
      {
        title: 'Setup production environment',
        body: `## Description
Configure production environment for deployment.

## Requirements
- Docker containers
- PostgreSQL database
- Redis cache
- Asterisk servers
- Monitoring

## Tasks
- [ ] Create Docker configurations
- [ ] Setup database migrations
- [ ] Configure Redis
- [ ] Deploy monitoring stack`,
        labels: ['devops', 'infrastructure'],
        milestone: createdMilestones['Phase 4: Production Deploy']
      }
    ];

    for (const issue of issues) {
      try {
        await makeRequest('POST', `/repos/${owner}/${repo}/issues`, issue);
        console.log(`✓ Created issue: ${issue.title}`);
      } catch (error) {
        console.log(`✗ Failed to create issue ${issue.title}: ${error.message}`);
      }
    }

    console.log('\nGitHub project setup complete!');
    console.log(`Visit: https://github.com/${owner}/${repo}/issues to view the issues`);
    console.log(`Visit: https://github.com/${owner}/${repo}/milestones to view the milestones`);

  } catch (error) {
    console.error('Error setting up GitHub project:', error);
  }
}

// Run the setup
setupGitHubProject();