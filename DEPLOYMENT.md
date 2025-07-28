# Go Tic-Tac-Toe AWS Lightsail Deployment

This document provides step-by-step instructions for deploying your Go Tic-Tac-Toe multiplayer game to AWS Lightsail using Docker Hub.

## Prerequisites

- AWS Account with Lightsail access
- Docker Desktop installed on your local machine
- Docker Hub account (for image registry)

## Quick Deployment

### Step 1: Build and Push Docker Image

1. Build and push the Docker image to Docker Hub:
   ```cmd
   docker build --no-cache -t pakasfand/go-tic-tac-toe:latest .
   docker push pakasfand/go-tic-tac-toe:latest
   ```

This will:
- Build the WebAssembly client inside Docker
- Build the Go server
- Create optimized Docker image
- Push to Docker Hub registry

### Step 2: Create Lightsail Instance

1. Go to [AWS Lightsail Console](https://lightsail.aws.amazon.com/)
2. Click "Create instance"
3. Choose:
   - **Platform**: Linux/Unix
   - **Blueprint**: Ubuntu 20.04 LTS or Ubuntu 22.04 LTS
   - **Instance plan**: At least $5/month (1GB RAM, 1 vCPU)
   - **Instance name**: `go-tic-tac-toe-server`
4. Click "Create instance"
5. Wait for the instance to be running

### Step 3: Deploy to Lightsail

1. Connect to your instance via SSH (use the Lightsail browser-based SSH)

2. Install Docker and Docker Compose:
   ```bash
   # Update system
   sudo apt-get update -y
   
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker ubuntu
   
   # Install Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   
   # Logout and login again for group changes to take effect
   exit
   ```

3. Reconnect and deploy:
   ```bash
   # Create application directory
   sudo mkdir -p /opt/go-tic-tac-toe
   cd /opt/go-tic-tac-toe
   
   # Create docker-compose.yml (copy content from your local file)
   sudo nano docker-compose.yml
   
   # Create nginx.conf (copy content from your local file)  
   sudo nano nginx.conf
   
   # Pull and start containers
   sudo docker-compose pull
   sudo docker-compose up -d
   
   # Configure firewall
   sudo ufw allow 22
   sudo ufw allow 80
   sudo ufw allow 443
   sudo ufw --force enable
   ```

### Step 4: Configure Domain (Optional)

1. In Lightsail, go to your instance
2. Click "Networking" tab
3. Create a static IP address
4. Attach it to your instance
5. If you have a domain, create DNS records pointing to your static IP

For SSL/HTTPS setup:
```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Get SSL certificate (replace with your domain)
sudo certbot --nginx -d yourdomain.com

# Update nginx.conf to use HTTPS
sudo nano /opt/go-tic-tac-toe/nginx.conf
# Uncomment the HTTPS server block and update the domain name

# Restart services
sudo systemctl restart tic-tac-toe.service
```

## Manual Deployment (Alternative)

If you prefer to deploy manually without Docker Compose:

### 1. Prepare Your Lightsail Instance

```bash
# Update system
sudo apt-get update -y
sudo apt-get upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu

# Logout and login again for group changes to take effect
```

### 2. Deploy Your Application

```bash
# Pull the Docker image
sudo docker pull pakasfand/go-tic-tac-toe:latest

# Run the container
sudo docker run -d --name go-tic-tac-toe -p 80:80 -p 8080:8080 pakasfand/go-tic-tac-toe:latest

# Enable firewall
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw --force enable
```

## Service Management

### View Application Status
```bash
sudo docker ps
sudo docker-compose ps  # if using docker-compose
```

### View Logs
```bash
# Container logs
sudo docker logs go-tic-tac-toe

# Docker compose logs
cd /opt/go-tic-tac-toe
sudo docker-compose logs -f
```

### Restart Application
```bash
# Single container
sudo docker restart go-tic-tac-toe

# Docker compose
cd /opt/go-tic-tac-toe
sudo docker-compose restart
```

### Update Application
```bash
# Pull latest image
sudo docker pull pakasfand/go-tic-tac-toe:latest

# For single container deployment:
sudo docker stop go-tic-tac-toe
sudo docker rm go-tic-tac-toe
sudo docker run -d --name go-tic-tac-toe -p 80:80 -p 8080:8080 pakasfand/go-tic-tac-toe:latest

# For docker-compose deployment:
cd /opt/go-tic-tac-toe
sudo docker-compose pull
sudo docker-compose up -d
```

## Monitoring and Maintenance

### Health Checks
The Docker containers include health checks. Monitor them with:
```bash
sudo docker ps
sudo docker-compose ps
```

### Backup
```bash
# Backup application files
sudo tar -czf tic-tac-toe-backup-$(date +%Y%m%d).tar.gz /opt/go-tic-tac-toe/

# Store backup in a safe location
```

### Updates
```bash
# System updates
sudo apt-get update && sudo apt-get upgrade -y

# Docker updates
sudo docker-compose pull
sudo docker-compose up -d
```

## Troubleshooting

### Application Won't Start
```bash
# Check Docker status
sudo systemctl status docker

# Check container logs
sudo docker logs go-tic-tac-toe

# Restart Docker service
sudo systemctl restart docker
```

### Can't Connect to Game
1. Check if ports 80 and 443 are open in Lightsail firewall
2. Verify the container is running: `sudo docker ps`
3. Check nginx logs: `sudo docker logs go-tic-tac-toe`
4. Verify the public IP address in Lightsail console

### WebSocket Connection Issues
1. Ensure nginx is properly proxying WebSocket connections
2. Check the `/ws` endpoint specifically  
3. Verify the WebAssembly client is using dynamic host detection
4. Clear browser cache or try incognito mode to ensure latest WASM is loaded

## Security Considerations

1. **Firewall**: Only ports 22 (SSH), 80 (HTTP), and 443 (HTTPS) should be open
2. **SSL**: Use Let's Encrypt for free SSL certificates
3. **Updates**: Regularly update the system and Docker images
4. **Monitoring**: Set up CloudWatch or other monitoring for production use
5. **Backups**: Regular backups of application data and configuration

## Cost Optimization

- **Instance Size**: Start with the $5/month plan (1GB RAM)
- **Static IP**: Free with Lightsail instance
- **Data Transfer**: 1TB included with $5 plan
- **Monitoring**: Use Lightsail metrics (included)

For higher traffic, consider upgrading to larger instance sizes or migrating to EC2 with Auto Scaling.
