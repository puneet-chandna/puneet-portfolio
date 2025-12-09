# GCP VM Setup Guide for Terminal Portfolio

## Quick Start

### 1. Create GCP VM (Free Tier)

```bash
# Using gcloud CLI (or use GCP Console)
gcloud compute instances create terminal-portfolio \
  --machine-type=e2-micro \
  --zone=us-central1-a \
  --image-family=debian-12 \
  --image-project=debian-cloud \
  --boot-disk-size=10GB \
  --tags=ssh-server

# Create firewall rule for SSH on port 22
gcloud compute firewall-rules create allow-ssh-portfolio \
  --allow=tcp:22 \
  --target-tags=ssh-server \
  --source-ranges=0.0.0.0/0

# Reserve static IP
gcloud compute addresses create portfolio-ip --region=us-central1
gcloud compute instances add-access-config terminal-portfolio \
  --zone=us-central1-a \
  --address=$(gcloud compute addresses describe portfolio-ip --region=us-central1 --format='get(address)')
```

### 2. Initial VM Setup

SSH into your VM and run:

```bash
# Download and run startup script
curl -sL https://raw.githubusercontent.com/puneet-chandna/puneet-portfolio/main/terminal/deploy/startup-script.sh | sudo bash
```

### 3. Upload Binary (Manual - First Time)

From your local machine:

```bash
# Build the binary
cd terminal
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o portfolio ./cmd/portfolio

# Upload to VM
gcloud compute scp portfolio YOUR_VM_USER@terminal-portfolio:/tmp/
gcloud compute ssh terminal-portfolio --command="sudo mv /tmp/portfolio /opt/terminal-portfolio/ && sudo chmod +x /opt/terminal-portfolio/portfolio"
```

### 4. Configure Environment

On the VM:

```bash
# Edit environment file
sudo nano /opt/terminal-portfolio/.env

# Add your WEB3FORMS_KEY:
# WEB3FORMS_KEY=your-actual-key-here
# SSH_PORT=22
```

### 5. Start the Service

```bash
sudo systemctl start portfolio
sudo systemctl status portfolio
```

### 6. Configure DNS

Add an A record in your domain registrar:

- **Host**: `@` (or `ssh` for subdomain)
- **Value**: Your GCP VM's external IP
- **TTL**: 300

### 7. Test

```bash
ssh puneetchandna.com
# or
ssh ssh.puneetchandna.com
```

---

## GitHub Secrets Required

Add these secrets to your GitHub repository (Settings > Secrets > Actions):

| Secret Name           | Description                      |
| --------------------- | -------------------------------- |
| `GCP_SSH_PRIVATE_KEY` | Private SSH key to connect to VM |
| `GCP_VM_IP`           | External IP of your GCP VM       |
| `GCP_VM_USER`         | Your username on the VM          |
| `WEB3FORMS_KEY`       | (Optional) For contact form      |

### Generate SSH Key for GitHub Actions

```bash
# On your local machine
ssh-keygen -t ed25519 -f ~/.ssh/gcp-deploy -N "" -C "github-actions-deploy"

# Copy public key to VM
gcloud compute ssh terminal-portfolio --command="mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys" < ~/.ssh/gcp-deploy.pub

# Add private key content to GitHub Secrets as GCP_SSH_PRIVATE_KEY
cat ~/.ssh/gcp-deploy
```

---

## Useful Commands

```bash
# Check service status
sudo systemctl status portfolio

# View logs
sudo journalctl -u portfolio -f

# Restart service
sudo systemctl restart portfolio

# Check if port 22 is listening
sudo ss -tulpn | grep :22
```
