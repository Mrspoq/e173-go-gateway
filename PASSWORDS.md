# E173 Gateway - Password Management
*Development Environment - Temporary Password Storage*

⚠️  **SECURITY NOTE**: This file contains sensitive credentials. 
- Remove from production environments
- Add to .gitignore before committing
- Store securely and delete after setup

---

## 🔐 Master Password
**Default Password (All Systems):** `3omartel580`

---

## 📊 Database Credentials
- **Host:** localhost
- **Port:** 5432
- **Database:** e173_gateway
- **Username:** e173_user
- **Password:** `3omartel580`
- **Connection String:** `postgresql://e173_user:3omartel580@localhost:5432/e173_gateway`

---

## 📞 Asterisk AMI Credentials
- **Host:** localhost
- **Port:** 5038
- **Username:** admin
- **Password:** `3omartel580`
- **Connection String:** `admin:3omartel580@localhost:5038`

---

## 🖥️ System Access
- **Root Password:** `3omartel580`
- **System User:** root
- **SSH Access:** root@192.168.1.40

---

## 🌐 Application Credentials
- **Admin User:** admin
- **Admin Password:** `3omartel580`
- **API Key:** (auto-generated)
- **Session Secret:** `e173_secret_3omartel580`

---

## 🔄 Service Accounts
- **Database Service User:** e173_user
- **Database Service Password:** `3omartel580`
- **Application Service User:** e173app
- **Application Service Password:** `3omartel580`

---

## 📡 External Service Credentials
- **SMS Gateway API Key:** (to be configured)
- **WhatsApp API Token:** (to be configured)
- **Voice Recognition API Key:** (to be configured)

---

## 🛡️ Security Tokens
- **JWT Secret:** `e173_jwt_secret_3omartel580_$(date +%Y%m%d)`
- **Encryption Key:** `e173_encrypt_3omartel580`
- **Session Cookie Secret:** `e173_cookie_3omartel580`

---

## 🔧 Configuration Files Using These Credentials
- `.env` - Main application configuration
- `manager.conf` - Asterisk AMI configuration
- `pg_hba.conf` - PostgreSQL authentication
- `scripts/*.sh` - Backup and deployment scripts

---

## 🔄 Password Rotation Schedule
- **Development:** Manual (as needed)
- **Staging:** Monthly
- **Production:** Weekly + after incidents

---

## 🚨 Emergency Procedures
If passwords are compromised:
1. Run `./scripts/rotate_passwords.sh` (to be created)
2. Update all configuration files
3. Restart all services
4. Create new snapshots with updated credentials

---

**Last Updated:** $(date)
**Environment:** Development
**Status:** Active
