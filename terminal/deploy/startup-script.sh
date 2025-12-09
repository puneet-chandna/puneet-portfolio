#!/bin/bash
# GCP VM Startup Script for Terminal Portfolio
# This script runs on first boot or can be run manually to set up the server

set -e

echo "=== Setting up Terminal Portfolio Server ==="

# Create application directory
mkdir -p /opt/terminal-portfolio
mkdir -p /opt/terminal-portfolio/.ssh

# Install dependencies (if needed)
apt-get update -qq
apt-get install -y -qq curl

# Generate SSH host key if it doesn't exist
if [ ! -f /opt/terminal-portfolio/.ssh/id_ed25519 ]; then
    echo "Generating SSH host key..."
    ssh-keygen -t ed25519 -f /opt/terminal-portfolio/.ssh/id_ed25519 -N "" -C "portfolio-host-key"
fi

# Create environment file template if it doesn't exist
if [ ! -f /opt/terminal-portfolio/.env ]; then
    echo "Creating environment file template..."
    cat > /opt/terminal-portfolio/.env << 'EOF'
# Terminal Portfolio Environment Variables
SSH_PORT=22
WEB3FORMS_KEY=your-web3forms-key-here
EOF
    echo "WARNING: Update /opt/terminal-portfolio/.env with your WEB3FORMS_KEY"
fi

# Copy systemd service file
cp /tmp/portfolio.service /etc/systemd/system/portfolio.service 2>/dev/null || true

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable portfolio.service

echo "=== Setup Complete ==="
echo "Next steps:"
echo "1. Upload the portfolio binary to /opt/terminal-portfolio/"
echo "2. Update /opt/terminal-portfolio/.env with your WEB3FORMS_KEY"
echo "3. Run: systemctl start portfolio"
echo "4. Check status: systemctl status portfolio"
