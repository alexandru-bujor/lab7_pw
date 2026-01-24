#!/bin/bash

# MegaMobile Backend - Docker Deployment Script
# This script helps build and deploy the backend using Docker

set -e

echo "🚀 MegaMobile Backend - Docker Deployment"
echo "=========================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Function to build and start
build_and_start() {
    echo -e "${YELLOW}📦 Building Docker image...${NC}"
    docker-compose build
    
    echo -e "${YELLOW}🚀 Starting containers...${NC}"
    docker-compose up -d
    
    echo -e "${GREEN}✅ Backend is starting...${NC}"
    echo "📋 View logs with: docker-compose logs -f"
    echo "🔍 Check status with: docker-compose ps"
}

# Function to stop
stop() {
    echo -e "${YELLOW}🛑 Stopping containers...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ Containers stopped${NC}"
}

# Function to restart
restart() {
    echo -e "${YELLOW}🔄 Restarting containers...${NC}"
    docker-compose restart
    echo -e "${GREEN}✅ Containers restarted${NC}"
}

# Function to view logs
logs() {
    docker-compose logs -f
}

# Function to rebuild
rebuild() {
    echo -e "${YELLOW}🔨 Rebuilding and restarting...${NC}"
    docker-compose up -d --build
    echo -e "${GREEN}✅ Rebuild complete${NC}"
}

# Function to show status
status() {
    echo -e "${YELLOW}📊 Container Status:${NC}"
    docker-compose ps
    echo ""
    echo -e "${YELLOW}📋 Recent Logs:${NC}"
    docker-compose logs --tail=20
}

# Main menu
case "$1" in
    start|up)
        build_and_start
        ;;
    stop|down)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs
        ;;
    rebuild)
        rebuild
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|logs|rebuild|status}"
        echo ""
        echo "Commands:"
        echo "  start, up    - Build and start containers"
        echo "  stop, down   - Stop containers"
        echo "  restart      - Restart containers"
        echo "  logs         - View logs (follow mode)"
        echo "  rebuild      - Rebuild and restart containers"
        echo "  status       - Show container status and recent logs"
        exit 1
        ;;
esac


