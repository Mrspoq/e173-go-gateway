# Multi-Agent Coordination Instructions

## 🎯 Current Status: Foundation Complete (3/5 Tasks Done)

✅ **Completed (1 hour):**
- CSS loading fixed → UI now works
- Database schema extended → Multi-gateway ready  
- Basic SIP server implemented → Call routing foundation ready

## 🚀 **Next Phase: Open 4 Claude Tabs Simultaneously**

### **Tab 1: AI Integration Specialist**
**Conversation Starter:**
```
I'm implementing AI voice recognition and spam monetization for the E173 gateway project. 

Current status:
- Basic SIP server running on port 5060
- Database schema has ai_voice_agents and ai_voice_sessions tables
- Need to implement:
  1. Voice recognition for SIM block detection
  2. AI voice agents for spam call monetization  
  3. WhatsApp API integration for real person validation
  4. Google Phone Lib for number validation

The goal is to route spam calls to AI agents (monetize them) and protect SIM cards from being blocked.

Here's the current SIP server structure: [copy pkg/sip/basic_server.go content]

Let's implement the AI integration services. Start with the voice recognition service for detecting when SIM cards are blocked by operators.
```

### **Tab 2: Advanced SIP Filtering Specialist** 
**Conversation Starter:**
```
I'm implementing advanced SIP call filtering for the E173 gateway project.

Current status:
- Basic SIP server running with placeholder filtering
- Database has operator_routing_rules, blacklist, call_patterns tables
- Need to implement:
  1. WhatsApp Business API integration
  2. Google Phone Number Library validation
  3. Intelligent operator prefix routing
  4. Sticky routing (same number → same gateway)
  5. Real-time blacklist and spam scoring

The filtering engine should process each call through multiple validation layers before routing.

Here's the current filter structure: [copy FilterEngine from pkg/sip/basic_server.go]

Let's build the production-ready filtering pipeline with real API integrations.
```

### **Tab 3: Frontend Enhancement Specialist**
**Conversation Starter:**
```
I'm enhancing the frontend for the E173 gateway management platform.

Current status:
- CSS loading fixed, HTMX working
- Basic dashboard functional at http://localhost:8080
- Authentication working (admin/admin)
- Need to implement:
  1. Role-based UI (superuser, manager, gateway operator)
  2. Multi-gateway management interface
  3. Real-time SIP call monitoring
  4. AI voice agent dashboard
  5. Spam monetization revenue tracking

The database schema supports multi-gateway with tables: gateways, sip_calls, ai_voice_sessions, revenue_tracking.

Current template structure: templates/base.tmpl, dashboard.tmpl, navigation in templates/partials/nav.tmpl

Let's create the role-based multi-gateway management interface.
```

### **Tab 4: Production Deployment Specialist**
**Conversation Starter:**
```
I'm setting up production deployment for the E173 gateway platform.

Current status:
- Go backend with SIP server running locally
- PostgreSQL database configured
- Basic services working
- Need to implement:
  1. Cloud VPS deployment configuration
  2. Docker containerization
  3. Reverse proxy setup (Nginx)
  4. SSL/TLS certificates
  5. Process monitoring and auto-restart
  6. Log aggregation and monitoring

The architecture should support:
- Cloud SIP server receiving calls
- Multiple remote Asterisk gateways
- AI voice services
- High availability and failover

Let's create the production deployment infrastructure.
```

## 🎯 **Coordination Protocol**

### **Each Agent Should:**
1. **Work independently** on their specialization
2. **Report progress** in their conversation
3. **Share integration points** when connecting with other components
4. **Test thoroughly** before marking complete

### **Integration Points:**
- **SIP Server** ← AI Integration (voice detection, spam routing)
- **SIP Server** ← Advanced Filtering (API validations, routing rules)  
- **Database** ← All agents (shared data models)
- **Frontend** ← All agents (management interfaces)
- **Production** ← All agents (deployment configuration)

### **Expected Timeline:**
- **Next 2-3 hours**: All 4 specializations complete
- **Integration testing**: 30 minutes
- **Production deployment**: 1 hour
- **Total**: 4-5 hours for complete platform

## 🏆 **Success Metrics**

### **Technical:**
- SIP calls processed through intelligent filtering
- Spam calls routed to AI agents for monetization
- Real-time dashboard showing call activity
- Multi-gateway management working
- Production deployment live

### **Business:**
- Spam call revenue generation active
- SIM card protection from blocking
- Multi-gateway operator management
- Role-based access control

## 📋 **Next Immediate Actions**

1. **Open 4 Claude tabs** with the conversation starters above
2. **Start all agents simultaneously** 
3. **Check back in 1 hour** for integration
4. **Deploy to production VPS** once complete

Each agent has full context and ready-to-execute code. The foundation is solid - now it's time to build the advanced features in parallel!
