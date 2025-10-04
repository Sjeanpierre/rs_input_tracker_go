# Quick Reference Guide

This guide provides quick commands and references for common development tasks.

## Initial Setup

### 1. Clone and Setup
```bash
git clone https://github.com/sjeanpierre/rs_input_tracker_go.git
cd rs_input_tracker_go
```

### 2. Environment Configuration
```bash
# Copy example env files
cp .env.example .env
cp worker/.env.example worker/.env

# Edit .env files with your credentials
nano .env
nano worker/.env
```

### 3. Install Dependencies
```bash
# Go dependencies (automatically downloaded when building)
go mod download

# Frontend dependencies
cd frontend
npm install
cd ..
```

## Development Commands

### API Server

```bash
# Run API server locally
go run main.go

# Build API server binary
go build -o bin/input_tracker_app

# Run with hot reload (using air or similar)
air
```

### Frontend

```bash
cd frontend

# Development mode with hot reload
npm run dev

# Build for production
npm run build

# Build and analyze bundle
npm run build --report
```

### Lambda Worker

```bash
cd worker

# Build Lambda function
make build

# Deploy to AWS (requires AWS credentials and ACCOUNT_ID env var)
ACCOUNT_ID=your_account_id sls deploy

# View logs
sls logs -f fetch_elsm_acct --tail

# Remove deployment
sls remove
```

## Docker Commands

### API Server

```bash
cd build

# Build Docker image
docker build -t input_tracker .

# Run Docker container
docker run -p 9080:9080 --rm -it --env-file ../.env input_tracker

# Or use the build script
./build
```

## Database Operations

### Connect to MySQL
```bash
mysql -h <host> -u <user> -p <database>
```

### Common Queries
```sql
-- List all accounts
SELECT * FROM accounts;

-- List all arrays for an account
SELECT * FROM server_arrays WHERE account_id = 123;

-- List current arrays (latest version)
SELECT * FROM server_arrays sa
WHERE created_at = (
    SELECT MAX(created_at) 
    FROM server_arrays 
    WHERE array_id = sa.array_id
);

-- List inputs for an array (from account-specific table)
SELECT * FROM `123` 
WHERE array_id = 456 
  AND account_id = 123;

-- Get input version history
SELECT * FROM `123`
WHERE array_id = 456 
  AND input_name = 'MY_INPUT'
ORDER BY version DESC;
```

## Testing

### Manual API Testing

```bash
# Test with curl
curl http://localhost:9080/api/accounts

# Get specific account
curl http://localhost:9080/api/accounts/123

# List arrays
curl http://localhost:9080/api/accounts/123/arrays

# Get array inputs
curl http://localhost:9080/api/accounts/123/arrays/456/inputs

# Compare arrays
curl http://localhost:9080/api/accounts/123/array_data/compare/456/789

# Compare with ignore list
curl -X POST http://localhost:9080/api/accounts/123/array_data/compare/456/789 \
  -H "Content-Type: application/json" \
  -d '{"ignored": ["INPUT_TO_IGNORE"]}'
```

### Using Postman/Insomnia
Import these endpoints:
- Base URL: `http://localhost:9080/api`
- All GET endpoints don't require body
- POST compare endpoint requires: `{"ignored": ["INPUT1", "INPUT2"]}`

## Troubleshooting

### API Server Issues

**Problem**: Server won't start
```bash
# Check if port is in use
lsof -i :9080

# Check environment variables
printenv | grep -E 'MYSQL|RS_'

# Check MySQL connection
mysql -h <host> -u <user> -p
```

**Problem**: Database connection errors
- Verify MYSQL_CONNECTION_STRING format: `user:password@tcp(host:port)/database`
- Check MySQL is running and accessible
- Verify credentials are correct

### Lambda Issues

**Problem**: Lambda timeout
- Increase timeout in serverless.yml (current: 300s)
- Check MySQL connection limits
- Optimize parallel processing

**Problem**: Lambda can't connect to database
- Verify VPC configuration in serverless.yml
- Check security groups allow MySQL traffic
- Ensure subnet has NAT gateway for internet access

**Problem**: Deployment fails
- Verify AWS credentials: `aws sts get-caller-identity`
- Check ACCOUNT_ID environment variable
- Ensure serverless framework is installed: `npm install -g serverless`

### Frontend Issues

**Problem**: Build fails
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

**Problem**: API calls fail from frontend
- Check CORS settings in main.go
- Verify API_URL in frontend config
- Check browser console for errors

**Problem**: Development server won't start
```bash
# Check if port 8080 is in use
lsof -i :8080

# Use different port
npm run dev -- --port 8081
```

## Code Quality

### Linting and Formatting

```bash
# Go formatting
go fmt ./...

# Go linting (requires golangci-lint)
golangci-lint run

# Frontend linting
cd frontend
npm run lint
```

### Building and Testing

```bash
# Build Go code
go build ./...

# Run Go tests (if any)
go test ./...

# Build frontend
cd frontend && npm run build

# Frontend tests (if any)
npm run test
```

## Useful Git Commands

```bash
# Check status
git status

# Create feature branch
git checkout -b feature/my-feature

# Commit changes
git add .
git commit -m "Description of changes"

# Push changes
git push origin feature/my-feature

# Pull latest changes
git pull origin master

# View commit history
git log --oneline --graph
```

## Environment Variables Reference

### API Server (.env)
| Variable | Description | Example |
|----------|-------------|---------|
| MYSQL_CONNECTION_STRING | MySQL connection string | `user:pass@tcp(host:3306)/db` |
| RS_REFRESH_TOKEN | RightScale refresh token | `<token>` |
| RS_ACCOUNT_ID | RightScale account ID | `12345` |
| RS_ACCOUNT_ENDPOINT | RightScale API endpoint | `https://us-3.rightscale.com` |

### Lambda Worker (worker/.env)
| Variable | Description |
|----------|-------------|
| LIVE_MYSQL_CONNECTION_STRING | Production MySQL connection |
| RS_ELSM_REFRESH_TOKEN | ELSM account token |
| RS_ELSM_ACCOUNT_ID | ELSM account ID |
| RS_GCC_REFRESH_TOKEN | GCC account token |
| RS_GCC_ACCOUNT_ID | GCC account ID |
| RS_S1NAPROD_REFRESH_TOKEN | S1NA Production token |
| RS_S1NAPROD_ACCOUNT_ID | S1NA Production account ID |
| ACCOUNT_ID | AWS account ID for deployment |

## Common Development Workflows

### Adding a New API Endpoint

1. Add handler function in `app/controllers.go`:
```go
func myNewEndpoint(c *gin.Context) {
    // implementation
    c.JSON(200, result)
}
```

2. Register route in `app/routes.go`:
```go
api.GET("/my-endpoint", myNewEndpoint)
```

3. Add model functions if needed in `app/models/`

4. Test the endpoint:
```bash
curl http://localhost:9080/api/my-endpoint
```

### Adding a New Frontend Page

1. Create component in `frontend/src/components/MyPage.vue`
2. Add route in `frontend/src/router/index.js`
3. Test in development mode: `npm run dev`
4. Build for production: `npm run build`

### Debugging Lambda Issues

1. Check CloudWatch logs:
```bash
cd worker
sls logs -f fetch_elsm_acct --tail
```

2. Test locally (if possible):
```bash
# Set environment variables
export RS_REFRESH_TOKEN="..."
export RS_ACCOUNT_ID="..."
export MYSQL_CONNECTION_STRING="..."

# Run the code
go run fetcher.go
```

3. Deploy with verbose logging:
```bash
sls deploy --verbose
```

## Performance Tips

1. **Database**:
   - Create indexes on frequently queried columns
   - Use connection pooling (already configured)
   - Monitor slow queries

2. **API**:
   - Cache frequently accessed data
   - Use pagination for large result sets
   - Minimize database round trips

3. **Frontend**:
   - Lazy load routes
   - Optimize bundle size
   - Use production build for deployment

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Vue.js Guide](https://vuejs.org/guide/)
- [Serverless Framework](https://www.serverless.com/framework/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [RightScale API](https://docs.rightscale.com/api/)
