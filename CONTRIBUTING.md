# Contributing to RightScale Input Tracker

Thank you for considering contributing to this project! This document provides guidelines and information for contributors.

## Code Organization

### Directory Structure

- **`app/`** - API server code
  - `models/` - Database models and business logic
  - `controllers.go` - HTTP request handlers
  - `routes.go` - Route definitions
  
- **`worker/`** - Lambda worker for data collection
  - Keep Lambda-specific code here
  - Include serverless.yml for deployment config
  
- **`frontend/`** - Vue.js frontend application
  - Follow Vue.js best practices
  - Keep component logic in the component files, not separate script files
  
- **`build/`** - Docker and build configurations

### Go Code Guidelines

1. **Package Structure**
   - Use meaningful package names
   - Keep packages focused on a single responsibility
   - Avoid circular dependencies

2. **Naming Conventions**
   - Use camelCase for variables and functions
   - Use PascalCase for exported types and functions
   - Use descriptive names, avoid abbreviations

3. **Error Handling**
   - Always handle errors explicitly
   - Use meaningful error messages
   - Log errors appropriately

4. **Database Operations**
   - Keep all database logic in the `models` package
   - Use parameterized queries to prevent SQL injection
   - Close database connections properly

### Frontend Guidelines

1. **Vue.js Components**
   - Keep component logic within the component file
   - Use computed properties for derived data
   - Emit events for parent communication

2. **Code Style**
   - Follow Vue.js style guide
   - Use 2 spaces for indentation
   - Use single quotes for strings

3. **File Organization**
   - Component files should be PascalCase
   - Keep components small and focused
   - Extract reusable logic into mixins or composables

## Development Workflow

### Setting Up Development Environment

1. Clone the repository
2. Install Go dependencies
3. Set up environment variables (see README.md)
4. Install frontend dependencies: `cd frontend && npm install`

### Making Changes

1. Create a feature branch from `master`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes
   - Write clean, readable code
   - Add comments for complex logic
   - Follow existing code patterns

3. Test your changes
   - Ensure the API server runs without errors
   - Test frontend functionality
   - Verify database operations

4. Commit your changes
   ```bash
   git add .
   git commit -m "Description of changes"
   ```

5. Push and create a pull request
   ```bash
   git push origin feature/your-feature-name
   ```

### Code Review Process

- All changes require code review
- Address review comments promptly
- Keep pull requests focused and small
- Update documentation as needed

## Testing

### API Testing
- Test all endpoints manually or with tools like Postman
- Verify error handling
- Check database operations

### Frontend Testing
- Test UI functionality in different browsers
- Verify API integration
- Check responsive design

## Common Tasks

### Adding a New API Endpoint

1. Add the handler function in `app/controllers.go`
2. Register the route in `app/routes.go`
3. Add any required model functions in `app/models/`
4. Update API documentation in README.md

### Adding a New Model

1. Create a struct in the appropriate file under `app/models/`
2. Add database operations as methods or functions
3. Use the existing DB connection pattern
4. Export types and functions that need to be accessed from controllers

### Modifying Lambda Worker

1. Update `worker/fetcher.go`
2. Test locally if possible
3. Update `serverless.yml` if deployment config changes
4. Document any new environment variables

### Frontend Changes

1. Modify components in `frontend/src/components/`
2. Update routes if needed in `frontend/src/router/`
3. Test in development mode: `npm run dev`
4. Build and verify: `npm run build`

## Code Quality

### Best Practices

- **Don't Repeat Yourself (DRY)**: Extract common functionality
- **Keep It Simple (KISS)**: Avoid over-engineering
- **Single Responsibility**: Each function/module should do one thing well
- **Separation of Concerns**: Keep business logic separate from presentation

### Performance Considerations

- Use database indexes appropriately
- Avoid N+1 queries
- Cache when appropriate
- Optimize frontend bundle size

## Security

- Never commit sensitive data (API keys, passwords, etc.)
- Use environment variables for configuration
- Validate and sanitize all user inputs
- Use parameterized database queries
- Keep dependencies up to date

## Documentation

- Update README.md for significant changes
- Add inline comments for complex logic
- Document API changes
- Include examples where helpful

## Questions or Issues?

If you have questions or run into issues:
1. Check existing issues on GitHub
2. Review the README.md
3. Create a new issue with detailed information

Thank you for contributing!
