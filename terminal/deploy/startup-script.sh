#!/bin/bash
# AWS User Data Script for Terminal Portfolio
# Works on Amazon Linux 2023, Ubuntu, and Debian

set -e

# Redirect output for debugging (view at /var/log/user-data.log)
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

echo "=== Setting up Terminal Portfolio Server ==="

# Detect Package Manager
if command -v dnf &> /dev/null; then
    PKG_MGR="dnf"   # Amazon Linux 2023 / Fedora
elif command -v yum &> /dev/null; then
    PKG_MGR="yum"   # Amazon Linux 2 / RHEL
elif command -v apt-get &> /dev/null; then
    PKG_MGR="apt-get"  # Ubuntu / Debian
else
    echo "Unsupported OS"
    exit 1
fi

echo "Using Package Manager: $PKG_MGR"

# Install dependencies
if [ "$PKG_MGR" = "apt-get" ]; then
    $PKG_MGR update -qq
    $PKG_MGR install -y -qq curl git
else
    $PKG_MGR update -y
    $PKG_MGR install -y curl git
fi

# Create application directory
mkdir -p /opt/terminal-portfolio
mkdir -p /opt/terminal-portfolio/.ssh

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
RESEND_API_KEY=your-resend-key-here
EOF
fi

# Configure System SSH on Port 2222 (Critical step)
SSHD_CONFIG="/etc/ssh/sshd_config"
if ! grep -q "Port 2222" $SSHD_CONFIG; then
    echo "Configuring system SSH to listen on Port 2222..."
    
    # Backup config
    cp $SSHD_CONFIG "$SSHD_CONFIG.bak"
    
    # Add Port 2222 configuration if not present
    echo "Port 2222" >> $SSHD_CONFIG
    
    # Ensure policy allows binding to new port (SELinux)
    if command -v semanage &> /dev/null; then
        semanage port -a -t ssh_port_t -p tcp 2222 || true
    fi
    
    # Restart SSH service
    systemctl restart sshd || systemctl restart ssh
    echo "System SSH moved to Port 2222. Please reconnect on new port."
fi

# Copy systemd service file (assumes it was uploaded or created)
if [ -f /tmp/portfolio.service ]; then
    cp /tmp/portfolio.service /etc/systemd/system/portfolio.service
    systemctl daemon-reload
    systemctl enable portfolio.service
fi

echo "=== Setup Complete ==="
echo "WARNING: System SSH is now on Port 2222."
echo "Application will listen on Port 22 (Standard SSH)."
