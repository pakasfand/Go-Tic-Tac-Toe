#!/bin/bash

echo "================================================"
echo "  Setting up Go Tic-Tac-Toe on AWS Lightsail"
echo "  (Using pre-built Docker image)"
echo "================================================"

# Exit on any error
set -e

# Update system
echo "Updating system packages..."
apt-get update -y
apt-get upgrade -y

# Install required packages
echo "Installing required packages..."
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    wget \
    unzip

# Install Docker
echo "Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get update -y
    apt-get install -y docker-ce docker-ce-cli containerd.io
    systemctl enable docker
    systemctl start docker
    usermod -aG docker ubuntu
    echo "Docker installed successfully"
else
    echo "Docker is already installed"
fi

# Install Docker Compose
echo "Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose installed successfully"
else
    echo "Docker Compose is already installed"
fi

# Create application directory
echo "Setting up application directory..."
mkdir -p /opt/go-tic-tac-toe
cp -r ./* /opt/go-tic-tac-toe/
cd /opt/go-tic-tac-toe

# Load Docker image from file (if exists)
if [ -f "go-tic-tac-toe-image.tar.gz" ]; then
    echo "Loading Docker image from file..."
    gunzip go-tic-tac-toe-image.tar.gz
    docker load -i go-tic-tac-toe-image.tar
    docker tag go-tic-tac-toe:test go-tic-tac-toe:latest
    echo "Docker image loaded successfully"
elif command -v docker &> /dev/null; then
    echo "Pulling Docker image from Docker Hub..."
    docker pull pakasfand/go-tic-tac-toe:latest || {
        echo "Failed to pull from Docker Hub, building locally..."
        docker build -t go-tic-tac-toe:latest .
    }
    docker tag pakasfand/go-tic-tac-toe:latest go-tic-tac-toe:latest 2>/dev/null || true
else
    echo "Building Docker image locally..."
    docker build -t go-tic-tac-toe:latest .
fi

# Install systemd service
echo "Installing systemd service..."
cp tic-tac-toe.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable tic-tac-toe.service

# Configure firewall
echo "Configuring firewall..."
ufw allow 22    # SSH
ufw allow 80    # HTTP
ufw allow 443   # HTTPS
ufw --force enable

# Start the application
echo "Starting application..."
systemctl start tic-tac-toe.service

# Wait for services to start
echo "Waiting for services to start..."
sleep 30

# Check if services are running
echo "Checking service status..."
systemctl status tic-tac-toe.service
docker-compose ps

# Get public IP
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)

echo "================================================"
echo "  Deployment Complete!"
echo "================================================"
echo "Your Go Tic-Tac-Toe game is now running at:"
echo "  http://$PUBLIC_IP"
echo ""
echo "Service management commands:"
echo "  Start:   sudo systemctl start tic-tac-toe.service"
echo "  Stop:    sudo systemctl stop tic-tac-toe.service"
echo "  Restart: sudo systemctl restart tic-tac-toe.service"
echo "  Status:  sudo systemctl status tic-tac-toe.service"
echo "  Logs:    sudo journalctl -u tic-tac-toe.service -f"
echo ""
echo "Docker commands:"
echo "  View logs:     sudo docker-compose logs -f"
echo "  Restart app:   sudo docker-compose restart"
echo "  Stop app:      sudo docker-compose down"
echo "  Start app:     sudo docker-compose up -d"
echo ""
echo "Nginx configuration: /opt/go-tic-tac-toe/nginx.conf"
echo "Application files:   /opt/go-tic-tac-toe/"
echo "================================================"
