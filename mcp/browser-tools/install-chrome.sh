#!/bin/bash

echo "Installing Google Chrome for Browser Tools MCP..."

# Download and add Google's signing key
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add -

# Add Chrome repository
echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" | sudo tee /etc/apt/sources.list.d/google-chrome.list

# Update package list
sudo apt update

# Install Google Chrome
sudo apt install -y google-chrome-stable

# Verify installation
if which google-chrome > /dev/null; then
    echo "✅ Google Chrome installed successfully"
    google-chrome --version
else
    echo "❌ Chrome installation failed"
    exit 1
fi