#!/bin/bash

# E173 Go Gateway - Database Setup Script
# This script creates the PostgreSQL database and user for the E173 Gateway

set -e

echo "🔧 Setting up E173 Gateway Database..."

# Database configuration
DB_NAME="e173_gateway"
DB_USER="e173_user"
DB_PASS="e173_pass"

# Check if PostgreSQL is running
if ! systemctl is-active --quiet postgresql; then
    echo "⚠️ PostgreSQL is not running. Starting PostgreSQL..."
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
fi

# Create database and user
echo "📊 Creating database and user..."
sudo -u postgres psql << EOF
-- Create user if it doesn't exist
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN
        CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASS';
    END IF;
END
\$\$;

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE $DB_NAME OWNER $DB_USER'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
ALTER USER $DB_USER CREATEDB;
EOF

echo "✅ Database '$DB_NAME' and user '$DB_USER' created successfully!"

# Test connection
echo "🔍 Testing database connection..."
if PGPASSWORD="$DB_PASS" psql -h localhost -U "$DB_USER" -d "$DB_NAME" -c "SELECT version();" > /dev/null 2>&1; then
    echo "✅ Database connection test successful!"
else
    echo "❌ Database connection test failed!"
    exit 1
fi

echo "🎉 Database setup complete!"
echo ""
echo "Your database configuration:"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASS"
echo "  Host: localhost"
echo "  Port: 5432"
echo ""
echo "💡 Next steps:"
echo "1. Copy .env.example to .env"
echo "2. Update .env with your configuration if needed"
echo "3. Run database migrations: make migrate"
echo "4. Start the server: make run"
