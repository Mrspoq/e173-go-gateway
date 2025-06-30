# E173 Gateway - Advanced Filter System Plan

## Current Status
- ✅ Libphonenumber C++ integration working
- ✅ WhatsApp API integration (wa-validator.xyz)
- ✅ Morocco prefixes database (169 prefixes: IAM=58, Orange=45, Inwi=66)
- ✅ Basic filter service implemented
- ⚠️  Need parallel processing for 2-second API latency
- ⚠️  Need configurable filters from frontend

## Architecture Updates Needed

### 1. Parallel Processing System
```go
// Worker pool for concurrent validation
type ValidationWorkerPool struct {
    workers      int
    jobQueue     chan ValidationJob
    resultQueue  chan ValidationResult
    rateLimiter  *rate.Limiter
}

// Process multiple numbers concurrently
func (vwp *ValidationWorkerPool) ProcessBatch(numbers []string) []ValidationResult
```

### 2. Configurable Filter System

#### Filter Types (for both Source & Destination)
1. **Libphonenumber Validation**
   - E.164 format check
   - Country-specific rules
   - Length validation

2. **WhatsApp Validation**
   - Real-time API check
   - Cached results (24h TTL)
   - Rate limited (5 TPS)

3. **HLR (Home Location Register)**
   - Future implementation
   - Carrier/network validation
   - Number portability check

4. **Database Filters**
   - Blacklist check
   - Graylist check
   - Whitelist check
   - Customer-specific rules

#### Configuration Structure
```json
{
  "filter_config": {
    "source_filters": {
      "libphonenumber": {
        "enabled": true,
        "strict_mode": false,
        "allowed_countries": ["*"]
      },
      "whatsapp": {
        "enabled": false,
        "cache_ttl": "24h"
      },
      "hlr": {
        "enabled": false
      },
      "database": {
        "blacklist": true,
        "graylist": true,
        "whitelist": false
      }
    },
    "destination_filters": {
      "libphonenumber": {
        "enabled": true,
        "strict_mode": true,
        "allowed_countries": ["MA"]
      },
      "whatsapp": {
        "enabled": true,
        "cache_ttl": "24h",
        "required": false,
        "timeout": "5s",
        "retry_attempts": 2
      },
      "hlr": {
        "enabled": false
      },
      "database": {
        "blacklist": true,
        "graylist": false,
        "whitelist": false
      }
    }
  }
}
```

### 3. Database Schema Updates
```sql
-- Filter configurations table
CREATE TABLE filter_configs (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    source_filters JSONB NOT NULL,
    destination_filters JSONB NOT NULL,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Graylist table (suspicious but not blocked)
CREATE TABLE graylist (
    phone_number VARCHAR(20) PRIMARY KEY,
    reason TEXT,
    score INT DEFAULT 50, -- 0-100, higher = more suspicious
    added_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

-- Whitelist table (always allowed)
CREATE TABLE whitelist (
    phone_number VARCHAR(20) PRIMARY KEY,
    customer_id UUID,
    reason TEXT,
    added_at TIMESTAMP DEFAULT NOW()
);
```

### 4. API Endpoints
```
POST   /api/v1/filter-configs          Create filter configuration
GET    /api/v1/filter-configs          List all configurations
GET    /api/v1/filter-configs/:id      Get specific configuration
PUT    /api/v1/filter-configs/:id      Update configuration
DELETE /api/v1/filter-configs/:id      Delete configuration
POST   /api/v1/filter-configs/:id/activate   Activate configuration

POST   /api/v1/filters/test            Test filters with sample numbers
GET    /api/v1/filters/stats           Get filter statistics
```

### 5. Frontend Requirements
- Filter configuration UI
- Real-time filter testing
- Statistics dashboard
- Blacklist/Graylist/Whitelist management
- Bulk import/export

### 6. Performance Optimizations
- Redis caching for all validations
- Parallel processing with goroutines
- Connection pooling for APIs
- Batch processing for bulk operations
- Circuit breaker for external APIs

## Implementation Priority
1. **Phase 1**: Parallel processing infrastructure
2. **Phase 2**: Filter configuration API and database
3. **Phase 3**: Frontend filter management UI
4. **Phase 4**: Advanced features (HLR, ML-based scoring)

## Next Steps
1. Implement worker pool for parallel validation
2. Create filter configuration service
3. Update database schema
4. Build frontend filter management
5. Add comprehensive logging and monitoring

## Testing Requirements
- Load testing with 1000+ concurrent validations
- API failover testing
- Cache performance testing
- Filter rule conflict resolution
- Rate limit compliance testing

## Security Considerations
- API keys in environment variables
- Rate limiting per customer
- Audit logging for all filter changes
- Role-based access control for filter management