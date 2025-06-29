#!/bin/bash

# Setup environment variables for E173 Gateway
# This script helps create a .env file with secure credentials

echo "E173 Gateway Environment Setup"
echo "=============================="

# Check if .env exists
if [ -f ".env" ]; then
    echo "Warning: .env file already exists!"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborting..."
        exit 1
    fi
fi

# Copy from example
cp .env.example .env

# Get database password
read -sp "Enter database password (or press Enter for default): " db_pass
echo
if [ -z "$db_pass" ]; then
    db_pass="3omartel580"
fi

# Get AMI password
read -sp "Enter Asterisk AMI password (or press Enter for default): " ami_pass
echo
if [ -z "$ami_pass" ]; then
    ami_pass="3omartel580"
fi

# Get admin password
read -sp "Enter admin password (or press Enter for default): " admin_pass
echo
if [ -z "$admin_pass" ]; then
    admin_pass="admin"
fi

# Get JWT secret
jwt_secret=$(openssl rand -hex 32 2>/dev/null || echo "e173-gateway-secret-key-change-in-production")

# Update .env file
sed -i "s/YOUR_DB_PASSWORD_HERE/$db_pass/g" .env
sed -i "s/YOUR_AMI_PASSWORD_HERE/$ami_pass/g" .env
sed -i "s/YOUR_ADMIN_PASSWORD_HERE/$admin_pass/g" .env
sed -i "s/YOUR_JWT_SECRET_HERE/$jwt_secret/g" .env
sed -i "s/DB_USER:DB_PASSWORD/e173_user:$db_pass/g" .env

echo ""
echo "Environment file created successfully!"
echo "Remember to:"
echo "1. Add your GitHub token if using MCP"
echo "2. Update Redis password if needed"
echo "3. Change to production settings before deployment"