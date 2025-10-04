# RightScale Input Tracker

A comprehensive tracking and auditing system for RightScale server array inputs. This application monitors and tracks changes to RightScale server array inputs over time, providing historical versioning and comparison capabilities.

## Project Structure

```
.
├── app/                    # Web API server application
│   ├── models/            # Data models and database operations
│   ├── controllers.go     # HTTP request handlers
│   └── routes.go          # API route definitions
├── worker/                # AWS Lambda worker for data fetching
│   ├── fetcher.go         # Lambda function for fetching RightScale data
│   ├── serverless.yml     # Serverless deployment configuration
│   └── Makefile           # Build automation for Lambda
├── frontend/              # Vue.js web interface
│   ├── src/               # Source files
│   ├── build/             # Webpack build configuration
│   └── config/            # Environment configuration
├── build/                 # Docker build configuration
│   ├── Dockerfile         # Container definition for API server
│   └── build              # Build script for Docker
└── main.go               # API server entry point
```

## Components

### 1. API Server (`app/` and `main.go`)
- RESTful API built with Gin framework
- Serves the Vue.js frontend as static files
- Provides endpoints for:
  - Account management
  - Array listing and versioning
  - Input tracking and version history
  - Array comparison

### 2. Lambda Worker (`worker/`)
- AWS Lambda functions that periodically fetch RightScale data
- Deployed using Serverless Framework
- Configured to run every 45 minutes per account
- Stores fetched data in MySQL database

### 3. Frontend (`frontend/`)
- Vue.js single-page application
- Provides UI for:
  - Viewing accounts and arrays
  - Tracking input changes over time
  - Comparing array inputs
  - Visualizing version history

## Prerequisites

- Go 1.x or higher
- Node.js and npm (for frontend)
- MySQL database
- AWS account (for Lambda deployment)
- RightScale API credentials

## Environment Variables

Create a `.env` file with the following variables:

```bash
# Database
MYSQL_CONNECTION_STRING=user:password@tcp(host:3306)/database

# RightScale API
RS_REFRESH_TOKEN=your_refresh_token
RS_ACCOUNT_ID=your_account_id
RS_ACCOUNT_ENDPOINT=https://us-3.rightscale.com
```

## Development Setup

### API Server

```bash
# Build the API server
go build -o bin/input_tracker_app

# Run locally
./bin/input_tracker_app
# Server runs on http://localhost:9080
```

### Lambda Worker

```bash
cd worker
make build
# Deploy with serverless
sls deploy
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## Deployment

### Docker Deployment (API Server)

```bash
cd build
./build
# This builds the Docker image and runs the container
```

### Lambda Deployment

```bash
cd worker
ACCOUNT_ID=your_account_id ./build
# Builds and deploys Lambda functions to AWS
```

## API Endpoints

- `GET /api/accounts` - List all accounts
- `GET /api/accounts/:account_id` - Get account details
- `GET /api/accounts/:account_id/arrays` - List arrays for account
- `GET /api/accounts/:account_id/arrays/:array_id` - Get array details
- `GET /api/accounts/:account_id/arrays/:array_id/inputs` - List current inputs
- `GET /api/accounts/:account_id/arrays/:array_id/inputs/:input_name` - Get input version history
- `GET /api/accounts/:account_id/arrays/:array_id/history` - Get array version history
- `GET /api/accounts/:account_id/array_data/compare/:array_1/:array_2` - Compare arrays
- `POST /api/accounts/:account_id/array_data/compare/:array_1/:array_2` - Compare arrays with ignore list

## Database Schema

The application uses MySQL with the following main tables:
- `accounts` - RightScale account information
- `server_arrays` - Server array records with versioning
- Account-specific tables for input tracking (named by account_id)

## Architecture

1. **Data Collection**: Lambda functions periodically fetch data from RightScale API
2. **Storage**: Data is versioned and stored in MySQL
3. **API Layer**: Gin-based REST API serves data to frontend
4. **Frontend**: Vue.js SPA provides user interface

## Quick Start

1. **Clone and setup**:
   ```bash
   git clone https://github.com/sjeanpierre/rs_input_tracker_go.git
   cd rs_input_tracker_go
   cp .env.example .env
   # Edit .env with your credentials
   ```

2. **Install dependencies**:
   ```bash
   make install-deps
   ```

3. **Run locally**:
   ```bash
   # Terminal 1: API Server
   make run-api
   
   # Terminal 2: Frontend
   make run-frontend
   ```

4. **Access the application**:
   - Frontend: http://localhost:8080
   - API: http://localhost:9080/api

For detailed commands, see [QUICKREF.md](QUICKREF.md).

## Documentation

- **[README.md](README.md)** (this file) - Project overview and quick start
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - System architecture and design
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Development guidelines and best practices
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Deployment guide for all environments
- **[QUICKREF.md](QUICKREF.md)** - Quick reference for common tasks
- **[ORGANIZATIONAL_IMPROVEMENTS.md](ORGANIZATIONAL_IMPROVEMENTS.md)** - Recommended organizational improvements

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Code organization guidelines
- Development workflow
- Code quality standards
- How to submit changes

## License

[Add license information]

## Support

- **Documentation**: See the docs listed above
- **Issues**: Open an issue on GitHub
- **Questions**: Check [QUICKREF.md](QUICKREF.md) for common tasks

## Roadmap

See [ORGANIZATIONAL_IMPROVEMENTS.md](ORGANIZATIONAL_IMPROVEMENTS.md) for planned improvements to project structure and organization.
