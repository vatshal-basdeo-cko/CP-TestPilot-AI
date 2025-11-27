# Validation Service

Response validation service for TestPilot AI using C# and .NET 8.

## Features

- JSON schema validation (NJsonSchema)
- Status code validation
- Custom validation rules
- CRUD operations for validation rules
- PostgreSQL integration

## Endpoints

### Validation
- `POST /api/v1/validate` - Validate API response

### Rules Management
- `GET /api/v1/rules` - List all rules
- `GET /api/v1/rules/{id}` - Get rule by ID
- `POST /api/v1/rules` - Create validation rule
- `PUT /api/v1/rules/{id}` - Update rule
- `DELETE /api/v1/rules/{id}` - Delete rule

### Health
- `GET /health` - Health check

## Development

```bash
# Restore packages
dotnet restore

# Run locally
dotnet run

# Build
dotnet build

# Publish
dotnet publish -c Release
```

## Environment Variables

- `ConnectionStrings__DefaultConnection` - PostgreSQL connection string

