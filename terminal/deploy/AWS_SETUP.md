# AWS EC2 Setup Guide for Terminal Portfolio

## Quick Start

### 1. Launch EC2 Instance (Free Tier)

1.  **Go to AWS Console** > EC2 > Launch Instance.
2.  **Name**: `terminal-portfolio`.
3.  **OS Image**: **Amazon Linux 2023 AMI** (Free Tier eligible).
4.  **Instance Type**: `t2.micro` or `t3.micro` (Free Tier eligible).
5.  **Key Pair**: Create a new key pair or select an existing one. **Save the `.pem` file!**
6.  **Network Settings**:
    - Create Security Group allowing:
      - **SSH (TCP 22)**: Anywhere (`0.0.0.0/0`) - _For Portfolio User Access_
      - **Custom TCP 2222**: Anywhere (`0.0.0.0/0`) - _For Admin SSH Access_
7.  **Advanced Details** > **User Data**:
    Paste the content of `startup-script.sh` here. This will automatically configure the server on boot.

---

### 2. Accessing Your Instance

After launch, the **system SSH** will verify keys and move to **Port 2222** (configured by the startup script).

**Admin SSH Access:**

```bash
ssh -p 2222 -i /path/to/key.pem ec2-user@YOUR_INSTANCE_IP
```

> **Note:** If the startup script hasn't finished, it might still remain on port 22 temporarily. Wait 1-2 minutes.

---

### 3. Configure Environment

SSH into your VM and edit the environment file:

```bash
sudo nano /opt/terminal-portfolio/.env
```

Update your API keys:

```ini
SSH_PORT=22
WEB3FORMS_KEY=your-actual-key
RESEND_API_KEY=your-actual-key
```

---

### 4. Configure GitHub Secrets for CI/CD

In your GitHub Repository > Settings > Secrets > Actions:

| Secret Name           | Value                                             |
| --------------------- | ------------------------------------------------- |
| `EC2_HOST`            | Public IP or DNS of your EC2 instance             |
| `EC2_USER`            | `ec2-user` (for Amazon Linux)                     |
| `EC2_SSH_PRIVATE_KEY` | Content of your private SSH key (generated below) |

#### Generate Deployment Key

To allow GitHub Actions to deploy, generate a new key pair:

```bash
# Local machine
ssh-keygen -t ed25519 -f terminal-deploy -C "github-actions"
```

1.  Copy content of `terminal-deploy` (private key) to **GitHub Secret `EC2_SSH_PRIVATE_KEY`**.
2.  Add content of `terminal-deploy.pub` to the server:

```bash
# On EC2
echo "YOUR_PUBLIC_KEY_CONTENT" >> ~/.ssh/authorized_keys
```

---

### 5. DNS Configuration

In your domain registrar (e.g., Namecheap, GoDaddy, Route53, Hostinger):

- **Type**: `A Record`
- **Host**: `@` (or `ssh` for subdomain)
- **Value**: Your EC2 Public IP
- **TTL**: `300` (5 mins)

**Test:**

```bash
ssh yourdomain.com
```

---

## Troubleshooting

**Service Status:**

```bash
sudo systemctl status portfolio
```

**View Logs:**

```bash
sudo journalctl -u portfolio -f
```

**Check Ports:**

```bash
sudo ss -tulpn | grep -E ":22|:2222"
```
