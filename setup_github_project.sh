#!/bin/bash

# GitHub Project Setup Script for E173 Gateway
# User: mrspoq
# Repository: e173-go-gateway

echo "Setting up GitHub Project for E173 Gateway..."

# 1. Create GitHub Project
echo "Creating GitHub Project..."
gh project create "E173-Intelligent-Gateway" --owner mrspoq --public

# 2. Create Milestones
echo "Creating milestones..."
gh api repos/mrspoq/e173-go-gateway/milestones \
  --method POST \
  --field title="Phase 1: Core Platform" \
  --field description="CSS fixes, SIP server, database, basic filtering" \
  --field due_on="2025-07-15T00:00:00Z"

gh api repos/mrspoq/e173-go-gateway/milestones \
  --method POST \
  --field title="Phase 2: Voice Recognition" \
  --field description="Dual-direction voice detection, spam classification, SIM monitoring" \
  --field due_on="2025-08-01T00:00:00Z"

gh api repos/mrspoq/e173-go-gateway/milestones \
  --method POST \
  --field title="Phase 3: AI Integration" \
  --field description="AI voice agents, spam monetization, automated responses" \
  --field due_on="2025-08-15T00:00:00Z"

gh api repos/mrspoq/e173-go-gateway/milestones \
  --method POST \
  --field title="Phase 4: Production Deploy" \
  --field description="Cloud deployment, monitoring, multi-gateway management" \
  --field due_on="2025-09-01T00:00:00Z"

# 3. Create Epic Issues
echo "Creating epic issues..."

# Epic 1: Cloud-Optimized SIP Platform
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "[EPIC] Cloud-Optimized SIP Platform" \
  --body "## Epic Description
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
- [ ] Add real-time monitoring" \
  --label "epic" \
  --milestone "Phase 1: Core Platform"

# Epic 2: Multi-Gateway Voice Management
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "[EPIC] Multi-Gateway Voice Management" \
  --body "## Epic Description
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
- [ ] Failover implementation" \
  --label "epic" \
  --milestone "Phase 2: Voice Recognition"

# Epic 3: Voice Recognition & AI Integration
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "[EPIC] Voice Recognition & AI Integration" \
  --body "## Epic Description
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
- [ ] WhatsApp API integration" \
  --label "epic" \
  --milestone "Phase 3: AI Integration"

# 4. Create Feature Issues for Current Sprint
echo "Creating current sprint issues..."

# Fix import cycle issue
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "Fix import cycle preventing compilation" \
  --body "## Description
There's an import cycle preventing the server from compiling with new features.

## Error
\`\`\`
package command-line-arguments
        imports github.com/e173-gateway/e173_go_gateway/internal/handlers
        imports github.com/e173-gateway/e173_go_gateway/internal/services
        imports github.com/e173-gateway/e173_go_gateway/internal/handlers: import cycle not allowed
\`\`\`

## Tasks
- [ ] Identify circular dependencies
- [ ] Refactor package structure
- [ ] Test compilation
- [ ] Deploy new binary" \
  --label "bug,high-priority" \
  --milestone "Phase 1: Core Platform"

# Fix CDR/Blacklist pages
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "Fix CDR and Blacklist pages calling modems API" \
  --body "## Description
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
- [ ] Test both pages function correctly" \
  --label "bug,frontend" \
  --milestone "Phase 1: Core Platform"

# Customer SIP Management
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "Add customer SIP account management" \
  --body "## Description
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
- [ ] Integration with Asterisk" \
  --label "feature,backend" \
  --milestone "Phase 1: Core Platform"

# Setup Production Environment
gh issue create \
  --repo mrspoq/e173-go-gateway \
  --title "Setup production environment" \
  --body "## Description
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
- [ ] Deploy monitoring stack" \
  --label "devops,infrastructure" \
  --milestone "Phase 4: Production Deploy"

echo "GitHub Project setup complete!"
echo "Visit: https://github.com/mrspoq/e173-go-gateway/projects to view the project board"