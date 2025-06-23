# E173 Gateway Deployment & Backup System

## 🎯 One-Command Operations

### Create Snapshot
```bash
# Create coordinated Git + Database snapshot
./scripts/snapshot_create.sh "Stable production version v1.2"
```

### Deploy Anywhere
```bash
# Deploy to new server with one command
./scripts/deploy_fresh.sh a1b2c3d production-server.com production

# Deploy locally
./scripts/deploy_fresh.sh a1b2c3d localhost production
```

### Quick Rollback
```bash
# Rollback to previous snapshot
./scripts/rollback.sh

# Rollback multiple steps
./scripts/rollback.sh 3
```

### Restore Specific Snapshot
```bash
# Restore to any snapshot
./scripts/snapshot_restore.sh a1b2c3d
```

## 📋 Snapshot Management

### List All Snapshots
```bash
./scripts/snapshot_list.sh
```

### Cleanup Old Snapshots  
```bash
# Remove snapshots older than 30 days
./scripts/snapshot_cleanup.sh

# Custom retention period
./scripts/snapshot_cleanup.sh 7   # Keep last 7 days
./scripts/snapshot_cleanup.sh 90  # Keep last 90 days
```

## 🏗️ System Architecture

### Coordinated Snapshots
Each snapshot consists of:
- **Git Commit**: Exact code state
- **Database Backup**: Complete PostgreSQL dump
- **Metadata**: Timestamp, size, record count
- **Verification**: Integrity checks

### Directory Structure
```
snapshots/
├── commits.log              # Snapshot registry
└── backups/
    ├── e173_gateway_a1b2c3d.sql.gz      # Database backup
    └── e173_gateway_a1b2c3d.meta.json   # Metadata
```

### State Synchronization
- Git commits and database backups are perfectly synchronized
- Each snapshot represents a complete, deployable system state
- Zero data loss between snapshots

## 🚀 Deployment Scenarios

### 1. Fresh Server Deployment
```bash
# Complete setup on new Ubuntu/CentOS server
./scripts/deploy_fresh.sh a1b2c3d your-server.com production
```

**What happens:**
1. Installs dependencies (PostgreSQL, Go, etc.)
2. Clones repository to `/opt/e173_go_gateway`  
3. Sets up database and user
4. Configures environment variables
5. Builds application
6. Restores snapshot data
7. Starts services
8. Verifies deployment

### 2. Local Development Reset
```bash
# Reset local environment to clean state
./scripts/snapshot_restore.sh a1b2c3d
```

### 3. Emergency Rollback
```bash
# Instant rollback when something breaks
./scripts/rollback.sh
```

### 4. Multi-Environment Management
```bash
# Production deployment
./scripts/deploy_fresh.sh v1.2.0 prod-server.com production

# Staging deployment  
./scripts/deploy_fresh.sh v1.3.0-beta staging-server.com staging

# Development reset
./scripts/snapshot_restore.sh develop-branch
```

## 🔧 Manual Operations

### Create Database Backup Only
```bash
./scripts/backup_database.sh
```

### Restore Database Only
```bash
./scripts/restore_database.sh /path/to/backup.sql.gz
```

### Setup Automated Backups
```bash
./scripts/automated_backup_setup.sh
```

## 🛡️ Safety Features

### Emergency Backups
- Automatic emergency backup before any destructive operation
- Uncommitted changes preserved before rollbacks
- Emergency commits tagged for easy recovery

### Verification Checks
- Database connectivity verification
- Server response validation
- Integrity checks on all operations

### Rollback Protection
- Confirmation prompts for destructive operations
- Multiple rollback levels (1 step, 3 steps, etc.)
- Emergency restore procedures

## 📊 Best Practices

### Regular Snapshots
```bash
# Before major changes
./scripts/snapshot_create.sh "Before implementing customer billing"

# After successful deployments
./scripts/snapshot_create.sh "Production v1.2 - stable"

# Before risky operations
./scripts/snapshot_create.sh "Before database migration"
```

### Naming Convention
- Use semantic versioning: `v1.2.0`, `v1.2.1-hotfix`
- Include feature descriptions: `customer-management-mvp`
- Mark stability: `stable-production`, `beta-testing`

### Retention Management
```bash
# Weekly cleanup (recommended in cron)
0 2 * * 0 /opt/e173_go_gateway/scripts/snapshot_cleanup.sh 30
```

## 🔍 Troubleshooting

### Deployment Fails
```bash
# Check server logs
ssh your-server.com 'cd /opt/e173_go_gateway && make logs'

# Manual verification
ssh your-server.com 'curl localhost:8080/ping'
```

### Rollback Issues
```bash
# Force restore to known good state
./scripts/snapshot_restore.sh <last_known_good_commit>

# Check emergency backups
git log --oneline | grep "Emergency backup"
```

### Database Problems
```bash
# Verify database connectivity
psql -h localhost -U e173_user -d e173_gateway -c "SELECT 1;"

# Check database size
./scripts/snapshot_list.sh
```

### Space Management
```bash
# Check snapshot storage usage
du -sh snapshots/

# Aggressive cleanup if needed
./scripts/snapshot_cleanup.sh 7
```

## 🌐 GitHub Integration

### Private Repository
- Repository URL: `https://github.com/Mrspoq/e173-go-gateway`
- Make private in GitHub settings for security
- All snapshots automatically synchronized

### Collaboration Workflow
```bash
# Team member setup
git clone https://github.com/Mrspoq/e173-go-gateway.git
cd e173-go-gateway
./scripts/snapshot_restore.sh <latest_stable>
```

### CI/CD Integration
- Automated snapshots on successful builds
- Deployment scripts ready for CI/CD pipelines
- Zero-downtime deployment capability

## 💡 Advanced Usage

### Multi-Server Deployment
```bash
# Deploy same snapshot to multiple servers
for server in prod1.com prod2.com prod3.com; do
    ./scripts/deploy_fresh.sh v1.2.0 $server production
done
```

### Snapshot Comparison
```bash
# Compare snapshots
./scripts/snapshot_list.sh | head -5
git diff commit1..commit2
```

### Disaster Recovery
```bash
# Complete disaster recovery procedure
./scripts/deploy_fresh.sh <last_good_commit> <new_server> production
```

---

**🎉 Result: Enterprise-grade deployment system with zero-downtime capabilities, perfect synchronization between code and data, and one-command operations for all scenarios.**
