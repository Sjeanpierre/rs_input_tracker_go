# Deployment Guide

This guide covers deployment options for the RightScale Input Tracker application.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [AWS Lambda Deployment](#aws-lambda-deployment)
- [Production Deployment](#production-deployment)
- [Environment Configuration](#environment-configuration)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### All Deployments
- MySQL database (accessible from deployment environment)
- RightScale API credentials (refresh tokens and account IDs)

### Docker Deployment
- Docker installed and running
- Access to Docker registry (for production)

### AWS Lambda Deployment
- AWS account with appropriate permissions
- AWS CLI configured
- Serverless Framework installed: `npm install -g serverless`
- VPC configured with MySQL access

### Production Deployment
- Domain name (optional)
- SSL certificate (recommended)
- Load balancer or reverse proxy (recommended)

## Local Development

### 1. Setup Environment

```bash
# Clone repository
git clone https://github.com/sjeanpierre/rs_input_tracker_go.git
cd rs_input_tracker_go

# Copy environment files
cp .env.example .env
cp worker/.env.example worker/.env

# Edit with your credentials
nano .env
nano worker/.env
```

### 2. Install Dependencies

```bash
# Using Makefile
make install-deps

# Or manually
go mod download
cd frontend && npm install
```

### 3. Run Components

**Option A: Using Makefile**
```bash
# Run API server
make run-api

# Run frontend (in another terminal)
make run-frontend
```

**Option B: Manual Commands**
```bash
# Terminal 1: API Server
go run main.go

# Terminal 2: Frontend
cd frontend && npm run dev
```

### 4. Access Application
- Frontend: http://localhost:8080
- API: http://localhost:9080/api

## Docker Deployment

### Building Docker Image

**Using Makefile:**
```bash
make build-docker
```

**Using Docker directly:**
```bash
cd build
docker build -t input_tracker .
```

**Using build script:**
```bash
cd build
./build
```

### Running Docker Container

**Basic run:**
```bash
docker run -p 9080:9080 --rm -it --env-file .env input_tracker
```

**With custom environment:**
```bash
docker run -p 9080:9080 --rm -it \
  -e MYSQL_CONNECTION_STRING="user:pass@tcp(host:3306)/db" \
  -e RS_REFRESH_TOKEN="your_token" \
  -e RS_ACCOUNT_ID="your_account_id" \
  -e RS_ACCOUNT_ENDPOINT="https://us-3.rightscale.com" \
  input_tracker
```

**Detached mode (background):**
```bash
docker run -d -p 9080:9080 --name input_tracker --env-file .env input_tracker
```

**View logs:**
```bash
docker logs -f input_tracker
```

**Stop container:**
```bash
docker stop input_tracker
docker rm input_tracker
```

### Docker Compose (Recommended)

Create `docker-compose.yml` in project root:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: build/Dockerfile
    ports:
      - "9080:9080"
    environment:
      - MYSQL_CONNECTION_STRING=${MYSQL_CONNECTION_STRING}
      - RS_REFRESH_TOKEN=${RS_REFRESH_TOKEN}
      - RS_ACCOUNT_ID=${RS_ACCOUNT_ID}
      - RS_ACCOUNT_ENDPOINT=${RS_ACCOUNT_ENDPOINT}
    restart: unless-stopped
    depends_on:
      - db
  
  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=rs_input_tracker
      - MYSQL_USER=tracker
      - MYSQL_PASSWORD=trackerpass
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

Run with Docker Compose:
```bash
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## AWS Lambda Deployment

### 1. Configure AWS

```bash
# Configure AWS CLI
aws configure
# Enter Access Key ID, Secret Key, Region, Output format

# Verify configuration
aws sts get-caller-identity
```

### 2. Configure VPC (One-time Setup)

Ensure your Lambda functions can access:
- MySQL database (via VPC connection)
- Internet (for RightScale API access - via NAT Gateway)

Update `worker/serverless.yml`:
```yaml
vpc:
  securityGroupIds:
    - sg-xxxxxxxxx  # Security group with MySQL access
  subnetIds:
    - subnet-xxxxxxxx  # Private subnet with NAT gateway
    - subnet-yyyyyyyy  # Another subnet for HA
```

### 3. Configure Environment Variables

Edit `worker/.env`:
```bash
# Database
LIVE_MYSQL_CONNECTION_STRING=user:pass@tcp(db-host:3306)/database

# RightScale credentials for each account
RS_ELSM_REFRESH_TOKEN=token_for_elsm
RS_ELSM_ACCOUNT_ID=account_id_for_elsm

RS_GCC_REFRESH_TOKEN=token_for_gcc
RS_GCC_ACCOUNT_ID=account_id_for_gcc

RS_S1NAPROD_REFRESH_TOKEN=token_for_s1naprod
RS_S1NAPROD_ACCOUNT_ID=account_id_for_s1naprod

# AWS account for deployment
ACCOUNT_ID=123456789012
```

### 4. Deploy Lambda Functions

**Using Makefile:**
```bash
ACCOUNT_ID=123456789012 make deploy-worker
```

**Using build script:**
```bash
cd worker
ACCOUNT_ID=123456789012 ./build
```

**Manual deployment:**
```bash
cd worker
make build
sls deploy
```

### 5. Verify Deployment

```bash
# Check deployed functions
sls info

# View logs
sls logs -f fetch_elsm_acct --tail
sls logs -f fetch_gcc_acct --tail
sls logs -f fetch_s1naprod_acct --tail

# Invoke function manually
sls invoke -f fetch_elsm_acct
```

### 6. Monitor Execution

```bash
# CloudWatch Logs
aws logs tail /aws/lambda/rs-input-fetcher-global-fetch_elsm_acct --follow

# View metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda \
  --metric-name Invocations \
  --dimensions Name=FunctionName,Value=rs-input-fetcher-global-fetch_elsm_acct \
  --start-time 2024-01-01T00:00:00Z \
  --end-time 2024-01-02T00:00:00Z \
  --period 3600 \
  --statistics Sum
```

## Production Deployment

### Architecture Options

**Option 1: Docker on EC2/ECS**
- Deploy API container to EC2 or ECS
- Use Application Load Balancer
- Auto Scaling Group for high availability

**Option 2: Kubernetes**
- Deploy using Kubernetes manifests
- Use Ingress for routing
- Horizontal Pod Autoscaler

**Option 3: Cloud Run / App Engine**
- Deploy containerized app to managed service
- Automatic scaling
- Built-in load balancing

### Recommended Production Setup

```
┌─────────────┐
│   Route 53  │ (DNS)
└──────┬──────┘
       │
┌──────▼──────────┐
│ Application LB  │ (SSL termination)
└──────┬──────────┘
       │
┌──────▼──────────┐
│   Auto Scaling  │
│      Group      │
│                 │
│ ┌─────────────┐ │
│ │  EC2 + App  │ │ (Multiple instances)
│ └─────────────┘ │
└─────────────────┘
       │
┌──────▼──────────┐
│   RDS MySQL     │ (Managed database)
└─────────────────┘
```

### 1. Build Production Frontend

```bash
cd frontend

# Update API URL in config
# Edit frontend/config/prod.env.js

# Build
npm run build

# Output in frontend/dist/
```

### 2. Build Production API Binary

```bash
# Build optimized binary
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/input_tracker_app

# Or use make
make build-api
```

### 3. Prepare Production Image

```bash
# Build and tag
docker build -t input_tracker:1.0.0 -f build/Dockerfile .

# Tag for registry
docker tag input_tracker:1.0.0 your-registry/input_tracker:1.0.0

# Push to registry
docker push your-registry/input_tracker:1.0.0
```

### 4. Deploy to Production

**EC2 Example:**
```bash
# SSH to EC2 instance
ssh ec2-user@your-server

# Pull image
docker pull your-registry/input_tracker:1.0.0

# Stop old container
docker stop input_tracker
docker rm input_tracker

# Start new container
docker run -d \
  --name input_tracker \
  -p 9080:9080 \
  --restart unless-stopped \
  --env-file /etc/input_tracker/.env \
  your-registry/input_tracker:1.0.0

# Check logs
docker logs -f input_tracker
```

**ECS Example:**
1. Create ECS task definition with image
2. Create ECS service with ALB
3. Deploy new task revision
4. Monitor deployment

### 5. Configure Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:9080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Environment Configuration

### Development (.env)
```bash
MYSQL_CONNECTION_STRING=user:pass@tcp(localhost:3306)/rs_input_tracker_dev
RS_REFRESH_TOKEN=dev_token
RS_ACCOUNT_ID=dev_account
RS_ACCOUNT_ENDPOINT=https://us-3.rightscale.com
```

### Staging (.env)
```bash
MYSQL_CONNECTION_STRING=user:pass@tcp(staging-db:3306)/rs_input_tracker_staging
RS_REFRESH_TOKEN=staging_token
RS_ACCOUNT_ID=staging_account
RS_ACCOUNT_ENDPOINT=https://us-3.rightscale.com
```

### Production (.env)
```bash
MYSQL_CONNECTION_STRING=user:pass@tcp(prod-db:3306)/rs_input_tracker
RS_REFRESH_TOKEN=prod_token
RS_ACCOUNT_ID=prod_account
RS_ACCOUNT_ENDPOINT=https://us-3.rightscale.com
```

## Health Checks

Add health check endpoint to `app/routes.go`:
```go
api.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
})
```

Configure load balancer health check:
- Path: `/api/health`
- Interval: 30 seconds
- Timeout: 5 seconds
- Healthy threshold: 2
- Unhealthy threshold: 3

## Monitoring

### Application Logs
```bash
# Docker
docker logs -f input_tracker

# Lambda
sls logs -f fetch_elsm_acct --tail

# Systemd (if using)
journalctl -u input_tracker -f
```

### Metrics to Monitor
- Request rate and latency
- Error rate
- Lambda execution time and errors
- Database connection pool usage
- Memory and CPU usage

### Recommended Tools
- **Logs**: CloudWatch, ELK Stack, Splunk
- **Metrics**: Prometheus + Grafana, DataDog
- **APM**: New Relic, DataDog APM
- **Errors**: Sentry, Rollbar

## Backup and Recovery

### Database Backups
```bash
# Manual backup
mysqldump -h db-host -u user -p rs_input_tracker > backup.sql

# Restore
mysql -h db-host -u user -p rs_input_tracker < backup.sql
```

### Automated Backups
- Use RDS automated backups (AWS)
- Set up cron job for mysqldump
- Store backups in S3 or similar

## Troubleshooting

### Docker Issues
```bash
# Check container status
docker ps -a

# View logs
docker logs input_tracker

# Execute commands in container
docker exec -it input_tracker sh

# Check resource usage
docker stats input_tracker
```

### Lambda Issues
```bash
# View recent logs
sls logs -f fetch_elsm_acct --startTime 1h

# Check function configuration
aws lambda get-function-configuration --function-name rs-input-fetcher-global-fetch_elsm_acct

# Test function
sls invoke -f fetch_elsm_acct -l
```

### Database Connection Issues
```bash
# Test connection from container
docker exec -it input_tracker sh
mysql -h db-host -u user -p

# Test from Lambda (using VPC)
# Ensure security groups allow MySQL (port 3306)
# Ensure subnets have route to database
```

## Security Checklist

- [ ] Use HTTPS/TLS in production
- [ ] Store credentials in secrets manager (AWS Secrets Manager, HashiCorp Vault)
- [ ] Rotate RightScale tokens regularly
- [ ] Enable database encryption at rest
- [ ] Use VPC for Lambda to database communication
- [ ] Implement API authentication/authorization
- [ ] Enable CloudTrail for AWS audit logs
- [ ] Regular security updates for dependencies
- [ ] Implement rate limiting
- [ ] Enable WAF for DDoS protection

## Rollback Procedure

### Docker Deployment
```bash
# Tag previous working image as :stable
docker tag input_tracker:1.0.0 input_tracker:stable

# On issues, rollback
docker stop input_tracker
docker rm input_tracker
docker run -d --name input_tracker ... input_tracker:stable
```

### Lambda Deployment
```bash
# Deploy specific version
sls deploy --stage global --version 1.0.0

# Or rollback to previous
aws lambda update-function-code \
  --function-name rs-input-fetcher-global-fetch_elsm_acct \
  --s3-bucket your-bucket \
  --s3-key previous-version.zip
```

## Performance Optimization

### Database
- Enable query cache
- Add indexes on frequently queried columns
- Use read replicas for read-heavy workloads

### API Server
- Enable response compression
- Implement caching (Redis)
- Use connection pooling (already implemented)

### Lambda
- Increase memory for faster execution
- Use provisioned concurrency for consistent performance
- Optimize cold starts

## Support

For deployment issues:
1. Check application logs
2. Review CloudWatch metrics
3. Verify environment configuration
4. Check database connectivity
5. Review this guide and QUICKREF.md
6. Open issue on GitHub
