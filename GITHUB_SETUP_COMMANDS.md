# GitHub Projects Setup Commands

## 🚀 Run These Commands Tomorrow to Set Up Complete Project Management

### **1. GitHub Project Creation**
```bash
# Install GitHub CLI if not already installed
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update && sudo apt install gh

# Authenticate with GitHub
gh auth login

# Create the main project
gh repo create e173-intelligent-gateway --public --description "Intelligent VoIP Gateway with AI spam monetization"
cd e173_go_gateway
git remote add origin https://github.com/YOUR_USERNAME/e173-intelligent-gateway.git

# Create GitHub Project board
gh project create "E173 Gateway Development" --owner YOUR_USERNAME
```

### **2. Create All Issues at Once**
```bash
# Issue 1: Frontend Template Fix
gh issue create \
  --title "[AGENT-FRONTEND] Fix Template Collisions - High Priority" \
  --body "$(cat << 'EOF'
## 🎯 Goal
Fix template collision causing dashboard to show settings content

## 📋 Tasks
- [ ] Create unique template names (dashboard_standalone.tmpl, etc.)
- [ ] Fix all navigation sections (dashboard, modems, SIMs, customers, CDRs, blacklist)
- [ ] Test UI on both localhost and LAN IP (192.168.1.35:8080)
- [ ] Ensure all HTMX live updates work

## 📂 Files to Work With
- templates/*.tmpl
- cmd/server/main.go (route handlers)
- web/static/* (CSS assets)

## ✅ Success Criteria
- All navigation links work correctly
- Dashboard shows real dashboard content
- Each section displays appropriate data
- UI responsive and styled properly

## 🚨 Current Problem
Dashboard (/) shows "System Settings" instead of dashboard content due to Go template name collision with {{define "content"}}.

## 💡 Solution Approach
Create standalone templates for each section with unique names to avoid conflicts.
EOF
)" \
  --label "agent-frontend,high-priority,ui-fix" \
  --assignee YOUR_USERNAME

# Issue 2: WhatsApp API Integration  
gh issue create \
  --title "[AGENT-BACKEND] Complete WhatsApp API Integration - High Priority" \
  --body "$(cat << 'EOF'
## 🎯 Goal
Integrate user's private WhatsApp API for real person validation in spam filtering

## 🔑 API Details
- **URL:** https://bulkvalidation.wa-validator.xyz/v2/validate/wa_id?number={number}
- **Auth:** Bearer e42f7c9b-2a8e-4b86-a7e4-8f1de2c01f53
- **Response:** {"status":true,"valid":true,"wa_id":"34674944456","chat_link":"https://wa.me/34674944456"}

## 📋 Tasks
- [ ] Complete pkg/validation/private_whatsapp.go implementation
- [ ] Add 24-hour result caching to database
- [ ] Integrate with SIP server filtering pipeline
- [ ] Add batch validation for multiple numbers
- [ ] Create validation statistics and monitoring

## 📂 Files to Work With
- pkg/validation/private_whatsapp.go (already started)
- pkg/sip/basic_server.go (integrate filtering)
- Database schema for caching (whatsapp_validation_cache table exists)

## ✅ Success Criteria
- SIP server validates callers via WhatsApp API
- Results cached to avoid repeated API calls
- Spam calls (non-WhatsApp numbers) routed to AI agents
- Dashboard shows validation statistics

## 💰 Revenue Impact
This enables spam monetization by identifying non-human callers and routing them to AI agents instead of blocking.
EOF
)" \
  --label "agent-backend,high-priority,api-integration" \
  --assignee YOUR_USERNAME

# Issue 3: Voice Recognition
gh issue create \
  --title "[AGENT-AI] Implement Dual-Direction Voice Recognition - High Priority" \
  --body "$(cat << 'EOF'
## 🎯 Goal
Implement voice recognition for both spam detection (incoming) and SIM status monitoring (outgoing)

## 🎙️ Two-Direction System
**A) Source-side (Incoming):** Detect robocaller IVRs and route to AI agents for monetization
**B) SIM-side (Outgoing):** Detect operator messages (\"SIM blocked\", \"Low credit\", voicemail)

## 📋 Tasks
- [ ] Set up audio capture from SIP streams
- [ ] Implement speech-to-text (Whisper or Google Speech API)
- [ ] Create LLM classification system (GPT-4/Claude for transcript analysis)
- [ ] Build action engine (route to AI vs flag SIM vs normal routing)
- [ ] Create recording system for manual review and training

## 🤖 LLM Classification Prompt
```
Analyze this call transcript and classify:
Audio: "{transcript}"
Categories:
- SPAM_ROBOCALL: Route to AI agent for monetization
- SIM_BLOCKED: Flag SIM for replacement  
- VOICEMAIL: Handle appropriately
- NORMAL_CALL: Allow through
Response: {category} | confidence: {0.0-1.0} | action: {route_to_ai|flag_sim|normal_routing}
```

## 📂 Files to Create
- pkg/voice/recognition.go
- pkg/voice/classification.go
- pkg/ai/voice_agents.go
- Integration with pkg/sip/basic_server.go

## ✅ Success Criteria
- Real-time voice analysis on all calls
- Spam calls automatically routed to AI agents
- SIM status automatically detected and flagged
- Revenue generated from spam call monetization

## 💰 Revenue Impact
Primary monetization strategy - convert spam costs into revenue while protecting SIM infrastructure.
EOF
)" \
  --label "agent-ai,high-priority,voice-recognition,monetization" \
  --assignee YOUR_USERNAME

# Issue 4: Database Performance
gh issue create \
  --title "[AGENT-DATABASE] Optimize for High-Volume Operations - Medium Priority" \
  --body "$(cat << 'EOF'
## 🎯 Goal
Optimize database and caching for 200+ modems handling 1000+ concurrent calls

## 📋 Tasks
- [ ] Implement Redis caching layer for validation results
- [ ] Optimize database indexes for call_patterns, sip_calls tables
- [ ] Add connection pooling and query optimization
- [ ] Create analytics dashboard with real-time metrics
- [ ] Implement call pattern analysis for spam detection

## 🎯 Performance Targets
- **Response Time:** <100ms for API calls
- **Concurrent Calls:** 1000+ simultaneous
- **Database Load:** Handle 200+ modems continuous operation
- **Cache Hit Rate:** >80% for validation requests

## 📂 Files to Work With
- internal/repository/* (optimize queries)
- pkg/cache/* (create Redis caching)
- pkg/analytics/* (create analytics)
- migrations/* (add performance indexes)

## ✅ Success Criteria
- All API endpoints respond under 100ms
- System handles target concurrent load
- Real-time analytics dashboard functional
- Automated performance monitoring alerts
EOF
)" \
  --label "agent-database,medium-priority,performance,caching" \
  --assignee YOUR_USERNAME

# Issue 5: Production Deployment
gh issue create \
  --title "[AGENT-DEVOPS] Production Deployment & Monitoring - Medium Priority" \
  --body "$(cat << 'EOF'
## 🎯 Goal
Create production-ready deployment with monitoring and auto-scaling

## 🏗️ Architecture
- **Cloud VPS:** SIP server + Web interface + Database
- **Remote Gateways:** Multiple Asterisk servers with E173 modems
- **Communication:** Secure VPN/direct IP connections
- **Scaling:** Auto-scale based on call volume

## 📋 Tasks
- [ ] Create Docker containers for all services
- [ ] Set up cloud VPS (DigitalOcean/AWS/Linode)
- [ ] Implement monitoring (Prometheus + Grafana)
- [ ] Create CI/CD pipeline with GitHub Actions
- [ ] Set up log aggregation and alerting
- [ ] Configure SSL/TLS and security hardening

## 📂 Files to Create
- Dockerfile
- docker-compose.yml
- .github/workflows/deploy.yml
- deployment/production.yml
- monitoring/prometheus.yml

## ✅ Success Criteria
- One-command production deployment
- Real-time monitoring dashboard
- Automated scaling based on load
- 99.9% uptime with alerting
- Secure HTTPS access

## 🔒 Security Requirements
- SSL/TLS encryption
- VPN for gateway communication
- Database encryption at rest
- API rate limiting and DDoS protection
EOF
)" \
  --label "agent-devops,medium-priority,deployment,monitoring" \
  --assignee YOUR_USERNAME
```

### **3. Set Up Project Board Automation**
```bash
# Add automation rules
gh api graphql -f query='
mutation {
  updateProject(input: {
    projectId: "PROJECT_ID"
    title: "E173 Gateway Development"
  }) {
    project {
      id
    }
  }
}'

# Create project fields
gh project field-create PROJECT_ID --name "Status" --type "single_select" --option "Todo" --option "In Progress" --option "Review" --option "Done"
gh project field-create PROJECT_ID --name "Priority" --type "single_select" --option "High" --option "Medium" --option "Low"
gh project field-create PROJECT_ID --name "Agent Type" --type "single_select" --option "Frontend" --option "Backend" --option "AI" --option "Database" --option "DevOps"
```

### **4. Push All Code to GitHub**
```bash
# Add all project files
git add .
git commit -m "Initial E173 Gateway implementation with agent coordination system

- SIP server with intelligent filtering
- WhatsApp API integration framework  
- Database schema for multi-gateway support
- Agent coordination and GitHub Projects setup
- Voice recognition architecture planned"

git push -u origin main

# Create development branch
git checkout -b development
git push -u origin development
```

## 🎯 **Tomorrow's Agent Instructions**

When you wake up, follow these steps:

### **Step 1: Run GitHub Setup (5 minutes)**
```bash
cd /root/e173_go_gateway
# Run all commands from GITHUB_SETUP_COMMANDS.md
```

### **Step 2: Open 5 Claude Conversations (10 minutes)**
For each Claude instance, provide:

1. **MCP GitHub Server Access**: Connect to your repository
2. **Agent Role Assignment**: Copy instructions from AGENT_COORDINATION_MASTER_PLAN.md
3. **Specific Issue Assignment**: Assign each agent their GitHub issue
4. **Start Command**: "Begin working on your assigned issue and update progress every 2 hours"

### **Step 3: Monitor Progress**
- Check GitHub Project board hourly
- Review agent updates in issue comments
- Coordinate integration points
- Activate debug specialists as needed

This system will give you complete visibility and coordination of all development work!
