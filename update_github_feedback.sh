#!/bin/bash

# Source GitHub configuration
source .github_config

# GitHub API base URL
API_BASE="https://api.github.com"
REPO_PATH="repos/${GITHUB_USER}/${GITHUB_REPO}"

# Create headers for authentication
AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
ACCEPT_HEADER="Accept: application/vnd.github.v3+json"

echo "Creating GitHub issues for UI feedback and features..."

# Issue 1: Critical UI Bugs
BUGS_BODY=$(cat <<EOF
# Critical UI Bugs to Fix

## High Priority Bugs

### 1. Gateway Section Login Bug 🚨
- **Issue**: Gateway page shows authentication prompt instead of gateway list
- **Expected**: Should show list of gateways with management options
- **Impact**: Cannot manage gateways at all

### 2. Customer Management Buttons 🚨
- **Issue**: Edit and Add Customer buttons are not working
- **Expected**: Should open edit form or create new customer form
- **Impact**: Cannot modify customer data

### 3. Dashboard Layout Issue 🚨
- **Issue**: 5 cards showing as 2 per row in 3 rows (table-like)
- **Expected**: 5 small cards in a single row at the top
- **Current**: Looks like a table instead of dashboard cards

### 4. Authentication Display Bug 🚨
- **Issue**: Shows "Login" in navigation even when user is authenticated
- **Expected**: Should show "Welcome, [Username]" with logout option
- **Note**: Backend is working, frontend not updating properly

### 5. Dark Theme Not Working
- **Issue**: Dark theme toggle doesn't apply dark styles
- **Expected**: Should toggle between light and dark themes
- **Impact**: User preference not respected

### 6. Export Report Button
- **Issue**: Export report button in customers section doesn't work
- **Expected**: Should show export options dialog

## Reproduction Steps
1. Login with admin/admin
2. Navigate to each section
3. Try the mentioned actions

## Environment
- Server: Latest build running on port 8080
- Browser: All modern browsers affected
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"🐛 Critical UI Bugs - Gateway, Customer, Dashboard\",
    \"body\": $(echo "$BUGS_BODY" | jq -Rs .),
    \"labels\": [\"bug\", \"high-priority\", \"ui\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created critical bugs issue"

# Issue 2: Gateway Management Feature
GATEWAY_BODY=$(cat <<EOF
# Gateway Management System

## Overview
Implement comprehensive gateway management for Asterisk AMI connections.

## Features Required

### 1. Gateway List View
- Show existing gateways in a grid/list
- Display empty slots for adding new gateways
- Show default Asterisk private installation
- Status indicators (online/offline)

### 2. Add Gateway Form
- **Fields Required**:
  - Gateway Name
  - IP Address
  - AMI Port (default: 5038)
  - AMI Username
  - AMI Password
  - SSL/TLS Enable checkbox
  - Description/Location
  - Weight/Priority (for routing)

### 3. Gateway Statistics Dashboard
- Number of total ports
- Free dongles count
- Occupied dongles
- Total credits across all dongles
- Active calls per gateway
- Gateway health status

### 4. Advanced Features
- IMEI change script integration
- Python/Shell script execution
- API endpoint for configuration changes
- Real-time status via WebSocket

### 5. Gateway Actions
- Edit gateway settings
- Delete gateway (with confirmation)
- Test connection
- View detailed logs
- Restart/reload gateway

## Technical Implementation
- Backend: AMI connection manager
- Frontend: Real-time status updates
- Database: Gateway configurations
- Security: Encrypted password storage

## UI Mockup
\`\`\`
[Gateway 1]        [Gateway 2]        [+ Add Gateway]
Asterisk Local     Remote GW-01        Empty Slot
✅ Online          ❌ Offline
Ports: 50/200      Ports: 0/100
Credits: €1,250    Credits: €0
[Edit] [Delete]    [Edit] [Delete]    [Configure]
\`\`\`
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"🚀 Feature: Gateway Management System\",
    \"body\": $(echo "$GATEWAY_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"high-priority\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created gateway management issue"

# Issue 3: SIM Card Recharge System
RECHARGE_BODY=$(cat <<EOF
# SIM Card Recharge Management System

## Overview
Comprehensive prepaid SIM card recharge system with bulk operations and SMS parsing.

## Core Features

### 1. Recharge Code Database
- Store prepaid recharge codes
- Track used/unused codes
- Operator-specific codes
- Expiry date tracking

### 2. SIM Card Information Display
- **Essential Fields**:
  - IMSI (International Mobile Subscriber Identity)
  - MSISDN (Phone Number - via USSD if needed)
  - ICCID (SIM Card Serial Number)
  - Operator Name
  - Balance in Minutes
  - Balance in SMS
  - Balance in MB/GB (Data)
  - Expiry Dates (for each balance type)
  - Last Recharge Date
  - Status (Active/Inactive/Suspended)

### 3. Bulk Recharge Operations
- Multi-select with Shift/Cmd
- Sort by credits (ascending/descending)
- Filter by operator
- Filter by balance threshold
- One-click bulk recharge
- Progress indicator for bulk operations

### 4. Operator-Specific Methods
- **SMS-based recharge**:
  - Send recharge code via SMS
  - Parse confirmation SMS
  - Auto-update balance
- **USSD Commands**:
  - Execute USSD for balance check
  - Get phone number if not in DB
- **API Integration** (if available)

### 5. Promotion Handling
- Track promotion types (e.g., €1 = 1 hour)
- Expiry periods (2-3 days, 1 week)
- Auto-calculate remaining time
- Alert before expiry

### 6. SMS Parsing Engine
- Operator-specific parsers
- Extract balance information
- Extract expiry dates
- Handle multiple languages
- Error handling for failed parses

## Technical Requirements
- AT command integration for SMS/USSD
- Background job for SMS monitoring
- Real-time balance updates
- WebSocket for live updates
- Audit trail for all recharges

## UI Concept
\`\`\`
[Select] | MSISDN      | Operator | Minutes | SMS  | Data   | Expires    | Action
☐       | +1234567890 | Orange   | 45 min  | 100  | 500MB  | 2 days     | [Recharge]
☑       | +0987654321 | Vodafone | 5 min   | 10   | 0MB    | 12 hours   | [Recharge]
☑       | +1122334455 | O2       | 0 min   | 0    | 0MB    | Expired    | [Recharge]

[Bulk Recharge Selected] [Check All Balances] [Import Codes]
\`\`\`
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"💳 Feature: SIM Card Recharge Management System\",
    \"body\": $(echo "$RECHARGE_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"high-priority\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created SIM recharge system issue"

# Issue 4: Real-time Dashboard
DASHBOARD_BODY=$(cat <<EOF
# Real-time Dashboard Enhancement

## Overview
Transform dashboard into real-time monitoring center with active calls and live statistics.

## Features

### 1. Active Calls Widget
- **Display Requirements**:
  - Caller Number
  - Called Number
  - Duration (live counter)
  - Gateway/Port being used
  - Customer name
  - Rate per minute
  - Current cost (updating)
- **Similar to**: Call Me Soft billing system
- **Update frequency**: Every second

### 2. Dashboard Layout Redesign
- **Top Row**: 5 small stat cards
  - Total Active Calls
  - Total Gateways Online
  - Total Credits Available
  - Today's Revenue
  - System Health
- **Main Area**: Active calls table
- **Side Panel**: Quick actions

### 3. Real-time Statistics
- Calls per minute graph
- Revenue per hour
- Gateway utilization
- Top customers by usage
- Failed call alerts

### 4. WebSocket Integration
- Live call events
- Real-time balance updates
- Gateway status changes
- System alerts

## Technical Implementation
- Backend: AMI event monitoring
- WebSocket server for real-time updates
- Frontend: React/Vue components for live data
- Caching: Redis for performance

## UI Mockup
\`\`\`
[📞 Active: 47] [🌐 Gateways: 3/5] [💰 Credits: €5,432] [📈 Today: €234] [✅ Health: 98%]

Active Calls:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Caller       | Called      | Duration | Gateway | Customer    | Cost
+123456789   | +987654321  | 00:03:45 | GW-01   | TechCorp   | €0.15
+111222333   | +444555666  | 00:01:23 | GW-02   | CallPlus   | €0.05
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
\`\`\`
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"📊 Feature: Real-time Dashboard with Active Calls\",
    \"body\": $(echo "$DASHBOARD_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"dashboard\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created real-time dashboard issue"

# Issue 5: Customer Management Enhancements
CUSTOMER_BODY=$(cat <<EOF
# Customer Management Enhancements

## Overview
Enhance customer management with prepaid/postpaid support, real-time billing, and advanced features.

## Features Required

### 1. Customer Types
- **Prepaid Customers**:
  - Balance deduction in real-time
  - Low balance alerts
  - Auto-suspend on zero balance
- **Postpaid Customers**:
  - Credit limit management
  - Monthly invoicing
  - Overage alerts

### 2. Payment Management
- Add Credits button (working)
- Payment methods tracking
- Payment history
- Invoice generation
- Receipt printing

### 3. Real-time Billing
- Per-second billing calculation
- Live balance updates during calls
- WebSocket for instant updates
- Balance history graph

### 4. Customer Statistics
- **ACD** (Average Call Duration)
- **ASR** (Answer Seizure Ratio)
- **Daily/Weekly/Monthly** usage
- Top destinations
- Call quality metrics

### 5. SIP Configuration
- IP Whitelist management
- Multiple IPs per customer
- Codec preferences:
  - G.729
  - G.711 μ-law
  - G.711 A-law
- Concurrent call limits
- Password management

### 6. Export Functionality
- **Export Options**:
  - Payment history
  - CDRs (with filters)
  - Statistics report
  - Invoice PDF
- **Filters**:
  - Date range
  - Call status
  - Destination

## Technical Requirements
- Real-time balance calculation engine
- WebSocket for live updates
- PDF generation for invoices
- Advanced SQL queries for statistics

## UI Enhancements
- Tab-based customer view
- Real-time balance indicator
- Quick actions menu
- Bulk operations support
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"👥 Feature: Customer Management Enhancements\",
    \"body\": $(echo "$CUSTOMER_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"customers\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created customer enhancements issue"

# Issue 6: CDR and Recording System
CDR_BODY=$(cat <<EOF
# CDR Enhancement with Call Recordings

## Overview
Comprehensive CDR system with filtering, recordings, and export capabilities.

## Features

### 1. Advanced Filtering
- **Filter Options**:
  - Customer (dropdown)
  - Date range (date pickers)
  - Operator/Gateway
  - Call duration (min/max)
  - Call status (answered/failed/busy)
  - Destination country
  - Price range

### 2. Call Recordings
- **For each CDR**:
  - Play button for in-browser playback
  - Download as MP3/WAV
  - Recording duration
  - File size
- **Storage**:
  - Organized by date/customer
  - Retention policy settings
  - Compression options

### 3. Export Features
- **Formats**: CSV, Excel, PDF
- **Custom columns** selection
- **Include recordings** option (zip file)
- **Scheduled exports** via email

### 4. CDR Details
- Caller ID and name
- Called number
- Start time
- Answer time
- End time
- Duration
- Bill duration
- Rate applied
- Total cost
- Gateway used
- Termination cause
- Recording available (yes/no)

## Technical Implementation
- Recording via Asterisk MixMonitor
- Storage in organized directory structure
- Database references to recording files
- Audio streaming server
- Batch export processing

## UI Mockup
\`\`\`
Filters: [Customer ▼] [From: ___] [To: ___] [Status ▼] [Apply]

CDR List:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Time      | From        | To         | Duration | Cost  | Rec | Actions
09:45:12  | +123456789  | +987654321 | 00:03:45 | €0.15 | 🔊  | [▶️] [⬇️]
09:43:30  | +111222333  | +444555666 | 00:01:23 | €0.05 | 🔊  | [▶️] [⬇️]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[Export Selected] [Export All] [Email Report]
\`\`\`
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"📞 Feature: CDR Enhancement with Call Recordings\",
    \"body\": $(echo "$CDR_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"cdr\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created CDR enhancement issue"

# Issue 7: Unified Modem/SIM Management
MODEM_SIM_BODY=$(cat <<EOF
# Unified Modem/SIM Management

## Overview
Merge modems and SIM cards into unified management interface.

## Design Concept

### 1. Unified View
- Show modems with their SIM cards
- Empty modems (ready for IMEI change)
- Filter by gateway
- Search by IMSI/MSISDN/IMEI

### 2. Modem Information
- IMEI
- Model
- Signal strength
- Temperature
- Network status
- Current SIM (if any)

### 3. SIM Operations
- Insert/Remove SIM (virtual assignment)
- Swap SIMs between modems
- Bulk SIM operations
- Track SIM history

### 4. IMEI Change Management
- Modems marked for IMEI change
- Script execution interface
- Change history tracking
- Validation before/after change

### 5. Network Management
- Network selection
- Signal quality monitoring
- Roaming settings
- APN configuration

## Technical Implementation
- AT commands for modem info
- Database relationships: modem ↔ SIM
- Background monitoring service
- WebSocket for real-time updates

## UI Concept
\`\`\`
Gateway: [All Gateways ▼] Network: [All ▼] Status: [All ▼]

┌─────────────────────────────────────┐
│ Modem: 353123456789012 📶 85%      │
│ SIM: +1234567890 (Orange)          │
│ Balance: 45min | 100SMS | 500MB    │
│ [Recharge] [Remove SIM] [Details]  │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ Modem: 353987654321098 📶 92%      │
│ No SIM - Ready for IMEI change     │
│ [Insert SIM] [Change IMEI]         │
└─────────────────────────────────────┘
\`\`\`
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"📱 Feature: Unified Modem/SIM Management\",
    \"body\": $(echo "$MODEM_SIM_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"modems\"],
    \"milestone\": 2
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created unified modem/SIM issue"

echo "All GitHub issues created successfully!"

# Create a summary issue for the roadmap
ROADMAP_BODY=$(cat <<EOF
# E173 Gateway Development Roadmap

## Overview
Based on user feedback from 2025-06-29, here's our development roadmap.

## Phase 1: Critical Fixes (Week 1)
- [ ] Fix gateway authentication bug (#10)
- [ ] Fix customer management buttons (#10)
- [ ] Fix dashboard layout (#10)
- [ ] Fix authentication display (#10)

## Phase 2: Core Features (Week 2-3)
- [ ] Gateway Management System (#11)
- [ ] SIM Card Recharge System (#12)
- [ ] Real-time Dashboard (#13)
- [ ] Customer Enhancements (#14)

## Phase 3: Advanced Features (Week 4-5)
- [ ] CDR with Recordings (#15)
- [ ] Unified Modem/SIM (#16)
- [ ] Blacklist Management
- [ ] Advanced Analytics

## Technical Debt
- [ ] WebSocket implementation
- [ ] Performance optimization
- [ ] Security hardening
- [ ] API documentation

## Success Metrics
- All critical bugs fixed
- Real-time call monitoring working
- SIM recharge system operational
- Customer can manage gateways
- CDR exports functional

## Dependencies
- Asterisk AMI integration
- AT command library
- Audio processing for recordings
- SMS parsing engine
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"🗺️ Development Roadmap - UI Feedback Implementation\",
    \"body\": $(echo "$ROADMAP_BODY" | jq -Rs .),
    \"labels\": [\"roadmap\", \"planning\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created roadmap issue!"
echo "GitHub project fully updated with all feedback!"