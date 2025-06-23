# 🚀 E173 Gateway: Deployment & Backup Strategy

## 📋 **Overview**

This document describes the **coordinated snapshot system** that ties Git commits to database backups, enabling complete system state management for deployment, recovery, and rollback operations.

## 🎯 **Core Concept: Coordinated Snapshots**

Each **Git commit** is paired with a **database backup** using the same reference ID, ensuring complete system state preservation:

```
Git Commit: 4c7cec7 ←→ Database Backup: e173_gateway_4c7cec7.sql.gz
```

This enables:
- ✅ **Complete rollback** to any previous state
- ✅ **One-click deployment** with matching code + data
- ✅ **Zero-configuration recovery** with all settings preserved
- ✅ **Cross-server migration** with identical functionality

## 📁 **File Structure**

```
e173_go_gateway/
├── docs/
│   └── DEPLOYMENT_AND_BACKUP_STRATEGY.md   # This document
├── scripts/
│   ├── snapshot_create.sh                  # Create coordinated snapshot
│   ├── snapshot_restore.sh                 # Restore to specific snapshot
│   ├── deploy_fresh.sh                     # Deploy to new server
│   ├── rollback.sh                         # Rollback to previous state
│   └── backup_database.sh                  # Manual database backup
├── snapshots/
│   ├── commits.log                         # Commit-backup mapping
│   └── backups/                            # Database backups by commit
│       ├── e173_gateway_4c7cec7.sql.gz
│       └── e173_gateway_780049e.sql.gz
└── deployment/
    ├── .env.template                       # Environment template
    └── config/                             # Deployment configurations
```

## 🔧 **System Components**

### **1. Snapshot Creation**
Creates coordinated Git commit + database backup:
```bash
./scripts/snapshot_create.sh "Added customer management"
# → Git commit: a1b2c3d + Database: e173_gateway_a1b2c3d.sql.gz
```

### **2. Snapshot Restoration**  
Restores complete system state to specific commit:
```bash
./scripts/snapshot_restore.sh a1b2c3d
# → Restores code + database + configuration to exact state
```

### **3. Fresh Deployment**
Deploys complete system to new server:
```bash
./scripts/deploy_fresh.sh a1b2c3d
# → Sets up everything from scratch with specified snapshot
```

### **4. Rollback System**
Quick rollback to previous working state:
```bash
./scripts/rollback.sh
# → Automatically reverts to last known good snapshot
```

## 📊 **Snapshot Lifecycle**

### **Development Workflow**
```bash
# 1. Make changes to code/database
vim pkg/api/handlers.go
make run

# 2. Test functionality
curl http://localhost:8080/api/stats

# 3. Create coordinated snapshot
./scripts/snapshot_create.sh "Feature: Enhanced gateway stats"
# → Commit: f5e4d3c + Backup: e173_gateway_f5e4d3c.sql.gz

# 4. Push to GitHub (optional)
git push origin master
```

### **Deployment Workflow**
```bash
# Deploy specific snapshot to production
./scripts/deploy_fresh.sh f5e4d3c production-server.com

# Or rollback if issues found
./scripts/rollback.sh
```

## 🛠️ **Command Reference**

### **Snapshot Management**

| Command | Description | Example |
|---------|-------------|---------|
| `snapshot_create.sh "<message>"` | Create coordinated snapshot | `./scripts/snapshot_create.sh "Added auth system"` |
| `snapshot_restore.sh <commit>` | Restore to specific snapshot | `./scripts/snapshot_restore.sh a1b2c3d` |
| `snapshot_list.sh` | List all snapshots | `./scripts/snapshot_list.sh` |
| `snapshot_cleanup.sh <days>` | Clean old snapshots | `./scripts/snapshot_cleanup.sh 30` |

### **Deployment Operations**

| Command | Description | Example |
|---------|-------------|---------|
| `deploy_fresh.sh <commit> [server]` | Deploy to new server | `./scripts/deploy_fresh.sh a1b2c3d prod.example.com` |
| `deploy_update.sh <commit>` | Update existing deployment | `./scripts/deploy_update.sh a1b2c3d` |
| `rollback.sh [steps]` | Rollback to previous state | `./scripts/rollback.sh 2` |

### **Maintenance Operations**

| Command | Description | Example |
|---------|-------------|---------|
| `backup_database.sh` | Manual database backup | `./scripts/backup_database.sh` |
| `verify_snapshot.sh <commit>` | Verify snapshot integrity | `./scripts/verify_snapshot.sh a1b2c3d` |
| `migrate_server.sh <source> <target>` | Migrate between servers | `./scripts/migrate_server.sh old.com new.com` |

## 🔄 **Deployment Scenarios**

### **Scenario 1: Fresh Production Deployment**
```bash
# On new production server
git clone https://github.com/Mrspoq/e173-go-gateway.git
cd e173-go-gateway

# Deploy specific snapshot (includes all setup)
./scripts/deploy_fresh.sh f5e4d3c

# System is ready with:
# ✅ Database created and populated
# ✅ Environment configured  
# ✅ Services started
# ✅ All dependencies installed
```

### **Scenario 2: Update Existing System**
```bash
# On existing server
cd e173-go-gateway

# Update to new snapshot
./scripts/deploy_update.sh g6f7e8d

# System updated with:
# ✅ Code updated to new version
# ✅ Database migrated/updated
# ✅ Services restarted
# ✅ Configuration preserved
```

### **Scenario 3: Emergency Rollback**
```bash
# If issues found in production
./scripts/rollback.sh

# Or rollback specific number of versions
./scripts/rollback.sh 3

# System reverted with:
# ✅ Code reverted to previous version
# ✅ Database restored to matching state
# ✅ Services restarted with old config
# ✅ All functionality restored
```

### **Scenario 4: Cross-Server Migration**
```bash
# Migrate from old to new server
./scripts/migrate_server.sh old-server.com new-server.com

# Complete migration includes:
# ✅ All snapshots copied
# ✅ Database migrated
# ✅ Services configured
# ✅ DNS updated (if configured)
```

## ⚙️ **Configuration Management**

### **Environment Templates**
Each snapshot can include environment configuration:
```bash
deployment/
├── .env.template                    # Base environment template
├── config/
│   ├── development.env             # Development settings
│   ├── staging.env                 # Staging settings  
│   └── production.env              # Production settings
```

### **Automatic Configuration**
During deployment, the system automatically:
1. **Detects environment** (dev/staging/production)
2. **Applies appropriate config** from templates
3. **Generates secrets** if needed
4. **Updates DNS/services** if configured

## 🔒 **Security & Best Practices**

### **Backup Encryption**
```bash
# Database backups are automatically encrypted
# Encryption key derived from deployment environment
export BACKUP_ENCRYPTION_KEY="your-secure-key"
```

### **Credential Management**
```bash
# Credentials stored separately from code
snapshots/
├── credentials/
│   ├── development.credentials     # Dev credentials
│   ├── staging.credentials         # Staging credentials
│   └── production.credentials      # Prod credentials (encrypted)
```

### **Access Control**
```bash
# Snapshot access controlled by user permissions
chmod 700 scripts/snapshot_*.sh     # Owner only
chmod 600 snapshots/credentials/*   # Credentials protected
```

## 📈 **Monitoring & Verification**

### **Snapshot Integrity**
Each snapshot includes verification checksums:
```bash
# Verify snapshot integrity
./scripts/verify_snapshot.sh a1b2c3d
# ✅ Git commit exists and matches
# ✅ Database backup exists and is valid
# ✅ Configuration files are present
# ✅ Checksums match expected values
```

### **Deployment Health**
Post-deployment verification:
```bash
# Health check after deployment/rollback
./scripts/health_check.sh
# ✅ Database connectivity
# ✅ API endpoints responding
# ✅ Services running
# ✅ Asterisk AMI connected
```

## 🚨 **Emergency Procedures**

### **System Recovery**
If system is completely broken:
```bash
# Nuclear option: complete rebuild from snapshot
./scripts/emergency_rebuild.sh a1b2c3d

# This will:
# ⚠️  Completely wipe current installation
# 🔄 Rebuild from scratch using snapshot
# ✅ Restore all data and configuration
# ✅ Restart all services
```

### **Data Loss Recovery**
If only database is corrupted:
```bash
# Restore only database from snapshot
./scripts/restore_database_only.sh a1b2c3d

# This will:
# 🗄️  Drop and recreate database
# 📥 Restore from snapshot backup
# 🔄 Run any necessary migrations
# ✅ Verify data integrity
```

## 📋 **Snapshot Metadata**

Each snapshot includes comprehensive metadata:
```json
{
  "commit_id": "a1b2c3d",
  "timestamp": "2025-01-23T17:30:00Z",
  "message": "Added customer management system",
  "author": "team@e173gateway.com",
  "database_backup": "e173_gateway_a1b2c3d.sql.gz",
  "backup_size": "1.2MB",
  "tables_count": 10,
  "records_count": 15847,
  "environment": "production",
  "dependencies": {
    "go_version": "1.18",
    "postgresql_version": "13.7",
    "asterisk_version": "18.x"
  },
  "checksums": {
    "code": "sha256:abc123...",
    "database": "sha256:def456...",
    "config": "sha256:ghi789..."
  },
  "tags": ["stable", "feature-complete"],
  "notes": "Customer management system fully implemented and tested"
}
```

---

## 🎯 **Quick Start Guide**

### **For Developers:**
```bash
# Create your first coordinated snapshot
./scripts/snapshot_create.sh "Initial stable version"

# List all snapshots
./scripts/snapshot_list.sh

# Test rollback (safe - can rollback again)
./scripts/rollback.sh
```

### **For Deployment:**
```bash
# Deploy to production server
./scripts/deploy_fresh.sh <latest-commit> production-server.com

# Monitor deployment
./scripts/health_check.sh

# Rollback if needed
./scripts/rollback.sh
```

### **For Maintenance:**
```bash
# Clean old snapshots (keep last 30 days)
./scripts/snapshot_cleanup.sh 30

# Verify system integrity
./scripts/verify_snapshot.sh $(git rev-parse HEAD)

# Check backup status
./scripts/backup_status.sh
```

---

**🎉 With this system, you have industrial-grade deployment and recovery capabilities that ensure your E173 Gateway system can be perfectly reproduced, updated, and rolled back with complete confidence.**
