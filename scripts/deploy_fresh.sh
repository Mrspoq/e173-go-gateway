#!/bin/bash

# E173 Gateway: Fresh Deployment Script  
# Deploys complete system to new server with one command

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SNAPSHOTS_DIR="$PROJECT_DIR/snapshots"
BACKUPS_DIR="$SNAPSHOTS_DIR/backups"

# Check arguments
if [ $# -eq 0 ]; then
    echo "❌ Usage: $0 <commit_id> [server_address] [environment]"
    echo ""
    echo "📋 Examples:"
    echo "   $0 a1b2c3d                           # Deploy locally"
    echo "   $0 a1b2c3d production-server.com     # Deploy to remote server"
    echo "   $0 a1b2c3d localhost production      # Deploy locally as production"
    echo ""
    exit 1
fi

COMMIT_ID="$1"
SERVER_ADDRESS="${2:-localhost}"
ENVIRONMENT="${3:-production}"

echo "🚀 E173 Gateway Fresh Deployment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Commit:      $COMMIT_ID"
echo "🖥️  Server:      $SERVER_ADDRESS"
echo "🏷️  Environment: $ENVIRONMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Verify snapshot exists
BACKUP_FILE="$BACKUPS_DIR/e173_gateway_${COMMIT_ID}.sql.gz"
if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Snapshot backup not found: $BACKUP_FILE"
    echo "💡 Available snapshots:"
    ls -la "$BACKUPS_DIR"/*.sql.gz 2>/dev/null | head -5
    exit 1
fi

# Generate deployment script
DEPLOY_SCRIPT="/tmp/e173_deploy_${COMMIT_ID}.sh"
cat > "$DEPLOY_SCRIPT" << 'EODEPLOY'
#!/bin/bash

set -e

COMMIT_ID="$1"
REPO_URL="$2"
ENVIRONMENT="$3"
PROJECT_NAME="e173_go_gateway"

echo "🏗️  Starting fresh deployment on $(hostname)..."

# Install dependencies
echo "📦 Installing dependencies..."
if command -v apt-get > /dev/null; then
    sudo apt-get update
    sudo apt-get install -y postgresql postgresql-contrib git golang-go make curl jq
elif command -v yum > /dev/null; then
    sudo yum update -y
    sudo yum install -y postgresql postgresql-server git golang make curl jq
fi

# Start PostgreSQL if not running
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create project directory
PROJECT_DIR="/opt/$PROJECT_NAME"
sudo mkdir -p "$PROJECT_DIR"
sudo chown $(whoami):$(whoami) "$PROJECT_DIR"

# Clone repository
if [ -d "$PROJECT_DIR/.git" ]; then
    echo "📂 Updating existing repository..."
    cd "$PROJECT_DIR"
    git fetch origin
    git checkout "$COMMIT_ID"
else
    echo "📥 Cloning repository..."
    git clone "$REPO_URL" "$PROJECT_DIR"
    cd "$PROJECT_DIR"
    git checkout "$COMMIT_ID"
fi

# Setup PostgreSQL database
echo "🗄️  Setting up database..."
sudo -u postgres psql << EOF
CREATE DATABASE e173_gateway;
CREATE USER e173_user WITH PASSWORD 'e173_pass';
GRANT ALL PRIVILEGES ON DATABASE e173_gateway TO e173_user;
ALTER USER e173_user CREATEDB;
\q
EOF

# Configure environment
echo "⚙️  Configuring environment..."
cat > .env << EOF
# E173 Gateway Configuration - $ENVIRONMENT
ENVIRONMENT=$ENVIRONMENT
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=e173_gateway
DB_USER=e173_user
DB_PASSWORD=e173_pass

# Asterisk AMI Configuration
AMI_HOST=localhost
AMI_PORT=5038
AMI_USER=admin
AMI_PASSWORD=secret
EOF

# Build application
echo "🔨 Building application..."
go mod tidy
make build

# Run database migrations
echo "📊 Running database migrations..."
make migrate-up

echo "✅ Fresh deployment base setup completed"
echo "🔄 Ready for snapshot restore..."

EODEPLOY

chmod +x "$DEPLOY_SCRIPT"

# Deploy to server
if [ "$SERVER_ADDRESS" = "localhost" ]; then
    # Local deployment
    echo "🏠 Deploying locally..."
    bash "$DEPLOY_SCRIPT" "$COMMIT_ID" "$(git remote get-url origin 2>/dev/null || echo 'local')" "$ENVIRONMENT"
    
    # Restore snapshot locally
    echo "📥 Restoring snapshot..."
    cd "$PROJECT_DIR"
    ./scripts/snapshot_restore.sh "$COMMIT_ID"
    
else
    # Remote deployment
    echo "🌐 Deploying to remote server: $SERVER_ADDRESS"
    
    # Copy deployment script to server
    scp "$DEPLOY_SCRIPT" "$SERVER_ADDRESS:/tmp/"
    
    # Copy snapshot backup to server
    scp "$BACKUP_FILE" "$SERVER_ADDRESS:/tmp/"
    
    # Execute deployment on remote server
    ssh "$SERVER_ADDRESS" bash /tmp/$(basename "$DEPLOY_SCRIPT") "$COMMIT_ID" "$(git remote get-url origin)" "$ENVIRONMENT"
    
    # Restore snapshot on remote server
    ssh "$SERVER_ADDRESS" << EOSSH
cd /opt/e173_go_gateway
mkdir -p snapshots/backups
cp /tmp/$(basename "$BACKUP_FILE") snapshots/backups/
./scripts/snapshot_restore.sh "$COMMIT_ID"
EOSSH
fi

# Cleanup
rm -f "$DEPLOY_SCRIPT"

# Verify deployment
echo "🔍 Verifying deployment..."
if [ "$SERVER_ADDRESS" = "localhost" ]; then
    HEALTH_URL="http://localhost:8080"
else
    HEALTH_URL="http://$SERVER_ADDRESS:8080"
fi

sleep 5
if curl -s "$HEALTH_URL/ping" > /dev/null 2>&1; then
    echo "✅ Deployment successful! Server is responding."
    echo "🌐 Access your application at: $HEALTH_URL"
else
    echo "⚠️  Deployment completed but server may need more time to start"
    echo "🔧 Check logs with: ssh $SERVER_ADDRESS 'cd /opt/e173_go_gateway && make logs'"
fi

echo ""
echo "🎉 Fresh deployment completed!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Snapshot:     $COMMIT_ID"
echo "🖥️  Server:       $SERVER_ADDRESS"
echo "🌐 URL:          $HEALTH_URL"
echo "🏷️  Environment: $ENVIRONMENT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
