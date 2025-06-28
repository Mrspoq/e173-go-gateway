# E173 Go Gateway - Codebase Overview

## 🚨 CURRENT PROJECT STATUS (June 26, 2025)

### ✅ WORKING COMPONENTS
1. **Authentication System (FULLY FUNCTIONAL)**
   - JWT-based authentication implemented and tested
   - Login endpoint `/login` working with `admin/admin` credentials
   - Session management via cookies (`session_token`)
   - Protected route middleware functioning correctly
   - Admin user exists in `e173_gateway` database with correct password hash

2. **Backend API**
   - All CRUD endpoints implemented (modems, SIMs, CDRs, customers, blacklist)
   - PostgreSQL database properly configured with all tables
   - Repository pattern implemented for data access
   - Asterisk AMI integration configured

3. **Basic Routing**
   - All main routes defined and accessible
   - Authentication redirects working (unauthenticated → /login)

### ❌ FRONTEND ISSUES REQUIRING IMMEDIATE FIX
1. **CSS/Styling Not Loading**
   - `/static/bundle.css` exists and serves (24KB) but doesn't render
   - No colors, layouts, or styling visible in browser
   - Likely path issue when accessing via LAN IP (192.168.1.35)

2. **Template Content Collision**
   - Dashboard route (`/`) shows "System Settings" content instead of Dashboard
   - Multiple templates use same `{{define "content"}}` block name
   - Go template engine uses last loaded template, causing overrides

3. **Missing Dashboard Features**
   - Real-time statistics not displaying
   - No modem status visualization
   - No call activity charts
   - System health metrics not shown

### 🔧 IMMEDIATE NEXT STEPS
1. Fix CSS loading issue for proper UI rendering
2. Resolve template collisions (use unique block names or standalone templates)
3. Implement dashboard data visualization
4. Complete HTMX integration for dynamic updates
5. Add real-time WebSocket updates for live data

## Project Overview

The E173 Go Gateway is a telecommunications management platform designed to handle ~200 Huawei E173 USB modems for voice call routing, SMS handling, and SIM card management. It's built with Go (Gin framework) for the backend and HTMX + Tailwind CSS for the frontend.

## Architecture

### Core Technologies
- **Backend**: Go 1.18+ with Gin web framework
- **Frontend**: HTMX + Tailwind CSS (server-side rendered templates)
- **Database**: PostgreSQL 13+
- **VoIP Integration**: Asterisk 18 with chan_dongle
- **Message Queue**: Redis (planned)
- **SIP Signaling**: OpenSIPS (future phase)

### Project Structure

```
e173_go_gateway/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── database/           # Database models and migrations
│   ├── handlers/           # HTTP request handlers
│   ├── repository/         # Data access layer
│   └── service/           # Business logic layer
├── pkg/
│   ├── ami/               # Asterisk Manager Interface integration
│   ├── auth/              # Authentication & authorization
│   ├── config/            # Configuration management
│   ├── models/            # Shared data models
│   └── repository/        # Repository interfaces
├── templates/             # HTMX templates
│   ├── base.tmpl         # Base layout
│   ├── dashboard.tmpl    # Main dashboard
│   ├── customers/        # Customer management views
│   ├── modems/          # Modem management views
│   ├── sims/            # SIM card management views
│   ├── cdrs/            # Call detail records views
│   ├── blacklist/       # Blacklist management views
│   └── settings/        # Settings views
├── web/
│   └── static/          # Static assets (CSS, JS, images)
├── migrations/          # Database migrations
└── scripts/            # Utility scripts

```

## Key Components

### 1. Backend Services

#### Authentication Service (`pkg/auth/`)
- JWT-based authentication
- Role-based access control (Admin/User)
- Session management

#### AMI Service (`pkg/ami/service.go`)
- Connects to Asterisk via AMI protocol
- Monitors call events
- Manages modem channels
- Ingests CDR data in real-time

#### Repository Layer (`internal/repository/`)
- `CustomerRepository`: Customer CRUD operations
- `ModemRepository`: Modem state management
- `SimCardRepository`: SIM card tracking
- `CdrRepository`: Call detail records
- `GatewayRepository`: Gateway configuration
- `UserRepository`: User management

### 2. Frontend Components

#### Dashboard (`templates/dashboard.tmpl`)
- Real-time statistics cards
- Live updates via HTMX polling
- Displays: Active modems, SIM cards, calls today, spam blocked

#### Customer Management
- List view with search/filter
- Create/Edit forms
- Balance tracking
- Rate sheet management

#### Modem Management
- USB modem status monitoring
- Enable/disable controls
- Signal strength indicators
- Carrier information

#### SIM Card Management
- Bulk recharge capabilities
- Balance tracking
- Usage monitoring
- Operator prefix routing

### 3. Database Schema

#### Core Tables
- `users`: Authentication and user profiles
- `customers`: Customer accounts with rate sheets
- `modems`: Physical E173 modem devices
- `sim_cards`: SIM card inventory and status
- `call_detail_records`: CDR data from Asterisk
- `gateways`: Gateway configuration
- `blacklist`: Blocked phone numbers
- `routing_rules`: Operator prefix routing (planned)

## API Endpoints

### Stats API
- `GET /api/stats` - Dashboard statistics

### Customer API
- `GET /api/customers` - List customers
- `POST /api/customers` - Create customer
- `PUT /api/customers/:id` - Update customer
- `DELETE /api/customers/:id` - Delete customer

### Modem API
- `GET /api/modems` - List modems
- `GET /api/modems/:id` - Get modem details
- `POST /api/modems/:id/enable` - Enable modem
- `POST /api/modems/:id/disable` - Disable modem

### SIM Card API
- `GET /api/sims` - List SIM cards
- `POST /api/sims/bulk-recharge` - Bulk recharge

## Configuration

### Environment Variables (.env)
```
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=e173_gateway
DB_USER=e173_user
DB_PASSWORD=your_password
DB_SSLMODE=disable

# Asterisk AMI Configuration
AMI_HOST=localhost
AMI_PORT=5038
AMI_USERNAME=admin
AMI_PASSWORD=secret

# JWT Configuration
JWT_SECRET=your_jwt_secret
```

## Build & Development

### Prerequisites
- Go 1.18+
- PostgreSQL 13+
- Node.js 16+ (for Tailwind CSS)
- Make

### Build Commands
```bash
make build        # Build the binary
make run          # Run in development mode
make migrate      # Run database migrations
make tailwind     # Build Tailwind CSS
make clean        # Clean build artifacts
```

## Current Status

### Completed Features
- ✅ Basic authentication system
- ✅ Customer CRUD operations
- ✅ Database schema and migrations
- ✅ AMI integration for Asterisk
- ✅ Real-time dashboard stats
- ✅ HTMX template structure
- ✅ Repository pattern implementation

### In Progress
- 🔄 Frontend template rendering fixes
- 🔄 Navigation and UI polish
- 🔄 WebSocket real-time updates

### Planned Features
- 📋 Operator prefix routing
- 📋 Bulk SIM recharge UI
- 📋 Voice recognition for bot detection
- 📋 YAML-based USSD/SMS automation
- 📋 Short-call spam detection
- 📋 IVR autoresponder
- 📋 OpenSIPS integration
- 📋 WhatsApp validation API

## Development Guidelines

### Code Organization
- Follow standard Go project layout
- Use dependency injection
- Implement repository pattern for data access
- Keep business logic in service layer
- Use HTMX for dynamic UI updates

### Testing
- Unit tests for repository layer
- Integration tests for API endpoints
- Use testify for assertions
- Mock external dependencies

### Error Handling
- Use custom error types
- Log errors with context
- Return appropriate HTTP status codes
- Provide user-friendly error messages

## Debugging Tips

### Common Issues
1. **Template not rendering**: Check template directive `{{template "base" .}}`
2. **AMI connection failed**: Verify Asterisk credentials in .env
3. **Database errors**: Check PostgreSQL connection and migrations
4. **HTMX not updating**: Verify hx-* attributes and endpoints

### Useful Commands
```bash
# Check server logs
tail -f server.log

# Test API endpoints
curl http://localhost:8080/api/stats

# Verify database connection
psql -U e173_user -d e173_gateway

# Check Asterisk AMI
telnet localhost 5038
```

## Multi-Agent Collaboration Notes

### Working Areas
- **Agent 1**: Frontend templates and HTMX integration
- **Agent 2**: Backend API development
- **Agent 3**: Database and repository layer
- **Agent 4**: AMI and telephony integration

### Coordination Points
- Database schema changes require migration files
- API changes need corresponding frontend updates
- New features should include both backend and frontend
- Always update this documentation when adding major features

### Branch Strategy
- `main`: Stable production code
- `develop`: Integration branch
- `feature/*`: Individual feature branches
- `fix/*`: Bug fix branches

## Next Steps

1. Fix template rendering issues in all UI sections
2. Implement WebSocket for real-time updates
3. Complete SIM bulk recharge functionality
4. Add comprehensive filtering capabilities
5. Integrate OpenSIPS for advanced call routing
