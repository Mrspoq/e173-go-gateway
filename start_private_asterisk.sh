#!/bin/bash

# Start the private Asterisk instance for E173 Gateway

ASTERISK_BASE="/root/e173_go_gateway/asterisk_private"

echo "Starting private Asterisk instance..."

# Check if already running
if pgrep -f "$ASTERISK_BASE/sbin/asterisk" > /dev/null; then
    echo "Private Asterisk is already running"
    exit 0
fi

# Start Asterisk
$ASTERISK_BASE/sbin/asterisk -C $ASTERISK_BASE/etc/asterisk/asterisk.conf

# Wait a moment for it to start
sleep 2

# Check if it started successfully
if pgrep -f "$ASTERISK_BASE/sbin/asterisk" > /dev/null; then
    echo "Private Asterisk started successfully"
    
    # Show dongle status
    echo "Checking dongle devices..."
    $ASTERISK_BASE/sbin/asterisk -rx "dongle show devices"
else
    echo "Failed to start private Asterisk"
    exit 1
fi