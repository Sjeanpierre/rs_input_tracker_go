# Organizational Structure Recommendations

This document outlines recommended improvements to the project's organizational structure.

## Current Issues

### 1. Package Structure
- **Issue**: The `worker/` directory contains `package main`, which is confusing since it's not at the root level
- **Impact**: Makes it unclear that this is a separate executable/Lambda function
- **Recommendation**: Consider renaming to `cmd/worker/` or `lambda/` to better indicate it's a separate binary

### 2. Frontend Component Organization
- **Issue**: JavaScript logic is separated from Vue components in a `script/` subdirectory
- **Files affected**: 
  - `frontend/src/components/script/array_inputs.js`
  - `frontend/src/components/script/array_list.js`
  - `frontend/src/components/script/input_versions.js`
  - `frontend/src/components/script/home.js`
- **Impact**: Reduces code cohesion and makes components harder to maintain
- **Recommendation**: Move logic back into the `.vue` component files in the `<script>` section

### 3. Build Scripts Location
- **Issue**: Build scripts scattered across multiple locations:
  - `/build/build` - Docker build script
  - `/worker/build` - Lambda build script
  - `/worker/Makefile` - Lambda make targets
  - `/build/Dockerfile` - Docker configuration
- **Impact**: Inconsistent build process, unclear what builds what
- **Recommendation**: 
  - Create a root-level `Makefile` or `build/` directory with all build scripts
  - Clearly name scripts by purpose (e.g., `build-api.sh`, `build-lambda.sh`, `build-docker.sh`)

### 4. Missing Go Modules
- **Issue**: No `go.mod` file for dependency management
- **Impact**: 
  - Unclear what dependencies are required
  - Difficult to reproduce builds
  - No version pinning
- **Recommendation**: Initialize Go modules with `go mod init`

### 5. Environment Configuration
- **Issue**: Multiple `.env` symlinks pointing to user-specific paths
  - `/worker/.env` -> `/Users/jstevenson/.secure/sls_fetcher_creds`
  - Root `.env` -> `/Users/jstevenson/.secure/input_tracker_app_env`
- **Impact**: Won't work for other developers
- **Recommendation**: 
  - Create `.env.example` files with dummy values
  - Add `.env` to `.gitignore`
  - Document required environment variables in README

### 6. Inconsistent Naming Conventions
- **Issue**: Mix of naming styles:
  - Go files: mix of snake_case and camelCase
  - Directories: lowercase
  - Frontend: various styles
- **Examples**:
  - `input_tracker_go` (snake_case)
  - `arrayInputs` (camelCase)
  - `rs_input_tracker_go` (snake_case with prefix)
- **Recommendation**: Establish and document naming conventions:
  - Go: use camelCase for variables, PascalCase for exports
  - Files: use snake_case or kebab-case consistently
  - Directories: use lowercase with underscores or hyphens

### 7. Missing Documentation
- **Issue**: No main README.md at project root (now added)
- **Impact**: New developers don't know how to set up or use the project
- **Recommendation**: Maintain comprehensive README with:
  - Project overview
  - Setup instructions
  - Architecture diagram
  - API documentation
  - Deployment guides

## Recommended Directory Structure

Here's an improved directory structure to consider:

```
rs_input_tracker_go/
├── README.md                    # Project overview and setup (✓ Added)
├── CONTRIBUTING.md              # Contribution guidelines (✓ Added)
├── .env.example                # Example environment configuration
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies lock file
├── Makefile                    # Root-level build commands
│
├── cmd/                        # Main applications
│   ├── server/                # API server
│   │   └── main.go
│   └── worker/                # Lambda worker
│       └── main.go
│
├── internal/                   # Private application code
│   ├── models/                # Data models
│   ├── handlers/              # HTTP handlers (controllers)
│   └── routes/                # Route definitions
│
├── pkg/                       # Public library code (if any)
│
├── web/                       # Frontend application
│   ├── src/
│   ├── build/
│   └── config/
│
├── deployments/               # Deployment configurations
│   ├── docker/
│   │   └── Dockerfile
│   └── lambda/
│       └── serverless.yml
│
└── scripts/                   # Build and deployment scripts
    ├── build-server.sh
    ├── build-worker.sh
    └── deploy-lambda.sh
```

## Priority Recommendations

### High Priority (Do First)
1. ✅ Add comprehensive README.md
2. ✅ Add CONTRIBUTING.md
3. 🔲 Initialize Go modules (`go mod init`)
4. 🔲 Create `.env.example` files
5. 🔲 Move frontend component logic into .vue files

### Medium Priority (Do Soon)
6. 🔲 Reorganize build scripts into `scripts/` directory
7. 🔲 Restructure packages following Go conventions (`cmd/`, `internal/`)
8. 🔲 Document all API endpoints (OpenAPI/Swagger)
9. 🔲 Add architecture diagram to README

### Low Priority (Nice to Have)
10. 🔲 Add unit tests with coverage reporting
11. 🔲 Set up CI/CD pipeline
12. 🔲 Add linting configuration (golangci-lint, ESLint)
13. 🔲 Create developer setup script

## Implementation Steps

### Step 1: Documentation (✓ Complete)
- [x] Create README.md with project overview
- [x] Create CONTRIBUTING.md with guidelines
- [ ] Add .env.example files

### Step 2: Go Modules
```bash
go mod init github.com/sjeanpierre/rs_input_tracker_go
go mod tidy
```

### Step 3: Frontend Reorganization
For each component, move the script logic into the component:
```javascript
// Instead of importing from script/
export default {
  // ... component definition here
}
```

### Step 4: Build Scripts
Consolidate into a `scripts/` directory with clear naming:
- `scripts/build-api.sh` - Build API server
- `scripts/build-lambda.sh` - Build Lambda functions
- `scripts/build-docker.sh` - Build Docker image
- `scripts/deploy-lambda.sh` - Deploy to AWS

### Step 5: Package Restructuring
Move files to follow Go standard project layout:
- `app/` → `cmd/server/` and `internal/`
- `worker/` → `cmd/worker/`
- Keep `frontend/` or rename to `web/`

## Benefits of Reorganization

1. **Clarity**: Clear separation between different components
2. **Maintainability**: Easier to find and modify code
3. **Onboarding**: New developers can understand structure quickly
4. **Best Practices**: Follows Go and Vue.js community standards
5. **Scalability**: Easier to add new features and components

## Migration Strategy

To avoid breaking changes:

1. Create new structure alongside existing
2. Gradually move code with full test coverage
3. Update build scripts incrementally
4. Maintain both structures temporarily
5. Deprecate old structure once stable
6. Remove old structure after verification

## Additional Resources

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Vue.js Style Guide](https://vuejs.org/style-guide/)
- [The Twelve-Factor App](https://12factor.net/)
- [Semantic Versioning](https://semver.org/)
