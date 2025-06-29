# UI Feedback and Feature Requests - 2025-06-29

## UI Issues Observed

### 1. Dashboard
- **Issue**: 5 cards showing as 2 per row in 3 rows (looks like a table)
- **Expected**: Small cards on top in a single row
- **Missing**: Active/real-time calls display

### 2. Gateways Section
- **Bug**: Shows login prompt instead of gateway list
- **Expected**: List of gateways with empty slots to add more
- **Missing**: Gateway management interface

### 3. Modems Section
- **Issue**: Poor organization, nested boxes
- **Issue**: Inconsistent styling (first card has box inside box, second doesn't)
- **Needs**: Better organization and merge with SIM cards

### 4. Customers Section
- **Bug**: Edit button doesn't work
- **Bug**: Add customer button doesn't work
- **Missing**: Prepaid/Postpaid selection
- **Missing**: Add credits functionality
- **Missing**: Real-time balance updates

### 5. General UI Issues
- **Bug**: Dark theme doesn't work
- **Bug**: Shows "Login" in top right even when logged in
- **Issue**: Export report button non-functional

### 6. CDR Section
- **Missing**: Filtering options
- **Missing**: Call recordings

### 7. Blacklist Section
- **Missing**: Add individual numbers
- **Missing**: CSV import

## Feature Requests

### 1. Dashboard Enhancements
- **Active Calls Display**
  - Real-time connected calls
  - Similar to Call Me Soft billing system
  - Live updates
  
### 2. Gateway Management
- **Gateway List View**
  - Show existing gateways
  - Empty slots for adding new gateways
  - Display Asterisk private installation as first gateway
  
- **Add Gateway Form**
  - IP address
  - AMI port
  - SSL/TLS option
  - Connection to Asterisk AMI interface
  
- **Gateway Statistics**
  - Number of ports
  - Free dongles count
  - Total credits across all dongles
  - Gateway weight/priority for routing
  
- **Integration Features**
  - IMEI change scripts (Python/Shell)
  - API calls for configuration changes

### 3. Modem/SIM Unified Management
- **Merge Modems and SIMs sections**
  - SIMs are inside modems
  - Show modems without SIMs (for IMEI changes)
  - Filter by gateway
  - Network management capabilities

### 4. Customer Management Enhancements
- **Balance Types**
  - Prepaid customers
  - Postpaid customers
  - Setting during customer creation
  
- **Payment Management**
  - Add credits button
  - Payment history
  - Real-time balance deduction (per second billing)
  
- **Statistics**
  - ACD (Average Call Duration)
  - ASR (Answer Seizure Ratio)
  - Per customer analytics
  
- **SIP Configuration**
  - IP whitelist (PBX/SIP server IPs)
  - Codec selection (G729, G711 ulaw/alaw)
  - SIP credentials management

### 5. CDR Enhancements
- **Filtering Options**
  - By customer
  - By date range
  - By operator
  - By duration
  - Failed vs answered calls
  
- **Call Recordings**
  - MP3/WAV export
  - Evidence for successful calls
  - IVR recordings
  
- **Export Features**
  - Custom date ranges
  - Filter by call status
  - CSV/Excel formats

### 6. Blacklist Management
- **Input Methods**
  - Add individual numbers
  - Import CSV files
  - Bulk operations
  
- **Integration**
  - Phone number structure validation
  - WhatsApp API integration
  - Automatic updates from filter system

### 7. SIM Card Management System
- **Recharge System**
  - Database of prepaid recharge codes
  - Bulk selection (shift/cmd for multiple)
  - Sort by credits (minutes/money)
  - One-click bulk recharge
  
- **Operator-Specific Methods**
  - SMS-based recharge
  - USSD commands
  - Parse SMS confirmations
  - Update balance automatically
  
- **SIM Card Details**
  - IMSI (International Mobile Subscriber Identity)
  - MSISDN (Phone Number)
  - ICCID (SIM Card Serial Number)
  - Operator name
  - Balance in minutes
  - Balance in SMS
  - Balance in MB/GB (data)
  - Expiry dates for all balances
  
- **Promotions Handling**
  - Example: €1 = 1 hour talk time
  - Expiry in 2-3 days or 1 week
  - Track promotion types

## Technical Requirements

### Backend Features Needed
1. Real-time call monitoring via AMI
2. SIM card USSD command execution
3. SMS parsing for balance confirmations
4. Recharge code database management
5. Call recording storage and retrieval
6. Advanced CDR filtering queries
7. Gateway health monitoring

### Frontend Features Needed
1. Real-time WebSocket for active calls
2. Drag-and-drop for gateway priorities
3. Multi-select for bulk operations
4. Advanced filtering UI
5. Audio player for recordings
6. Export functionality
7. Real-time balance updates

## Priority Classification

### High Priority (Critical for Operations)
1. Fix gateway login bug
2. Fix customer edit/add buttons
3. Gateway management interface
4. SIM card recharge system
5. Real-time balance updates
6. CDR filtering

### Medium Priority (Important Features)
1. Merge modems/SIMs sections
2. Active calls display
3. Call recordings
4. Blacklist CSV import
5. Customer statistics (ACD/ASR)
6. Dark theme fix

### Low Priority (Nice to Have)
1. IMEI change scripts
2. Gateway weight/routing
3. Promotion management
4. Advanced analytics

## Implementation Plan

### Phase 1: Fix Critical Bugs (Week 1)
- Gateway login issue
- Customer management buttons
- User authentication display
- Dashboard card layout

### Phase 2: Core Features (Week 2-3)
- Gateway management UI
- SIM card recharge system
- Real-time balance updates
- CDR filtering

### Phase 3: Enhanced Features (Week 4-5)
- Active calls dashboard
- Merge modems/SIMs
- Call recordings
- Bulk operations

### Phase 4: Advanced Features (Week 6+)
- Analytics and statistics
- IMEI management
- Promotion handling
- Advanced routing