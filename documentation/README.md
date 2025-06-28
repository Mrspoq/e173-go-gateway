# E173 Go Gateway

A telecommunications management platform for handling ~200 Huawei E173 USB modems with VoIP routing, SMS handling, and SIM card management.

## Quick Start

### Prerequisites
- Go 1.18+
- PostgreSQL 13+
- Node.js 16+ (for Tailwind CSS)
- Asterisk 18 with chan_dongle

### Setup
```bash
# 1. Clone and enter directory
git clone <repository>
cd e173_go_gateway

# 2. Setup database
sudo -u postgres psql
CREATE DATABASE e173_gateway;
CREATE USER e173_user WITH PASSWORD 'e173_pass';
GRANT ALL PRIVILEGES ON DATABASE e173_gateway TO e173_user;
\q

# 3. Configure environment
cp .env.example .env
# Edit .env with your database credentials

# 4. Build and run
make migrate     # Run database migrations
make tailwind    # Build CSS
make build       # Build binary
make run         # Start server
```

### Access
- Web UI: http://localhost:8080
- API: http://localhost:8080/api

## Project Structure

```
e173_go_gateway/
├── cmd/server/main.go      # Application entry point
├── internal/               # Private application code
│   ├── database/          # Models and migrations
│   ├── handlers/          # HTTP handlers
│   ├── repository/        # Data access layer
│   └── service/          # Business logic
├── pkg/                   # Public packages
│   ├── ami/              # Asterisk integration
│   ├── auth/             # Authentication
│   ├── config/           # Configuration
│   ├── models/           # Shared models
│   └── repository/       # Repository interfaces
├── templates/             # HTMX templates
├── web/static/           # Static assets
└── migrations/           # Database migrations
```

## Current Status

### Working
- Backend API endpoints
- Database schema and migrations
- AMI integration with Asterisk
- Real-time dashboard statistics
- Basic HTMX template structure

### In Progress
- Frontend template rendering fixes
- Navigation between UI sections
- WebSocket real-time updates

### Planned
- Operator prefix routing
- Bulk SIM recharge
- Voice recognition bot detection
- YAML-based automation
- OpenSIPS integration

## Development

### API Endpoints
- `GET /api/stats` - Dashboard statistics
- `GET /api/customers` - Customer management
- `GET /api/modems` - Modem management
- `GET /api/sims` - SIM card management
- `GET /api/cdrs` - Call detail records

### Make Commands
```bash
make build       # Build binary
make run         # Run development server
make test        # Run tests
make migrate     # Run database migrations
make tailwind    # Build Tailwind CSS
make clean       # Clean build artifacts
```

### Testing
```bash
# Test API
curl http://localhost:8080/api/stats

# Test UI sections
curl http://localhost:8080/customers
curl http://localhost:8080/modems
curl http://localhost:8080/sims
```

## Multi-Agent Collaboration

See [CODEBASE_OVERVIEW.md](CODEBASE_OVERVIEW.md) for detailed architecture documentation.

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

### Key Files
- **Server**: `cmd/server/main.go`
- **Routes**: Look for `router.GET/POST` in main.go
- **Templates**: `templates/` directory
- **Database**: `internal/database/models/`
- **Config**: `.env` and `pkg/config/`

### Current Issues
1. Some templates missing `{{template "base" .}}` directive
2. Navigation links need fixing
3. WebSocket handlers not implemented

## License

Private project - All rights reserved
