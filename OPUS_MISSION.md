# Claude Opus - Sequential Mission Plan

## 🎯 PRIMARY MISSION
Execute the E173 Gateway completion plan in sequential order. Fix UI first, then work through each agent's responsibilities one by one.

## 🚨 IMMEDIATE PRIORITY: Fix Template Collision Issue
**Problem:** Dashboard shows "System Settings" instead of dashboard content
**Root Cause:** Multiple templates use `{{define "content"}}` - Go uses the last loaded template

**Quick Fix Strategy:**
1. Use the standalone templates already created:
   - `templates/dashboard_standalone.tmpl` 
   - `templates/modems_standalone.tmpl`
   - `templates/sims_standalone.tmpl`
   - `templates/customers_standalone.tmpl`

2. Update route handlers in `cmd/server/main.go` to use standalone templates

3. Test that all navigation works: Dashboard, Modems, SIMs, Customers, CDRs, Blacklist

**Success Criteria:** All UI sections accessible and showing correct content

## 📋 SEQUENTIAL AGENT MISSIONS (Execute After UI Fix)

### **Mission 1: Backend API Integration (2-3 hours)**
**Goal:** Complete WhatsApp API integration for spam detection

**Key Files:**
- `pkg/validation/private_whatsapp.go` (partially implemented)
- `pkg/sip/basic_server.go` (integrate with filtering)

**API Details:**
- URL: https://bulkvalidation.wa-validator.xyz/v2/validate/wa_id?number={number}
- Auth: Bearer e42f7c9b-2a8e-4b86-a7e4-8f1de2c01f53
- Response: `{"status":true,"valid":true,"wa_id":"34674944456","chat_link":"https://wa.me/34674944456"}`

**Tasks:**
- [ ] Complete private WhatsApp API implementation
- [ ] Add 24-hour caching to avoid repeated API calls
- [ ] Integrate with SIP server filtering pipeline
- [ ] Test spam detection and routing to AI agents

### **Mission 2: Voice Recognition System (4-6 hours)**
**Goal:** Implement dual-direction voice recognition

**Architecture:** See `VOICE_RECOGNITION_PLAN.md`

**Tasks:**
- [ ] Set up audio capture from SIP streams
- [ ] Implement speech-to-text (Whisper or Google Speech API)
- [ ] Create LLM classification for spam vs legitimate calls
- [ ] Build action engine (route to AI vs normal routing)

### **Mission 3: Database Performance (2-3 hours)**
**Goal:** Optimize for high-volume operations

**Tasks:**
- [ ] Add Redis caching layer
- [ ] Optimize database indexes
- [ ] Create analytics dashboard
- [ ] Test performance with simulated load

### **Mission 4: Production Deployment (3-4 hours)**
**Goal:** Deploy to cloud with monitoring

**Tasks:**
- [ ] Create Docker containers
- [ ] Set up cloud VPS deployment
- [ ] Implement monitoring and alerting
- [ ] Create CI/CD pipeline

### **Mission 5: Advanced Features (2-3 hours)**
**Goal:** Complete advanced spam detection and monetization

**Tasks:**
- [ ] Implement advanced spam pattern detection
- [ ] Create AI voice agent integration
- [ ] Add revenue tracking and reporting
- [ ] Test complete spam monetization flow

## 🔄 EXECUTION PROTOCOL

**For Each Mission:**
1. **Start:** Update this file with "Mission X: IN PROGRESS"
2. **Work:** Focus on that mission's tasks only
3. **Test:** Verify everything works before moving on
4. **Complete:** Update this file with "Mission X: COMPLETED"
5. **Next:** Move to next mission

**Status Tracking:**
- ✅ UI Fix: COMPLETED (2025-06-28 03:15)
- ✅ Mission 1 (Backend): COMPLETED (2025-06-28 04:10)
- ❌ Mission 2 (Voice): PENDING
- ❌ Mission 3 (Database): PENDING
- ❌ Mission 4 (Deploy): PENDING
- ❌ Mission 5 (Advanced): PENDING

## 🆘 IF STUCK
- Update this file with "BLOCKED: [reason]"
- Move to next mission and come back later
- Focus on completing 80% rather than perfecting 100%

## 🎯 SUCCESS CRITERIA
- UI fully functional (all sections work)
- WhatsApp API detecting and routing spam calls
- Voice recognition identifying robocallers
- System deployed and monitoring in production
- Revenue being generated from spam call monetization

**Expected Total Time:** 12-18 hours of focused work
**Expected Result:** Production-ready intelligent VoIP gateway with AI spam monetization

---
**Start with UI fix - everything else depends on having a working interface!**
