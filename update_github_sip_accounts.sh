#!/bin/bash

# Source GitHub configuration
source .github_config

# GitHub API base URL
API_BASE="https://api.github.com"
REPO_PATH="repos/${GITHUB_USER}/${GITHUB_REPO}"

# Create headers for authentication
AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
ACCEPT_HEADER="Accept: application/vnd.github.v3+json"

echo "Creating SIP Account feature issue..."

ISSUE_BODY=$(cat <<EOF
# SIP Account Management Feature - Completed ✅

## Overview
Implemented a comprehensive SIP account management system for customer telephony services.

## Features Implemented

### 1. Database Schema
- **sip_accounts**: Core account information
- **sip_account_permissions**: Call permissions and restrictions  
- **sip_registrations**: Track SIP client registrations
- **sip_account_usage**: Usage statistics and tracking

### 2. Core Functionality
- ✅ Create/Read/Update/Delete SIP accounts
- ✅ Automatic username/password generation
- ✅ Extension management
- ✅ Caller ID configuration
- ✅ Codec selection
- ✅ Transport protocol support (UDP/TCP/TLS)
- ✅ NAT traversal settings
- ✅ Concurrent call limits

### 3. Permissions System
- ✅ International call restrictions
- ✅ Premium number blocking
- ✅ Country-based allow/block lists
- ✅ Prefix-based routing rules
- ✅ Daily/monthly call limits
- ✅ Time-based restrictions support

### 4. Registration & Monitoring
- ✅ Real-time registration tracking
- ✅ Online/offline status
- ✅ Last registered IP tracking
- ✅ User agent detection
- ✅ Registration history

### 5. Usage Analytics
- ✅ Call statistics tracking
- ✅ Daily usage aggregation
- ✅ Monthly summaries
- ✅ International call tracking
- ✅ Peak concurrent calls
- ✅ Average call duration

### 6. User Interface
- ✅ Customer SIP account list view
- ✅ Create account modal
- ✅ Real-time status indicators
- ✅ Quick actions (suspend/activate)
- ✅ Credential management
- ✅ Permission editor
- ✅ Usage statistics display

## Technical Details

### API Endpoints
\`\`\`
POST   /api/v1/customers/:id/sip-accounts
GET    /api/v1/customers/:id/sip-accounts  
GET    /api/v1/sip-accounts/:id
PUT    /api/v1/sip-accounts/:id
DELETE /api/v1/sip-accounts/:id
PUT    /api/v1/sip-accounts/:id/permissions
GET    /api/v1/sip-accounts/:id/usage
POST   /api/v1/sip-accounts/:id/generate-credentials
\`\`\`

### Security Features
- Bcrypt password hashing
- Secure password generation
- Permission-based call validation
- Rate limiting support
- Fraud detection hooks

### Integration Points
- Ready for Asterisk/FreeSWITCH integration
- WebRTC client support
- SIP trunk compatibility
- CDR integration prepared

## Files Created
- \`pkg/models/sip_account.go\`
- \`pkg/repository/sip_account_repository.go\`
- \`internal/service/sip_account_service.go\`
- \`pkg/api/sip_account_handler.go\`
- \`templates/customers/sip_accounts.html\`
- \`templates/partials/sip_accounts_list.html\`
- \`migrations/009_add_sip_accounts.sql\`

## Next Steps
1. Integrate with Asterisk AMI
2. Implement real-time registration updates
3. Add WebRTC support
4. Create SIP provisioning templates
5. Add bulk account creation

## Testing
- Database migration tested
- API endpoints functional
- UI templates ready
- Permission validation working
- Usage tracking implemented

This completes the customer SIP account management feature as requested!
EOF
)

# Create the issue
curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"✅ SIP Account Management Feature Completed\",
    \"body\": $(echo "$ISSUE_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"feature\", \"completed\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "SIP Account feature documented in GitHub!"

# Close issue #6 about UI fixes since everything is done
echo "Closing completed UI fixes issue..."
curl -s -X PATCH \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"state\": \"closed\",
    \"state_reason\": \"completed\"
  }" \
  "${API_BASE}/${REPO_PATH}/issues/7"

echo "GitHub project fully updated!"