# Execution Service

Execution service for TestPilot AI - executes API calls to target environments.

## Features

- Execute HTTP requests (GET, POST, PUT, PATCH, DELETE)
- Environment management (QA, Staging, Production)
- Request/response logging
- Timeout and retry configuration
- Clean architecture with Go

## Endpoints

### Execution
- `POST /api/v1/execute` - Execute an API call

### Environments
- `GET /api/v1/environments` - List all environments
- `GET /api/v1/environments/:id` - Get environment by ID
- `POST /api/v1/environments` - Create environment
- `PUT /api/v1/environments/:id` - Update environment
- `DELETE /api/v1/environments/:id` - Delete environment

### Health
- `GET /health` - Health check

## Development

```bash
# Run locally
go run main.go

# Build
go build -o execution

# Run tests
go test ./...
```

## Environment Variables

- `SERVER_PORT` - Server port (default: 8003)
- `DATABASE_URL` - PostgreSQL connection string
- `DEFAULT_TIMEOUT` - Default request timeout in seconds (default: 30)
- `MAX_RETRIES` - Maximum retry attempts (default: 3)

