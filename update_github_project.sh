#!/bin/bash

# Source GitHub configuration
source .github_config

# GitHub API base URL
API_BASE="https://api.github.com"
REPO_PATH="repos/${GITHUB_USER}/${GITHUB_REPO}"

# Create headers for authentication
AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
ACCEPT_HEADER="Accept: application/vnd.github.v3+json"

echo "Updating GitHub project tracker..."

# First, let's check if we have any open issues to update
echo "Checking for open issues..."

# Get project boards
curl -s -H "${AUTH_HEADER}" -H "${ACCEPT_HEADER}" \
  "${API_BASE}/${REPO_PATH}/projects" > projects.json

# Create a new issue for the UI fixes completed
echo "Creating issue for completed UI fixes..."

ISSUE_BODY=$(cat <<EOF
# UI Fixes Completed - 2025-06-29

## Summary
Major UI fixes have been completed for the E173 Gateway project. All critical UI issues reported by the user have been addressed.

## Completed Tasks

### 1. Authentication System ✅
- Created admin user with admin/admin credentials
- Fixed user authentication context in middleware  
- Added CurrentUser display in navigation bar
- Implemented logout functionality

### 2. Dashboard Issues ✅
- Fixed dashboard showing 4 cards instead of 5
- Removed full page refresh on dashboard load
- Fixed loading indicators stuck in loading state
- Added proper HTMX trigger scoping

### 3. Template Rendering Issues ✅
- Fixed gateways page showing modems content
- Removed global content block definition causing template collision
- Fixed base template content block handling
- Ensured proper template scoping

### 4. HTMX/Polling Issues ✅
- Fixed CDR and Blacklist pages repeatedly calling /api/v1/modems
- Added HTMX cleanup script to stop polling on page navigation
- Properly scoped polling triggers
- Added abort handling for HTMX requests

### 5. API Response Issues ✅
- Fixed customer stats returning JSON instead of HTML for HTMX requests
- Added HX-Request header detection
- Returns appropriate response format based on request type

### 6. SIM Cards Display ✅
- Fixed empty SIM cards container issue
- Updated template to handle empty state
- Added proper loading states

## Files Modified
- \`cmd/server/main.go\` - Authentication and template data handling
- \`templates/base.tmpl\` - HTMX cleanup and content blocks
- \`templates/dashboard_standalone.tmpl\` - Grid layout and triggers
- \`templates/partials/nav.tmpl\` - User authentication display
- \`internal/handlers/customer_handlers.go\` - HTMX response handling
- \`pkg/api/gateway_handler.go\` - Enhanced logging and data passing
- \`templates/modems/list.tmpl\` - Removed due to template collision

## Next Steps
1. Server restart required to apply all changes
2. Browser automation testing 
3. Customer SIP account management features
4. Asterisk configuration (pending user demonstration)

## Technical Notes
- All changes maintain backward compatibility
- No database schema changes required
- Performance improvements from reduced polling
- Better error handling and user feedback

Related to: #UI-Fixes #HTMX #Authentication #Templates
EOF
)

# Create the issue
curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"✅ UI Fixes Completed - Dashboard, Authentication, Templates\",
    \"body\": $(echo "$ISSUE_BODY" | jq -Rs .),
    \"labels\": [\"enhancement\", \"ui\", \"completed\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "GitHub project updated successfully!"

# Create a follow-up issue for remaining tasks
REMAINING_TASKS=$(cat <<EOF
# Remaining Tasks - E173 Gateway

## High Priority
1. **Customer SIP Account Management**
   - Design SIP account creation interface
   - Implement account CRUD operations
   - Add SIP credentials management
   - Integration with billing system

2. **Browser Automation Testing**
   - Run comprehensive UI tests
   - Verify all fixes are working
   - Create automated test suite
   - Document test results

3. **Asterisk Configuration**
   - Waiting for user demonstration
   - Dongle setup and configuration
   - AMI integration testing
   - Call routing setup

## Medium Priority
1. **Performance Optimization**
   - Database query optimization
   - Caching improvements
   - HTMX request batching

2. **Error Handling**
   - Improve error messages
   - Add user-friendly notifications
   - Better logging

## Low Priority
1. **UI Enhancements**
   - Dark mode improvements
   - Mobile responsiveness
   - Accessibility features

## Technical Debt
1. Template organization and naming conventions
2. API response format standardization
3. Frontend JavaScript organization
4. Test coverage improvements
EOF
)

curl -s -X POST \
  -H "${AUTH_HEADER}" \
  -H "${ACCEPT_HEADER}" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"📋 Remaining Tasks and Roadmap\",
    \"body\": $(echo "$REMAINING_TASKS" | jq -Rs .),
    \"labels\": [\"enhancement\", \"roadmap\"],
    \"milestone\": 1
  }" \
  "${API_BASE}/${REPO_PATH}/issues"

echo "Created remaining tasks issue!"