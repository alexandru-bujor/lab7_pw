# Production Deployment Guide - MegaMobile Backend

## Overview

This guide covers deploying the MegaMobile backend to a Linux VPS with:
- MySQL as a system service (not Docker)
- Backend API as Docker container
- Automatic migrations on startup
- CORS configured for production frontend

## Prerequisites

- Linux VPS (Ubuntu/Debian recommended)
- Docker and Docker Compose installed
- Root or sudo access
- Domain name configured (megamobile.md)

## Step 1: Install MySQL Server

### 1.1 Install MySQL

```bash
# Update package list
sudo apt update

# Install MySQL server
sudo apt install mysql-server -y

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Check status
sudo systemctl status mysql
```

### 1.2 Secure MySQL Installation

```bash
# Run MySQL secure installation
sudo mysql_secure_installation
```

Follow the prompts:
- Set root password (or leave empty for socket authentication)
- Remove anonymous users: **Yes**
- Disallow root login remotely: **Yes** (we'll create a separate user)
- Remove test database: **Yes**
- Reload privilege tables: **Yes**

### 1.3 Create Database and User

```bash
# Login to MySQL as root
sudo mysql -u root -p
```

Then run these SQL commands:

```sql
-- Create database
CREATE DATABASE IF NOT EXISTS mega_mobile CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create dedicated user for the application
CREATE USER 'megamobile_user'@'localhost' IDENTIFIED BY 'your_secure_password_here';

-- Grant privileges on the database
GRANT ALL PRIVILEGES ON mega_mobile.* TO 'megamobile_user'@'localhost';

-- Allow connection from Docker containers (if needed)
CREATE USER 'megamobile_user'@'%' IDENTIFIED BY 'your_secure_password_here';
GRANT ALL PRIVILEGES ON mega_mobile.* TO 'megamobile_user'@'%';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Verify user was created
SELECT user, host FROM mysql.user WHERE user = 'megamobile_user';

-- Exit MySQL
EXIT;
```

**Important:** Replace `'your_secure_password_here'` with a strong password!

### 1.4 Configure MySQL to Accept Docker Connections

Edit MySQL configuration:

```bash
sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

Find the `bind-address` line and change it:

```ini
# Change from:
bind-address = 127.0.0.1

# To:
bind-address = 0.0.0.0
```

**Security Note:** This allows connections from any IP. For better security, you can:
- Keep `bind-address = 127.0.0.1` and use `host.docker.internal` (works on Docker Desktop)
- Or configure firewall to only allow local connections

Restart MySQL:

```bash
sudo systemctl restart mysql
```

### 1.5 Test MySQL Connection

```bash
# Test local connection
mysql -u megamobile_user -p -h localhost mega_mobile

# If successful, you should see MySQL prompt
# Type EXIT; to leave
```

## Step 2: Prepare Backend Deployment

### 2.1 Upload Code to VPS

```bash
# On your local machine, create a deployment package
cd c:\dev\MegaMobile\MegaMobileBack
tar -czf megamobile-backend.tar.gz \
  --exclude='.git' \
  --exclude='node_modules' \
  --exclude='tmp' \
  --exclude='build' \
  --exclude='.env' \
  --exclude='.env.local' \
  .

# Upload to VPS
scp megamobile-backend.tar.gz root@164.68.98.158:/opt/

# SSH into VPS
ssh root@164.68.98.158
```

### 2.2 Extract and Setup on VPS

```bash
# Create application directory
mkdir -p /opt/megamobile-backend
cd /opt/megamobile-backend

# Extract files
tar -xzf /opt/megamobile-backend.tar.gz

# Create uploads directory
mkdir -p uploads/products
chmod -R 755 uploads
```

### 2.3 Create Production Environment File

```bash
nano .env.production
```

Add the following (replace with your actual values):

```env
# Database Configuration
DB_USER=megamobile_user
DB_PASS=your_secure_password_here
DB_HOST=host.docker.internal
DB_PORT=3306
DB_NAME=mega_mobile

# Server Configuration
PORT=8080
GIN_MODE=release

# Telegram Bot (optional)
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here

# CORS Configuration
CORS_ORIGINS=https://megamobile.md,https://www.megamobile.md
```

**Important:** 
- Replace `your_secure_password_here` with the password you set for `megamobile_user`
- If `host.docker.internal` doesn't work, try `172.17.0.1` (Docker bridge IP) or the host's actual IP

### 2.4 Test Database Connection

```bash
# Test from host (should work)
mysql -u megamobile_user -p -h localhost mega_mobile -e "SELECT 1;"

# If using host.docker.internal, test from a temporary container
docker run --rm --add-host=host.docker.internal:host-gateway \
  mysql:8.0 mysql -h host.docker.internal -u megamobile_user -p mega_mobile -e "SELECT 1;"
```

## Step 3: Deploy with Docker

### 3.1 Build and Start Container

```bash
cd /opt/megamobile-backend

# Build and start
docker-compose -f docker-compose.prod.yml up -d --build

# Check logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 3.2 Verify Deployment

```bash
# Check container status
docker ps

# Check health endpoint
curl http://localhost:8080/api/health

# Check logs for migration messages
docker logs megamobile-api-prod
```

You should see messages like:
```
✅ products tables ready
✅ categories table ready
✅ services table ready
...
🚀 Server running on http://localhost:8080
```

## Step 4: Configure Nginx Reverse Proxy (Optional but Recommended)

### 4.1 Install Nginx

```bash
sudo apt install nginx -y
```

### 4.2 Create Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/megamobile-api
```

Add:

```nginx
server {
    listen 80;
    server_name api.megamobile.md;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 4.3 Enable Site

```bash
sudo ln -s /etc/nginx/sites-available/megamobile-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 4.4 Setup SSL with Let's Encrypt

```bash
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d api.megamobile.md
```

## Step 5: Run Manual Migrations (if needed)

The application runs automatic migrations on startup. However, if you have SQL migration files:

```bash
# Connect to MySQL
mysql -u megamobile_user -p mega_mobile

# Run migration files
source /opt/megamobile-backend/migrations/20251229_create_service_categories.sql;
# ... other migrations

EXIT;
```

Or from command line:

```bash
mysql -u megamobile_user -p mega_mobile < /opt/megamobile-backend/migrations/20251229_create_service_categories.sql
```

## Step 6: Update Frontend API URL

Update your frontend `.env.production`:

```env
VITE_API_URL=https://api.megamobile.md
```

Or if using same domain:

```env
VITE_API_URL=https://megamobile.md
```

## Maintenance Commands

### View Logs

```bash
docker-compose -f docker-compose.prod.yml logs -f
```

### Restart Container

```bash
docker-compose -f docker-compose.prod.yml restart
```

### Update Application

```bash
cd /opt/megamobile-backend
# Pull latest code or upload new files
docker-compose -f docker-compose.prod.yml up -d --build
```

### Backup Database

```bash
# Create backup
mysqldump -u megamobile_user -p mega_mobile > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore backup
mysql -u megamobile_user -p mega_mobile < backup_20250124_120000.sql
```

### Check Database Connection from Container

```bash
docker exec -it megamobile-api-prod sh
# Inside container, if you have mysql client:
# mysql -h host.docker.internal -u megamobile_user -p mega_mobile
```

## Troubleshooting

### Container Can't Connect to MySQL

1. **Check MySQL is running:**
   ```bash
   sudo systemctl status mysql
   ```

2. **Check MySQL bind address:**
   ```bash
   sudo grep bind-address /etc/mysql/mysql.conf.d/mysqld.cnf
   ```

3. **Test connection from host:**
   ```bash
   mysql -u megamobile_user -p -h localhost mega_mobile
   ```

4. **If host.docker.internal doesn't work, try:**
   - Use `172.17.0.1` (Docker bridge IP)
   - Use the host's actual IP address
   - Use `host-gateway` in docker-compose (already configured)

### Migrations Not Running

- Check container logs: `docker logs megamobile-api-prod`
- Verify database connection in logs
- Check that AutoMigrate is running (it runs automatically in `main.go`)

### CORS Issues

- Verify `CORS_ORIGINS` in `.env.production` includes your frontend domain
- Check backend logs for CORS-related messages
- Verify frontend is using correct API URL

### Port Already in Use

```bash
# Check what's using port 8080
sudo lsof -i :8080

# Or
sudo netstat -tulpn | grep 8080
```

## Security Checklist

- [ ] MySQL root password is set
- [ ] Application user has strong password
- [ ] MySQL only allows necessary connections
- [ ] Firewall configured (UFW recommended)
- [ ] SSL/HTTPS enabled for API
- [ ] `.env.production` has secure passwords
- [ ] Uploads directory has proper permissions
- [ ] Regular database backups configured

## Firewall Setup (UFW)

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow MySQL only from localhost (default)
# No need to open port 3306 if using host.docker.internal

# Enable firewall
sudo ufw enable
sudo ufw status
```

## Production Environment Variables Summary

Required variables in `.env.production`:

- `DB_USER` - MySQL username
- `DB_PASS` - MySQL password  
- `DB_HOST` - MySQL host (use `host.docker.internal` for Docker)
- `DB_PORT` - MySQL port (usually 3306)
- `DB_NAME` - Database name
- `PORT` - Backend port (8080)
- `GIN_MODE` - Set to `release` for production
- `CORS_ORIGINS` - Comma-separated list of allowed frontend origins

Optional:
- `TELEGRAM_BOT_TOKEN` - For Telegram bot functionality
