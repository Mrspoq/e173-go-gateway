# Next Steps - E173 Gateway Development

## High Priority Features to Implement

### 1. Gateway Management Interface
**User Requirements**:
- Display existing gateways in grid/card layout
- Show empty slots for adding new gateways
- Default Asterisk private installation as first gateway
- Online/offline status indicators
- Gateway statistics (ports, dongles, credits)
- Add/Edit/Delete functionality
- AMI connection configuration

**Technical Tasks**:
- Create gateway list UI with card layout
- Implement create/edit forms with AMI settings
- Add real-time status monitoring via AMI
- Implement gateway health checks
- Create gateway statistics aggregation

### 2. SIM Card Recharge System
**User Requirements**:
- Database of prepaid recharge codes
- Bulk selection with Shift/Cmd
- Sort by credits (ascending/descending)
- One-click bulk recharge
- SMS-based recharge with confirmation parsing
- Balance tracking (minutes, SMS, data)
- Expiry date management
- Promotion handling (e.g., €1 = 1 hour)

**Technical Tasks**:
- Create recharge codes table and model
- Implement AT command interface for SMS
- Build SMS parser for operator responses
- Create bulk operation UI with multi-select
- Implement balance update system
- Add expiry date tracking and alerts

### 3. Real-time Dashboard with Active Calls
**User Requirements**:
- Display active calls with live duration counter
- Show caller/called numbers, customer, rate, cost
- Similar to Call Me Soft billing system
- Update every second
- 5 stat cards in top row

**Technical Tasks**:
- Implement WebSocket server
- Create AMI event listener for call events
- Build real-time call tracking system
- Design active calls table UI
- Implement live cost calculation

### 4. CDR Enhancement with Filtering
**User Requirements**:
- Filter by customer, date, operator, duration, status
- Call recordings with MP3/WAV export
- Export to CSV/Excel
- Custom date ranges

**Technical Tasks**:
- Create advanced filter UI
- Implement recording storage system
- Build export functionality
- Add recording playback UI

### 5. Customer Management Enhancements
**User Requirements**:
- Prepaid/Postpaid customer types
- Real-time balance deduction
- IP whitelist management
- Codec selection (G.729, G.711)
- Statistics (ACD, ASR)
- Payment history

**Technical Tasks**:
- Extend customer model for types
- Implement real-time billing engine
- Create IP management UI
- Add codec configuration
- Build statistics calculation

## Medium Priority Features

### 6. Unified Modem/SIM Management
- Merge modems and SIMs sections
- Show modems with their SIM cards
- Filter by gateway
- IMEI change management

### 7. Blacklist CSV Import
- CSV file upload
- Bulk number import
- Validation and error handling

### 8. Settings Enhancement
- Tabbed interface (Filter, SIP, AMI, General)
- System configuration management
- User preferences

## Technical Infrastructure

### WebSocket Implementation
- Real-time updates for calls, balances, status
- Event-driven architecture
- Client reconnection handling

### AT Command Integration
- SMS sending/receiving
- USSD command execution
- Balance checking

### AMI Integration Enhancement
- Event monitoring
- Call control
- Gateway status tracking

## Database Migrations Needed
- Recharge codes table
- Call recordings references
- Customer types and limits
- IP whitelist table

## API Endpoints to Create
- `/api/v1/gateways/*/status` - Real-time status
- `/api/v1/recharge/bulk` - Bulk recharge
- `/api/v1/calls/active` - Active calls list
- `/api/v1/cdr/export` - CDR export
- `/api/v1/customers/*/ips` - IP management

## UI Components to Build
- Gateway card component
- Multi-select table with Shift/Cmd
- Real-time call table
- Advanced filter component
- Audio player for recordings

## Performance Considerations
- Redis caching for active calls
- WebSocket connection pooling
- Efficient AMI event handling
- Bulk operation optimization

## Security Requirements
- Secure recording storage
- IP whitelist enforcement
- Rate limiting for bulk operations
- Audit trail for recharges

## Testing Requirements
- Gateway connection tests
- SMS parser unit tests
- WebSocket stress testing
- Bulk operation performance tests

## Documentation Needs
- Gateway setup guide
- Recharge system manual
- API documentation
- WebSocket protocol spec