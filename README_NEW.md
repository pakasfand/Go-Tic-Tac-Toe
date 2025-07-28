# Go Tic-Tac-Toe

A multiplayer Tic-Tac-Toe game built with Go and WebAssembly, featuring real-time gameplay through WebSockets.

## Development Commands

### Build WebAssembly Client
To update main.wasm run the following in PowerShell:
```powershell
$env:GOOS="js"
$env:GOARCH="wasm"
cd client
go build -o main.wasm .
```

### Run Locally
```bash
# Start the server
cd server
go run .

# Open client/index.html in browser (serve from local web server)
```

## Docker Commands

### Build and Push Docker Image
```bash
# Build new docker image
docker build --no-cache -t pakasfand/go-tic-tac-toe:latest .

# Push to Docker Hub
docker push pakasfand/go-tic-tac-toe:latest
```

### Run Locally with Docker
```bash
# Build and run locally
docker build -t go-tic-tac-toe:local .
docker run -p 80:80 -p 8080:8080 go-tic-tac-toe:local

# Or use docker-compose
docker-compose up --build
```

## AWS Lightsail Deployment Commands

### Initial Setup (on Lightsail instance)
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Configure firewall
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw --force enable
```

### Deploy/Update Application
```bash
# Pull latest docker image
sudo docker pull pakasfand/go-tic-tac-toe:latest

# If using single container:
sudo docker stop go-tic-tac-toe 2>/dev/null || true
sudo docker rm go-tic-tac-toe 2>/dev/null || true
sudo docker run -d --name go-tic-tac-toe -p 80:80 -p 8080:8080 pakasfand/go-tic-tac-toe:latest

# If using docker-compose:
cd /opt/go-tic-tac-toe
sudo docker-compose pull
sudo docker-compose up -d
```

### Management Commands
```bash
# Check status
sudo docker ps
sudo docker logs go-tic-tac-toe

# Restart
sudo docker restart go-tic-tac-toe

# Stop/remove
sudo docker stop go-tic-tac-toe
sudo docker rm go-tic-tac-toe
```

## Complete Deployment Workflow

### 1. Local Development & Testing
```bash
# Make code changes
# Test locally

# Build and test with Docker
docker build -t go-tic-tac-toe:test .
docker run -p 80:80 -p 8080:8080 go-tic-tac-toe:test
```

### 2. Build & Push to Registry
```bash
docker build --no-cache -t pakasfand/go-tic-tac-toe:latest .
docker push pakasfand/go-tic-tac-toe:latest
```

### 3. Deploy to Lightsail
```bash
# SSH into Lightsail instance
ssh ubuntu@YOUR_LIGHTSAIL_IP

# Pull and deploy
sudo docker pull pakasfand/go-tic-tac-toe:latest
sudo docker stop go-tic-tac-toe 2>/dev/null || true
sudo docker rm go-tic-tac-toe 2>/dev/null || true
sudo docker run -d --name go-tic-tac-toe -p 80:80 -p 8080:8080 pakasfand/go-tic-tac-toe:latest
```

## Troubleshooting Commands

### Debug WebSocket Issues
```bash
# Check if WebSocket endpoint is accessible
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Host: YOUR_DOMAIN" -H "Origin: http://YOUR_DOMAIN" http://YOUR_DOMAIN/ws

# Check nginx proxy
sudo docker exec go-tic-tac-toe cat /etc/nginx/nginx.conf

# Check container logs
sudo docker logs go-tic-tac-toe -f
```

### Performance Monitoring
```bash
# Check container resource usage
sudo docker stats go-tic-tac-toe

# Check disk usage
sudo docker system df

# Clean up unused images/containers
sudo docker system prune -f
```

## Architecture

- **Client**: WebAssembly (Go + Ebiten) served as static files
- **Server**: Go WebSocket server for real-time communication
- **Proxy**: Nginx for static file serving and WebSocket proxying
- **Deployment**: Docker container on AWS Lightsail
