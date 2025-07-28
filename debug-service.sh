#!/bin/bash

echo "================================================"
echo "  Debugging tic-tac-toe.service"
echo "================================================"

echo "1. Checking systemd service status..."
systemctl status tic-tac-toe.service --no-pager

echo ""
echo "2. Checking systemd service logs..."
journalctl -u tic-tac-toe.service --no-pager -n 20

echo ""
echo "3. Checking if docker-compose exists and is executable..."
ls -la /usr/local/bin/docker-compose
ls -la /usr/bin/docker-compose

echo ""
echo "4. Testing docker-compose manually..."
cd /opt/go-tic-tac-toe
which docker-compose
docker-compose --version

echo ""
echo "5. Checking working directory and files..."
pwd
ls -la

echo ""
echo "6. Testing docker-compose up manually..."
docker-compose up -d

echo ""
echo "7. Checking container status..."
docker ps

echo ""
echo "8. Checking if ports are available..."
netstat -tlnp | grep -E ':(80|8080|443)'

echo ""
echo "================================================"
echo "  Debug Complete"
echo "================================================"
