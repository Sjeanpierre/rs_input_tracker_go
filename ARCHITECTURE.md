# Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         RightScale API                           │
│              (Multiple accounts: ELSM, GCC, S1NAPROD)            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │ HTTP/HTTPS
                             │
              ┌──────────────▼──────────────┐
              │   AWS Lambda Workers        │
              │  (Serverless Framework)     │
              │                             │
              │  - Runs every 45 minutes    │
              │  - Fetches server arrays    │
              │  - Fetches array inputs     │
              │  - Versions data            │
              └──────────────┬──────────────┘
                             │
                             │ MySQL Protocol
                             │
              ┌──────────────▼──────────────┐
              │        MySQL Database       │
              │                             │
              │  - accounts table           │
              │  - server_arrays table      │
              │  - Account-specific tables  │
              │  - Versioned data storage   │
              └──────────────┬──────────────┘
                             │
                             │ MySQL Protocol
                             │
              ┌──────────────▼──────────────┐
              │      API Server (Go)        │
              │     Gin Framework           │
              │                             │
              │  - REST API endpoints       │
              │  - Serves static frontend   │
              │  - Port 9080                │
              └──────────────┬──────────────┘
                             │
                             │ HTTP/JSON
                             │
              ┌──────────────▼──────────────┐
              │   Frontend (Vue.js SPA)     │
              │                             │
              │  - Account management       │
              │  - Array browsing           │
              │  - Input version history    │
              │  - Array comparison         │
              └─────────────────────────────┘
```

## Component Details

### 1. Data Collection Layer (Lambda Workers)

**Technology**: Go, AWS Lambda, Serverless Framework

**Responsibilities**:
- Authenticate with RightScale API using refresh tokens
- Fetch server array data from multiple RightScale accounts
- Fetch input configurations for each array
- Store data with versioning in MySQL
- Run on scheduled intervals (every 45 minutes)

**Key Files**:
- `worker/fetcher.go` - Main Lambda function
- `worker/serverless.yml` - Deployment configuration
- `worker/Makefile` - Build automation

**Flow**:
1. Lambda triggered by CloudWatch scheduled event
2. Authenticate with RightScale using refresh token
3. Fetch all server arrays for the account
4. For each array, fetch current input configurations
5. Compare with existing data in database
6. Create new version records for changed inputs
7. Mark removed inputs as inactive

### 2. Data Storage Layer (MySQL)

**Technology**: MySQL

**Schema Design**:
- **accounts**: Stores RightScale account information
- **server_arrays**: Stores server array metadata with versioning
- **Account-specific tables**: One table per account for input data (named by account_id)

**Versioning Strategy**:
- Each input change creates a new record with incremented version
- Deleted inputs are marked with `inactive=true` flag
- Queries typically fetch latest version using MAX(version) or MAX(created_at)

### 3. API Layer (Go/Gin)

**Technology**: Go, Gin Framework

**Responsibilities**:
- Expose REST API for frontend consumption
- Serve static Vue.js frontend files
- Handle CORS for development
- Query and aggregate database data

**Key Endpoints**:
```
GET  /api/accounts
GET  /api/accounts/:account_id
GET  /api/accounts/:account_id/arrays
GET  /api/accounts/:account_id/arrays/:array_id
GET  /api/accounts/:account_id/arrays/:array_id/inputs
GET  /api/accounts/:account_id/arrays/:array_id/inputs/:input_name
GET  /api/accounts/:account_id/arrays/:array_id/history
GET  /api/accounts/:account_id/array_data/compare/:array_1/:array_2
POST /api/accounts/:account_id/array_data/compare/:array_1/:array_2
```

**Key Files**:
- `main.go` - Application entry point
- `app/routes.go` - Route definitions
- `app/controllers.go` - Request handlers
- `app/models/` - Data models and database operations

### 4. Presentation Layer (Vue.js)

**Technology**: Vue.js 2.x, Vuetify, Vue Router

**Responsibilities**:
- Display accounts and arrays
- Show input configurations and version history
- Provide array comparison functionality
- Visualize changes over time

**Key Components**:
- `Home.vue` - Account selection
- `ArrayList.vue` - List of arrays for an account
- `ArrayInputs.vue` - Current inputs for an array
- `InputVersions.vue` - Version history for specific input

**Key Files**:
- `frontend/src/main.js` - Application entry
- `frontend/src/router/index.js` - Route definitions
- `frontend/src/components/` - Vue components

## Data Flow

### Input Tracking Flow

1. **Data Collection** (Every 45 minutes)
   ```
   Lambda → RightScale API → Fetch Arrays
                           → Fetch Inputs
                           → Compare with DB
                           → Store Versions
   ```

2. **Version Detection**
   ```
   New Inputs → Compare with Latest Version
             → If Different → Create New Version Record
             → If Removed → Mark as Inactive
   ```

3. **User Query** (When viewing input history)
   ```
   User Request → API Server → Query DB for all versions
                            → Sort by version
                            → Return to Frontend
   ```

### Array Comparison Flow

1. **Simple Comparison** (GET)
   ```
   Frontend → API → Fetch Array 1 Latest Inputs
                 → Fetch Array 2 Latest Inputs
                 → Compare and diff
                 → Return differences
   ```

2. **Comparison with Ignore List** (POST)
   ```
   Frontend → API (with ignore list)
           → Fetch inputs
           → Filter out ignored inputs
           → Compare remaining
           → Return 200 if match, 417 if different
   ```

## Deployment Architecture

### Development
```
Local Machine:
  - Frontend: npm run dev (port 8080)
  - API: go run main.go (port 9080)
  - Database: Local MySQL or remote
```

### Production (API Server)
```
Docker Container:
  - Alpine Linux base
  - Compiled Go binary
  - Exposes port 9080
  - Includes startup script
```

### Production (Lambda Workers)
```
AWS Lambda:
  - Multiple functions (one per account)
  - Scheduled via CloudWatch Events
  - VPC-connected for database access
  - Environment-specific configs
```

## Security Considerations

1. **API Authentication**: Currently no auth (consider adding)
2. **Database Access**: Connection string in environment variables
3. **RightScale Credentials**: Stored as environment variables, injected via serverless
4. **CORS**: Open for development (should restrict in production)
5. **SQL Injection**: Using parameterized queries via GORM

## Scalability Considerations

1. **Database**: 
   - Connection pooling (max 5 connections)
   - Connection lifetime management (10 minutes)
   - Consider read replicas for heavy load

2. **API Server**:
   - Stateless design allows horizontal scaling
   - Can deploy multiple instances behind load balancer

3. **Lambda Workers**:
   - Automatically scales with AWS Lambda
   - Concurrent execution limits apply
   - Consider rate limiting for RightScale API

## Performance Optimizations

1. **Database Queries**:
   - Use of subqueries to fetch latest versions
   - Indexed queries on account_id, array_id, version
   - Raw SQL for complex queries

2. **Lambda Execution**:
   - Parallel array processing with goroutines
   - Connection pooling to avoid thundering herd
   - Batch operations where possible

3. **Frontend**:
   - Code splitting in webpack
   - Production build with minification
   - Vendor chunk separation

## Monitoring and Logging

**Current State**:
- Basic logging with Go's `log` package
- Lambda logs to CloudWatch
- GORM debug mode for SQL logging

**Recommendations**:
- Add structured logging (e.g., logrus, zap)
- Implement metrics collection
- Add application performance monitoring (APM)
- Set up error tracking (e.g., Sentry)
